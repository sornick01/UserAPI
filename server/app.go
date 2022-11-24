package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sornick01/UserAPI/internal/user"
	"github.com/sornick01/UserAPI/internal/user/delivery"
	"github.com/sornick01/UserAPI/internal/user/repository"
	"github.com/sornick01/UserAPI/internal/user/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	httpServer *http.Server

	useCase user.UseCase
}

func NewApp() *App {
	repo := repository.NewJsonRepo("users.json")

	return &App{
		useCase: usecase.NewDefaultUC(repo),
	}
}

func (a *App) Run(port string) error {
	r := chi.NewRouter()

	delivery.RegisterEndpoints(r, a.useCase)

	a.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := a.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-serverCtx.Done()

	return nil
}
