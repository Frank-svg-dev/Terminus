package quota

type QuotaManager interface {
	SetProjectID(path string, projectID uint32) error
	SetQuota(projectID uint32, limitBytes uint64) error
	FetchAllReports(mountPoint string, typeFlag string) (map[uint32]QuotaReport, error)
	RemoveQuota(dirPath string, projectID uint32) error
}

var (
	ContainerdRootPath = "/var/lib/containerd"
)
