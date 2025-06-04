package web

import (
	"douyin-backend/app/global/consts"
	"douyin-backend/app/global/variable"
	"douyin-backend/app/model/user"
	"douyin-backend/app/model/video"
	userstoken "douyin-backend/app/service/users/token"
	"douyin-backend/app/utils/auth"
	"douyin-backend/app/utils/email_utils"
	"douyin-backend/app/utils/md5_encrypt"
	"douyin-backend/app/utils/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type UserController struct {
}

func (u *UserController) Register(ctx *gin.Context) {
	var phone = ctx.GetString(consts.ValidatorPrefix + "phone")
	var password = ctx.GetString(consts.ValidatorPrefix + "password")
	var userIp = ctx.ClientIP()
	if user.CreateUserFactory("").Register(phone, md5_encrypt.Base64Md5(password), userIp) {
		response.Success(ctx, consts.CurdStatusOkMsg, consts.CurdRegisterOkMsg)
	} else {
		response.Fail(ctx, consts.CurdRegisterFailCode, consts.CurdRegisterFailMsg, "")
	}
}

func (u *UserController) Login(ctx *gin.Context) {
	var phone = ctx.GetString(consts.ValidatorPrefix + "phone")
	var password = ctx.GetString(consts.ValidatorPrefix + "password")
	userModel, ok := user.CreateUserFactory("").Login(phone, password)
	if ok {
		userTokenFactory := userstoken.CreateUserFactory()
		if userToken, err := userTokenFactory.GenerateToken(userModel.UID, userModel.Nickname, userModel.Phone, variable.ConfigYml.GetInt64("Token.JwtTokenCreatedExpireAt")); err == nil {
			if userTokenFactory.RecordLoginToken(userToken, ctx.ClientIP()) {
				response.Success(ctx, consts.CurdStatusOkMsg, gin.H{
					"isExist": true,
					"uid":     userModel.UID,
					"token":   userToken,
				})
			}

		} else {
			variable.ZapLog.Error("生成token出错!")
		}
	} else {
		response.Fail(ctx, consts.CurdLoginFailCode, consts.CurdLoginFailMsg, gin.H{
			"isExist": false,
			"uid":     userModel.UID,
			"token":   "",
		})
	}
}

func (u *UserController) UpdateInfo(ctx *gin.Context) {
	uid := auth.GetUidFromToken(ctx)
	var operationType = ctx.GetFloat64(consts.ValidatorPrefix + "operation_type")
	var data = ctx.GetString(consts.ValidatorPrefix + "data")
	updateState := user.CreateUserFactory("").UpdateInfo(uid, int(operationType), data)
	if updateState {
		response.Success(ctx, consts.CurdStatusOkMsg, gin.H{
			"data": updateState,
			"msg":  "修改成功!",
		})
	} else {
		response.Fail(ctx, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, gin.H{
			"data": updateState,
			"msg":  "修改失败!",
		})
	}
}

func (u *UserController) Attention(ctx *gin.Context) {
	var uid = auth.GetUidFromToken(ctx)
	var followingId = ctx.GetString(consts.ValidatorPrefix + "following_id")
	var action = ctx.GetBool(consts.ValidatorPrefix + "action")
	var followingIdInt64, _ = strconv.ParseInt(followingId, 10, 64)
	actionStatus := user.CreateUserFactory("").Attention(uid, followingIdInt64, action)
	if actionStatus {
		if action {
			response.Success(ctx, consts.CurdStatusOkMsg, gin.H{
				"data": actionStatus,
				"msg":  "关注成功",
			})
		} else {
			response.Success(ctx, consts.CurdStatusOkMsg, gin.H{
				"data": actionStatus,
				"msg":  "取消关注成功",
			})
		}
	} else {
		if action {
			response.Fail(ctx, consts.CurdInsertFailCode, consts.CurdInsertFailMsg, gin.H{
				"data": actionStatus,
				"msg":  "关注失败",
			})
		} else {
			response.Fail(ctx, consts.CurdDeleteFailCode, consts.CurdDeleteFailMsg, gin.H{
				"data": actionStatus,
				"msg":  "取消关注失败",
			})
		}
	}
}

func (u *UserController) AwemeStatus(ctx *gin.Context) {
	var uid = auth.GetUidFromToken(ctx)
	awemeStatus, success := user.CreateUserFactory("").AwemeStatus(uid)
	if success {
		response.Success(ctx, consts.CurdStatusOkMsg, awemeStatus)
	} else {
		response.Fail(ctx, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, awemeStatus)
	}
}

func (u *UserController) JsonInBlacklist(ctx *gin.Context) {
	// TODO
}

func (u *UserController) GetUserInfo(context *gin.Context) {
	// TODO 具体业务逻辑实现
	response.Success(context, consts.CurdStatusOkMsg, "GetUserInfo-ok")
	//var id = context.GetFloat64(consts.ValidatorPrefix + "id")
	//video := sv_home.CreateShortVideoFactory("").GetVideoById(int(id))
	//if video.Id != 0 {
	//	response.Success(context, consts.CurdStatusOkMsg, video)
	//} else {
	//	response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	//}
}

func (u *UserController) GetUserVideoList(ctx *gin.Context) {
	uid, _ := strconv.Atoi(ctx.Query("uid"))
	videoList, ok := video.CreateVideoFactory("").GetUserVideoList(int64(uid))
	if ok {
		response.Success(ctx, consts.CurdStatusOkMsg, videoList)
	} else {
		response.Fail(ctx, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
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

func (u *UserController) GetFriends(ctx *gin.Context) {
	var uid = auth.GetUidFromToken(ctx)
	friends, ok := user.CreateUserFactory("").GetFriends(uid)
	if ok {
		response.Success(ctx, consts.CurdStatusOkMsg, friends)
	} else {
		response.Fail(ctx, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

func (u *UserController) GetFollow(ctx *gin.Context) {
	var uid = auth.GetUidFromToken(ctx)
	follow, ok := user.CreateUserFactory("").GetFollow(uid)
	if ok {
		response.Success(ctx, consts.CurdStatusOkMsg, follow)
	} else {
		response.Fail(ctx, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

func (u *UserController) GetFans(ctx *gin.Context) {
	var uid = auth.GetUidFromToken(ctx)
	fans, ok := user.CreateUserFactory("").GetFans(uid)
	if ok {
		response.Success(ctx, consts.CurdStatusOkMsg, fans)
	} else {
		response.Fail(ctx, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

func (u *UserController) GetMyVideo(ctx *gin.Context) {
	var uid = auth.GetUidFromToken(ctx)
	var PageNo = ctx.GetFloat64(consts.ValidatorPrefix + "pageNo")
	var PageSize = ctx.GetFloat64(consts.ValidatorPrefix + "pageSize")
	list, total, ok := video.CreateVideoFactory("").GetMyVideo(uid, int64(PageNo), int64(PageSize))
	if ok {
		response.Success(ctx, consts.CurdStatusOkMsg, gin.H{
			"pageNo": PageNo,
			"total":  total,
			"list":   list,
		})
	} else {
		response.Fail(ctx, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

func (u *UserController) GetMyPrivateVideo(ctx *gin.Context) {
	var uid = auth.GetUidFromToken(ctx)
	var PageNo = ctx.GetFloat64(consts.ValidatorPrefix + "pageNo")
	var PageSize = ctx.GetFloat64(consts.ValidatorPrefix + "pageSize")
	list, total, ok := video.CreateVideoFactory("").GetMyPrivateVideo(uid, int64(PageNo), int64(PageSize))
	if ok {
		response.Success(ctx, consts.CurdStatusOkMsg, gin.H{
			"pageNo": PageNo,
			"total":  total,
			"list":   list,
		})
	} else {
		response.Fail(ctx, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

func (u *UserController) GetMyLikeVideo(ctx *gin.Context) {
	var uid = auth.GetUidFromToken(ctx)
	var PageNo = ctx.GetFloat64(consts.ValidatorPrefix + "pageNo")
	var PageSize = ctx.GetFloat64(consts.ValidatorPrefix + "pageSize")
	list, total, ok := video.CreateVideoFactory("").GetMyLikeVideo(uid, int64(PageNo), int64(PageSize))
	if ok {
		response.Success(ctx, consts.CurdStatusOkMsg, gin.H{
			"pageNo": PageNo,
			"total":  total,
			"list":   list,
		})
	} else {
		response.Fail(ctx, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

func (u *UserController) GetMyCollectVideo(ctx *gin.Context) {
	var uid = auth.GetUidFromToken(ctx)
	var PageNo = ctx.GetFloat64(consts.ValidatorPrefix + "pageNo")
	var PageSize = ctx.GetFloat64(consts.ValidatorPrefix + "pageSize")
	list, total, ok := video.CreateVideoFactory("").GetMyCollectVideo(uid, int64(PageNo), int64(PageSize))

	if ok {
		response.Success(ctx, consts.CurdStatusOkMsg, gin.H{
			"video": gin.H{
				"pageNo": PageNo,
				"total":  total,
				"list":   list,
			},
			"music": gin.H{
				"list":  []interface{}{},
				"total": 0,
			},
		})
	} else {
		response.Fail(ctx, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

func (u *UserController) GetMyHistoryVideo(ctx *gin.Context) {
	var uid = auth.GetUidFromToken(ctx)
	var PageNo = ctx.GetFloat64(consts.ValidatorPrefix + "pageNo")
	var PageSize = ctx.GetFloat64(consts.ValidatorPrefix + "pageSize")

	list, total, ok := video.CreateVideoFactory("").GetMyHistoryVideo(uid, int64(PageNo), int64(PageSize))
	if ok {
		response.Success(ctx, consts.CurdStatusOkMsg, gin.H{
			"pageNo": PageNo,
			"total":  total,
			"list":   list,
		})
	} else {
		response.Fail(ctx, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

func (u *UserController) GetMyHistoryOther(ctx *gin.Context) {
	response.Success(ctx, consts.CurdStatusOkMsg, "GetMyHistoryOther-ok")
}

// SendEmailCode 发送邮箱验证码
func (u *UserController) SendEmailCode(ctx *gin.Context) {
	email := ctx.GetString(consts.ValidatorPrefix + "email")
	variable.ZapLog.Info("收到发送验证码请求", zap.String("email", email))

	// 生成并发送验证码
	code, err := user.CreateUserFactory("").SendEmailVerificationCode(email)
	if err != nil {
		variable.ZapLog.Error("生成验证码失败", zap.Error(err))
		response.Fail(ctx, consts.CurdLoginFailCode, "发送验证码失败", err.Error())
		return
	}
	variable.ZapLog.Info("验证码生成成功", zap.String("code", code))

	// 发送验证码邮件
	if err := email_utils.SendVerificationCode(email, code); err != nil {
		variable.ZapLog.Error("发送验证码邮件失败", zap.Error(err))
		response.Fail(ctx, consts.CurdLoginFailCode, "发送验证码邮件失败", err.Error())
		return
	}

	response.Success(ctx, consts.CurdStatusOkMsg, "验证码已发送")
}

// EmailLogin 邮箱验证码登录
func (u *UserController) EmailLogin(ctx *gin.Context) {
	email := ctx.GetString(consts.ValidatorPrefix + "email")
	code := ctx.GetString(consts.ValidatorPrefix + "code")

	// 验证验证码
	if !user.CreateUserFactory("").VerifyEmailCode(email, code) {
		response.Fail(ctx, consts.CurdLoginFailCode, "验证码错误或已过期", "")
		return
	}

	// 获取用户信息
	userModel, exists := user.CreateUserFactory("").GetUserByEmail(email)
	var uid int64
	var nickname string

	if !exists {
		// 如果用户不存在，进行注册
		var ok bool
		uid, ok = user.CreateUserFactory("").RegisterByEmail(email, ctx.ClientIP())
		if !ok {
			response.Fail(ctx, consts.CurdLoginFailCode, "注册失败", "")
			return
		}
		// 获取新注册用户的信息
		userModel, _ = user.CreateUserFactory("").GetUserByEmail(email)
		nickname = userModel.Nickname
	} else {
		uid = userModel.UID
		nickname = userModel.Nickname
	}

	// 生成token
	userTokenFactory := userstoken.CreateUserFactory()
	if userToken, err := userTokenFactory.GenerateToken(uid, nickname, email, variable.ConfigYml.GetInt64("Token.JwtTokenCreatedExpireAt")); err == nil {
		if userTokenFactory.RecordLoginToken(userToken, ctx.ClientIP()) {
			response.Success(ctx, consts.CurdStatusOkMsg, gin.H{
				"isExist": exists, // 返回是否是新注册用户
				"uid":     uid,
				"token":   userToken,
			})
			return
		}
	}

	variable.ZapLog.Error("生成token出错!")
	response.Fail(ctx, consts.CurdLoginFailCode, "登录失败", "")
}
