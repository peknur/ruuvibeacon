package ruuvibeacon

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func httpListenAndServer(addr string) chan bool {
	srv := &http.Server{Addr: addr}
	done := make(chan bool)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		// We received an interrupt/kill signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Println(err)
		}
		close(done)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatal(err)
	}
	return done
}
