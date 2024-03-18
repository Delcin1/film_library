package all_test

import (
	"encoding/json"
	"errors"
	"film_library/internal/storage/postgres"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	searchAll "film_library/internal/http-server/handlers/actor/all"
	"film_library/internal/http-server/handlers/actor/all/mocks"
	"film_library/internal/lib/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		sortBy    string
		respError string
		mockError error
	}{
		{
			name: "Success",
		},
		{
			name:      "GetActors Error",
			respError: "actors search failed",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		// tc := tc // go version < 1.22

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actorsAllGetterMock := mocks.NewActorsAllGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				actorsAllGetterMock.On("GetActors").
					Return([]postgres.Actor{}, tc.mockError).
					Once()
			}

			handler := searchAll.New(slogdiscard.NewDiscardLogger(), actorsAllGetterMock)

			req, err := http.NewRequest(http.MethodGet, "/actor/all", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp searchAll.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
