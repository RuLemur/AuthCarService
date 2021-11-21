package controller

import (
	"fmt"
	"github.com/RuLemur/CarService/pkg/endpoint"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CarSearch godoc
// @Summary Search car by brand and model
// @Description get string by ID
// @Tags Cars
// @Accept  json
// @Produce  json
// @Param cars body endpoint.CarSearchRequest true "Car Search"
// @Success 200 {object} endpoint.CarSearchResponse
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /cars/search [post]
func (c *Controller) CarSearch(ctx *gin.Context) {
	var model endpoint.CarSearchRequest
	err := ctx.ShouldBindJSON(&model)
	if err != nil {
		response := Error{
			Message: fmt.Sprint("Fail parse!"),
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	carSearchResponse, err := c.grpcClient.CarSearch(c.ctx, &model)
	if err != nil {
		response := Error{
			Message: fmt.Sprint("Fail server!"),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, carSearchResponse)
	return
}

// AddCar godoc
// @Summary Add car to user
// @Description add car to user
// @Tags Cars
// @Accept  json
// @Produce  json
// @Param cars body endpoint.AddCarRequest true "Add Car"
// @Success 200 {object} endpoint.AddCarResponse
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /cars/add [post]
func (c *Controller) AddCar(ctx *gin.Context) {
	var model endpoint.AddCarRequest
	err := ctx.ShouldBindJSON(&model)
	if err != nil {
		response := Error{
			Message: fmt.Sprintf("Fail parse!%s", err.Error()),
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	addCarResponse, err := c.grpcClient.AddCar(c.ctx, &model)
	if err != nil {
		response := Error{
			Message: fmt.Sprintf("Fail server! %s", err.Error()),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, addCarResponse)
	return
}
