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

	"film_library/internal/http-server/handlers/actor/update"
	"film_library/internal/http-server/handlers/actor/update/mocks"
	"film_library/internal/lib/logger/handlers/slogdiscard"
)

const bigName = "Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea co"

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		actorId   int
		actorName string
		gender    string
		birthdate string
		respError string
		mockError error
	}{
		{
			name:      "Success update name",
			actorId:   1,
			actorName: "Nikita",
		},
		{
			name:      "Invalid name",
			actorId:   1,
			actorName: bigName,
			respError: "field name is not valid",
		},
		{
			name:    "Success update gender",
			actorId: 1,
			gender:  "male",
		},
		{
			name:      "Invalid gender",
			actorId:   1,
			gender:    "abc",
			respError: "field gender is not valid",
		},
		{
			name:      "Success update birthdate",
			actorId:   1,
			birthdate: "2000-01-01",
		},
		{
			name:      "Invalid birthdate",
			actorId:   1,
			birthdate: "1900.01.01",
			respError: "field birthdate is not valid",
		},
		{
			name:      "Success update all",
			actorId:   1,
			actorName: "Nikita",
			gender:    "male",
			birthdate: "2000-01-01",
		},
		{
			name:      "Invalid actorId",
			actorId:   0,
			respError: "field actor_id is not valid",
		},
		{
			name:      "UpdateActor Error",
			actorId:   1,
			actorName: "Nikita",
			gender:    "male",
			birthdate: "2000-01-01",
			respError: "failed to update actor name",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		// tc := tc // go version < 1.22

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actorUpdaterMock := mocks.NewActorUpdater(t)

			if tc.respError == "" || tc.mockError != nil {
				actorUpdaterMock.On("UpdateActorName", tc.actorId, tc.actorName).
					Return(tc.mockError).
					Maybe()
				actorUpdaterMock.On("UpdateActorGender", tc.actorId, tc.gender).
					Return(tc.mockError).
					Maybe()
				actorUpdaterMock.On("UpdateActorBirthdate", tc.actorId, tc.birthdate).
					Return(tc.mockError).
					Maybe()
			}

			handler := update.New(slogdiscard.NewDiscardLogger(), actorUpdaterMock)

			//input := fmt.Sprintf(`{"actor_id": %d, "name": "%s", "gender": "%s", "birthdate": "%s"}`,
			//	tc.actorId, tc.actorName, tc.gender, tc.birthdate)
			input := `{"actor_id": ` + strconv.Itoa(tc.actorId) + ``
			if tc.actorName != "" {
				input += `,"name": "` + tc.actorName + `"`
			}
			if tc.gender != "" {
				input += `,"gender": "` + tc.gender + `"`
			}
			if tc.birthdate != "" {
				input += `,"birthdate": "` + tc.birthdate + `"`
			}
			input += "}"

			req, err := http.NewRequest(http.MethodPost, "/actor/update", bytes.NewReader([]byte(input)))
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
