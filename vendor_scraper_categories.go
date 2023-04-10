package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getCategoriesList() (map[string]map[string]string, error) {
	doc, closer, err := urlToGoQuery("categories.php")
	if err != nil {
		return nil, err
	}
	defer closer()

	categorySubcategoriesList := map[string]map[string]string{}
	doc.Find(".collection-details a").Each(func(i int, anchor *goquery.Selection) {
		categoryHref, _ := anchor.Attr("href")
		title, _ := anchor.Attr("title")

		subcategoryList, err := requestSubcategories(categoryHref)
		if err != nil {
			fmt.Println("ERR iterating categories", err.Error())
		}

		categorySubcategoriesList[title] = subcategoryList
	})

	return categorySubcategoriesList, nil
}

func requestSubcategories(categoryPath string) (map[string]string, error) {
	doc, closer, err := urlToGoQuery(categoryPath)
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

func requestCatalogue(subcategoryPath string, isLastPage bool) ([]string, error) {
	totalItems := -1
	var productUrls []string

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
		if pageStart >= totalItems || isLastPage {
			break
		}
	}

	return productUrls, nil
}

func iterateCataloguePage(subcategoryPath string, start, end int) (int, []string, error) {
	form := url.Values{}
	form.Add("listar_catalogo_url", "/"+subcategoryPath)
	form.Add("pos_i", fmt.Sprintf("%d", start))
	form.Add("pos_f", fmt.Sprintf("%d", end))
	form.Add("carga_i", fmt.Sprintf("%t", start == 0))
	req, err := http.NewRequest("POST", "https://www.argyros.com.pa/fn/fn-cataloquer.php", strings.NewReader(form.Encode()))
	if err != nil {
		return 0, nil, err
	}
	req = addRequestHeaders(req)

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

	var productsUrls []string
	doc.Find(".container_item").Each(func(i int, anchor *goquery.Selection) {
		productUrl, _ := anchor.Attr("href")
		productsUrls = append(productsUrls, productUrl)
	}).Attr("href")

	return totalItems, productsUrls, nil
}
