package user

import (
	"douyin-backend/app/global/consts"
	"douyin-backend/app/http/controller/web"
	"douyin-backend/app/http/validator/core/data_transfer"
	"douyin-backend/app/utils/response"

	"github.com/gin-gonic/gin"
)

type SendEmailCode struct {
	Email
}

func (s SendEmailCode) CheckParams(ctx *gin.Context) {
	if err := ctx.ShouldBind(&s); err != nil {
		response.ValidatorError(ctx, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(s, consts.ValidatorPrefix, ctx)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(ctx, "send_email_code 表单验证器json化失败", "")
	} else {
		(&web.UserController{}).SendEmailCode(extraAddBindDataContext)
	}
}
