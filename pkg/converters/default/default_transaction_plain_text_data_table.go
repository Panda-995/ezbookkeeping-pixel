package _default

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

const (
	legacyEzbookkeepingDescriptionColumnName   = "Comment"
	currentEzbookkeepingDescriptionColumnName  = "Description"
	ezbookkeepingTransactionTimeColumnName     = "Time"
	ezbookkeepingTransactionTimezoneColumnName = "Timezone"
	ezbookkeepingTransactionTimeFormat         = "2006-01-02 15:04:05"
)

var compatibleEzbookkeepingTransactionTimeFormats = []string{
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-1-2 15:04:05",
	"2006-1-2 15:04",
	"2006/01/02 15:04:05",
	"2006/01/02 15:04",
	"2006/1/2 15:04:05",
	"2006/1/2 15:04",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04",
	"2006.01.02 15:04:05",
	"2006.01.02 15:04",
	"2006.1.2 15:04:05",
	"2006.1.2 15:04",
	"2006年1月2日 15:04:05",
	"2006年1月2日 15:04",
	time.RFC3339Nano,
}

// defaultPlainTextDataTable defines the structure of ezbookkeeping default plain text data table
type defaultPlainTextDataTable struct {
	columnSeparator          string
	lineSeparator            string
	allLines                 [][]string
	headerLineColumnNames    []string
	transactionTimeIndex     int
	transactionTimezoneIndex int
}

// defaultPlainTextDataRow defines the structure of ezbookkeeping default plain text data row
type defaultPlainTextDataRow struct {
	allItems []string
}

// defaultPlainTextDataRowIterator defines the structure of ezbookkeeping default plain text data row iterator
type defaultPlainTextDataRowIterator struct {
	dataTable    *defaultPlainTextDataTable
	currentIndex int
}

// defaultTransactionPlainTextDataTableBuilder defines the structure of ezbookkeeping default transaction plain text data table builder
type defaultTransactionPlainTextDataTableBuilder struct {
	columnSeparator       string
	lineSeparator         string
	columns               []datatable.TransactionDataTableColumn
	dataColumnNameMapping map[datatable.TransactionDataTableColumn]string
	dataLineFormat        string
	builder               *strings.Builder
}

// DataRowCount returns the total count of data row
func (t *defaultPlainTextDataTable) DataRowCount() int {
	if len(t.allLines) < 1 {
		return 0
	}

	return len(t.allLines) - 1
}

// HeaderColumnNames returns the header column name list
func (t *defaultPlainTextDataTable) HeaderColumnNames() []string {
	return t.headerLineColumnNames
}

// DataRowIterator returns the iterator of data row
func (t *defaultPlainTextDataTable) DataRowIterator() datatable.BasicDataTableRowIterator {
	return &defaultPlainTextDataRowIterator{
		dataTable:    t,
		currentIndex: 0,
	}
}

// ColumnCount returns the total count of column in this data row
func (r *defaultPlainTextDataRow) ColumnCount() int {
	return len(r.allItems)
}

// GetData returns the data in the specified column index
func (r *defaultPlainTextDataRow) GetData(columnIndex int) string {
	if columnIndex >= len(r.allItems) {
		return ""
	}

	return r.allItems[columnIndex]
}

// HasNext returns whether the iterator does not reach the end
func (t *defaultPlainTextDataRowIterator) HasNext() bool {
	return t.currentIndex+1 < len(t.dataTable.allLines)
}

// CurrentRowId returns current index
func (t *defaultPlainTextDataRowIterator) CurrentRowId() string {
	return fmt.Sprintf("line#%d", t.currentIndex)
}

// Next returns the next basic data row
func (t *defaultPlainTextDataRowIterator) Next() datatable.BasicDataTableRow {
	if t.currentIndex+1 >= len(t.dataTable.allLines) {
		return nil
	}

	t.currentIndex++

	rowItems := append([]string(nil), t.dataTable.allLines[t.currentIndex]...)

	if t.dataTable.transactionTimeIndex >= 0 && t.dataTable.transactionTimeIndex < len(rowItems) {
		timezoneValue := ""

		if t.dataTable.transactionTimezoneIndex >= 0 && t.dataTable.transactionTimezoneIndex < len(rowItems) {
			timezoneValue = rowItems[t.dataTable.transactionTimezoneIndex]
		}

		rowItems[t.dataTable.transactionTimeIndex] = normalizeEzbookkeepingTransactionTime(
			rowItems[t.dataTable.transactionTimeIndex],
			timezoneValue,
		)
	}

	return &defaultPlainTextDataRow{
		allItems: rowItems,
	}
}

// AppendTransaction appends the specified transaction to data builder
func (b *defaultTransactionPlainTextDataTableBuilder) AppendTransaction(data map[datatable.TransactionDataTableColumn]string) {
	dataRowParams := make([]any, len(b.columns))

	for i := 0; i < len(b.columns); i++ {
		dataRowParams[i] = data[b.columns[i]]
	}

	b.builder.WriteString(fmt.Sprintf(b.dataLineFormat, dataRowParams...))
}

// ReplaceDelimiters returns the text after removing the delimiters
func (b *defaultTransactionPlainTextDataTableBuilder) ReplaceDelimiters(text string) string {
	text = strings.Replace(text, "\r\n", " ", -1)
	text = strings.Replace(text, "\r", " ", -1)
	text = strings.Replace(text, "\n", " ", -1)
	text = strings.Replace(text, b.columnSeparator, " ", -1)
	text = strings.Replace(text, b.lineSeparator, " ", -1)

	return text
}

// String returns the textual representation of this data
func (b *defaultTransactionPlainTextDataTableBuilder) String() string {
	return b.builder.String()
}

func (b *defaultTransactionPlainTextDataTableBuilder) generateHeaderLine() string {
	var ret strings.Builder

	for i := 0; i < len(b.columns); i++ {
		if ret.Len() > 0 {
			ret.WriteString(b.columnSeparator)
		}

		dataColumn := b.columns[i]
		columnName := b.dataColumnNameMapping[dataColumn]

		ret.WriteString(columnName)
	}

	ret.WriteString(b.lineSeparator)

	return ret.String()
}

func (b *defaultTransactionPlainTextDataTableBuilder) generateDataLineFormat() string {
	var ret strings.Builder

	for i := 0; i < len(b.columns); i++ {
		if ret.Len() > 0 {
			ret.WriteString(b.columnSeparator)
		}

		ret.WriteString("%s")
	}

	ret.WriteString(b.lineSeparator)

	return ret.String()
}

func createNewDefaultPlainTextDataTable(content string, columnSeparator string, lineSeparator string) (*defaultPlainTextDataTable, error) {
	separatorRunes := []rune(columnSeparator)

	if len(separatorRunes) != 1 {
		return nil, errs.ErrInvalidCSVFile
	}

	reader := csv.NewReader(strings.NewReader(content))
	reader.Comma = separatorRunes[0]
	reader.FieldsPerRecord = -1
	reader.ReuseRecord = false

	allLines := make([][]string, 0)

	for {
		items, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, errs.ErrInvalidCSVFile
		}

		if len(items) == 1 && strings.TrimSpace(items[0]) == "" {
			continue
		}

		allLines = append(allLines, items)
	}

	if len(allLines) < 2 {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	headerLine := allLines[0]
	headerLineItems := append([]string(nil), headerLine...)
	transactionTimeIndex := -1
	transactionTimezoneIndex := -1
	legacyDescriptionIndex := -1
	hasCurrentDescriptionColumn := false

	for i := 0; i < len(headerLineItems); i++ {
		headerLineItems[i] = strings.TrimSpace(strings.TrimPrefix(headerLineItems[i], "\uFEFF"))

		switch headerLineItems[i] {
		case ezbookkeepingTransactionTimeColumnName:
			transactionTimeIndex = i
		case ezbookkeepingTransactionTimezoneColumnName:
			transactionTimezoneIndex = i
		case legacyEzbookkeepingDescriptionColumnName:
			legacyDescriptionIndex = i
		case currentEzbookkeepingDescriptionColumnName:
			hasCurrentDescriptionColumn = true
		}
	}

	// ezBookkeeping v0.1.0-v0.4.1 called the transaction description
	// column "Comment". Normalize it to the current column name so the
	// existing import pipeline preserves legacy descriptions.
	if legacyDescriptionIndex >= 0 && !hasCurrentDescriptionColumn {
		headerLineItems[legacyDescriptionIndex] = currentEzbookkeepingDescriptionColumnName
	}

	return &defaultPlainTextDataTable{
		columnSeparator:          columnSeparator,
		lineSeparator:            lineSeparator,
		allLines:                 allLines,
		headerLineColumnNames:    headerLineItems,
		transactionTimeIndex:     transactionTimeIndex,
		transactionTimezoneIndex: transactionTimezoneIndex,
	}, nil
}

func normalizeEzbookkeepingTransactionTime(value string, timezoneValues ...string) string {
	trimmedValue := normalizeEzbookkeepingTransactionTimeText(value)

	for _, timeFormat := range compatibleEzbookkeepingTransactionTimeFormats {
		parsedTime, err := time.Parse(timeFormat, trimmedValue)

		if err == nil {
			return parsedTime.Format(ezbookkeepingTransactionTimeFormat)
		}
	}

	if parsedTime, ok := parseEzbookkeepingNumericTransactionTime(trimmedValue, timezoneValues...); ok {
		return parsedTime.Format(ezbookkeepingTransactionTimeFormat)
	}

	return trimmedValue
}

func normalizeEzbookkeepingTransactionTimeText(value string) string {
	value = strings.Map(func(char rune) rune {
		switch char {
		case '\u200B', '\u200C', '\u200D', '\u2060', '\uFEFF':
			return -1
		default:
			return char
		}
	}, value)

	value = strings.TrimSpace(value)

	if strings.HasPrefix(value, "=") {
		value = strings.TrimSpace(strings.TrimPrefix(value, "="))
	}

	if len(value) >= 2 {
		firstChar := value[0]
		lastChar := value[len(value)-1]

		if (firstChar == '"' && lastChar == '"') || (firstChar == '\'' && lastChar == '\'') {
			value = strings.TrimSpace(value[1 : len(value)-1])
		}
	}

	return value
}

func parseEzbookkeepingNumericTransactionTime(value string, timezoneValues ...string) (time.Time, bool) {
	if value == "" {
		return time.Time{}, false
	}

	if serial, err := strconv.ParseFloat(value, 64); err == nil && serial > 0 && serial < 2958466 {
		excelEpoch := time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)
		duration := time.Duration(serial * float64(24*time.Hour))
		return excelEpoch.Add(duration), true
	}

	unixValue, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		return time.Time{}, false
	}

	var parsedTime time.Time

	switch {
	case unixValue >= 1_000_000_000_000_000_000:
		parsedTime = time.Unix(0, unixValue)
	case unixValue >= 1_000_000_000_000_000:
		parsedTime = time.Unix(0, unixValue*int64(time.Microsecond))
	case unixValue >= 1_000_000_000_000:
		parsedTime = time.UnixMilli(unixValue)
	case unixValue >= 1_000_000_000:
		parsedTime = time.Unix(unixValue, 0)
	default:
		return time.Time{}, false
	}

	if len(timezoneValues) > 0 {
		if timezone, ok := parseEzbookkeepingTransactionTimezone(timezoneValues[0]); ok {
			parsedTime = parsedTime.In(timezone)
		}
	}

	return parsedTime, true
}

func parseEzbookkeepingTransactionTimezone(value string) (*time.Location, bool) {
	value = strings.TrimSpace(value)

	if value == "" {
		return time.UTC, false
	}

	if minuteOffset, err := strconv.Atoi(value); err == nil {
		return time.FixedZone("Transaction Timezone", minuteOffset*60), true
	}

	normalizedValue := strings.ReplaceAll(value, ":", "")

	if len(normalizedValue) != 5 || (normalizedValue[0] != '+' && normalizedValue[0] != '-') {
		return time.UTC, false
	}

	hours, hourErr := strconv.Atoi(normalizedValue[1:3])
	minutes, minuteErr := strconv.Atoi(normalizedValue[3:5])

	if hourErr != nil || minuteErr != nil || hours > 23 || minutes > 59 {
		return time.UTC, false
	}

	offset := (hours*60 + minutes) * 60

	if normalizedValue[0] == '-' {
		offset = -offset
	}

	return time.FixedZone("Transaction Timezone", offset), true
}

func createNewDefaultTransactionPlainTextDataTableBuilder(transactionCount int, columns []datatable.TransactionDataTableColumn, dataColumnNameMapping map[datatable.TransactionDataTableColumn]string, columnSeparator string, lineSeparator string) *defaultTransactionPlainTextDataTableBuilder {
	var builder strings.Builder
	builder.Grow(transactionCount * 100)

	dataTableBuilder := &defaultTransactionPlainTextDataTableBuilder{
		columnSeparator:       columnSeparator,
		lineSeparator:         lineSeparator,
		columns:               columns,
		dataColumnNameMapping: dataColumnNameMapping,
		builder:               &builder,
	}

	headerLine := dataTableBuilder.generateHeaderLine()
	dataLineFormat := dataTableBuilder.generateDataLineFormat()

	dataTableBuilder.builder.WriteString(headerLine)
	dataTableBuilder.dataLineFormat = dataLineFormat

	return dataTableBuilder
}
