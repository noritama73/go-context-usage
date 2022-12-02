package db

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type DB struct {
	Data map[int]string
}

var DefaultDB = DB{
	Data: make(map[int]string),
}

func init() {
	DefaultDB.Data[4] = "hogeData"
}

func (db DB) Search(ctx context.Context, id int) chan string {
	out := make(chan string)
	go func() {
		rand.Seed(time.Now().UnixNano())
		randTime := rand.Intn(4)
		fmt.Printf("wait %d second\n", randTime)
		select {
		case <-time.After(time.Duration(randTime) * time.Second):
			out <- db.Data[id]
			return
		case <-ctx.Done():
			break
		}
		close(out)
	}()
	return out
}
