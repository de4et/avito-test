package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	readTimeout  = 10
	writeTimeout = 30
)

func NewServer(routes http.Handler) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes,
		IdleTimeout:  time.Minute,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
	}

	return server
}
