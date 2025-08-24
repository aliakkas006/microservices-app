package email

import (
	"log"
	"time"

	"github.com/aliakkas006/email-service/models"
)

// SendEmail simulates sending an email
func SendEmail(todo models.Todo) error {
	email := todo.UserEmail
	if email == "" {
		email = "aliakkas006@gmail.com"
	}

	// email := "aliakkas006@gmail.com"

	// Simulate delay for sending email
	time.Sleep(1 * time.Second)

	// Print log instead of real email
	log.Printf("📨 Email sent to the %s for Todo: %s - %s", email, todo.Title, todo.Description)
	return nil
}
