package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type Message struct {
	Domain    string `json:"domain"`
	Address   string `json:"address"`
	Statement string `json:"statement"`
	Uri       string `json:"uri"`
	Version   string `json:"version"`
	ChainId   int    `json:"chainId"`
	Nonce     string `json:"nonce"`
	IssuedAt  string `json:"issuedAt"`
}

type UserLoginReq struct {
	g.Meta    `path:"/user/login" tags:"User" method:"get" x-group:"user" summary:"用户登陆"`
	Username  string  `json:"username" v:"required|length:42,42#Please Input Wallet Address" dc:"用户名,钱包地址"`
	Message   Message `json:"message" dc:"签名信息"`
	Signature string  `json:"signature" dc:"签名数组"`
}

type UserLoginRes struct {
	Token string `json:"token" dc:"Jwt token: "`
	// Expire string `json:"expire" dc:"TOKEN"`
}

type UserLogoutReq struct {
	g.Meta `path:"/user/logout" tags:"User" method:"post" x-group:"user" summary:"用户退出登陆"`

	User  string `json:"user" v:"required|length:42,42#Please Connect Wallet"`
	Token string `json:"token" v:"required#User not login: Please disconnect you wallet try again"`
}

type UserLogoutRes struct {
}
