package config

import (
	"log"
	//import email channel package
	//import resend provider
	"github.com/r1i2t3/agni/pkg/notification/channels/email"
	"github.com/r1i2t3/agni/pkg/notification/channels/email/EmailProviders"
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
		log.Fatalf("Failed to initialize Resend notifier: %v", err)
	}

	log.Println("✅ Resend channel initialized successfully")
}
