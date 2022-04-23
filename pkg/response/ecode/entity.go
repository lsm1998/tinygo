package ecode

import "fmt"

//	状态码规则	来源(一位) 二级业务码（两位） 三级业务码（四位） 如：5 01 0100
//	     来源   1：默认 2：客户端异常 4：服务端提示	5：服务端异常
//	   二级码 	00~98：业务码		99：未知错误
//	   三级吗 	0000~9998：业务码 	9999：未知错误

// Errno 基础定义错误码
type Errno struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *Errno) Error() string {
	return fmt.Sprintf("Ext - code: %d, msg: %s", e.Code, e.Msg)
}

var (
	Success = &Errno{Code: 1000000, Msg: "成功"}
	Error   = &Errno{Code: 2000000, Msg: "失败"}

	ClientParamLack       = &Errno{Code: 2000100, Msg: "参数缺失"}
	ClientParamTypeFail   = &Errno{Code: 2000101, Msg: "参数类型错误"}
	ClientParamValueIlleg = &Errno{Code: 2000102, Msg: "参数值非法"}
	ClientFileExtIlleg    = &Errno{Code: 2000200, Msg: "非法文件类型"}
	ClientRequestRefuse   = &Errno{Code: 2010100, Msg: "请求拒绝"}
	ClientRequestRate     = &Errno{Code: 2010200, Msg: "请求过频"}
	ClientUnknownError    = &Errno{Code: 2999999, Msg: "未知客户端异常"}

	TipsDbNull              = &Errno{Code: 4000100, Msg: "DB查询结果为空"}
	TipsRedisNull           = &Errno{Code: 4000200, Msg: "redis查询结果为空"}
	TipsEsNull              = &Errno{Code: 4000300, Msg: "es查询结果为空"}
	TipsMongodbNull         = &Errno{Code: 4000400, Msg: "mongodb查询结果为空"}
	TipsUserNoLogin         = &Errno{Code: 4010100, Msg: "用户未登录"}
	TipsUserIsNoWpsVip      = &Errno{Code: 4010101, Msg: "用户非wps会员"}
	TipsUserIsNoDocerVip    = &Errno{Code: 4010102, Msg: "用户非稻壳会员"}
	TipsUserIsNoSuperVip    = &Errno{Code: 4010103, Msg: "用户非超级会员"}
	TipsUserIsNoPrivilege   = &Errno{Code: 4010104, Msg: "用户没有特权包"}
	TipsUserIsNoInFreeLimit = &Errno{Code: 4010105, Msg: "资源限免过期"}
	TipsUserNoPermission    = &Errno{Code: 4010106, Msg: "用户无权限"}
	TipsUnknownError        = &Errno{Code: 4999999, Msg: "未知服务提示"}

	ErrorDbConnectFail        = &Errno{Code: 5000100, Msg: "数据库链接失败"}
	ErrorDbSqlFail            = &Errno{Code: 5000101, Msg: "数据库SQL执行语句失败"}
	ErrorDbBreakFail          = &Errno{Code: 5000102, Msg: "数据库异常退出"}
	ErrorRedisConnectFail     = &Errno{Code: 5000200, Msg: "redis链接失败"}
	ErrorRedisCheckFail       = &Errno{Code: 5000201, Msg: "redis执行语句失败"}
	ErrorMongodbConnectFail   = &Errno{Code: 5000300, Msg: "mongodb链接失败"}
	ErrorMongodbCheckFail     = &Errno{Code: 5000301, Msg: "mongodb执行语句失败"}
	ErrorEsConnectFail        = &Errno{Code: 5000400, Msg: "es链接失败"}
	ErrorEsCheckFail          = &Errno{Code: 5000401, Msg: "es执行语句失败"}
	ErrorOutApiConnectFail    = &Errno{Code: 5010100, Msg: "三方接口链接失败"}
	ErrorOutApiTimeoutFail    = &Errno{Code: 5010101, Msg: "三方接口响应超时"}
	ErrorInsideApiConnectFail = &Errno{Code: 5010200, Msg: "内部接口链接失败"}
	ErrorInsideApiTimeoutFail = &Errno{Code: 5010201, Msg: "内部接口响应超时"}
	ErrorUnknownError         = &Errno{Code: 5999999, Msg: "未知服务异常"}
)
