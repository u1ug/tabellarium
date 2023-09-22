package user_service

import (
	"github.com/valyala/fastjson"
	"regexp"
)

type RegisterDevicePayload struct {
	UserID string `json:"userID"`
	Token  string `json:"token"`
}

func (p *RegisterDevicePayload) IsValid() bool {
	return !p.isEmpty() && p.tokenIsValid()
}

func (p *RegisterDevicePayload) isEmpty() bool {
	return p.UserID == "" || p.Token == ""
}

func (p *RegisterDevicePayload) tokenIsValid() bool {
	pattern := `ExponentPushToken\[[\x20-\x7E]{10}-[\x20-\x7E]{11}\]`
	matched, err := regexp.MatchString(pattern, p.Token)
	return err == nil && matched
}

func (p *RegisterDevicePayload) Serialize() []byte {
	var b fastjson.Arena
	obj := b.NewObject()
	obj.Set("userID", b.NewString(p.UserID))
	obj.Set("token", b.NewString(p.Token))
	return obj.MarshalTo(nil)
}

func ParseRegisterDevicePayload(body []byte) (*RegisterDevicePayload, error) {
	p := fastjson.Parser{}
	val, err := p.ParseBytes(body)
	if err != nil {
		return nil, err
	}
	return &RegisterDevicePayload{
		UserID: string(val.GetStringBytes("userID")),
		Token:  string(val.GetStringBytes("token")),
	}, err
}
