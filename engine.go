package main

import "C"
import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Engine struct {
	client_dist_chan chan<- Order
}

func InitEngine(ctx context.Context) *Engine {
	new_client_dist_chan := make(chan Order, 1000)
	engine := &Engine{client_dist_chan: new_client_dist_chan}
	InitDistributor(new_client_dist_chan, ctx)
	return engine
}

func (e *Engine) accept(ctx context.Context, conn net.Conn) {
	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	go handleConn(conn, e.client_dist_chan)
}

func handleConn(conn net.Conn, client_dist_chan chan<- Order) {
	defer conn.Close()
	done_client_chan := make(chan struct{})
	for {
		in, err := readInput(conn)
		if err != nil {
			if err != io.EOF {
				_, _ = fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			}
			return
		}
		switch in.orderType {
		case inputCancel:
			cancel_order := Order{inputCancel, in.orderId, "", 0, 0, 0, 0, done_client_chan}
			client_dist_chan <- cancel_order
		case inputBuy:
			buy_order := Order{inputBuy, in.orderId, in.instrument, in.price, in.count, GetCurrentTimestamp(), 0, done_client_chan}
			client_dist_chan <- buy_order
		case inputSell:
			sell_order := Order{inputSell, in.orderId, in.instrument, in.price, in.count, GetCurrentTimestamp(), 0, done_client_chan}
			client_dist_chan <- sell_order
		}
		// wait for the order to be processed
		<-done_client_chan
	}
}

func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano()
}
