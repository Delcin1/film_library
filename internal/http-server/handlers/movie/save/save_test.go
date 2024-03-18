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

	"film_library/internal/http-server/handlers/movie/save"
	"film_library/internal/http-server/handlers/movie/save/mocks"
	"film_library/internal/lib/logger/handlers/slogdiscard"
)

const bigDesription = "Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat. Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi.Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat. Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui"

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name        string
		title       string
		description string
		releaseDate string
		rating      int
		actorsIds   []int
		respError   string
		mockError   error
	}{
		{
			name:        "Success",
			title:       "Best movie",
			description: "Best of the best",
			releaseDate: "2000-01-01",
			rating:      10,
			actorsIds:   []int{1, 2, 3},
		},
		{
			name:        "Empty Title",
			title:       "",
			description: "Best of the best",
			releaseDate: "2000-01-01",
			rating:      10,
			actorsIds:   []int{1, 2, 3},
			respError:   "field title is not valid",
		},
		{
			name:        "Invalid Description",
			title:       "Best movie",
			description: bigDesription,
			releaseDate: "2000-01-01",
			rating:      10,
			actorsIds:   []int{1, 2, 3},
			respError:   "field description is not valid",
		},
		{
			name:        "Invalid Release Date",
			title:       "Best movie",
			description: "Best of the best",
			releaseDate: "2000-01-01 00:00:00",
			rating:      10,
			actorsIds:   []int{1, 2, 3},
			respError:   "field release_date is not valid",
		},
		{
			name:        "Invalid Rating",
			title:       "Best movie",
			description: "Best of the best",
			releaseDate: "2000-01-01",
			rating:      11,
			actorsIds:   []int{1, 2, 3},
			respError:   "field rating is not valid",
		},
		{
			name:        "Invalid Actors",
			title:       "Best movie",
			description: "Best of the best",
			releaseDate: "2000-01-01",
			rating:      10,
			actorsIds:   []int{1, -2, 3},
			respError:   "field actors_ids is not valid",
		},
		{
			name:        "SaveMovie Error",
			title:       "Best movie",
			description: "Best of the best",
			releaseDate: "2000-01-01",
			rating:      10,
			actorsIds:   []int{1, 2, 3},
			respError:   "failed to save movie",
			mockError:   errors.New("failed to save movie"),
		},
	}

	for _, tc := range cases {
		// tc := tc // go version < 1.22

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			movieSaverMock := mocks.NewMovieSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				movieSaverMock.On("SaveMovie", tc.title, tc.description, tc.releaseDate, tc.rating, tc.actorsIds).
					Return(1, tc.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), movieSaverMock)

			actorsStr := fmt.Sprintf("%v", tc.actorsIds)
			actorsJSONStr := "[" + strings.Join(strings.Split(actorsStr[1:len(actorsStr)-1], " "), ", ") + "]"
			input := fmt.Sprintf(`{"title": "%s", "description": "%s", "release_date": "%s", "rating": %d, "actors_ids": %s}`,
				tc.title, tc.description, tc.releaseDate, tc.rating, actorsJSONStr)

			req, err := http.NewRequest(http.MethodPost, "/movie/save", bytes.NewReader([]byte(input)))
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
