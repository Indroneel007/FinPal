package api

import (
	"encoding/json"
	"examples/SimpleBankProject/prompt"
	"examples/SimpleBankProject/util"
	"fmt"
	"net/http"
	"os"

	//"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/revrost/go-openrouter"
	"github.com/spf13/viper"
)

type PromptRequest struct {
	Mindset string `json:"saving_mindset" binding:"required"` // low, medium, high
}

type PromptResponse struct {
	Prompt string `json:"prompt"`
}

func (server *Server) PromptAPI(c *gin.Context) {
	paths := []string{".env", "../.env", "../../.env"}
	for _, path := range paths {
		_ = godotenv.Load(path)
	}

	viper.AutomaticEnv()
	cwd, _ := os.Getwd()
	fmt.Println("[DEBUG] CWD:", cwd)
	secret := viper.GetString("OPENROUTER_API_KEY")
	fmt.Println("[DEBUG] OPENROUTER_API_KEY from viper:", secret)
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GEMINI_API_KEY is not set in the environment variables"})
		return
	}

	client := openrouter.NewClient(secret, openrouter.WithXTitle("FinPal"), openrouter.WithHTTPReferer("http://finpal.com"))

	var req PromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	payloadData, exists := c.Get(authorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization payload not found"})
		return
	}

	payload, ok := payloadData.(*util.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid authorization payload"})
		return
	}

	// Get total expenses by type for this owner
	expenses, err := server.store.GetTotalByOwnerAndType(c, payload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed: " + err.Error()})
		return
	}

	// Initialize total salary and expenses struct
	var totalSalary float64
	userExpenses := prompt.ExpenseCategories{}

	for _, exp := range expenses {
		switch exp.Type {
		case "salary":
			totalSalary += float64(exp.TotalBalance)
		case "rent":
			userExpenses.Rent = int64(exp.TotalBalance)
		case "food":
			userExpenses.Food = int64(exp.TotalBalance)
		case "travel":
			userExpenses.Travel = int64(exp.TotalBalance)
		case "savings":
			userExpenses.Savings = int64(exp.TotalBalance)
		case "bills":
			userExpenses.Bills = int64(exp.TotalBalance)
		case "medical":
			userExpenses.Medical = int64(exp.TotalBalance)
		case "shopping":
			userExpenses.Shopping = int64(exp.TotalBalance)
		}
	}

	location, err := server.store.GetLocationByUsername(c, payload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed: " + err.Error()})
		return
	}

	getUser, err := server.store.GetUser(c, payload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed: " + err.Error()})
		return
	}

	existsInRedis, err := server.redis.Exists(c, util.PromptStorePrefix+location.Address).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "redis query failed: " + err.Error()})
		return
	}

	var standardCosts prompt.ExpenseCategories

	if existsInRedis == 1 {
		redisData, err := server.redis.Get(c, util.PromptStorePrefix+location.Address).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "redis query failed: " + err.Error()})
			return
		}

		if err := json.Unmarshal([]byte(redisData), &standardCosts); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unmarshal failed: " + err.Error()})
			return
		}
	} else {
		standardCosts, err = prompt.GetCostOfLivingFromAI(c, location.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Prompt query failed: " + err.Error()})
			return
		}
	}

	// Compose prompt
	promptData := prompt.PromptData{
		Location:      location.Address, // Change if dynamic
		Salary:        getUser.Salary,
		Mindset:       req.Mindset,
		Expenses:      userExpenses,
		StandardCosts: standardCosts,
	}

	promptText := prompt.BuildPrompt(promptData)

	resp, err := client.CreateChatCompletion(c, openrouter.ChatCompletionRequest{
		Model: openrouter.DeepseekV3,
		Messages: []openrouter.ChatCompletionMessage{
			{
				Role:    openrouter.ChatMessageRoleUser,
				Content: openrouter.Content{Text: promptText},
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Prompt query failed: " + err.Error()})
		return
	}

	textPart := resp.Choices[0].Message.Content

	/*rawText := strings.TrimSpace(textPart.Text)

	//fmt.Println(rawText)

	categories, err := prompt.ParseExpenseCategories(rawText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Prompt query failed: " + err.Error()})
		return
	}*/

	c.JSON(http.StatusOK, PromptResponse{
		Prompt: textPart.Text,
	})
}
