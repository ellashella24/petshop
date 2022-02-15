package service

import (
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"petshop/constants"
)

func SendMail(file, emailUser string) error {
	server := mail.NewSMTPClient()
	server.Host = constants.CONFIG_SMTP_HOST
	server.Port = 587
	server.Username = constants.CONFIG_AUTH_EMAIL
	server.Password = constants.CONFIG_AUTH_PASSWORD

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Create email
	email := mail.NewMSG()
	email.SetFrom(constants.CONFIG_SENDER_NAME)
	email.AddTo(emailUser)
	email.SetSubject("Data Transaksi Produk 1")
	email.AddAttachment(file)

	// Send email
	err = email.Send(smtpClient)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
