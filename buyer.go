package main

import (
	"container/heap"
	"context"
	"math"
)

type BuyWorker struct {
	masterNseller_buyer_chan  chan Order
	masterNbuyer_seller_chan  chan Order
	delete_worker_master_chan chan<- uint32
	done_buyer_master_chan    chan<- struct{}
	id_order_map              map[uint32]*Order
	buy_orderbook             *OrderBook
}

func InitBuyWorker(masterNseller_buyer_chan chan Order, masterNbuyer_seller_chan chan Order, delete_worker_master_chan chan<- uint32, done_buyer_master_chan chan<- struct{}) *BuyWorker {
	buy_orderbook := InitOrderBook()
	bw := &BuyWorker{
		masterNseller_buyer_chan:  masterNseller_buyer_chan,
		masterNbuyer_seller_chan:  masterNbuyer_seller_chan,
		delete_worker_master_chan: delete_worker_master_chan,
		done_buyer_master_chan:    done_buyer_master_chan,
		id_order_map:              make(map[uint32]*Order),
		buy_orderbook:             buy_orderbook,
	}
	return bw
}

func (bw *BuyWorker) work(ctx context.Context) {
	for {
		select {
		case o := <-bw.masterNseller_buyer_chan:
			switch o.order_type {
			case inputBuy: // sent from seller worker
				bw.addBuyOrder(o)
				o.done_client_chan <- struct{}{}
			case inputSell: // sent from master
				bw.matchSellOrder(o)
				bw.done_buyer_master_chan <- struct{}{}

			default: // cancel from master
				bw.cancelBuyOrder(o.order_id)
				o.done_client_chan <- struct{}{}
				// Signal done to master
				bw.done_buyer_master_chan <- struct{}{}
			}

		case <-ctx.Done():
			return
		}
	}
}

func (bw *BuyWorker) addBuyOrder(buy_order Order) {
	bo := &buy_order
	bo.timestamp = GetCurrentTimestamp()
	bw.id_order_map[bo.order_id] = bo
	heap.Push(bw.buy_orderbook.order_pqueue, bo)
	in := input{inputBuy, bo.order_id, bo.price, bo.count, bo.instrument}
	outputOrderAdded(in, bo.timestamp)
}

func (bw *BuyWorker) matchSellOrder(sell_order Order) {
	for bw.buy_orderbook.order_pqueue.Len() > 0 && bw.buy_orderbook.order_pqueue.Top().price >= sell_order.price {
		// get the top buy order
		buy_order := heap.Pop(bw.buy_orderbook.order_pqueue).(*Order)
		// get the minimum count of the two orders to match
		match_count := uint32(math.Min(float64(buy_order.count), float64(sell_order.count)))
		// increase execution id of the resting order
		buy_order.execution_id += 1
		outputOrderExecuted(buy_order.order_id, sell_order.order_id, buy_order.execution_id, buy_order.price, match_count, GetCurrentTimestamp())
		// update the count of the two orders
		buy_order.count -= match_count
		sell_order.count -= match_count
		// if the buy order still has count, push it back to the heap
		if buy_order.count > 0 {
			heap.Push(bw.buy_orderbook.order_pqueue, buy_order)
		} else {
			// else, remove from order_map
			delete(bw.id_order_map, buy_order.order_id)
			// send delete order_id to master
			bw.delete_worker_master_chan <- buy_order.order_id
		}
		// if the sell order has no count, break
		if sell_order.count == 0 {
			break
		}
	}
	if sell_order.count > 0 {
		// send to sell worker
		bw.masterNbuyer_seller_chan <- sell_order
	} else {
		// send delete order_id to master
		bw.delete_worker_master_chan <- sell_order.order_id
		// signal order done being matched
		sell_order.done_client_chan <- struct{}{}
	}
}

func (bw *BuyWorker) cancelBuyOrder(id uint32) {
	if _, ok := bw.id_order_map[id]; ok {
		delete(bw.id_order_map, id)
		for i, buy_order := range *bw.buy_orderbook.order_pqueue {
			if buy_order.order_id == id {
				heap.Remove(bw.buy_orderbook.order_pqueue, i)
				break
			}
		}
		in := input{inputCancel, id, 0, 0, ""}
		outputOrderDeleted(in, true, GetCurrentTimestamp())
		// send delete order_id to master
		bw.delete_worker_master_chan <- id

	} else {
		in := input{inputCancel, id, 0, 0, ""}
		outputOrderDeleted(in, false, GetCurrentTimestamp())

	}
}
