package infrastructure

import (
	"github.com/valyala/fasthttp"
	"tabellarium/internal/entities"
)

const NotificationURL = "https://exp.host/--/api/v2/push/send"

func SendNotification(n *entities.Notification) (error, []byte) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(NotificationURL)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	payload := n.Serialize()
	req.SetBody(payload)

	err := fasthttp.Do(req, resp)
	if err != nil {
		return err, nil
	}
	return err, resp.Body()
}
