package middlewares

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/syyongx/php2go"
	"go-cms/common"
	"go-cms/pkg/e"
	"go-cms/pkg/util"
	"strings"
	"time"
)


//map[string]interface{}{"code": 400, "msg": "no user exists!", "data": nil}
type Response struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	TimeStamp int64       `json:"timestamp"`
}

func OutResponse(code int, data interface{}, msg string) Response {
	Resp := Response{
		Code:      code,
		Msg:       msg,
		Data:      data,
		TimeStamp: time.Now().Unix(), //time.Now().Format("2006-01-02 15:04:05")
	}
	return Resp
}

var supportMethod = [6]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
//配置不需要登录的url
var urlMapping = []string{"common:page_not_found","user:login","captcha:check","wechat:connect"}
// 支持伪造restful风格的http请求
// _method = "DELETE" 即将http的POST请求改为DELETE请求
func RestfulHandler() func(ctx *context.Context) {
	var restfulHandler = func(ctx *context.Context) {
		// 获取隐藏请求
		requestMethod := ctx.Input.Query("_method")
		
		if requestMethod ==  ""{
			// 正常请求
			requestMethod = ctx.Input.Method()
			logs.Debug("requestMethod:",requestMethod)
		}
		
		// 判断当前请求是否在允许请求内
		flag := false
		for _, method := range supportMethod{
			if method == requestMethod {
				flag = true
				break
			}
		}
		
		// 方法请求
		if flag == false {
			ctx.Output.Header("Content-Type", "application/json")
			resBody, err := json.Marshal(OutResponse(e.ERROR, nil, "Method Not Allow"))
			ctx.Output.Body(resBody)
			if err != nil {
				panic(err)
			}
			return
		}
		
		// 伪造请求方式
		if requestMethod != "" && ctx.Input.IsPost() {
			ctx.Request.Method = requestMethod
		}

		current_url := ctx.Request.URL.RequestURI()
		controllerName, actionName := getControllerAndAction(current_url)
		is_pass := php2go.InArray(php2go.Strtolower(controllerName+":"+actionName),urlMapping)
		if is_pass == false {
			token := ctx.Input.Header(beego.AppConfig.String("jwt::token_name"))
			allow, message,code := util.CheckToken(token)
			if(allow == false){
				ctx.Output.Header("Content-Type", "application/json")
				resBody, err := json.Marshal(OutResponse(code, nil, message))
				ctx.Output.Body(resBody)
				if err != nil {
					panic(err)
				}

				//_, ok := ctx.Input.Session("uid").(string)
				//ok2 := strings.Contains(ctx.Request.RequestURI, "/login")
				//if !ok && !ok2{
				//	ctx.Redirect(302, "/login/index")
				//}
			}else{
				common.UserId = code
			}
		}
	}
	return restfulHandler
}

func getControllerAndAction(url string)  (controllerName, actionName string){
	newStr := strings.ReplaceAll(strings.TrimLeft(url,"/api"),"/","|")

	tmp :=strings.Split(newStr, "|")
	var tow = ""
	if len(tmp) >= 2 {
		tow = tmp[1]
	}
	return tmp[0],tow
}