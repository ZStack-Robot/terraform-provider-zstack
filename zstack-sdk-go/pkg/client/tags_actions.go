// Copyright (c) ZStack.io, Inc.

package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// CreateUserTag Create a user tag
func (cli *ZSClient) CreateUserTag(params param.CreateTagParam) (view.UserTagInventoryView, error) {
	var resp view.UserTagInventoryView
	return resp, cli.Post("v1/user-tags", params, &resp)
}

// QueryUserTag Query user tags
func (cli *ZSClient) QueryUserTag(params param.QueryParam) ([]view.UserTagInventoryView, error) {
	var tags []view.UserTagInventoryView
	return tags, cli.List("v1/user-tags", &params, &tags)
}

// QueryUserTag Query all user tags
func (cli *ZSClient) ListAllUserTags() ([]view.UserTagInventoryView, error) {
	params := param.NewQueryParam()
	var tags []view.UserTagInventoryView
	return tags, cli.ListAll("v1/user-tags", &params, &tags)
}

// GetTag Get tag by uuid
func (cli *ZSClient) GetTag(uuid string) (view.TagInventoryView, error) {
	var resp view.TagInventoryView
	return resp, cli.Get("v1/tags/", uuid, nil, &resp)
}

func (cli *ZSClient) GetUserTag(uuid string) ([]view.TagInventory, error) {
	var tags []view.TagInventory
	queryParam := param.NewQueryParam()
	queryParam.AddQ("tagPatternUuid=" + uuid)
	return tags, cli.ListAll("v1/user-tags", &queryParam, &tags)
}

// QueryTag
func (cli *ZSClient) QueryTag(params param.QueryParam) ([]view.TagInventoryView, error) {
	var tags []view.TagInventoryView
	return tags, cli.List("v1/tags", &params, &tags)
}

// QueryTag Query all tags
func (cli *ZSClient) ListAllTags() ([]view.TagInventoryView, error) {
	params := param.NewQueryParam()
	var tags []view.TagInventoryView
	return tags, cli.ListAll("v1/tags", &params, &tags)
}

// CreateTag Create a tag
func (cli *ZSClient) CreateTag(params param.CreateResourceTagParam) (view.TagInventoryView, error) {
	var resp view.TagInventoryView
	return resp, cli.Post("v1/tags", params, &resp)
}

// UpdateTag updates a tag's properties by its UUID.
// Note: Tags with type 'simple' cannot have their values updated.
// Returns the updated tag inventory and any error encountered.
func (cli *ZSClient) UpdateTag(uuid string, params param.UpdateResourceTagParam) (view.TagInventoryView, error) {
	var resp view.TagInventoryView
	return resp, cli.Put("v1/tags", uuid, params, &resp)
}

// AttachTagToResource attaches a tag to one or more resources with optional token support.
//
// Parameters:
//   - tagUuid: The UUID of the tag to be attached
//   - resourceUuids: A slice of resource UUIDs to attach the tag to
//   - attachType: Optional parameters:
//   - First parameter: "withToken" to enable token support
//   - Second parameter: JSON string of tokens (e.g. {"key":"value"})
//
// Returns the result of the attachment operation and any error encountered.
func (cli *ZSClient) AttachTagToResource(tagUuid string, resourceUuids []string, attachType ...string) (view.AttachTagToResourceResult, error) {
	params := param.AttachTagToResourceParam{
		BaseParam: param.BaseParam{},
		Params: param.AttachTagToResourceDetailParam{
			ResourceUuids: resourceUuids,
		},
	}

	if len(attachType) > 0 && attachType[0] == "withToken" {
		if len(attachType) > 1 {
			var tokens map[string]interface{}
			tokenStr := attachType[1]

			//tokenStr = strings.ReplaceAll(tokenStr, "", "\"")

			err := json.Unmarshal([]byte(tokenStr), &tokens)
			if err != nil {
				return view.AttachTagToResourceResult{}, fmt.Errorf("invalid tokens format: %v", err)
			}
			params.Params.Tokens = tokens
		} else {
			return view.AttachTagToResourceResult{}, fmt.Errorf("attachType 'withToken' requires token string")
		}

	}

	var tag view.AttachTagToResourceResult

	err := cli.PostWithRespKey(fmt.Sprintf("v1/tags/%s/resources", tagUuid), "", params, &tag)
	if err != nil {
		return tag, err
	}
	return tag, nil
}

// DetachTagFromResource detaches a tag from one or more resources.
//
// Parameters:
//   - tagUuid: The UUID of the tag to be detached
//   - resourceUuids: A slice of resource UUIDs to detach the tag from
//
// Returns an error if the detachment operation fails.
func (cli *ZSClient) DetachTagFromResource(tagUuid string, resourceUuids []string) error {
	uuidsStr := strings.Join(resourceUuids, ",")

	if err := cli.DeleteWithSpec("v1/tags", tagUuid, "resources", fmt.Sprintf("resourceUuids=%s", uuidsStr), nil); err != nil {
		return err
	}

	return nil
}
