package exporter

import (
	"github.com/Frank-svg-dev/Terminus/pkg/metadata"
	"github.com/Frank-svg-dev/Terminus/pkg/quota"
)

func SaveState(qm quota.QuotaManager, path string, store *metadata.AsyncStore) {
	qutoaReport, err := qm.FetchAllReports(path, "b")
	if err != nil {
		return
	}

	for k := range qutoaReport {
		store.TriggerUpdate(k, metadata.ContainerInfo{})
	}
}
