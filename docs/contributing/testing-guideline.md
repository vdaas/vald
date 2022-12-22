# Testing guideline

Testing your changes is very important to ensure your implementation is working as expected.

## Unit test

Unit test is used to check whether the unit is implemented correctly in various cases.

We suggest you implement or update the unit test code when making logical changes or implementing new functionality in Vald.

Before implementing the unit test code, we suggest you read our [unit test guideline](./docs/contributing/unit-test-guideline.md) to guide you to create good unit tests and [coding guideline for unit test](./docs/contributing/coding-style.md#test) to guide you to implement unit tests.

If you want to execute the unit test on only part of the code, you can use `go test` command to execute the unit test on the specific package/function.
For example, if you want to execute the unit test on a specific package, use the following command.

```bash
go test -race [package]
```

This command will execute the unit test on the package, and also enable the race detector to check if any race occurs in the implementation.

If you want to execute the unit test on the whole Vald implementation, Vald provides the following command to do that.

```bash
make test
```

This command will execute all unit tests of `*target*_test.go` files on `cmd`, `internal` and `pkg` packages. It is useful to ensure that your changes will not affect the behavior of other components and packages.

## End-To-End (E2E) testing

End-To-End (E2E) testing is used to test the application flow of Vald is working as expected from beginning to end.

If you want to execute E2E test on Vald, Vald provides the following commands to test the implementation.

```bash
make e2e
```

The command `make e2e` execute E2E tests to ensure whether the functionality is working as expected. It will perform the actual CRUD action on a cluster and verify the result.

E2E tests require deploying Vald on a Kubernetes cluster beforehand. You can deploy Vald on your Kubernetes cluster, or you can create a Kubernetes cluster on your local machine easily by using the tools like [k3d](https://k3d.io/) or [kind](https://kind.sigs.k8s.io/).

Please refer to our [get started](./docs/tutorial/get-started.md) to create the cluster and deploy Vald on a Kubernetes cluster.

If you want to execute E2E test on your Kubernetes cluster, you may need to modify the configuration on [Makefile](https://github.com/vdaas/vald/blob/main/Makefile) before executing the E2E test.

| Config name                        | Description                                                                   |
| :--------------------------------- | :---------------------------------------------------------------------------- |
| E2E_BIND_HOST                      | The target host of Kubernetes cluster                                         |
| E2E_BIND_PORT                      | The target port of Kubernetes cluster                                         |
| E2E_TIMEOUT                        | The timeout of E2E test                                                       |
| E2E_DATASET_NAME                   | The dataset name of the E2E test                                              |
| E2E_INSERT_COUNT                   | The number of index insert in E2E test                                        |
| E2E_SEARCH_COUNT                   | The number of search request in E2E test                                      |
| E2E_SEARCH_BY_ID_COUNT             | The number of search by ID request in E2E test                                |
| E2E_GET_OBJECT_COUNT               | The number of get object request in E2E test                                  |
| E2E_UPDATE_COUNT                   | The number of update request in E2E test                                      |
| E2E_UPSERT_COUNT                   | The number of upsert request in E2E test                                      |
| E2E_REMOVE_COUNT                   | The number of remove request in E2E test                                      |
| E2E_WAIT_FOR_CREATE_INDEX_DURATION | The wait time of create index operation after insert is completed in E2E test |
| E2E_TARGET_NAME                    | The target pod name in the Vald cluster                                       |
| E2E_TARGET_NAMESPACE               | The target namespace of the Vald cluster                                      |
| E2E_TARGET_PORT                    | The pod forward port of the target pod to the local host                      |
| E2E_PORTFORWARD_ENABLED            | Enable/Disable port forwarding                                                |
