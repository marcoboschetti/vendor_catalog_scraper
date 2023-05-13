package service

import (
	"bitbucket.org/marcoboschetti/catalogscraper/src/data"
	"bitbucket.org/marcoboschetti/catalogscraper/src/entities"
)

func GetSubcategories() ([]*entities.CategoryDTO, error) {
	categories, err := data.GetAll[entities.Category]()
	if err != nil {
		return nil, err
	}

	var ans []*entities.CategoryDTO
	m := map[string]*entities.CategoryDTO{}

	for _, c := range categories {
		cDTO := entities.CategoryDTO{
			ID:   c.ID,
			Name: c.Name,
			URI:  c.URI,
		}
		m[c.ID] = &cDTO
		ans = append(ans, &cDTO)
	}

	subcategories, err := data.GetAll[entities.Subcategory]()
	if err != nil {
		return nil, err
	}
	for _, s := range subcategories {
		sDTO := entities.SubcategoryDTO{
			ID:         s.ID,
			CategoryID: s.CategoryID,
			Name:       s.Name,
			URI:        s.URI,
		}
		m[s.CategoryID].Subcategories = append(m[s.CategoryID].Subcategories, sDTO)
	}

	return ans, nil
}
