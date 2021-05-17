package users

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestValidateSuccess(t *testing.T) {
	u := User{
		FirstName: " John  ",
		LastName:  "Doe",
		Email:     "    johndoe@gmail.com   ",
		Password:  "   12345   ",
	}

	err := u.Validate()

	assert.Nil(t, err)
	assert.EqualValues(t, "John", u.FirstName)
	assert.EqualValues(t, "Doe", u.LastName)
	assert.EqualValues(t, "johndoe@gmail.com", u.Email)
	assert.EqualValues(t, "12345", u.Password)
}

func TestValidateNameError(t *testing.T) {
	u := User{
		FirstName: "",
		LastName:  "LastName",
	}

	err := u.Validate()

	assert.NotNil(t, err)

	u.FirstName = "FirstName"
	u.LastName = "     "

	err = u.Validate()

	assert.NotNil(t, err)
}

func TestValidateEmailError(t *testing.T) {
	u := User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "",
	}

	err := u.Validate()

	assert.NotNil(t, err)

	u.Email = "johndoe#email.com"

	u.Validate()

	assert.NotNil(t, err)
}

func TestValidatePasswordError(t *testing.T) {
	u := User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@gmail.com",
		Password:  "",
	}

	err := u.Validate()

	assert.NotNil(t, err)
}
