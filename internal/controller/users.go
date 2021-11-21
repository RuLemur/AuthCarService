package controller

import (
	"context"
	"fmt"
	"github.com/RuLemur/CarService/pkg/endpoint"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Controller example
type Controller struct {
	ctx        context.Context
	grpcClient endpoint.CarServiceClient
}

// NewController example
func NewController(ctx context.Context, grpcClient endpoint.CarServiceClient) *Controller {
	return &Controller{ctx: ctx, grpcClient: grpcClient}
}

type Error struct {
	Message string
}

// GetUser godoc
// @Summary Get User
// @Description get user by id
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} endpoint.GetUserResponse
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /users/{id} [get]
func (c *Controller) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		response := Error{
			Message: fmt.Sprintf("Fail to parse userID"),
		}
		ctx.JSON(http.StatusBadRequest, response)

		return
	}
	getUserResponse, err := c.grpcClient.GetUser(c.ctx, &endpoint.GetUserRequest{
		Id: int64(userID),
	})
	if err != nil {
		response := Error{
			Message: fmt.Sprintf("Fail server %d!", userID),
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	ctx.JSON(http.StatusOK, getUserResponse)
}

// AddUser godoc
// @Summary Add new user
// @Description Add new user
// @Tags Users
// @Accept  json
// @Produce  json
// @Param userInfo body endpoint.AddUserRequest true "Add User"
// @Success 200 {object} endpoint.AddUserResponse
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /users/add [post]
func (c *Controller) AddUser(ctx *gin.Context) {

	var user endpoint.AddUserRequest
	err := ctx.BindJSON(&user)
	if err != nil {
		errResponse := Error{
			Message: fmt.Sprintf("Fail to parse request!"),
		}
		ctx.JSON(http.StatusBadRequest, errResponse)
	}
	addUserResponse, err := c.grpcClient.AddUser(c.ctx, &user)
	if err != nil {
		errResponse := Error{
			Message: fmt.Sprintf("Fail server! %s", err.Error()),
		}
		ctx.JSON(http.StatusServiceUnavailable, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, addUserResponse)
}
