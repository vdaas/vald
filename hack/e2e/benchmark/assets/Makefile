DATASETS = \
	fashion-mnist-784-euclidean \
	glove-25-angular \
	glove-50-angular \
	glove-100-angular \
	glove-200-angular \
	mnist-784-euclidean \
	nytimes-256-angular \
	sift-128-euclidean

HDF5 = $(addsuffix .hdf5,$(DATASETS))

%.hdf5:
	./download.sh $@

download: $(HDF5)
