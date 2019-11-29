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

# $1 = directory
# $2 = schema
# $3 = target

package=$(echo $1 | sed -e 's:/$::' | awk -F "/" '{ print $NF }')
config=apis/graphql/$package/gqlgen.yml

if [ ! -f $config ]; then

    cat >$config <<EOF
schema: $2
exec:
  filename: $3
  package: $package
model:
  filename: apis/graphql/$package/models.generated.go
  package: $package
models:
  Controll_CreateIndexRequestInput:
    model: github.com/vdaas/vald/apis/grpc/payload.Controll_CreateIndexRequest
  Insert_MultiRequestInput:
    model: github.com/vdaas/vald/apis/grpc/payload.Insert_MultiRequest
  Insert_RequestInput:
    model: github.com/vdaas/vald/apis/grpc/payload.Insert_Request
  Object_IDInput:
    model: github.com/vdaas/vald/apis/grpc/payload.Object_ID
  Object_IDsInput:
    model: github.com/vdaas/vald/apis/grpc/payload.Object_IDs
  Object_VectorInput:
    model: github.com/vdaas/vald/apis/grpc/payload.Object_Vector
  Search_RequestInput:
    model: github.com/vdaas/vald/apis/grpc/payload.Search_Request
  Update_MultiRequestInput:
    model: github.com/vdaas/vald/apis/grpc/payload.Update_MultiRequest
  Update_RequestInput:
    model: github.com/vdaas/vald/apis/grpc/payload.Update_Request
  Empty:
    model: github.com/vdaas/vald/apis/grpc/payload.Empty
  Empty:
    model: github.com/vdaas/vald/apis/grpc/payload.Empty
  Object_Data:
    model: github.com/vdaas/vald/apis/grpc/payload.Object_Data
  Object_Distance:
    model: github.com/vdaas/vald/apis/grpc/payload.Object_Distance
  Object_ID:
    model: github.com/vdaas/vald/apis/grpc/payload.Object_ID
  Object_Vector:
    model: github.com/vdaas/vald/apis/grpc/payload.Object_Vector
  Search_Response:
    model: github.com/vdaas/vald/apis/grpc/payload.Search_Response
EOF

fi

gqlgen generate -c $config
