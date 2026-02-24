package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestPasswordStrength(t *testing.T) {
	t.Parallel()

	validate := validator.New()
	err := validate.RegisterValidation("password_strength", PasswordStrength)
	assert.NoError(t, err)

	type request struct {
		Password string `validate:"password_strength"`
	}

	testCases := []struct {
		name string

		inputPassword string

		expectedValid bool
	}{
		{
			name: "Valid password with all requirements",

			inputPassword: "Secure@1",

			expectedValid: true,
		},
		{
			name: "Missing uppercase letter",

			inputPassword: "secure@1",

			expectedValid: false,
		},
		{
			name: "Missing lowercase letter",

			inputPassword: "SECURE@1",

			expectedValid: false,
		},
		{
			name: "Missing number",

			inputPassword: "Secure@!",

			expectedValid: false,
		},
		{
			name: "Missing special character",

			inputPassword: "Secure123",

			expectedValid: false,
		},
		{
			name: "Empty password",

			inputPassword: "",

			expectedValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := validate.Struct(request{Password: tc.inputPassword})
			if tc.expectedValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
