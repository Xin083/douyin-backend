package gin_release

import (
	variable "douyin-backend/app/global"
	"douyin-backend/app/global/consts"
	"douyin-backend/app/utils/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
)

func ReleaseRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	engine := gin.New()
	// 载入gin的中间件，关键是第二个中间件，我们对它进行了自定义重写，将可能的 panic 异常等，统一使用 zaplog 接管，保证全局日志打印统一

	engine.Use(gin.Logger(), CustomRecovery())
	return engine
}

func CustomRecovery() gin.HandlerFunc {
	DefaultErrorWriter := &PanicExceptionRecord{}
	return gin.RecoveryWithWriter(DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		// 这里针对发生的panic等异常进行统一响应即可
		// 这里的 err 数据类型为 ：runtime.boundsError  ，需要转为普通数据类型才可以输出
		response.ErrorSystem(c, "", fmt.Sprintf("%s", err))
	})
}

// PanicExceptionRecord  panic等异常记录
type PanicExceptionRecord struct{}

func (p *PanicExceptionRecord) Write(b []byte) (n int, err error) {
	errStr := string(b)
	err = errors.New(errStr)
	variable.ZapLog.Error(consts.ServerOccurredErrorMsg, zap.String("errStrace", errStr))
	return len(errStr), err
}
