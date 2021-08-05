# Remail

## Installation 

`go get -u github.com/rahman-teja/remail`

## Usage

`import "github.com/rahman-teja/remail"`

### Example
#### SMTP
```
// generate auth
sender := "<your-email-account>"
password := "<your-email-password>"
host := "<your-email-host>"
port := "<your-email-port>"

auth := smtp.PlainAuth(
	"",
	sender,
	password,
	host,
)

message := remail.Message{
	Subject: "<subject>",
	To: []remail.Recepient{
		{
			Name:    "",
			Address: "",
		},
	},
	Body: remail.MessageBody{
		ContentType: remail.ContentTypeHTML,
		Body:        []byte("my body"),
	},
}

mailer := remail.NewSMTPSender(auth, sender, host, port)
err := mailer.Send(context.Background(), message)
if err != nil {
	log.Fatal(err)
	return
}

fmt.Println("Email sent")
```

#### AWS Ses
```
region := ""
awsCharset := ""
config := &aws.Config{}

message := remail.Message{
	Subject: "<subject>",
	To: []remail.Recepient{
		{
			Name:    "",
			Address: "",
		},
	},
	Body: remail.MessageBody{
		ContentType: remail.ContentTypeHTML,
		Body:        []byte("my body"),
	},
}

mailer := remail.NewAWSSes(region, awsCharset, sender, config)
err := mailer.Send(context.Background(), message)
if err != nil {
	log.Fatal(err)
	return
}

fmt.Println("Email sent")
```