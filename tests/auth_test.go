package tests

import (
	"testing"

	"github.com/RianNegreiros/go-graphql-api/models"
	"github.com/stretchr/testify/require"
)

func TestRegisterInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input models.RegisterInput
		err   error
	}{
		{
			name: "valid input",
			input: models.RegisterInput{
				Username:        "john",
				Email:           "johndoe@mail.com",
				Password:        "123456",
				ConfirmPassword: "123456",
			},
			err: nil,
		},
		{
			name: "invalid email",
			input: models.RegisterInput{
				Username:        "john",
				Email:           "john",
				Password:        "123456",
				ConfirmPassword: "123456",
			},
			err: models.ErrValidation,
		}, {
			name: "invalid username",
			input: models.RegisterInput{
				Username:        "j",
				Email:           "johndoe@mail.com",
				Password:        "123456",
				ConfirmPassword: "123456",
			},
			err: models.ErrValidation,
		}, {
			name: "invalid password",
			input: models.RegisterInput{
				Username:        "john",
				Email:           "johndoe@mail.com",
				Password:        "12345",
				ConfirmPassword: "12345",
			},
			err: models.ErrValidation,
		}, {
			name: "password and confirm password don't match",
			input: models.RegisterInput{
				Username:        "john",
				Email:           "johndoe@mail.com",
				Password:        "123456",
				ConfirmPassword: "1234567",
			},
			err: models.ErrValidation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRegisterInput_Sanitize(t *testing.T) {
	input := models.RegisterInput{
		Username:        "  john  ",
		Email:           "  JOHNDOE@mail.com ",
		Password:        " 123456 ",
		ConfirmPassword: " 123456 ",
	}

	want := models.RegisterInput{
		Username:        "john",
		Email:           "johndoe@mail.com",
		Password:        "123456",
		ConfirmPassword: "123456",
	}

	input = input.Sanitize()

	require.Equal(t, want, input)
}

func TestLoginInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input models.LoginInput
		err   error
	}{
		{
			name: "valid input",
			input: models.LoginInput{
				Email:    "johndoe@mail.com",
				Password: "123456",
			},
			err: nil,
		},
		{
			name: "invalid email",
			input: models.LoginInput{
				Email:    "invalid_email",
				Password: "123456",
			},
			err: models.ErrValidation,
		},
		{
			name: "password required",
			input: models.LoginInput{
				Email:    "johndoe@mail.com",
				Password: "",
			},
			err: models.ErrValidation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLoginInput_Sanitize(t *testing.T) {
	input := models.LoginInput{
		Email:    "   JOHNDOE@mail.com   ",
		Password: " 123456 ",
	}

	want := models.LoginInput{
		Email:    "johndoe@mail.com",
		Password: "123456",
	}

	input = input.Sanitize()

	require.Equal(t, want, input)
}
