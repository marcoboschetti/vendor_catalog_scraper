package service

import (
	"strings"

	"bitbucket.org/marcoboschetti/catalogscraper/src/data"
	"bitbucket.org/marcoboschetti/catalogscraper/src/entities"
)

const pageSize = 12

func GetPagedProducts(subcategoryID string, page int) ([]entities.Product, int, int, error) {
	products, count, err := data.GetPagedProducs(subcategoryID, page-1, pageSize)
	if err != nil {
		return nil, -1, -1, err
	}

	return products, count, pageSize, nil
}

func GetProductByID(productID string) (*entities.Product, error) {
	product, err := data.GetByID[entities.Product](productID)
	if err != nil {
		return nil, err
	}

	sanitiseProduct(product)

	return product, nil
}

func sanitiseProduct(prod *entities.Product) {
	prod.Description2 = strings.Replace(prod.Description2, "<br/>", "|", -1)
	var desc2Arr = strings.Split(prod.Description2, "|")

	if len(desc2Arr) > 1 {
		desc2Arr = desc2Arr[:len(desc2Arr)-2]
		prod.Description2 = strings.Join(desc2Arr, " | ")
	}
}
