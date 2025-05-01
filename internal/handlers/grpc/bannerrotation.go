package grpchandler

import (
	"context"
	"errors"

	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/app"
	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/model"
	"github.com/AndreyNagorskiy/otus-banner-rotation/pb/api/bannerrotation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BannerRotationHandler struct {
	pb.UnimplementedBannerRotationServiceServer
	app app.Application
}

func NewBannerRotationHandler(app app.Application) *BannerRotationHandler {
	return &BannerRotationHandler{app: app}
}

func (h *BannerRotationHandler) AddBanner(
	ctx context.Context,
	req *pb.AddBannerRequest,
) (*pb.AddBannerResponse, error) {
	err := h.app.AddBannerToSlot(ctx, req.SlotId, req.BannerId)
	if err != nil {
		if errors.Is(err, model.ErrSlotNotFound) || errors.Is(err, model.ErrBannerNotFound) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddBannerResponse{}, nil
}

func (h *BannerRotationHandler) RemoveBanner(
	ctx context.Context,
	req *pb.RemoveBannerRequest,
) (*pb.RemoveBannerResponse, error) {
	err := h.app.AddBannerToSlot(ctx, req.SlotId, req.BannerId)
	if err != nil {
		if errors.Is(err, model.ErrSlotNotFound) || errors.Is(err, model.ErrBannerNotFound) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemoveBannerResponse{}, nil
}

func (h *BannerRotationHandler) RegisterClick(
	ctx context.Context,
	req *pb.RegisterClickRequest,
) (*pb.RegisterClickResponse, error) {
	err := h.app.RegisterClick(ctx, req.SlotId, req.BannerId, req.SocialGroupId)
	if err != nil {
		if errors.Is(err, model.ErrSlotNotFound) ||
			errors.Is(err, model.ErrBannerNotFound) ||
			errors.Is(err, model.ErrSocialGroupNotFound) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterClickResponse{}, nil
}

func (h *BannerRotationHandler) GetBannerForSlot(
	ctx context.Context,
	req *pb.GetBannerRequest,
) (*pb.GetBannerResponse, error) {
	//TODO implement
	return nil, status.Error(codes.Unimplemented, "")
}
