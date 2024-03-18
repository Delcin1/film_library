package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"film_library/internal/http-server/handlers/actor/save"
	"film_library/internal/http-server/handlers/actor/save/mocks"
	"film_library/internal/lib/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		actorName string
		gender    string
		birthdate string
		respError string
		mockError error
	}{
		{
			name:      "Success",
			actorName: "Nikita",
			gender:    "male",
			birthdate: "2000-01-01",
		},
		{
			name:      "Empty name",
			actorName: "",
			gender:    "female",
			birthdate: "1900-01-01",
			respError: "field name is not valid",
		},
		{
			name:      "Invalid gender",
			actorName: "Nikita",
			gender:    "abc",
			birthdate: "1900-01-01",
			respError: "field gender is not valid",
		},
		{
			name:      "Invalid birthdate",
			actorName: "Nikita",
			gender:    "male",
			birthdate: "1900.01.01",
			respError: "field birthdate is not valid",
		},
		{
			name:      "SaveActor Error",
			actorName: "Nikita",
			gender:    "male",
			birthdate: "2000-01-01",
			respError: "failed to save actor",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		// tc := tc // go version < 1.22

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actorSaverMock := mocks.NewActorSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				actorSaverMock.On("SaveActor", tc.actorName, tc.gender, tc.birthdate).
					Return(1, tc.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), actorSaverMock)

			input := fmt.Sprintf(`{"name": "%s", "gender": "%s", "birthdate": "%s"}`,
				tc.actorName, tc.gender, tc.birthdate)

			req, err := http.NewRequest(http.MethodPost, "/actor/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
