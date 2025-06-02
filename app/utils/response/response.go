package response

import (
	"douyin-backend/app/global/consts"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ReturnJson(Context *gin.Context, httpCode int, dataCode int, msg string, data interface{}) {
	Context.JSON(httpCode, gin.H{
		"code": dataCode,
		"msg":  msg,
		"data": data,
	})
}

func ErrorSystem(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, data)
}
