package data

import (
	"bitbucket.org/marcoboschetti/catalogscraper/src/entities"
)

func GetPagedProducs(subcategoryID string, page, pageSize int) ([]entities.Product, int, error) {
	var res []entities.Product
	err := pgConnection.Model(&res).
		Where("subcategory_id = ?", subcategoryID).
		Column("id", "info_block", "page_title", "image_urls").
		Offset(page * pageSize).
		Limit(pageSize).
		Order("id DESC").
		Select()

	if err != nil {
		return nil, -1, err
	}

	count, err := pgConnection.Model(&res).
		Where("subcategory_id = ?", subcategoryID).
		Count()

	if err != nil {
		return nil, -1, err
	}

	return res, count, nil
}
