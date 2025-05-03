package internalgrpc

import (
	"context"
	"time"

	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func loggingUnaryInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Логируем входящий запрос
		logger.Info("incoming gRPC request",
			"method", info.FullMethod,
			"request", req,
		)

		// Вызываем обработчик
		resp, err := handler(ctx, req)

		// Логируем результат
		duration := time.Since(start)
		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		logger.Info("request processed",
			"method", info.FullMethod,
			"status", statusCode.String(),
			"duration", duration,
			"error", err,
		)

		return resp, err
	}
}
