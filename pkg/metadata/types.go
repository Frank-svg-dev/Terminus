package metadata

// ContainerInfo 定义了容器的元数据信息
// 这些字段主要用于：
// 1. 在 NRI 阶段被记录 (Set)
// 2. 在 Exporter 阶段被读取并转化为 Prometheus Labels (Get)
type ContainerInfo struct {
	// --- 核心索引 ---
	ProjectID uint32 `json:"project_id"` // 关联 XFS Quota 的唯一 ID

	// --- K8s 身份信息 (用于监控 Labels) ---
	Namespace     string `json:"namespace"` // e.g. "default"
	PodName       string `json:"pod"`       // e.g. "nginx-demo-xxx"
	ContainerName string `json:"container"` // e.g. "nginx"

	// --- 扩展信息 (可选，视监控需求决定是否开启) ---
	// Image         string `json:"image,omitempty"`           // 镜像名，方便统计不同镜像的磁盘消耗
	// PodUID        string `json:"pod_uid,omitempty"`         // Pod 唯一 ID，防止同名 Pod 混淆
	// QoSClass      string `json:"qos_class,omitempty"`       // e.g. "Burstable", "BestEffort"
	// WorkloadName  string `json:"workload_name,omitempty"`   // 所属 Deployment/StatefulSet 名称
}
