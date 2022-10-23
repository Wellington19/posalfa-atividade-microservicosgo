package main

import (
	"encoding/json"
	"fmt"
	"order-checkout/db"
	"order-checkout/queue"
	"time"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/streadway/amqp"
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

func createOrder(payload []byte) Order {
	var order Order
	json.Unmarshal(payload, &order)

	uuid, _ := uuid.NewV4()
	order.Uuid = uuid.String()
	order.Status = "pendente"
	order.CreatedAt = time.Now()
	saveOrder(order)
	return order
}

func notifyOrderCreated(order Order, ch *amqp.Channel) {
	json, _ := json.Marshal(order)
	queue.Notify(json, "order_ex", "", ch)
}

func main() {
	in := make(chan []byte)
	connection := queue.Connect("checkout_ex")

	queue.StartConsuming("checkout_queue", connection, in)
	for payload := range in {
		notifyOrderCreated(createOrder(payload), connection)
		fmt.Println(string(payload))
	}
}
