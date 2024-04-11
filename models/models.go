package models

import "time"

type Ads struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Creation    time.Time `json:"creation"`
}
