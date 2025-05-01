package algorithms

import (
	"math"
)

type BanditUCB1 struct{}

func NewBanditUCB1() *BanditUCB1 {
	return &BanditUCB1{}
}

// calculateUCB вычисляет значение UCB1 для баннера.
func (b *BanditUCB1) calculateUCB(stat BannerStat, totalImpressions int64) float64 {
	ctr := float64(stat.Clicks) / float64(stat.Impressions)

	return ctr + math.Sqrt(2*math.Log(float64(totalImpressions))/float64(stat.Impressions))
}

// SelectBanner выбирает баннер по алгоритму UCB1.
func (b *BanditUCB1) SelectBanner(stats []BannerStat) (int64, bool) {
	if len(stats) == 0 {
		return 0, false
	}

	var (
		totalImpressions int64
		bestBannerID     int64
		maxUCB           float64 = -1
	)

	// Считаем общее количество показов
	for _, stat := range stats {
		totalImpressions += stat.Impressions
	}

	for _, stat := range stats {
		// Если баннер ещё не показывался, выбираем его сразу
		if stat.Impressions == 0 {
			return stat.BannerID, true
		}

		ucb := b.calculateUCB(stat, totalImpressions)

		// Выбираем баннер с максимальным UCB1
		if ucb > maxUCB {
			maxUCB = ucb
			bestBannerID = stat.BannerID
		}
	}

	return bestBannerID, true
}
