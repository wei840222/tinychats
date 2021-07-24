package graph

import (
	"context"
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

	userLoader *UserLoader
}

func NewGraphQLHandler(userClient proto.UserClient, messageClient proto.MessageClient) http.Handler {
	resolvers := &Resolver{
		messageCreatedChans: make(map[string]chan *model.Message),
		userClient:          userClient,
		messageClient:       messageClient,
		userLoader: NewUserLoader(UserLoaderConfig{
			Fetch: func(ids []string) ([]*model.User, []error) {
				res, err := userClient.MutiGetUsers(context.Background(), &proto.MutiGetUsersRequest{Ids: ids})
				if err != nil {
					return nil, []error{err}
				}
				var users []*model.User
				for _, userResponse := range res.GetUsers() {
					users = append(users, &model.User{
						ID:        userResponse.GetId(),
						Name:      userResponse.GetName(),
						AvatarURL: pointer.ToStringOrNil(userResponse.GetAvatarUrl()),
					})
				}
				return users, nil
			},
			Wait:     2 * time.Millisecond,
			MaxBatch: 100,
		}),
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
	return gqlHandler
}
