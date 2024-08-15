//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>
#include <iostream>
#include <faiss/IndexBinaryFlat.h>
#include <faiss/IndexBinaryIVF.h>
#include <faiss/IndexFlat.h>
#include <faiss/IndexIVFPQ.h>
#include <faiss/impl/AuxIndexStructures.h>
#include <faiss/index_io.h>
#include <faiss/MetricType.h>
#include "Capi.h"

enum MethodType {
  IVFPQ = 0,
  BINARYIVF = 1,
};

FaissStruct* faiss_create_index(
    const int d,
    const int nlist,
    const int m,
    const int nbits_per_idx,
    const int method_type,
    const int metric_type) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  switch (method_type) {
    case IVFPQ:
      return faiss_create_index_ivfpq(d, nlist, m, nbits_per_idx, metric_type);
    case BINARYIVF:
      return faiss_create_index_binaryivf(d*8, nlist);
    default:
      std::stringstream ss;
      ss << "Capi : " << __FUNCTION__ << "() : Error: no method type.";
      std::cerr << ss.str() << std::endl;
      return NULL;
  }
}

FaissStruct* faiss_create_index_ivfpq(
    const int d,
    const int nlist,
    const int m,
    const int nbits_per_idx,
    const int metric_type) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  FaissStruct *st = NULL;
  try {
    faiss::IndexFlat *quantizer;
    switch (metric_type) {
      case faiss::METRIC_INNER_PRODUCT:
        quantizer = new faiss::IndexFlat(d, faiss::METRIC_INNER_PRODUCT);
      case faiss::METRIC_L2:
        quantizer = new faiss::IndexFlat(d, faiss::METRIC_L2);
      default:
        std::stringstream ss;
        ss << "Capi : " << __FUNCTION__ << "() : Error: no metric type.";
        std::cerr << ss.str() << std::endl;
        return NULL;
    }
    faiss::IndexIVFPQ *index = new faiss::IndexIVFPQ(quantizer, d, nlist, m, nbits_per_idx);
    //index->verbose = true;
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

FaissStruct* faiss_create_index_binaryivf(
    const int d,
    const int nlist) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  FaissStruct *st = NULL;
  try {
    faiss::IndexBinaryFlat *quantizer;
    quantizer = new faiss::IndexBinaryFlat(d);
    faiss::IndexBinaryIVF *index = new faiss::IndexBinaryIVF(quantizer, d, nlist);
    //index->verbose = true;
    st = new FaissStruct{
      static_cast<FaissQuantizer>(quantizer),
      static_cast<FaissIndex>(index),
    };
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
  }

  return st;
}

FaissStruct* faiss_read_index(const char* fname, const int method_type) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  switch (method_type) {
    case IVFPQ:
      return faiss_read_index_ivfpq(fname);
    case BINARYIVF:
      return faiss_read_index_binaryindex(fname);
    default:
      std::stringstream ss;
      ss << "Capi : " << __FUNCTION__ << "() : Error: no method type.";
      std::cerr << ss.str() << std::endl;
      return NULL;
  }
}

FaissStruct* faiss_read_index_ivfpq(const char* fname) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

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

FaissStruct* faiss_read_index_binaryindex(const char* fname) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  FaissStruct *st = NULL;
  try {
    st = new FaissStruct{
      static_cast<FaissQuantizer>(NULL),
      static_cast<FaissIndex>(faiss::read_index_binary(fname))
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
    const char* fname,
    const int method_type) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  switch (method_type) {
    case IVFPQ:
      return faiss_write_index_ivfpq(st, fname);
    case BINARYIVF:
      return faiss_write_index_binaryivf(st, fname);
    default:
      std::stringstream ss;
      ss << "Capi : " << __FUNCTION__ << "() : Error: no method type.";
      std::cerr << ss.str() << std::endl;
      return false;
  }
}

bool faiss_write_index_ivfpq(
    const FaissStruct* st,
    const char* fname) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

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

bool faiss_write_index_binaryivf(
    const FaissStruct* st,
    const char* fname) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  try {
    faiss::write_index_binary(static_cast<faiss::IndexBinaryIVF*>(st->faiss_index), fname);
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
    const float* xb,
    const int method_type) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  switch (method_type) {
    case IVFPQ:
      return faiss_train_ivfpq(st, nb, xb);
    case BINARYIVF:
      return faiss_train_binaryivf(st, nb, reinterpret_cast<const uint8_t*>(xb));
    default:
      std::stringstream ss;
      ss << "Capi : " << __FUNCTION__ << "() : Error: no method type.";
      std::cerr << ss.str() << std::endl;
      return false;
  }
}

bool faiss_train_ivfpq(
    const FaissStruct* st,
    const int nb,
    const float* xb) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  try {
    //printf("is_trained: %d\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->is_trained);
    //printf("ntotal: %ld\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal);
    (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->train(nb, xb);
    //printf("is_trained: %d\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->is_trained);
    //printf("ntotal: %ld\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal);
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
    return false;
  }

  fflush(stdout);
  return true;
}

bool faiss_train_binaryivf(
    const FaissStruct* st,
    const int nb,
    const uint8_t* xb) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  try {
    (static_cast<faiss::IndexBinaryIVF*>(st->faiss_index))->train(nb, xb);
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
    return false;
  }

  fflush(stdout);
  return true;
}

int faiss_add(
    const FaissStruct* st,
    const int nb,
    const float* xb,
    const long int* xids,
    const int method_type) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  switch (method_type) {
    case IVFPQ:
      return faiss_add_ivfpq(st, nb, xb, xids);
    case BINARYIVF:
      return faiss_add_binaryivf(st, nb, reinterpret_cast<const uint8_t*>(xb), xids);
    default:
      std::stringstream ss;
      ss << "Capi : " << __FUNCTION__ << "() : Error: no method type.";
      std::cerr << ss.str() << std::endl;
      return -1;
  }
}

int faiss_add_ivfpq(
    const FaissStruct* st,
    const int nb,
    const float* xb,
    const long int* xids) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  try {
    //printf("is_trained: %d\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->is_trained);
    //printf("ntotal: %ld\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal);
    (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->add_with_ids(nb, xb, xids);
    //printf("is_trained: %d\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->is_trained);
    //printf("ntotal: %ld\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal);
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
    return -1;
  }

  fflush(stdout);
  return (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal;
}

int faiss_add_binaryivf(
    const FaissStruct* st,
    const int nb,
    const uint8_t* xb,
    const long int* xids) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  try {
    (static_cast<faiss::IndexBinaryIVF*>(st->faiss_index))->add_with_ids(nb, xb, xids);
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
    return -1;
  }

  fflush(stdout);
  return (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal;
}

bool faiss_search(
    const FaissStruct* st,
    const int k,
    const int nprobe,
    const int nq,
    const float* xq,
    long* I,
    float* D,
    const int method_type) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  switch (method_type) {
    case IVFPQ:
      return faiss_search_ivfpq(st, k, nprobe, nq, xq, I, D);
    case BINARYIVF:
      return faiss_search_binaryivf(st, k, nprobe, nq, reinterpret_cast<const uint8_t*>(xq), I, D);
    default:
      std::stringstream ss;
      ss << "Capi : " << __FUNCTION__ << "() : Error: no method type.";
      std::cerr << ss.str() << std::endl;
      return false;
  }
}

bool faiss_search_ivfpq(
    const FaissStruct* st,
    const int k,
    const int nprobe,
    const int nq,
    const float* xq,
    long* I,
    float* D) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  try {
    //printf("is_trained: %d\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->is_trained);
    //printf("ntotal: %ld\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal);
    (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->nprobe = nprobe;
    (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->search(nq, xq, k, D, I);
    //printf("I=\n");
    //for(int i = 0; i < nq; i++) {
    //    for(int j = 0; j < k; j++) {
    //        printf("%5ld ", I[i * k + j]);
    //    }
    //    printf("\n");
    //}
    //printf("D=\n");
    //for(int i = 0; i < nq; i++) {
    //    for(int j = 0; j < k; j++) {
    //        printf("%7g ", D[i * k + j]);
    //    }
    //    printf("\n");
    //}
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
    return false;
  }

  return true;
}

bool faiss_search_binaryivf(
    const FaissStruct* st,
    const int k,
    const int nprobe,
    const int nq,
    const uint8_t* xq,
    long* I,
    float* D) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  int32_t* tmpD = new int32_t[nq*k];
  try {
    (static_cast<faiss::IndexBinaryIVF*>(st->faiss_index))->nprobe = nprobe;
    (static_cast<faiss::IndexBinaryIVF*>(st->faiss_index))->search(nq, xq, k, tmpD, I);
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
    return false;
  }
  for(int i = 0; i < nq*k; i++) {
    D[i] = tmpD[i];
  }
  delete[] tmpD;

  return true;
}

int faiss_remove(
    const FaissStruct* st,
    const int size,
    const long int* ids,
    const int method_type) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  switch (method_type) {
    case IVFPQ:
      return faiss_remove_ivfpq(st, size, ids);
    case BINARYIVF:
      return faiss_remove_binaryivf(st, size, ids);
    default:
      std::stringstream ss;
      ss << "Capi : " << __FUNCTION__ << "() : Error: no method type.";
      std::cerr << ss.str() << std::endl;
      return -1;
  }
}

int faiss_remove_ivfpq(
    const FaissStruct* st,
    const int size,
    const long int* ids) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  try {
    //printf("is_trained: %d\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->is_trained);
    //printf("ntotal: %ld\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal);
    faiss::IDSelectorArray sel(size, ids);
    (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->remove_ids(sel);
    //printf("is_trained: %d\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->is_trained);
    //printf("ntotal: %ld\n", (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal);
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
    return -1;
  }

  return (static_cast<faiss::IndexIVFPQ*>(st->faiss_index))->ntotal;
}

int faiss_remove_binaryivf(
    const FaissStruct* st,
    const int size,
    const long int* ids) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  try {
    faiss::IDSelectorArray sel(size, ids);
    (static_cast<faiss::IndexBinaryIVF*>(st->faiss_index))->remove_ids(sel);
  } catch(std::exception &err) {
    std::stringstream ss;
    ss << "Capi : " << __FUNCTION__ << "() : Error: " << err.what();
    std::cerr << ss.str() << std::endl;
    return -1;
  }

  return (static_cast<faiss::IndexBinaryIVF*>(st->faiss_index))->ntotal;
}

void faiss_free(FaissStruct* st) {
  //printf(__FUNCTION__);
  //printf("\n");
  //fflush(stdout);

  free(st);
  return;
}
