package main

import (
	"context"
	"fmt"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/scarlet0725/insta-caption-crawler/components/notification/pkg/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelToken := os.Getenv("LINE_CHANNEL_TOKEN")

	if channelSecret == "" || channelToken == "" {
		log.Fatalf("channel secret or channelToken is not set")
	}
	line, err := linebot.New(channelSecret, channelToken)

	if err != nil {
		log.Fatalf("cloud not initialize linebot client")
	}

	lineWebhookHandler := handler.NewLineWebhookHandler(line)
	lineNotificationHandler := handler.NewLineNotificationHandler(line)

	mux := http.NewServeMux()
	mux.Handle("/line", handler.PostOnlyMiddleware(lineWebhookHandler))
	mux.Handle("/notification/line", handler.PostOnlyMiddleware(lineNotificationHandler))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	go srv.ListenAndServe()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	srv.Shutdown(ctx)

}
