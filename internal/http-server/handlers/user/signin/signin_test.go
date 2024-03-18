package signin_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"film_library/internal/http-server/handlers/user/signin"
	"film_library/internal/http-server/handlers/user/signin/mocks"
	"film_library/internal/lib/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		username  string
		password  string
		respError string
		mockError error
	}{
		{
			name:     "Success",
			username: "admin",
			password: "admin",
		},
		{
			name:      "Empty username",
			username:  "",
			password:  "admin",
			respError: "field username is required",
		},
		{
			name:      "Empty password",
			username:  "admin",
			password:  "",
			respError: "field password is required",
		},
		{
			name:      "SaveUser Error",
			username:  "admin",
			password:  "admin",
			respError: "failed to authenticate user",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		// tc := tc // go version < 1.22

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userAuthenticatorMock := mocks.NewUserAuthenticator(t)

			if tc.respError == "" || tc.mockError != nil {
				userAuthenticatorMock.On("GetUser", tc.username, tc.password).
					Return(1, tc.mockError).
					Once()
			}

			handler := signin.New(slogdiscard.NewDiscardLogger(), userAuthenticatorMock, "jwtKey")

			input := fmt.Sprintf(`{"username": "%s", "password": "%s"}`,
				tc.username, tc.password)

			req, err := http.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp signin.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
