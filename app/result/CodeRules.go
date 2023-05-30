package result

// 错误码规则:
// (1) 错误码需为 > 0 的数;
//
// (2) 错误码为 5 位数:
//              ----------------------------------------------------------
//                  第1位               2、3位                  4、5位
//              ----------------------------------------------------------
//                服务级错误码          模块级错误码	         具体错误码
//

var (
	Ok  = result(200, "success") // 通用成功
	Err = result(500, "error")   // 通用错误

	// 服务级错误码
	ErrBind       = result(10001, "请求参数错误")
	ErrValidation = result(10002, "参数验证失败")
	ErrDatabase   = result(10003, "数据库错误")

	NoToken      = result(10004, "请求未携带token，无权限访问")
	TokenExpired = result(10005, "token已过期")
	TokenInvalid = result(10006, "token不合法")
	ErrToken     = result(10007, "token生成失败")

	ErrUser       = result(10010, "用户不存在")
	ErrPassword   = result(10011, "密码错误")
	ErrUserExist  = result(10012, "用户已存在")
	ErrUserAdd    = result(10013, "用户添加失败")
	ErrUserFind   = result(10014, "用户查找失败")
	ErrUserUpdate = result(10015, "用户更新失败")
	ErrUserDelete = result(10016, "用户删除失败")

	ErrPostAdd    = result(10021, "帖子添加失败")
	ErrPostFind   = result(10022, "帖子查找失败")
	ErrPostUpdate = result(10023, "帖子更新失败")
	ErrPostDelete = result(10024, "帖子删除失败")

	ErrCommentAdd    = result(10031, "评论添加失败")
	ErrCommentFind   = result(10032, "评论查找失败")
	ErrCommentUpdate = result(10033, "评论更新失败")
	ErrCommentDelete = result(10034, "评论删除失败")

	ErrLikeAdd    = result(10041, "点赞添加失败")
	ErrLikeFind   = result(10042, "点赞查找失败")
	ErrLikeUpdate = result(10043, "点赞更新失败")
)

var MsgFlags = map[int]string{
	200: "success", // 通用成功
	500: "error",   // 通用错误

	// 服务级错误码
	10001: "请求参数错误",
	10002: "参数验证失败",
	10003: "数据库错误",

	10004: "请求未携带token，无权限访问",
	10005: "token已过期",
	10006: "token不合法",
	10007: "token生成失败",

	10010: "用户不存在",
	10011: "密码错误",
	10012: "用户已存在",
	10013: "用户添加失败",
	10014: "用户查找失败",
	10015: "用户更新失败",
	10016: "用户删除失败",

	10021: "帖子添加失败",
	10022: "帖子查找失败",
	10023: "帖子更新失败",
	10024: "帖子删除失败",

	10031: "评论添加失败",
	10032: "评论查找失败",
	10033: "评论更新失败",
	10034: "评论删除失败",

	10041: "点赞添加失败",
	10042: "点赞查找失败",
	10043: "点赞更新失败",
}

// GetMsg 获取错误码对应的错误信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code] //MsgFlags是一个map[int]string
	if ok {
		return msg
	}
	return MsgFlags[Err.Code]
}
