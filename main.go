package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"bitbucket.org/marcoboschetti/catalogscraper/src/data"
	"bitbucket.org/marcoboschetti/catalogscraper/src/service"
)

func main() {
	go data.SetDbConnection()

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
	registerLegacyEndpoints(public)
	registerAPI(public)

	r.Run()
}

func registerAPI(public *gin.RouterGroup) {
	// public.Use(scrappers.VendorAuthMiddleware)
	admin := public.Group("/admin")

	// For cron
	admin.POST("/scrap_categories", ScrapAndPersistCategories)
	admin.POST("/save_all_new_products", SaveAllNewProducts)

	// For pdf generator
	admin.GET("/categories", GetFullSubcategories)
	admin.GET("/products/subcategory/:subcategory_id/:page", GetPagedProducts)
	admin.GET("/products/:product_id", GetProduct)

}

func ScrapAndPersistCategories(c *gin.Context) {
	err := service.ScrapAndSaveSubcategories()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "categories map populated"})
}

func SaveAllNewProducts(c *gin.Context) {
	err := service.ScrapSaveAllNewProducts()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "all new products populated"})
}

func GetFullSubcategories(c *gin.Context) {
	categories, err := service.GetSubcategories()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

func GetPagedProducts(c *gin.Context) {
	subcategoryID := c.Param("subcategory_id")
	pageStr := c.Param("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	products, count, pageSize, err := service.GetPagedProducts(subcategoryID, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products, "total": count, "page_size": pageSize})
}

func GetProduct(c *gin.Context) {
	productID := c.Param("product_id")

	product, err := service.GetProductByID(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}
