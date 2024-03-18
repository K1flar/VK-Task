package domains

var Roles = map[string]struct{}{"admin": struct{}{}, "viewer": struct{}{}}

type User struct {
	ID       uint32 `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type Role string

func (r Role) IsValidRole() bool {
	if _, ok := Roles[string(r)]; !ok {
		return false
	}
	return true
}
