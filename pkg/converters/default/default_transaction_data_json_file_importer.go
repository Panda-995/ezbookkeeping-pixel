package _default

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var allJsonDataSupportedColumns = []datatable.TransactionDataTableColumn{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME,
	datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT,
	datatable.TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION,
	datatable.TRANSACTION_DATA_TABLE_TAGS,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION,
}

// defaultTransactionDataJsonImporter defines the structure of ezbookkeeping default json importer for transaction data
type defaultTransactionDataJsonImporter struct{}

type compatibleJsonString string

type compatibleJsonImportTransactionRequest struct {
	Transactions []*compatibleJsonImportTransactionRequestItem `json:"transactions"`
}

type compatibleJsonImportTransactionRequestItem struct {
	Time                   compatibleJsonString `json:"time"`
	UtcOffset              compatibleJsonString `json:"utcOffset"`
	Type                   compatibleJsonString `json:"type"`
	CategoryName           compatibleJsonString `json:"categoryName,omitempty"`
	SourceAccountName      compatibleJsonString `json:"sourceAccountName,omitempty"`
	DestinationAccountName compatibleJsonString `json:"destinationAccountName,omitempty"`
	SourceAmount           compatibleJsonString `json:"sourceAmount"`
	DestinationAmount      compatibleJsonString `json:"destinationAmount,omitempty"`
	GeoLocation            compatibleJsonString `json:"geoLocation,omitempty"`
	TagNames               compatibleJsonString `json:"tagNames,omitempty"`
	Comment                compatibleJsonString `json:"comment,omitempty"`
}

// Initialize an ezbookkeeping default transaction data json file importer singleton instance
var (
	DefaultTransactionDataJsonFileImporter = &defaultTransactionDataJsonImporter{}
)

// ParseImportedData returns the imported data by parsing the transaction json data
func (c *defaultTransactionDataJsonImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezone *time.Location, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	var importRequest compatibleJsonImportTransactionRequest

	if err := json.Unmarshal(data, &importRequest); err != nil {
		return nil, nil, nil, nil, nil, nil, errs.ErrInvalidJSONFile
	}

	transactionDataTable, err := c.createNewDefaultTransactionDataTable(importRequest)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTableImporter := converter.CreateNewImporterWithTypeNameMapping(
		ezbookkeepingTransactionTypeNameMapping,
		ezbookkeepingGeoLocationSeparator,
		ezbookkeepingGeoLocationOrder,
		ezbookkeepingTagSeparator,
	)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezone, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}

func (c *defaultTransactionDataJsonImporter) createNewDefaultTransactionDataTable(importRequest compatibleJsonImportTransactionRequest) (datatable.TransactionDataTable, error) {
	transactionDataTable := datatable.CreateNewWritableTransactionDataTable(allJsonDataSupportedColumns)

	if importRequest.Transactions == nil || len(importRequest.Transactions) < 1 {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	for i := 0; i < len(importRequest.Transactions); i++ {
		transaction := importRequest.Transactions[i]

		utcOffsetText := string(transaction.UtcOffset)
		utcOffset, err := utils.StringToInt(utcOffsetText)

		if err != nil {
			return nil, errs.ErrTransactionTimeZoneInvalid
		}

		timezone := time.FixedZone("Transaction Timezone", utcOffset*60)
		timezoneOffset := utils.FormatTimezoneOffset(time.Now().Unix(), timezone)

		row := make(map[datatable.TransactionDataTableColumn]string, len(allJsonDataSupportedColumns))
		row[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = normalizeEzbookkeepingTransactionTime(string(transaction.Time), utcOffsetText)
		row[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE] = timezoneOffset
		row[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = string(transaction.Type)
		row[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = string(transaction.CategoryName)
		row[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = string(transaction.SourceAccountName)
		row[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = string(transaction.SourceAmount)
		row[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = string(transaction.DestinationAccountName)
		row[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = string(transaction.DestinationAmount)
		row[datatable.TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION] = string(transaction.GeoLocation)
		row[datatable.TRANSACTION_DATA_TABLE_TAGS] = string(transaction.TagNames)
		row[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = string(transaction.Comment)

		transactionDataTable.Add(row)
	}

	return transactionDataTable, nil
}

func (s *compatibleJsonString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		*s = ""
		return nil
	}

	if len(data) > 0 && data[0] == '"' {
		var value string

		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}

		*s = compatibleJsonString(value)
		return nil
	}

	var value json.Number

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*s = compatibleJsonString(value.String())
	return nil
}
