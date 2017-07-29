package service

import (
	"github.com/julienschmidt/httprouter"
)

type Context struct {
	RepositoryFactory *SquadRepositoryFactory
}

func newContext(config Configuration) (*Context, error) {
	repositoryFactory := SquadRepositoryFactory{config, nil}

	squadService := Context{&repositoryFactory}

	return &squadService, nil
}

type ContextualHandler interface {
	With(service *Context) httprouter.Handle
}

func (context *Context) with(contextual ContextualHandler) httprouter.Handle {
	return contextual.With(context)
}

func (context Context) Close() {
	context.RepositoryFactory.Close()
}
