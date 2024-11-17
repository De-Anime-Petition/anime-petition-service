// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure for table users.
type Users struct {
	Id         int         `json:"id"         orm:"id"          ` //
	Wallet     string      `json:"wallet"     orm:"wallet"      ` //
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" ` //
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" ` //
	Token      string      `json:"token"      orm:"token"       ` //
}
