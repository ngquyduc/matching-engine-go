package main

import "container/heap"

// order book for a particular instrument
type OrderBook struct {
	order_pqueue *OrderPQueue
}

func InitOrderBook() *OrderBook {
	order_pqueue := make(OrderPQueue, 0)
	new_orderbook := &OrderBook{
		order_pqueue: &order_pqueue,
	}
	heap.Init(new_orderbook.order_pqueue)
	return new_orderbook
}

// implementation of a priority queue to store buy/sell orders
type OrderPQueue []*Order

func (pq *OrderPQueue) Len() int {
	return len(*pq)
}

func (pq *OrderPQueue) Less(i, j int) bool {
	return (*pq)[i].Less((*pq)[j])
}

func (pq *OrderPQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

func (pq *OrderPQueue) Push(x any) {
	order := x.(*Order)
	*pq = append(*pq, order)
}

func (pq *OrderPQueue) Pop() any {
	old := *pq
	n := len(old)
	pop_order := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return pop_order
}

func (pq *OrderPQueue) Top() *Order {
	if len(*pq) > 0 {
		return (*pq)[0]
	}
	return nil
}
