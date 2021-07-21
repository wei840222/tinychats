package main

import (
	"github.com/wei840222/todo/graph"
	"github.com/wei840222/todo/graph/generated"
	"github.com/wei840222/todo/pkg/todo"
	"github.com/wei840222/todo/proto"

	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func main() {
	log := zapLogger()

	host := flag.String("host", "", "server listening host")
	port := flag.Int("port", 8080, "server listening port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcC, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}

	m := http.NewServeMux()
	m.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))
	m.Handle("/graphql", handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					Logger:     log,
					TodoClient: proto.NewTodoClient(grpcC),
				},
			}),
	))
	httpS := &http.Server{Handler: newAccessLogMiddleware(newRecoveryMiddleware(m))}

	// Shared options for the logger, with a custom gRPC code to log level function.
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(err interface{}) error {
			log.Error(fmt.Sprint(err))
			return status.Errorf(codes.Unknown, "panic triggered: %v", err)
		}),
	}

	grpcS := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_zap.UnaryServerInterceptor(log),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_zap.StreamServerInterceptor(log),
			grpc_recovery.StreamServerInterceptor(recoveryOpts...),
		),
	)
	proto.RegisterTodoServer(grpcS, &todo.Server{})
	reflection.Register(grpcS)

	cm := cmux.New(lis)

	httpL := cm.Match(cmux.HTTP1Fast())
	grpcL := cm.Match(cmux.Any())
	go httpS.Serve(httpL)
	go grpcS.Serve(grpcL)

	if err := cm.Serve(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err.Error())
	}
}
