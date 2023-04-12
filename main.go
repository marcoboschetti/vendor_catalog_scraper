package main

import (
	"fmt"
	"html"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("$PORT must be set. Defaulted to 8080")
		port = "8080"
	}
	gin.Default()
	r := gin.New()

	// *************** SITE **************
	r.StaticFS("/site", gin.Dir("site", false))
	r.StaticFile("/", "./site/index.html")

	// *************** PROXY **************
	r.Any("/assets/*proxyPath", proxyAssets("/assets"))
	r.Any("/fn/*proxyPath", proxyAssets("/fn"))

	// *************** API **************
	r.POST("/auth/vendor_login", loginVendor)

	public := r.Group("/api")
	public.Use(VendorAuthMiddleware)
	public.GET("/all_categories", listAllCategories)
	public.GET("/subcategory_products", getSubcategoryProducts)
	public.GET("/category_products", listCategoryProducts)
	public.POST("/get_product", getProduct)

	r.Run()
}

func proxyAssets(subPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse("https://www.argyros.com.pa")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			for k, v := range c.Request.Header {
				req.Header.Add(k, v[0])
			}
			req.URL.Scheme = "https"
			req.URL.Host = "www.argyros.com.pa"
			req.Host = "www.argyros.com.pa"
			req.URL.Path = subPath + c.Param("proxyPath")
			req.URL.RawQuery = c.Request.URL.RawQuery
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func listAllCategories(c *gin.Context) {
	categoryMap, err := getCategoriesList()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category_map": categoryMap})
}

func getSubcategoryProducts(c *gin.Context) {
	subcategoryURL := c.Request.URL.Query().Get("subcat_url")
	numberOfPagesStr := c.Request.URL.Query().Get("number_of_pages")

	numberOfPages, err := strconv.Atoi(numberOfPagesStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	productsUrls, err := requestCatalogue(subcategoryURL, numberOfPages)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products_urls": productsUrls})
}

func listCategoryProducts(c *gin.Context) {
	categoryMap, err := getCategoriesList()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category_map": categoryMap})
}

func loginVendor(c *gin.Context) {
	input := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	err = login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func getProduct(c *gin.Context) {
	input := struct {
		ProductUrl string `json:"product_url"`
	}{}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// Retrieve items
	itemPayload, err := getOneItem(input.ProductUrl)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product_http": html.EscapeString(itemPayload)})
}
