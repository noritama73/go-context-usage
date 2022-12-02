package main

import (
	"context_usage/internal/server"
)

func main() {
	srv := server.DefaultServer()
	srv.ListenAndServe()
}
