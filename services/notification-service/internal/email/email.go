package email

import (
	"fmt"
	"net/smtp"
)

// SendEmail sends a simple email using an SMTP server.
func SendEmail(to, subject, body string) error {
	// Replace with your SMTP server details
	smtpHost := "smtp.your-email-provider.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", "your-email@example.com", "your-email-password", smtpHost)

	// Format the email message
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, "your-email@example.com", []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
