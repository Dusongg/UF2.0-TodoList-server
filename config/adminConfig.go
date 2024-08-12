package config

// 管理员用户：修改其他用户的任务
var Admin = map[string]bool{
	"dusong": true,
}

func IsAdmin(admin string) bool {
	_, ok := Admin[admin]
	return ok
}
