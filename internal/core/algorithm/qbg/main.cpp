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
#
# Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

int
main(int argc, char **argv)
{
  std::string   indexPath =     "index";        // Index.

  QBGError err = ngt_create_error_object();

  NGTIndex index;

  {
    /// just create an empty index.
    std::cerr << "create" << std::endl;
    QBGConstructionParameters constructionParameters;
    qbg_initialize_construction_parameters(&constructionParameters);
    constructionParameters.dimension = 3;
    constructionParameters.extended_dimension = 16;
    constructionParameters.number_of_subvectors = 1;
    constructionParameters.number_of_blobs = 0;
    if (!qbg_create(indexPath.c_str(), &constructionParameters, err)) {
      std::cerr << "Cannot create" << std::endl;
      std::cerr << ngt_get_error_string(err) << std::endl;
    }
  }

  index = qbg_open_index(indexPath.c_str(), err);

  const size_t numberOfObjects = 64;
  for (int i = 0; i < 64; i++) {
    float vector[3] = {1.0 + i, 2.0 + i, 3.0 + 1};
    if (qbg_append_object(index, vector, 3, err) == 0) {
      std::cerr << ngt_get_error_string(err) << std::endl;
      exit(1);
    }
  }

  qbg_save_index(index, err);
  qbg_close_index(index);
  {
    std::cerr << "build..." << std::endl;
    QBGBuildParameters buildParameters;
    qbg_initialize_build_parameters(&buildParameters);
    buildParameters.number_of_objects = 0;              /// optimizer
    auto status = qbg_build_index(indexPath.c_str(), &buildParameters, err);
    if (!status) {
      std::cerr << "Cannot build. " << ngt_get_error_string(err) << std::endl;
      exit(1);
    }
    std::cerr << "end of build" << std::endl;
  }

  index = qbg_open_index(indexPath.c_str(), err);
  if (index == 0) {
    std::cerr << "Cannot open. " << ngt_get_error_string(err) << std::endl;
    exit(1);
  }

  QBGQuery query;
  qbg_initialize_query(&query);
  float queryVector[3] = {1.0, 2.0, 3.0};
  query.query = &queryVector[0];
  NGTObjectDistances results = ngt_create_empty_results(err);

  std::cerr << "search..." << std::endl;
  auto status = qbg_search_index(index, query, results, err);
  std::cerr << "end of search" << std::endl;
  if (!status) {
    std::cerr << "Cannot search. " << ngt_get_error_string(err) << std::endl;
    exit(1);
  }

  {
    auto rsize = ngt_get_result_size(results, err);
    // output resultant objects.
    std::cout << "# of results=" << rsize << std::endl;
    std::cout << "Rank\tID\tDistance" << std::endl;
    for (size_t i = 0; i < rsize; i++) {
      NGTObjectDistance object = ngt_get_result(results, i, err);
      std::cout << i + 1 << "\t" << object.id << "\t" << object.distance << std::endl;
    }

    ngt_destroy_results(results);
  }

  qbg_close_index(index);

}
