//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package errors

var (
	ErrNGTIndexStatisticsDisabled = New("ngt get statistics is disabled")

	ErrNGTIndexStatisticsNotReady = New("ngt get statistics is not ready")
)

type NGTError struct {
	Msg string
}

func NewNGTError(msg string) error {
	return NGTError{
		Msg: msg,
	}
}

func (n NGTError) Error() string {
	return n.Msg
}
