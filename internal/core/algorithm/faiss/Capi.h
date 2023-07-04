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

#ifdef __cplusplus
extern "C" {
#endif
  #include <stdio.h>
  #include <stdlib.h>
  #include <stdbool.h>

  typedef void* FaissQuantizer;
  typedef void* FaissIndex;
  typedef struct {
    FaissQuantizer  faiss_quantizer;
    FaissIndex      faiss_index;
  } FaissStruct;

  FaissStruct* faiss_create_index(
      const int d,
      const int nlist,
      const int m,
      const int nbits_per_idx,
      const int metric_type);
  FaissStruct* faiss_read_index(const char* fname);
  bool faiss_write_index(
      const FaissStruct* st,
      const char* fname);
  bool faiss_train(
      const FaissStruct* st,
      const int nb,
      const float* xb);
  int faiss_add(
      const FaissStruct* st,
      const int nb,
      const float* xb,
      const long int* xids);
  bool faiss_search(
      const FaissStruct* st,
      const int k,
      const int nq,
      const float* xq,
      long* I,
      float* D);
  int faiss_remove(
      const FaissStruct* st,
      const int size,
      const long int* ids);
  void faiss_free(FaissStruct* st);
#ifdef __cplusplus
}
#endif
