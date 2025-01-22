package server

type Server interface {
	IsEnabled() bool
	Serve()
	Shutdown()
}
