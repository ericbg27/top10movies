package authorization

import (
	"errors"
	"strconv"
	"strings"
)

type AuthorizationMock struct {
	CanCreate bool
}

func (a AuthorizationMock) CreateToken(userId int64) (string, error) {
	if !a.CanCreate {
		return "", errors.New("failed to create token")
	}

	var sb strings.Builder

	sb.WriteString("token_")
	sb.WriteString(strconv.Itoa(int(userId)))

	s := sb.String()

	return s, nil
}
