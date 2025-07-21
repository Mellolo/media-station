package enum

var PermissionLevels = []string{
	PermissionPublic,
	PermissionLogin,
	PermissionVIP,
	PermissionPrivate,
	PermissionForbidden,
}

const (
	PermissionPublic    = "public "   // 公开：所有用户可访问
	PermissionLogin     = "login"     // 登录用户：仅登录用户可访问
	PermissionVIP       = "vip"       // VIP 用户：需为 VIP 身份
	PermissionPrivate   = "private"   // 私密：仅上传者或指定用户可访问
	PermissionForbidden = "forbidden" // 禁止：任何用户无法访问
)
