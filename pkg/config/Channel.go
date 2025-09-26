package config

import (
	"log"

	//import email channel package
	//import resend provider
	"github.com/r1i2t3/agni/pkg/notification/channels/email"
	"github.com/r1i2t3/agni/pkg/notification/channels/email/EmailProviders"
	inapp "github.com/r1i2t3/agni/pkg/notification/channels/in-app"
	"github.com/r1i2t3/agni/pkg/notification/channels/webpush"

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

func InitializeWebPushProvider(WebPushEnvConfig *WebPushEnvConfig) {
	if WebPushEnvConfig == nil {
		log.Fatal("WebPush configuration is required")
	}
	log.Printf("Initializing WebPush provider with config: %+v", WebPushEnvConfig)
	_, err := webpush.NewPushNotifier(
		WebPushEnvConfig.VAPID_PUBLIC_KEY,
		WebPushEnvConfig.VAPID_PRIVATE_KEY,
		WebPushEnvConfig.VAPID_SUBJECT,
	)
	if err != nil {
		log.Printf("Failed to initialize WebPush notifier: %v", err)
	}

	log.Println("✅ WebPush channel initialized successfully")
}

func InitializeInAppProvider(InAppConfig *InAppConfig) {
	if InAppConfig == nil {
		log.Fatal("InApp configuration is required")
	}
	inapp.NewInAppNotifier(InAppConfig.stream)
	log.Println("✅ InApp channel initialized successfully")
}
