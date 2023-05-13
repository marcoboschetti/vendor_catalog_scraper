package data

import (
	"bitbucket.org/marcoboschetti/catalogscraper/src/entities"
)

func GetAllProductURIsInSubcategory(subcategoryID string) (map[string]interface{}, error) {
	var res []struct {
		URI string
	}
	err := pgConnection.Model((*entities.Product)(nil)).
		Column("uri").
		Where("subcategory_id = ? ", subcategoryID).
		Select(&res)

	if err != nil {
		return nil, err
	}

	ans := map[string]interface{}{}
	for _, r := range res {
		ans[r.URI] = true
	}

	return ans, nil
}
