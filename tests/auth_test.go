package tests

import (
	"github.com/RianNegreiros/go-graphql-api/internal/transport/http"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegisterInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input http.RegisterInput
		err   error
	}{
		{
			name: "valid input",
			input: http.RegisterInput{
				Username:        "john",
				Email:           "johndoe@mail.com",
				Password:        "123456",
				ConfirmPassword: "123456",
			},
			err: nil,
		},
		{
			name: "invalid email",
			input: http.RegisterInput{
				Username:        "john",
				Email:           "john",
				Password:        "123456",
				ConfirmPassword: "123456",
			},
			err: http.ErrValidation,
		}, {
			name: "invalid username",
			input: http.RegisterInput{
				Username:        "j",
				Email:           "johndoe@mail.com",
				Password:        "123456",
				ConfirmPassword: "123456",
			},
			err: http.ErrValidation,
		}, {
			name: "invalid password",
			input: http.RegisterInput{
				Username:        "john",
				Email:           "johndoe@mail.com",
				Password:        "12345",
				ConfirmPassword: "12345",
			},
			err: http.ErrValidation,
		}, {
			name: "password and confirm password don't match",
			input: http.RegisterInput{
				Username:        "john",
				Email:           "johndoe@mail.com",
				Password:        "123456",
				ConfirmPassword: "1234567",
			},
			err: http.ErrValidation,
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
	input := http.RegisterInput{
		Username:        "  john  ",
		Email:           "  JOHNDOE@mail.com ",
		Password:        " 123456 ",
		ConfirmPassword: " 123456 ",
	}

	want := http.RegisterInput{
		Username:        "john",
		Email:           "johndoe@mail.com",
		Password:        "123456",
		ConfirmPassword: "123456",
	}

	input.Sanitize()

	require.Equal(t, want, input)
}
