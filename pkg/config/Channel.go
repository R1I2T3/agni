package config

import (
	"fmt"
	"log"

	//import email channel package
	//import resend provider
	"github.com/r1i2t3/agni/pkg/notification/channels/email"
	"github.com/r1i2t3/agni/pkg/notification/channels/email/EmailProviders"

	// SMS provider
	smsproviders "github.com/r1i2t3/agni/pkg/notification/channels/sms/SMSProviders"
)

func InitializeEmailChannel(EmailEnvConfig *EmailEnvConfig) {
	if EmailEnvConfig == nil {
		log.Fatal("Email configuration is required")
	}
	log.Printf("Initializing email channel with config: %+v", EmailEnvConfig)
	err := email.NewEmailNotifier(
		EmailEnvConfig.SMTPHost,
		EmailEnvConfig.SMTPPort,
		EmailEnvConfig.SMTPUsername,
		EmailEnvConfig.SMTPPassword,
	)
	if err != nil {
		log.Fatalf("Failed to initialize email notifier: %v", err)
	}

	log.Println("✅ Email channel initialized successfully")
}

func InitializeResendProvider(ResendEnvConfig *ResendEnvConfig) {
	if ResendEnvConfig == nil {
		log.Fatal("Resend configuration is required")
	}
	log.Printf("Initializing Resend provider with config: %+v", ResendEnvConfig)
	err := EmailProviders.NewResendNotifier(
		ResendEnvConfig.APIKey,
		ResendEnvConfig.FromAddress,
	)
	if err != nil {
		log.Printf("Failed to initialize Resend notifier: %v", err)
	}

	log.Println("✅ Resend channel initialized successfully")
}

func InitializeTwilioProvider(TwilioEnvConfig *TwilioEnvConfig) {
	fmt.Println(TwilioEnvConfig)
	if TwilioEnvConfig == nil {
		log.Fatal("Twilio configuration is required")
	}
	log.Printf("Initializing Twilio provider with config: %+v", TwilioEnvConfig)
	_, err := smsproviders.NewTwilioSender(
		TwilioEnvConfig.TWILIO_PHONE_NUMBER,
		TwilioEnvConfig.ACCOUNT_SID,
		TwilioEnvConfig.AUTH_TOKEN,
	)
	if err != nil {
		log.Printf("Failed to initialize Twilio notifier: %v", err)
	}

	log.Println("✅ Twilio channel initialized successfully")
}
