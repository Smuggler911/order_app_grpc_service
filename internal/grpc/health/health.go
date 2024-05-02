package health

import (
	"golang.org/x/net/context"
	grpcHealthV1 "google.golang.org/grpc/health/grpc_health_v1"
	"log/slog"
)

type HServer struct {
	status grpcHealthV1.HealthCheckResponse_ServingStatus
}

func NewHServer() *HServer {
	return &HServer{
		status: grpcHealthV1.HealthCheckResponse_SERVING,
	}
}

func (h *HServer) Check(ctx context.Context, req *grpcHealthV1.HealthCheckRequest) (*grpcHealthV1.HealthCheckResponse, error) {
	return &grpcHealthV1.HealthCheckResponse{
		Status: h.status,
	}, nil
}
func (h *HServer) Watch(req *grpcHealthV1.HealthCheckRequest, stream grpcHealthV1.Health_WatchServer) error {

	status := h.status

	if err := stream.Send(&grpcHealthV1.HealthCheckResponse{
		Status: status,
	}); err != nil {

		slog.Error("Error sending health check response:", err)
		return err
	}
	<-stream.Context().Done()
	if stream.Context().Err() != nil {
		slog.Error("Stream context error:", stream.Context().Err())
	}

	return nil
}
