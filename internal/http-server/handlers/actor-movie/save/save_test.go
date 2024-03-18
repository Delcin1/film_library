package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"film_library/internal/http-server/handlers/actor-movie/save"
	"film_library/internal/http-server/handlers/actor-movie/save/mocks"
	"film_library/internal/lib/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		movieId   int
		actorsIds []int
		respError string
		mockError error
	}{
		{
			name:      "Success",
			movieId:   1,
			actorsIds: []int{1, 2},
		},
		{
			name:      "Invalid movie_id",
			movieId:   -1,
			actorsIds: []int{1, 2},
			respError: "field movie_id is not valid",
		},
		{
			name:      "Invalid actors_ids",
			movieId:   1,
			actorsIds: []int{1, -2},
			respError: "field actors_ids is not valid",
		},
		{
			name:      "SaveActorMovie Error",
			movieId:   1,
			actorsIds: []int{1, 2},
			respError: "failed to save actor-movie",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		// tc := tc // go version < 1.22

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actorMovieSaverMock := mocks.NewActorMovieSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				actorMovieSaverMock.On("SaveActorMovie", tc.movieId, tc.actorsIds).
					Return(tc.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), actorMovieSaverMock)

			actorsStr := fmt.Sprintf("%v", tc.actorsIds)
			actorsJSONStr := "[" + strings.Join(strings.Split(actorsStr[1:len(actorsStr)-1], " "), ", ") + "]"
			input := fmt.Sprintf(`{"movie_id": %d, "actors_ids": %s}`, tc.movieId, actorsJSONStr)

			req, err := http.NewRequest(http.MethodPost, "/actor-movie/save", bytes.NewReader([]byte(input)))
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
