package main

import (
	"github.com/wei840222/todo/graph"
	"github.com/wei840222/todo/graph/generated"
	"github.com/wei840222/todo/pkg/todo"
	"github.com/wei840222/todo/proto"

	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/lib/pq"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed web
var web embed.FS

func main() {
	log := zapLogger()

	var port int
	var databaseURL string
	flag.IntVar(&port, "port", 8080, "server listening port")
	flag.StringVar(&databaseURL, "databaseUrl", "", "database url (ex: postgres://username:password@localhost:5432/todo)")
	flag.Parse()

	if envPort := os.Getenv("PORT"); port == 8080 && envPort != "" {
		port, _ = strconv.Atoi(envPort)
	}
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcC, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}

	public, err := fs.Sub(web, "web")
	if err != nil {
		log.Fatal(err.Error())
	}
	m := http.NewServeMux()
	m.Handle("/", http.FileServer(http.FS(public)))
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

	var db *gorm.DB
	if databaseURL == "" {
		sqliteDB, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
		if err != nil {
			log.Fatal(err.Error())
		}
		db = sqliteDB
	} else {
		dsn, err := pq.ParseURL(databaseURL)
		if err != nil {
			log.Fatal(err.Error())
		}
		postgresDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal(err.Error())
		}
		db = postgresDB
	}

	proto.RegisterTodoServer(grpcS, &todo.Server{DB: db})
	reflection.Register(grpcS)

	cm := cmux.New(lis)

	httpL := cm.Match(cmux.HTTP1Fast())
	grpcL := cm.Match(cmux.Any())
	go httpS.Serve(httpL)
	go grpcS.Serve(grpcL)

	log.Sugar().Infof("server start and listen on %s", lis.Addr().String())
	if err := cm.Serve(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err.Error())
	}
}
