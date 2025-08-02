package config

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	//"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	//"github.com/subosito/gotenv"
)

var SMTPClient *smtp.Client

func SMTPConnection() {
	/*err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file in SMTP connection:", err)
		return
	}*/

	paths := []string{".env", "../.env", "../../.env"}
	for _, path := range paths {
		_ = godotenv.Load(path)
	}

	viper.AutomaticEnv()
	host := viper.GetString("SMTP_HOST")
	port := viper.GetString("SMTP_PORT")
	email := viper.GetString("SMTP_EMAIL")
	password := viper.GetString("SMTP_PASSWORD")

	/*host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	email := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")*/

	if host == "" || port == "" || email == "" || password == "" {
		fmt.Println("SMTP configuration is not set in the environment variables")
	}

	smtpAuth := smtp.PlainAuth("", email, password, host)

	client, err := smtp.Dial(host + ":" + port)

	if err != nil {
		panic(err)
	}

	SMTPClient = client
	client = nil

	// initiate TLS handshake
	if ok, _ := SMTPClient.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: host}
		if err = SMTPClient.StartTLS(config); err != nil {
			panic(err)
		}
	}

	// authenticate
	err = SMTPClient.Auth(smtpAuth)
	if err != nil {
		panic(err)
	}

	fmt.Println("SMTP Connected")

}
