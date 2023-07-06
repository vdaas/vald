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

#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <iostream>
#include <faiss/IndexFlat.h>
#include <faiss/IndexIVFPQ.h>
#include <faiss/impl/AuxIndexStructures.h>
#include <faiss/index_io.h>
#include <faiss/MetricType.h>
#include "Capi.h"

FaissStruct* faiss_create_index(
    const int d,
    const int nlist,
    const int m,
    const int nbits_per_idx,
    const int metric_type) {
#if DEBUG
  printf(__FUNCTION__);
  printf("\n");
  fflush(stdout);
#endif

  FaissStruct *st = NULL;
  try {
    faiss::IndexFlat *quantizer;
    switch (metric_type) {
      case faiss::METRIC_INNER_PRODUCT:
        quantizer = new faiss::IndexFlat(d, faiss::METRIC_INNER_PRODUCT);
        break;
      case faiss::METRIC_L2:
        quantizer = new faiss::IndexFlat(d, faiss::METRIC_L2);
        break;
      default:
        std::stringstream ss;
        ss << "Capi : " << __FUNCTION__ << "() : Error: no metric type.";
        std::cerr << ss.str() << std::endl;
        return NULL;
    }
    faiss::IndexIVFPQ *index = new faiss::IndexIVFPQ(quantizer, d, nlist, m, nbits_per_idx);
#if DEBUG
    index->verbose = true;
#endif
    st = new FaissStruct{
      static_cast<FaissQuantizer>(quantizer),
      static_cast<FaissIndex>(index)
    };
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
  }

  return st;
}

FaissStruct* faiss_read_index(const char* fname) {
#if DEBUG
  printf(__FUNCTION__);
  printf("\n");
  fflush(stdout);
#endif

  FaissStruct *st = NULL;
  try {
    st = new FaissStruct{
      static_cast<FaissQuantizer>(NULL),
      static_cast<FaissIndex>(faiss::read_index(fname))
    };
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
  }

  return st;
}

bool faiss_write_index(
    const FaissStruct* st,
    const char* fname) {
#if DEBUG
  printf(__FUNCTION__);
  printf("\n");
  fflush(stdout);
#endif

  try {
    faiss::write_index(static_cast<faiss::IndexIVFPQ*>(st->faiss_index), fname);
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
    return false;
  }

  fflush(stdout);
  return true;
}

bool faiss_train(
    const FaissStruct* st,
    const int nb,
    const float* xb) {
#if DEBUG
  printf(__FUNCTION__);
  printf("\n");
  fflush(stdout);
#endif

  try {
#if DEBUG
    printf("is_trained: %d\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->is_trained);
    printf("ntotal: %ld\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal);
#endif
    (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->train(nb, xb);
#if DEBUG
    printf("is_trained: %d\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->is_trained);
    printf("ntotal: %ld\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal);
#endif
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
    return false;
  }

#if DEBUG
  fflush(stdout);
#endif
  return true;
}

int faiss_add(
    const FaissStruct* st,
    const int nb,
    const float* xb,
    const long int* xids ) {
#if DEBUG
  printf(__FUNCTION__);
  printf("\n");
  fflush(stdout);
#endif

  try {
#if DEBUG
    printf("is_trained: %d\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->is_trained);
    printf("ntotal: %ld\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal);
#endif
    (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->add_with_ids(nb, xb, xids);
#if DEBUG
    printf("is_trained: %d\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->is_trained);
    printf("ntotal: %ld\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal);
#endif
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
    return -1;
  }

#if DEBUG
  fflush(stdout);
#endif
  return (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal;
}

bool faiss_search(
    const FaissStruct* st,
    const int k,
    const int nq,
    const float* xq,
    long* I,
    float* D) {
#if DEBUG
  printf(__FUNCTION__);
  printf("\n");
  fflush(stdout);
#endif

  try {
#if DEBUG
    printf("is_trained: %d\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->is_trained);
    printf("ntotal: %ld\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal);
#endif
    (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->search(nq, xq, k, D, I);
#if DEBUG
    printf("I=\n");
    for(int i = 0; i < nq; i++) {
        for(int j = 0; j < k; j++) {
            printf("%5ld ", I[i * k + j]);
        }
        printf("\n");
    }
    printf("D=\n");
    for(int i = 0; i < nq; i++) {
        for(int j = 0; j < k; j++) {
            printf("%7g ", D[i * k + j]);
        }
        printf("\n");
    }
#endif
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
    return false;
  }

#if DEBUG
  fflush(stdout);
#endif
  return true;
}

int faiss_remove(
    const FaissStruct* st,
    const int size,
    const long int* ids) {
#if DEBUG
  printf(__FUNCTION__);
  printf("\n");
  fflush(stdout);
#endif

  try {
#if DEBUG
    printf("is_trained: %d\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->is_trained);
    printf("ntotal: %ld\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal);
#endif
    faiss::IDSelectorArray sel(size, ids);
    (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->remove_ids(sel);
#if DEBUG
    printf("is_trained: %d\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->is_trained);
    printf("ntotal: %ld\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal);
#endif
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
    return -1;
  }

#if DEBUG
  fflush(stdout);
#endif
  return (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal;
}

void faiss_free(FaissStruct* st) {
#if DEBUG
  printf(__FUNCTION__);
  printf("\n");
  fflush(stdout);
#endif

  free(st);
  return;
}
