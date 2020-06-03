//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package compress provides compressor interface
package compress

import "io"

type Compressor interface {
	CompressVector(vector []float32) (bytes []byte, err error)
	DecompressVector(bytes []byte) (vector []float32, err error)
	Reader(src io.Reader) (io.Reader, error)
	Writer(dst io.Writer) (io.WriteCloser, error)
}
