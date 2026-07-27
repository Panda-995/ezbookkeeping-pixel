package mcp

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// MCPDeleteTransactionRequest represents a transaction deletion.
type MCPDeleteTransactionRequest struct {
	Id     string `json:"id" jsonschema_description:"Transaction id returned by query_transactions"`
	DryRun bool   `json:"dry_run,omitempty" jsonschema_description:"Validate without deleting"`
}

// MCPDeleteTransactionResponse represents the result of a transaction deletion.
type MCPDeleteTransactionResponse struct {
	Success       bool   `json:"success"`
	DryRun        bool   `json:"dry_run,omitempty"`
	TransactionId string `json:"transaction_id"`
}

type mcpDeleteTransactionToolHandler struct{}

// MCPDeleteTransactionToolHandler is the MCP handler for transaction deletion.
var MCPDeleteTransactionToolHandler = &mcpDeleteTransactionToolHandler{}

func (h *mcpDeleteTransactionToolHandler) Name() string {
	return "delete_transaction"
}

func (h *mcpDeleteTransactionToolHandler) Description() string {
	return "Delete a transaction and its related transfer record after applying the user's edit-scope rules."
}

func (h *mcpDeleteTransactionToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPDeleteTransactionRequest{})
}

func (h *mcpDeleteTransactionToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPDeleteTransactionResponse{})
}

func (h *mcpDeleteTransactionToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPDeleteTransactionRequest

	if callToolReq.Arguments == nil {
		return nil, nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	if err := json.Unmarshal(callToolReq.Arguments, &request); err != nil {
		return nil, nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	transactionId, err := utils.StringToInt64(request.Id)

	if err != nil || transactionId <= 0 {
		return nil, nil, errs.ErrTransactionIdInvalid
	}

	transaction, err := services.GetTransactionService().GetTransactionByTransactionId(c, user.Uid, transactionId)

	if err != nil {
		return nil, nil, err
	}

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		return nil, nil, errs.ErrTransactionTypeInvalid
	}

	allAccounts, err := services.GetAccountService().GetAllAccountsByUid(c, user.Uid)

	if err != nil {
		return nil, nil, err
	}

	accountMap := services.GetAccountService().GetAccountMapByList(allAccounts)
	clientTimezone, err := c.GetClientTimezone()

	if err != nil {
		clientTimezone = time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)
	}

	if !user.CanEditTransactionByTransactionTime(transaction.TransactionTime, clientTimezone, accountMap[transaction.AccountId], accountMap[transaction.RelatedAccountId]) {
		return nil, nil, errs.ErrCannotModifyTransactionWithThisTransactionTime
	}

	if !request.DryRun {
		err = services.GetTransactionService().DeleteTransaction(c, user.Uid, transaction.TransactionId)

		if err != nil {
			log.Errorf(c, "[delete_transaction.Handle] failed to delete transaction \"id:%d\" for user \"uid:%d\", because %s", transactionId, user.Uid, err.Error())
			return nil, nil, err
		}

		log.Infof(c, "[delete_transaction.Handle] user \"uid:%d\" deleted transaction \"id:%d\"", user.Uid, transactionId)
	}

	response := MCPDeleteTransactionResponse{
		Success:       true,
		DryRun:        request.DryRun,
		TransactionId: request.Id,
	}
	content, err := json.Marshal(response)

	if err != nil {
		return nil, nil, err
	}

	return response, []*MCPTextContent{NewMCPTextContent(string(content))}, nil
}
