package app

import (
	"context"

	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/logger"
	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/model"
)

type App struct {
	logger logger.Logger
	rep    Repository
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
}

type Application interface {
	AddBannerToSlot(ctx context.Context, slotID, bannerID int64) error
}

func New(logger logger.Logger, rep Repository) *App {
	return &App{
		logger: logger,
		rep:    rep,
	}
}

func (a *App) AddBannerToSlot(ctx context.Context, slotID, bannerID int64) error {
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

	err = a.rep.AddBannerToSlot(ctx, slotID, bannerID)
	if err != nil {
		a.logger.Error("failed to add banner to slot", "slotID", slotID, "bannerID", bannerID, "error", err.Error())
	}

	return err
}
