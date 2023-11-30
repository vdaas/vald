//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package version provides version comparison functionality
package version

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/vdaas/vald/internal/errors"
)

func Check(cur, max, min string) error {
	curv, err := version.NewSemver(cur)
	if err != nil {
		return err
	}

	// Constraints example.
	constraints, err := version.NewConstraint(fmt.Sprintf(">= %s, <= %s", min, max))
	if err != nil {
		return err
	}

	if !constraints.Check(curv) {
		return errors.ErrInvalidConfigVersion(curv.String(), constraints.String())
	}

	return nil
}
