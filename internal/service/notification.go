package service

import "encoding/json"

type Notification struct {
	To     string `json:"to"`
	Header string `json:"header"`
	Body   string `json:"body"`
}

func (n Notification) Serialize() ([]byte, error) {
	jsonData, err := json.Marshal(n)
	return jsonData, err
}

func ParseNotification(data []byte) (Notification, error) {
	n := Notification{}
	err := json.Unmarshal(data, &n)
	return n, err
}
