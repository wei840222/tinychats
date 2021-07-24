package main

import (
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
	line_login_sdk "github.com/kkdai/line-login-sdk-go"
	"github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/soheilhy/cmux"
	"github.com/wei840222/tinychats/graph"
	"github.com/wei840222/tinychats/graph/generated"
	"github.com/wei840222/tinychats/graph/model"
	"github.com/wei840222/tinychats/pkg"
	"github.com/wei840222/tinychats/pkg/message"
	"github.com/wei840222/tinychats/pkg/user"
	"github.com/wei840222/tinychats/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

//go:embed web/dist
var webDist embed.FS

func main() {
	log := pkg.ZapLogger()

	var port int
	var databaseURL, lineChannelID, lineChannelSecret string
	flag.IntVar(&port, "port", 8080, "server listening port")
	flag.StringVar(&databaseURL, "database.url", "", "database url (ex: postgres://username:password@localhost:5432/todo)")
	flag.StringVar(&lineChannelID, "line.channel.id", "", "LINE Login Channel Id")
	flag.StringVar(&lineChannelSecret, "line.channel.secret", "", "LINE Login Channel Secret")
	flag.Parse()

	if envPort := os.Getenv("PORT"); port == 8080 && envPort != "" {
		port, _ = strconv.Atoi(envPort)
	}
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if lineChannelID == "" {
		lineChannelID = os.Getenv("LINE_CHANNEL_ID")
	}
	if lineChannelSecret == "" {
		lineChannelSecret = os.Getenv("LINE_CHANNEL_SECRET")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err.Error())
	}

	grpcC, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}

	public, err := fs.Sub(webDist, "web/dist")
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
					Logger:             log,
					UserClient:         proto.NewUserClient(grpcC),
					MessageClient:      proto.NewMessageClient(grpcC),
					MessageCreatedChan: make(chan *model.Message),
				},
			}),
	))

	lineLoginClient, err := line_login_sdk.New(lineChannelID, lineChannelSecret)
	if err != nil {
		log.Fatal(err.Error())
	}

	httpS := &http.Server{Handler: cors.AllowAll().
		Handler(
			pkg.NewAccessLogMiddleware(
				pkg.NewRecoveryMiddleware(
					pkg.NewLINELoginMiddleware(m, lineLoginClient),
				),
			),
		),
	}

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

	gormLogger := zapgorm2.New(log.With(zap.String("system", "gorm")))
	gormLogger.SetAsDefault()
	gormLogger.LogLevel = gorm_logger.Info
	var db *gorm.DB
	if databaseURL == "" {
		sqliteDB, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{Logger: gormLogger})
		if err != nil {
			log.Fatal(err.Error())
		}
		db = sqliteDB
	} else {
		dsn, err := pq.ParseURL(databaseURL)
		if err != nil {
			log.Fatal(err.Error())
		}
		postgresDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
		if err != nil {
			log.Fatal(err.Error())
		}
		db = postgresDB
	}

	userS, err := user.NewServer(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	messageS, err := message.NewServer(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	proto.RegisterUserServer(grpcS, userS)
	proto.RegisterMessageServer(grpcS, messageS)

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
