package main

import (
	"foodlive/component"
	"foodlive/middleware"
	"foodlive/modules/authsso/fbsso/fbssotransport/ginfbsso"
	"foodlive/modules/authsso/googlesso/googlessotransport/gingooglesso"
	"foodlive/modules/cart/carttransport/gincart"
	"foodlive/modules/category/categorytransport/gincategory"
	"foodlive/modules/food/foodtransport/ginfood"
	"foodlive/modules/order/ordertransport/ginorder"
	"foodlive/modules/restaurant/restauranttransport/ginrestaurant"
	"foodlive/modules/restaurantlike/restaurantliketransport/ginrestaurantlike"
	"foodlive/modules/restaurantowner/restaurantownertransport/ginrestaurantowner"
	"foodlive/modules/restaurantrating/restaurantratingtransport/ginrestaurantrating"
	"foodlive/modules/upload/uploadtransport/ginupload"
	"foodlive/modules/user/usertransport/ginuser"
	"foodlive/modules/useraddress/useraddresstransport/ginuseraddress"
	"foodlive/modules/userdevicetoken/userdevicetokentransport/ginuserdevicetoken"
	"github.com/gin-gonic/gin"
)

func setupRouter(r *gin.Engine, appCtx component.AppContext) {
	r.Use(middleware.Recover(appCtx))
	v1Route(r, appCtx)
}

func v1Route(r *gin.Engine, appCtx component.AppContext) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", ginuser.UserReigster(appCtx))
		v1.POST("/activate", ginuser.UserActiveAccount(appCtx))
		v1.POST("/resend-otp-active", ginuser.ResendOTPActive(appCtx))
		v1.POST("/forgot-password", ginuser.UserForgotPassword(appCtx))
		v1.POST("/reset-password", ginuser.UserResetPassword(appCtx))
		v1.POST("/login", ginuser.UserLogin(appCtx))

		v1.POST("/upload", ginupload.Upload(appCtx))

		sso := v1.Group("/sso")
		{
			sso.POST("/register-google", gingooglesso.UserGoogleLogin(appCtx))
			sso.POST("/login-google", gingooglesso.UserGoogleLogin(appCtx))

			sso.POST("/register-facebook", ginfbsso.UserFacebookRegister(appCtx))
			sso.POST("/login-facebook", ginfbsso.UserFacebookRegister(appCtx))
		}

		restaurant := v1.Group("/restaurant", middleware.RequireAuth(appCtx))
		{
			restaurant.GET("", ginrestaurant.ListRestaurant(appCtx))
			restaurant.GET("/:id", ginrestaurant.FindRestaurant(appCtx))

			//get food if restaurant
			restaurant.GET("/:id/food", ginfood.UserListFoodOfRestaurant(appCtx))

			//Like restaurant
			restaurant.POST("/:id/like", ginrestaurantlike.UserLikeRestaurant(appCtx))
			restaurant.DELETE("/:id/unlike", ginrestaurantlike.UserUnLikeRestaurant(appCtx))

			//List user like restaurant
			restaurant.GET("/:id/like", ginrestaurantlike.ListUserLikeRestaurant(appCtx))

			//User rating restaurant
			restaurant.POST("/:id/rating", ginrestaurantrating.CreateRestaurantRating(appCtx))
			restaurant.GET("/:id/rating", ginrestaurantrating.ListRestaurantRating(appCtx))
			restaurant.PUT("/rating/:id_rating", ginrestaurantrating.UpdateRestaurantRating(appCtx))
		}

		category := v1.Group("/category", middleware.RequireAuth(appCtx))
		{
			category.GET("", gincategory.UserListCategory(appCtx))
		}

		userAddress := v1.Group("address", middleware.RequireAuth(appCtx))
		{
			userAddress.POST("", ginuseraddress.CreateUserAddress(appCtx))
			userAddress.PUT("/:id", ginuseraddress.UpdateUserAddress(appCtx))
			userAddress.DELETE("/:id", ginuseraddress.DeleteUserAddress(appCtx))
			userAddress.GET("", ginuseraddress.ListUserAddress(appCtx))
		}

		cart := v1.Group("/cart", middleware.RequireAuth(appCtx))
		{
			cart.POST("", gincart.CreateItemInCart(appCtx))
			cart.PUT("", gincart.UpdateItemInCart(appCtx))
			cart.DELETE("/:food_id", gincart.DeleteAItemInCart(appCtx))
			cart.DELETE("", gincart.DeleteAllItemInCart(appCtx))
			cart.GET("", gincart.ListItemInCart(appCtx))
		}

		userDeviceToken := v1.Group("/user-device-token", middleware.RequireAuth(appCtx))
		{
			userDeviceToken.GET("/my-device", ginuserdevicetoken.FindUserDeviceToken(appCtx))
			userDeviceToken.POST("", ginuserdevicetoken.CreateUserDeviceToken(appCtx))
		}

		order := v1.Group("/order", middleware.Recover(appCtx))
		{
			order.POST("", ginorder.CreateOrder(appCtx))
			order.GET("/:order_id", ginorder.FindOrder(appCtx))
			order.GET("", ginorder.ListOrder(appCtx))

		}

		//========================================================================================================

		admin := v1.Group("/admin")
		{
			admin.POST("/register-owner-restaurant", ginrestaurantowner.OwnerRestaurantRegister(appCtx))

			adminRestaurant := admin.Group("/restaurant")
			{
				adminRestaurant.DELETE("/:id", ginrestaurant.DeleteRestaurant(appCtx))
				adminRestaurant.PUT("/:id/status", ginrestaurant.UpdateRestaurantStatus(appCtx))
			}

			adminCategory := admin.Group("/category")
			{
				adminCategory.POST("", gincategory.CreateCategory(appCtx))
				adminCategory.PUT("/:id", gincategory.UpdateCategory(appCtx))
				adminCategory.DELETE("/:id", gincategory.DeleteCategory(appCtx))
				adminCategory.GET("", gincategory.AdminListCategory(appCtx))
			}

			adminUserDeviceToken := admin.Group("/user-device-token")
			{
				adminUserDeviceToken.GET("", ginuserdevicetoken.ListUserDeviceToken(appCtx))
			}
		}

		ownerRestaurant := v1.Group("/owner-restaurant")
		{
			ownerRestaurant.POST("/login", ginrestaurantowner.OwnerRestaurantLogin(appCtx))
			ownerRestaurant.POST("/active", ginrestaurantowner.OwnerRestaurantActive(appCtx))
			ownerRestaurant.POST("/send-otp", ginrestaurantowner.SendOTPActiveOwnerRestaurant(appCtx))

			ownerRestaurant1 := ownerRestaurant.Group("/restaurant", middleware.RequireAuthOwnerRestaurant(appCtx))
			{
				ownerRestaurant1.POST("", ginrestaurant.CreateRestaurant(appCtx))
				ownerRestaurant1.PUT("/:id", ginrestaurant.UpdateRestaurant(appCtx))
				ownerRestaurant1.GET("", ginrestaurant.ListRestaurantOwner(appCtx))
				ownerRestaurant1.GET("/:id/food", ginfood.ListFoodOfRestaurant(appCtx)) // Get food of restaurant
			}

			ownerFood := ownerRestaurant.Group("/food", middleware.RequireAuthOwnerRestaurant(appCtx))
			{
				ownerFood.POST("", ginfood.CreateFood(appCtx))
				ownerFood.PUT("/:id", ginfood.UpdateRestaurant(appCtx))
				ownerFood.DELETE("/:id", ginfood.DeleteRestaurant(appCtx))
				ownerFood.GET("", ginfood.ListFoodOfRestaurant(appCtx))
			}

			ownerCategory := ownerRestaurant.Group("/category")
			{
				ownerCategory.GET("", gincategory.UserListCategory(appCtx))
			}
		}

	}
}
