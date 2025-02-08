package util

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"time"
)

// httpServerManager - manages the lifecycle of httpServers
type httpServerManager struct {
	Server   *http.Server
	Listener net.Listener
}

var (
	serverManager *httpServerManager
)

func StartHTTPServer(
	addr string,
	handler http.Handler,
	readTimeout,
	writeTimeout,
	readHeaderTimeout time.Duration,
) (lsnr net.Listener, srv *http.Server, err error) {
	// server config with timeouts
	srv = &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
	}

	// if tlsConfig != nil {
	// 	srv.TLSConfig = tlsConfig
	// }

	lsnr, err = net.Listen("tcp", addr)
	if err != nil {
		return
	}
	err = srv.Serve(lsnr)
	if err != nil {
		return
	}

	serverManager.Listener = lsnr
	serverManager.Server = srv

	return
}

func ShutdownHTTPServer() {
	//Stop the current server
	lsnr := serverManager.Listener
	svr := serverManager.Server

	if svr != nil {
		fmt.Println("Stopping old server")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := svr.Shutdown(ctx); err != nil {
			fmt.Printf("error: server failed to exit gracefully: %v\n", err)
		}
	} else {
		fmt.Println("error: svr is nil in shutdown")
	}
	if lsnr != nil {
		//stopping the listener as well
		if err := lsnr.Close(); err != nil {
			fmt.Printf("erorr: stopping listener: %v\n", err)
		}
	} else {
		fmt.Println("error: nil listener provided to stop")
	}
	fmt.Println("Successfully stopped server")
}

func PanicHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	w.WriteHeader(500)
	GetGlobalLogger(r.Context()).Println("panic stack trace", string(debug.Stack()))
	fmt.Fprintln(w, "500 - Internal error")
}
