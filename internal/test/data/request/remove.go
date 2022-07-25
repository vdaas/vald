//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package request

import (
	"strconv"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

// GenMultiRemoveReq generates Remove_MultiRequest request.
func GenMultiRemoveReq(num int, cfg *payload.Remove_Config) *payload.Remove_MultiRequest {
	req := &payload.Remove_MultiRequest{
		Requests: make([]*payload.Remove_Request, num),
	}
	for i := 0; i < num; i++ {
		req.Requests[i] = &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: "uuid-" + strconv.Itoa(i+1),
			},
			Config: cfg,
		}
	}

	return req
}
