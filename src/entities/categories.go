package entities

import "time"

type (
	Category struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		URI         string    `json:"uri"`
		DateCreated time.Time `json:"date_created"`
	}

	Subcategory struct {
		ID         string `json:"id"`
		CategoryID string `json:"category_id"`
		Name       string `json:"name"`
		URI        string `json:"uri"`

		DateCreated         time.Time `json:"date_created"`
		LastProductsUpdated time.Time `json:"last_products_updated"`
		LastProductURI      string    `json:"last_product_uri"`
	}

	Product struct {
		ID            string    `json:"id"`
		SubcategoryID string    `json:"subcategory_id"`
		URI           string    `json:"-"`
		DateCreated   time.Time `json:"date_created"`

		// Data to render
		ProductID    string   `json:"product_id"`
		ProductIDDet string   `json:"product_id_det"`
		PageTitle    string   `json:"page_title"`
		InfoBlock    string   `json:"info_block"`
		Ref          string   `json:"ref"`
		Description  string   `json:"description"`
		Description2 string   `json:"description_2"`
		Price        string   `json:"-"`
		Sizes        string   `json:"sizes"`
		ImageUrls    []string `json:"image_urls"`
	}
)
