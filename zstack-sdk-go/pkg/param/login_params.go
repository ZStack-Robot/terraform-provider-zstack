package param

type LoginByAccountParam struct {
	BaseParam

	LoginByAccount LoginByAccountDetailParam `json:"logInByAccount"`
}

type LoginByAccountDetailParam struct {
	AccountName string                 `json:"accountName"` //账户名称
	Password    string                 `json:"password"`    //密码
	AccountType string                 `json:"accountType"` //账户类型
	CaptchaUuid string                 `json:"captchaUuid"` //验证码UUID
	VerifyCode  string                 `json:"verifyCode"`  //验证码
	ClientInfo  map[string]interface{} `json:"clientInfo"`  //客户端信息
}

type LogInByUserParam struct {
	BaseParam

	LogInByUser LogInByUserDetailParam `json:"logInByUser"`
}

type LogInByUserDetailParam struct {
	AccountUuid string                 `json:"accountUuid"` //账户UUID
	AccountName string                 `json:"accountName"` //账户名称
	UserName    string                 `json:"userName"`    //用户名称
	Password    string                 `json:"password"`    //密码
	ClientInfo  map[string]interface{} `json:"clientInfo"`  //客户端信息
}