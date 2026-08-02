package mcp

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mayswind/ezbookkeeping/pkg/models"
)

func TestUpdateAccountContractIncludesAllMutableFields(t *testing.T) {
	requestType := reflect.TypeOf(MCPUpdateAccountRequest{})
	requiredJSONFields := []string{
		"id", "account_name", "name", "category", "icon", "color", "currency", "balance",
		"comment", "hidden", "credit_card_statement_date", "last_reconciled_time",
		"clear_last_reconciled_time", "dry_run",
	}

	assertStructHasJSONFields(t, requestType, requiredJSONFields)
}

func TestUpdateTransactionContractIncludesAllMutableFields(t *testing.T) {
	requestType := reflect.TypeOf(MCPUpdateTransactionRequest{})
	requiredJSONFields := []string{
		"id", "type", "time", "category_name", "account_name", "amount",
		"destination_account_name", "destination_amount", "tags", "picture_ids",
		"comment", "hide_amount", "geo_location", "clear_geo_location", "dry_run",
	}

	assertStructHasJSONFields(t, requestType, requiredJSONFields)
}

func TestAccountModelsEqualForMCPIncludesCurrencyAndBalance(t *testing.T) {
	base := &models.Account{
		Name:     "Cash",
		Category: models.ACCOUNT_CATEGORY_CASH,
		Icon:     1,
		Color:    "176b5b",
		Currency: "CNY",
		Balance:  100,
		Extend:   &models.AccountExtend{},
	}

	currencyChanged := *base
	currencyChanged.Currency = "USD"
	if accountModelsEqualForMCP(base, &currencyChanged) {
		t.Fatal("accounts with different currencies must not be equal")
	}

	balanceChanged := *base
	balanceChanged.Balance = 200
	if accountModelsEqualForMCP(base, &balanceChanged) {
		t.Fatal("accounts with different balances must not be equal")
	}
}

func assertStructHasJSONFields(t *testing.T, requestType reflect.Type, requiredFields []string) {
	t.Helper()

	actualFields := make(map[string]bool, requestType.NumField())
	for index := 0; index < requestType.NumField(); index++ {
		jsonName := strings.Split(requestType.Field(index).Tag.Get("json"), ",")[0]
		actualFields[jsonName] = true
	}

	for _, field := range requiredFields {
		if !actualFields[field] {
			t.Fatalf("%s is missing JSON field %q", requestType.Name(), field)
		}
	}
}
