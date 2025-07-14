package util

import (
	"context"
	"crypto/rand"
	"examples/SimpleBankProject/config"
	"fmt"
	"net/smtp"

	//"fmt"
	"log"
	"math/big"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	//"github.com/subosito/gotenv"
	//"golang.org/x/crypto/bcrypt"
)

func GenerateOTP() string {
	result := make([]byte, OTPLength)
	charsetLength := big.NewInt(int64(len(otpCharSet)))

	for i := range result {
		// generate a secure random number in the range of the charset length
		num, _ := rand.Int(rand.Reader, charsetLength)
		result[i] = otpCharSet[num.Int64()]
	}

	return string(result)
}

func AddOTPToRedis(otp string, email string, c context.Context) error {
	key := otpKeyPrefix + email
	hashedOTP, err := HashPassword(otp)
	if err != nil {
		log.Printf("Error hashing OTP: %v", err)
		return err
	}

	err = config.Redis.Client.Set(c, key, hashedOTP, otpExp).Err()
	if err != nil {
		log.Printf("Error setting OTP in Redis: %v", err)
		return err
	}

	return nil
}

func SendOTPEmail(otp, recepient string) error {
	/*err := godotenv.Load()
	if err != nil {
		return err
	}

	viper.AutomaticEnv()

	sender := viper.GetString("SMTP_EMAIL")
	if sender == "" {
		log.Println("SMTP_EMAIL is not set in the environment variables")
		return err
	}

	client := config.SMTPClient

	// set recipient
	err = client.Rcpt(recepient)
	if err != nil {
		fmt.Printf("Error setting recipient: %v\n", err)
		return err
	}

	// start writing email
	writeCloser, err := client.Data()
	if err != nil {
		return err
	}

	// contents of the email
	msg := fmt.Sprintf(emailTemplate, recepient, otp)

	// write the email
	_, err = writeCloser.Write([]byte(msg))
	if err != nil {
		return err
	}

	// close writecloser and send email
	err = writeCloser.Close()
	if err != nil {
		return err
	}

	return nil*/

	/*err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env:", err)
		return err
	}*/

	paths := []string{".env", "../.env", "../../.env"}
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			break
		}
	}

	var err error

	viper.AutomaticEnv()

	smtpHost := viper.GetString("SMTP_HOST") // e.g., smtp.gmail.com
	smtpPort := viper.GetString("SMTP_PORT") // e.g., 587
	smtpEmail := viper.GetString("SMTP_EMAIL")
	smtpPassword := viper.GetString("SMTP_PASSWORD")

	if smtpHost == "" || smtpPort == "" || smtpEmail == "" || smtpPassword == "" {
		return fmt.Errorf("missing SMTP credentials")
	}

	auth := smtp.PlainAuth("", smtpEmail, smtpPassword, smtpHost)

	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP code is: %s", otp)

	msg := []byte("To: " + recepient + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpEmail, []string{recepient}, msg)
	if err != nil {
		log.Println("Error sending mail:", err)
		return err
	}

	log.Println("OTP email sent to:", recepient)
	return nil
}
