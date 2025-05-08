package controllers

import (
	"context"
	"fmt"
	"net/http"
	"tesodev_interview/models"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var timeOut = 10 * time.Second

func CreateProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	var product models.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if product.Name == "" || product.Description == "" || product.Price == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "All feild must be required"})
	}

	product.Created_at = time.Now()
	product.Updated_at = time.Now()

	newProduct := models.Product{
		Id:          primitive.NewObjectID(),
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Created_at:  product.Created_at,
		Updated_at:  product.Updated_at,
	}

	result, err := productCollection.InsertOne(ctx, newProduct)
	if err != nil {
		fmt.Println("did not create product")
	}

	fmt.Println(result)

	return c.JSON(http.StatusOK, result)
}
