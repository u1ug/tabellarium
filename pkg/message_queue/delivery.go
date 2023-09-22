package message_queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

// Delivery is a struct that can be passed to producer and received from consumer. Contains all required data for message sending.
type Delivery struct {
	Exchange  string                 // The name of the exchange to which the message should be published.
	Key       string                 // Routing key.
	Mandatory bool                   // If set to true, the message will be returned to the producer if it cannot be routed to a queue (e.g., there are no bindings that match the routing key).
	Immediate bool                   // If set to true, the message will be returned to the producer if it cannot be immediately delivered to a consumer (e.g., there are no consumers ready to accept the message).
	Headers   map[string]interface{} // Additional message metadata

	// Copy of amqp.Publishing fields.
	ContentType     string    // MIME content type
	ContentEncoding string    // MIME content encoding
	DeliveryMode    uint8     // Transient (0 or 1) or Persistent (2)// MIME content type
	Priority        uint8     // 0 to 9
	CorrelationId   string    // correlation identifier
	ReplyTo         string    // address to reply to (ex: RPC)
	Expiration      string    // message expiration spec
	MessageId       string    // message identifier
	Timestamp       time.Time // message timestamp
	Type            string    // message type name
	UserId          string    // creating user id - ex: "guest"
	AppId           string    // creating application id

	// The application specific payload of the message.
	Body []byte
}

// FormPublishing converts Delivery and given message body to amqp.Publishing.
func FormPublishing(d Delivery) amqp.Publishing {
	return amqp.Publishing{
		Headers:         amqp.Table(d.Headers),
		ContentType:     d.ContentType,
		ContentEncoding: d.ContentEncoding,
		DeliveryMode:    d.DeliveryMode,
		Priority:        d.Priority,
		CorrelationId:   d.CorrelationId,
		ReplyTo:         d.ReplyTo,
		Expiration:      d.Expiration,
		MessageId:       d.MessageId,
		Timestamp:       d.Timestamp,
		Type:            d.Type,
		UserId:          d.UserId,
		AppId:           d.AppId,
		Body:            d.Body,
	}
}
