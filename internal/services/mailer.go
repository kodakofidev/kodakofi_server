package services

import (
	"fmt"
	"net/smtp"
	"os"
)

type MailerService interface {
	SendVerificationEmail(recipientEmail string, otp string) error
	SendOTPEmail(recipientEmail string, otp string, otpType int) error
}

type GmailSMTP struct {
	SMTPHost    string
	SMTPPort    string
	SenderEmail string
	SenderPass  string
	SenderName  string
}

func NewGmailMailer() *GmailSMTP {
	return &GmailSMTP{
		SMTPHost:    "smtp.gmail.com",
		SMTPPort:    "587",
		SenderEmail: os.Getenv("GMAIL_EMAIL"),
		SenderPass:  os.Getenv("GMAIL_APP_PASSWORD"),
		SenderName:  "KodaKofi",
	}
}

func (m *GmailSMTP) SendVerificationEmail(recipientEmail string, otp string) error {
	return m.SendOTPEmail(recipientEmail, otp, 1) // Type 1 is email verification
}

func (m *GmailSMTP) SendOTPEmail(recipientEmail string, otp string, otpType int) error {
	var subject string
	var body string

	switch otpType {
	case 1: // Email verification
		subject = "KodaKofi Email Verification"
		body = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Email Verification</title>
</head>
<body>
    <h1>Email Verification</h1>
    <p>Thank you for registering with KodaKofi. Please use the following verification code to complete your registration:</p>
    <h2>%s</h2>
    <p>This code will expire in 15 minutes.</p>
    <p>If you did not request this verification, please ignore this email.</p>
</body>
</html>`, otp)
	case 2: // Password reset
		subject = "KodaKofi Password Reset"
		body = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Password Reset</title>
</head>
<body>
    <h1>Password Reset Request</h1>
    <p>We received a request to reset your password. Please use the following verification code to reset your password:</p>
    <h2>%s</h2>
    <p>This code will expire in 15 minutes.</p>
    <p>If you did not request a password reset, please ignore this email.</p>
</body>
</html>`, otp)
	default:
		subject = "KodaKofi OTP Code"
		body = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>OTP Code</title>
</head>
<body>
    <h1>Your OTP Code</h1>
    <p>You requested an OTP code. Please use the following code:</p>
    <h2>%s</h2>
    <p>This code will expire in 15 minutes.</p>
    <p>If you did not request this code, please ignore this email.</p>
</body>
</html>`, otp)
	}

	// Format the MIME message
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := fmt.Sprintf("Subject: %s\r\n%s\r\n%s", subject, mime, body)

	// Set up authentication information
	auth := smtp.PlainAuth("", m.SenderEmail, m.SenderPass, m.SMTPHost)

	// Send the email
	addr := fmt.Sprintf("%s:%s", m.SMTPHost, m.SMTPPort)
	err := smtp.SendMail(addr, auth, m.SenderEmail, []string{recipientEmail}, []byte(message))
	return err
}
