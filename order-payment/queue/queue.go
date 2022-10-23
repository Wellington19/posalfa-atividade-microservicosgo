package queue

import (
	"fmt"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

func Connect(exchange string) *amqp.Channel {
	dsn := "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + ":" + os.Getenv("RABBITMQ_DEFAULT_PORT")

	conn, err := amqp.Dial(dsn)
	if err != nil {
		panic(err.Error())
	}

	channel, err := conn.Channel()
	if err != nil {
		panic(err.Error())
	}

	err = channel.ExchangeDeclare(exchange, "direct", false, false, false, false, nil)
	if err != nil {
		panic(err.Error())
	}

	return channel
}

func Notify(payload []byte, exchange string, routingKey string, ch *amqp.Channel) {
	err := ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(payload),
		})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Message sent")
}

func StartConsuming(queue string, ch *amqp.Channel, in chan []byte) {
	q, err := ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err.Error())
	}

	err = ch.QueueBind(queue, "", strings.Replace(queue, "_queue", "_ex", 1), false, nil)
	if err != nil {
		panic(err.Error())
	}

	msgs, err := ch.Consume(
		q.Name,
		"checkout",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		panic(err.Error())
	}

	go func() {
		for m := range msgs {
			in <- []byte(m.Body)
		}
		close(in)
	}()
}
