package all_test

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

	searchAll "film_library/internal/http-server/handlers/movie/all"
	"film_library/internal/http-server/handlers/movie/all/mocks"
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
			name:   "Success",
			sortBy: "title_asc",
		},
		{
			name:      "Invalid sort_by",
			sortBy:    "title",
			respError: "field sort_by is not valid",
		},
		{
			name:      "GetMovies Error",
			sortBy:    "title_asc",
			respError: "movies search failed",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		// tc := tc // go version < 1.22

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			moviesAllGetterMock := mocks.NewMoviesAllGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				moviesAllGetterMock.On("GetMovies", tc.sortBy).
					Return([]postgres.Movie{}, tc.mockError).
					Once()
			}

			handler := searchAll.New(slogdiscard.NewDiscardLogger(), moviesAllGetterMock)

			input := fmt.Sprintf(`{"sort_by": "%s"}`, tc.sortBy)

			req, err := http.NewRequest(http.MethodGet, "/movie/all", bytes.NewReader([]byte(input)))
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
