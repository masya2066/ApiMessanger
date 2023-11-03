package consumers

import (
	"ApiMessenger/models"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func SendJSON(model models.RMQMessage) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("RabbitMQ connection error:", err)
		panic(err)
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			fmt.Println("RabbitMQ close connection error:", err)
			panic(err)
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("RabbitMQ connect channel error:", err)
		panic(err)
	}

	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			fmt.Println("RabbitMQ close channel error:", err)
		}
	}(ch)

	q, err := ch.QueueDeclare(
		"connector",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(q)

	jsonData, err := json.Marshal(model)

	if err != nil {
		log.Fatal("JSON Marshalling error:", err)
	}

	err = ch.Publish(
		"",
		"connector",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
}
