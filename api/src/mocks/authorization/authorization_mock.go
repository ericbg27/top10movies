package authorization

import (
	"errors"
	"strconv"
	"strings"

	auth "github.com/ericbg27/top10movies-api/src/utils/authorization"
)

type AuthorizationMock struct {
	CanCreate  bool
	Authorized bool
	WrongID    bool
}

func (a AuthorizationMock) CreateToken(userId int64) (*auth.TokenDetails, error) {
	if !a.CanCreate {
		return nil, errors.New("failed to create token")
	}

	var sb strings.Builder

	sb.WriteString("token_")
	sb.WriteString(strconv.Itoa(int(userId)))

	s := sb.String()

	tokenInfo := &auth.TokenDetails{
		AccessToken: s,
	}

	return tokenInfo, nil
}

func (a AuthorizationMock) FetchAuth(bearToken string) (uint64, error) {
	if !a.Authorized {
		return 0, errors.New("not authorized")
	}

	if a.WrongID {
		return 0, nil
	}

	id, _ := strconv.Atoi(strings.Split(bearToken, "_")[1])

	return uint64(id), nil
}
