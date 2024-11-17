package controller

import (
	v1 "anime_petition/api/v1"
	"anime_petition/internal/dao"
	"anime_petition/internal/model/entity"
	"anime_petition/utility"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	siwe "github.com/spruceid/siwe-go"
)

var User = &cUser{}

type cUser struct{}

func (c *cUser) Login(ctx context.Context, req *v1.UserLoginReq) (res *v1.UserLoginRes, err error) {
	lock := utility.Lock{}
	if !lock.Lock(ctx, req.Message.Address) {
		return nil, errors.New("lock failed")
	}
	defer lock.Unlock(ctx, req.Message.Address)

	message, err := siwe.InitMessage(req.Message.Domain, req.Message.Address, req.Message.Uri, req.Message.Nonce, map[string]interface{}{"chainId": req.Message.ChainId, "statement": req.Message.Statement, "issuedAt": req.Message.IssuedAt})
	if err != nil {
		g.Log().Warning(ctx, "Failed to init message:", err)
		return nil, err
	}

	pubKey, err := message.VerifyEIP191(req.Signature)
	if err != nil {
		g.Log().Warning(ctx, "Failed to verify message:", err)
		return nil, err
	}

	// Convert the public key to an address
	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()
	// Check if the recovered address matches the expected address
	if strings.ToLower(recoveredAddress) == strings.ToLower(req.Message.Address) {
		token := uuid.New().String()

		tx, err := g.DB().Begin(ctx)
		if err != nil {
			return nil, err
		}

		md := dao.Users.Ctx(ctx)
		refId := req.Ref
		res, err := md.Where(dao.Users.Columns().Wallet, strings.ToLower(req.Message.Address)).One()
		if err != nil {
			return nil, err
		}

		if res == nil {
			var refUser *entity.Users
			_ = md.Where(dao.Users.Columns().Id, refId).Scan(&refUser)

			if refUser == nil {
				var team *entity.TeamMembers
				// Add team member
				teamMd := dao.TeamMembers.Ctx(ctx)
				err = teamMd.Where(dao.TeamMembers.Columns().UserId, 0).Scan(&team)
				if err != nil {
					tx.Rollback()
					return nil, errors.New("can't find ref team")
				}

				user := entity.Users{
					Wallet:     req.Message.Address,
					RefId:      0,
					Token:      token,
					CreateTime: gtime.Now(),
					UpdateTime: gtime.Now(),
				}

				res, err := md.Insert(user)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				id, err := res.LastInsertId()
				if err != nil {
					tx.Rollback()
					return nil, err
				}

				teamMember := entity.TeamMembers{
					UserId:    int(id),
					TeamId:    team.TeamId,
					CreaterId: team.CreaterId,
				}

				_, err = teamMd.Insert(teamMember)
				if err != nil {
					tx.Rollback()
					return nil, err
				}

				tx.Commit()
				return &v1.UserLoginRes{
					Token: token,
				}, nil

			} else {
				var team []entity.TeamMembers
				// Add team member
				teamMd := dao.TeamMembers.Ctx(ctx)
				err = teamMd.Where(dao.TeamMembers.Columns().UserId, refUser.Id).Scan(&team)
				if len(team) == 0 {
					tx.Rollback()
					return nil, errors.New("can't find ref team")
				}

				user := entity.Users{
					Wallet:     req.Message.Address,
					RefId:      refUser.Id,
					Token:      token,
					CreateTime: gtime.Now(),
					UpdateTime: gtime.Now(),
				}

				res, err := md.Insert(user)
				if err != nil {
					tx.Rollback()
					return nil, err
				}

				id, err := res.LastInsertId()
				if err != nil {
					tx.Rollback()
					return nil, err
				}

				var team1 entity.TeamMembers
				if len(team) == 1 {
					team1 = team[0]
				} else {
					for _, t := range team {
						if t.CreaterId == refUser.Id {
							team1 = t
							break
						}
					}
				}

				teamMember := entity.TeamMembers{
					UserId:    int(id),
					TeamId:    team1.TeamId,
					CreaterId: team1.CreaterId,
				}

				_, err = teamMd.Insert(teamMember)
				if err != nil {
					tx.Rollback()
					return nil, err
				}

				tx.Commit()
				return &v1.UserLoginRes{
					Token: token,
				}, nil

			}

		} else {
			if _, err := md.Where("wallet", strings.ToLower(req.Message.Address)).Update(fmt.Sprintf("token='%s'", token)); err != nil {
				return nil, err
			}
			tx.Commit()
			return &v1.UserLoginRes{
				Token: token,
			}, nil
		}
	} else {
		return nil, errors.New("Signature does not match known address!")
	}
}

func (c *cUser) Logout(ctx context.Context, req *v1.UserLogoutReq) (res *v1.UserLogoutRes, err error) {
	var user entity.Users
	md := dao.Users.Ctx(ctx)
	err = md.Where("token", req.Token).Where("wallet", strings.ToLower(req.User)).Scan(&user)
	if err != nil {
		return res, nil
	}

	_, err = md.Where("token", req.Token).Where("wallet", strings.ToLower(req.User)).Update("token=''")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *cUser) MakeChoose(ctx context.Context, req *v1.MakeChooseReq) (res *v1.MakeChooseRes, err error) {
	if req.Choose != 1 && req.Choose != 2 {
		return nil, errors.New("invalid choice: must be 1 or 2")
	}

	var user entity.Users
	md := dao.Users.Ctx(ctx)
	err = md.Where("token", req.Token).Where("wallet", strings.ToLower(req.User)).Scan(&user)
	if err != nil {
		return nil, errors.New("user not Login")
	}
	choose, err := md.Where("choose", 0).Where("token", req.Token).Where("wallet", strings.ToLower(req.User)).Count("id")
	if err != nil {
		return nil, errors.New("query choose failed")
	}

	if choose == 0 {
		return nil, errors.New("already choose")
	}

	_, err = md.Where("token", req.Token).Where("wallet", strings.ToLower(req.User)).Update("choose=" + strconv.Itoa(req.Choose))
	if err != nil {
		return nil, err
	}

	return res, nil
}
