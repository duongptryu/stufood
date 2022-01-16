package ginuser

import (
	"foodlive/common"
	"foodlive/component"
	"foodlive/modules/user/userbiz"
	"foodlive/modules/user/usermodel"
	"foodlive/modules/user/userstorage"
	"github.com/gin-gonic/gin"
)

func ResendOTPActive(appCtx component.AppContext) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var data usermodel.ResendOTP

		if err := c.ShouldBindJSON(&data); err != nil {
			panic(common.ErrParseJson(err))
		}

		store := userstorage.NewSQLStore(appCtx.GetDatabase())

		biz := userbiz.NewResendOTPActiveBiz(store, appCtx.GetMyCache(), appCtx.GetMySms())
		err := biz.ResendOTPActiveAccount(c.Request.Context(), &data)
		if err != nil {
			panic(err)
		}

		c.JSON(200, common.NewSimpleSuccessResponse(true))
	}
}
