package server

type Interface interface {
	Serve()
	Shutdown()
}
