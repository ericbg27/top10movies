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

	result, err := u.Validate()

	validatedUser := result.(User)

	assert.Nil(t, err)
	assert.EqualValues(t, "John", validatedUser.FirstName)
	assert.EqualValues(t, "Doe", validatedUser.LastName)
	assert.EqualValues(t, "johndoe@gmail.com", validatedUser.Email)
	assert.EqualValues(t, "12345", validatedUser.Password)
}

func TestValidateNameError(t *testing.T) {
	u := User{
		FirstName: "",
		LastName:  "LastName",
	}

	result, err := u.Validate()

	assert.NotNil(t, err)
	assert.Nil(t, result)

	u.FirstName = "FirstName"
	u.LastName = "     "

	result, err = u.Validate()

	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestValidateEmailError(t *testing.T) {
	u := User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "",
	}

	result, err := u.Validate()

	assert.NotNil(t, err)
	assert.Nil(t, result)

	u.Email = "johndoe#email.com"

	result, err = u.Validate()

	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestValidatePasswordError(t *testing.T) {
	u := User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@gmail.com",
		Password:  "",
	}

	result, err := u.Validate()

	assert.NotNil(t, err)
	assert.Nil(t, result)
}
