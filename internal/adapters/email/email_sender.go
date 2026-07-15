package email

import "net/smtp"

type EmailSender interface {
	SendVerificationCode(email, code string) error
}

type SmtpEmailSender struct {
	address  string
	port     string
	login    string
	password string
}

func NewSmtpEmailSender(address string, port string, login string, password string) *SmtpEmailSender {
	return &SmtpEmailSender{
		address:  address,
		port:     port,
		login:    login,
		password: password,
	}
}

func (s *SmtpEmailSender) SendVerificationCode(email, code string) error {
	subject := "TEXT"
	body := code
	msg := []byte(subject + "\n" + body)

	err := smtp.SendMail(s.address+":"+s.port, smtp.PlainAuth("", s.login, s.password, s.address), s.login, []string{email}, msg)
	return err
}
