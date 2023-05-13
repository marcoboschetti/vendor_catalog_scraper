package entities

type (
	CategoryDTO struct {
		ID            string           `json:"id"`
		Name          string           `json:"name"`
		URI           string           `json:"uri"`
		Subcategories []SubcategoryDTO `json:"subcategories"`
	}

	SubcategoryDTO struct {
		ID         string `json:"id"`
		CategoryID string `json:"category_id"`
		Name       string `json:"name"`
		URI        string `json:"uri"`
	}
)
