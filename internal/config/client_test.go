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

package config

import (
	"testing"

	"github.com/vdaas/vald/internal/test"
)

func TestClient_Bind(t *testing.T) {
	if err := test.Run(t.Context(), t, func(tt *testing.T, args *Client) (c *Client, err error) {
		tt.Helper()
		if args == nil {
			args = new(Client)
		}
		c = args.Bind()
		return c, nil
	}, []test.Case[*Client, *Client]{
		{
			Name: "returns Client when no argument found",
			Want: test.Result[*Client]{
				Val: new(Client).Bind(),
			},
		},
		{
			Name: "return Client when the bind successes and net and transport is not nil",
			Args: (&Client{
				Net:       new(Net),
				Transport: new(Transport),
			}),
			Want: test.Result[*Client]{
				Val: (&Client{
					Net: new(Net),
					Transport: &Transport{
						RoundTripper: new(RoundTripper),
						Backoff:      new(Backoff),
					},
				}).Bind(),
				Err: nil,
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

// NOT IMPLEMENTED BELOW
