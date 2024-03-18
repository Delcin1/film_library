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

	"film_library/internal/http-server/handlers/movie/delete"
	"film_library/internal/http-server/handlers/movie/delete/mocks"
	"film_library/internal/lib/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		movieId   int
		respError string
		mockError error
	}{
		{
			name:    "Success",
			movieId: 1,
		},
		{
			name:      "Invalid movie_id",
			movieId:   -1,
			respError: "field movie_id is not valid",
		},
		{
			name:      "DeleteMovie Error",
			movieId:   1,
			respError: "failed to delete movie",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		// tc := tc // go version < 1.22

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			movieDeleterMock := mocks.NewMovieDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				movieDeleterMock.On("DeleteMovie", tc.movieId).
					Return(tc.mockError).
					Once()
			}

			handler := delete.New(slogdiscard.NewDiscardLogger(), movieDeleterMock)

			input := fmt.Sprintf(`{"movie_id": %d}`, tc.movieId)

			req, err := http.NewRequest(http.MethodPost, "/movie/delete", bytes.NewReader([]byte(input)))
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
