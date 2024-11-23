package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"centrifuge_sample_server/api"
)

func main() {
	httpPort := 7000

	httpServerMux := http.NewServeMux()

	cenApi, err := api.New(
		httpServerMux,
		nil,
	)
	if err != nil {
		slog.Error("Centrifugo API failed")
		return
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: httpServerMux,
	}
	go func() {
		slog.Info("HTTP Server try listen", slog.Int("port", httpPort))
		if e := httpServer.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
			slog.Error("HTTP Server listen failed", slog.Any("msg", e))
			panic(e)
		}
	}()

	for {
		time.Sleep(50 * time.Millisecond)
		e := cenApi.BroadcastKV()
		if e != nil {
			slog.Error(e.Error())
		}
	}
}
