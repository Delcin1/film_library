package search_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"film_library/internal/storage/postgres"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"film_library/internal/http-server/handlers/actor/search"
	"film_library/internal/http-server/handlers/actor/search/mocks"
	"film_library/internal/lib/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		actorId   int
		respError string
		mockError error
	}{
		{
			name:    "Success",
			actorId: 1,
		},
		{
			name:      "Invalid actor_id",
			actorId:   -1,
			respError: "field actor_id is not valid",
		},
		{
			name:      "GetActor Error",
			actorId:   1,
			respError: "actor search failed",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		// tc := tc // go version < 1.22

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actorSearcherMock := mocks.NewActorSearcher(t)

			if tc.respError == "" || tc.mockError != nil {
				actorSearcherMock.On("GetActor", tc.actorId).
					Return(postgres.Actor{}, tc.mockError).
					Once()
			}

			handler := search.New(slogdiscard.NewDiscardLogger(), actorSearcherMock)

			input := fmt.Sprintf(`{"actor_id": %d}`, tc.actorId)

			req, err := http.NewRequest(http.MethodGet, "/actor/search", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp search.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
