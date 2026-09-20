package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/lniklison/bill-splitter-go/internal/bills"
	"github.com/lniklison/bill-splitter-go/internal/bills/controller"
	"github.com/lniklison/bill-splitter-go/internal/bills/repository"
	"github.com/lniklison/bill-splitter-go/internal/database"
	appweb "github.com/lniklison/bill-splitter-go/internal/web"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := database.Open(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := database.Migrate(ctx, pool); err != nil {
		log.Fatalf("apply database migrations: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	billRepository := repository.NewPostgres(pool)
	billService := bills.NewService(billRepository)
	billHTTP := controller.NewHTTP(billService)
	mux.HandleFunc("GET /bills/{id}/shares", billHTTP.GetShares)
	mux.HandleFunc("PUT /bills/{id}/shares", billHTTP.ReplaceShares)
	mux.HandleFunc("/bills/{id}/shares", billHTTP.MethodNotAllowed)
	mux.HandleFunc("/bills", billHTTP.InvalidPath)
	mux.HandleFunc("/bills/", billHTTP.InvalidPath)
	mux.Handle("/", appweb.Handler())
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/bills") && path.Clean(r.URL.Path) != r.URL.Path {
			billHTTP.InvalidPath(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("shut down server: %v", err)
		}
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("serve HTTP: %v", err)
		}
	}
}
