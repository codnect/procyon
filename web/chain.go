package web

import "context"

type chainFunc func(ctx any) error

type HandlerChain struct {
	pathVariableNames []string
	interceptors      []func(ctx context.Context) (bool, error)
	function          *HandlerFunction

	handlerIndex         int
	afterCompletionIndex int
}

func createExecutionChain() *HandlerChain {
	return nil
}
