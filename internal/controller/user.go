package controller

import (
	v1 "anime_petition/api/v1"
	"anime_petition/internal/dao"
	"anime_petition/internal/model/do"
	"anime_petition/internal/model/entity"
	"anime_petition/utility"
	"context"
	"fmt"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
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
		return nil, fmt.Errorf("failed to verify message")
	}
	redis := g.Redis()
	value, err := redis.Get(ctx, req.Message.Address)
	if err != nil || value.IsEmpty() {
		g.Log().Infof(ctx, "address:%s not in redis, user not login, err: %v", req.Message.Address, err)
	} else {
		if err := g.Validator().Rules("regex:^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$").Data(value.String()).Messages("invalid address and token").Run(ctx); err != nil {
			g.Log().Infof(ctx, "address:%s in redis, the value:%s is error: %v", req.Message.Address, value.String(), err)
		} else {
			g.Log().Infof(ctx, "address:%s in redis, the value:%s is valid", req.Message.Address, value.String())
			return &v1.UserLoginRes{
				Message: "login success",
				Token:   value.String(),
			}, nil
		}
	}

	token := uuid.New().String()
	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, err
	}

	md := dao.Users.Ctx(ctx)
	dbres, err := md.Where(do.Users{Wallet: req.Message.Address}).One()
	if err != nil {
		return nil, err
	}

	err = redis.SetEX(ctx, req.Message.Address, token, 600)
	if err != nil {
		tx.Rollback()
		g.Log().Errorf(ctx, "user:%s, token:%s, Failed to set redis: %v", req.Message.Address, token, err)
		return nil, err
	}
	if dbres == nil {
		user := do.Users{
			Wallet:     req.Message.Address,
			Token:      token,
			CreateTime: gtime.Now(),
			UpdateTime: gtime.Now(),
		}

		_, err := md.Data(user).Insert()
		if err != nil {
			tx.Rollback()
			g.Log().Errorf(ctx, "user:%s, token:%s, Failed to insert user: %v", req.Message.Address, token, err)
			return nil, err
		}
		tx.Commit()
	} else {
		updates := do.Users{
			Token:      token,
			UpdateTime: gtime.Now(),
		}
		if _, err := md.Where(do.Users{Wallet: req.Message.Address}).Data(updates).Update(); err != nil {
			return nil, err
		}
		err = redis.SetEX(ctx, req.Message.Address, token, 600)
		if err != nil {
			tx.Rollback()
			g.Log().Errorf(ctx, "user:%s, token:%s, Failed to set redis: %v", req.Message.Address, token, err)
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
	err = md.Where(do.Users{
		Token:  req.Token,
		Wallet: req.User,
	}).Scan(&user)
	if err != nil {
		g.Log().Errorf(ctx, "user:%s, token:%s, Failed to logout: %v", req.User, req.Token, err)
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	updates := do.Users{
		Token: "",
	}
	_, err = md.Where(do.Users{
		Token:  req.Token,
		Wallet: req.User,
	}).Data(updates).Update()
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
