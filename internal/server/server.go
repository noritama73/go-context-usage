package server

import (
	"context"
	"fmt"

	"context_usage/internal/auth"
	"context_usage/internal/handlers"
	"context_usage/internal/session"
)

type MyServer struct {
	router map[string]handlers.MyHandleFunc
}

func (srv *MyServer) ListenAndServe() {
	for {
		var path, token string
		fmt.Print("path: ")
		fmt.Scan(&path)
		fmt.Print("token: ")
		fmt.Scan(&token)

		ctx := session.SetSessionID(context.Background())
		go srv.Request(ctx, path, token)
	}
}

func (srv *MyServer) Request(ctx context.Context, path string, token string) {
	// リクエストオブジェクト作成
	var req handlers.MyRequest
	req.SetPath(path)

	// (key:authToken <=> value:token)をcontextに入れる
	ctx = auth.SetAuthToken(ctx, token)

	// ルーティング操作
	if handler, ok := srv.router[req.GetPath()]; ok {
		handler(ctx, req)
	} else {
		handlers.NotFoundHandler(ctx, req)
	}
}

func DefaultServer() *MyServer {
	r := make(map[string]handlers.MyHandleFunc)
	r["greet"] = handlers.GetGreeting
	return &MyServer{
		router: r,
	}
}
