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

#ifdef __cplusplus
extern "C" {
#endif
  #include <stdio.h>
  #include <stdlib.h>
  #include <stdbool.h>
  #include <stdint.h>

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
      const int method_type,
      const int metric_type);
  FaissStruct* faiss_create_index_ivfpq(
      const int d,
      const int nlist,
      const int m,
      const int nbits_per_idx,
      const int metric_type);
  FaissStruct* faiss_create_index_binaryivf(
      const int d,
      const int nlist);
  FaissStruct* faiss_read_index(const char* fname, const int method_type);
  FaissStruct* faiss_read_index_ivfpq(const char* fname);
  FaissStruct* faiss_read_index_binaryindex(const char* fname);
  bool faiss_write_index(
      const FaissStruct* st,
      const char* fname,
      const int method_type);
  bool faiss_write_index_ivfpq(
      const FaissStruct* st,
      const char* fname);
  bool faiss_write_index_binaryivf(
      const FaissStruct* st,
      const char* fname);
  bool faiss_train(
      const FaissStruct* st,
      const int nb,
      const float* xb,
      const int method_type);
  bool faiss_train_ivfpq(
      const FaissStruct* st,
      const int nb,
      const float* xb);
  bool faiss_train_binaryivf(
      const FaissStruct* st,
      const int nb,
      const uint8_t* xb);
  int faiss_add(
      const FaissStruct* st,
      const int nb,
      const float* xb,
      const long int* xids,
      const int method_type);
  int faiss_add_ivfpq(
      const FaissStruct* st,
      const int nb,
      const float* xb,
      const long int* xids);
  int faiss_add_binaryivf(
      const FaissStruct* st,
      const int nb,
      const uint8_t* xb,
      const long int* xids);
  bool faiss_search(
      const FaissStruct* st,
      const int k,
      const int nprobe,
      const int nq,
      const float* xq,
      long* I,
      float* D,
      const int method_type);
  bool faiss_search_ivfpq(
      const FaissStruct* st,
      const int k,
      const int nprobe,
      const int nq,
      const float* xq,
      long* I,
      float* D);
  bool faiss_search_binaryivf(
      const FaissStruct* st,
      const int k,
      const int nprobe,
      const int nq,
      const uint8_t* xq,
      long* I,
      float* D);
  int faiss_remove(
      const FaissStruct* st,
      const int size,
      const long int* ids,
      const int method_type);
  int faiss_remove_ivfpq(
      const FaissStruct* st,
      const int size,
      const long int* ids);
  int faiss_remove_binaryivf(
      const FaissStruct* st,
      const int size,
      const long int* ids);
  void faiss_free(FaissStruct* st);
#ifdef __cplusplus
}
#endif
