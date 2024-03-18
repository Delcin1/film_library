package update_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"film_library/internal/http-server/handlers/movie/update"
	"film_library/internal/http-server/handlers/movie/update/mocks"
	"film_library/internal/lib/logger/handlers/slogdiscard"
)

const bigTitle = "Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea co"
const bigDesription = "Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat. Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi.Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat. Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui"

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name        string
		movieId     int
		title       string
		description string
		releaseDate string
		rating      int
		respError   string
		mockError   error
	}{
		{
			name:    "Success update title",
			movieId: 1,
			title:   "Best movie",
		},
		{
			name:      "Invalid title",
			movieId:   1,
			title:     bigTitle,
			respError: "field title is not valid",
		},
		{
			name:        "Success update description",
			movieId:     1,
			description: "some description",
		},
		{
			name:        "Invalid description",
			movieId:     1,
			description: bigDesription,
			respError:   "field description is not valid",
		},
		{
			name:        "Success update release_date",
			movieId:     1,
			releaseDate: "2000-01-01",
		},
		{
			name:        "Invalid release_date",
			movieId:     1,
			releaseDate: "1900.01.01",
			respError:   "field release_date is not valid",
		},
		{
			name:    "Success update rating",
			movieId: 1,
			rating:  10,
		},
		{
			name:      "Invalid rating",
			movieId:   1,
			rating:    11,
			respError: "field rating is not valid",
		},
		{
			name:        "Success update all",
			movieId:     1,
			title:       "Best movie",
			description: "some description",
			releaseDate: "2000-01-01",
			rating:      10,
		},
		{
			name:      "Invalid movieId",
			movieId:   0,
			respError: "field movie_id is not valid",
		},
		{
			name:        "UpdateActor Error",
			movieId:     1,
			title:       "Best movie",
			description: "some description",
			releaseDate: "2000-01-01",
			rating:      10,
			respError:   "failed to update movie title",
			mockError:   errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		// tc := tc // go version < 1.22

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			movieUpdaterMock := mocks.NewMovieUpdater(t)

			if tc.respError == "" || tc.mockError != nil {
				movieUpdaterMock.On("UpdateMovieTitle", tc.movieId, tc.title).
					Return(tc.mockError).
					Maybe()
				movieUpdaterMock.On("UpdateMovieReleaseDate", tc.movieId, tc.releaseDate).
					Return(tc.mockError).
					Maybe()
				movieUpdaterMock.On("UpdateMovieRating", tc.movieId, tc.rating).
					Return(tc.mockError).
					Maybe()
				movieUpdaterMock.On("UpdateMovieDescription", tc.movieId, tc.description).
					Return(tc.mockError).
					Maybe()
			}

			handler := update.New(slogdiscard.NewDiscardLogger(), movieUpdaterMock)

			//input := fmt.Sprintf(`{"actor_id": %d, "name": "%s", "gender": "%s", "birthdate": "%s"}`,
			//	tc.actorId, tc.actorName, tc.gender, tc.birthdate)
			input := `{"movie_id": ` + strconv.Itoa(tc.movieId) + ``
			if tc.title != "" {
				input += `,"title": "` + tc.title + `"`
			}
			if tc.description != "" {
				input += `,"description": "` + tc.description + `"`
			}
			if tc.releaseDate != "" {
				input += `,"release_date": "` + tc.releaseDate + `"`
			}
			if tc.rating != 0 {
				input += `,"rating": ` + strconv.Itoa(tc.rating)
			}
			input += "}"

			req, err := http.NewRequest(http.MethodPost, "/movie/update", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp update.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
