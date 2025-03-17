package service

import (
	"net/smtp"
	"os"
	"sync"

	"vote/app/utils"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func Send(title string, body string, to string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	envErr := godotenv.Load()
	if envErr != nil {
		panic(envErr)
	}

	from := os.Getenv("SMTP_FROM")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	port := os.Getenv("SMTP_PORT")
	host := os.Getenv("SMTP_HOST")

	msg := "From: " + from + "\n" +
		   "To: " + to + "\n" +
		   "Subject: " + title + "\n" +
		   body

	err := smtp.SendMail(host + ":" + port,
		   smtp.PlainAuth("", username, password, host),
		   from, []string{to}, []byte(msg))
		   
	if err != nil{
		utils.Logger().WithFields(logrus.Fields{
			"name": "Smtp",
		}).Error("error: ", err)
		return
	}

	utils.Logger().WithFields(logrus.Fields{
		"name": "Smtp",
	}).Info("Send from: ", from + ", To: ", to)
}

func MultiSend(email string) {
	var wg sync.WaitGroup
	wg.Add(2)
	go Send("Register Notification", "Welcome to become our membership", email, &wg)
	go Send("Please review the rules", "Rules1:..........", email, &wg)
	wg.Wait()
	utils.Logger().WithFields(logrus.Fields{
		"name": "Smtp",
	}).Info("Finished all tasks")
}