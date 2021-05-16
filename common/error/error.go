package error

//自定义错误
type Errno int
var (
	ERROR_OK             Errno = 200
	ERROR_USER_NOTEXISTS Errno = 400
	ERROR_USER_EXISTS    Errno = 401
	ERROR_USER_PWD       Errno = 402
	ERROR_USER_LOGINED   Errno = 403
	ERROR_SEVER_ERR      Errno = 404
	ERROR_UNKNOW         Errno = 405
)

var ErrorNoMsgMap map[Errno]string

func init(){
	ErrorNoMsgMap = make(map[Errno]string, 0)
	ErrorNoMsgMap[ERROR_OK]					= "成功"
	ErrorNoMsgMap[ERROR_USER_NOTEXISTS] 	= "用户不存在"
	ErrorNoMsgMap[ERROR_USER_EXISTS] 		= "用户已存在"
	ErrorNoMsgMap[ERROR_USER_PWD] 			= "用户名或密码错误"
	ErrorNoMsgMap[ERROR_USER_LOGINED] 		= "用户已登录"
	ErrorNoMsgMap[ERROR_SEVER_ERR] 			= "内部服务错误"
	ErrorNoMsgMap[ERROR_UNKNOW] 			= "未知错误"
}