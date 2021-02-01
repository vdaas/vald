podname=`kubectl get pods --selector=app=vald-meta-gateway | tail -1 | awk '{print $1}'`
# NAME=TestE2EInsert
NAME=TestE2EUpdate
# NAME=TestE2ERemove
# NAME=TestE2ERemove
# NAME=TestE2ESearch
go test \
  -run ${NAME} \
  -v tests/e2e/crud_test.go \
  -tags "e2e" \
  -timeout 15m \
  -host=127.0.0.1 \
  -port=8082 \
  -dataset=${GOPATH}/src/github.com/vdaas/vald/hack/benchmark/assets/dataset/fashion-mnist-784-euclidean.hdf5 \
  -insert-num=1000 \
  -search-num=1000 \
  -search-by-id-num=10 \
  -get-object-num=10 \
  -update-num=3 \
  -remove-num=2 \
  -wait-after-insert=2m \
  -portforward \
  -portforward-ns=default \
  -portforward-pod-name=${podname} \
  -portforward-pod-port=8081

  # -run TestE2ERemove \
