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

// MCPUpdateTransactionRequest represents a partial transaction update.
// Omitted fields preserve their existing value. An explicitly empty tags array clears all tags.
type MCPUpdateTransactionRequest struct {
	Id                     string              `json:"id" jsonschema_description:"Transaction id returned by query_transactions"`
	Type                   *string             `json:"type,omitempty" jsonschema:"enum=balance_modification,enum=income,enum=expense,enum=transfer" jsonschema_description:"New transaction type. Balance modifications cannot be converted to another type or vice versa."`
	Time                   *string             `json:"time,omitempty" jsonschema:"format=date-time" jsonschema_description:"New transaction time in RFC 3339 format"`
	SecondaryCategoryName  *string             `json:"category_name,omitempty" jsonschema_description:"New secondary category name"`
	AccountName            *string             `json:"account_name,omitempty" jsonschema_description:"New source account name"`
	Amount                 *string             `json:"amount,omitempty" jsonschema_description:"New source amount or absolute balance for balance_modification"`
	DestinationAccountName *string             `json:"destination_account_name,omitempty" jsonschema_description:"New destination account for a transfer"`
	DestinationAmount      *string             `json:"destination_amount,omitempty" jsonschema_description:"New destination amount for a transfer"`
	Tags                   *[]string           `json:"tags,omitempty" jsonschema_description:"Replacement tag list. Pass an empty array to clear all tags."`
	PictureIds             *[]string           `json:"picture_ids,omitempty" jsonschema_description:"Replacement list of existing or previously uploaded picture ids. Pass an empty array to clear pictures."`
	Comment                *string             `json:"comment,omitempty" jsonschema_description:"Replacement transaction description"`
	HideAmount             *bool               `json:"hide_amount,omitempty" jsonschema_description:"Whether to hide this amount from statistics"`
	GeoLocation            *MCPGeoLocationInfo `json:"geo_location,omitempty" jsonschema_description:"Replacement geographic location"`
	ClearGeoLocation       bool                `json:"clear_geo_location,omitempty" jsonschema_description:"Clear the existing geographic location"`
	DryRun                 bool                `json:"dry_run,omitempty" jsonschema_description:"Validate and preview without saving"`
}

// MCPUpdateTransactionResponse represents the result of a transaction update.
type MCPUpdateTransactionResponse struct {
	Success     bool                `json:"success"`
	DryRun      bool                `json:"dry_run,omitempty"`
	Transaction *MCPTransactionInfo `json:"transaction"`
}

type mcpUpdateTransactionToolHandler struct{}

// MCPUpdateTransactionToolHandler is the MCP handler for complete transaction field updates.
var MCPUpdateTransactionToolHandler = &mcpUpdateTransactionToolHandler{}

func (h *mcpUpdateTransactionToolHandler) Name() string {
	return "update_transaction"
}

func (h *mcpUpdateTransactionToolHandler) Description() string {
	return "Update transaction type, time, category, accounts, amounts, tags, description, visibility, or geographic location. Omitted fields are preserved."
}

func (h *mcpUpdateTransactionToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPUpdateTransactionRequest{})
}

func (h *mcpUpdateTransactionToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPUpdateTransactionResponse{})
}

func (h *mcpUpdateTransactionToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPUpdateTransactionRequest

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

	uid := user.Uid
	oldTransaction, err := services.GetTransactionService().GetTransactionByTransactionId(c, uid, transactionId)

	if err != nil {
		return nil, nil, err
	}

	if oldTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		return nil, nil, errs.ErrTransactionTypeInvalid
	}

	allAccounts, err := services.GetAccountService().GetAllAccountsByUid(c, uid)

	if err != nil {
		return nil, nil, err
	}

	accountMap := services.GetAccountService().GetAccountMapByList(allAccounts)
	visibleAccountNameMap := services.GetAccountService().GetVisibleAccountNameMapByList(allAccounts)
	newTransaction := *oldTransaction

	if request.Type != nil {
		newType, typeErr := textualTransactionTypeToDbType(*request.Type)

		if typeErr != nil {
			return nil, nil, typeErr
		}

		if (oldTransaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE) != (newType == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE) {
			return nil, nil, errs.ErrTransactionTypeInvalid
		}

		newTransaction.Type = newType
	}

	if request.Time != nil {
		transactionTime, timeErr := utils.ParseFromLongDateTimeWithTimezoneRFC3339Format(*request.Time)

		if timeErr != nil {
			return nil, nil, timeErr
		}

		newTransaction.TransactionTime = utils.GetMinTransactionTimeFromUnixTime(transactionTime.Unix())
		newTransaction.TimezoneUtcOffset = utils.GetTimezoneOffsetMinutes(transactionTime.Unix(), transactionTime.Location())
	}

	if request.AccountName != nil {
		account := visibleAccountNameMap[*request.AccountName]

		if account == nil {
			return nil, nil, errs.ErrSourceAccountNotFound
		}

		newTransaction.AccountId = account.AccountId
	}

	if request.Amount != nil {
		amount, amountErr := utils.ParseAmount(*request.Amount)

		if amountErr != nil {
			return nil, nil, amountErr
		}

		newTransaction.Amount = amount
	}

	if request.Comment != nil {
		newTransaction.Comment = *request.Comment
	}

	if request.HideAmount != nil {
		newTransaction.HideAmount = *request.HideAmount
	}

	if request.ClearGeoLocation {
		newTransaction.GeoLongitude = 0
		newTransaction.GeoLatitude = 0
	} else if request.GeoLocation != nil {
		newTransaction.GeoLongitude = request.GeoLocation.Longitude
		newTransaction.GeoLatitude = request.GeoLocation.Latitude
	}

	if newTransaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		if request.SecondaryCategoryName != nil || request.DestinationAccountName != nil || request.DestinationAmount != nil {
			return nil, nil, errs.ErrTransactionTypeInvalid
		}

		newTransaction.CategoryId = 0
		newTransaction.RelatedAccountId = 0
		newTransaction.RelatedAccountAmount = 0

		displayedAmount := newTransaction.Amount
		if request.Amount == nil {
			if oldAccount := accountMap[oldTransaction.AccountId]; oldAccount != nil && oldAccount.Category.IsLiability() {
				displayedAmount = -displayedAmount
			}
		}
		newTransaction.Amount = displayedAmount
		if account := accountMap[newTransaction.AccountId]; account != nil && account.Category.IsLiability() {
			newTransaction.Amount = -displayedAmount
		}
	} else {
		categoryType := transactionDbTypeToCategoryType(newTransaction.Type)

		if request.SecondaryCategoryName != nil {
			category, categoryErr := findVisibleSecondaryCategoryByName(c, uid, *request.SecondaryCategoryName, categoryType, services)

			if categoryErr != nil {
				return nil, nil, categoryErr
			}

			newTransaction.CategoryId = category.CategoryId
		} else {
			category, categoryErr := services.GetTransactionCategoryService().GetCategoryByCategoryId(c, uid, newTransaction.CategoryId)

			if categoryErr != nil || category.Type != categoryType || category.ParentCategoryId == models.LevelOneTransactionCategoryParentId {
				return nil, nil, errs.ErrTransactionCategoryNotFound
			}
		}

		if newTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			if request.DestinationAccountName != nil {
				destinationAccount := visibleAccountNameMap[*request.DestinationAccountName]

				if destinationAccount == nil {
					return nil, nil, errs.ErrDestinationAccountNotFound
				}

				newTransaction.RelatedAccountId = destinationAccount.AccountId
			}

			if newTransaction.RelatedAccountId <= 0 {
				return nil, nil, errs.ErrDestinationAccountNotFound
			}

			if request.DestinationAmount != nil {
				destinationAmount, amountErr := utils.ParseAmount(*request.DestinationAmount)

				if amountErr != nil {
					return nil, nil, amountErr
				}

				newTransaction.RelatedAccountAmount = destinationAmount
			}
		} else {
			if request.DestinationAccountName != nil || request.DestinationAmount != nil {
				return nil, nil, errs.ErrTransactionTypeInvalid
			}

			newTransaction.RelatedAccountId = 0
			newTransaction.RelatedAccountAmount = 0
		}
	}

	allTransactionTagIds, err := services.GetTransactionTagService().GetAllTagIdsOfTransactions(c, uid, []int64{oldTransaction.TransactionId})

	if err != nil {
		return nil, nil, err
	}

	oldTagIds := allTransactionTagIds[oldTransaction.TransactionId]

	if oldTagIds == nil {
		oldTagIds = make([]int64, 0)
	}

	newTagIds := oldTagIds
	newTagNames := make([]string, 0)

	if request.Tags != nil {
		if len(*request.Tags) > models.MaximumTagsCountOfTransaction {
			return nil, nil, errs.ErrTransactionHasTooManyTags
		}

		allTags, tagErr := services.GetTransactionTagService().GetAllTagsByUid(c, uid)

		if tagErr != nil {
			return nil, nil, tagErr
		}

		tagNameMap := services.GetTransactionTagService().GetVisibleTagNameMapByList(allTags)
		newTagIds = make([]int64, 0, len(*request.Tags))

		for _, tagName := range *request.Tags {
			tag := tagNameMap[tagName]

			if tag == nil {
				return nil, nil, errs.ErrTransactionTagNotFound
			}

			newTagIds = append(newTagIds, tag.TagId)
			newTagNames = append(newTagNames, tag.Name)
		}
	} else {
		if len(oldTagIds) > 0 {
			allTags, tagErr := services.GetTransactionTagService().GetTagsByTagIds(c, uid, oldTagIds)

			if tagErr != nil {
				return nil, nil, tagErr
			}

			for _, tagId := range oldTagIds {
				if tag := allTags[tagId]; tag != nil {
					newTagNames = append(newTagNames, tag.Name)
				}
			}
		}
	}

	oldPictureInfos, err := services.GetTransactionPictureService().GetPictureInfosByTransactionId(c, uid, oldTransaction.TransactionId)
	if err != nil {
		return nil, nil, err
	}
	oldPictureIds := services.GetTransactionPictureService().GetTransactionPictureIds(oldPictureInfos)
	newPictureIds := oldPictureIds

	if request.PictureIds != nil {
		if len(*request.PictureIds) > models.MaximumPicturesCountOfTransaction {
			return nil, nil, errs.ErrTransactionHasTooManyPictures
		}
		newPictureIds, err = utils.StringArrayToInt64Array(*request.PictureIds)
		if err != nil {
			return nil, nil, errs.ErrTransactionPictureIdInvalid
		}
		addPictureIds := utils.Int64SliceMinus(newPictureIds, oldPictureIds)
		if len(addPictureIds) > 0 {
			newPictureInfos, pictureErr := services.GetTransactionPictureService().GetNewPictureInfosByPictureIds(c, uid, addPictureIds)
			if pictureErr != nil {
				return nil, nil, pictureErr
			}
			if len(utils.Int64SliceMinus(addPictureIds, services.GetTransactionPictureService().GetTransactionPictureIds(newPictureInfos))) > 0 {
				return nil, nil, errs.ErrTransactionPictureNotFound
			}
		}
	}

	if transactionModelsEqualForMCP(oldTransaction, &newTransaction) &&
		utils.Int64SliceEquals(oldTagIds, newTagIds) &&
		utils.Int64SliceEquals(oldPictureIds, newPictureIds) {
		return nil, nil, errs.ErrNothingWillBeUpdated
	}

	clientTimezone, err := c.GetClientTimezone()

	if err != nil {
		clientTimezone = time.FixedZone("Transaction Timezone", int(newTransaction.TimezoneUtcOffset)*60)
	}

	oldEditable := user.CanEditTransactionByTransactionTime(oldTransaction.TransactionTime, clientTimezone, accountMap[oldTransaction.AccountId], accountMap[oldTransaction.RelatedAccountId])
	newEditable := user.CanEditTransactionByTransactionTime(newTransaction.TransactionTime, clientTimezone, accountMap[newTransaction.AccountId], accountMap[newTransaction.RelatedAccountId])

	if !oldEditable || !newEditable {
		return nil, nil, errs.ErrCannotModifyTransactionWithThisTransactionTime
	}

	if !request.DryRun {
		changeToTransfer := newTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT && oldTransaction.Type != models.TRANSACTION_DB_TYPE_TRANSFER_OUT
		var addTagIds []int64
		var removeTagIds []int64

		if !utils.Int64SliceEquals(oldTagIds, newTagIds) {
			addTagIds = utils.Int64SliceMinus(newTagIds, oldTagIds)
			removeTagIds = utils.Int64SliceMinus(oldTagIds, newTagIds)
		}

		addPictureIds := utils.Int64SliceMinus(newPictureIds, oldPictureIds)
		removePictureIds := utils.Int64SliceMinus(oldPictureIds, newPictureIds)
		err = services.GetTransactionService().ModifyTransaction(c, &newTransaction, changeToTransfer, len(oldTagIds), addTagIds, removeTagIds, addPictureIds, removePictureIds)

		if err != nil {
			log.Errorf(c, "[update_transaction.Handle] failed to update transaction \"id:%d\" for user \"uid:%d\", because %s", transactionId, uid, err.Error())
			return nil, nil, err
		}

		log.Infof(c, "[update_transaction.Handle] user \"uid:%d\" updated transaction \"id:%d\"", uid, transactionId)
	}

	transactionInfo := createMCPTransactionInfo(&newTransaction, accountMap, services, c, uid, newTagNames)
	transactionInfo.PictureIds = utils.Int64ArrayToStringArray(newPictureIds)
	response := MCPUpdateTransactionResponse{
		Success:     true,
		DryRun:      request.DryRun,
		Transaction: transactionInfo,
	}
	content, err := json.Marshal(response)

	if err != nil {
		return nil, nil, err
	}

	return response, []*MCPTextContent{NewMCPTextContent(string(content))}, nil
}

func textualTransactionTypeToDbType(transactionType string) (models.TransactionDbType, error) {
	switch transactionType {
	case transactionTypeModifyBalance:
		return models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, nil
	case transactionTypeIncome:
		return models.TRANSACTION_DB_TYPE_INCOME, nil
	case transactionTypeExpense:
		return models.TRANSACTION_DB_TYPE_EXPENSE, nil
	case transactionTypeTransfer:
		return models.TRANSACTION_DB_TYPE_TRANSFER_OUT, nil
	default:
		return 0, errs.ErrTransactionTypeInvalid
	}
}

func transactionDbTypeToText(transactionType models.TransactionDbType) string {
	switch transactionType {
	case models.TRANSACTION_DB_TYPE_MODIFY_BALANCE:
		return transactionTypeModifyBalance
	case models.TRANSACTION_DB_TYPE_INCOME:
		return transactionTypeIncome
	case models.TRANSACTION_DB_TYPE_EXPENSE:
		return transactionTypeExpense
	case models.TRANSACTION_DB_TYPE_TRANSFER_OUT:
		return transactionTypeTransfer
	default:
		return ""
	}
}

func transactionDbTypeToCategoryType(transactionType models.TransactionDbType) models.TransactionCategoryType {
	switch transactionType {
	case models.TRANSACTION_DB_TYPE_INCOME:
		return models.CATEGORY_TYPE_INCOME
	case models.TRANSACTION_DB_TYPE_EXPENSE:
		return models.CATEGORY_TYPE_EXPENSE
	case models.TRANSACTION_DB_TYPE_TRANSFER_OUT:
		return models.CATEGORY_TYPE_TRANSFER
	default:
		return 0
	}
}

func findVisibleSecondaryCategoryByName(c *core.WebContext, uid int64, name string, categoryType models.TransactionCategoryType, services MCPAvailableServices) (*models.TransactionCategory, error) {
	categories, err := services.GetTransactionCategoryService().GetAllCategoriesByUid(c, uid, categoryType, -1)

	if err != nil {
		return nil, err
	}

	for _, category := range categories {
		if !category.Hidden && category.ParentCategoryId != models.LevelOneTransactionCategoryParentId && category.Name == name {
			return category, nil
		}
	}

	return nil, errs.ErrTransactionCategoryNotFound
}

func transactionModelsEqualForMCP(left *models.Transaction, right *models.Transaction) bool {
	return left.Type == right.Type &&
		left.CategoryId == right.CategoryId &&
		utils.GetUnixTimeFromTransactionTime(left.TransactionTime) == utils.GetUnixTimeFromTransactionTime(right.TransactionTime) &&
		left.TimezoneUtcOffset == right.TimezoneUtcOffset &&
		left.AccountId == right.AccountId &&
		left.Amount == right.Amount &&
		left.RelatedAccountId == right.RelatedAccountId &&
		left.RelatedAccountAmount == right.RelatedAccountAmount &&
		left.HideAmount == right.HideAmount &&
		left.Comment == right.Comment &&
		left.GeoLongitude == right.GeoLongitude &&
		left.GeoLatitude == right.GeoLatitude
}

func createMCPTransactionInfo(transaction *models.Transaction, accountMap map[int64]*models.Account, services MCPAvailableServices, c *core.WebContext, uid int64, tagNames []string) *MCPTransactionInfo {
	info := &MCPTransactionInfo{
		Id:         utils.Int64ToString(transaction.TransactionId),
		Type:       transactionDbTypeToText(transaction.Type),
		Tags:       tagNames,
		HideAmount: transaction.HideAmount,
		Comment:    transaction.Comment,
		Time: utils.FormatUnixTimeToLongDateTimeWithTimezoneRFC3339Format(
			utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime),
			time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60),
		),
	}

	if account := accountMap[transaction.AccountId]; account != nil {
		info.AccountName = account.Name
		info.Currency = account.Currency

		if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE && account.Category.IsLiability() {
			info.Amount = utils.FormatAmount(-transaction.Amount)
		} else {
			info.Amount = utils.FormatAmount(transaction.Amount)
		}
	} else {
		info.Amount = utils.FormatAmount(transaction.Amount)
	}

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
		info.DestinationAmount = utils.FormatAmount(transaction.RelatedAccountAmount)

		if account := accountMap[transaction.RelatedAccountId]; account != nil {
			info.DestinationAccountName = account.Name
			info.DestinationCurrency = account.Currency
		}
	}

	if transaction.CategoryId > 0 {
		if category, err := services.GetTransactionCategoryService().GetCategoryByCategoryId(c, uid, transaction.CategoryId); err == nil {
			info.SecondaryCategoryName = category.Name
		}
	}

	if transaction.GeoLongitude != 0 || transaction.GeoLatitude != 0 {
		info.GeoLocation = &MCPGeoLocationInfo{
			Latitude:  transaction.GeoLatitude,
			Longitude: transaction.GeoLongitude,
		}
	}

	return info
}
