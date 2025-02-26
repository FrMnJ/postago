package email

type EmailService interface {
	SendEmail(toEmail, subject, template string, data map[string]any) error
}
