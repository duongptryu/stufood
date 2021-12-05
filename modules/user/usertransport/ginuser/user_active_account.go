package ginuser

import (
	"fooddelivery/common"
	"fooddelivery/component"
	"fooddelivery/modules/user/userbiz"
	"fooddelivery/modules/user/usermodel"
	"fooddelivery/modules/user/userstorage"
	"github.com/gin-gonic/gin"
)

func UserActiveAccount(appCtx component.AppContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var data usermodel.UserActive

		if err := c.ShouldBindJSON(&data); err != nil {
			panic(common.ErrParseJson(err))
		}

		userActiveStore := userstorage.NewSQLStore(appCtx.GetDatabase())
		userActiveBiz := userbiz.NewActiveUserBiz(userActiveStore, appCtx.GetMyCache())

		if err := userActiveBiz.UserActiveBiz(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		c.JSON(200, common.NewSimpleSuccessResponse(true))
	}
}