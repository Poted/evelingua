package model

type Category struct {
	ID          string `json:"id"`          // Unique identifier for the category
	Name        string `json:"name"`        // Name of the category
	Description string `json:"description"` // Optional description
	Words       []Word `json:"words"`       // List of words associated with the category
	CreatedAt   string `json:"createdAt"`   // Timestamp of category creation
	UpdatedAt   string `json:"updatedAt"`   // Timestamp of last update
	UserID      string `json:"userId"`      // ID of the user who owns the category
}

type Word struct {
	Word        string `json:"word"`        // The word
	Language    string `json:"language"`    // Language of the word
	Translation string `json:"translation"` // Translation of the word
}
