package main

import (
	"context"
	"errors"
	"os"
	"strings"

	chama_app "github.com/gidyon/machama-app/internal/chama"
	"github.com/gidyon/machama-app/internal/chamamember"
	loan_app "github.com/gidyon/machama-app/internal/loan"
	loanproduct "github.com/gidyon/machama-app/internal/loanplan"
	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/internal/moneyaccount"
	transaction_app "github.com/gidyon/machama-app/internal/transaction"
	"github.com/gidyon/machama-app/pkg/api/chama"
	"github.com/gidyon/machama-app/pkg/api/loan"
	"github.com/gidyon/machama-app/pkg/api/transaction"
	"github.com/gidyon/micro/utils/encryption"
	"github.com/gidyon/micro/v2"
	"github.com/gidyon/micro/v2/pkg/config"
	"github.com/gidyon/micro/v2/pkg/healthcheck"
	"github.com/gidyon/micro/v2/pkg/middleware/grpc/auth"
	"github.com/gidyon/micro/v2/pkg/middleware/grpc/zaplogger"
	"github.com/gidyon/micro/v2/utils/errs"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"go.uber.org/zap"

	app_grpc_middleware "github.com/gidyon/micro/v2/pkg/middleware/grpc"
)

func main() {
	ctx := context.Background()

	// Config
	cfg, err := config.New()
	errs.Panic(err)

	// Initialize logger
	errs.Panic(zaplogger.Init(cfg.LogLevel(), ""))

	zaplogger.Log = zaplogger.Log.WithOptions(zap.WithCaller(true))

	appLogger := zaplogger.ZapGrpcLoggerV2(zaplogger.Log)

	app, err := micro.NewService(ctx, cfg, appLogger)
	errs.Panic(err)

	// Recovery middleware
	recoveryUIs, recoverySIs := app_grpc_middleware.AddRecovery()
	app.AddGRPCUnaryServerInterceptors(recoveryUIs...)
	app.AddGRPCStreamServerInterceptors(recoverySIs...)

	// Logging middleware
	logginUIs, loggingSIs := app_grpc_middleware.AddLogging(zaplogger.Log)
	app.AddGRPCUnaryServerInterceptors(logginUIs...)
	app.AddGRPCStreamServerInterceptors(loggingSIs...)

	jwtKey := []byte(strings.TrimSpace(os.Getenv("JWT_SIGNING_KEY")))

	if len(jwtKey) == 0 {
		errs.Panic(errors.New("missing jwt key"))
	}

	// Authentication API
	authAPI, err := auth.NewAPI(&auth.Options{
		SigningKey: jwtKey,
		Issuer:     "MACHAMA APP",
		Audience:   "apis",
	})
	errs.Panic(err)

	// Authentication middleware
	app.AddGRPCUnaryServerInterceptors(grpc_auth.UnaryServerInterceptor(authAPI.AuthorizeFunc))
	app.AddGRPCStreamServerInterceptors(grpc_auth.StreamServerInterceptor(authAPI.AuthorizeFunc))

	app.AddEndpoint("/api/health/ready", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service: app,
		Type:    healthcheck.ProbeReadiness,
	}))

	app.AddEndpoint("/api/health/alive", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service: app,
		Type:    healthcheck.ProbeLiveNess,
	}))

	app.Start(ctx, func() error {
		sqlDB := app.GormDB()
		logger := app.Logger()

		// Automigrations
		if !sqlDB.Migrator().HasTable(&models.Chama{}) {
			errs.Panic(sqlDB.Migrator().AutoMigrate(&models.Chama{}))
		}

		if !sqlDB.Migrator().HasTable(&models.ChamaAccount{}) {
			errs.Panic(sqlDB.Migrator().AutoMigrate(&models.ChamaAccount{}))
		}

		if !sqlDB.Migrator().HasTable(&models.ChamaMember{}) {
			errs.Panic(sqlDB.Migrator().AutoMigrate(&models.ChamaMember{}))
		}

		if !sqlDB.Migrator().HasTable(&models.Loan{}) {
			errs.Panic(sqlDB.Migrator().AutoMigrate(&models.Loan{}))
		}

		if !sqlDB.Migrator().HasTable(&models.LoanProduct{}) {
			errs.Panic(sqlDB.Migrator().AutoMigrate(&models.LoanProduct{}))
		}

		if !sqlDB.Migrator().HasTable(&models.Transaction{}) {
			errs.Panic(sqlDB.Migrator().AutoMigrate(&models.Transaction{}))
		}

		pageHasher, err := encryption.NewHasher(string(jwtKey))
		errs.Panic(err)

		// CHAMA API
		chamaAPI, err := chama_app.NewChamaAPI(ctx, &chama_app.Options{
			SQLDB:         sqlDB,
			PageHasher:    pageHasher,
			Logger:        logger,
			Auth:          authAPI,
			AllowedGroups: append(authAPI.AdminGroups(), "TREASURER"),
		})
		errs.Panic(err)

		chama.RegisterChamaAPIServer(app.GRPCServer(), chamaAPI)
		errs.Panic(chama.RegisterChamaAPIHandler(ctx, app.RuntimeMux(), app.ClientConn()))

		// CHAMA MEMBERS API
		chamaMemberAPI, err := chamamember.NewChamaMemberAPI(ctx, &chamamember.Options{
			SQLDB:         sqlDB,
			PageHasher:    pageHasher,
			Logger:        logger,
			Auth:          authAPI,
			AllowedGroups: append(authAPI.AdminGroups(), "TREASURER", "CHAIRMAN"),
		})
		errs.Panic(err)

		chama.RegisterChamaMemberAPIServer(app.GRPCServer(), chamaMemberAPI)
		errs.Panic(chama.RegisterChamaMemberAPIHandler(ctx, app.RuntimeMux(), app.ClientConn()))

		// TRANSACTIONS API
		transactionAPI, err := transaction_app.NewTransactionAPI(ctx, &transaction_app.Options{
			SQLDB:         sqlDB,
			PageHasher:    pageHasher,
			Logger:        logger,
			Auth:          authAPI,
			AllowedGroups: append(authAPI.AdminGroups(), "TREASURER", "CHAIRMAN"),
		})
		errs.Panic(err)

		transaction.RegisterTransactionAPIServer(app.GRPCServer(), transactionAPI)
		errs.Panic(transaction.RegisterTransactionAPIHandler(ctx, app.RuntimeMux(), app.ClientConn()))

		// CHAMA ACCOUNTS API
		chamaAccountsAPI, err := moneyaccount.NewChamaAccountAPI(ctx, &moneyaccount.Options{
			SQLDB:         sqlDB,
			PageHasher:    pageHasher,
			Logger:        logger,
			Auth:          authAPI,
			AllowedGroups: append(authAPI.AdminGroups(), "TREASURER", "CHAIRMAN"),
		})
		errs.Panic(err)

		transaction.RegisterChamaAccountAPIServer(app.GRPCServer(), chamaAccountsAPI)
		errs.Panic(transaction.RegisterChamaAccountAPIHandler(ctx, app.RuntimeMux(), app.ClientConn()))

		// LOAN API
		loanAPI, err := loan_app.NewLoanAPI(ctx, &loan_app.Options{
			MoneyAccountAPI: chamaAccountsAPI,
			TransactionAPI:  transactionAPI,
			SQLDB:           sqlDB,
			PageHasher:      pageHasher,
			Logger:          logger,
			Auth:            authAPI,
			AllowedGroups:   append(authAPI.AdminGroups(), "TREASURER", "CHAIRMAN"),
		})
		errs.Panic(err)

		loan.RegisterLoanAPIServer(app.GRPCServer(), loanAPI)
		errs.Panic(loan.RegisterLoanAPIHandler(ctx, app.RuntimeMux(), app.ClientConn()))

		// LOAN PRODUCT API
		LoanProductAPI, err := loanproduct.NewLoanProductAPI(ctx, &loanproduct.Options{
			SQLDB:         sqlDB,
			PageHasher:    pageHasher,
			Logger:        logger,
			Auth:          authAPI,
			AllowedGroups: append(authAPI.AdminGroups(), "TREASURER", "CHAIRMAN"),
		})
		errs.Panic(err)

		loan.RegisterLoanProductAPIServer(app.GRPCServer(), LoanProductAPI)
		errs.Panic(loan.RegisterLoanProductAPIHandler(ctx, app.RuntimeMux(), app.ClientConn()))

		return nil
	})

}
