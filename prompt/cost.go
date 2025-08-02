package prompt

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/revrost/go-openrouter"
	"github.com/spf13/viper"
	//"google.golang.org/api/option"
)

func GetCostOfLivingFromAI(ctx context.Context, location string) (ExpenseCategories, error) {
	paths := []string{".env", "../.env", "../../.env"}
	for _, path := range paths {
		_ = godotenv.Load(path)
	}

	viper.AutomaticEnv()
	secret := viper.GetString("OPENROUTER_API_KEY")
	if secret == "" {
		return ExpenseCategories{}, fmt.Errorf("GEMINI_API_KEY is not set in the environment variables")
	}

	client := openrouter.NewClient(secret, openrouter.WithXTitle("FinPal"), openrouter.WithHTTPReferer("http://finpal.com"))

	prompt := fmt.Sprintf(
		"Give me a rough monthly cost of living in %s, India with categories and only INR numbers (no commas or currency symbols). Categories should be: Rent, Food, Travel, Savings, Bills, Medical, Shopping. Reply in this format:\nRent: 15000\nFood: 5000\nTravel: 2000\nSavings: 5000\nBills: 3000\nMedical: 1000\nShopping: 2000",
		location,
	)

	resp, err := client.CreateChatCompletion(ctx, openrouter.ChatCompletionRequest{
		Model: openrouter.DeepseekV3,
		Messages: []openrouter.ChatCompletionMessage{
			{
				Role:    openrouter.ChatMessageRoleUser,
				Content: openrouter.Content{Text: prompt},
			},
		},
	})
	if err != nil {
		return ExpenseCategories{}, fmt.Errorf("gemini API error: %v", err)
	}

	//textPart, ok := resp.Candidates[0].Content.Parts[0].(*genai.Text)
	textPart := resp.Choices[0].Message.Content

	rawText := strings.TrimSpace(textPart.Text)

	//fmt.Println(rawText)

	categories, err := ParseExpenseCategories(rawText)
	if err != nil {
		fmt.Println(err)
		return ExpenseCategories{}, err
	}

	fmt.Printf("Categories %+v\n", categories)

	return categories, err
}

func ParseExpenseCategories(text string) (ExpenseCategories, error) {
	lines := strings.Split(text, "\n")
	categories := ExpenseCategories{}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and lines without a colon
		if line == "" || !strings.Contains(line, ":") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue // skip invalid lines
		}

		key := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])

		if valueStr == "" {
			continue
		}

		// Remove any commas, ₹, etc.
		valueStr = strings.ReplaceAll(valueStr, "₹", "")
		valueStr = strings.ReplaceAll(valueStr, ",", "")

		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			continue // skip lines that can't be parsed
		}

		switch strings.ToLower(key) {
		case "rent":
			categories.Rent = value
		case "food":
			categories.Food = value
		case "travel":
			categories.Travel = value
		case "savings":
			categories.Savings = value
		case "bills":
			categories.Bills = value
		case "medical":
			categories.Medical = value
		case "shopping":
			categories.Shopping = value
		default:
			// Unknown category — skip
		}
	}

	return categories, nil
}
