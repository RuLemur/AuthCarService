package pkg

import (
	"AuthCarService/internal/app"
	"AuthCarService/internal/config"
	"context"
	"github.com/RuLemur/CarService/pkg/endpoint"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	log *logrus.Logger
}

func NewApp() *App {
	return &App{}
}

func (a *App) RunApp(config *config.Config) {
	a.log = logrus.New()
	a.initLogger(a.log)

	a.log.Infof("Create GRPC Client...")
	conn, err := grpc.Dial(config.GrpcServer.Host, grpc.WithInsecure())
	if err != nil {
		a.log.Fatalln("Failed to create client", err)
	}
	defer conn.Close()
	client := endpoint.NewCarServiceClient(conn)

	a.log.Infof("Created")

	api := app.NewApi(client)

	a.log.Infof("Create HTTP Server...")
	handler := api.NewRouter()

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	// Start the server
	//go func() {
	a.log.Infof("Start Serve host %s", config.HttpServer.Host)
	err = srv.ListenAndServe()
	if err != nil {
		a.log.Fatalf("Fail to Serve HTTP Server %e", err)
	}
	//}()
	a.log.Infof("Created")

	// Wait for an interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Attempt a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func (a *App) initLogger(log *logrus.Logger) {
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}
