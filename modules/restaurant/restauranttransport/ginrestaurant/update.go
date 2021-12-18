package ginrestaurant

import (
	"fooddelivery/common"
	"fooddelivery/component"
	"fooddelivery/modules/restaurant/restaurantbiz"
	"fooddelivery/modules/restaurant/restaurantmodel"
	"fooddelivery/modules/restaurant/restaurantstore"
	"github.com/gin-gonic/gin"
	"strconv"
)

func UpdateRestaurant(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			panic(common.ErrParseJson(err))
		}

		var data restaurantmodel.RestaurantUpdate

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrParseJson(err))
		}

		store := restaurantstore.NewSqlStore(appCtx.GetDatabase())
		biz := restaurantbiz.NewUpdateRestaurantBiz(store)

		if err := biz.UpdateRestaurantBiz(c.Request.Context(), id, &data); err != nil {
			panic(err)
		}

		c.JSON(200, common.NewSimpleSuccessResponse(true))
	}
}