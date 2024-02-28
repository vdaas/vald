package vald

import "time"

const (
	TimestampLayout                              = time.RFC3339Nano
	UncommittedAnnotationsKey                    = "vald.vdaas.org/uncommitted"
	UnsavedProcessedVqAnnotationsKey             = "vald.vdaas.org/unsaved-processed-vq"
	UnsavedCreateIndexExecutionNumAnnotationsKey = "vald.vdaas.org/unsaved-create-index-execution"
	LastTimeSaveIndexTimestampAnnotationsKey     = "vald.vdaas.org/last-time-save-index-timestamp"
	IndexCountAnnotationsKey                     = "vald.vdaas.org/index-count"
	LastTimeSnapshotTimestampAnnotationsKey      = "vald.vdaas.org/last-time-snapshot-timestamp"
)
