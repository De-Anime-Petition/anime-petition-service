package controller

import (
	v1 "anime_petition/api/v1"
	"anime_petition/internal/dao"
	"anime_petition/internal/model/entity"
	"anime_petition/utility"
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

var User = &cUser{}

type cUser struct{}

func (c *cUser) Login(ctx context.Context, req *v1.UserLoginReq) (res *v1.UserLoginRes, err error) {
	message := req.Message.Domain + req.Message.Address + req.Message.Statement + req.Message.Type + req.Message.ChainId + req.Message.Nonce + req.Message.Timestamp

	_, err = utility.VerifyEIP191Signature(req.Message.Address, message, req.Signature)
	if err != nil {
		g.Log().Warningf(ctx, "address:%s, message: %s, signature: %s, Failed to verify message: %v", req.Message.Address, message, req.Signature, err)
		return nil, fmt.Errorf("Failed to verify message")
	}

	token := uuid.New().String()
	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, err
	}

	md := dao.Users.Ctx(ctx)
	dbres, err := md.Where(dao.Users.Columns().Wallet, req.Message.Address).One()
	if err != nil {
		return nil, err
	}

	if dbres == nil {
		user := entity.Users{
			Wallet:     req.Message.Address,
			Token:      token,
			CreateTime: gtime.Now(),
			UpdateTime: gtime.Now(),
		}

		_, err := md.Insert(user)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
	} else {
		updates := map[string]interface{}{
			"Token":       token,
			"update_time": gtime.Now(),
		}
		if _, err := md.Where("wallet", req.Message.Address).Update(updates); err != nil {
			return nil, err
		}
		tx.Commit()
	}
	return &v1.UserLoginRes{
		Message: "login success",
		Token:   token,
	}, nil
}

func (c *cUser) Logout(ctx context.Context, req *v1.UserLogoutReq) (res *v1.UserLogoutRes, err error) {
	var user *entity.Users
	md := dao.Users.Ctx(ctx)
	err = md.Where("token", req.Token).Where("wallet", req.User).Scan(&user)
	if err != nil {
		g.Log().Errorf(ctx, "user:%s, token:%s, Failed to logout: %v", req.User, req.Token, err)
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	_, err = md.Where("token", req.Token).Where("wallet", req.User).Update("token=''")
	if err != nil {
		g.Log().Errorf(ctx, "user:%s, token:%s, update wallet data failed: %v", req.User, req.Token, err)
		return nil, err
	}

	return &v1.UserLogoutRes{
		Message: "logout success",
	}, nil
}

func (c *cUser) ActiveUserInfo(ctx context.Context, req *v1.ActiveUserInfoReq) (res *v1.ActiveUserInfoRes, err error) {
	md := dao.Users.Ctx(ctx)
	total, err := md.Count()
	if err != nil {
		g.Log().Errorf(ctx, "get total user count failed: %v", err)
		return nil, err
	}
	return &v1.ActiveUserInfoRes{
		Message:    "get active user info success",
		TotalUsers: total,
	}, nil
}
