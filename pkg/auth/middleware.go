package auth

import (
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// extractToken 从不同位置提取 token
func extractToken(r *http.Request) string {
	// 检查是否为 WebSocket 升级请求
	if r.Header.Get("Connection") == "upgrade" ||
		r.Header.Get("Upgrade") == "websocket" {
		// WebSocket 请求优先从查询参数获取，其次从头部获取
		if token := r.URL.Query().Get("Authorization"); token != "" {
			return token
		}
		if token := r.Header.Get("Sec-WebSocket-Protocol"); token != "" {
			return token
		}
		return ""
	}

	// HTTP 请求从 Authorization 头部获取
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		return authHeader
	}

	// 也可以从查询参数获取token作为备选方案
	return r.URL.Query().Get("Authorization")
}

func AuthMiddleware(enforcer *AuthEnforcer) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 检查是否是 API 请求
			if !strings.HasPrefix(r.URL.Path, "/api") {
				next.ServeHTTP(w, r)
				return
			}
			// 检查是否是登陆请求
			if r.URL.Path == "/api/v1/login" && r.Method == http.MethodPost {
				next.ServeHTTP(w, r)
				return
			}
			// 从请求头获取token
			token := extractToken(r)
			if token == "" {
				httpx.WriteJson(w, ErrNoAuthor.Code, ErrNoAuthor.Reply())
				return
			}
			ctx := r.Context()
			// 身份认证
			info, err := enforcer.Authentication(ctx, token)
			if err != nil {
				httpx.WriteJson(w, err.Code, err.Reply())
				return
			}
			// 访问鉴权
			hasPerm, err := enforcer.Authorization(info.Role, r.URL.Path, r.Method)
			if err != nil {
				httpx.WriteJson(w, err.Code, err.Reply())
				return
			}
			if !hasPerm {
				httpx.WriteJson(w, ErrForbidden.Code, ErrForbidden.Reply())
				return
			}

			// 将用户信息存储到 context 中
			nctx := SetUserClaims(ctx, info)
			// 创建一个新的请求，使用更新后的 context
			newReq := r.WithContext(nctx)

			next.ServeHTTP(w, newReq)
		}
	}
}
