package scrappers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/marcoboschetti/catalogscraper/src/entities"
	"bitbucket.org/marcoboschetti/catalogscraper/src/utils"
	"github.com/PuerkitoBio/goquery"
)

func GetOneProduct(subcategoryID, productUrl string) (*entities.Product, error) {
	req, err := http.NewRequest("GET", "https://www.argyros.com.pa/"+productUrl, nil)
	if err != nil {
		return nil, err
	}

	req = utils.AddRequestHeaders(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc := sanitizePayload(string(body))
	product := parseProductFromDoc(doc)

	product.SubcategoryID = subcategoryID
	product.DateCreated = time.Now()
	product.URI = productUrl

	return product, nil
}

func parseProductFromDoc(doc *goquery.Document) *entities.Product {
	prodId, _ := doc.Find("#idprod").Attr("value")
	prodIdDet, _ := doc.Find("#iddetalle").Attr("value")
	pageTitle, _ := doc.Find("#page-title").Find("span").Html()

	infoBlock, _ := doc.Find("#bloque-info-producto").Find("div").Find("span").Html()
	ref, _ := doc.Find("#referencia-producto-pvd").Find("span").Html() // Remove span, /span, spaces
	ref = strings.ReplaceAll(ref, "\t", "")
	ref = strings.ReplaceAll(ref, "\n", "")
	ref = strings.ReplaceAll(ref, " ", "")
	ref = strings.ReplaceAll(ref, "<span>", "")
	ref = strings.ReplaceAll(ref, "</span>", "")

	description, _ := doc.Find(".row .description").Html()
	description2, _ := doc.Find("#description-2").Html()

	price, _ := doc.Find("#vprice_visible").Html()

	var sizesArr []string
	doc.Find(".swatch-element.medium.available.seltdp").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("data-value")
		sizesArr = append(sizesArr, val)
	})
	sizes := strings.Join(sizesArr, "|")

	var images []string
	doc.Find("ul.exzoom_img_ul").Find("li").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Find("img").Attr("src")
		images = append(images, val)
	})

	product := entities.Product{
		ID:           fmt.Sprintf("%s_%s", prodId, prodIdDet),
		PageTitle:    pageTitle,
		ProductID:    prodId,
		ProductIDDet: prodIdDet,

		InfoBlock:    infoBlock,
		Ref:          ref,
		Description:  description,
		Description2: description2,

		Price:     price,
		Sizes:     sizes,
		ImageUrls: images,
	}

	return &product
}
