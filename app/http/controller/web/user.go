package web

import (
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
