package authorization

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	redisdb "github.com/ericbg27/top10movies-api/src/datasources/redis"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

type AccessDetails struct {
	accessUuid string
	userId     uint64
}

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
	FetchAuth(bearToken string) (uint64, error)
}

type authorizationManager struct {
	accessSecret  string
	refreshSecret string
	redisClient   *redisdb.RedisClient
}

var (
	AuthManager AuthorizationManagerInterface = &authorizationManager{}
)

// TODO: Get environment variable names from config
func NewAuthorizationManager(redisClient *redisdb.RedisClient) *authorizationManager {
	authManager := &authorizationManager{
		accessSecret:  os.Getenv("TOP10MOVIES_ACCESS_SECRET"),
		refreshSecret: os.Getenv("TOP10MOVIES_REFRESH_SECRET"),
		redisClient:   redisClient,
	}

	return authManager
}

func (a authorizationManager) CreateToken(userId int64) (*TokenDetails, error) {
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

func (a authorizationManager) saveTokenMetadata(userId int64, tokenInfo *TokenDetails) error {
	at := time.Unix(tokenInfo.AtExpires, 0)
	rt := time.Unix(tokenInfo.RtExpires, 0)
	now := time.Now()

	errAccess := a.redisClient.Client.Set(tokenInfo.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := a.redisClient.Client.Set(tokenInfo.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

func (a authorizationManager) extractToken(bearToken string) string {
	bearTokenArgs := strings.Split(bearToken, " ")
	if len(bearTokenArgs) == 2 {
		return bearTokenArgs[1]
	}

	return ""
}

func (a authorizationManager) verifyToken(bearToken string) (*jwt.Token, error) {
	tokenString := a.extractToken(bearToken)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.accessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, err
	}

	return token, nil
}

func (a authorizationManager) extractTokenMetadata(bearToken string) (*AccessDetails, error) {
	token, err := a.verifyToken(bearToken)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}

		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}

		return &AccessDetails{
			accessUuid: accessUuid,
			userId:     userId,
		}, nil
	}

	return nil, err
}

func (a authorizationManager) FetchAuth(bearToken string) (uint64, error) {
	accessDetails, err := a.extractTokenMetadata(bearToken)
	if err != nil {
		return 0, err
	}

	userId, err := a.redisClient.Client.Get(accessDetails.accessUuid).Result()
	if err != nil {
		return 0, err
	}

	userID, _ := strconv.ParseUint(userId, 10, 64)

	return userID, nil
}
