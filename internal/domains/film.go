package domains

type Film struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ReleaseDate Time   `json:"releaseDate" format:"2006-01-02"`
	Rating      int    `json:"rating"`
}
