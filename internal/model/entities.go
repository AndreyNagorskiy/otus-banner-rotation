package model

import "errors"

var (
	ErrSlotNotFound        = errors.New("slot not found")
	ErrBannerNotFound      = errors.New("banner not found")
	ErrSocialGroupNotFound = errors.New("social group not found")
)

type Slot struct {
	ID          int64  `db:"id"`
	Description string `db:"description"`
}

type Banner struct {
	ID          int64  `db:"id"`
	Description string `db:"description"`
}

type SocialGroup struct {
	ID          int64  `db:"id"`
	Description string `db:"description"`
}

type SlotBanner struct {
	ID       int64 `db:"id"`
	SlotID   int64 `db:"slot_id"`
	BannerID int64 `db:"banner_id"`
}

type BannerStat struct {
	ID            int64 `db:"id"`
	SlotID        int64 `db:"slot_id"`
	BannerID      int64 `db:"banner_id"`
	SocialGroupID int64 `db:"social_group_id"`
	Impressions   int64 `db:"impressions"`
	Clicks        int64 `db:"clicks"`
}
