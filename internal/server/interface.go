package server

type Interface interface {
	IsEnabled() bool
	Serve()
	Shutdown()
}
