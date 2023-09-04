package post

import (
	"testing"

	"github.com/RianNegreiros/go-graphql-api/internal/user"
	"github.com/RianNegreiros/go-graphql-api/tests/faker"

	"github.com/RianNegreiros/go-graphql-api/internal/post"
	"github.com/stretchr/testify/require"
)

func TestCreatePostInput_Sanitize(t *testing.T) {
	input := post.CreatePostInput{
		Body: " test     ",
	}

	want := post.CreatePostInput{
		Body: "test",
	}

	input.Sanitize()

	require.Equal(t, want, input)
}

func TestCreatePostInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input post.CreatePostInput
		err   error
	}{
		{
			name: "valid",
			input: post.CreatePostInput{
				Body: "test",
			},
			err: nil,
		},
		{
			name: "post not long enough",
			input: post.CreatePostInput{
				Body: "t",
			},
			err: user.ErrValidation,
		},
		{
			name: "post too long",
			input: post.CreatePostInput{
				Body: faker.RandStr(300),
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
