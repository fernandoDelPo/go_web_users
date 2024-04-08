package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/fernandoDelPo/go-web-users-response/response"
	"github.com/fernandoDelPo/go_web_users/internal/user"
	"github.com/fernandoDelPo/go_web_users/pkg/transport"
	"github.com/gin-gonic/gin"
)

func NewUserHTPPServer(endpoints user.Endpoints) http.Handler {

	r := gin.Default()

	r.POST("/users", transport.GinServer(
		transport.Endpoint(endpoints.Create),
		decodeCreateUser,
		encodeResponse,
		encondeError))

	r.GET("/users", transport.GinServer(
		transport.Endpoint(endpoints.GetAll),
		decodeGetAllUsers,
		encodeResponse,
		encondeError))

	r.GET("/users/:id", transport.GinServer(
		transport.Endpoint(endpoints.Get),
		decodeGetUser,
		encodeResponse,
		encondeError))

	r.PATCH("/users/:id", transport.GinServer(
		transport.Endpoint(endpoints.Update),
		decoUpdateUser,
		encodeResponse,
		encondeError))

	return  r

}

func decodeGetUser(c *gin.Context) (interface{}, error) {

	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return user.GetReq{
		ID: id,
	}, nil
}

func decodeGetAllUsers(c *gin.Context) (interface{}, error) {
	
	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}
	return nil, nil
}

func decodeCreateUser(c *gin.Context) (interface{}, error) {
	c.Request.Header.Get("Authorization")

	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	var req user.CreateRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("error decoding request body:  '%v'", err.Error()))
	}
	return req, nil
}

func tokenVerify(token string) error {
	if os.Getenv("TOKEN") != token {
		return errors.New("invalid token")
	}
	return nil
}

func encodeResponse(c *gin.Context, resp interface{}) {
	r := resp.(response.Response)
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(r.StatusCode(), resp)

}

func encondeError(c *gin.Context, err error) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	c.JSON(resp.StatusCode(), resp)
}

func decoUpdateUser(c *gin.Context) (interface{}, error) {
	var req user.UpdateRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("error decoding request body:  '%v'", err.Error()))
	}

	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)

	if err != nil {
		return nil, response.BadRequest(err.Error())
	}
	req.ID = id

	return req, nil

}
