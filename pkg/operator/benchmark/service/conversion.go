// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
)

// jobsByAppNamePool pools the grouping maps handed to jobReconcile so that
// every reconcile does not reallocate the map and its slices.
//
//nolint:gochecknoglobals // process-wide sync.Pool shared across reconciles by design
var jobsByAppNamePool = sync.Pool{
	New: func() any {
		return make(map[string][]k8s.Job)
	},
}

// jobsByAppName converts a JobList into a map grouped by the "app" label
// (falling back to the job name without its generated suffix). The returned
// map is pooled: callers must hand it back via releaseJobsByAppName after the
// consumer is done with it.
func jobsByAppName(list *k8s.JobList) map[string][]k8s.Job {
	jobs, ok := jobsByAppNamePool.Get().(map[string][]k8s.Job)
	if !ok {
		jobs = make(map[string][]k8s.Job)
	}
	for idx := range list.Items {
		job := list.Items[idx]
		name, ok := job.GetObjectMeta().GetLabels()["app"]
		if !ok {
			jns := strings.Split(job.GetName(), "-")
			name = strings.Join(jns[:len(jns)-1], "-")
		}

		// No capacity hint: pre-sizing every group's slice to len(list.Items)
		// over-allocated O(groups x items) capacity per reconcile; groups are
		// typically small, so let append grow them.
		jobs[name] = append(jobs[name], job)
	}
	return jobs
}

// releaseJobsByAppName resets the map and returns it to the pool. It uses
// clear (rather than truncating each slice to [:0:cap]) so that stale app-name
// keys do not accumulate across pooled reuses and the value slices' backing
// arrays — which retain full k8s.Job objects — are released for GC. The map
// allocation itself is still reused via the pool.
func releaseJobsByAppName(jobs map[string][]k8s.Job) {
	clear(jobs)
	jobsByAppNamePool.Put(jobs)
}
