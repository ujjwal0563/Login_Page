package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type brevoEmailPayload struct {
	Sender struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"sender"`
	To []struct {
		Email string `json:"email"`
	} `json:"to"`
	Subject     string `json:"subject"`
	HTMLContent string `json:"htmlContent"`
}

func SendOTPEmailBrevo(toEmail string, otp string) error {
	apiKey := os.Getenv("BREVO_API_KEY")

	payload := brevoEmailPayload{}
	payload.Sender.Email = os.Getenv("BREVO_FROM_EMAIL")
	payload.Sender.Name = os.Getenv("BREVO_FROM_NAME")
	payload.Subject = "Password Reset OTP"

	payload.To = append(payload.To, struct {
		Email string `json:"email"`
	}{Email: toEmail})

	payload.HTMLContent = `
		<h2>Password Reset OTP</h2>
		<p>Your OTP is:</p>
		<h1>` + otp + `</h1>
		<p>This OTP is valid for <b>5 minutes</b>.</p>
	`

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		"POST",
		"https://api.brevo.com/v3/smtp/email",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", apiKey)
	req.Header.Set("content-type", "application/json")

	_, err = http.DefaultClient.Do(req)
	return err
}
