package main

import (
	"flag"

	"k8s.io/klog/v2"
)

func main() {
	var (
		// listenAddress = flag.String("web.listen-address", ":9100", "Address to listen on for web interface and metrics.")
		mountPoint = flag.String("storage.mount-point", "/", "The XFS mount point to monitor (usually / or /var/lib/containerd).")
	)
	klog.InitFlags(nil)
	flag.Parse()

	klog.InfoS("Starting Terminus Exporter", "mount", *mountPoint)

	// 1. 创建 Collector
	// xfsColl := exporter.NewXFSCollector(*mountPoint)

	// // 2. 注册到 Prometheus Registry
	// // 我们使用 NewRegistry 而不是默认的，这样可以去掉 Go 运行时的那些默认指标，保持清爽
	// reg := prometheus.NewRegistry()
	// reg.MustRegister(xfsColl)

	// // 3. 启动 HTTP Handler
	// http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	// klog.InfoS("Listening on", "address", *listenAddress)
	// if err := http.ListenAndServe(*listenAddress, nil); err != nil {
	// 	klog.Fatal(err)
	// }
}
