package dao

import (
	"airport-match-system/models"
	"context"
	"encoding/json"

	"airport-match-system/initial"

	log "airport-match-system/log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func CreateDriverOrderIfNotExists(order *models.DriverPublishOrder) (bool, error) {
	result := initial.MysqlDB.Where(models.DriverPublishOrder{DriverOrderID: order.DriverOrderID}).FirstOrCreate(&order)
	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected > 0 {
		return true, nil
	}

	return false, nil
}

func CreatePagIfNotExists(order *models.PassengerPublishOrder) (bool, error) {
	result := initial.MysqlDB.Where(models.PassengerPublishOrder{PassengerOrderID: order.PassengerOrderID}).FirstOrCreate(&order)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected > 0 {
		return true, nil
	}

	return false, nil
}

func FindMatch(passengerOrder *models.PassengerPublishOrder) (*models.DriverPublishOrder, error) {
	var driverOrder models.DriverPublishOrder
	err := initial.MysqlDB.Where("airport = ? AND available_time = ? AND vehicle_type = ? AND price_per_km <= ? AND status = 0",
		passengerOrder.Airport,
		passengerOrder.AvailableTime,
		passengerOrder.VehicleType,
		passengerOrder.MaxPricePerKm).Order("price_per_km ASC, rating DESC, created_at ASC").First(&driverOrder).Error
	if err != nil {
		return nil, err
	}

	return &driverOrder, nil
}

func CreateMatchRelation(passengerOrderId string, driverOrderId string) error {
	var matchRelation models.MatchRelation
	result := initial.MysqlDB.Where(models.MatchRelation{DriverOrderID: driverOrderId, PassengerOrderID: passengerOrderId}).FirstOrCreate(&matchRelation)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func CheckMatchRelation(passengerOrderId string, driverOrderId string) (bool, error) {
	var matchRelation models.MatchRelation
	err := initial.MysqlDB.Where(models.MatchRelation{DriverOrderID: driverOrderId, PassengerOrderID: passengerOrderId, Status: 0}).Order("updated_at DESC").First(&matchRelation).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func ExecuteMatch(passengerOrder *models.PassengerPublishOrder, driverOrder *models.DriverPublishOrder, feePerKm float64) (*models.MatchOrderResult, error) {
	passengerOrderID := passengerOrder.PassengerOrderID
	driverOrderID := driverOrder.DriverOrderID
	matchResult := models.MatchOrderResult{
		MatchOrderID:        uuid.New().String(),
		PassengerID:         passengerOrderID,
		DriverID:            driverOrderID,
		DriverRating:        float64(driverOrder.Rating),
		VehicleType:         driverOrder.VehicleType,
		DriverPricePerKm:    driverOrder.PricePerKm,
		PassengerPricePerKm: driverOrder.PricePerKm + feePerKm,
		FeePerKm:            feePerKm,
		Airport:             passengerOrder.Airport,
		AvailableTime:       passengerOrder.AvailableTime,
	}

	tx := initial.MysqlDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Model(&models.PassengerPublishOrder{}).Where("passenger_order_id = ?", passengerOrderID).Update("status", 1)
	if err.Error != nil {
		tx.Rollback()
		return nil, err.Error
	}

	err = tx.Model(&models.DriverPublishOrder{}).Where("driver_order_id = ?", driverOrderID).Update("status", 1)
	if err.Error != nil {
		tx.Rollback()
		return nil, err.Error
	}

	err = tx.Model(&models.MatchRelation{}).Where("driver_order_id = ? AND passenger_order_id = ?", driverOrderID, passengerOrderID).Update("status", 1)
	if err.Error != nil {
		tx.Rollback()
		return nil, err.Error
	}

	createErr := tx.Create(&matchResult).Error
	if createErr != nil {
		tx.Rollback()
		return nil, createErr
	}

	commitErr := tx.Commit().Error
	if commitErr != nil {
		tx.Rollback()
		return nil, commitErr
	}

	return &matchResult, nil
}

func GetFeeConfig() (*models.FeeConfig, error) {
	var feeConfig models.FeeConfig
	err := initial.MysqlDB.Where("status = ?", 1).Order("fee_per_km ASC").First(&feeConfig).Error
	if err != nil {
		return nil, err
	}

	return &feeConfig, nil
}

func FindDriverOrder(orderId string) (*models.DriverPublishOrder, error) {
	var order models.DriverPublishOrder
	err := initial.MysqlDB.Where(models.DriverPublishOrder{DriverOrderID: orderId}).First(&order).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func FindPassengerOrder(orderId string) (*models.PassengerPublishOrder, error) {
	var order models.PassengerPublishOrder
	err := initial.MysqlDB.Where(models.PassengerPublishOrder{PassengerOrderID: orderId}).First(&order).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func CheckPassengerOrderId(orderId string) (*models.PassengerPublishOrder, error) {
	result, err := initial.RDB.Get(context.Background(), orderId).Result()

	if err != nil && err != redis.Nil {
		return nil, err
	}

	var order models.PassengerPublishOrder
	if result != "" {
		err := json.Unmarshal([]byte(result), &order)
		if err != nil {
			return nil, err
		}
		log.LogInfo("[passenger redis cache]", order)
		return &order, nil
	}

	err = initial.MysqlDB.Where(models.PassengerPublishOrder{PassengerOrderID: orderId}).First(&order).Error
	if err != nil {
		return nil, err
	}
	log.LogInfo("[passenger mysql get]", order)
	return &order, nil
}

func CheckDriverOrderId(orderId string) (*models.DriverPublishOrder, error) {
	result, err := initial.RDB.Get(context.Background(), orderId).Result()

	if err != nil && err != redis.Nil {
		return nil, err
	}

	var order models.DriverPublishOrder
	if result != "" {
		err := json.Unmarshal([]byte(result), &order)
		if err != nil {
			return nil, err
		}
		log.LogInfo("[driver redis cache]", order)
		return &order, nil
	}

	err = initial.MysqlDB.Where(models.DriverPublishOrder{DriverOrderID: orderId}).First(&order).Error
	if err != nil {
		return nil, err
	}
	log.LogInfo("[driver mysql get]", order)
	return &order, nil
}

func DeleteCacheOrder(passengerOrderId string, driverOrderId string) (bool, error) {
	err := initial.RDB.Del(context.Background(), passengerOrderId, driverOrderId).Err()
	if err != nil {
		log.LogError("[delete order cache]", err)
		return false, err
	}
	return true, nil
}
