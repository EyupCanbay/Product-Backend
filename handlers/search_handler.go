package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"tesodev_interview/models"
	"tesodev_interview/responses"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SearchProducts(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	name := c.QueryParam("name")
	exact := c.QueryParam("exact") == "true"
	priceMin, _ := strconv.ParseFloat(c.QueryParam("price_min"), 64)
	priceMax, _ := strconv.ParseFloat(c.QueryParam("price_max"), 64)
	sort := c.QueryParam("sort") // "asc" or "desc"
	limitParam := c.QueryParam("limit")
	pageParam := c.QueryParam("page")

	limit := 1
	page := 1

	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
		if page < 0 {
			page = 1
		}
	}
	fmt.Println(page)

	if limitParam != "" {
		limit, _ = strconv.Atoi(limitParam)
		if limit < 0 {
			limit = 1
		}
	}

	skip := (page - 1) * limit

	filter := bson.M{}

	if priceMin > 0 || priceMax > 0 {
		priceFilter := bson.M{}
		if priceMin > 0 {
			priceFilter["$gte"] = priceMin
		}
		if priceMax > 0 {
			priceFilter["$lte"] = priceMax
		}
		filter["price"] = priceFilter
	}

	var sortOrder int
	if sort == "desc" {
		sortOrder = -1
	} else if sort == "asc" {
		sortOrder = 1
	} else {
		sortOrder = 0 // do not sort
	}

	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	if sortOrder != 0 {
		findOptions.SetSort(bson.D{{Key: "price", Value: sortOrder}})
	}

	if name != "" {
		if exact {
			filter["name"] = name // exact matches
		} else {
			filter["name"] = bson.M{"$regex": name, "$options": "i"} // partial matches
		}
	}

	cursor, err := productCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ResponseHandler{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, responses.ResponseHandler{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &echo.Map{"products": products},
	})
}
