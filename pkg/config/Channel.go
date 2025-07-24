package config

import (
	"log"
	//import email channel package

	"github.com/r1i2t3/agni/pkg/notification/channels/email"
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
