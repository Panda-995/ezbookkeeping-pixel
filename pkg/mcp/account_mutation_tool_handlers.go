package mcp

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

var mcpHexColorPattern = regexp.MustCompile(`^[0-9a-fA-F]{6}$`)

// MCPAccountRecord is the non-recursive representation shared by parent and child accounts.
type MCPAccountRecord struct {
	Id                      string `json:"id"`
	ParentId                string `json:"parent_id,omitempty"`
	Name                    string `json:"name"`
	Category                string `json:"category"`
	Type                    string `json:"type" jsonschema:"enum=single,enum=multi_sub_accounts"`
	Icon                    string `json:"icon"`
	Color                   string `json:"color"`
	Currency                string `json:"currency,omitempty"`
	Balance                 string `json:"balance,omitempty"`
	Comment                 string `json:"comment,omitempty"`
	CreditCardStatementDate *int   `json:"credit_card_statement_date,omitempty"`
	LastReconciledTime      *int64 `json:"last_reconciled_time,omitempty"`
	Hidden                  bool   `json:"hidden"`
}

// MCPAccountInfo is the complete MCP representation of an account hierarchy.
// ezBookkeeping supports one child level, so this finite shape is also safe
// for JSON Schema generation.
type MCPAccountInfo struct {
	MCPAccountRecord
	SubAccounts []*MCPAccountRecord `json:"sub_accounts,omitempty"`
}

// MCPCreateSubAccountRequest represents a sub-account created with a parent account.
type MCPCreateSubAccountRequest struct {
	Name        string `json:"name"`
	Currency    string `json:"currency"`
	Balance     string `json:"balance,omitempty" jsonschema_description:"Displayed opening balance"`
	BalanceTime string `json:"balance_time,omitempty" jsonschema:"format=date-time"`
	Icon        int64  `json:"icon,omitempty"`
	Color       string `json:"color,omitempty"`
	Comment     string `json:"comment,omitempty"`
}

// MCPCreateAccountRequest represents account creation.
type MCPCreateAccountRequest struct {
	Name                    string                       `json:"name"`
	Category                string                       `json:"category" jsonschema:"enum=cash,enum=checking,enum=credit_card,enum=virtual,enum=debt,enum=receivables,enum=investment,enum=savings,enum=certificate_of_deposit"`
	Currency                string                       `json:"currency,omitempty" jsonschema_description:"ISO 4217 currency. Defaults to the user's currency for a single account."`
	Balance                 string                       `json:"balance,omitempty" jsonschema_description:"Displayed opening balance"`
	BalanceTime             string                       `json:"balance_time,omitempty" jsonschema:"format=date-time"`
	Icon                    int64                        `json:"icon,omitempty" jsonschema_description:"Account icon id, defaults to 1"`
	Color                   string                       `json:"color,omitempty" jsonschema_description:"Six-digit RGB hex without #, defaults to 176b5b"`
	Comment                 string                       `json:"comment,omitempty"`
	CreditCardStatementDate int                          `json:"credit_card_statement_date,omitempty" jsonschema_description:"0-28, only for credit_card"`
	SubAccounts             []MCPCreateSubAccountRequest `json:"sub_accounts,omitempty" jsonschema_description:"If present, creates a multi-sub-account parent"`
	DryRun                  bool                         `json:"dry_run,omitempty"`
}

// MCPAccountMutationResponse represents create or update results.
type MCPAccountMutationResponse struct {
	Success bool            `json:"success"`
	DryRun  bool            `json:"dry_run,omitempty"`
	Account *MCPAccountInfo `json:"account"`
}

type mcpCreateAccountToolHandler struct{}

// MCPCreateAccountToolHandler is the MCP handler for account creation.
var MCPCreateAccountToolHandler = &mcpCreateAccountToolHandler{}

func (h *mcpCreateAccountToolHandler) Name() string {
	return "create_account"
}

func (h *mcpCreateAccountToolHandler) Description() string {
	return "Create a single account or a parent account with sub-accounts, including an auditable opening balance."
}

func (h *mcpCreateAccountToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPCreateAccountRequest{})
}

func (h *mcpCreateAccountToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPAccountMutationResponse{})
}

func (h *mcpCreateAccountToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPCreateAccountRequest

	if callToolReq.Arguments == nil {
		return nil, nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	if err := json.Unmarshal(callToolReq.Arguments, &request); err != nil {
		return nil, nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	category, err := textualAccountCategoryToModel(request.Category)

	if err != nil {
		return nil, nil, err
	}

	if err = validateMCPAccountTextFields(request.Name, request.Comment, request.Color); err != nil {
		return nil, nil, err
	}

	if request.CreditCardStatementDate < 0 || request.CreditCardStatementDate > 28 {
		return nil, nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	if category != models.ACCOUNT_CATEGORY_CREDIT_CARD && request.CreditCardStatementDate != 0 {
		return nil, nil, errs.ErrCannotSetStatementDateForNonCreditCard
	}

	icon := request.Icon
	if icon <= 0 {
		icon = 1
	}

	color := strings.ToLower(request.Color)
	if color == "" {
		color = "176b5b"
	}

	accountType := models.ACCOUNT_TYPE_SINGLE_ACCOUNT
	currency := strings.ToUpper(request.Currency)

	if len(request.SubAccounts) > 0 {
		accountType = models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS
		currency = validators.ParentAccountCurrencyPlaceholder
	} else if currency == "" {
		currency = user.DefaultCurrency
	}

	if accountType == models.ACCOUNT_TYPE_SINGLE_ACCOUNT && !validators.AllCurrencyNames[currency] {
		return nil, nil, errs.ErrAccountCurrencyInvalid
	}

	balance, balanceTime, err := parseMCPAccountOpeningBalance(request.Balance, request.BalanceTime, category)

	if err != nil {
		return nil, nil, err
	}

	if accountType == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS {
		balance = 0
		balanceTime = 0
	}

	maxOrder, err := services.GetAccountService().GetMaxDisplayOrder(c, user.Uid, category)

	if err != nil {
		return nil, nil, err
	}

	statementDate := request.CreditCardStatementDate
	mainAccount := &models.Account{
		Uid:          user.Uid,
		Name:         strings.TrimSpace(request.Name),
		DisplayOrder: maxOrder + 1,
		Category:     category,
		Type:         accountType,
		Icon:         icon,
		Color:        color,
		Currency:     currency,
		Balance:      balance,
		Comment:      request.Comment,
		Extend:       &models.AccountExtend{},
	}

	if category == models.ACCOUNT_CATEGORY_CREDIT_CARD {
		mainAccount.Extend.CreditCardStatementDate = &statementDate
	}

	subAccounts := make([]*models.Account, 0, len(request.SubAccounts))
	subAccountBalanceTimes := make([]int64, 0, len(request.SubAccounts))

	for index, subAccountRequest := range request.SubAccounts {
		if err = validateMCPAccountTextFields(subAccountRequest.Name, subAccountRequest.Comment, subAccountRequest.Color); err != nil {
			return nil, nil, err
		}

		subCurrency := strings.ToUpper(subAccountRequest.Currency)
		if !validators.AllCurrencyNames[subCurrency] {
			return nil, nil, errs.ErrAccountCurrencyInvalid
		}

		subBalance, subBalanceTime, parseErr := parseMCPAccountOpeningBalance(subAccountRequest.Balance, subAccountRequest.BalanceTime, category)
		if parseErr != nil {
			return nil, nil, parseErr
		}

		subIcon := subAccountRequest.Icon
		if subIcon <= 0 {
			subIcon = icon
		}

		subColor := strings.ToLower(subAccountRequest.Color)
		if subColor == "" {
			subColor = color
		}

		subAccounts = append(subAccounts, &models.Account{
			Uid:          user.Uid,
			Name:         strings.TrimSpace(subAccountRequest.Name),
			DisplayOrder: int32(index + 1),
			Category:     category,
			Type:         models.ACCOUNT_TYPE_SINGLE_ACCOUNT,
			Icon:         subIcon,
			Color:        subColor,
			Currency:     subCurrency,
			Balance:      subBalance,
			Comment:      subAccountRequest.Comment,
			Extend:       &models.AccountExtend{},
		})
		subAccountBalanceTimes = append(subAccountBalanceTimes, subBalanceTime)
	}

	if !request.DryRun {
		clientTimezone, timezoneErr := c.GetClientTimezone()
		if timezoneErr != nil {
			clientTimezone = time.UTC
		}

		err = services.GetAccountService().CreateAccounts(c, mainAccount, balanceTime, subAccounts, subAccountBalanceTimes, clientTimezone)
		if err != nil {
			log.Errorf(c, "[create_account.Handle] failed to create account for user \"uid:%d\", because %s", user.Uid, err.Error())
			return nil, nil, err
		}

		log.Infof(c, "[create_account.Handle] user \"uid:%d\" created account \"id:%d\"", user.Uid, mainAccount.AccountId)
	}

	accountInfo := createMCPAccountInfo(mainAccount)
	for _, subAccount := range subAccounts {
		accountInfo.SubAccounts = append(accountInfo.SubAccounts, createMCPAccountRecord(subAccount))
	}

	return newMCPAccountMutationResponse(accountInfo, request.DryRun)
}

// MCPUpdateAccountRequest represents a partial update of every mutable account field.
type MCPUpdateAccountRequest struct {
	Id                      string  `json:"id,omitempty" jsonschema_description:"Account id returned by query_all_accounts"`
	AccountName             string  `json:"account_name,omitempty" jsonschema_description:"Current exact account name when id is omitted"`
	Name                    *string `json:"name,omitempty"`
	Category                *string `json:"category,omitempty" jsonschema:"enum=cash,enum=checking,enum=credit_card,enum=virtual,enum=debt,enum=receivables,enum=investment,enum=savings,enum=certificate_of_deposit"`
	Icon                    *int64  `json:"icon,omitempty"`
	Color                   *string `json:"color,omitempty" jsonschema_description:"Six-digit RGB hex without #"`
	Comment                 *string `json:"comment,omitempty"`
	Currency                *string `json:"currency,omitempty" jsonschema_description:"Replacement ISO 4217 currency for a single account"`
	Balance                 *string `json:"balance,omitempty" jsonschema_description:"Replacement displayed balance for a single account"`
	Hidden                  *bool   `json:"hidden,omitempty"`
	CreditCardStatementDate *int    `json:"credit_card_statement_date,omitempty"`
	LastReconciledTime      *int64  `json:"last_reconciled_time,omitempty" jsonschema_description:"Transaction sequence time used by reconciliation"`
	ClearLastReconciledTime bool    `json:"clear_last_reconciled_time,omitempty"`
	DryRun                  bool    `json:"dry_run,omitempty"`
}

type mcpUpdateAccountToolHandler struct{}

// MCPUpdateAccountToolHandler is the MCP handler for complete mutable account updates.
var MCPUpdateAccountToolHandler = &mcpUpdateAccountToolHandler{}

func (h *mcpUpdateAccountToolHandler) Name() string {
	return "update_account"
}

func (h *mcpUpdateAccountToolHandler) Description() string {
	return "Update any mutable account field, including name, category, icon, color, currency, balance, description, visibility, statement date, or reconciliation time. Omitted fields are preserved."
}

func (h *mcpUpdateAccountToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPUpdateAccountRequest{})
}

func (h *mcpUpdateAccountToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPAccountMutationResponse{})
}

func (h *mcpUpdateAccountToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPUpdateAccountRequest

	if callToolReq.Arguments == nil {
		return nil, nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	if err := json.Unmarshal(callToolReq.Arguments, &request); err != nil {
		return nil, nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	account, err := findMCPAccount(c, user.Uid, request.Id, request.AccountName, services)
	if err != nil {
		return nil, nil, err
	}

	updatedAccount := *account
	updatedAccount.Extend = cloneMCPAccountExtend(account.Extend)

	if request.Name != nil {
		updatedAccount.Name = strings.TrimSpace(*request.Name)
	}
	if request.Category != nil {
		updatedAccount.Category, err = textualAccountCategoryToModel(*request.Category)
		if err != nil {
			return nil, nil, err
		}
	}
	if request.Icon != nil {
		updatedAccount.Icon = *request.Icon
	}
	if request.Color != nil {
		updatedAccount.Color = strings.ToLower(*request.Color)
	}
	if request.Comment != nil {
		updatedAccount.Comment = *request.Comment
	}
	if request.Currency != nil {
		if updatedAccount.Type != models.ACCOUNT_TYPE_SINGLE_ACCOUNT {
			return nil, nil, errs.ErrNotSupportedChangeCurrency
		}

		currency := strings.ToUpper(strings.TrimSpace(*request.Currency))
		if !validators.AllCurrencyNames[currency] {
			return nil, nil, errs.ErrAccountCurrencyInvalid
		}
		updatedAccount.Currency = currency
	}
	if request.Balance != nil {
		if updatedAccount.Type != models.ACCOUNT_TYPE_SINGLE_ACCOUNT {
			return nil, nil, errs.ErrNotSupportedChangeBalance
		}

		balance, balanceErr := utils.ParseAmount(*request.Balance)
		if balanceErr != nil {
			return nil, nil, balanceErr
		}
		if updatedAccount.Category.IsLiability() {
			balance = -balance
		}
		updatedAccount.Balance = balance
	} else if updatedAccount.Type == models.ACCOUNT_TYPE_SINGLE_ACCOUNT &&
		account.Category.IsLiability() != updatedAccount.Category.IsLiability() {
		updatedAccount.Balance = -updatedAccount.Balance
	}
	if request.Hidden != nil {
		updatedAccount.Hidden = *request.Hidden
	}

	if err = validateMCPAccountTextFields(updatedAccount.Name, updatedAccount.Comment, updatedAccount.Color); err != nil {
		return nil, nil, err
	}
	if updatedAccount.Icon <= 0 {
		return nil, nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	if updatedAccount.ParentAccountId != models.LevelOneAccountParentId {
		parentAccount, parentErr := services.GetAccountService().GetAccountByAccountId(c, user.Uid, updatedAccount.ParentAccountId)
		if parentErr != nil {
			return nil, nil, parentErr
		}
		if updatedAccount.Category != parentAccount.Category {
			return nil, nil, errs.ErrSubAccountCategoryNotEqualsToParent
		}
		if request.CreditCardStatementDate != nil {
			return nil, nil, errs.ErrCannotSetStatementDateForSubAccount
		}
	}

	if request.CreditCardStatementDate != nil {
		if updatedAccount.Category != models.ACCOUNT_CATEGORY_CREDIT_CARD {
			return nil, nil, errs.ErrCannotSetStatementDateForNonCreditCard
		}
		if *request.CreditCardStatementDate < 0 || *request.CreditCardStatementDate > 28 {
			return nil, nil, errs.ErrIncompleteOrIncorrectSubmission
		}
		updatedAccount.Extend.CreditCardStatementDate = request.CreditCardStatementDate
	} else if updatedAccount.Category != models.ACCOUNT_CATEGORY_CREDIT_CARD {
		updatedAccount.Extend.CreditCardStatementDate = nil
	}

	if request.ClearLastReconciledTime {
		updatedAccount.Extend.LastReconciledTime = nil
	} else if request.LastReconciledTime != nil {
		updatedAccount.Extend.LastReconciledTime = request.LastReconciledTime
	}

	updateAccounts := []*models.Account{&updatedAccount}
	mainAccount := account
	var accountAndSubAccounts []*models.Account

	if account.ParentAccountId != models.LevelOneAccountParentId {
		mainAccount, err = services.GetAccountService().GetAccountByAccountId(c, user.Uid, account.ParentAccountId)
		if err != nil {
			return nil, nil, err
		}
	} else if account.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS && updatedAccount.Category != account.Category {
		accountAndSubAccounts, err = services.GetAccountService().GetAccountAndSubAccountsByAccountId(c, user.Uid, account.AccountId)
		if err != nil {
			return nil, nil, err
		}
		for _, subAccount := range accountAndSubAccounts {
			if subAccount.AccountId == account.AccountId {
				continue
			}
			updatedSubAccount := *subAccount
			updatedSubAccount.Category = updatedAccount.Category
			if subAccount.Category.IsLiability() != updatedSubAccount.Category.IsLiability() {
				updatedSubAccount.Balance = -updatedSubAccount.Balance
			}
			updateAccounts = append(updateAccounts, &updatedSubAccount)
		}
	}

	if accountModelsEqualForMCP(account, &updatedAccount) && len(updateAccounts) == 1 {
		return nil, nil, errs.ErrNothingWillBeUpdated
	}

	if !request.DryRun {
		clientTimezone, timezoneErr := c.GetClientTimezone()
		if timezoneErr != nil {
			clientTimezone = time.UTC
		}
		err = services.GetAccountService().ModifyAccounts(c, mainAccount, updateAccounts, nil, nil, nil, clientTimezone)
		if err != nil {
			log.Errorf(c, "[update_account.Handle] failed to update account \"id:%d\" for user \"uid:%d\", because %s", account.AccountId, user.Uid, err.Error())
			return nil, nil, err
		}
		log.Infof(c, "[update_account.Handle] user \"uid:%d\" updated account \"id:%d\"", user.Uid, account.AccountId)
	}

	info := createMCPAccountInfo(&updatedAccount)
	if updatedAccount.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS {
		children := accountAndSubAccounts
		if len(children) == 0 {
			children, err = services.GetAccountService().GetSubAccountsByAccountId(c, user.Uid, updatedAccount.AccountId)
		}

		updatedAccountMap := make(map[int64]*models.Account, len(updateAccounts))
		for _, changedAccount := range updateAccounts {
			updatedAccountMap[changedAccount.AccountId] = changedAccount
		}

		if err == nil {
			for _, child := range children {
				if child.AccountId == updatedAccount.AccountId {
					continue
				}
				if changedChild := updatedAccountMap[child.AccountId]; changedChild != nil {
					child = changedChild
				}
				info.SubAccounts = append(info.SubAccounts, createMCPAccountRecord(child))
			}
		}
	}

	return newMCPAccountMutationResponse(info, request.DryRun)
}

// MCPDeleteAccountRequest represents an account deletion.
type MCPDeleteAccountRequest struct {
	Id          string `json:"id,omitempty"`
	AccountName string `json:"account_name,omitempty"`
	DryRun      bool   `json:"dry_run,omitempty"`
}

// MCPDeleteAccountResponse represents the result of account deletion.
type MCPDeleteAccountResponse struct {
	Success   bool   `json:"success"`
	DryRun    bool   `json:"dry_run,omitempty"`
	AccountId string `json:"account_id"`
}

type mcpDeleteAccountToolHandler struct{}

// MCPDeleteAccountToolHandler is the MCP handler for account deletion.
var MCPDeleteAccountToolHandler = &mcpDeleteAccountToolHandler{}

func (h *mcpDeleteAccountToolHandler) Name() string {
	return "delete_account"
}

func (h *mcpDeleteAccountToolHandler) Description() string {
	return "Delete an unused account or sub-account. Accounts referenced by ordinary transactions remain protected by service-layer rules."
}

func (h *mcpDeleteAccountToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPDeleteAccountRequest{})
}

func (h *mcpDeleteAccountToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPDeleteAccountResponse{})
}

func (h *mcpDeleteAccountToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPDeleteAccountRequest
	if callToolReq.Arguments == nil {
		return nil, nil, errs.ErrIncompleteOrIncorrectSubmission
	}
	if err := json.Unmarshal(callToolReq.Arguments, &request); err != nil {
		return nil, nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	account, err := findMCPAccount(c, user.Uid, request.Id, request.AccountName, services)
	if err != nil {
		return nil, nil, err
	}

	if !request.DryRun {
		if account.ParentAccountId == models.LevelOneAccountParentId {
			err = services.GetAccountService().DeleteAccount(c, user.Uid, account.AccountId)
		} else {
			err = services.GetAccountService().DeleteSubAccount(c, user.Uid, account.AccountId)
		}
		if err != nil {
			return nil, nil, err
		}
		log.Infof(c, "[delete_account.Handle] user \"uid:%d\" deleted account \"id:%d\"", user.Uid, account.AccountId)
	}

	response := MCPDeleteAccountResponse{
		Success:   true,
		DryRun:    request.DryRun,
		AccountId: utils.Int64ToString(account.AccountId),
	}
	content, err := json.Marshal(response)
	if err != nil {
		return nil, nil, err
	}
	return response, []*MCPTextContent{NewMCPTextContent(string(content))}, nil
}

func textualAccountCategoryToModel(category string) (models.AccountCategory, error) {
	switch category {
	case "cash":
		return models.ACCOUNT_CATEGORY_CASH, nil
	case "checking":
		return models.ACCOUNT_CATEGORY_CHECKING_ACCOUNT, nil
	case "credit_card":
		return models.ACCOUNT_CATEGORY_CREDIT_CARD, nil
	case "virtual":
		return models.ACCOUNT_CATEGORY_VIRTUAL, nil
	case "debt":
		return models.ACCOUNT_CATEGORY_DEBT, nil
	case "receivables":
		return models.ACCOUNT_CATEGORY_RECEIVABLES, nil
	case "investment":
		return models.ACCOUNT_CATEGORY_INVESTMENT, nil
	case "savings":
		return models.ACCOUNT_CATEGORY_SAVINGS_ACCOUNT, nil
	case "certificate_of_deposit":
		return models.ACCOUNT_CATEGORY_CERTIFICATE_OF_DEPOSIT, nil
	default:
		return 0, errs.ErrAccountCategoryInvalid
	}
}

func accountCategoryToText(category models.AccountCategory) string {
	switch category {
	case models.ACCOUNT_CATEGORY_CASH:
		return "cash"
	case models.ACCOUNT_CATEGORY_CHECKING_ACCOUNT:
		return "checking"
	case models.ACCOUNT_CATEGORY_CREDIT_CARD:
		return "credit_card"
	case models.ACCOUNT_CATEGORY_VIRTUAL:
		return "virtual"
	case models.ACCOUNT_CATEGORY_DEBT:
		return "debt"
	case models.ACCOUNT_CATEGORY_RECEIVABLES:
		return "receivables"
	case models.ACCOUNT_CATEGORY_INVESTMENT:
		return "investment"
	case models.ACCOUNT_CATEGORY_SAVINGS_ACCOUNT:
		return "savings"
	case models.ACCOUNT_CATEGORY_CERTIFICATE_OF_DEPOSIT:
		return "certificate_of_deposit"
	default:
		return ""
	}
}

func validateMCPAccountTextFields(name string, comment string, color string) error {
	if strings.TrimSpace(name) == "" || utf8.RuneCountInString(name) > 64 || utf8.RuneCountInString(comment) > 255 {
		return errs.ErrIncompleteOrIncorrectSubmission
	}
	if color != "" && !mcpHexColorPattern.MatchString(color) {
		return errs.ErrIncompleteOrIncorrectSubmission
	}
	return nil
}

func parseMCPAccountOpeningBalance(balanceText string, balanceTimeText string, category models.AccountCategory) (int64, int64, error) {
	balance := int64(0)
	var err error

	if balanceText != "" {
		balance, err = utils.ParseAmount(balanceText)
		if err != nil {
			return 0, 0, err
		}
	}
	if category.IsLiability() {
		balance = -balance
	}

	balanceTime := int64(0)
	if balanceTimeText != "" {
		parsedTime, parseErr := utils.ParseFromLongDateTimeWithTimezoneRFC3339Format(balanceTimeText)
		if parseErr != nil {
			return 0, 0, parseErr
		}
		balanceTime = parsedTime.Unix()
	}
	if balance != 0 && balanceTime <= 0 {
		return 0, 0, errs.ErrAccountBalanceTimeNotSet
	}
	return balance, balanceTime, nil
}

func findMCPAccount(c *core.WebContext, uid int64, id string, name string, services MCPAvailableServices) (*models.Account, error) {
	if id != "" {
		accountId, err := utils.StringToInt64(id)
		if err != nil || accountId <= 0 {
			return nil, errs.ErrAccountIdInvalid
		}
		return services.GetAccountService().GetAccountByAccountId(c, uid, accountId)
	}
	if name == "" {
		return nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	accounts, err := services.GetAccountService().GetAllAccountsByUid(c, uid)
	if err != nil {
		return nil, err
	}
	for _, account := range accounts {
		if account.Name == name {
			return account, nil
		}
	}
	return nil, errs.ErrAccountNotFound
}

func cloneMCPAccountExtend(extend *models.AccountExtend) *models.AccountExtend {
	if extend == nil {
		return &models.AccountExtend{}
	}
	cloned := *extend
	return &cloned
}

func accountModelsEqualForMCP(left *models.Account, right *models.Account) bool {
	leftStatementDate := 0
	rightStatementDate := 0
	var leftReconciledTime *int64
	var rightReconciledTime *int64
	if left.Extend != nil {
		if left.Extend.CreditCardStatementDate != nil {
			leftStatementDate = *left.Extend.CreditCardStatementDate
		}
		leftReconciledTime = left.Extend.LastReconciledTime
	}
	if right.Extend != nil {
		if right.Extend.CreditCardStatementDate != nil {
			rightStatementDate = *right.Extend.CreditCardStatementDate
		}
		rightReconciledTime = right.Extend.LastReconciledTime
	}

	reconciledEqual := (leftReconciledTime == nil && rightReconciledTime == nil) ||
		(leftReconciledTime != nil && rightReconciledTime != nil && *leftReconciledTime == *rightReconciledTime)

	return left.Name == right.Name &&
		left.Category == right.Category &&
		left.Icon == right.Icon &&
		left.Color == right.Color &&
		left.Comment == right.Comment &&
		left.Currency == right.Currency &&
		left.Balance == right.Balance &&
		left.Hidden == right.Hidden &&
		leftStatementDate == rightStatementDate &&
		reconciledEqual
}

func createMCPAccountInfo(account *models.Account) *MCPAccountInfo {
	return &MCPAccountInfo{
		MCPAccountRecord: *createMCPAccountRecord(account),
	}
}

func createMCPAccountRecord(account *models.Account) *MCPAccountRecord {
	balance := account.Balance
	if account.Category.IsLiability() {
		balance = -balance
	}

	accountType := "single"
	if account.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS {
		accountType = "multi_sub_accounts"
	}

	info := &MCPAccountRecord{
		Id:       utils.Int64ToString(account.AccountId),
		Name:     account.Name,
		Category: accountCategoryToText(account.Category),
		Type:     accountType,
		Icon:     utils.Int64ToString(account.Icon),
		Color:    account.Color,
		Comment:  account.Comment,
		Hidden:   account.Hidden,
	}
	if account.ParentAccountId > 0 {
		info.ParentId = utils.Int64ToString(account.ParentAccountId)
	}
	if account.Type == models.ACCOUNT_TYPE_SINGLE_ACCOUNT {
		info.Currency = account.Currency
		info.Balance = utils.FormatAmount(balance)
	}
	if account.Extend != nil {
		info.CreditCardStatementDate = account.Extend.CreditCardStatementDate
		info.LastReconciledTime = account.Extend.LastReconciledTime
	}
	return info
}

func newMCPAccountMutationResponse(accountInfo *MCPAccountInfo, dryRun bool) (any, []*MCPTextContent, error) {
	response := MCPAccountMutationResponse{
		Success: true,
		DryRun:  dryRun,
		Account: accountInfo,
	}
	content, err := json.Marshal(response)
	if err != nil {
		return nil, nil, err
	}
	return response, []*MCPTextContent{NewMCPTextContent(string(content))}, nil
}
