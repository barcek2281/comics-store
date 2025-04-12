package model

type Comics struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`       // Comic title
	Author      string  `json:"author"`      // Author's name
	Description string  `json:"description"` // Description of the comic
	ReleaseDate string  `json:"releaseDate"` // Release date (can be formatted time.Time if needed)
	Price       float32 `json:"price"`       // Price of the comic
	Quantity    int32   `json:"quantity"`    // Available quantity
}
