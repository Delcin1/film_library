package search_by_id_test

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

	searchById "film_library/internal/http-server/handlers/movie/search_by_id"
	"film_library/internal/http-server/handlers/movie/search_by_id/mocks"
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
			name:      "GetMovie Error",
			movieId:   1,
			respError: "movie search failed",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		// tc := tc // go version < 1.22

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			movieSearcherByIdMock := mocks.NewMovieSearcherById(t)

			if tc.respError == "" || tc.mockError != nil {
				movieSearcherByIdMock.On("GetMovie", tc.movieId).
					Return(postgres.Movie{}, tc.mockError).
					Once()
			}

			handler := searchById.New(slogdiscard.NewDiscardLogger(), movieSearcherByIdMock)

			input := fmt.Sprintf(`{"movie_id": %d}`, tc.movieId)

			req, err := http.NewRequest(http.MethodGet, "/movie/search_by_id", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp searchById.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
