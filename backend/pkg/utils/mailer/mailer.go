package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Mailer interface {
	SendResetPassword(toEmail, token string) error
}

type resendMailer struct {
	apiKey string
	from   string
}

func NewResend(apiKey, from string) Mailer {
	return &resendMailer{apiKey: apiKey, from: from}
}

func (m *resendMailer) SendResetPassword(toEmail, token string) error {
	url := "https://api.resend.com/emails"

	// Di industri, biasanya reset_link mengarah ke web frontend
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), token)

	payload := map[string]interface{}{
		"from":    m.from,
		"to":      []string{toEmail},
		"subject": "Reset Password Anda - School Payment",
		"html":    fmt.Sprintf("<strong>Halo!</strong><p>Anda meminta reset password. Silakan klik link berikut untuk melanjutkan:</p><a href='%s'>Reset Password</a><p>Link ini akan hangus dalam 15 menit.</p>", resetLink),
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("resend error: status code %d", resp.StatusCode)
	}

	return nil
}
