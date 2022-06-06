package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

func main() {
	r, connectErr := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if connectErr != nil {
		panic(connectErr)
	}

	channel, channelErr := r.Channel()
	if channelErr != nil {
		panic(channelErr)
	}

	que, queErr := channel.QueueDeclare(
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

	go consume("c1", r, que.Name)
	go consume("c2", r, que.Name)

	i := 0
	for {
		err := channel.Publish(
			"",
			que.Name,
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

func consume(consumeName string, conn *amqp.Connection, q string) {
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
