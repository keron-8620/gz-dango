package database

import (
	"strconv"
	"strings"
)

var (
	DefaultPage int = 1
	DefaultSize int = 10
)

func StringToListUint(pks string) []uint {
	pks = strings.TrimSpace(pks)
	if pks == "" {
		return make([]uint, 0)
	}
	pkList := strings.Split(pks, ",")
	var ids []uint
	for _, pk := range pkList {
		pk = strings.TrimSpace(pk)
		if pk == "" {
			continue
		}
		value, err := strconv.ParseUint(pk, 10, 32)
		if err != nil {
			continue
		}
		ids = append(ids, uint(value))
	}
	return ids
}

func CountPages(count int64, size int64) int64 {
	if count == 0 || size == 0 {
		return 0
	}
	return (count + size - 1) / size
}
