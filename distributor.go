package main

import (
	"context"
)

// Distribute jobs to different goroutines
type Distributor struct {
	client_dist_chan        <-chan Order
	id_instrument_map       map[uint32]string
	instrument_chan_map     map[string]chan<- Order // contains dist_master_chan for each instrument
	delete_master_dist_chan chan uint32
}

func InitDistributor(client_dist_chan <-chan Order, ctx context.Context) {
	d := &Distributor{
		client_dist_chan:        client_dist_chan,
		id_instrument_map:       make(map[uint32]string),
		instrument_chan_map:     make(map[string]chan<- Order, 1000),
		delete_master_dist_chan: make(chan uint32, 1000),
	}

	go d.RunDistributor(ctx)
}

func (d *Distributor) RunDistributor(ctx context.Context) {
	for {
		select {
		case o := <-d.client_dist_chan:
			if o.order_type != inputCancel {
				// add new buy/sell order to order_id_map
				d.id_instrument_map[o.order_id] = o.instrument
				// check if there is a worker for the instrument
				if _, ok := d.instrument_chan_map[o.instrument]; ok {
					// exist -> send order to the corresponding channel
					dist_master_chan := d.instrument_chan_map[o.instrument]
					dist_master_chan <- o
				} else {
					// not exist -> create a new channel -> send order to the created channel
					new_dist_master_chan := make(chan Order)
					d.instrument_chan_map[o.instrument] = new_dist_master_chan
					InitWorkers(ctx, new_dist_master_chan, d.delete_master_dist_chan)
					new_dist_master_chan <- o
				}
			} else {
				// check if the order exists
				if _, ok := d.id_instrument_map[o.order_id]; !ok {
					// not exist -> cancel failed
					in := input{inputCancel, o.order_id, 0, 0, ""}
					outputOrderDeleted(in, false, GetCurrentTimestamp())
					o.done_client_chan <- struct{}{}
				} else {
					// exist -> send cancel order to the corresponding channel
					order_instrument := d.id_instrument_map[o.order_id]
					dist_master_chan := d.instrument_chan_map[order_instrument]
					dist_master_chan <- o
				}
			}
		case order_id := <-d.delete_master_dist_chan:
			delete(d.id_instrument_map, order_id)

		case <-ctx.Done():
			return
		}
	}
}
