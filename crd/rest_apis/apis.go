package restapis

import (
	"net"
	"net/http"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("controller-runtime.probe-api")

func ServeHTTP() {

	mux := http.NewServeMux()
	mux.Handle(GET_PROBES_PATH, GetProbesHandler())
	server := http.Server{
		Handler: mux,
	}
	// Run the server
	go func() {
		log.Info("starting probes server", "path", GET_PROBES_PATH)
		listener, err := net.Listen("tcp", ":2709")
		if err != nil {
			log.Error(err, "Could not listen the given address")
		}
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Error(err, "Could not serve endpoint")
		}
	}()

}
