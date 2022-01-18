package restaurantratingmodel

import "foodlive/common"

const (
	EntityName = "RestaurantRating"
)

type RestaurantRating struct {
	common.SQLModel
	UserId       int     `json:"user_id" gorm:"user_id"`
	RestaurantId int     `json:"restaurant_id" gorm:"restaurant_id"`
	Point        float64 `json:"point" gorm:"point"`
	Comment      string  `json:"comment" gorm:"comment"`
	Status       bool    `json:"status" gorm:"status"`
}

func (RestaurantRating) TableName() string {
	return "restaurant_ratings"
}

type RestaurantRatingCreate struct {
	common.SQLModelCreate
	UserId       int     `json:"user_id" gorm:"user_id"`
	RestaurantId int     `json:"restaurant_id" gorm:"restaurant_id"`
	Point        float64 `json:"point" gorm:"point"`
	Comment      string  `json:"comment" gorm:"comment"`
	Status       bool    `json:"-" gorm:"status"`
}

func (RestaurantRatingCreate) TableName() string {
	return RestaurantRating{}.TableName()
}

func (data *RestaurantRatingCreate) Validate() error {
	return nil
}

type RestaurantRatingUpdate struct {
	common.SQLModelUpdate
	Point   float64 `json:"point" gorm:"point"`
	Comment string  `json:"comment" gorm:"comment"`
	Status  bool    `json:"-" gorm:"status"`
}

func (RestaurantRatingUpdate) TableName() string {
	return RestaurantRating{}.TableName()
}

func (data *RestaurantRatingUpdate) Validate() error {
	return nil
}
