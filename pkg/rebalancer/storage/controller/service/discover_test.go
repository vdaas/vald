package service

import (
	"fmt"
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"
)

var (
	jobTpl = `
apiVersion: batch/v1
kind: Job
metadata:
  name: cassandra-init
spec:
  template:
    spec:
      containers:
        - name: cassandra-init
          image: "cassandra:latest"
          imagePullPolicy: Always
  `

	jobTpl2 = `
apiVersion: batch/v1
kind: Job
metadata:
  name: cassandra-init---
spec:
  template:
    spec:
      containers:
        - name: cassandra-init
          image: "cassandra:latest"
          imagePullPolicy: Always
  `
)

func TestDecode(t *testing.T) {
	decover, err := conversion.NewDecoder(runtime.NewScheme())
	if err != nil {
		t.Fatal(err)
	}

	var job batchv1.Job
	if err := decover.DecodeInto([]byte(jobTpl), &job); err != nil {
		t.Fatal(err)
	}

	var job2 batchv1.Job
	if err := decover.DecodeInto([]byte(jobTpl), &job2); err != nil {
		t.Fatal(err)
	}

	fmt.Println(equality.Semantic.DeepEqual(job, job2))
}
