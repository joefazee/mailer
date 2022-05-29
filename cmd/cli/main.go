package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joefazee/mailer"
)

func main() {

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	// set your config here
	config := mailer.MailerConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Timeout:  5 * time.Second,
		Sender:   os.Getenv("SMTP_SENDER"),
	}

	sender := mailer.New(config)

	err = sender.Send("jondoe@gmail.com", "welcome.html", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent!")
}
