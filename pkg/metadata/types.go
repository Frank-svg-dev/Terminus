package metadata

type ContainerInfo struct {
	// --- 核心索引 ---
	ProjectID uint32 `json:"project_id"` 
	Namespace     string `json:"namespace"` 
	PodName       string `json:"pod"`      
	ContainerName string `json:"container"`
}
