package vald

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/manager"
)

func (b *VrsBuilder) buildManager() manager.Manager {
	indexer := b.CR.Spec.VectorEngine.Vald.Indexer

	m := manager.Manager{
		Index: manager.Index{
			Logging: b.buildLogging(indexer.LogLevel),
			Enabled: indexer.Manager,
		},
	}

	// TODO: Discuss with the OSS team whether manager mode is ever disabled in practice.
	// If it is never disabled, this branch can be removed and the indexer setup hardcoded.
	if m.Index.Enabled {
		m.Index.Indexer = manager.Indexer{
			AutoIndexDurationLimit:     "1h",
			AutoSaveIndexDurationLimit: "-1h",
			AutoIndexCheckDuration:     indexer.IndexDuration,
			AutoSaveIndexWaitDuration:  indexer.SaveDuration,
			Concurrency:                &indexer.Concurrency,
		}
		return m
	}

	m.Index.Creator = &manager.Creator{
		Enabled:     true,
		Schedule:    indexer.IndexSchedule,
		Suspend:     indexer.IndexSuspend,
		Concurrency: &indexer.Concurrency,
	}
	m.Index.Saver = &manager.Saver{
		Enabled:  true,
		Schedule: indexer.SaveSchedule,
		Suspend:  indexer.SaveSuspend,
	}
	return m
}
