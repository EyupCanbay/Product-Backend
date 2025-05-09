package controllers

import (
	"context"
	"fmt"
	"net/http"
	"tesodev_interview/models"
	"tesodev_interview/responses"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var timeOut = 10 * time.Second

func validateObjectId(c echo.Context, id string) (string, error) {
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "id is not a object id", Data: &echo.Map{"data": err.Error()}})
	}
	return id, nil
}

func CreateProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	var product models.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "Invalid request body", Data: &echo.Map{"data": "Have an error"}})

	}

	if product.Name == "" || product.Description == "" || product.Price == "" {
		return c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "All feild must be required", Data: &echo.Map{"data": "Have an error"}})
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

	return c.JSON(http.StatusCreated, responses.ResponseHandler{Status: http.StatusCreated, Message: "Successfuly create product", Data: &echo.Map{"data": result}})
}

func UpdateProduct(c echo.Context) error {

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	var product models.Product
	userId, _ := validateObjectId(c, c.Param("product_id"))
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "Invalid data type only JSON", Data: &echo.Map{"data": err.Error()}})
	}

	if product.Name == "" || product.Description == "" || product.Price == "" {
		return c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "All feild must be required", Data: &echo.Map{"data": "Have an error"}})
	}

	product.Updated_at = time.Now()

	updateProduct := bson.M{
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"update_at":   product.Updated_at,
	}

	result, err := productCollection.UpdateOne(ctx, bson.M{"id": userId}, bson.M{"$set": updateProduct})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{Status: http.StatusInternalServerError, Message: "product did not update", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, responses.ResponseHandler{Status: http.StatusCreated, Message: "Successfuly update product", Data: &echo.Map{"data": result}})

}
