package consumers

import (
	"ApiMessenger/models"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

func SendJSON(model models.RMQMessage) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("RabbitMQ connection error:", err)
		return
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			fmt.Println("RabbitMQ close connection error:", err)
			return
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("RabbitMQ connect channel error:", err)
		return
	}

	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			fmt.Println("RabbitMQ close channel error:", err)
			return
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
		fmt.Println(err)
		return
	}

	fmt.Println(q)

	jsonData, err := json.Marshal(model)

	if err != nil {
		fmt.Println("JSON Marshalling error:", err)
		return
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
