package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/The-Ogulgozel/Banking-system/internal/config"
	"github.com/The-Ogulgozel/Banking-system/internal/db"
	"github.com/The-Ogulgozel/Banking-system/internal/handlers"
	"github.com/The-Ogulgozel/Banking-system/internal/repository"
	"github.com/The-Ogulgozel/Banking-system/internal/router"
	"github.com/The-Ogulgozel/Banking-system/internal/usecase"
)

func Run() error {
	cfg, err := config.Load("config.yaml")

	if err != nil {
		return err
	}

	if err := db.RunMigrations(cfg.Database.DSN(), "migrations"); err != nil {
		return err
	}

	ctx := context.Background()
	pool, err := db.NewPostgresPool(ctx, cfg.Database.DSN())
	if err != nil {
		return err
	}
	defer pool.Close()

	accountRepo := repository.NewAccountRepo(pool)
	transactionRepo := repository.NewTransactionRepo(pool)
	accountUsecase := usecase.NewAccountUsecase(accountRepo)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo, accountRepo)

	accountHandler := handlers.NewAccountsHandler(accountUsecase)
	transactionHandler := handlers.NewTransactionHandler(transactionUsecase)

	r := router.NewRouter(cfg, &router.RouterDeps{
		AccountsHandler:    accountHandler,
		TransactionHandler: transactionHandler,
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	log.Printf("server started on port %s", srv.Addr)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Printf("server shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	log.Printf("server exited cleanly")
	return nil

}
