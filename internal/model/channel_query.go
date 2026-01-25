package model

import (
	"uni-search-hub/pkg/database"
)

type providerCountRow struct {
	Provider string
	Cnt      int64
}

// CountEnabledChannelsByProvider 统计各 provider 启用渠道数（enabled=true）。
func CountEnabledChannelsByProvider(providers []string) (map[string]int64, error) {
	out := map[string]int64{}
	if len(providers) == 0 {
		return out, nil
	}

	var rows []providerCountRow
	err := database.DB.Model(&Channel{}).
		Select("provider, count(*) as cnt").
		Where("enabled = ?", true).
		Where("provider IN ?", providers).
		Group("provider").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.Provider] = r.Cnt
	}
	return out, nil
}
