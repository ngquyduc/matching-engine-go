package main

import (
	"container/heap"
	"context"
	"math"
)

type SellWorker struct {
	masterNbuyer_seller_chan  chan Order
	masterNseller_buyer_chan  chan Order
	delete_worker_master_chan chan<- uint32
	done_seller_master_chan   chan<- struct{}
	id_order_map              map[uint32]*Order
	sell_orderbook            *OrderBook
}

func InitSellWorker(masterNbuyer_seller_chan chan Order, masterNseller_buyer_chan chan Order, delete_worker_master_chan chan<- uint32, done_seller_master_chan chan<- struct{}) *SellWorker {
	sell_orderbook := InitOrderBook()
	sw := &SellWorker{
		masterNbuyer_seller_chan:  masterNbuyer_seller_chan,
		masterNseller_buyer_chan:  masterNseller_buyer_chan,
		delete_worker_master_chan: delete_worker_master_chan,
		done_seller_master_chan:   done_seller_master_chan,
		id_order_map:              make(map[uint32]*Order),
		sell_orderbook:            sell_orderbook,
	}
	return sw
}

func (sw *SellWorker) work(ctx context.Context) {
	for {
		select {
		case o := <-sw.masterNbuyer_seller_chan:
			switch o.order_type {
			case inputSell: // sent from buyer worker
				sw.addSellOrder(o)
				o.done_client_chan <- struct{}{}
				// I'm done master
			case inputBuy: // sent from master
				sw.matchBuyOrder(o)
				sw.done_seller_master_chan <- struct{}{}

			default: // cancel from master
				sw.cancelSellOrder(o.order_id)
				o.done_client_chan <- struct{}{}
				// Signal done to master
				sw.done_seller_master_chan <- struct{}{}
			}

		case <-ctx.Done():
			return
		}
	}
}

func (sw *SellWorker) addSellOrder(sell_order Order) {
	so := &sell_order
	so.timestamp = GetCurrentTimestamp()
	sw.id_order_map[so.order_id] = so
	heap.Push(sw.sell_orderbook.order_pqueue, so)
	in := input{inputSell, so.order_id, so.price, so.count, so.instrument}
	outputOrderAdded(in, so.timestamp)
}

func (sw *SellWorker) matchBuyOrder(buy_order Order) {
	for sw.sell_orderbook.order_pqueue.Len() > 0 && sw.sell_orderbook.order_pqueue.Top().price <= buy_order.price {
		// get the top sell order
		sell_order := heap.Pop(sw.sell_orderbook.order_pqueue).(*Order)
		// get the minimum match_count of the two orders to match
		match_count := uint32(math.Min(float64(buy_order.count), float64(sell_order.count)))
		// increase execution id of the resting order
		sell_order.execution_id += 1
		outputOrderExecuted(sell_order.order_id, buy_order.order_id, sell_order.execution_id, sell_order.price, match_count, GetCurrentTimestamp())
		// update the count of the two orders
		sell_order.count -= match_count
		buy_order.count -= match_count
		// if the sell order still has count, push it back to the heap
		if sell_order.count > 0 {
			heap.Push(sw.sell_orderbook.order_pqueue, sell_order)
		} else {
			// else, remove from order_map
			delete(sw.id_order_map, sell_order.order_id)
			// send delete signal to master to remove from map
			sw.delete_worker_master_chan <- sell_order.order_id
		}
		// if the buy order has no count, break
		if buy_order.count == 0 {
			break
		}
	}
	if buy_order.count > 0 {
		// send the buy_order with existing count to buy worker
		sw.masterNseller_buyer_chan <- buy_order
	} else {
		// send delete signal to master to remove from map
		sw.delete_worker_master_chan <- buy_order.order_id
		// signal order done being matched
		buy_order.done_client_chan <- struct{}{}
	}
}

func (sw *SellWorker) cancelSellOrder(id uint32) {

	if _, ok := sw.id_order_map[id]; ok {
		delete(sw.id_order_map, id)
		for i, sell_order := range *sw.sell_orderbook.order_pqueue {
			if sell_order.order_id == id {
				heap.Remove(sw.sell_orderbook.order_pqueue, i)
				break
			}
		}
		in := input{inputCancel, id, 0, 0, ""}
		outputOrderDeleted(in, true, GetCurrentTimestamp())
		sw.delete_worker_master_chan <- id
	} else {
		in := input{inputCancel, id, 0, 0, ""}
		outputOrderDeleted(in, false, GetCurrentTimestamp())
	}
}
