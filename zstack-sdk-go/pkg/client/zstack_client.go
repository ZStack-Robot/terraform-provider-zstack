// Copyright (c) ZStack.io, Inc.

package client

import (
	"crypto/sha512"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/errors"
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

type ZSClient struct {
	*ZSHttpClient
}

func NewZSClient(config *ZSConfig) *ZSClient {
	return &ZSClient{
		ZSHttpClient: NewZSHttpClient(config),
	}
}

func (cli *ZSClient) Login() (*view.SessionView, error) {
	if cli.authType != AuthTypeAccountUser && cli.authType != AuthTypeAccount {
		return nil, errors.ErrNotSupported
	}

	var sessionView *view.SessionView
	var err error
	if cli.authType == AuthTypeAccountUser {
		sessionView, err = cli.logInByAccountUser()
	} else {
		sessionView, err = cli.logInByAccount()
	}

	if err != nil {
		golog.Errorf("ZSClient.Login error:%v", err)
		return nil, err
	}

	cli.LoadSession(sessionView.UUID)
	return sessionView, nil
}

func (cli *ZSClient) logInByAccountUser() (*view.SessionView, error) {
	if cli.authType != AuthTypeAccountUser {
		return nil, errors.ErrNotSupported
	}

	if len(cli.accountName) == 0 || len(cli.accountUserName) == 0 || len(cli.password) == 0 {
		return nil, errors.ErrParameter
	}

	params := param.LogInByUserParam{
		LogInByUser: param.LogInByUserDetailParam{
			AccountName: cli.accountName,
			UserName:    cli.accountUserName,
			Password:    fmt.Sprintf("%x", sha512.Sum512([]byte(cli.password))),
		},
	}
	sessionView := view.SessionView{}
	err := cli.Put("v1/accounts/users/login", "", params, &sessionView)
	if err != nil {
		golog.Errorf("ZSClient.logInByAccountUser Account[%s] User[%s] error:%v",
			cli.accountName, cli.accountUserName, err)
		return nil, err
	}

	return &sessionView, nil
}

func (cli *ZSClient) logInByAccount() (*view.SessionView, error) {
	if cli.authType != AuthTypeAccount {
		return nil, errors.ErrNotSupported
	}

	if len(cli.accountName) == 0 || len(cli.password) == 0 {
		return nil, errors.ErrParameter
	}

	params := param.LoginByAccountParam{
		LoginByAccount: param.LoginByAccountDetailParam{
			AccountName: cli.accountName,
			Password:    fmt.Sprintf("%x", sha512.Sum512([]byte(cli.password))),
		},
	}
	sessionView := view.SessionView{}
	err := cli.Put("v1/accounts/login", "", params, &sessionView)
	if err != nil {
		golog.Errorf("ZSClient.logInByAccount Account[%s] error:%v", cli.accountName, err)
		return nil, err
	}

	return &sessionView, nil
}

func (cli *ZSClient) ValidateSession() (map[string]bool, error) {
	if cli.authType != AuthTypeAccountUser && cli.authType != AuthTypeAccount {
		return nil, errors.ErrNotSupported
	}

	if len(cli.sessionId) == 0 {
		return nil, errors.ErrNotSupported
	}

	return cli.ValidateSessionId(cli.sessionId)
}

func (cli *ZSClient) ValidateSessionId(sessionId string) (map[string]bool, error) {
	validSession := make(map[string]bool)
	err := cli.GetWithSpec("v1/accounts/sessions", sessionId, "valid", "", nil, &validSession)
	if err != nil {
		golog.Errorf("ZSClient.ValidateSession sessionId[%s] error:%v", sessionId, err)
		return nil, err
	}

	golog.Debugf("ZSClient.ValidateSession sessionId[%s]:%v", sessionId, validSession)
	return validSession, nil
}

func (cli *ZSClient) Logout() error {
	if cli.authType != AuthTypeAccountUser && cli.authType != AuthTypeAccount {
		return errors.ErrNotSupported
	}

	if len(cli.sessionId) == 0 {
		return errors.ErrNotSupported
	}

	err := cli.Delete("v1/accounts/sessions", cli.sessionId, "")
	if err != nil {
		golog.Errorf("ZSClient.Logout sessionId[%s] error:%v", cli.sessionId, err)
		return err
	}

	cli.unloadSession()
	return nil
}

func (cli *ZSClient) WebLogin() (*view.WebUISessionView, error) {
	if cli.authType != AuthTypeAccountUser && cli.authType != AuthTypeAccount {
		return nil, errors.ErrNotSupported
	}

	var operationName, username, loginType, query string
	var input map[string]interface{}
	if cli.authType == AuthTypeAccount {
		operationName, username, loginType = "loginByAccount", cli.accountName, "iam1"
		input = map[string]interface{}{
			"accountName": cli.accountName,
			"password":    fmt.Sprintf("%x", sha512.Sum512([]byte(cli.password))),
		}
		query = `mutation loginByAccount($input:LoginByAccountInput!) { 
			loginByAccount(input: $input) { 
			  sessionId,
			  accountUuid,
			  userUuid,
			  currentIdentity
			}
		  }`
	} else {
		operationName, username, loginType = "loginIAM2VirtualID", cli.accountUserName, "iam2"
		input = map[string]interface{}{
			"name":     cli.accountUserName,
			"password": fmt.Sprintf("%x", sha512.Sum512([]byte(cli.password))),
		}
		query = `mutation loginIAM2VirtualID($input:LoginIAM2VirtualIDInput!) { 
			loginIAM2VirtualID(input: $input) { 
			  sessionId,
			  accountUuid,
			  userUuid,
			  currentIdentity
			}
		  }`
	}

	result := new(view.WebUISessionView)
	params := param.HqlParam{
		OperationName: operationName,
		Query:         query,
		Variables: param.Variables{
			Input: input,
		},
	}
	respHeader, err := cli.hql(params, result, responseKeyData, operationName)
	if err != nil {
		return nil, err
	}
	result.UserName = username
	result.LoginType = loginType
	result.ZSVersion = respHeader.Get("Zs-Version")
	return result, nil
}

func (cli *ZSClient) hql(params param.HqlParam, retVal interface{}, unMarshalKeys ...string) (http.Header, error) {
	urlStr := fmt.Sprintf("http://%s:%d/graphql", cli.hostname, WebZStackPort)
	_, respHeader, resp, err := cli.httpPost(urlStr, jsonMarshal(params), false)
	if err != nil {
		return nil, err
	}

	if retVal == nil {
		return nil, nil
	}

	return respHeader, resp.Unmarshal(retVal, unMarshalKeys...)
}

func (cli *ZSClient) Zql(querySt string, retVal interface{}, unMarshalKeys ...string) (http.Header, error) {
	encodedQuery := url.QueryEscape(querySt)
	baseUrl := cli.getRequestURL("v1/zql")
	urlStr := fmt.Sprintf("%s?zql=%s", baseUrl, encodedQuery)
	_, respHeader, resp, err := cli.httpGet(urlStr, false)
	if err != nil {
		return nil, err
	}

	if retVal == nil {
		return nil, nil
	}

	return respHeader, resp.Unmarshal(retVal, unMarshalKeys...)
}
