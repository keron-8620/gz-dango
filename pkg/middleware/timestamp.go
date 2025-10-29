// pkg/middleware/timestamp.go
package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"gz-dango/pkg/errors"
)

var (
	ErrNoTimestamp = errors.New(
		http.StatusBadRequest,
		"no_timestamp",
		"请求标头中缺少时间戳",
		nil,
	)
	ErrInvalidTimestamp = errors.New(
		http.StatusBadRequest,
		"invalid_timestamp",
		"无效的时间戳",
		nil,
	)
	ErrTimestampExpired = errors.New(
		http.StatusBadRequest,
		"timestamp_expired",
		"时间戳已过期, 请检查客户端时间同步",
		nil,
	)
)

// TimestampMiddleware 创建一个时间戳校验中间件
func TimestampMiddleware(maxDiffSecs int64) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 检查是否是 API 请求
			if !strings.HasPrefix(r.URL.Path, "/api") {
				next(w, r)
				return
			}
			logger := logx.WithContext(r.Context())

			// 从请求头获取 X-Timestamp
			timestampStr := r.Header.Get("X-Timestamp")
			if timestampStr == "" {
				logger.Errorw("请求缺少 X-Timestamp 头")
				httpx.WriteJson(w, ErrNoTimestamp.Code, ErrNoTimestamp.Reply())
				return
			}

			// 解析时间戳
			timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
			if err != nil {
				logger.Errorw(
					"请求时间戳解释失败",
					logx.Field("timestamp", timestampStr),
					logx.Field(errors.ErrKey, err.Error()),
				)
				httpx.WriteJson(w, ErrInvalidTimestamp.Code, ErrInvalidTimestamp.Reply())
				return
			}

			// 检查时间戳是否过期
			now := time.Now().Unix()
			diff := now - timestamp
			if diff < 0 {
				diff = -diff
			}
			if diff > maxDiffSecs { // 默认 300 秒 = 5 分钟
				logger.Errorw("X-Timestamp expired",
					logx.Field("current", now),
					logx.Field("received", timestamp))
				httpx.WriteJson(w, ErrTimestampExpired.Code, ErrTimestampExpired.Reply())
				return
			}

			next(w, r)
		}
	}
}
