package domains

var Genders = map[string]struct{}{"male": struct{}{}, "female": struct{}{}}

type Actor struct {
	ID       uint32 `json:"id"`
	FullName string `json:"fullName"`
	Gender   Gender `json:"gender"`
	Birthday Time   `json:"birthday" format:"2006-01-02"`
}

type Gender string

func (g Gender) IsValid() bool {
	_, ok := Genders[string(g)]
	return ok
}

type ActorWithFilms struct {
	Actor
	Films []*Film `json:"films"`
}
