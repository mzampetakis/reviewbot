package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/exp/slog"
	"os"
	"reviewbot/cmd/reviewbot/api"
	"reviewbot/internal/database"
	"reviewbot/internal/domain/orders"
	"reviewbot/internal/env"
	"reviewbot/internal/version"
	"reviewbot/pkg/responsegenerator/dummygenerator"
	"reviewbot/pkg/sentimentanalyzer/dummyanalyzer"
	"runtime/debug"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	var cfg api.ApplicationConfig

	dbHost := env.GetString("DB_HOST", "myreviewbotdb")
	dbUser := env.GetString("DB_USER", "user")
	dbPass := env.GetString("DB_PASSWORD", "pass")
	dbName := env.GetString("DB_NAME", "myreviewbot")
	dbPort := env.GetString("DB_PORT", "3306")

	cfg.BaseURL = env.GetString("BASE_URL", "")
	cfg.HttpPort = env.GetInt("HTTP_PORT", 4444)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	cfg.DB.DSN = env.GetString("DB_DSN", dsn)
	cfg.DB.Automigrate = env.GetBool("DB_AUTOMIGRATE", true)

	showVersion := flag.Bool("version", false, "display version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	db, err := database.New(cfg.DB.DSN, cfg.DB.Automigrate)
	if err != nil {
		logger.Error("Could not connect to DB " + dbName + " at host: " + dbHost + ":" + dbPort + ". Error: " + err.Error())
		return err
	}
	defer db.Close()

	logger.Info("version: " + version.Get())

	app := api.Application{
		Config: cfg,
		DB:     db,
		Logger: logger,
	}
	logger.Info("Starting...")
	ordersRepo := orders.NewDatabaseRepository(db.DB)

	ordersService := orders.NewService(ordersRepo, dummygenerator.NewDummyGenerator(),
		dummyganalyzer.NewDummyAnalyzer(), logger)
	srv := api.NewServer(ordersService, &app)
	logger.Info("Running...")
	return srv.Serve()
}
