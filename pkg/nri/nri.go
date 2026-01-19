package nri

import (
	"context"

	"github.com/containerd/nri/pkg/stub"
	"k8s.io/klog/v2"
)

func NewEnforcer(opts ...Option) (*Enforcer, error) {
	e := &Enforcer{
		SocketPath: "/var/run/nri/nri.sock",
		PluginName: "terminus",
		PluginIdx:  "00",
	}
	for _, opt := range opts {
		opt(e)
	}
	return e, nil
}

func (e *Enforcer) Run(ctx context.Context) error {
	klog.InfoS("Starting Enforcer", "hooks_count", len(e.Hooks))

	opts := []stub.Option{
		stub.WithPluginName(e.PluginName),
		stub.WithPluginIdx(e.PluginIdx),
		stub.WithSocketPath(e.SocketPath),
	}

	// e 实现了 Handler 接口
	st, err := stub.New(e, opts...)
	if err != nil {
		return err
	}
	return st.Run(ctx)
}
