package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type Message struct {
	Domain    string `p:"domain"`
	Address   string `p:"address"`
	Statement string `p:"statement"`
	Type      string `p:"type"`
	ChainId   string `p:"chainId"`
	Nonce     string `p:"nonce"`
	Timestamp string `p:"timestamp"`
}

type UserLoginReq struct {
	g.Meta    `path:"/user/login" tags:"User" method:"post" x-group:"user" summary:"用户登陆"`
	Message   Message `p:"message" v:"required#Please Input messages" dc:"签名信息原文"`
	Signature string  `p:"signature" v:"required#Please Input signature" dc:"签名数据"`
}

type UserLoginRes struct {
	Message string `json:"message" dc:"提示信息"`
	Token   string `json:"token" dc:"uuid token"`
}

type UserLogoutReq struct {
	g.Meta `path:"/user/logout" tags:"User" method:"post" x-group:"user" summary:"用户退出登陆"`

	User  string `json:"user" v:"required|length:42,42#Please Connect Wallet"`
	Token string `json:"token" v:"required#User not login: Please disconnect you wallet try again"`
}

type UserLogoutRes struct {
	Message string `json:"message" dc:"提示信息"`
}

type ActiveUserInfoReq struct {
	g.Meta `path:"/user/active_users_info" tags:"User" method:"get" x-group:"user" summary:"查看站点的活跃度"`
}

type ActiveUserInfoRes struct {
	Message    string `json:"message" dc:"提示信息"`
	TotalUsers int    `json:"total_users" dc:"登录过该站点的用户数量"`
}
