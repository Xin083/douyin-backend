package web

import (
	"douyin-backend/app/global/consts"
	"douyin-backend/app/model/user"
	"douyin-backend/app/utils/auth"
	"douyin-backend/app/utils/response"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (u *UserController) Login(ctx *gin.Context) {
	//var phone = ctx.GetString(consts.ValidatorPrefix + "phone")
	//var password = ctx.GetString(consts.ValidatorPrefix + "password")
	//ctx.ClientIP()

}

func (u *UserController) JsonInBlacklist(ctx *gin.Context) {

}

func (u *UserController) GetPanel(ctx *gin.Context) {
	var uid = auth.GetUidFromToken(ctx)
	userinfo, ok := user.CreateUserFactory("").GetPanel(uid)
	if ok {
		response.Success(ctx, consts.CurdStatusOkMsg, userinfo)
	} else {
		response.Fail(ctx, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}
