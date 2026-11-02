package main

import (
	"context"
	"errors"
	"github.com/Vazul15/go-lib-logger"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func Entry(app *app) {
	// ── Application middleware and endpoint registration
	app.router.Register(app.server.Router())
	
	// ── Start under one shutdown context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error { return app.server.Entry(ctx) })

	if err := group.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		logger.FromContext(ctx).Error("server exited with error", "err", err)
	}
	
	// Both runnable units have stopped before the shared Oracle pool is
	// closed, so the reconciler cannot start another database read after the
	// application begins teardown.
	if err := app.database.Close(); err != nil {
		logger.FromContext(ctx).Error("database pool close failed", "err", err)
	}
}