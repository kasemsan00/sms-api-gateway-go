package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"api-gateway-go/internal/config"
)

var (
	ErrSMSSendFailed = errors.New("SMS send failed")
	ErrSMSDisabled   = errors.New("SMS service is disabled")
)

// SMSService handles SMS sending
type SMSService struct {
	cfg        *config.Config
	httpClient *http.Client
}

// NewSMSService creates a new SMSService
func NewSMSService(cfg *config.Config) *SMSService {
	return &SMSService{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendSMS sends an SMS message
func (s *SMSService) SendSMS(ctx context.Context, phoneNumber, message string) error {
	if !s.cfg.SMSEnable {
		return ErrSMSDisabled
	}

	if s.cfg.SMSAPIURL == "" {
		return ErrSMSDisabled
	}

	// Prepare request body
	body := map[string]string{
		"phoneNumber": phoneNumber,
		"message":     message,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.SMSAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("SMS API error: " + string(body))
	}

	return nil
}

// SendCustomMessage sends a custom SMS message
func (s *SMSService) SendCustomMessage(ctx context.Context, phoneNumber, message string) error {
	return s.SendSMS(ctx, phoneNumber, message)
}

// IsEnabled returns whether SMS is enabled
func (s *SMSService) IsEnabled() bool {
	return s.cfg.SMSEnable && s.cfg.SMSAPIURL != ""
}
