package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ExistBlankParameter     = NewError(10001, "exist blank parameter")
	InvalidFormat           = NewError(10002, "invalid format")
	InvalidAvailableTime    = NewError(10003, "Invalid available time")
	InvalidVehicleType      = NewError(10004, "Invalid vehicle type")
	InvalidPerKm            = NewError(10005, "Invalid per km parameter")
	CreateOrderFail         = NewError(10006, "Failed to create order")
	DuplicatedOrder         = NewError(10007, "Duplicated order")
	MarshalOrderFail        = NewError(10008, "Failed to marshal order")
	CheckOrderFail          = NewError(10009, "Failed to check order")
	GetFeeFail              = NewError(10010, "Failed to get fee")
	MatchOrderFail          = NewError(10011, "Failed to match order")
	CreateMatchRelationFail = NewError(10012, "Failed to create match relation")
	OrderHasMatched         = NewError(10013, "Order has matched")
	CheckRelationFail       = NewError(10014, "Failed to check relation")
	InvalidPassengerOrder   = NewError(10015, "Invalid passenger order")
	InvalidDriverOrder      = NewError(10016, "Invalid driver order")
	ExecuteMatchFail        = NewError(10017, "Failed to execute match orders")
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"Message"`
}

func NewError(code int, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

type ResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(http.StatusBadRequest, ResponseData{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func SuccessResponse(c *gin.Context, code int, message string, data any) {
	if message == "" {
		message = "success"
	}
	c.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
