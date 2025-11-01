package service

import (
	"airport-match-system/dao"
	"airport-match-system/initial"
	"airport-match-system/models"
	"airport-match-system/utils"
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	log "airport-match-system/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PassengerCreateRoute(c *gin.Context) {
	userName := c.PostForm("user_name")
	airport := c.PostForm("airport")
	availableTimeStr := c.PostForm("available_time")
	vehicleType := c.PostForm("vehicle_type")
	MaxPricePerKmStr := c.PostForm("max_price_per_km")

	if userName == "" || airport == "" || availableTimeStr == "" || vehicleType == "" || MaxPricePerKmStr == "" {
		utils.ErrorResponse(c, utils.ExistBlankParameter.Code, utils.ExistBlankParameter.Message)
		return
	}

	availableTime, err := time.ParseInLocation("2006-01-02 15:04:05", availableTimeStr, time.Local)
	if err != nil {
		utils.ErrorResponse(c, utils.InvalidFormat.Code, utils.InvalidFormat.Message)
		return
	}

	now := time.Now()
	log.LogInfo("[availableTime]", availableTime)
	if availableTime.Before(now.Add(1 * time.Hour)) {
		utils.ErrorResponse(c, utils.InvalidAvailableTime.Code, utils.InvalidAvailableTime.Message)
		return
	}

	vehicleTypeInt, err := strconv.Atoi(vehicleType)
	if err != nil {
		utils.ErrorResponse(c, utils.InvalidVehicleType.Code, utils.InvalidVehicleType.Message)
		return
	}

	maxPricePerKm, err := strconv.ParseFloat(MaxPricePerKmStr, 64)
	if err != nil || maxPricePerKm <= 0 {
		utils.ErrorResponse(c, utils.InvalidPerKm.Code, utils.InvalidPerKm.Message)
		return
	}

	orderId := "passenger-" + uuid.New().String()

	passengerOrder := models.PassengerPublishOrder{
		PassengerOrderID: orderId,
		PassengerName:    userName,
		VehicleType:      int8(vehicleTypeInt),
		MaxPricePerKm:    maxPricePerKm,
		Airport:          airport,
		AvailableTime:    availableTime,
		Status:           0,
	}

	res, err := dao.CreatePagIfNotExists(&passengerOrder)
	if err != nil {
		utils.ErrorResponse(c, utils.CreateOrderFail.Code, utils.CreateOrderFail.Message)
		return
	}
	if !res {
		utils.ErrorResponse(c, utils.DuplicatedOrder.Code, utils.DuplicatedOrder.Message)
		return
	}

	rdbRes, err := json.Marshal(&passengerOrder)
	if err != nil {
		log.LogError(err)
		utils.ErrorResponse(c, utils.MarshalOrderFail.Code, utils.MarshalOrderFail.Message)
		return
	}

	err = initial.RDB.Set(context.Background(), orderId, rdbRes, 72*time.Hour).Err()
	if err != nil {
		log.LogError(err)
	}

	data := map[string]interface{}{
		"passenger_order_id": orderId,
		"status":             0, // 0: not matched, 1: matched
	}

	utils.SuccessResponse(c, 200, "success", data)
}

func DriverCreateRoute(c *gin.Context) {
	userName := c.PostForm("user_name")
	airport := c.PostForm("airport")
	vehicleType := c.PostForm("vehicle_type")
	pricePerKmStr := c.PostForm("price_per_km")
	availableTimeStr := c.PostForm("available_time")
	if userName == "" || airport == "" || vehicleType == "" || pricePerKmStr == "" || availableTimeStr == "" {
		utils.ErrorResponse(c, utils.ExistBlankParameter.Code, utils.ExistBlankParameter.Message)
		return
	}

	availableTime, err := time.ParseInLocation("2006-01-02 15:04:05", availableTimeStr, time.Local)
	if err != nil {
		utils.ErrorResponse(c, utils.InvalidFormat.Code, utils.InvalidFormat.Message)
		return
	}

	now := time.Now()
	log.LogInfo("[availableTime]", availableTime)
	if availableTime.Before(now.Add(1 * time.Hour)) {
		utils.ErrorResponse(c, utils.InvalidAvailableTime.Code, utils.InvalidAvailableTime.Message)
		return
	}

	vehicleTypeInt, err := strconv.Atoi(vehicleType)
	if err != nil {
		utils.ErrorResponse(c, utils.InvalidVehicleType.Code, utils.InvalidVehicleType.Message)
		return
	}

	rating := rand.Intn(5) + 1
	pricePerKm, err := strconv.ParseFloat(pricePerKmStr, 64)
	if err != nil || pricePerKm <= 0 {
		utils.ErrorResponse(c, utils.InvalidPerKm.Code, utils.InvalidPerKm.Message)
		return
	}

	orderId := "driver-" + uuid.New().String()

	driverOrder := models.DriverPublishOrder{
		DriverOrderID: orderId,
		DriverName:    userName,
		VehicleType:   int8(vehicleTypeInt),
		Rating:        int8(rating),
		PricePerKm:    pricePerKm,
		Airport:       airport,
		AvailableTime: availableTime,
		Status:        0,
	}

	res, err := dao.CreateDriverOrderIfNotExists(&driverOrder)
	if err != nil {
		utils.ErrorResponse(c, utils.CreateOrderFail.Code, utils.CreateOrderFail.Message)
		return
	}
	if !res {
		utils.ErrorResponse(c, utils.DuplicatedOrder.Code, utils.DuplicatedOrder.Message)
		return
	}

	rdbRes, err := json.Marshal(&driverOrder)
	if err != nil {
		log.LogError(err)
		utils.ErrorResponse(c, utils.MarshalOrderFail.Code, utils.MarshalOrderFail.Message)
		return
	}

	err = initial.RDB.Set(context.Background(), orderId, rdbRes, 72*time.Hour).Err()
	if err != nil {
		log.LogError(err)
	}

	data := map[string]interface{}{
		"driver_order_id": orderId,
		"status":          0, // 0: not matched, 1: matched
	}

	utils.SuccessResponse(c, 200, "success", data)
}

func FindMatchOrder(c *gin.Context) {
	passengerOrderId := c.PostForm("passenger_order_id")
	passengerOrder, err := dao.CheckPassengerOrderId(passengerOrderId)
	if err != nil {
		utils.ErrorResponse(c, utils.CheckOrderFail.Code, utils.CheckOrderFail.Message)
		return
	}
	if passengerOrder.Status != 0 {
		utils.ErrorResponse(c, utils.OrderHasMatched.Code, utils.OrderHasMatched.Message)
		return
	}

	log.LogInfo("[passengerOrder]", passengerOrder)

	FeeConfig, err := dao.GetFeeConfig()
	if err != nil {
		utils.ErrorResponse(c, utils.GetFeeFail.Code, utils.GetFeeFail.Message)
		return
	}
	feePerKm := FeeConfig.FeePerKm

	passengerPublicOrder := models.PassengerPublishOrder{
		Airport:       passengerOrder.Airport,
		AvailableTime: passengerOrder.AvailableTime,
		VehicleType:   passengerOrder.VehicleType,
		MaxPricePerKm: passengerOrder.MaxPricePerKm - feePerKm,
	}
	matchedDriverOrder, err := dao.FindMatch(&passengerPublicOrder)
	if err != nil {
		utils.ErrorResponse(c, utils.MatchOrderFail.Code, utils.MatchOrderFail.Message)
		return
	}
	log.LogInfo("[matchedDriverOrder]", matchedDriverOrder)

	driverSidePrice := matchedDriverOrder.PricePerKm
	passengerSideShowPrice := driverSidePrice + feePerKm

	showData := map[string]interface{}{
		"driver_order_id": matchedDriverOrder.DriverOrderID,
		"driver_name":     matchedDriverOrder.DriverName,
		"vehicle_type":    matchedDriverOrder.VehicleType,
		"rating":          matchedDriverOrder.Rating,
		"price_per_km":    passengerSideShowPrice,
		"airport":         matchedDriverOrder.Airport,
		"available_time":  matchedDriverOrder.AvailableTime,
	}

	err = dao.CreateMatchRelation(passengerOrderId, matchedDriverOrder.DriverOrderID)
	if err != nil {
		utils.ErrorResponse(c, utils.CreateMatchRelationFail.Code, utils.CreateMatchRelationFail.Message)
		return
	}

	utils.SuccessResponse(c, 200, "success", showData)
}

func ExecuteMatchOrder(c *gin.Context) {
	passengerOrderId := c.PostForm("passenger_order_id")
	driverOrderId := c.PostForm("driver_order_id")

	_, err := dao.CheckMatchRelation(passengerOrderId, driverOrderId)
	if err != nil {
		utils.ErrorResponse(c, utils.CheckRelationFail.Code, utils.CheckRelationFail.Message)
		return
	}

	passengerOrder, err := dao.CheckPassengerOrderId(passengerOrderId)
	if err != nil {
		utils.ErrorResponse(c, utils.InvalidPassengerOrder.Code, utils.InvalidPassengerOrder.Message)
		return
	}
	if passengerOrder.Status != 0 {
		utils.ErrorResponse(c, utils.OrderHasMatched.Code, utils.OrderHasMatched.Message)
		return
	}

	driverOrder, err := dao.CheckDriverOrderId(driverOrderId)
	if err != nil {
		utils.ErrorResponse(c, utils.InvalidDriverOrder.Code, utils.InvalidDriverOrder.Message)
		return
	}
	if driverOrder.Status != 0 {
		utils.ErrorResponse(c, utils.OrderHasMatched.Code, utils.OrderHasMatched.Message)
		return
	}

	FeeConfig, err := dao.GetFeeConfig()
	if err != nil {
		utils.ErrorResponse(c, utils.GetFeeFail.Code, utils.GetFeeFail.Message)
		return
	}
	feePerKm := FeeConfig.FeePerKm

	matchResult, err := dao.ExecuteMatch(passengerOrder, driverOrder, feePerKm)
	if err != nil {
		utils.ErrorResponse(c, utils.ExecuteMatchFail.Code, utils.ExecuteMatchFail.Message)
		return
	}

	err = initial.RDB.Del(context.Background(), passengerOrderId, driverOrderId).Err()
	if err != nil {
		log.LogError("[delete order cache]", err)
	}

	utils.SuccessResponse(c, 200, "success", matchResult)
}
