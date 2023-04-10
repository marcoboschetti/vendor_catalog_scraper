package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getOneItem(productUrl string) (string, error) {
	req, err := http.NewRequest("GET", "https://www.argyros.com.pa/"+productUrl, nil)
	if err != nil {
		return "", err
	}

	req = addRequestHeaders(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	doc := sanitizePayload(string(body))

	doc = modifyProductToCatalog(doc)
	return doc.Html()
}

func modifyProductToCatalog(doc *goquery.Document) *goquery.Document {

	doc.Find(".left-slidebar").First().Remove()
	doc.Find("#top").First().Remove()
	doc.Find("#breadcrumb").First().Remove()
	doc.Find("#footer").First().Remove()
	doc.Find(".exzoom_btn").Remove()
	doc.Find(".col-sm-18").First().RemoveClass("col-sm-18")
	doc.Find(".rel-container").First().Remove()
	doc.Find(".afavorito").Remove()
	doc.Find("#content-wrapper-parent").SetAttr("style", "margin:0em;")
	doc.Find("#page-title").SetAttr("style", "font-family: 'Lato', sans-serif;")
	doc.Find("br").Remove() // Avoid <br> elements

	doc.Find("#exzoom").RemoveClass("hidden")

	// Remove all images
	doc.Find("#gen-img-ord").Remove()
	doc.Find("#carr-modelos-disponibles").Remove()
	for idx, img := range doc.Find(".exzoom_img_ul").Children().Nodes {
		if idx != 0 {
			img.Parent.RemoveChild(img)
		}
	}
	doc.Find("img").First().SetAttr("width", "100%")

	// Remove prices
	doc.Find("#seleccion-cant").Remove()
	doc.Find("#historial_pedidos").Remove()
	doc.Find("#product-info-left").Remove()
	doc.Find(".detail-price").Parent().Remove()
	doc.Find("#description-2").Children().Last().Remove() // Price

	descText, descTextOk := doc.Find("#description-2").Html()
	if descTextOk == nil {
		// descText := descText[:strings.LastIndex(descText, "|")] // Remove last "|"
		doc.Find("#description-2").SetHtml(descText)

		weightPrice, exists := doc.Find(".swatch-element").Attr("data-peso")
		if exists {
			doc.Find("#description-2").Children().Last().SetHtml("Peso: " + weightPrice) // Weight price
		}
	}
	doc.Find(".swatch-element.available.seltdp label").SetAttr("style", "background: var(--color-id); color: #fff;") // "Color selected for all sizes"

	// Fix sizes
	doc.Find("#seleccion-talla").Children().First().Remove()
	return doc
}

func sanitizePayload(body string) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))

	if err != nil {
		log.Fatal(err)
	}
	removeScripts(doc)
	return doc
}

func removeScripts(n *goquery.Document) {
	// if note is script tag
	for _, script := range n.Find("script").Nodes {
		script.Parent.RemoveChild(script)
	}
}
