package signup_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"film_library/internal/http-server/handlers/user/signup"
	"film_library/internal/http-server/handlers/user/signup/mocks"
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
			respError: "failed to save user",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		// tc := tc // go version < 1.22

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userSaverMock := mocks.NewUserSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				userSaverMock.On("SaveUser", tc.username, tc.password).
					Return(tc.mockError).
					Once()
			}

			handler := signup.New(slogdiscard.NewDiscardLogger(), userSaverMock)

			input := fmt.Sprintf(`{"username": "%s", "password": "%s"}`,
				tc.username, tc.password)

			req, err := http.NewRequest(http.MethodPost, "/user/signup", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp signup.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
