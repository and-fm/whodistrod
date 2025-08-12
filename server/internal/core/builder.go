package core

import (
	"log"
	"reflect"

	"go.uber.org/dig"
)

type builder struct {
	container     *dig.Container
	serverStarted bool
}

type Builder interface {
	// AddSingleton registers a singleton dependency in the container.
	//
	// The only parameter is a function that returns the type of the dependency.
	//
	// This function will log and panic if the registration fails.
	//
	// Example usage:
	//
	//	app.AddSingleton(NewDatabase)
	//	app.AddSingleton(func() *Config { return &Config{} })
	AddSingleton(newSingletonFunc any)

	// AddHandler registers a handler with the web server.
	//
	// The handler must implement the Handler interface.
	// The Register() method of the handler will be called automatically.
	// It is up to the handler to register its routes with the web server within the Register() method.
	//
	// Similar to AddSingleton, it takes a single parameter, a function that returns the handler type.
	//
	// Since handlers are invoked immediately, you must register all dependencies that the handler requires before calling this method.
	// This function will log and panic if the registration fails.
	//
	// Example usage:
	//
	//	app.AddHandler(NewHealthcheckHandler)
	AddHandler(newHandlerFunc any)

	// Run starts the server using the provided function.
	//
	// The function must return a type that implements the Server interface.
	// The Start() method of the server will be called automatically, this should contain the logic to start the server and listen for http requests.
	//
	// This function will log and panic if the invocation fails.
	//
	// Example usage:
	//
	//	app.Run(NewServer)
	Run(newServerFunc any)
}

func NewBuilder() Builder {
	return &builder{container: dig.New(), serverStarted: false}
}

func (b *builder) AddSingleton(newSingletonFunc any) {
	err := b.container.Provide(newSingletonFunc)
	if err != nil {
		log.Fatalf("Error providing dependency: %v. ", err)
	}
}

func (b *builder) AddHandler(newHandlerFunc any) {
	funcType := reflect.TypeOf(newHandlerFunc)

	scope := b.container.Scope("")

	err := scope.Provide(newHandlerFunc)
	if err != nil {
		log.Fatalf("Error providing dependency: %v. ", err)
	}

	if !funcType.Out(0).Implements(reflect.TypeOf((*Handler)(nil)).Elem()) {
		log.Fatalf("AddHandler expected a function that returns a type implementing the Handler interface, got %s", funcType.Out(0))
	}

	err = scope.Invoke(func(h Handler) {
		h.Register()
	})
	if err != nil {
		log.Fatalf("Error invoking handler: %v. Did you remember to register all dependencies before the handlers in main.go?: %v", funcType.In(0), err)
	}
}

func (b *builder) Run(newServerFunc any) {
	if b.serverStarted {
		log.Fatalf("Server has already been started. You can only call Run once.")
	}

	funcType := reflect.TypeOf(newServerFunc)

	err := b.container.Provide(newServerFunc)
	if err != nil {
		log.Fatalf("Error providing server: %v. ", err)
	}

	if !funcType.Out(0).Implements(reflect.TypeOf((*Server)(nil)).Elem()) {
		log.Fatalf("Run expected a function that returns a type implementing the Server interface, got %s", funcType.Out(0))
	}

	err = b.container.Invoke(func(s Server) {
		s.Start()
	})
	if err != nil {
		log.Fatalf("Failed to invoke the server. Did you remember to register all dependencies in main.go?: %v", err)
	}

	b.serverStarted = true
}
