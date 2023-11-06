.PHONY: minikube/install
minikube/install: $(BINDIR)/minikube

$(BINDIR)/minikube:
	mkdir -p $(BINDIR)
	curl -L https://storage.googleapis.com/minikube/releases/latest/minikube-$(shell echo $(UNAME) | tr '[:upper:]' '[:lower:]')-$(subst x86_64,amd64,$(shell echo $(ARCH) | tr '[:upper:]' '[:lower:]')) -o $(BINDIR)/minikube
	chmod a+x $(BINDIR)/minikube

# Start minikube with CSI Driver and Volume Snapshots support
# Only use this for development related to Volume Snapshots. Usually k3d is faster.
.PHONY: minikube/start
minikube/start:
	minikube start --force
	minikube addons enable volumesnapshots
	minikube addons enable csi-hostpath-driver
	minikube addons disable storage-provisioner
	minikube addons disable default-storageclass
	kubectl patch storageclass csi-hostpath-sc -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'

.PHONY: minikube/delete
minikube/delete:
	minikube delete

.PHONY: minikube/restart
minikube/restart:
	@make minikube/delete
	@make minikube/start
