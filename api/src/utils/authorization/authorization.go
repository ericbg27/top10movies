package authorization

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthorizationManagerInterface interface {
	CreateToken(int64) (string, error)
}

type AuthorizationManager struct {
	accessSecret string
}

var (
	AuthManager AuthorizationManagerInterface = &AuthorizationManager{}
)

func init() {
	AuthManager.(*AuthorizationManager).accessSecret = os.Getenv("TOP10MOVIES_ACCESS_SECRET")
}

func (a AuthorizationManager) CreateToken(userId int64) (string, error) {
	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token, err := at.SignedString([]byte(a.accessSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}
