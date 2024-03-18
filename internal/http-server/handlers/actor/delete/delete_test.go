package delete_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"film_library/internal/http-server/handlers/actor/delete"
	"film_library/internal/http-server/handlers/actor/delete/mocks"
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
			name:      "DeleteActor Error",
			actorId:   1,
			respError: "failed to delete actor",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		// tc := tc // go version < 1.22

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actorDeleterMock := mocks.NewActorDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				actorDeleterMock.On("DeleteActor", tc.actorId).
					Return(tc.mockError).
					Once()
			}

			handler := delete.New(slogdiscard.NewDiscardLogger(), actorDeleterMock)

			input := fmt.Sprintf(`{"actor_id": %d}`, tc.actorId)

			req, err := http.NewRequest(http.MethodPost, "/actor/delete", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp delete.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
