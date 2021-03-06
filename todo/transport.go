package todo

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type EndPoints struct {
	GetTodoEndpoint    endpoint.Endpoint
	TodoPagesEndpoint  endpoint.Endpoint
	TodoSearchEndpoint endpoint.Endpoint
	CreateTodoEndpoint endpoint.Endpoint
	UpdateTodoEndpoint endpoint.Endpoint
}

func transportLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("msg", "calling endpoint")
			defer logger.Log("msg", "called endpoint")
			return next(ctx, request)
		}
	}
}

func MakeEndPoints(svc TodoService, logger log.Logger) EndPoints {
	var getTodoEndpoint endpoint.Endpoint
	getTodoEndpoint = makeGetTodoEndpoint(svc)
	getTodoEndpoint = transportLoggingMiddleware(log.With(logger, "method", "getTodo"))(getTodoEndpoint)

	var todoPagesEndpoint endpoint.Endpoint
	todoPagesEndpoint = makeTodoPagesEndpoint(svc)
	todoPagesEndpoint = transportLoggingMiddleware(log.With(logger, "method", "TodoPages"))(todoPagesEndpoint)

	var todoSearchEndpoint endpoint.Endpoint
	todoSearchEndpoint = makeTodoSearchEndpoint(svc)
	todoSearchEndpoint = transportLoggingMiddleware(log.With(logger, "method", "TodoSearch"))(todoSearchEndpoint)

	var createTodoEndpoint endpoint.Endpoint
	createTodoEndpoint = makeCreateTodoEndpoint(svc)
	createTodoEndpoint = transportLoggingMiddleware(log.With(logger, "method", "createTodo"))(createTodoEndpoint)

	var updateTodoEndpoint endpoint.Endpoint
	updateTodoEndpoint = makeUpdateTodoEndpoint(svc)
	updateTodoEndpoint = transportLoggingMiddleware(log.With(logger, "method", "updateTodo"))(updateTodoEndpoint)

	return EndPoints{
		GetTodoEndpoint:    getTodoEndpoint,
		TodoPagesEndpoint:  todoPagesEndpoint,
		TodoSearchEndpoint: todoSearchEndpoint,
		CreateTodoEndpoint: createTodoEndpoint,
		UpdateTodoEndpoint: updateTodoEndpoint,
	}
}
