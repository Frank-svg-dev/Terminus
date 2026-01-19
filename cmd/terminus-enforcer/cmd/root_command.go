package cmd

import (
	"context"
	"flag"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/Frank-svg-dev/Terminus/pkg/exporter"
	"github.com/Frank-svg-dev/Terminus/pkg/hooks"
	"github.com/Frank-svg-dev/Terminus/pkg/metadata"
	"github.com/Frank-svg-dev/Terminus/pkg/nri"
	"github.com/Frank-svg-dev/Terminus/pkg/quota/xfs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	// 引入业务逻辑包
	// "terminus/pkg/enforcer"
)

const (
	socketPath = "/var/run/nri/nri.sock"
	pluginName = "Terminus-Enforcer"
)

// rootCmd 定义根命令
var rootCmd = &cobra.Command{
	Use:   "terminus-enforcer",
	Short: "Terminus NRI Plugin",
	Long:  `Terminus Enforcer listens to NRI events and applies Project Quota limits.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 确保 klog 能够解析 flags
		flag.Parse()
	},
	// 真正的执行入口
	RunE: func(cmd *cobra.Command, args []string) error {

		store := metadata.NewAsyncStore(1000)
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		go store.Run(ctx)

		qm := xfs.NewXFSCLI()

		storageHook := hooks.NewStorageHook(qm, store)

		enforcer, err := nri.NewEnforcer(
			nri.WithSocketPath(socketPath),
			nri.WithPluginName(pluginName),
			nri.WithHook(storageHook),
			// nri.WithHook(logHook),
		)
		if err != nil {
			return err
		}

		go func() {
			collector := exporter.NewXFSCollector("/", store)
			reg := prometheus.NewRegistry()
			reg.MustRegister(collector)

			// 3. 启动 HTTP Handler
			http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

			klog.InfoS("Listening on", "address", ":9201")
			if err := http.ListenAndServe(":9201", nil); err != nil {
				klog.Fatal(err)
			}
		}()

		// 4. 启动
		return enforcer.Run(cmd.Context())
	},
}

// Execute 是 main.go 调用的函数
func Execute() error {
	return rootCmd.Execute()
}

// init 初始化 Flags
func init() {
	klog.InitFlags(nil)
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	_ = flag.Set("logtostderr", "true")
}
