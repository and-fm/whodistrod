package core

type Handler interface {
	// This method should hold the logic to register the handlers routes with the web server or router.
	// It is automatically invoked by the builder when the new handler function is passed to the AddHandler method.
	Register()
}
