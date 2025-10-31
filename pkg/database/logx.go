package database

import (
	"github.com/zeromicro/go-zero/core/logx"
)

func MapToLogFields(data map[string]any) []logx.LogField {
	fields := make([]logx.LogField, 0, len(data))
	for k, v := range data {
		fields = append(fields, logx.Field(k, v))
	}
	return fields
}

func QPToLogFields(qp QueryParams) []logx.LogField {
	fileds := MapToLogFields(qp.Query)
	fileds = append(fileds, logx.Field(PreloadsKey, qp.Preloads))
	fileds = append(fileds, logx.Field(OrderByKey, qp.OrderBy))
	fileds = append(fileds, logx.Field(LimitKey, qp.Limit))
	fileds = append(fileds, logx.Field(OffsetKey, qp.Offset))
	fileds = append(fileds, logx.Field(IsCountKey, qp.IsCount))
	return fileds
}
