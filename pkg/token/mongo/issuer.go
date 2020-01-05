package mongo

import (
	"time"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/types/ftime"

	"github.com/infraboard/keyauth/pkg/application"
	"github.com/infraboard/keyauth/pkg/token"
	"github.com/infraboard/keyauth/pkg/user"
)

func (s *service) newTokenIssuer(req *token.IssueTokenRequest) (*TokenIssuer, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	issuer := &TokenIssuer{
		IssueTokenRequest: req,
		clientChecker:     newClientChecker(s.app),
		user:              s.user,
	}
	return issuer, nil
}

// TokenIssuer 基于该数据进行扩展
type TokenIssuer struct {
	*token.IssueTokenRequest
	*clientChecker
	user user.Service
}

func (i *TokenIssuer) checkUser() (*user.User, error) {
	req := user.NewDescriptAccountRequest()
	req.Account = i.Username
	u, err := i.user.DescribeAccount(req)
	if err != nil {
		return nil, err
	}
	if err := u.HashedPassword.CheckPassword(i.Password); err != nil {
		return nil, err
	}

	return u, nil
}

// IssueToken 颁发token
func (i *TokenIssuer) IssueToken() (tk *token.Token, err error) {
	app, err := i.CheckClient(i.ClientID, i.ClientSecret)
	if err != nil {
		err = exception.NewUnauthorized(err.Error())
		return
	}

	switch i.GrantType {
	case token.PASSWORD:
		u, checkErr := i.checkUser()
		if checkErr != nil {
			err = exception.NewUnauthorized("user or password not connrect")
			return
		}

		tk = i.issuePasswordToken(app, u)
		return
	case token.CLIENT:
	case token.AUTHCODE:
	case token.REFRESH:
	default:
		err = exception.NewInternalServerError("unknown grant type %s", i.GrantType)
		return
	}

	return
}

func (i *TokenIssuer) issuePasswordToken(app *application.Application, u *user.User) *token.Token {
	tk := i.newBearToken(app)
	tk.UserID = u.ID
	return tk
}

func (i *TokenIssuer) newBearToken(app *application.Application) *token.Token {
	now := time.Now()
	expire := now.Add(time.Duration(app.TokenExpireSecond) * time.Second)
	return &token.Token{
		Type:          token.Bearer,
		AccessToken:   token.MakeBearer(24),
		RefreshToken:  token.MakeBearer(32),
		CreatedAt:     ftime.T(now),
		ClientID:      i.ClientID,
		GrantType:     i.GrantType,
		ExpiresAt:     ftime.T(expire),
		ExpiresIn:     app.TokenExpireSecond,
		ApplicationID: app.ID,
	}
}

func newClientChecker(app application.Service) *clientChecker {
	return &clientChecker{app}
}

// clientChecker 检测client正确性
type clientChecker struct {
	application.Service
}

func (ck *clientChecker) CheckClient(clientID, clientSecret string) (*application.Application, error) {
	req := application.NewDescriptApplicationRequest()
	req.ClientID = clientID
	app, err := ck.DescriptionApplication(req)
	if err != nil {
		return nil, err
	}

	if err := app.CheckClientSecret(clientSecret); err != nil {
		return nil, err
	}

	return app, nil
}
