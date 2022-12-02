package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"context_usage/internal/auth"
	"context_usage/internal/db"
)

type MyRequest struct {
	path string
}

func (req *MyRequest) SetPath(path string) {
	req.path = path
}

func (req *MyRequest) GetPath() string {
	return req.path
}

type MyResponse struct {
	Code int
	Body string
	Err  error
}

type MyHandleFunc func(context.Context, MyRequest)

var GetGreeting MyHandleFunc = func(ctx context.Context, req MyRequest) {
	var res MyResponse

	fmt.Println("access")

	// トークンからユーザー検証→ダメなら即return
	userID, err := auth.VerifyAuthToken(ctx)
	if err != nil {
		res = MyResponse{Code: 403, Err: err}
		fmt.Println(res)
		return
	}
	fmt.Println("authorized ID:", userID)

	// DBリクエストをいつタイムアウトさせるかcontext経由で設定
	dbReqCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	//DBからデータ取得
	rcvChan := db.DefaultDB.Search(dbReqCtx, userID)
	data, ok := <-rcvChan

	// DBリクエストがタイムアウトしていたら408で返す
	if !ok {
		res = MyResponse{Code: 408, Err: errors.New("DB request timeout")}
		fmt.Println(res)
		return
	}
	fmt.Println("fetch data")

	// レスポンスの作成
	res = MyResponse{
		Code: 200,
		Body: fmt.Sprintf("From path %s, Hello! your ID is %d\ndata → %s", req.path, userID, data),
	}

	// レスポンス内容を標準出力(=本物ならnet.Conn)に書き込み
	fmt.Println(res)
}

var NotFoundHandler MyHandleFunc = func(ctx context.Context, req MyRequest) {
	var res MyResponse

	res = MyResponse{
		Code: 404,
		Body: fmt.Sprintf("path \"%s\" was not found ", req.path),
	}

	fmt.Println(res)
}
