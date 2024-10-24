package main

import (
	"context"
)

// MasterWorker order worker to handle incoming orders from distributor
type MasterWorker struct {
	dist_master_chan          <-chan Order
	masterNseller_buyer_chan  chan Order
	masterNbuyer_seller_chan  chan Order
	delete_worker_master_chan chan uint32
	delete_master_dist_chan   chan<- uint32
	done_buyer_master_chan    chan struct{}
	done_seller_master_chan   chan struct{}
	id_inputtype_map          map[uint32]inputType
	current_sell_order_price  int64
	current_buy_order_price   int64
	current_sell_order_id     int64
	current_buy_order_id      int64
}

func InitWorkers(ctx context.Context, dist_master_chan <-chan Order, delete_master_dist_chan chan<- uint32) {
	master := InitMaster(dist_master_chan, delete_master_dist_chan)
	buyer := InitBuyWorker(master.masterNseller_buyer_chan, master.masterNbuyer_seller_chan, master.delete_worker_master_chan, master.done_buyer_master_chan)
	seller := InitSellWorker(master.masterNbuyer_seller_chan, master.masterNseller_buyer_chan, master.delete_worker_master_chan, master.done_seller_master_chan)

	go master.work(ctx)
	go buyer.work(ctx)
	go seller.work(ctx)
}

func InitMaster(dist_master_chan <-chan Order, delete_master_dist_chan chan<- uint32) *MasterWorker {
	master_seller_chan := make(chan Order, 1000)
	done_seller_master_chan := make(chan struct{}, 1)
	done_seller_master_chan <- struct{}{}

	master_buyer_chan := make(chan Order, 1000)
	done_buyer_master_chan := make(chan struct{}, 1)
	done_buyer_master_chan <- struct{}{}

	master := &MasterWorker{
		dist_master_chan:          dist_master_chan,
		masterNbuyer_seller_chan:  master_seller_chan,
		done_buyer_master_chan:    done_buyer_master_chan,
		masterNseller_buyer_chan:  master_buyer_chan,
		done_seller_master_chan:   done_seller_master_chan,
		delete_worker_master_chan: make(chan uint32, 1000),
		delete_master_dist_chan:   delete_master_dist_chan,
		id_inputtype_map:          make(map[uint32]inputType),
		current_sell_order_price:  -1,
		current_buy_order_price:   -1,
		current_sell_order_id:     -1,
		current_buy_order_id:      -1,
	}

	return master
}

func (master_worker *MasterWorker) work(ctx context.Context) {
	for {
		select {
		case order := <-master_worker.dist_master_chan:
			switch order.order_type {
			case inputBuy:
				<-master_worker.done_seller_master_chan // wait for previous seller worker to finish
				// add new buy order to the map
				master_worker.id_inputtype_map[order.order_id] = inputBuy
				// update current buy order price
				master_worker.current_buy_order_price = int64(order.price)
				// check (and update) to see if there is a current sell order being matched
				select {
				case <-master_worker.done_buyer_master_chan:
					// reset (clear) the current sell order price (no active sell order being matched)
					master_worker.current_sell_order_price = -1
					// release buyer worker channel for next order
					master_worker.done_buyer_master_chan <- struct{}{}
				default: // just fall through
				}
				// if there is current sell (buyer worker) being matched
				if master_worker.current_sell_order_price != -1 {
					// if the current sell order price is higher than the buy order price
					if master_worker.current_sell_order_price > int64(order.price) {
						// send the buy order to seller worker (can send directly without waiting)
						master_worker.masterNbuyer_seller_chan <- order
						continue
					} else {
						// wait for the buyer worker to finish (because buyer worker is matching the order)
						<-master_worker.done_buyer_master_chan
						// send the buy order to seller worker once the resting sell order is added
						master_worker.masterNbuyer_seller_chan <- order
						// release buyer worker channel for next order
						master_worker.done_buyer_master_chan <- struct{}{}
					}
				} else {
					// no current sell order being matched -> send the buy order to seller worker
					master_worker.masterNbuyer_seller_chan <- order
				}

			case inputSell:
				<-master_worker.done_buyer_master_chan // wait for previous buyer worker to finish
				// add new sell order to the map
				master_worker.id_inputtype_map[order.order_id] = inputSell
				// update current sell order price
				master_worker.current_sell_order_price = int64(order.price)
				// check (and update) to see if there is a current buy order being matched
				select {
				case <-master_worker.done_seller_master_chan:
					// reset (clear) the current buy order price (no active buy order being matched)
					master_worker.current_buy_order_price = -1
					// release seller worker channel for next order
					master_worker.done_seller_master_chan <- struct{}{}
				default: // just fall through
				}
				// if there is current buy (seller worker) being matched
				if master_worker.current_buy_order_price != -1 {
					// if the current buy order price is lower than the sell order price
					if master_worker.current_buy_order_price < int64(order.price) {
						// send the sell order to buyer worker (can send directly without waiting)
						master_worker.masterNseller_buyer_chan <- order
						continue
					} else {
						// wait for the seller worker to finish (because seller worker is matching the order)
						<-master_worker.done_seller_master_chan
						// send the sell order to buyer worker once the resting buy order is added
						master_worker.masterNseller_buyer_chan <- order
						// release seller worker channel for next order
						master_worker.done_seller_master_chan <- struct{}{}
					}
				} else {
					// no current buy order being matched -> send the sell order to buyer worker
					master_worker.masterNseller_buyer_chan <- order
				}

			default: // cancel from distributor
				// check if the order exists in the map
				if input_type, ok := master_worker.id_inputtype_map[order.order_id]; !ok {
					in := input{order.order_type, order.order_id, 0, 0, ""}
					outputOrderDeleted(in, false, GetCurrentTimestamp())
					// send done to client
					order.done_client_chan <- struct{}{}
					continue
				} else {
					// delete from map
					delete(master_worker.id_inputtype_map, order.order_id)
					switch input_type {
					case inputBuy:
						// check if the order is being matched in seller worker
						if master_worker.current_buy_order_id == int64(order.order_id) {
							// wait for the seller worker to finish (because seller worker is matching the order)
							<-master_worker.done_seller_master_chan
							master_worker.current_buy_order_id = -1
						}
						<-master_worker.done_buyer_master_chan
						master_worker.masterNseller_buyer_chan <- order
					case inputSell:
						// check if the order is being matched in buyer worker
						if master_worker.current_sell_order_id == int64(order.order_id) {
							// wait for the buyer worker to finish (because buyer worker is matching the order)
							<-master_worker.done_buyer_master_chan
							master_worker.current_sell_order_id = -1
						}
						<-master_worker.done_seller_master_chan
						master_worker.masterNbuyer_seller_chan <- order
					}
					// send order_id to distributor through delete chan
					master_worker.delete_master_dist_chan <- order.order_id
				}
			}
		case order_id := <-master_worker.delete_worker_master_chan:
			// always valid delete
			delete(master_worker.id_inputtype_map, order_id)
			// send order_id to distributor through delete chan
			master_worker.delete_master_dist_chan <- order_id
		case <-ctx.Done():
			return
		}
	}
}
