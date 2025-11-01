package models

import "time"

type DriverPublishOrder struct {
	ID            uint      `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	DriverOrderID string    `gorm:"column:driver_order_id;NOT NULL"`
	DriverName    string    `gorm:"column:driver_name;NOT NULL"`
	VehicleType   int8      `gorm:"column:vehicle_type;default:1;NOT NULL"` // Vehicle type 1 - Economy 2 - Comfort 3 - Business 4 - Luxury
	Rating        int8      `gorm:"column:rating;default:1;NOT NULL"`       // Driver level 1-5, with the lowest level being 1 and the highest level being 5 (the driver's level is randomly generated)
	PricePerKm    float64   `gorm:"column:price_per_km;NOT NULL"`
	Airport       string    `gorm:"column:airport;NOT NULL"`
	AvailableTime time.Time `gorm:"column:available_time;NOT NULL"`
	Status        int8      `gorm:"column:status;default:0;NOT NULL"` // Order status: 0-Unmatched 1-Matched
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

type PassengerPublishOrder struct {
	ID               uint      `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	PassengerOrderID string    `gorm:"column:passenger_order_id;NOT NULL"`
	PassengerName    string    `gorm:"column:passenger_name;NOT NULL"`
	VehicleType      int8      `gorm:"column:vehicle_type;default:1;NOT NULL"`
	MaxPricePerKm    float64   `gorm:"column:max_price_per_km;NOT NULL"`
	Airport          string    `gorm:"column:airport;NOT NULL"`
	AvailableTime    time.Time `gorm:"column:available_time;NOT NULL"`
	Status           int8      `gorm:"column:status;default:0;NOT NULL"` // Order status: 0-Unmatched 1-Matched
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

type MatchRelation struct {
	ID               uint      `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	DriverOrderID    string    `gorm:"column:driver_order_id;NOT NULL"`
	PassengerOrderID string    `gorm:"column:passenger_order_id;NOT NULL"`
	Status           int8      `gorm:"column:status;default:0;NOT NULL"` // 0-Initialize matching 1-Matching completed
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

type MatchOrderResult struct {
	ID                  uint      `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	MatchOrderID        string    `gorm:"column:match_order_id;NOT NULL"`
	PassengerID         string    `gorm:"column:passenger_id;NOT NULL"`
	DriverID            string    `gorm:"column:driver_id;NOT NULL"`
	DriverRating        float64   `gorm:"column:rating;default:1;NOT NULL"`
	VehicleType         int8      `gorm:"column:vehicle_type;default:1;NOT NULL"`
	DriverPricePerKm    float64   `gorm:"column:driver_price_per_km;NOT NULL"`
	PassengerPricePerKm float64   `gorm:"column:passenger_price_per_km;NOT NULL"`
	FeePerKm            float64   `gorm:"column:fee_per_km;NOT NULL"`
	Airport             string    `gorm:"column:air_port;NOT NULL"`
	AvailableTime       time.Time `gorm:"column:available_time;NOT NULL"`
	CreatedAt           time.Time `gorm:"column:created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at"`
}

type FeeConfig struct {
	ID        uint      `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	FeePerKm  float64   `gorm:"column:fee_per_km;NOT NULL"`
	Status    int8      `gorm:"column:status;default:1;NOT NULL"` // 0-Disabled 1-Enabled
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
