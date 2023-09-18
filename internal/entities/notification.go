package entities

import "github.com/valyala/fastjson"

// Notification represents message payload from RabbitMQ.
type Notification struct {
	To    []string
	Title string
	Body  string
	Data  string
}

// Serialize converts Notification to JSON string.
func (n *Notification) Serialize() []byte {
	var b fastjson.Arena
	obj := b.NewObject()

	toArr := b.NewArray()
	for i, _ := range n.To {
		toArr.SetArrayItem(i, b.NewString(n.To[i]))
	}
	obj.Set("to", toArr)
	obj.Set("title", b.NewString(n.Title))
	obj.Set("body", b.NewString(n.Body))
	obj.Set("data", b.NewString(n.Data))

	return obj.MarshalTo(nil)
}

// ParseNotification parses Notification from JSON string.
func ParseNotification(data []byte) (*Notification, error) {
	p := fastjson.Parser{}
	val, err := p.ParseBytes(data)
	if err != nil {
		return nil, err
	}

	// Parse the 'To' field from a JSON array
	toArr := val.GetArray("to")
	toSlice := make([]string, len(toArr))
	for i, v := range toArr {
		toSlice[i] = string(v.GetStringBytes())
	}

	return &Notification{
		To:    toSlice,
		Title: string(val.GetStringBytes("title")),
		Body:  string(val.GetStringBytes("body")),
		Data:  string(val.GetStringBytes("data")),
	}, nil
}
