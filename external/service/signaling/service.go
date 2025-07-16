package signaling

import (
	internalsignaling "github.com/kizuwe1lf/fuel/internal/service/signaling"
	signalingv1 "github.com/kizuwe1lf/fuel/proto/gen/signaling/v1"
)

type Service struct {
	signalingv1.UnimplementedSignalingServer
	internalsignaling *internalsignaling.Service
}

func NewService(service *internalsignaling.Service) *Service {
	return &Service{internalsignaling: service}
}
