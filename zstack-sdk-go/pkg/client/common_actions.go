// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

func (cli *ZSClient) GuestOsTypeList() ([]view.GuestOsTypeView, error) {
	result := new([]view.GuestOsTypeView)
	params := param.HqlParam{
		OperationName: "guestOsTypeList",
		Query: `query guestOsTypeList {
  guestOsTypeList {
    list {
      platform
      children {
        guestName
        children {
          uuid
          platform
          name
          osRelease
          version
          __typename
        }
        __typename
      }
      __typename
    }
    __typename
  }
}`,
	}

	if _, err := cli.hql(params, result, responseKeyData, "guestOsTypeList", "list"); err != nil {
		return nil, err
	}

	return *result, nil
}

func (cli *ZSClient) GuestNameList() ([]view.GuestOsNameView, error) {
	result := new([]view.GuestOsNameView)
	params := param.HqlParam{
		OperationName: "guestNameList",
		Query: `query guestNameList {
  guestNameList {
    name
    __typename
  }
}
`,
	}

	if _, err := cli.hql(params, result, responseKeyData, "guestNameList"); err != nil {
		return nil, err
	}

	return *result, nil
}
