package providers

type Provider struct {
	Receiver []string
}

func New() *Provider {
	return &Provider{
		Receiver: []string{"tgobp@mailto.plus"},
	}
}

//SMTP
//func (provider *Provider) Authenticate() {
//	//// Sender data.
//	from := "roman.kocenko.2004@gmail.com"
//	password := "vebghujogiettinp"
//
//	//// smtp server configuration.
//	smtpHost := "smtp.gmail.com"
//	smtpPort := "587"
//	//
//	//// Authentication.
//	auth := smtp.PlainAuth("", from, password, smtpHost)
//	//
//	t, _ := template.ParseFiles("templates/message.html")
//
//	var body bytes.Buffer
//
//	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
//	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))
//
//	t.Execute(&body, struct {
//		Name    string
//		Message string
//	}{
//		Name:    "Puneet Singh",
//		Message: "This is a test message in a HTML template",
//	})
//
//	// Sending email.
//	err := smtp.MailHandler(smtpHost+":"+smtpPort, auth, from, provider.Receiver, body.Bytes())
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println("Email Sent!")
//}
