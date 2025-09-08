package smsproviders

import (
	"fmt"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioSender struct {
	Client     *twilio.RestClient
	FromNumber string
}

var TwilioClient *TwilioSender

func NewTwilioSender(fromNumber, accountSID, authToken string) (*TwilioSender, error) {
	if fromNumber == "" || accountSID == "" || authToken == "" {
		return nil, fmt.Errorf("twilio from number, account SID, and auth token are required")
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})
	TwilioClient = &TwilioSender{
		Client:     client,
		FromNumber: fromNumber,
	}
	return TwilioClient, nil
}

func (s *TwilioSender) TwilioSend(to, message string) (string, error) {
	if TwilioClient == nil {
		return "", fmt.Errorf("twilio sender is not initialized")
	}
	params := openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(s.FromNumber)
	params.SetBody(message)

	resp, err := s.Client.Api.CreateMessage(&params)
	if err != nil {
		return "", fmt.Errorf("failed to send SMS via Twilio: %w", err)
	}

	fmt.Printf("SMS sent via Twilio to %s\n", to)

	return *resp.Sid, nil
}
