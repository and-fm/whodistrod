package core

type Server interface {
	// This method should hold the logic to start the server and start listening for http requests.
	// It is automatically invoked by the builder when the new server function is passed to the Run method.
	Start()
}
