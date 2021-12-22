package foodstore

import (
	"context"
	"fooddelivery/common"
	"fooddelivery/modules/food/foodmodel"
)

func (s *sqlStore) UpdateFood(ctx context.Context, id int, data *foodmodel.FoodUpdate) error {
	db := s.db

	if err := db.Where("id = ?", id).Updates(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}