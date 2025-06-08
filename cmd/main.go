package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/iamstep4ik/TestTaskOzonBank/graph"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/config"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/log"
	commentservice "github.com/iamstep4ik/TestTaskOzonBank/internal/service/comment_service"
	postservice "github.com/iamstep4ik/TestTaskOzonBank/internal/service/post_service"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/storage"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const defaultPort = "8080"

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	logDir := "./logs"

	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("Error creating log directory: " + err.Error())
	}
	logFile := filepath.Join(logDir, "app.log")
	err := log.Initialize(log.Config{
		LogFile:    logFile,
		LogLevel:   os.Getenv("LOG_LEVEL"),
		MaxSizeMB:  100,
		MaxBackups: 10,
		MaxAgeDays: 30,
		Compress:   true,
		Console:    true,
	})
	if err != nil {
		panic("Error initializing logger: " + err.Error())
	}
	defer log.Sync()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Error("Error loading configuration", zap.Error(err))
		return
	}
	ctx := context.Background()
	dbpool, err := cfg.ConnectDatabase(ctx, cfg)
	if err != nil {
		log.Error("Error connecting to database", zap.Error(err))
		return
	}
	log.Info("Database connection established", zap.String("host", cfg.Database.Host), zap.String("port", cfg.Database.Port), zap.String("name", cfg.Database.Name))
	defer dbpool.Close()
	storageName := os.Getenv("STORAGE_TYPE")
	storage := storage.NewStorage(ctx, dbpool)

	log.Info("Using storage", zap.String("type", storageName))

	postService := postservice.NewPostService(storage, log.GetLogger())
	commentService := commentservice.NewCommentService(storage, log.GetLogger())
	resolver := graph.NewResolver(postService, commentService)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.Websocket{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	log.Info("Starting server", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Error("Error starting server", zap.Error(err))
		return
	}
	log.Info("Server started successfully", zap.String("port", port))

}
