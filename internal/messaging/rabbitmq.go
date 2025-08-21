package messaging

import (
	"User-api/internal/config"
	"User-api/internal/models"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ(cfg *config.Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to RabbitMQ successfully")
	return conn, nil
}

func PublishUserEvent(conn *amqp.Connection, eventType string, payload interface{}) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"user_events", 
		"topic",       
		true,          
		false,         
		false,         
		false,         
		nil,           
	)
	if err != nil {
		return err
	}

	message := models.Message{
		Type:    eventType,
		Payload: payload,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"user_events",      
		"user."+eventType,  
		false,              
		false,              
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	log.Printf("Published %s event to RabbitMQ", eventType)
	return nil
}