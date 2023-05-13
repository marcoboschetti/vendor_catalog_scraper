package scrappers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"bitbucket.org/marcoboschetti/catalogscraper/src/utils"
	"github.com/PuerkitoBio/goquery"
)

func GetOnlyCategoriesList() (map[string]string, error) {
	doc, closer, err := utils.URLToGoQuery("categories.php")
	if err != nil {
		return nil, err
	}
	defer closer()

	categoryURIList := map[string]string{}
	doc.Find(".collection-details a").Each(func(i int, anchor *goquery.Selection) {
		categoryHref, _ := anchor.Attr("href")
		title, _ := anchor.Attr("title")
		categoryURIList[title] = categoryHref
	})

	return categoryURIList, nil
}

func GetCategoriesList() (map[string]map[string]string, error) {
	doc, closer, err := utils.URLToGoQuery("categories.php")
	if err != nil {
		return nil, err
	}
	defer closer()

	categorySubcategoriesList := map[string]map[string]string{}
	doc.Find(".collection-details a").Each(func(i int, anchor *goquery.Selection) {
		categoryHref, _ := anchor.Attr("href")
		title, _ := anchor.Attr("title")

		subcategoryList, err := RequestSubcategories(categoryHref)
		if err != nil {
			fmt.Println("ERR iterating categories", err.Error())
		}

		categorySubcategoriesList[title] = subcategoryList
	})

	return categorySubcategoriesList, nil
}

func RequestSubcategories(categoryPath string) (map[string]string, error) {
	doc, closer, err := utils.URLToGoQuery(categoryPath)
	if err != nil {
		return nil, err
	}
	defer closer()

	subcategoriesMap := map[string]string{}
	doc.Find(".collection-details a").Each(func(i int, anchor *goquery.Selection) {
		subcategoryPath, _ := anchor.Attr("href")
		subCategoryTitle, _ := anchor.Attr("title")

		subcategoriesMap[subCategoryTitle] = subcategoryPath
	})

	return subcategoriesMap, nil
}

const pageSize = 32

func RequestCataloguePage(subcategoryPath string, page int) ([]ProductUrl, int, error) {
	var productUrls []ProductUrl

	totalProducts, pageProductsUrl, err := iterateCataloguePage(subcategoryPath, page*pageSize, (page+1)*pageSize)
	if err != nil {
		return nil, -1, err
	}

	productUrls = append(productUrls, pageProductsUrl...)
	return productUrls, totalProducts, nil
}

func RequestCatalogue(subcategoryPath string, numberOfPages int) ([]ProductUrl, error) {
	totalItems := -1
	var productUrls []ProductUrl

	pageCount := 0
	for pageStart := 0; ; {
		totalItemsInPage, pageProductsUrl, err := iterateCataloguePage(subcategoryPath, pageStart, pageStart+pageSize)
		if err != nil {
			return nil, err
		}

		productUrls = append(productUrls, pageProductsUrl...)

		if pageStart == 0 {
			totalItems = totalItemsInPage
		}

		pageStart += pageSize
		pageCount += 1
		if pageStart >= totalItems || pageCount >= numberOfPages {
			break
		}
	}

	return productUrls, nil
}

func iterateCataloguePage(subcategoryPath string, start, end int) (int, []ProductUrl, error) {
	form := url.Values{}
	fmt.Println("Page from ", start, "TO", end)
	form.Add("listar_catalogo_url", "/"+subcategoryPath)
	form.Add("pos_i", fmt.Sprintf("%d", start))
	form.Add("pos_f", fmt.Sprintf("%d", end))
	form.Add("carga_i", fmt.Sprintf("%t", start == 0))
	req, err := http.NewRequest("POST", "https://www.argyros.com.pa/fn/fn-cataloquer.php", strings.NewReader(form.Encode()))
	if err != nil {
		return 0, nil, err
	}
	req = utils.AddRequestHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, nil, fmt.Errorf("error getting %s with status code %d: %v", subcategoryPath, resp.StatusCode, resp)
	}

	ansBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	ansPayload := map[string]interface{}{}
	err = json.Unmarshal(ansBody, &ansPayload)
	if err != nil {
		return 0, nil, err
	}

	var totalItems = 0
	if start == 0 {
		totalItems, err = strconv.Atoi(fmt.Sprintf("%v", ansPayload["contador"]))
		if err != nil {
			return 0, nil, err
		}
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fmt.Sprintf("%v", ansPayload["resultados"])))
	if err != nil {
		return 0, nil, err
	}

	var productsUrls []ProductUrl
	doc.Find(".container_item").Each(func(i int, anchor *goquery.Selection) {
		productUrl, _ := anchor.Attr("href")
		isNew := len(anchor.Find(".sale_text").Nodes) > 0

		productsUrls = append(productsUrls, ProductUrl{productUrl, isNew})
	}).Attr("href")

	return totalItems, productsUrls, nil
}

type ProductUrl struct {
	URL   string `json:"url"`
	IsNew bool   `json:"is_new"`
}
