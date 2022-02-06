package main

import (
	"net/http"
	"os"

	"github.com/siongui/go-kit-gqlgen-postgres-todo-example/config"
	"github.com/siongui/go-kit-gqlgen-postgres-todo-example/todo"

	"github.com/99designs/gqlgen/graphql/playground"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Load config
	if err := config.LoadConfig(); err != nil {
		logger.Log("err", err.Error())
		os.Exit(1)
	}
	logger.Log("msg", "App config")
	logger.Log("port", config.Config.App.Port)
	logger.Log("timezone", config.Config.App.Timezone)
	/*
		logger.Log("msg", "Postgres DSN")
		logger.Log("host", config.Config.Database.Postgres.Host)
		logger.Log("port", config.Config.Database.Postgres.Port)
		logger.Log("user", config.Config.Database.Postgres.User)
		logger.Log("password", config.Config.Database.Postgres.Password)
		logger.Log("dbname", config.Config.Database.Postgres.Dbname)
	*/
	gormdsn := "host=" + config.Config.Database.Postgres.Host +
		" user=" + config.Config.Database.Postgres.User +
		" password=" + config.Config.Database.Postgres.Password +
		" dbname=" + config.Config.Database.Postgres.Dbname +
		" port=" + config.Config.Database.Postgres.Port +
		" sslmode=disable TimeZone=" + config.Config.App.Timezone
	logger.Log("gorm_dsn", gormdsn)

	// set up prometheus
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "todo_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "todo_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	// set up todo service and middleware
	var svc todo.TodoService
	svc, err := todo.NewService(gormdsn)
	if err != nil {
		logger.Log("err", err.Error())
		os.Exit(1)
	}
	svc = todo.NewLoggingMiddleware(logger, svc)
	svc = todo.NewInstrumentingMiddleware(requestCount, requestLatency, svc)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", todo.MakeGraphQLHandler(svc, logger))
	http.Handle("/metrics", promhttp.Handler())

	logger.Log("msg", "connect to http://localhost:"+config.Config.App.Port+"/ for GraphQL playground")
	logger.Log("err", http.ListenAndServe(":"+config.Config.App.Port, nil))
}
