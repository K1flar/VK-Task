package pagination

import "net/http"

const (
	QueryFilmName      = "film"
	QueryActorName     = "actor"
	QueryOrderByName   = "sort"
	QueryDirectionName = "direct"

	DefaultSortBy        = "rating"
	DefaultSortDirection = "desc"
)

var (
	fieldsForOrderFilms = map[string]struct{}{"name": struct{}{}, "rating": struct{}{}, "release_date": struct{}{}}
)

type FilmFilter struct {
	Pagination        *Pagination
	NameContains      string
	ActorNameContains string
	OrderBy           string
	Direction         string
}

type ActorsFilter struct {
	Pagination       *Pagination `json:"pagination"`
	FullNameContains string      `json:"fullNameContains"`
}

func (f *FilmFilter) Validate() {
	f.Pagination.ValidatePagination()
	if _, ok := fieldsForOrderFilms[f.OrderBy]; !ok {
		f.OrderBy = DefaultSortBy
		f.Direction = DefaultSortDirection
	}
	if f.Direction != "asc" && f.Direction != "desc" {
		f.Direction = "asc"
	}
}

func NewFilmFilterFromRequest(r *http.Request) *FilmFilter {
	nameContains := r.URL.Query().Get(QueryFilmName)
	actorNameContains := r.URL.Query().Get(QueryActorName)
	orderBy := r.URL.Query().Get(QueryOrderByName)
	direction := r.URL.Query().Get(QueryDirectionName)
	return &FilmFilter{
		Pagination:        NewFromRequest(r),
		NameContains:      nameContains,
		ActorNameContains: actorNameContains,
		OrderBy:           orderBy,
		Direction:         direction,
	}
}

func NewActorFilterFromRequest(r *http.Request) *ActorsFilter {
	fullNameContains := r.URL.Query().Get(QueryActorName)
	return &ActorsFilter{
		Pagination:       NewFromRequest(r),
		FullNameContains: fullNameContains,
	}
}
