package svc

import (
	"go-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

func mapToLogFields(data map[string]any) []logx.LogField {
	fields := make([]logx.LogField, 0, len(data))
	for k, v := range data {
		fields = append(fields, logx.Field(k, v))
	}
	return fields
}

func qpToLogFields(qp database.QueryParams) []logx.LogField {
	fileds := mapToLogFields(qp.Query)
	fileds = append(fileds, logx.Field(database.PreloadsKey, qp.Preloads))
	fileds = append(fileds, logx.Field(database.OrderByKey, qp.OrderBy))
	fileds = append(fileds, logx.Field(database.LimitKey, qp.Limit))
	fileds = append(fileds, logx.Field(database.OffsetKey, qp.Offset))
	fileds = append(fileds, logx.Field(database.IsCountKey, qp.IsCount))
	return fileds
}
