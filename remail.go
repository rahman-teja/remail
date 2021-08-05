package remail

import (
	"context"
)

type Recepient struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type MessageBody struct {
	ContentType ContentType
	Body        []byte
}

type Attachment struct {
	Filename string
	Body     []byte
}

type Message struct {
	From        string
	Subject     string
	To          []Recepient
	Cc          []Recepient
	Bcc         []Recepient
	Body        MessageBody
	Attachments []Attachment
}

type Remail interface {
	Send(ctx context.Context, messages Message) (err error)
}
