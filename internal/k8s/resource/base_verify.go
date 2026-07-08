//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package resource

import (
	"reflect"

	"github.com/vdaas/vald/internal/errors"
)

// VerifyBase validates the first-field embedding contract documented on Base
// for the type T: the first field of T must be the anonymous Base[T, PT]
// embed located at offset 0. Base.self relies on this layout to recover the
// outer object from the embedded Base pointer, so a violation corrupts every
// promoted DeepCopy call. VerifyBase is a test helper; call it once from a
// test in every package that embeds Base.
func VerifyBase[T any, PT DeepCopyIntoer[T]]() error {
	rt := reflect.TypeFor[T]()
	if rt.Kind() != reflect.Struct {
		return errors.Errorf("%s must be a struct to embed resource.Base", rt)
	}
	if rt.NumField() == 0 {
		return errors.Errorf("%s has no fields; resource.Base must be embedded as the first field", rt)
	}
	f := rt.Field(0)
	if bt := reflect.TypeFor[Base[T, PT]](); f.Type != bt {
		return errors.Errorf("%s field 0 has type %s, want %s embedded as the first field", rt, f.Type, bt)
	}
	if !f.Anonymous {
		return errors.Errorf("%s field 0 (%s) must be an anonymous resource.Base embed", rt, f.Name)
	}
	if f.Offset != 0 {
		return errors.Errorf("%s field 0 (%s) sits at offset %d, want 0", rt, f.Name, f.Offset)
	}
	return nil
}
