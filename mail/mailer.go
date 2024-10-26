package mail

import (
	"RJD02/job-portal/config"
	"fmt"
	"net/smtp"
	"sync"
)

func GenerateMagicLinkEmail(name, magicLink string) string {
	return fmt.Sprintf(`
        <html>
        <body>
            <h2>Hello %s,</h2>
            <p>You can log in using the link below:</p>
            <a href="%s">Click here to log in</a>
            <br />
            <b>If you're not able to see/click link, please paste below link into your browser, will work fine</b>
            <i>%s</i>
            <p>If you didn’t request this, please ignore this email.</p>
            <br>
            <p>Thank you,</p>
            <p>Your Company Team</p>
        </body>
        </html>
    `, name, magicLink, magicLink)
}

func SendMail(email, htmlMessage, subject string, wg *sync.WaitGroup) {
	defer wg.Done()
	from := config.AppConfig.FROM_GMAIL
	password := config.AppConfig.GMAIL_PASSWORD
	to := config.AppConfig.TO_GMAIL
	if config.AppConfig.ENVIRONMENT == config.Production {
		to = email
	}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, password, smtpHost)

	headers := make(map[string]string)

	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlMessage

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		fmt.Println("Something went wrong", err)
		return
	}
	fmt.Println("Success")
}
