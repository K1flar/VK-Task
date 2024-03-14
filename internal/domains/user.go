package domains

type User struct {
	ID       uint32
	Login    string
	Password string
	Role     string
}
