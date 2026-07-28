package _default

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestDefaultTransactionDataJsonFileImporter_NormalizesCompatibleTransactionTime(t *testing.T) {
	context := core.NewNullContext()
	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	jsonData := []byte(`{
		"transactions": [{
			"time": " 2024/9/1 12:34 ",
			"utcOffset": "480",
			"type": "Expense",
			"categoryName": "Dining",
			"sourceAccountName": "Wallet",
			"sourceAmount": "12.34"
		}]
	}`)

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, _, _, _, err :=
		DefaultTransactionDataJsonFileImporter.ParseImportedData(
			context,
			user,
			jsonData,
			time.UTC,
			converter.DefaultImporterOptions,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

	if !assert.NoError(t, err) ||
		!assert.Len(t, allNewTransactions, 1) ||
		!assert.Len(t, allNewAccounts, 1) ||
		!assert.Len(t, allNewSubExpenseCategories, 1) {
		return
	}

	transactionTimezone := time.FixedZone("UTC+8", 8*60*60)
	assert.Equal(
		t,
		"2024-09-01 12:34:00",
		utils.FormatUnixTimeToLongDateTime(
			utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime),
			transactionTimezone,
		),
	)
}

func TestDefaultTransactionDataJsonFileImporter_AcceptsNumericOriginalValues(t *testing.T) {
	context := core.NewNullContext()
	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}
	unixMilliseconds := time.Date(2024, time.September, 1, 12, 34, 56, 0, time.UTC).UnixMilli()
	jsonData := []byte(`{
		"transactions": [{
			"time": ` + utils.Int64ToString(unixMilliseconds) + `,
			"utcOffset": 480,
			"type": "Expense",
			"categoryName": "Dining",
			"sourceAccountName": "Wallet",
			"sourceAmount": 12.34
		}]
	}`)

	allNewTransactions, _, _, _, _, _, err := DefaultTransactionDataJsonFileImporter.ParseImportedData(
		context,
		user,
		jsonData,
		time.UTC,
		converter.DefaultImporterOptions,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	if !assert.NoError(t, err) || !assert.Len(t, allNewTransactions, 1) {
		return
	}

	transactionTimezone := time.FixedZone("UTC+8", 8*60*60)
	assert.Equal(
		t,
		"2024-09-01 20:34:56",
		utils.FormatUnixTimeToLongDateTime(
			utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime),
			transactionTimezone,
		),
	)
	assert.Equal(t, int64(1234), allNewTransactions[0].Amount)
}
