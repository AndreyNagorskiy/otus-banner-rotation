package grpchandler

import (
	"context"
	"errors"
	"testing"

	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/model"
	pb "github.com/AndreyNagorskiy/otus-banner-rotation/pb/api/bannerrotation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockApplication struct {
	mock.Mock
}

func (m *MockApplication) AddBannerToSlot(ctx context.Context, slotID, bannerID int64) error {
	args := m.Called(ctx, slotID, bannerID)
	return args.Error(0)
}

func (m *MockApplication) RemoveBannerFromSlot(ctx context.Context, slotID, bannerID int64) error {
	args := m.Called(ctx, slotID, bannerID)
	return args.Error(0)
}

func (m *MockApplication) RegisterClick(ctx context.Context, slotID, bannerID, socialGroupID int64) error {
	args := m.Called(ctx, slotID, bannerID, socialGroupID)
	return args.Error(0)
}

func (m *MockApplication) GetBannerForSlot(ctx context.Context, slotID, socialGroupID int64) (int64, error) {
	args := m.Called(ctx, slotID, socialGroupID)
	return args.Get(0).(int64), args.Error(1)
}

func TestAddBanner(t *testing.T) {
	tests := []struct {
		name          string
		mockError     error
		expectedCode  codes.Code
		expectedError bool
	}{
		{
			name:          "success",
			mockError:     nil,
			expectedCode:  codes.OK,
			expectedError: false,
		},
		{
			name:          "slot not found",
			mockError:     model.ErrSlotNotFound,
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name:          "banner not found",
			mockError:     model.ErrBannerNotFound,
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name:          "internal error",
			mockError:     errors.New("internal error"),
			expectedCode:  codes.Internal,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApp := new(MockApplication)
			handler := NewBannerRotationHandler(mockApp)

			mockApp.On("AddBannerToSlot", mock.Anything, int64(1), int64(2)).Return(tt.mockError)

			req := &pb.AddBannerRequest{SlotId: 1, BannerId: 2}
			resp, err := handler.AddBanner(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tt.expectedCode, status.Code(err))
				if tt.mockError != nil {
					assert.Contains(t, err.Error(), tt.mockError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestRemoveBanner(t *testing.T) {
	tests := []struct {
		name          string
		mockError     error
		expectedCode  codes.Code
		expectedError bool
	}{
		{
			name:          "success",
			mockError:     nil,
			expectedCode:  codes.OK,
			expectedError: false,
		},
		{
			name:          "slot not found",
			mockError:     model.ErrSlotNotFound,
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name:          "banner not found",
			mockError:     model.ErrBannerNotFound,
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name:          "internal error",
			mockError:     errors.New("internal error"),
			expectedCode:  codes.Internal,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApp := new(MockApplication)
			handler := NewBannerRotationHandler(mockApp)

			mockApp.On("RemoveBannerFromSlot", mock.Anything, int64(1), int64(2)).Return(tt.mockError)

			req := &pb.RemoveBannerRequest{SlotId: 1, BannerId: 2}
			resp, err := handler.RemoveBanner(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tt.expectedCode, status.Code(err))
				if tt.mockError != nil {
					assert.Contains(t, err.Error(), tt.mockError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestRegisterClick(t *testing.T) {
	tests := []struct {
		name          string
		mockError     error
		expectedCode  codes.Code
		expectedError bool
	}{
		{
			name:          "success",
			mockError:     nil,
			expectedCode:  codes.OK,
			expectedError: false,
		},
		{
			name:          "slot not found",
			mockError:     model.ErrSlotNotFound,
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name:          "banner not found",
			mockError:     model.ErrBannerNotFound,
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name:          "social group not found",
			mockError:     model.ErrSocialGroupNotFound,
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name:          "internal error",
			mockError:     errors.New("internal error"),
			expectedCode:  codes.Internal,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApp := new(MockApplication)
			handler := NewBannerRotationHandler(mockApp)

			mockApp.On("RegisterClick", mock.Anything, int64(1), int64(2), int64(3)).Return(tt.mockError)

			req := &pb.RegisterClickRequest{SlotId: 1, BannerId: 2, SocialGroupId: 3}
			resp, err := handler.RegisterClick(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tt.expectedCode, status.Code(err))
				if tt.mockError != nil {
					assert.Contains(t, err.Error(), tt.mockError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestGetBannerForSlot(t *testing.T) {
	tests := []struct {
		name          string
		mockBannerID  int64
		mockError     error
		expectedCode  codes.Code
		expectedError bool
	}{
		{
			name:          "success",
			mockBannerID:  123,
			mockError:     nil,
			expectedCode:  codes.OK,
			expectedError: false,
		},
		{
			name:          "slot not found",
			mockBannerID:  0,
			mockError:     model.ErrSlotNotFound,
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name:          "social group not found",
			mockBannerID:  0,
			mockError:     model.ErrSocialGroupNotFound,
			expectedCode:  codes.InvalidArgument,
			expectedError: true,
		},
		{
			name:          "internal error",
			mockBannerID:  0,
			mockError:     errors.New("internal error"),
			expectedCode:  codes.Internal,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApp := new(MockApplication)
			handler := NewBannerRotationHandler(mockApp)

			mockApp.On("GetBannerForSlot", mock.Anything, int64(1), int64(2)).Return(tt.mockBannerID, tt.mockError)

			req := &pb.GetBannerRequest{SlotId: 1, SocialGroupId: 2}
			resp, err := handler.GetBannerForSlot(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tt.expectedCode, status.Code(err))
				if tt.mockError != nil {
					assert.Contains(t, err.Error(), tt.mockError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.mockBannerID, resp.BannerId)
			}
		})
	}
}
