package search_by_part_test

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

	searchByPart "film_library/internal/http-server/handlers/movie/search_by_part"
	"film_library/internal/http-server/handlers/movie/search_by_part/mocks"
	"film_library/internal/lib/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		part      string
		respError string
		mockError error
	}{
		{
			name: "Success",
			part: "test",
		},
		{
			name:      "GetMoviesBySearchRequest Error",
			part:      "test",
			respError: "movies search failed",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		// tc := tc // go version < 1.22

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			movieSearcherByPartMock := mocks.NewMovieSearcherByPart(t)

			if tc.respError == "" || tc.mockError != nil {
				movieSearcherByPartMock.On("GetMoviesBySearchRequest", tc.part).
					Return([]postgres.Movie{}, tc.mockError).
					Once()
			}

			handler := searchByPart.New(slogdiscard.NewDiscardLogger(), movieSearcherByPartMock)

			input := fmt.Sprintf(`{"part": "%s"}`, tc.part)

			req, err := http.NewRequest(http.MethodGet, "/movie/search_by_part", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp searchByPart.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
