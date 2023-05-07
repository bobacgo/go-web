package app

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	admin_v1 "github.com/gogoclouds/go-web/api/admin/v1"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/r"
	"io"
)

type Resp struct {
	Code r.StatusCode `json:"code"`
	Msg  string       `json:"msg"`
}

type RespData[T any] struct {
	Resp
	Data T `json:"data"`
}

type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// LoggerErrMiddleware 日志记录中间件
// 1.只对 application/json
func LoggerErrMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.ContentType() != binding.MIMEJSON { // 只输出 application/json
			c.Next()
			return
		}
		reqBody, _ := c.GetRawData()
		if len(reqBody) > 0 { // 请求包体写回。
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}
		blw := &ResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		rspBody := blw.body.Bytes()
		var rsp RespData[any]
		if err := json.Unmarshal(rspBody, &rsp); err != nil {
			logger.Error(err)
			return
		}
		if rsp.Code != r.Ok {
			logger.Errorf("\nrequest: %s\nresponse: %s", reqBody, rspBody)
		}
	}
}

func loadRouter(e *gin.Engine) {
	e.Use(LoggerErrMiddleware())
	g := e.Group("v1")

	// sys menu
	menuApi_v1 := new(admin_v1.MenuApi)

	menu := g.Group("menu")
	menu.POST("", menuApi_v1.Create)
	menu.PUT("", menuApi_v1.Update)
	menu.DELETE("", menuApi_v1.Delete)
	menu.POST("tree", menuApi_v1.Tree)

	// sys role
	roleApi_v1 := new(admin_v1.RoleApi)

	role := g.Group("role")
	role.POST("details", roleApi_v1.Details)
	role.POST("", roleApi_v1.Create)
	role.PUT("", roleApi_v1.Update)
	role.DELETE("", roleApi_v1.Delete)
	role.POST("pageList", roleApi_v1.PageList)

	// sys user
	userApi_v1 := new(admin_v1.UserApi)

	// 获取用户列表、获取用户详情、创建用户、更新用户、更新状态、更新密码、删除用户
	user := g.Group("user")
	user.POST("details", userApi_v1.Details)
	user.POST("", userApi_v1.Create)
	user.PUT("", userApi_v1.Update)
	user.DELETE("", userApi_v1.Delete)
	user.PUT("updateStatus", userApi_v1.UpdateStatus)
	user.PUT("updatePassword", userApi_v1.UpdatePassword)
	user.POST("pageList", userApi_v1.PageList)

	// system
	systemApi_v1 := new(admin_v1.SystemApi)

	base := g.Group("base")
	base.POST("login", systemApi_v1.Login)
	base.GET("logout", systemApi_v1.Logout)
	base.GET("captcha", systemApi_v1.Captcha)
	base.POST("upload", systemApi_v1.Upload)

	// sys dictionary
}
