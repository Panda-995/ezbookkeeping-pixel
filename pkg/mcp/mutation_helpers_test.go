package mcp

import (
	"testing"

	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

func TestMCPAccountCategoryRoundTrip(t *testing.T) {
	categories := []string{
		"cash",
		"checking",
		"credit_card",
		"virtual",
		"debt",
		"receivables",
		"investment",
		"savings",
		"certificate_of_deposit",
	}

	for _, categoryName := range categories {
		category, err := textualAccountCategoryToModel(categoryName)
		if err != nil {
			t.Fatalf("textualAccountCategoryToModel(%q) returned error: %v", categoryName, err)
		}
		if actual := accountCategoryToText(category); actual != categoryName {
			t.Fatalf("account category round trip returned %q, expected %q", actual, categoryName)
		}
	}
}

func TestParseMCPAccountOpeningBalanceForLiability(t *testing.T) {
	balance, balanceTime, err := parseMCPAccountOpeningBalance(
		"123.45",
		"2026-07-27T12:00:00+08:00",
		models.ACCOUNT_CATEGORY_CREDIT_CARD,
	)
	if err != nil {
		t.Fatalf("parseMCPAccountOpeningBalance returned error: %v", err)
	}
	if balance != -12345 {
		t.Fatalf("balance was %d, expected -12345", balance)
	}
	if balanceTime <= 0 {
		t.Fatalf("balance time was %d, expected a positive unix time", balanceTime)
	}
}

func TestMCPTransactionTypeRoundTrip(t *testing.T) {
	transactionTypes := []string{
		transactionTypeModifyBalance,
		transactionTypeIncome,
		transactionTypeExpense,
		transactionTypeTransfer,
	}

	for _, transactionType := range transactionTypes {
		databaseType, err := textualTransactionTypeToDbType(transactionType)
		if err != nil {
			t.Fatalf("textualTransactionTypeToDbType(%q) returned error: %v", transactionType, err)
		}
		if actual := transactionDbTypeToText(databaseType); actual != transactionType {
			t.Fatalf("transaction type round trip returned %q, expected %q", actual, transactionType)
		}
	}
}

func TestMCPHandlersHaveUniqueNames(t *testing.T) {
	if err := InitializeMCPHandlers(&settings.Config{}); err != nil {
		t.Fatalf("InitializeMCPHandlers returned error: %v", err)
	}

	tools := Container.GetMCPTools()
	if len(tools) < 20 {
		t.Fatalf("registered %d MCP tools, expected at least 20", len(tools))
	}

	names := make(map[string]bool, len(tools))
	for _, tool := range tools {
		if names[tool.Name] {
			t.Fatalf("duplicate MCP tool name %q", tool.Name)
		}
		names[tool.Name] = true
	}

	requiredNames := []string{
		"create_account",
		"update_account",
		"adjust_account_balance",
		"delete_account",
		"update_transaction",
		"delete_transaction",
		"create_transaction_category",
		"update_transaction_tag",
		"delete_transaction_tag_group",
	}
	for _, name := range requiredNames {
		if !names[name] {
			t.Fatalf("required MCP tool %q was not registered", name)
		}
	}
}
