package actorhandler

import (
	"film_library/internal/domains"
	mock_services "film_library/internal/services/mocks"
	"film_library/pkg/mux"
	"film_library/pkg/pagination"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestActorHandlerGetActors(t *testing.T) {
	type mockBehavior func(r *mock_services.MockActorService, filter pagination.ActorsFilter)

	time, _ := time.Parse(time.DateOnly, "2022-06-23")

	tests := []struct {
		name                 string
		queryParams          string
		inputFilter          pagination.ActorsFilter
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Correct",
			queryParams: `page=1&size=5`,
			inputFilter: pagination.ActorsFilter{Pagination: pagination.New(1, 5)},
			mockBehavior: func(r *mock_services.MockActorService, filter pagination.ActorsFilter) {
				r.EXPECT().GetActorsWithFilms(&filter).Return([]*domains.ActorWithFilms{
					{
						Actor: domains.Actor{ID: 1, FullName: "Denis", Gender: "male", Birthday: domains.Time(time)},
						Films: []*domains.Film{
							{
								ID:          1,
								Name:        "Test",
								Description: "",
								ReleaseDate: domains.Time(time),
								Rating:      10,
							},
						},
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"id":1,"fullName":"Denis","gender":"male","birthday":"2022-06-23","films":[{"id":1,"name":"Test","description":"","releaseDate":"2022-06-23","rating":10}]}]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_services.NewMockActorService(c)
			handler := ActorHandler{service: service}
			tc.mockBehavior(service, tc.inputFilter)

			r := mux.New()
			r.HandleFunc("GET /api/actors", handler.GetActorsWithFilms)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/actors?"+tc.queryParams, nil)

			r.ServeHTTP(w, req)

			if tc.expectedStatusCode != w.Code {
				t.Errorf("expected: %d\ngot: %d", tc.expectedStatusCode, w.Code)
			}

			if tc.expectedResponseBody != w.Body.String() {
				t.Errorf("expected: %s\ngot: %s", tc.expectedResponseBody, w.Body.String())
			}
		})
	}
}
