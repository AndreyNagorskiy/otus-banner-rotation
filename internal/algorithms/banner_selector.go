package algorithms

import "github.com/AndreyNagorskiy/otus-banner-rotation/internal/model"

type BannerStat struct {
	BannerID    int64
	Impressions int64
	Clicks      int64
}

type BannerSelector interface {
	SelectBanner(stats []BannerStat) (int64, bool)
}

func FromModelBannerStat(stat model.BannerStat) BannerStat {
	return BannerStat{
		BannerID:    stat.BannerID,
		Impressions: stat.Impressions,
		Clicks:      stat.Clicks,
	}
}
