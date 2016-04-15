package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

func sleepSeconds(seconds int64) {
	time.Sleep(time.Second * 10)
}

type StopServer struct {
	http.Server
	listener net.Listener
	waiter   sync.WaitGroup
}

func (srv *StopServer) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	var err error
	srv.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(srv.listener)
}

func (srv *StopServer) Serve(l net.Listener) error {
	cur_handler := srv.Handler
	defer func() {
		srv.Handler = cur_handler
	}()
	new_handler := http.NewServeMux()
	new_handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		srv.waiter.Add(1)
		defer srv.waiter.Done()
		cur_handler.ServeHTTP(w, r)
	})
	srv.Handler = new_handler
	return srv.Server.Serve(l)
}

func (srv *StopServer) Stop() error {
	return srv.listener.Close()
}

func (srv *StopServer) WaitUnfinished() {
	srv.waiter.Wait()
}

func mai() {
	httpServer := &StopServer{
		http.Server: http.Server{
			Addr: ":8099",
		},
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Starting handler...<br />\n")
		w.(http.Flusher).Flush()
		sleepSeconds(3)
		io.WriteString(w, "Finishing handler...<br />\n")
	})
	handler.HandleFunc("/down", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Shutting down...<br />\n")
		w.(http.Flusher).Flush()
		httpServer.Stop()
		sleepSeconds(3)
		io.WriteString(w, "Done...<br />\n")
	})
	httpServer.Handler = handler

	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Println("--------------------------------------")
		fmt.Println(err)
		fmt.Printf("%v\n", err)
		fmt.Printf("%#v\n", err)
		fmt.Printf("%V\n", err)
	}
	fmt.Println("====================================")
	fmt.Println("Waiting for waiter.")
	httpServer.WaitUnfinished()
}
