package tests

import (
	"testing"

	"github.com/RianNegreiros/go-graphql-api/internal/user"
	"github.com/stretchr/testify/require"
)

func TestRegisterInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input user.RegisterInput
		err   error
	}{
		{
			name: "valid input",
			input: user.RegisterInput{
				Username:        "john",
				Email:           "johndoe@mail.com",
				Password:        "123456",
				ConfirmPassword: "123456",
			},
			err: nil,
		},
		{
			name: "invalid email",
			input: user.RegisterInput{
				Username:        "john",
				Email:           "john",
				Password:        "123456",
				ConfirmPassword: "123456",
			},
			err: user.ErrValidation,
		}, {
			name: "invalid username",
			input: user.RegisterInput{
				Username:        "j",
				Email:           "johndoe@mail.com",
				Password:        "123456",
				ConfirmPassword: "123456",
			},
			err: user.ErrValidation,
		}, {
			name: "invalid password",
			input: user.RegisterInput{
				Username:        "john",
				Email:           "johndoe@mail.com",
				Password:        "12345",
				ConfirmPassword: "12345",
			},
			err: user.ErrValidation,
		}, {
			name: "password and confirm password don't match",
			input: user.RegisterInput{
				Username:        "john",
				Email:           "johndoe@mail.com",
				Password:        "123456",
				ConfirmPassword: "1234567",
			},
			err: user.ErrValidation,
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
	input := user.RegisterInput{
		Username:        "  john  ",
		Email:           "  JOHNDOE@mail.com ",
		Password:        " 123456 ",
		ConfirmPassword: " 123456 ",
	}

	want := user.RegisterInput{
		Username:        "john",
		Email:           "johndoe@mail.com",
		Password:        "123456",
		ConfirmPassword: "123456",
	}

	input.Sanitize()

	require.Equal(t, want, input)
}

func TestLoginInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input user.LoginInput
		err   error
	}{
		{
			name: "valid input",
			input: user.LoginInput{
				Email:    "johndoe@mail.com",
				Password: "123456",
			},
			err: nil,
		},
		{
			name: "invalid email",
			input: user.LoginInput{
				Email:    "invalid_email",
				Password: "123456",
			},
			err: user.ErrValidation,
		},
		{
			name: "password required",
			input: user.LoginInput{
				Email:    "johndoe@mail.com",
				Password: "",
			},
			err: user.ErrValidation,
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
	input := user.LoginInput{
		Email:    "   JOHNDOE@mail.com   ",
		Password: " 123456 ",
	}

	want := user.LoginInput{
		Email:    "johndoe@mail.com",
		Password: "123456",
	}

	input.Sanitize()

	require.Equal(t, want, input)
}
