package datastruct

import (
	"time"
)

type Garage struct {
	GarageID int64 `json:"id"`
	CarID    int64 `json:"car_id"`
}

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	GarageID       int64     `json:"garage_id"`
	RegistrationAt time.Time `json:"registration_at"`
}

type UserCar struct {
	GarageId int64 `json:"id"`
	CarId    int64 `json:"car_id"`
}

type CarModel struct {
	ID         string `json:"id"`
	Brand      string `json:"brand"`
	Model      string `json:"model"`
	Equipment  string `json:"equipment"`
	EngineType string `json:"engine_type"`
}
