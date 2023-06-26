package wscode

type accountCode struct {

	//注册的操作码
	RegistCreq string "0" //client request //参数 accountDto
	RegistSres string "1" //server response
	//登录的操作码
	Login string "2" //参数 accountDto 账号密码
	//登录失败
	LoginFail      string "3" //参数 accountDto 账号密码
	NICK_NAME_CREQ string "4"
	NICK_NAME_SRES string "5"
}

var AccountCode = &accountCode{
	//注册的操作码
	RegistCreq: "0", //client request //参数 accountDto
	RegistSres: "1", //server response

	//登录的操作码
	Login: "2", //参数 accountDto 账号密码
	//登录失败
	LoginFail: "3",

	NICK_NAME_CREQ: "4", //修改昵称

	NICK_NAME_SRES: "5", //修改昵称
}
