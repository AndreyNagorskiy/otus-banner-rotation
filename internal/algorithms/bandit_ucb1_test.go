package algorithms

import (
	"math"
	"testing"
)

func TestNewBandit(t *testing.T) {
	b := NewBanditUCB1()
	if b == nil {
		t.Fatal("NewBanditUCB1() returned nil")
	}
}

func TestSelectBanner_EmptyStats(t *testing.T) {
	b := NewBanditUCB1()
	id, ok := b.SelectBanner([]BannerStat{})
	if ok || id != 0 {
		t.Errorf("Expected (0, false) for empty stats, got (%d, %v)", id, ok)
	}
}

func TestSelectBanner_SingleBanner(t *testing.T) {
	b := NewBanditUCB1()
	stats := []BannerStat{
		{BannerID: 1, Impressions: 10, Clicks: 2},
	}
	id, ok := b.SelectBanner(stats)
	if !ok || id != 1 {
		t.Errorf("Expected (1, true), got (%d, %v)", id, ok)
	}
}

func TestSelectBanner_NewBannerPriority(t *testing.T) {
	b := NewBanditUCB1()
	stats := []BannerStat{
		{BannerID: 1, Impressions: 100, Clicks: 50}, // CTR = 0.5
		{BannerID: 2, Impressions: 0, Clicks: 0},    // New banner
		{BannerID: 3, Impressions: 200, Clicks: 60}, // CTR = 0.3
	}

	id, ok := b.SelectBanner(stats)
	if !ok || id != 2 {
		t.Errorf("Expected new banner (2, true), got (%d, %v)", id, ok)
	}
}

func TestSelectBanner_MultipleNewBanners(t *testing.T) {
	b := NewBanditUCB1()
	stats := []BannerStat{
		{BannerID: 1, Impressions: 0, Clicks: 0},
		{BannerID: 2, Impressions: 0, Clicks: 0},
	}

	// Проверяем что возвращается один из новых баннеров
	id, ok := b.SelectBanner(stats)
	if !ok || (id != 1 && id != 2) {
		t.Errorf("Expected either (1 or 2, true), got (%d, %v)", id, ok)
	}
}

func TestSelectBanner_UCBSelection(t *testing.T) {
	b := NewBanditUCB1()
	stats := []BannerStat{
		{BannerID: 1, Impressions: 100, Clicks: 10}, // CTR = 0.1
		{BannerID: 2, Impressions: 50, Clicks: 15},  // CTR = 0.3
		{BannerID: 3, Impressions: 200, Clicks: 30}, // CTR = 0.15
	}

	// Общее количество показов = 350
	// Ожидаем что выберется баннер 2, так как у него:
	// - Достаточно высокая CTR
	// - Меньше показов чем у баннера 3 (больше exploration bonus)
	id, ok := b.SelectBanner(stats)
	if !ok || id != 2 {
		t.Errorf("Expected banner 2 to be selected, got %d", id)
	}
}

func TestSelectBanner_AllZeroClicks(t *testing.T) {
	b := NewBanditUCB1()
	stats := []BannerStat{
		{BannerID: 1, Impressions: 100, Clicks: 0},
		{BannerID: 2, Impressions: 50, Clicks: 0},
		{BannerID: 3, Impressions: 200, Clicks: 0},
	}

	// В этом случае выберется баннер с наименьшим количеством показов (2)
	// так как exploration term будет максимальным для него
	id, ok := b.SelectBanner(stats)
	if !ok || id != 2 {
		t.Errorf("Expected banner 2 (least impressions) to be selected, got %d", id)
	}
}

func TestCalculateUCB(t *testing.T) {
	b := NewBanditUCB1()
	tests := []struct {
		name             string
		stat             BannerStat
		totalImpressions int64
		expected         float64
	}{
		{
			name:             "High CTR",
			stat:             BannerStat{Impressions: 100, Clicks: 50},
			totalImpressions: 1000,
			expected:         0.5 + math.Sqrt(2*math.Log(1000)/100),
		},
		{
			name:             "Low impressions",
			stat:             BannerStat{Impressions: 5, Clicks: 1},
			totalImpressions: 100,
			expected:         0.2 + math.Sqrt(2*math.Log(100)/5),
		},
		{
			name:             "Zero clicks",
			stat:             BannerStat{Impressions: 200, Clicks: 0},
			totalImpressions: 500,
			expected:         math.Sqrt(2 * math.Log(500) / 200),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ucb := b.calculateUCB(tt.stat, tt.totalImpressions)
			if math.Abs(ucb-tt.expected) > 1e-9 {
				t.Errorf("Expected UCB %.10f, got %.10f", tt.expected, ucb)
			}
		})
	}
}

func TestSelectBanner_EqualUCB(t *testing.T) {
	// Создаем ситуацию, когда UCB одинаковые
	// Должен выбираться первый баннер с максимальным значением
	b := NewBanditUCB1()
	stats := []BannerStat{
		{BannerID: 1, Impressions: 100, Clicks: 10},
		{BannerID: 2, Impressions: 100, Clicks: 10},
		{BannerID: 3, Impressions: 100, Clicks: 10},
	}

	id, ok := b.SelectBanner(stats)
	if !ok || id != 1 {
		t.Errorf("Expected first banner (1) to be selected, got %d", id)
	}
}

func TestSelectBanner_Precision(t *testing.T) {
	b := NewBanditUCB1()
	stats := []BannerStat{
		{BannerID: 1, Impressions: 1000000, Clicks: 500000}, // CTR = 0.5
		{BannerID: 2, Impressions: 1000000, Clicks: 500001}, // CTR = 0.500001
	}

	// Проверяем, что маленькая разница в CTR корректно обрабатывается
	id, ok := b.SelectBanner(stats)
	if !ok || id != 2 {
		t.Errorf("Expected banner 2 to be selected, got %d", id)
	}
}
