package websocket

import (
	"douyin-backend/app/global/consts"
	"douyin-backend/app/global/variable"
	controllerWs "douyin-backend/app/http/controller/websocket"
	"douyin-backend/app/http/validator/core/data_transfer"
	userstoken "douyin-backend/app/service/users/token"
	"douyin-backend/app/utils/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Connect struct {
	Token string `form:"token" json:"token" binding:"required,min=20"`
}

// 验证器语法，参见 Register.go文件，有详细说明
// 注意：websocket 连接建立之前如果有错误，只能在服务端同构日志输出方式记录（因为使用response.Fail等函数，客户端是收不到任何信息的）

func (c Connect) CheckParams(context *gin.Context) {
	// 1. 首先检查是否开启websocket服务配置（在配置项中开启）
	if variable.ConfigYml.GetInt("Websocket.Start") != 1 {
		variable.ZapLog.Error(consts.WsServerNotStartMsg)
		return
	}
	//2.基本的验证规则没有通过
	if err := context.ShouldBind(&c); err != nil {
		variable.ZapLog.Error("客户端上线参数不合格", zap.Error(err))
		return
	}
	token := c.Token
	if len(token) >= 20 {
		tokenIsEffective := userstoken.CreateUserFactory().IsEffective(token)
		if tokenIsEffective {
			if customeToken, err := userstoken.CreateUserFactory().ParseToken(token); err == nil {
				key := variable.ConfigYml.GetString("Token.BindContextKeyName")
				// token 验证通过并绑定在请求的上下文中
				context.Set(key, customeToken)
			} else {
				response.TokenParseFail(context, err)
			}
		} else {
			response.ErrorTokenAuthFail(context)
		}
	} else {
		response.ErrorTokenBaseInfo(context)
	}
	extraAddBindDataContext := data_transfer.DataAddContext(c, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		variable.ZapLog.Error("websocket-Connect 表单验证器json化失败")
		context.Abort()
		return
	} else {
		if serviceWs, ok := (&controllerWs.Ws{}).OnOpen(context); ok == false {
			variable.ZapLog.Error(consts.WsOpenFailMsg)
		} else {
			(&controllerWs.Ws{}).OnMessage(serviceWs, context) // 注意这里传递的service_ws必须是调用open返回的，必须保证的ws对象的一致性
		}
	}

}
