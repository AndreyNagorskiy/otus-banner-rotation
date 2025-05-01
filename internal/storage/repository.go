package storage

import (
	"context"
	"fmt"

	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SlotExists(ctx context.Context, slotID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM slots WHERE id = $1)",
		slotID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check slot existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) BannerExists(ctx context.Context, bannerID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM banners WHERE id = $1)",
		bannerID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check banner existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) SocialGroupExists(ctx context.Context, groupID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM social_groups WHERE id = $1)",
		groupID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check social group existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) AddBannerToSlot(ctx context.Context, slotID, bannerID int64) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO slot_banners (slot_id, banner_id) 
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, slotID, bannerID)

	return err
}

func (r *Repository) RemoveBannerFromSlot(ctx context.Context, slotID, bannerID int64) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM slot_banners 
		WHERE slot_id = $1 AND banner_id = $2
	`, slotID, bannerID)

	return err
}

func (r *Repository) GetSlotBanners(ctx context.Context, slotID int64) ([]int64, error) {
	var ids []int64
	err := pgxscan.Select(ctx, r.db, &ids, `
		SELECT banner_id FROM slot_banners WHERE slot_id = $1
	`, slotID)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *Repository) IncrementImpression(ctx context.Context, slotID, bannerID, groupID int64) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO banner_stats (slot_id, banner_id, social_group_id, impressions, clicks)
		VALUES ($1, $2, $3, 1, 0)
		ON CONFLICT (slot_id, banner_id, social_group_id)
		DO UPDATE SET impressions = banner_stats.impressions + 1
	`, slotID, bannerID, groupID)

	return err
}

func (r *Repository) IncrementClick(ctx context.Context, slotID, bannerID, groupID int64) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO banner_stats (slot_id, banner_id, social_group_id, impressions, clicks)
		VALUES ($1, $2, $3, 0, 1)
		ON CONFLICT (slot_id, banner_id, social_group_id)
		DO UPDATE SET clicks = banner_stats.clicks + 1
	`, slotID, bannerID, groupID)

	return err
}

func (r *Repository) GetBannerStats(ctx context.Context, slotID, groupID int64) ([]model.BannerStat, error) {
	var stats []model.BannerStat
	err := pgxscan.Select(ctx, r.db, &stats, `
		SELECT id, slot_id, banner_id, social_group_id, impressions, clicks
		FROM banner_stats
		WHERE slot_id = $1 AND social_group_id = $2
	`, slotID, groupID)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
