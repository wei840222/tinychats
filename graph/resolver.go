package graph

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/AlekSi/pointer"
	"github.com/gorilla/websocket"
	"github.com/wei840222/tinychats/graph/generated"
	"github.com/wei840222/tinychats/graph/model"
	"github.com/wei840222/tinychats/pkg"
	"github.com/wei840222/tinychats/proto"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	sync.Mutex

	messageCreatedChans map[string]chan *model.Message

	userClient    proto.UserClient
	messageClient proto.MessageClient
}

func (r *Resolver) dataloadersFormContext(ctx context.Context) *dataloaders {
	return ctx.Value(dataloadersCtxKey).(*dataloaders)
}

func NewGraphQLHandler(userClient proto.UserClient, messageClient proto.MessageClient) http.Handler {
	resolvers := &Resolver{
		messageCreatedChans: make(map[string]chan *model.Message),
		userClient:          userClient,
		messageClient:       messageClient,
	}
	gqlHandler := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers}))
	gqlHandler.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		KeepAlivePingInterval: 10 * time.Second,
	})
	gqlHandler.AddTransport(transport.Options{})
	gqlHandler.AddTransport(transport.GET{})
	gqlHandler.AddTransport(transport.POST{})
	gqlHandler.AddTransport(transport.MultipartForm{})
	gqlHandler.SetQueryCache(lru.New(1000))
	gqlHandler.Use(extension.Introspection{})
	gqlHandler.Use(extension.AutomaticPersistedQuery{Cache: lru.New(100)})
	return newGraphQDataloadersMiddleware(userClient)(gqlHandler)
}

var dataloadersCtxKey = struct{}{}

type dataloaders struct {
	userLoader *UserLoader
}

func newGraphQDataloadersMiddleware(userClient proto.UserClient) pkg.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), dataloadersCtxKey, &dataloaders{
				userLoader: NewUserLoader(UserLoaderConfig{
					Fetch: func(ids []string) ([]*model.User, []error) {
						res, err := userClient.MutiGetUsers(r.Context(), &proto.MutiGetUsersRequest{Ids: ids})
						if err != nil {
							return nil, []error{err}
						}
						userResponseMap := make(map[string]*proto.UserResponse)
						for _, userResponse := range res.GetUsers() {
							userResponseMap[userResponse.GetId()] = userResponse
						}
						var users []*model.User
						var errors []error
						for _, id := range ids {
							userResponse, ok := userResponseMap[id]
							if !ok {
								errors = append(errors, fmt.Errorf("user: %s not fount", id))
								continue
							}
							users = append(users, &model.User{
								ID:        userResponse.GetId(),
								Name:      userResponse.GetName(),
								AvatarURL: pointer.ToStringOrNil(userResponse.GetAvatarUrl()),
							})
						}
						return users, errors
					},
					Wait:     2 * time.Millisecond,
					MaxBatch: 100,
				}),
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
