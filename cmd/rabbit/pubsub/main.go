package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

const exchange = "go_ex"

func main() {
	r, connectErr := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if connectErr != nil {
		panic(connectErr)
	}

	channel, channelErr := r.Channel()
	if channelErr != nil {
		panic(channelErr)
	}

	err := channel.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	_, queErr := channel.QueueDeclare(
		"go_q1",
		true,
		false,
		false,
		false,
		nil,
	)
	if queErr != nil {
		panic(queErr)
	}

	go subscribe(r, exchange)
	go subscribe(r, exchange)
	go subscribe(r, exchange)

	i := 0
	for {
		err := channel.Publish(
			exchange,
			"",
			false,
			false,
			amqp.Publishing{
				Body: []byte(fmt.Sprintf("message %d", i)),
			},
		)
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(200 * time.Millisecond)
		i++
	}
}

func subscribe(conn *amqp.Connection, ex string) {
	ch, chErr := conn.Channel()
	if chErr != nil {
		panic(chErr)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			panic(err)
		}
	}(ch)

	que, queErr := ch.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	if queErr != nil {
		panic(queErr)
	}

	defer func(ch *amqp.Channel, name string, ifUnused, ifEmpty, noWait bool) {
		_, err := ch.QueueDelete(name, ifUnused, ifEmpty, noWait)
		if err != nil {
			panic(err)
		}
	}(ch, que.Name, false, false, false)

	bindErr := ch.QueueBind(que.Name, "", ex, false, nil)
	if bindErr != nil {
		panic(bindErr)
	}

	consume("c", ch, que.Name)
}

func consume(consumeName string, ch *amqp.Channel, q string) {
	cons, consErr := ch.Consume(
		q,
		consumeName,
		true,
		false,
		false,
		false,
		nil,
	)
	if consErr != nil {
		panic(consErr)
	}

	for msg := range cons {
		fmt.Printf("%s: %s\n", consumeName, msg.Body)
	}
}
