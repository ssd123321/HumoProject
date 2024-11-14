package smtp

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"net/smtp"
	"os"
)

type Sender struct {
	email    string
	password string
}

var Sen = Sender{
	email:    "frosta123456@gmail.com",
	password: "vhmv atbf qudt wxmk\n\n",
}

type Receiver struct {
	emails []string
}

type ServerSMTP struct {
	Host string
	Port string
}

var Ser = ServerSMTP{
	Host: "smtp.gmail.com",
	Port: "587",
}

func SendFile(s *Sender, email string, ser *ServerSMTP, file *os.File) error {
	var r = Receiver{[]string{
		email,
	}}
	message := gomail.NewMessage()
	message.SetHeader("To", r.emails[0])
	message.SetHeader("From", s.email)
	message.SetHeader("Subject", "Email with file attached")
	message.SetBody("text/plain", "Find and read the attached file")
	message.Attach(file.Name())
	d := gomail.NewDialer("smtp.gmail.com", 25, "frosta123456@gmail.com", "vhmv atbf qudt wxmk\n\n")
	err := d.DialAndSend(message)
	if err != nil {
		return err
	}
	return nil
}

func SendMessage(s *Sender, email string, ser *ServerSMTP, msg []byte) error {
	var r = Receiver{[]string{
		email,
	}}
	msg1 := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: Person!\r\n"+
			"\r\n"+
			"%s.\r\n", email, string(msg)))
	auth := smtp.PlainAuth("", s.email, s.password, ser.Host)
	fmt.Println("3", auth)
	err := smtp.SendMail(ser.Host+":"+ser.Port, auth, s.email, r.emails, msg1)
	fmt.Println("4")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
