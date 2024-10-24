package main

// for cancel, all other fields than type and id are just dummy values
type Order struct {
	order_type       inputType
	order_id         uint32
	instrument       string
	price            uint32
	count            uint32
	timestamp        int64
	execution_id     uint32
	done_client_chan chan<- struct{} // signal order done to engine
}

type OrderHeader struct {
	order_type inputType
	instrument string
}

// less function
func (o *Order) Less(other *Order) bool {
	// buy order -> get the higher price
	if o.order_type == inputBuy {
		if o.price != other.price {
			return o.price > other.price
		} else { // same price -> get the older order
			return o.timestamp < other.timestamp
		}
	} else { // sell order -> get the lower price
		if o.price != other.price {
			return o.price < other.price
		} else { // same price -> get the older order
			return o.timestamp < other.timestamp
		}
	}
}
