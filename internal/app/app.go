package app

import (
	"context"

	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/algorithms"
	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/logger"
	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/model"
)

type App struct {
	logger logger.Logger
	rep    Repository
	bs     algorithms.BannerSelector
}

type Repository interface {
	AddBannerToSlot(ctx context.Context, slotID, bannerID int64) error
	RemoveBannerFromSlot(ctx context.Context, slotID, bannerID int64) error
	GetSlotBanners(ctx context.Context, slotID int64) ([]int64, error)
	IncrementImpression(ctx context.Context, slotID, bannerID, groupID int64) error
	IncrementClick(ctx context.Context, slotID, bannerID, groupID int64) error
	GetBannerStats(ctx context.Context, slotID, groupID int64) ([]model.BannerStat, error)
	SlotExists(ctx context.Context, slotID int64) (bool, error)
	BannerExists(ctx context.Context, bannerID int64) (bool, error)
	SocialGroupExists(ctx context.Context, groupID int64) (bool, error)
}

type Application interface {
	AddBannerToSlot(ctx context.Context, slotID, bannerID int64) error
	RemoveBannerFromSlot(ctx context.Context, slotID, bannerID int64) error
	RegisterClick(ctx context.Context, slotID, bannerID, groupID int64) error
	GetBannerForSlot(ctx context.Context, slotID, groupID int64) (int64, error)
}

func New(logger logger.Logger, rep Repository, bs algorithms.BannerSelector) *App {
	return &App{
		logger: logger,
		rep:    rep,
		bs:     bs,
	}
}

func (a *App) AddBannerToSlot(ctx context.Context, slotID, bannerID int64) error {
	err := a.checkSlotAndBannerExistence(ctx, slotID, bannerID)
	if err != nil {
		return err
	}

	err = a.rep.AddBannerToSlot(ctx, slotID, bannerID)
	if err != nil {
		a.logger.Error("failed to add banner to slot", "slotID", slotID, "bannerID", bannerID, "error", err.Error())
	}

	return err
}

func (a *App) RemoveBannerFromSlot(ctx context.Context, slotID, bannerID int64) error {
	err := a.checkSlotAndBannerExistence(ctx, slotID, bannerID)
	if err != nil {
		return err
	}

	err = a.rep.RemoveBannerFromSlot(ctx, slotID, bannerID)
	if err != nil {
		a.logger.Error("failed to remove banner from slot", "slotID", slotID, "bannerID", bannerID, "error", err.Error())
	}

	return err
}

func (a *App) RegisterClick(ctx context.Context, slotID, bannerID, groupID int64) error {
	err := a.checkSlotAndBannerExistence(ctx, slotID, bannerID)
	if err != nil {
		return err
	}

	socialGroupExists, err := a.rep.SocialGroupExists(ctx, groupID)
	if err != nil {
		a.logger.Error("failed to check social group existence", "groupID", groupID, "error", err.Error())
		return err
	}

	if !socialGroupExists {
		return model.ErrSocialGroupNotFound
	}

	err = a.rep.IncrementClick(ctx, slotID, bannerID, groupID)
	if err != nil {
		a.logger.Error(
			"failed to register click",
			"slotID", slotID,
			"bannerID", bannerID,
			"groupID", groupID,
			"error", err.Error(),
		)
	}

	return err
}

func (a *App) GetBannerForSlot(ctx context.Context, slotID, socialGroupID int64) (int64, error) {
	slotExists, err := a.rep.SlotExists(ctx, slotID)
	if err != nil {
		a.logger.Error("failed to check slot existence", "slotID", slotID, "error", err.Error())
		return 0, err
	}

	if !slotExists {
		return 0, model.ErrSlotNotFound
	}
	socialGroupExists, err := a.rep.SocialGroupExists(ctx, socialGroupID)
	if err != nil {
		a.logger.Error(
			"failed to check social group existence",
			"socialGroupID", socialGroupID,
			"error", err.Error(),
		)
		return 0, err
	}

	if !socialGroupExists {
		return 0, model.ErrSocialGroupNotFound
	}

	bannerStats, err := a.rep.GetBannerStats(ctx, slotID, socialGroupID)
	if err != nil {
		a.logger.Error("failed to get banner stats", "slotID", slotID, "error", err.Error())
		return 0, err
	}

	if len(bannerStats) == 0 {
		return 0, model.ErrBannerNotFound
	}

	stats := make([]algorithms.BannerStat, len(bannerStats))

	for i, s := range bannerStats {
		stats[i] = algorithms.FromModelBannerStat(s)
	}

	bannerID, exists := a.bs.SelectBanner(stats)
	if !exists {
		a.logger.Error("failed to select banner", "slotID", slotID, "error", "no suitable banner found")
		return 0, model.ErrBannerNotFound
	}

	return bannerID, nil
}

func (a *App) checkSlotAndBannerExistence(ctx context.Context, slotID, bannerID int64) error {
	slotExists, err := a.rep.SlotExists(ctx, slotID)
	if err != nil {
		a.logger.Error("failed to check slot existence", "slotID", slotID, "error", err.Error())
		return err
	}

	if !slotExists {
		return model.ErrSlotNotFound
	}

	bannerExists, err := a.rep.BannerExists(ctx, bannerID)
	if err != nil {
		a.logger.Error("failed to check banner existence", "bannerID", bannerID, "error", err.Error())
		return err
	}

	if !bannerExists {
		return model.ErrBannerNotFound
	}

	return nil
}
