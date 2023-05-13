package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"bitbucket.org/marcoboschetti/catalogscraper/src/data"
	"bitbucket.org/marcoboschetti/catalogscraper/src/entities"
	"bitbucket.org/marcoboschetti/catalogscraper/src/scrappers"
)

func loginServer() {
	// Login
	loginErr := scrappers.Login(os.Getenv("VENDOR_USERNAME"), os.Getenv("VENDOR_PASSWORD"))
	if loginErr != nil {
		fmt.Println("On init login error:", loginErr.Error())
	}
}

func ScrapAndSaveSubcategories() error {
	loginServer()

	categoryMap, err := scrappers.GetOnlyCategoriesList()
	if err != nil {
		return err
	}

	for categoryTitle, categoryURI := range categoryMap {

		// Insert
		category := entities.Category{
			ID:          strings.ToLower(categoryTitle), //uuid.NewV4().String(),
			Name:        categoryTitle,
			URI:         categoryURI,
			DateCreated: time.Now(),
		}
		err = data.InsertNew(category)
		if err != nil {
			return err
		}

		// Populate
		subcategoryList, err := scrappers.RequestSubcategories(categoryURI)
		if err != nil {
			fmt.Println("ERR iterating categories", err.Error())
		}
		for subcategoryTitle, subcategoryURI := range subcategoryList {
			subcategory := entities.Subcategory{
				ID:                  category.ID + "_" + strings.ToLower(subcategoryTitle), // Remove last spaces, replace space by underscore
				CategoryID:          category.ID,
				Name:                subcategoryTitle,
				URI:                 subcategoryURI,
				DateCreated:         time.Now(),
				LastProductsUpdated: time.Time{},
				LastProductURI:      "",
			}
			err = data.InsertNew(subcategory)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ScrapSaveAllNewProducts() error {
	loginServer()

	subcategories, err := data.GetAll[entities.Subcategory]()
	if err != nil {
		return err
	}

	for subcategoryIdx, subcategory := range subcategories {
		fmt.Printf("🟢 Downloading %d/%d - %s.\n", (subcategoryIdx + 1), len(subcategories), subcategory.Name)

		retrievedSubcategoryIDs, err := data.GetAllProductURIsInSubcategory(subcategory.ID)
		if err != nil {
			fmt.Println("Failed to get subcategory prev items:", err.Error())
			return err
		}
		fmt.Printf("	Prev count was %d\n", len(retrievedSubcategoryIDs))

		err = saveNewProductsFromSubcategory(subcategory, retrievedSubcategoryIDs) // TODO, goroutines?
		if err != nil {
			return err
		}
	}

	return nil
}

func saveNewProductsFromSubcategory(subcategory entities.Subcategory, retrievedItems map[string]interface{}) error {
	page := 0
	newestProductURI := ""
	totalProducts := -1

	for {
		productsURL, curTotalProducts, err := scrappers.RequestCataloguePage(subcategory.URI, page)
		if err != nil {
			return err
		}

		if totalProducts == -1 {
			totalProducts = curTotalProducts
		}

		// Check if product was already retrieved
		forceStop := false
		var productsFromPage []*entities.Product
		for _, productURI := range productsURL {
			if productURI.URL == subcategory.LastProductURI {
				fmt.Println("Finished because item already retrieved")
				forceStop = true
				break
			}

			if _, ok := retrievedItems[productURI.URL]; ok {
				fmt.Println("Skipping because already in DB", productURI.URL)
				continue
			}

			product, err := scrappers.GetOneProduct(subcategory.ID, productURI.URL)
			if err != nil {
				fmt.Println("Failed with", productURI.URL, ":", err.Error())
				// return err
				continue
			}

			productsFromPage = append(productsFromPage, product)
			retrievedItems[productURI.URL] = true
			fmt.Printf("		Got prod %d/%d\n", len(productsFromPage)+page*32, totalProducts)
		}

		if len(productsFromPage) > 0 {
			err = data.InsertMany(productsFromPage)
			if err != nil {
				fmt.Println("Failed to save in DB bulk:", err.Error())

				if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					// Only break all if error was not already persisted items
					return err
				}
			}
		}

		if page == 0 && newestProductURI == "" {
			newestProductURI = productsURL[0].URL // Newest for subcategory track
		}

		// Break for full scan
		fmt.Printf("	Persisted %d/%d\n", len(retrievedItems), totalProducts)
		if forceStop || len(retrievedItems) >= totalProducts {
			break
		}
		page++
	}

	// Update subcategory
	subcategory.LastProductURI = newestProductURI
	subcategory.LastProductsUpdated = time.Now()
	err := data.Update(subcategory)
	return err
}
