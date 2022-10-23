package main

import (
	"encoding/json"
	"fmt"
	"order-payment/db"
	"order-payment/queue"
	"time"
)

type Product struct {
	Uuid    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float32 `json:"price,string"`
}

type Order struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	ProductId string    `json:"product_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at,string"`
}

func saveOrder(order Order) {
	json, _ := json.Marshal(order)
	connection := db.Connect()

	err := connection.Set(order.Uuid, string(json), 0).Err()
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	in := make(chan []byte)
	connection := queue.Connect("payment_ex")

	queue.StartConsuming("payment_queue", connection, in)
	var order Order
	for payload := range in {
		json.Unmarshal(payload, &order)
		saveOrder(order)
		fmt.Println("Payment: ", string(payload))
	}
}
