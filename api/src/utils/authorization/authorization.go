package authorization

import (
	"os"
	"strconv"
	"time"

	redisdb "github.com/ericbg27/top10movies-api/src/datasources/redis"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AuthorizationManagerInterface interface {
	CreateToken(int64) (*TokenDetails, error)
}

type AuthorizationManager struct {
	accessSecret  string
	refreshSecret string
}

var (
	AuthManager AuthorizationManagerInterface = &AuthorizationManager{}
)

// TODO: Get environment variable names from config
func init() {
	AuthManager.(*AuthorizationManager).accessSecret = os.Getenv("TOP10MOVIES_ACCESS_SECRET")
	AuthManager.(*AuthorizationManager).refreshSecret = os.Getenv("TOP10MOVIES_REFRESH_SECRET")
}

func (a AuthorizationManager) CreateToken(userId int64) (*TokenDetails, error) {
	tokenInfo := &TokenDetails{}

	tokenInfo.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	newUuid, _ := uuid.NewV4()
	tokenInfo.AccessUuid = newUuid.String()

	tokenInfo.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	newUuid, _ = uuid.NewV4()
	tokenInfo.RefreshUuid = newUuid.String()

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = tokenInfo.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = tokenInfo.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	tokenInfo.AccessToken, err = at.SignedString([]byte(a.accessSecret))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = tokenInfo.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = tokenInfo.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	tokenInfo.RefreshToken, err = rt.SignedString([]byte(a.refreshSecret))
	if err != nil {
		return nil, err
	}

	err = a.saveTokenMetadata(userId, tokenInfo)
	if err != nil {
		return nil, err
	}

	return tokenInfo, nil
}

func (a AuthorizationManager) saveTokenMetadata(userId int64, tokenInfo *TokenDetails) error {
	at := time.Unix(tokenInfo.AtExpires, 0)
	rt := time.Unix(tokenInfo.RtExpires, 0)
	now := time.Now()

	errAccess := redisdb.Client.Set(tokenInfo.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := redisdb.Client.Set(tokenInfo.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}

	return nil
}
