package dto

import "github.com/valyala/fastjson"

// Notification represents message payload from RabbitMQ
type Notification struct {
	To    string
	Title string
	Body  string
	Icon  string
}

func SerializeNotification(n Notification) string {
	var b fastjson.Arena
	obj := b.NewObject()
	obj.Set("to", b.NewString(n.To))
	obj.Set("title", b.NewString(n.Title))
	obj.Set("body", b.NewString(n.Body))
	obj.Set("icon", b.NewString(n.Icon))

	return string(obj.MarshalTo(nil))
}

func ParseNotification(data string) (Notification, error) {
	p := fastjson.Parser{}
	val, err := p.Parse(data)
	if err != nil {
		return Notification{}, err
	}

	return Notification{
		To:    string(val.GetStringBytes("to")),
		Title: string(val.GetStringBytes("title")),
		Body:  string(val.GetStringBytes("body")),
		Icon:  string(val.GetStringBytes("icon")),
	}, nil
}
