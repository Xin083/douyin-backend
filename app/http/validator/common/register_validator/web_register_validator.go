package register_validator

import (
	"douyin-backend/app/core/container"
	"douyin-backend/app/global/consts"
	"douyin-backend/app/http/validator/web/douyin/user"
)

func WebRegisterValidator() {
	containers := container.CreateContainersFactory()

	var key string
	// jwt
	{
		key = consts.ValidatorPrefix + "JsonInBlacklist"
		containers.Set(key, user.JsonInBlacklist{})
	}

}
