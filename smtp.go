package remail

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/smtp"
)

type smtpSender struct {
	auth          smtp.Auth
	defaultSender string
	host          string
	port          string
}

func NewSMTPSender(auth smtp.Auth, sender, host, port string) Remail {
	return &smtpSender{
		auth,
		sender,
		host,
		port,
	}
}

func (s smtpSender) Send(ctx context.Context, messages Message) (err error) {
	if messages.From == "" {
		messages.From = s.defaultSender
	}

	if messages.From == "" {
		err = ErrSenderIsRequired
		return
	}

	if len(messages.To) == 0 {
		err = ErrDestinationIsRequired
		return
	}

	withAttachment := len(messages.Attachments) > 0

	buf := bytes.NewBuffer(nil)

	buf.WriteString("MIME-Version: 1.0\n")
	buf.WriteString(fmt.Sprintf("Subject: %s\n", messages.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", buildRecepient(messages.To)))

	if len(messages.Cc) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", buildRecepient(messages.Cc)))
	}

	if len(messages.Bcc) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", buildRecepient(messages.Bcc)))
	}

	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	if withAttachment {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
		buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	}

	buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=\"utf-8\"\r\n", messages.Body.ContentType))
	buf.WriteString("\r\n" + mustBuildBody(messages.Body))

	if withAttachment {
		for _, v := range messages.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v.Body)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", v.Filename))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v.Body)))
			base64.StdEncoding.Encode(b, v.Body)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	err = smtp.SendMail(
		fmt.Sprintf("%s:%s", s.host, s.port),
		s.auth,
		messages.From,
		buildRecepient(messages.To),
		buf.Bytes(),
	)
	return
}
