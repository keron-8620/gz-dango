package middleware

import (
	goerrors "errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"


	"gz-dango/pkg/errors"
)

// ErrorMiddleware 异常处理中间件
// 拦截所有panic和错误，进行统一处理和响应
func ErrorMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic信息和堆栈跟踪
				stack := debug.Stack()

				// 构造错误响应
				var errMsg string
				switch v := err.(type) {
				case error:
					errMsg = v.Error()
				case string:
					errMsg = v
				default:
					errMsg = fmt.Sprintf("%v", v)
				}

				logx.WithContext(r.Context()).Errorw("panic recovered",
					logx.Field("error", errMsg),
					logx.Field("method", r.Method),
					logx.Field("url", r.URL.Path),
					logx.Field("stack", string(stack)),
				)
				err := errors.FromError(goerrors.New(errMsg))
				httpx.WriteJson(w, err.Code, err.Reply())
			}
		}()

		// 继续处理请求
		next.ServeHTTP(w, r)
	}
}
