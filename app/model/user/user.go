package user

import (
	"context"
	"douyin-backend/app/global/variable"
	"douyin-backend/app/model"
	"douyin-backend/app/utils/md5_encrypt"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type UserModel struct {
	*gorm.DB                `gorm:"-" json:"-"`
	RedisClient             *redis.Client   `gorm:"-" json:"-"`
	UID                     int64           `json:"uid"`                       // bigint
	ShortID                 int             `json:"short_id"`                  // int
	UniqueID                string          `json:"unique_id"`                 // varchar(255)
	Gender                  string          `json:"gender"`                    // char(1)
	UserAge                 int             `json:"user_age"`                  // int
	Nickname                string          `json:"nickname"`                  // varchar(100)
	Country                 string          `json:"country"`                   // varchar(100)
	Province                string          `json:"province"`                  // varchar(100)
	District                string          `json:"district"`                  // varchar(255)
	City                    string          `json:"city"`                      // varchar(255)
	Signature               string          `json:"signature"`                 // text
	IPLocation              string          `json:"ip_location"`               // varchar(100)
	BirthdayHideLevel       int             `json:"birthday_hide_level"`       // int
	CanShowGroupCard        int             `json:"can_show_group_card"`       // int
	AwemeCount              int64           `json:"aweme_count"`               // bigint
	TotalFavorited          int64           `json:"total_favorited"`           // bigint
	FavoritingCount         int             `json:"favoriting_count"`          // int
	FollowerCount           int64           `json:"follower_count"`            // bigint
	FollowingCount          int64           `json:"following_count"`           // bigint
	ForwardCount            int             `json:"forward_count"`             // int
	PublicCollectsCount     int             `json:"public_collects_count"`     // int
	MPlatformFollowersCount int64           `json:"mplatform_followers_count"` // bigint
	MaxFollowerCount        int64           `json:"max_follower_count"`        // bigint
	FollowStatus            int             `json:"follow_status"`             // int
	FollowerStatus          int             `json:"follower_status"`           // int
	FollowerRequestStatus   int             `json:"follower_request_status"`   // int
	CoverColour             string          `json:"cover_colour"`              // varchar(100)
	CoverURL                json.RawMessage `json:"cover_url"`                 // json
	WhiteCoverURL           json.RawMessage `json:"white_cover_url"`           // json
	ShareInfo               json.RawMessage `json:"share_info"`                // json
	CommerceInfo            json.RawMessage `json:"commerce_info"`             // json
	CommerceUserInfo        json.RawMessage `json:"commerce_user_info"`        // json
	CommerceUserLevel       int             `json:"commerce_user_level"`       // int
	CardEntries             json.RawMessage `json:"card_entries"`              // json
	Avatar168x168           json.RawMessage `json:"avatar_168x168"`            // json
	Avatar300x300           json.RawMessage `json:"avatar_300x300"`            // json
}

func CreateUserFactory(sqlType string) *UserModel {
	return &UserModel{
		DB:          model.UseDbConn(sqlType),
		RedisClient: model.GetRedisClient(),
	}
}

func (u *UserModel) Register(phone, password, userIp string) bool {
	var createTime = time.Now().Unix()
	sql1 := `INSERT INTO tb_accounts(phone, password, last_login_ip, create_time) SELECT ?, ?, ?, ? FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM tb_accounts WHERE phone=?)`
	result := u.Exec(sql1, phone, password, userIp, createTime, phone)
	if result.Error != nil {
		variable.ZapLog.Error("Register SQL代码执行出错!", zap.Error(result.Error))
		return false
	}
	// 获取 UID
	var uid int64
	sql2 := `SELECT uid FROM tb_accounts WHERE phone = ?;`
	result = u.Raw(sql2, phone).Find(&uid)
	if result.Error != nil {
		variable.ZapLog.Error("Register-SELECT-uid SQL代码执行出错!", zap.Error(result.Error))
		return false
	}
	sql3 := `INSERT INTO tb_users(uid, card_entries) VALUES (?, ?);`
	result = u.Exec(sql3, uid, json.RawMessage{})
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (u *UserModel) Login(phone, password string) (account Account, ok bool) {
	sql := `
		SELECT ta.uid, ta.nickname, ta.phone, ta.password
		from tb_accounts as ta
		where phone=?
		limit 1;`
	result := u.Raw(sql, phone).Find(&account)
	if result.Error != nil {
		variable.ZapLog.Error("Login SQL代码执行出错!", zap.Error(result.Error))
		ok = false
		return
	}
	if account.Password == md5_encrypt.Base64Md5(password) {
		ok = true
	} else {
		ok = false
	}
	return
}

func (u *UserModel) Attention(uid, followingId int64, action bool) bool {
	currentTime := time.Now().Unix()
	attentionSql := `INSERT INTO tb_relations (follower_id, following_id, create_time) VALUES (?, ?, ?);`
	unattentionSql := `DELETE FROM tb_relations WHERE follower_id=? and following_id=?;`
	var result *gorm.DB
	if action {
		result = u.Exec(attentionSql, uid, followingId, currentTime)
	} else {
		result = u.Exec(unattentionSql, uid, followingId)
	}
	if result.Error != nil {
		variable.ZapLog.Error("Attention SQL执行出错!", zap.Error(result.Error))
	}
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (u *UserModel) UpdateInfo(uid int64, operationType int, data string) bool {
	// 目前支持三种修改类型(nickname/unique_id/signature)
	var sql string
	switch operationType {
	case 1:
		sql = `UPDATE tb_users SET nickname=? WHERE uid=?;`
	case 2:
		sql = `UPDATE tb_users SET unique_id=? WHERE uid=?;`
	case 3:
		sql = `UPDATE tb_users SET signature=? WHERE uid=?;`
	case 4:
		sql = `UPDATE tb_users SET gender=? WHERE uid=?;`
	case 5:
		sql = `UPDATE tb_users SET birthday=? WHERE uid=?;`

	}
	result := u.Exec(sql, data, uid)
	if result.Error == nil {
		return true
	} else {
		variable.ZapLog.Error("Update UserInfo failed!")
		return false
	}
}

func (u *UserModel) AwemeStatus(uid int64) (awemeStatus AwemeStatusModel, success bool) {
	attentionSql := `SELECT following_id FROM tb_relations WHERE follower_id=?`
	diggSql := `SELECT aweme_id FROM tb_diggs WHERE uid=?`
	collectSql := `SELECT aweme_id FROM tb_collects WHERE uid=?`
	u.Raw(attentionSql, uid).Find(&awemeStatus.Attentions)
	u.Raw(diggSql, uid).Find(&awemeStatus.Likes)
	u.Raw(collectSql, uid).Find(&awemeStatus.Collects)
	return awemeStatus, true
}

func (u *UserModel) GetPanel(uid int64) (userinfo model.User, ok bool) {
	sql := `
		SELECT *
		from tb_users as tu
		where uid=?
		limit 1;`
	result := u.Raw(sql, uid).Find(&userinfo)
	if result.Error != nil {
		variable.ZapLog.Error("GetPanel SQL执行出错!", zap.Error(result.Error))
		ok = false
		return
	}
	ok = true
	return
}

func (u *UserModel) GetFriends(uid int64) (userinfo []model.User, ok bool) {
	sql := `
		SELECT *
		from tb_users as tu
		where uid IN (
		    SELECT tr.following_id
		    FROM tb_relations as tr
		    WHERE  tr.follower_id=?) AND
		    uid IN (
		    SELECT tr.follower_id
		    FROM tb_relations as tr
		    WHERE  tr.following_id=?
		    );`
	result := u.Raw(sql, uid, uid).Find(&userinfo)
	if result.Error != nil {
		variable.ZapLog.Error("GetFriends SQL执行出错!", zap.Error(result.Error))
		ok = false
		return
	}
	ok = true
	return
}

func (u *UserModel) GetFollow(uid int64) (userinfo []model.User, ok bool) {
	sql := `
		SELECT *
		from tb_users as tu
		where uid IN (
		    SELECT tr.following_id
		    FROM tb_relations as tr
		    WHERE  tr.follower_id=?);`
	result := u.Raw(sql, uid).Find(&userinfo)
	if result.Error != nil {
		variable.ZapLog.Error("GetFollow SQL执行出错!", zap.Error(result.Error))
		ok = false
		return
	}
	ok = true
	return
}

func (u *UserModel) GetFans(uid int64) (userinfo []model.User, ok bool) {
	sql := `
		SELECT *
		from tb_users as tu
		where uid IN (
		    SELECT tr.follower_id
		    FROM tb_relations as tr
		    WHERE  tr.following_id=?);`
	result := u.Raw(sql, uid).Find(&userinfo)
	if result.Error != nil {
		variable.ZapLog.Error("GetFans SQL执行出错!", zap.Error(result.Error))
		ok = false
		return
	}
	ok = true
	return
}

func (u *UserModel) OauthCheckTokenIsOk(uid int64, token string) bool {
	sql := `SELECT token 
			FROM tb_auth_access_tokens 
			WHERE uid=? AND revoked=0 AND expires_at>UNIX_TIMESTAMP(NOW()) ORDER BY expires_at DESC, updated_at DESC LIMIT ?;`
	maxOnlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
	rows, err := u.Raw(sql, uid, maxOnlineUsers).Rows()
	defer func() {
		// 释放记录集
		_ = rows.Close()
	}()
	if err == nil && rows != nil {
		for rows.Next() {
			var tempToken string
			err := rows.Scan(&tempToken)
			fmt.Println(err)
			if err == nil {
				if tempToken == token {
					return true
				}
			}
		}
	}
	return false
}

func (u *UserModel) OauthLoginToken(uid int64, token string, expiresAt int64, clientIp string) bool {
	sql := `INSERT INTO tb_auth_access_tokens(uid, action_name, token, created_at, expires_at, client_ip)
			SELECT ?, 'login', ?, ?, ?, ? FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM tb_auth_access_tokens a WHERE a.uid=? AND a.action_name='login' AND a.token=?)
			`
	createdAt := time.Now().Unix()
	if u.Exec(sql, uid, token, createdAt, expiresAt, clientIp, uid, token).Error == nil {
		// 异步缓存用户有效的token到redis
		if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
			go u.ValidTokenCacheToRedis(uid)
		}
		return true
	}
	return false
}

// SendEmailVerificationCode 发送邮箱验证码
func (u *UserModel) SendEmailVerificationCode(email string) (string, error) {
	// 生成6位随机验证码
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 将验证码存储到Redis，设置过期时间
	redisKey := fmt.Sprintf("email_verify_code:%s", email)
	expireTime := variable.ConfigYml.GetInt("Email.VerificationCodeExpire")

	ctx := context.Background()
	if err := u.RedisClient.Set(ctx, redisKey, code, time.Duration(expireTime)*time.Second).Err(); err != nil {
		return "", err
	}

	return code, nil
}

// RegisterByEmail 通过邮箱注册
func (u *UserModel) RegisterByEmail(email, userIp string) (int64, bool) {
	// 将当前时间格式化为 "YYYY-MM-DD HH:MM:SS" 格式的字符串
	var createTime = time.Now().Format("2006-01-02 15:04:05")
	// 生成随机昵称
	nickname := fmt.Sprintf("用户%d", rand.Intn(1000000))

	sql1 := `INSERT INTO tb_accounts(email, nickname, last_login_ip, create_time)
			 SELECT ?, ?, ?, ? FROM DUAL
			 WHERE NOT EXISTS(SELECT 1 FROM tb_accounts WHERE email=?)`
	// 注意这里使用了格式化后的时间字符串 createTime
	result := u.Exec(sql1, email, nickname, userIp, createTime, email)
	if result.Error != nil {
		variable.ZapLog.Error("RegisterByEmail SQL代码执行出错!", zap.Error(result.Error))
		return 0, false
	}

	// 获取 UID
	var uid int64
	sql2 := `SELECT uid FROM tb_accounts WHERE email = ?;`
	result = u.Raw(sql2, email).Find(&uid)
	if result.Error != nil {
		variable.ZapLog.Error("RegisterByEmail-SELECT-uid SQL代码执行出错!", zap.Error(result.Error))
		return 0, false
	}

	// 创建用户信息
	sql3 := `INSERT INTO tb_users(uid, card_entries) VALUES (?, ?);`
	result = u.Exec(sql3, uid, json.RawMessage{})
	if result.RowsAffected > 0 {
		return uid, true
	}
	return 0, false
}

// VerifyEmailCode 验证邮箱验证码
func (u *UserModel) VerifyEmailCode(email, code string) bool {
	redisKey := fmt.Sprintf("email_verify_code:%s", email)
	ctx := context.Background()

	storedCode, err := u.RedisClient.Get(ctx, redisKey).Result()
	if err != nil {
		return false
	}

	if storedCode == code {
		// 验证成功后删除验证码
		u.RedisClient.Del(ctx, redisKey)
		return true
	}
	return false
}

// GetUserByEmail 通过邮箱获取用户信息
func (u *UserModel) GetUserByEmail(email string) (Account, bool) {
	var account Account
	sql := `
		SELECT ta.uid, ta.nickname, ta.phone, ta.password
		from tb_accounts as ta
		where email=?
		limit 1;`
	result := u.Raw(sql, email).Find(&account)
	if result.Error != nil {
		variable.ZapLog.Error("GetUserByEmail SQL代码执行出错!", zap.Error(result.Error))
		return account, false
	}
	return account, account.UID != 0
}
