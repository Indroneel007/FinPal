package prompt

import (
	"fmt"
	"strings"
)

type ExpenseCategories struct {
	Rent     int64 `json:"rent"`
	Food     int64 `json:"food"`
	Travel   int64 `json:"travel"`
	Savings  int64 `json:"savings"`
	Bills    int64 `json:"bills"`
	Medical  int64 `json:"medical"`
	Shopping int64 `json:"shopping"`
	Misc     int64 `json:"misc"`
}

type PromptData struct {
	Location      string            `json:"location"`
	Salary        int64             `json:"salary"`
	Mindset       string            `json:"saving_mindset"` // low, medium, high
	Expenses      ExpenseCategories `json:"expenses"`
	StandardCosts ExpenseCategories `json:"standard_costs"` // from API or static
}

func BuildPrompt(data PromptData) string {
	var sb strings.Builder

	sb.WriteString("You are a financial advisor. Your task is to analyze the user's budget and recommend where they can cut costs or are doing well.\n\n")

	sb.WriteString(fmt.Sprintf("Location: %s\n", data.Location))
	sb.WriteString(fmt.Sprintf("Monthly Salary: ₹%d\n", data.Salary))
	sb.WriteString(fmt.Sprintf("Saving Mindset: %s\n\n", data.Mindset))

	sb.WriteString("User's Expenses:\n")
	sb.WriteString(fmt.Sprintf("- Rent: ₹%d\n", data.Expenses.Rent))
	sb.WriteString(fmt.Sprintf("- Food: ₹%d\n", data.Expenses.Food))
	sb.WriteString(fmt.Sprintf("- Travel: ₹%d\n", data.Expenses.Travel))
	sb.WriteString(fmt.Sprintf("- Savings: ₹%d\n", data.Expenses.Savings))
	sb.WriteString(fmt.Sprintf("- Bills: ₹%d\n", data.Expenses.Bills))
	sb.WriteString(fmt.Sprintf("- Medical: ₹%d\n", data.Expenses.Medical))
	sb.WriteString(fmt.Sprintf("- Shopping: ₹%d\n\n", data.Expenses.Shopping))

	sb.WriteString("Standard Cost of Living Benchmarks in this location:\n")
	sb.WriteString(fmt.Sprintf("- Rent: ₹%d\n", data.StandardCosts.Rent))
	sb.WriteString(fmt.Sprintf("- Food: ₹%d\n", data.StandardCosts.Food))
	sb.WriteString(fmt.Sprintf("- Travel: ₹%d\n", data.StandardCosts.Travel))
	sb.WriteString(fmt.Sprintf("- Savings: ₹%d\n", data.StandardCosts.Savings))
	sb.WriteString(fmt.Sprintf("- Bills: ₹%d\n", data.StandardCosts.Bills))
	sb.WriteString(fmt.Sprintf("- Medical: ₹%d\n", data.StandardCosts.Medical))
	sb.WriteString(fmt.Sprintf("- Shopping: ₹%d\n\n", data.StandardCosts.Shopping))

	sb.WriteString("Based on the above, compare each category with the benchmark adjusted to the saving mindset. Dont include any extra information or explanations. Write in less than 200 words. Dont include any extra information or explanations. Only contain english letters and numbers, no symbols. \n")
	//sb.WriteString("Give numerical feedback for each category, and suggestions on where the user can reduce expenses.\n")
	//sb.WriteString("Dont include any extra information or explanations. Write in less than 200 words.\n")
	//sb.WriteString("Beatuify the output with markdown. No irrelavant information. \n")
	//sb.WriteString("End with a recommendation summary.\n")

	return sb.String()
}
