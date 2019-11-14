#
# Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

%.hdf5:
	test -e $@ || curl -LO http://vectors.erikbern.com/$@
	md5sum -c $*.md5

Fashion-MNIST: fashion-mnist-784-euclidean.hdf5
	@test -e $@ || go run ../indexing/main.go $@ $< -distance-type l2

Glove-25: glove-25-angular.hdf5
	@test -e $@ || go run ../indexing/main.go $@ $< -distance-type angle

Glove-50: glove-50-angular.hdf5
	@test -e $@ || go run ../indexing/main.go $@ $< -distance-type angle

Glove-100: glove-100-angular.hdf5
	@test -e $@ || go run ../indexing/main.go $@ $< -distance-type angle

Glove-200: glove-200-angular.hdf5
	@test -e $@ || go run ../indexing/main.go $@ $< -distance-type angle

MNIST: mnist-784-euclidean.hdf5
	@test -e $@ || go run ../indexing/main.go $@ $< -distance-type l2

NYTimes: nytimes-256-angular.hdf5
	@test -e $@ || go run ../indexing/main.go $@ $< -distance-type angle

SIFT: sift-128-euclidean.hdf5
	@test -e $@ || go run ../indexing/main.go $@ $< -distance-type l2

create: Fashion-MNIST Glove-25 Glove-50 Glove-100 MNIST NYTimes SIFT #Glove-200

clean:
	rm -rf ./Fashion-MNIST
	rm -rf ./Glove-25
	rm -rf ./Glove-50
	rm -rf ./Glove-100
	#rm -rf ./Glove-200
	rm -rf ./MNIST
	rm -rf ./NYTimes
	rm -rf ./SIFT

.PHONY: clean create
