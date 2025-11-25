package email

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aliakkas006/email-service/models"
	"gopkg.in/gomail.v2"
)

var (
	sentEmails []models.Todo
	mu         sync.Mutex
)

// GetAllEmails returns all sent emails
func GetAllEmails() ([]models.Todo, error) {
	mu.Lock()
	defer mu.Unlock()
	return sentEmails, nil
}

// SendEmail simulates sending an email or sends a real one if configured
func SendEmail(todo models.Todo) error {
	emailTo := todo.UserEmail
	if emailTo == "" {
		emailTo = "aliakkas006@gmail.com"
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpHost != "" && smtpPortStr != "" && smtpEmail != "" && smtpPassword != "" {
		smtpPort, err := strconv.Atoi(smtpPortStr)
		if err != nil {
			return fmt.Errorf("invalid SMTP port: %v", err)
		}

		m := gomail.NewMessage()
		m.SetHeader("From", smtpEmail)
		m.SetHeader("To", emailTo)
		m.SetHeader("Subject", "New Todo Created: "+todo.Title)
		m.SetBody("text/plain", fmt.Sprintf("A new todo has been created:\n\nTitle: %s\nDescription: %s", todo.Title, todo.Description))

		d := gomail.NewDialer(smtpHost, smtpPort, smtpEmail, smtpPassword)

		if err := d.DialAndSend(m); err != nil {
			return fmt.Errorf("failed to send email via SMTP: %v", err)
		}
		log.Printf("📨 Real email sent to %s for Todo: %s", emailTo, todo.Title)
		
		mu.Lock()
		sentEmails = append(sentEmails, todo)
		mu.Unlock()
		
		return nil
	}

	// Simulate delay for sending email
	time.Sleep(1 * time.Second)

	// Print log instead of real email
	log.Printf("📨 [SIMULATION] Email sent to %s for Todo: %s - %s", emailTo, todo.Title, todo.Description)
	
	mu.Lock()
	sentEmails = append(sentEmails, todo)
	mu.Unlock()

	return nil
}
