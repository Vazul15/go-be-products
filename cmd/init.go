package main

import (
	"fmt"

	dbadapter "products/internal/connector/db/productRepository"
	"products/internal/connector/db/sqlc"
	"products/internal/connector/restserver/endpoint"
	"products/internal/connector/restserver/router"
	"products/internal/logic/product"
	"github.com/Vazul15/go-lib-logger"
	"github.com/Vazul15/go-lib-restserver"

	starterdb "github.com/Vazul15/starter-lib/pkg/db"
)

func Init() (*app, error) {
	log := logger.MustNew(logger.EnvConfig())
	log.Info("initializing", "application", applicationName)

	// TODO: Ne adjanak vissza errort... ezt a starterlibben kell modositsam, ott panicot kell hivni
	// Database connection and migration
	dbCfg, err := starterdb.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load db config: %w", err)
	}

	db, err := starterdb.Connect(dbCfg)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	if err := starterdb.Migrate(db, dbMigrationPath); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("db migrate: %w", err)
	}

	server := restserver.MustNew(restserver.EnvConfig(applicationName))
	queries := sqlc.New(db)
	productRepository := dbadapter.NewProductRepository(db, queries)

	// ── Orchestration
	productOrchestrator := product.MustNew(productRepository)

	// ── HTTP endpoints
	getProductEndpoint := endpoint.MustNewGetProductEndpoint(productOrchestrator)
	listProductsEndpoint := endpoint.MustNewListProductsEndpoint(productOrchestrator)
	createProductEndpoint := endpoint.MustNewCreateProductEndpoint(productOrchestrator)
	updateProductEndpoint := endpoint.MustNewUpdateProductEndpoint(productOrchestrator)
	deleteProductEndpoint := endpoint.MustNewDeleteProductEndpoint(productOrchestrator)
	createProductVariantEndpoint := endpoint.MustNewCreateProductVariantEndpoint(productOrchestrator)
	listProductVariantsEndpoint := endpoint.MustNewListProductVariantsEndpoint(productOrchestrator)
	getProductVariantEndpoint := endpoint.MustNewGetProductVariantEndpoint(productOrchestrator)
	updateProductVariantEndpoint := endpoint.MustNewUpdateProductVariantEndpoint(productOrchestrator)
	deleteProductVariantEndpoint := endpoint.MustNewDeleteProductVariantEndpoint(productOrchestrator)

	rtr := router.MustNew(
		getProductEndpoint,
		listProductsEndpoint,
		createProductEndpoint,
		updateProductEndpoint,
		deleteProductEndpoint,
		createProductVariantEndpoint,
		listProductVariantsEndpoint,
		getProductVariantEndpoint,
		updateProductVariantEndpoint,
		deleteProductVariantEndpoint,
	)

	return &app{
		database: db,
		router:   rtr,
		server:   server,
	}, nil
}
