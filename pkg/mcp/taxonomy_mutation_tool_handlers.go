package mcp

import (
	"encoding/json"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// MCPTransactionCategoryInfo is a mutation-friendly transaction category record.
type MCPTransactionCategoryInfo struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	ParentId     string `json:"parent_id,omitempty"`
	Type         string `json:"type"`
	Icon         string `json:"icon"`
	Color        string `json:"color"`
	Comment      string `json:"comment,omitempty"`
	DisplayOrder int32  `json:"display_order"`
	Hidden       bool   `json:"hidden"`
}

// MCPCreateTransactionCategoryRequest represents category creation.
type MCPCreateTransactionCategoryRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type" jsonschema:"enum=income,enum=expense,enum=transfer"`
	ParentId string `json:"parent_id,omitempty" jsonschema_description:"Optional primary category id"`
	Icon     int64  `json:"icon,omitempty" jsonschema_description:"Icon id, defaults to 1"`
	Color    string `json:"color,omitempty" jsonschema_description:"Six-digit RGB hex without #, defaults to 176b5b"`
	Comment  string `json:"comment,omitempty"`
	DryRun   bool   `json:"dry_run,omitempty"`
}

// MCPUpdateTransactionCategoryRequest represents a partial category update.
type MCPUpdateTransactionCategoryRequest struct {
	Id       string  `json:"id"`
	Name     *string `json:"name,omitempty"`
	ParentId *string `json:"parent_id,omitempty" jsonschema_description:"Primary category id; empty string means top level"`
	Icon     *int64  `json:"icon,omitempty"`
	Color    *string `json:"color,omitempty"`
	Comment  *string `json:"comment,omitempty"`
	Hidden   *bool   `json:"hidden,omitempty"`
	DryRun   bool    `json:"dry_run,omitempty"`
}

// MCPTransactionCategoryMutationResponse represents a category mutation result.
type MCPTransactionCategoryMutationResponse struct {
	Success  bool                        `json:"success"`
	DryRun   bool                        `json:"dry_run,omitempty"`
	Category *MCPTransactionCategoryInfo `json:"category"`
}

type mcpCreateTransactionCategoryToolHandler struct{}

// MCPCreateTransactionCategoryToolHandler creates transaction categories.
var MCPCreateTransactionCategoryToolHandler = &mcpCreateTransactionCategoryToolHandler{}

func (h *mcpCreateTransactionCategoryToolHandler) Name() string {
	return "create_transaction_category"
}

func (h *mcpCreateTransactionCategoryToolHandler) Description() string {
	return "Create a primary or secondary income, expense, or transfer category."
}

func (h *mcpCreateTransactionCategoryToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPCreateTransactionCategoryRequest{})
}

func (h *mcpCreateTransactionCategoryToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPTransactionCategoryMutationResponse{})
}

func (h *mcpCreateTransactionCategoryToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPCreateTransactionCategoryRequest
	if err := decodeMCPArguments(callToolReq, &request); err != nil {
		return nil, nil, err
	}

	categoryType, err := textualTransactionCategoryTypeToModel(request.Type)
	if err != nil {
		return nil, nil, err
	}

	icon := request.Icon
	if icon <= 0 {
		icon = 1
	}
	color := strings.ToLower(request.Color)
	if color == "" {
		color = "176b5b"
	}
	if err = validateMCPCategoryFields(request.Name, request.Comment, icon, color); err != nil {
		return nil, nil, err
	}

	parentId, err := parseMCPOptionalId(request.ParentId, errs.ErrTransactionCategoryIdInvalid)
	if err != nil {
		return nil, nil, err
	}
	if parentId > 0 {
		parent, parentErr := services.GetTransactionCategoryService().GetCategoryByCategoryId(c, user.Uid, parentId)
		if parentErr != nil {
			return nil, nil, parentErr
		}
		if parent.ParentCategoryId != models.LevelOneTransactionCategoryParentId {
			return nil, nil, errs.ErrCannotAddToSecondaryTransactionCategory
		}
		if parent.Type != categoryType {
			return nil, nil, errs.ErrNotAllowChangePrimaryTransactionType
		}
	}

	var maxOrder int32
	if parentId > 0 {
		maxOrder, err = services.GetTransactionCategoryService().GetMaxSubCategoryDisplayOrder(c, user.Uid, categoryType, parentId)
	} else {
		maxOrder, err = services.GetTransactionCategoryService().GetMaxDisplayOrder(c, user.Uid, categoryType)
	}
	if err != nil {
		return nil, nil, err
	}

	category := &models.TransactionCategory{
		Uid:              user.Uid,
		Type:             categoryType,
		ParentCategoryId: parentId,
		Name:             strings.TrimSpace(request.Name),
		DisplayOrder:     maxOrder + 1,
		Icon:             icon,
		Color:            color,
		Comment:          request.Comment,
	}
	if !request.DryRun {
		if err = services.GetTransactionCategoryService().CreateCategory(c, category); err != nil {
			return nil, nil, err
		}
		log.Infof(c, "[create_transaction_category.Handle] user \"uid:%d\" created category \"id:%d\"", user.Uid, category.CategoryId)
	}
	return newMCPTransactionCategoryMutationResponse(category, request.DryRun)
}

type mcpUpdateTransactionCategoryToolHandler struct{}

// MCPUpdateTransactionCategoryToolHandler updates transaction categories.
var MCPUpdateTransactionCategoryToolHandler = &mcpUpdateTransactionCategoryToolHandler{}

func (h *mcpUpdateTransactionCategoryToolHandler) Name() string {
	return "update_transaction_category"
}

func (h *mcpUpdateTransactionCategoryToolHandler) Description() string {
	return "Update a category name, parent, icon, color, description, or visibility."
}

func (h *mcpUpdateTransactionCategoryToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPUpdateTransactionCategoryRequest{})
}

func (h *mcpUpdateTransactionCategoryToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPTransactionCategoryMutationResponse{})
}

func (h *mcpUpdateTransactionCategoryToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPUpdateTransactionCategoryRequest
	if err := decodeMCPArguments(callToolReq, &request); err != nil {
		return nil, nil, err
	}
	categoryId, err := parseMCPRequiredId(request.Id, errs.ErrTransactionCategoryIdInvalid)
	if err != nil {
		return nil, nil, err
	}
	category, err := services.GetTransactionCategoryService().GetCategoryByCategoryId(c, user.Uid, categoryId)
	if err != nil {
		return nil, nil, err
	}

	updated := *category
	if request.Name != nil {
		updated.Name = strings.TrimSpace(*request.Name)
	}
	if request.Icon != nil {
		updated.Icon = *request.Icon
	}
	if request.Color != nil {
		updated.Color = strings.ToLower(*request.Color)
	}
	if request.Comment != nil {
		updated.Comment = *request.Comment
	}
	if request.Hidden != nil {
		updated.Hidden = *request.Hidden
	}
	if err = validateMCPCategoryFields(updated.Name, updated.Comment, updated.Icon, updated.Color); err != nil {
		return nil, nil, err
	}

	if request.ParentId != nil {
		updated.ParentCategoryId, err = parseMCPOptionalId(*request.ParentId, errs.ErrTransactionCategoryIdInvalid)
		if err != nil {
			return nil, nil, err
		}
		if category.ParentCategoryId == 0 && updated.ParentCategoryId != 0 {
			return nil, nil, errs.ErrNotAllowChangePrimaryTransactionCategoryToSecondary
		}
		if category.ParentCategoryId != 0 && updated.ParentCategoryId == 0 {
			return nil, nil, errs.ErrNotAllowChangeSecondaryTransactionCategoryToPrimary
		}
		if updated.ParentCategoryId != category.ParentCategoryId {
			parent, parentErr := services.GetTransactionCategoryService().GetCategoryByCategoryId(c, user.Uid, updated.ParentCategoryId)
			if parentErr != nil {
				return nil, nil, parentErr
			}
			if parent.ParentCategoryId != models.LevelOneTransactionCategoryParentId {
				return nil, nil, errs.ErrNotAllowUseSecondaryTransactionAsPrimaryCategory
			}
			if parent.Type != category.Type {
				return nil, nil, errs.ErrNotAllowChangePrimaryTransactionType
			}
			updated.DisplayOrder, err = services.GetTransactionCategoryService().GetMaxSubCategoryDisplayOrder(c, user.Uid, category.Type, updated.ParentCategoryId)
			if err != nil {
				return nil, nil, err
			}
			updated.DisplayOrder++
		}
	}

	if transactionCategoryModelsEqualForMCP(category, &updated) {
		return nil, nil, errs.ErrNothingWillBeUpdated
	}
	if !request.DryRun {
		if err = services.GetTransactionCategoryService().ModifyCategory(c, &updated); err != nil {
			return nil, nil, err
		}
		log.Infof(c, "[update_transaction_category.Handle] user \"uid:%d\" updated category \"id:%d\"", user.Uid, categoryId)
	}
	return newMCPTransactionCategoryMutationResponse(&updated, request.DryRun)
}

// MCPDeleteTaxonomyRequest is shared by category, tag, and tag-group deletion.
type MCPDeleteTaxonomyRequest struct {
	Id     string `json:"id"`
	DryRun bool   `json:"dry_run,omitempty"`
}

// MCPDeleteTaxonomyResponse is a taxonomy deletion result.
type MCPDeleteTaxonomyResponse struct {
	Success bool   `json:"success"`
	DryRun  bool   `json:"dry_run,omitempty"`
	Id      string `json:"id"`
}

type mcpDeleteTransactionCategoryToolHandler struct{}

// MCPDeleteTransactionCategoryToolHandler deletes unused transaction categories.
var MCPDeleteTransactionCategoryToolHandler = &mcpDeleteTransactionCategoryToolHandler{}

func (h *mcpDeleteTransactionCategoryToolHandler) Name() string {
	return "delete_transaction_category"
}

func (h *mcpDeleteTransactionCategoryToolHandler) Description() string {
	return "Delete an unused transaction category. In-use categories remain protected."
}

func (h *mcpDeleteTransactionCategoryToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPDeleteTaxonomyRequest{})
}

func (h *mcpDeleteTransactionCategoryToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPDeleteTaxonomyResponse{})
}

func (h *mcpDeleteTransactionCategoryToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPDeleteTaxonomyRequest
	if err := decodeMCPArguments(callToolReq, &request); err != nil {
		return nil, nil, err
	}
	id, err := parseMCPRequiredId(request.Id, errs.ErrTransactionCategoryIdInvalid)
	if err != nil {
		return nil, nil, err
	}
	if _, err = services.GetTransactionCategoryService().GetCategoryByCategoryId(c, user.Uid, id); err != nil {
		return nil, nil, err
	}
	if !request.DryRun {
		if err = services.GetTransactionCategoryService().DeleteCategory(c, user.Uid, id); err != nil {
			return nil, nil, err
		}
		log.Infof(c, "[delete_transaction_category.Handle] user \"uid:%d\" deleted category \"id:%d\"", user.Uid, id)
	}
	return newMCPDeleteTaxonomyResponse(request.Id, request.DryRun)
}

// MCPTransactionTagInfo is a mutation-friendly tag record.
type MCPTransactionTagInfo struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	GroupId      string `json:"group_id,omitempty"`
	DisplayOrder int32  `json:"display_order"`
	Hidden       bool   `json:"hidden"`
}

// MCPCreateTransactionTagRequest represents tag creation.
type MCPCreateTransactionTagRequest struct {
	Name    string `json:"name"`
	GroupId string `json:"group_id,omitempty"`
	DryRun  bool   `json:"dry_run,omitempty"`
}

// MCPUpdateTransactionTagRequest represents a partial tag update.
type MCPUpdateTransactionTagRequest struct {
	Id      string  `json:"id"`
	Name    *string `json:"name,omitempty"`
	GroupId *string `json:"group_id,omitempty" jsonschema_description:"Empty string moves the tag out of a group"`
	Hidden  *bool   `json:"hidden,omitempty"`
	DryRun  bool    `json:"dry_run,omitempty"`
}

// MCPTransactionTagMutationResponse represents a tag mutation result.
type MCPTransactionTagMutationResponse struct {
	Success bool                   `json:"success"`
	DryRun  bool                   `json:"dry_run,omitempty"`
	Tag     *MCPTransactionTagInfo `json:"tag"`
}

type mcpCreateTransactionTagToolHandler struct{}

// MCPCreateTransactionTagToolHandler creates transaction tags.
var MCPCreateTransactionTagToolHandler = &mcpCreateTransactionTagToolHandler{}

func (h *mcpCreateTransactionTagToolHandler) Name() string {
	return "create_transaction_tag"
}

func (h *mcpCreateTransactionTagToolHandler) Description() string {
	return "Create a transaction tag, optionally inside a tag group."
}

func (h *mcpCreateTransactionTagToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPCreateTransactionTagRequest{})
}

func (h *mcpCreateTransactionTagToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPTransactionTagMutationResponse{})
}

func (h *mcpCreateTransactionTagToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPCreateTransactionTagRequest
	if err := decodeMCPArguments(callToolReq, &request); err != nil {
		return nil, nil, err
	}
	name := strings.TrimSpace(request.Name)
	if err := validateMCPTagName(name); err != nil {
		return nil, nil, err
	}
	groupId, err := validateMCPTagGroup(c, user.Uid, request.GroupId, services)
	if err != nil {
		return nil, nil, err
	}
	maxOrder, err := services.GetTransactionTagService().GetMaxDisplayOrder(c, user.Uid, groupId)
	if err != nil {
		return nil, nil, err
	}
	tag := &models.TransactionTag{
		Uid:          user.Uid,
		Name:         name,
		TagGroupId:   groupId,
		DisplayOrder: maxOrder + 1,
	}
	if !request.DryRun {
		if err = services.GetTransactionTagService().CreateTag(c, tag); err != nil {
			return nil, nil, err
		}
		log.Infof(c, "[create_transaction_tag.Handle] user \"uid:%d\" created tag \"id:%d\"", user.Uid, tag.TagId)
	}
	return newMCPTransactionTagMutationResponse(tag, request.DryRun)
}

type mcpUpdateTransactionTagToolHandler struct{}

// MCPUpdateTransactionTagToolHandler updates transaction tags.
var MCPUpdateTransactionTagToolHandler = &mcpUpdateTransactionTagToolHandler{}

func (h *mcpUpdateTransactionTagToolHandler) Name() string {
	return "update_transaction_tag"
}

func (h *mcpUpdateTransactionTagToolHandler) Description() string {
	return "Update a transaction tag name, group, or visibility."
}

func (h *mcpUpdateTransactionTagToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPUpdateTransactionTagRequest{})
}

func (h *mcpUpdateTransactionTagToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPTransactionTagMutationResponse{})
}

func (h *mcpUpdateTransactionTagToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPUpdateTransactionTagRequest
	if err := decodeMCPArguments(callToolReq, &request); err != nil {
		return nil, nil, err
	}
	tagId, err := parseMCPRequiredId(request.Id, errs.ErrTransactionTagIdInvalid)
	if err != nil {
		return nil, nil, err
	}
	tag, err := services.GetTransactionTagService().GetTagByTagId(c, user.Uid, tagId)
	if err != nil {
		return nil, nil, err
	}
	updated := *tag
	if request.Name != nil {
		updated.Name = strings.TrimSpace(*request.Name)
	}
	if err = validateMCPTagName(updated.Name); err != nil {
		return nil, nil, err
	}
	if request.GroupId != nil {
		updated.TagGroupId, err = validateMCPTagGroup(c, user.Uid, *request.GroupId, services)
		if err != nil {
			return nil, nil, err
		}
		if updated.TagGroupId != tag.TagGroupId {
			updated.DisplayOrder, err = services.GetTransactionTagService().GetMaxDisplayOrder(c, user.Uid, updated.TagGroupId)
			if err != nil {
				return nil, nil, err
			}
			updated.DisplayOrder++
		}
	}
	if request.Hidden != nil {
		updated.Hidden = *request.Hidden
	}

	nameChanged := updated.Name != tag.Name
	metadataChanged := nameChanged || updated.TagGroupId != tag.TagGroupId
	hiddenChanged := updated.Hidden != tag.Hidden
	if !metadataChanged && !hiddenChanged {
		return nil, nil, errs.ErrNothingWillBeUpdated
	}
	if !request.DryRun {
		if metadataChanged {
			if err = services.GetTransactionTagService().ModifyTag(c, &updated, nameChanged); err != nil {
				return nil, nil, err
			}
		}
		if hiddenChanged {
			if err = services.GetTransactionTagService().HideTag(c, user.Uid, []int64{tagId}, updated.Hidden); err != nil {
				return nil, nil, err
			}
		}
		log.Infof(c, "[update_transaction_tag.Handle] user \"uid:%d\" updated tag \"id:%d\"", user.Uid, tagId)
	}
	return newMCPTransactionTagMutationResponse(&updated, request.DryRun)
}

type mcpDeleteTransactionTagToolHandler struct{}

// MCPDeleteTransactionTagToolHandler deletes unused transaction tags.
var MCPDeleteTransactionTagToolHandler = &mcpDeleteTransactionTagToolHandler{}

func (h *mcpDeleteTransactionTagToolHandler) Name() string {
	return "delete_transaction_tag"
}

func (h *mcpDeleteTransactionTagToolHandler) Description() string {
	return "Delete an unused transaction tag. In-use tags remain protected."
}

func (h *mcpDeleteTransactionTagToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPDeleteTaxonomyRequest{})
}

func (h *mcpDeleteTransactionTagToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPDeleteTaxonomyResponse{})
}

func (h *mcpDeleteTransactionTagToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPDeleteTaxonomyRequest
	if err := decodeMCPArguments(callToolReq, &request); err != nil {
		return nil, nil, err
	}
	id, err := parseMCPRequiredId(request.Id, errs.ErrTransactionTagIdInvalid)
	if err != nil {
		return nil, nil, err
	}
	if _, err = services.GetTransactionTagService().GetTagByTagId(c, user.Uid, id); err != nil {
		return nil, nil, err
	}
	if !request.DryRun {
		if err = services.GetTransactionTagService().DeleteTag(c, user.Uid, id); err != nil {
			return nil, nil, err
		}
		log.Infof(c, "[delete_transaction_tag.Handle] user \"uid:%d\" deleted tag \"id:%d\"", user.Uid, id)
	}
	return newMCPDeleteTaxonomyResponse(request.Id, request.DryRun)
}

// MCPTransactionTagGroupInfo is a tag group record.
type MCPTransactionTagGroupInfo struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	DisplayOrder int32  `json:"display_order"`
}

// MCPCreateTransactionTagGroupRequest represents tag group creation.
type MCPCreateTransactionTagGroupRequest struct {
	Name   string `json:"name"`
	DryRun bool   `json:"dry_run,omitempty"`
}

// MCPUpdateTransactionTagGroupRequest represents tag group modification.
type MCPUpdateTransactionTagGroupRequest struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	DryRun bool   `json:"dry_run,omitempty"`
}

// MCPTransactionTagGroupMutationResponse represents a tag group mutation result.
type MCPTransactionTagGroupMutationResponse struct {
	Success bool                        `json:"success"`
	DryRun  bool                        `json:"dry_run,omitempty"`
	Group   *MCPTransactionTagGroupInfo `json:"group"`
}

type mcpCreateTransactionTagGroupToolHandler struct{}

// MCPCreateTransactionTagGroupToolHandler creates transaction tag groups.
var MCPCreateTransactionTagGroupToolHandler = &mcpCreateTransactionTagGroupToolHandler{}

func (h *mcpCreateTransactionTagGroupToolHandler) Name() string {
	return "create_transaction_tag_group"
}

func (h *mcpCreateTransactionTagGroupToolHandler) Description() string {
	return "Create a transaction tag group."
}

func (h *mcpCreateTransactionTagGroupToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPCreateTransactionTagGroupRequest{})
}

func (h *mcpCreateTransactionTagGroupToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPTransactionTagGroupMutationResponse{})
}

func (h *mcpCreateTransactionTagGroupToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPCreateTransactionTagGroupRequest
	if err := decodeMCPArguments(callToolReq, &request); err != nil {
		return nil, nil, err
	}
	name := strings.TrimSpace(request.Name)
	if err := validateMCPTagName(name); err != nil {
		return nil, nil, err
	}
	maxOrder, err := services.GetTransactionTagGroupService().GetMaxDisplayOrder(c, user.Uid)
	if err != nil {
		return nil, nil, err
	}
	group := &models.TransactionTagGroup{Uid: user.Uid, Name: name, DisplayOrder: maxOrder + 1}
	if !request.DryRun {
		if err = services.GetTransactionTagGroupService().CreateTagGroup(c, group); err != nil {
			return nil, nil, err
		}
		log.Infof(c, "[create_transaction_tag_group.Handle] user \"uid:%d\" created tag group \"id:%d\"", user.Uid, group.TagGroupId)
	}
	return newMCPTransactionTagGroupMutationResponse(group, request.DryRun)
}

type mcpUpdateTransactionTagGroupToolHandler struct{}

// MCPUpdateTransactionTagGroupToolHandler updates transaction tag groups.
var MCPUpdateTransactionTagGroupToolHandler = &mcpUpdateTransactionTagGroupToolHandler{}

func (h *mcpUpdateTransactionTagGroupToolHandler) Name() string {
	return "update_transaction_tag_group"
}

func (h *mcpUpdateTransactionTagGroupToolHandler) Description() string {
	return "Rename a transaction tag group."
}

func (h *mcpUpdateTransactionTagGroupToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPUpdateTransactionTagGroupRequest{})
}

func (h *mcpUpdateTransactionTagGroupToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPTransactionTagGroupMutationResponse{})
}

func (h *mcpUpdateTransactionTagGroupToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPUpdateTransactionTagGroupRequest
	if err := decodeMCPArguments(callToolReq, &request); err != nil {
		return nil, nil, err
	}
	id, err := parseMCPRequiredId(request.Id, errs.ErrTransactionTagGroupIdInvalid)
	if err != nil {
		return nil, nil, err
	}
	group, err := services.GetTransactionTagGroupService().GetTagGroupByTagGroupId(c, user.Uid, id)
	if err != nil {
		return nil, nil, err
	}
	name := strings.TrimSpace(request.Name)
	if err = validateMCPTagName(name); err != nil {
		return nil, nil, err
	}
	if name == group.Name {
		return nil, nil, errs.ErrNothingWillBeUpdated
	}
	updated := *group
	updated.Name = name
	if !request.DryRun {
		if err = services.GetTransactionTagGroupService().ModifyTagGroup(c, &updated); err != nil {
			return nil, nil, err
		}
		log.Infof(c, "[update_transaction_tag_group.Handle] user \"uid:%d\" updated tag group \"id:%d\"", user.Uid, id)
	}
	return newMCPTransactionTagGroupMutationResponse(&updated, request.DryRun)
}

type mcpDeleteTransactionTagGroupToolHandler struct{}

// MCPDeleteTransactionTagGroupToolHandler deletes empty transaction tag groups.
var MCPDeleteTransactionTagGroupToolHandler = &mcpDeleteTransactionTagGroupToolHandler{}

func (h *mcpDeleteTransactionTagGroupToolHandler) Name() string {
	return "delete_transaction_tag_group"
}

func (h *mcpDeleteTransactionTagGroupToolHandler) Description() string {
	return "Delete an empty transaction tag group."
}

func (h *mcpDeleteTransactionTagGroupToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPDeleteTaxonomyRequest{})
}

func (h *mcpDeleteTransactionTagGroupToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPDeleteTaxonomyResponse{})
}

func (h *mcpDeleteTransactionTagGroupToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, user *models.User, currentConfig *settings.Config, services MCPAvailableServices) (any, []*MCPTextContent, error) {
	var request MCPDeleteTaxonomyRequest
	if err := decodeMCPArguments(callToolReq, &request); err != nil {
		return nil, nil, err
	}
	id, err := parseMCPRequiredId(request.Id, errs.ErrTransactionTagGroupIdInvalid)
	if err != nil {
		return nil, nil, err
	}
	if _, err = services.GetTransactionTagGroupService().GetTagGroupByTagGroupId(c, user.Uid, id); err != nil {
		return nil, nil, err
	}
	if !request.DryRun {
		if err = services.GetTransactionTagGroupService().DeleteTagGroup(c, user.Uid, id); err != nil {
			return nil, nil, err
		}
		log.Infof(c, "[delete_transaction_tag_group.Handle] user \"uid:%d\" deleted tag group \"id:%d\"", user.Uid, id)
	}
	return newMCPDeleteTaxonomyResponse(request.Id, request.DryRun)
}

func textualTransactionCategoryTypeToModel(categoryType string) (models.TransactionCategoryType, error) {
	switch categoryType {
	case "income":
		return models.CATEGORY_TYPE_INCOME, nil
	case "expense":
		return models.CATEGORY_TYPE_EXPENSE, nil
	case "transfer":
		return models.CATEGORY_TYPE_TRANSFER, nil
	default:
		return 0, errs.ErrTransactionCategoryTypeInvalid
	}
}

func transactionCategoryTypeToText(categoryType models.TransactionCategoryType) string {
	switch categoryType {
	case models.CATEGORY_TYPE_INCOME:
		return "income"
	case models.CATEGORY_TYPE_EXPENSE:
		return "expense"
	case models.CATEGORY_TYPE_TRANSFER:
		return "transfer"
	default:
		return ""
	}
}

func validateMCPCategoryFields(name string, comment string, icon int64, color string) error {
	if strings.TrimSpace(name) == "" || utf8.RuneCountInString(name) > 64 ||
		utf8.RuneCountInString(comment) > 255 || icon <= 0 || !mcpHexColorPattern.MatchString(color) {
		return errs.ErrIncompleteOrIncorrectSubmission
	}
	return nil
}

func validateMCPTagName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errs.ErrTransactionTagNameIsEmpty
	}
	if utf8.RuneCountInString(name) > 64 {
		return errs.ErrIncompleteOrIncorrectSubmission
	}
	return nil
}

func validateMCPTagGroup(c *core.WebContext, uid int64, id string, services MCPAvailableServices) (int64, error) {
	groupId, err := parseMCPOptionalId(id, errs.ErrTransactionTagGroupIdInvalid)
	if err != nil {
		return 0, err
	}
	if groupId > 0 {
		if _, err = services.GetTransactionTagGroupService().GetTagGroupByTagGroupId(c, uid, groupId); err != nil {
			return 0, err
		}
	}
	return groupId, nil
}

func parseMCPRequiredId(id string, invalidError error) (int64, error) {
	parsed, err := utils.StringToInt64(id)
	if err != nil || parsed <= 0 {
		return 0, invalidError
	}
	return parsed, nil
}

func parseMCPOptionalId(id string, invalidError error) (int64, error) {
	if id == "" {
		return 0, nil
	}
	return parseMCPRequiredId(id, invalidError)
}

func decodeMCPArguments(callToolReq *MCPCallToolRequest, output any) error {
	if callToolReq.Arguments == nil {
		return errs.ErrIncompleteOrIncorrectSubmission
	}
	if err := json.Unmarshal(callToolReq.Arguments, output); err != nil {
		return errs.NewIncompleteOrIncorrectSubmissionError(err)
	}
	return nil
}

func transactionCategoryModelsEqualForMCP(left *models.TransactionCategory, right *models.TransactionCategory) bool {
	return left.ParentCategoryId == right.ParentCategoryId &&
		left.Name == right.Name &&
		left.DisplayOrder == right.DisplayOrder &&
		left.Icon == right.Icon &&
		left.Color == right.Color &&
		left.Comment == right.Comment &&
		left.Hidden == right.Hidden
}

func createMCPTransactionCategoryInfo(category *models.TransactionCategory) *MCPTransactionCategoryInfo {
	info := &MCPTransactionCategoryInfo{
		Id:           utils.Int64ToString(category.CategoryId),
		Name:         category.Name,
		Type:         transactionCategoryTypeToText(category.Type),
		Icon:         utils.Int64ToString(category.Icon),
		Color:        category.Color,
		Comment:      category.Comment,
		DisplayOrder: category.DisplayOrder,
		Hidden:       category.Hidden,
	}
	if category.ParentCategoryId > 0 {
		info.ParentId = utils.Int64ToString(category.ParentCategoryId)
	}
	return info
}

func createMCPTransactionTagInfo(tag *models.TransactionTag) *MCPTransactionTagInfo {
	info := &MCPTransactionTagInfo{
		Id:           utils.Int64ToString(tag.TagId),
		Name:         tag.Name,
		DisplayOrder: tag.DisplayOrder,
		Hidden:       tag.Hidden,
	}
	if tag.TagGroupId > 0 {
		info.GroupId = utils.Int64ToString(tag.TagGroupId)
	}
	return info
}

func createMCPTransactionTagGroupInfo(group *models.TransactionTagGroup) *MCPTransactionTagGroupInfo {
	return &MCPTransactionTagGroupInfo{
		Id:           utils.Int64ToString(group.TagGroupId),
		Name:         group.Name,
		DisplayOrder: group.DisplayOrder,
	}
}

func newMCPTransactionCategoryMutationResponse(category *models.TransactionCategory, dryRun bool) (any, []*MCPTextContent, error) {
	response := MCPTransactionCategoryMutationResponse{Success: true, DryRun: dryRun, Category: createMCPTransactionCategoryInfo(category)}
	return newMCPJSONResponse(response)
}

func newMCPTransactionTagMutationResponse(tag *models.TransactionTag, dryRun bool) (any, []*MCPTextContent, error) {
	response := MCPTransactionTagMutationResponse{Success: true, DryRun: dryRun, Tag: createMCPTransactionTagInfo(tag)}
	return newMCPJSONResponse(response)
}

func newMCPTransactionTagGroupMutationResponse(group *models.TransactionTagGroup, dryRun bool) (any, []*MCPTextContent, error) {
	response := MCPTransactionTagGroupMutationResponse{Success: true, DryRun: dryRun, Group: createMCPTransactionTagGroupInfo(group)}
	return newMCPJSONResponse(response)
}

func newMCPDeleteTaxonomyResponse(id string, dryRun bool) (any, []*MCPTextContent, error) {
	response := MCPDeleteTaxonomyResponse{Success: true, DryRun: dryRun, Id: id}
	return newMCPJSONResponse(response)
}

func newMCPJSONResponse(response any) (any, []*MCPTextContent, error) {
	content, err := json.Marshal(response)
	if err != nil {
		return nil, nil, err
	}
	return response, []*MCPTextContent{NewMCPTextContent(string(content))}, nil
}
