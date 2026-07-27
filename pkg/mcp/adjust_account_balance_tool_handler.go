package mcp

import (
	"encoding/json"
	"reflect"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// MCPAdjustAccountBalanceRequest represents all parameters for an auditable account balance adjustment.
type MCPAdjustAccountBalanceRequest struct {
	AccountName string   `json:"account_name" jsonschema_description:"Account name whose balance should be adjusted"`
	Balance     string   `json:"balance" jsonschema_description:"Absolute account balance after the adjustment"`
	Time        string   `json:"time" jsonschema:"format=date-time" jsonschema_description:"Adjustment time in RFC 3339 format"`
	Tags        []string `json:"tags,omitempty" jsonschema_description:"Optional transaction tags, maximum 10"`
	Comment     string   `json:"comment,omitempty" jsonschema_description:"Reason for the balance adjustment"`
	DryRun      bool     `json:"dry_run,omitempty" jsonschema_description:"Validate and preview without saving"`
}

type mcpAdjustAccountBalanceToolHandler struct{}

// MCPAdjustAccountBalanceToolHandler is the MCP handler for auditable balance adjustments.
var MCPAdjustAccountBalanceToolHandler = &mcpAdjustAccountBalanceToolHandler{}

func (h *mcpAdjustAccountBalanceToolHandler) Name() string {
	return "adjust_account_balance"
}

func (h *mcpAdjustAccountBalanceToolHandler) Description() string {
	return "Set an account to an absolute balance by creating an auditable balance-modification transaction."
}

func (h *mcpAdjustAccountBalanceToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPAdjustAccountBalanceRequest{})
}

func (h *mcpAdjustAccountBalanceToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPAddTransactionResponse{})
}

func (h *mcpAdjustAccountBalanceToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPAdjustAccountBalanceRequest

	if callToolReq.Arguments == nil {
		return nil, nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	if err := json.Unmarshal(callToolReq.Arguments, &request); err != nil {
		return nil, nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	delegatedArguments, err := json.Marshal(&MCPAddTransactionRequest{
		Type:        transactionTypeModifyBalance,
		Time:        request.Time,
		AccountName: request.AccountName,
		Amount:      request.Balance,
		Tags:        request.Tags,
		Comment:     request.Comment,
		DryRun:      request.DryRun,
	})

	if err != nil {
		return nil, nil, err
	}

	delegatedRequest := &MCPCallToolRequest{
		Name:      MCPAddTransactionToolHandler.Name(),
		Arguments: delegatedArguments,
	}

	return MCPAddTransactionToolHandler.Handle(c, delegatedRequest, user, currentConfig, services)
}
