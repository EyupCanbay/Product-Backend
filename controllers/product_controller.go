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

	return c.JSON(http.StatusCreated, responses.ResponseHandler{Status: http.StatusCreated, Message: "Successfuly create product", Data: &echo.Map{"data": result}})
}

func UpdateProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	var product models.Product
	productId, err := primitive.ObjectIDFromHex(c.Param("product_id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "id is not a object id", Data: &echo.Map{"data": err.Error()}})
	}

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

	result, err := productCollection.UpdateOne(ctx, bson.M{"id": productId}, bson.M{"$set": updateProduct})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{Status: http.StatusInternalServerError, Message: "product did not update", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, responses.ResponseHandler{Status: http.StatusOK, Message: "Successfuly update product", Data: &echo.Map{"data": result}})

}

func GetAProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	var product models.Product

	productId, err := primitive.ObjectIDFromHex(c.Param("product_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "id is not a object id", Data: &echo.Map{"data": err.Error()}})
	}

	err = productCollection.FindOne(ctx, bson.M{"_id": productId}).Decode(&product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{Status: http.StatusInternalServerError, Message: "product did not fetch", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, responses.ResponseHandler{Status: http.StatusOK, Message: "Successfuly fetch product", Data: &echo.Map{"data": product}})

}

func GetAllProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	var products []models.Product

	results, err := productCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{Status: http.StatusInternalServerError, Message: "data did not fetch", Data: &echo.Map{"data": err.Error()}})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleProduct models.Product
		if err = results.Decode(&singleProduct); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{Status: http.StatusInternalServerError, Message: "data did not fetch", Data: &echo.Map{"data": err.Error()}})
		}

		products = append(products, singleProduct)
	}

	return c.JSON(http.StatusOK, responses.ResponseHandler{Status: http.StatusOK, Message: "Successfuly fetch products", Data: &echo.Map{"data": products}})
}

func DeleteProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	productId, _ := primitive.ObjectIDFromHex(c.Param("product_id"))

	result, err := productCollection.DeleteOne(ctx, bson.M{"_id": productId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{Status: http.StatusInternalServerError, Message: "product did not delete", Data: &echo.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.JSON(http.StatusNotFound, responses.ResponseHandler{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": "product with id not found"}})
	}

	return c.JSON(http.StatusOK, responses.ResponseHandler{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "Successfuly delete products"}})

}
