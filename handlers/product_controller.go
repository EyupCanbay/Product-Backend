package handlers

import (
	"context"
	"net/http"
	"tesodev_interview/models"
	"tesodev_interview/responses"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var timeOut = 10 * time.Second

func CreateProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	var product models.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "Invalid body request"}})

	}

	if product.Name == "" || product.Description == "" || product.Price <= 0 {
		return c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "All feild must be required"}})
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
		return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, responses.ResponseHandler{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": result}})
}

func UpdateProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	var product models.Product
	productId, err := primitive.ObjectIDFromHex(c.Param("product_id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if product.Name == "" || product.Description == "" || product.Price <= 0 {
		return c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "All feild must be required"}})
	}

	product.Updated_at = time.Now()

	updateProduct := bson.M{
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"updated_at":  product.Updated_at,
	}

	result, err := productCollection.UpdateOne(ctx, bson.M{"_id": productId}, bson.M{"$set": updateProduct})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{Status: http.StatusInternalServerError, Message: "product did not update", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, responses.ResponseHandler{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": result}})

}

func GetAProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	var product models.Product

	productId, err := primitive.ObjectIDFromHex(c.Param("product_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	err = productCollection.FindOne(ctx, bson.M{"_id": productId}).Decode(&product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, responses.ResponseHandler{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": &product}})

}

func GetAllProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	var products []models.Product

	results, err := productCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleProduct models.Product
		if err = results.Decode(&singleProduct); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}

		products = append(products, singleProduct)
	}

	return c.JSON(http.StatusOK, responses.ResponseHandler{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": products}})
}

func DeleteProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	productId, _ := primitive.ObjectIDFromHex(c.Param("product_id"))

	result, err := productCollection.DeleteOne(ctx, bson.M{"_id": productId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.JSON(http.StatusNotFound, responses.ResponseHandler{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": "product with id not found"}})
	}

	return c.JSON(http.StatusOK, responses.ResponseHandler{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "Successfuly delete products"}})
}

func UpdateSingleFeild(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	productId, _ := primitive.ObjectIDFromHex(c.Param("product_id"))

	var product models.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "Invalid request body"}})
	}

	update := bson.M{}
	if product.Name != "" {
		update["name"] = product.Name
	}
	if product.Price <= 0 {
		update["price"] = product.Price
	}
	if product.Description != "" {
		update["description"] = product.Description
	}

	if len(update) == 0 {
		return c.JSON(http.StatusBadRequest, responses.ResponseHandler{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "No fields to update"}})
	}

	filter := bson.M{"_id": productId}
	updateDoc := bson.M{"$set": update}
	opts := options.Update().SetUpsert(false)

	result, err := productCollection.UpdateOne(ctx, filter, updateDoc, opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": "failed update product"}})
	}

	if result.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, responses.ResponseHandler{Status: http.StatusNotFound, Message: "error", Data: &echo.Map{"data": "product not found"}})
	}

	return c.JSON(http.StatusOK, responses.ResponseHandler{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": &result}})
}
