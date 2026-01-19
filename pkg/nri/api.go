package nri

import (
	"context"

	"github.com/containerd/nri/pkg/api"
)

// Hook 是所有业务逻辑必须实现的接口
type Hook interface {
	Name() string
	Process(ctx context.Context, pod *api.PodSandbox, container *api.Container) error
	Start(ctx context.Context, pod *api.PodSandbox, container *api.Container) error
	Stop(ctx context.Context, pod *api.PodSandbox, container *api.Container) error
}
