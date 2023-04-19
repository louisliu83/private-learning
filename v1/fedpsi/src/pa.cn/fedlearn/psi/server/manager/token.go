package manager

import (
	"context"
	"time"

	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/auth"
	"pa.cn/fedlearn/psi/log"
	"pa.cn/fedlearn/psi/model"
)

type TokenManager struct {
}

func (mgr *TokenManager) GenToken(ctx context.Context, userID string, d time.Duration, utype string) (token *types.TokenGenResponse, err error) {
	log.Debugln(ctx, "TokenManager.GenToken is called")
	user, err := model.GetUserByUserName(userID)
	if err != nil {
		return nil, err
	}
	customClaims := auth.CustomClaims{
		Party:  user.Party,
		UserID: user.UserName,
		Admin:  user.IsRoot,
		Type:   utype,
	}
	tokenStr, err := auth.GenerateToken(customClaims, d)
	if err != nil {
		return nil, err
	}
	token = &types.TokenGenResponse{
		Token: tokenStr,
	}
	return token, nil
}
