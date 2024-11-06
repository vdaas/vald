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
#include "qbg/src/input.h"
#include "qbg/src/lib.rs.h"

Property::Property()
{
    qbg_construction_parameters = new QBGConstructionParameters();
    qbg_build_parameters = new QBGBuildParameters();
    qbg_initialize_construction_parameters(qbg_construction_parameters);
    qbg_initialize_build_parameters(qbg_build_parameters);
}

Property::~Property()
{
    delete qbg_construction_parameters;
    delete qbg_build_parameters;
}

QBGConstructionParameters *Property::get_qbg_construction_parameters()
{
    return qbg_construction_parameters;
}

void Property::init_qbg_construction_parameters()
{
    qbg_initialize_construction_parameters(qbg_construction_parameters);
}

void Property::set_qbg_construction_parameters(
    rust::usize extended_dimension,
    rust::usize dimension,
    rust::usize number_of_subvectors,
    rust::usize number_of_blobs,
    rust::i32 internal_data_type,
    rust::i32 data_type,
    rust::i32 distance_type)
{
    qbg_initialize_construction_parameters(qbg_construction_parameters);
    qbg_construction_parameters->extended_dimension = extended_dimension;
    qbg_construction_parameters->dimension = dimension;
    qbg_construction_parameters->number_of_subvectors = number_of_subvectors;
    qbg_construction_parameters->number_of_blobs = number_of_blobs;
    qbg_construction_parameters->internal_data_type = internal_data_type;
    qbg_construction_parameters->data_type = data_type;
    qbg_construction_parameters->distance_type = distance_type;
}

void Property::set_extended_dimension(rust::usize extended_dimension)
{
    qbg_construction_parameters->extended_dimension = extended_dimension;
}

void Property::set_dimension(rust::usize dimension)
{
    qbg_construction_parameters->dimension = dimension;
}

void Property::set_number_of_subvectors(rust::usize number_of_subvectors)
{
    qbg_construction_parameters->number_of_subvectors = number_of_subvectors;
}

void Property::set_number_of_blobs(rust::usize number_of_blobs)
{
    qbg_construction_parameters->number_of_blobs = number_of_blobs;
}

void Property::set_internal_data_type(rust::i32 internal_data_type)
{
    qbg_construction_parameters->internal_data_type = internal_data_type;
}

void Property::set_data_type(rust::i32 data_type)
{
    qbg_construction_parameters->data_type = data_type;
}

void Property::set_distance_type(rust::i32 distance_type)
{
    qbg_construction_parameters->distance_type = distance_type;
}

QBGBuildParameters *Property::get_qbg_build_parameters()
{
    return qbg_build_parameters;
}

void Property::init_qbg_build_parameters()
{
    qbg_initialize_build_parameters(qbg_build_parameters);
}

void Property::set_qbg_build_parameters(
    rust::i32 hierarchical_clustering_init_mode,
    rust::usize number_of_first_objects,
    rust::usize number_of_first_clusters,
    rust::usize number_of_second_objects,
    rust::usize number_of_second_clusters,
    rust::usize number_of_third_clusters,
    rust::usize number_of_objects,
    rust::usize number_of_subvectors,
    rust::i32 optimization_clustering_init_mode,
    rust::usize rotation_iteration,
    rust::usize subvector_iteration,
    rust::usize number_of_matrices,
    bool rotation,
    bool repositioning)
{
    qbg_initialize_build_parameters(qbg_build_parameters);
    qbg_build_parameters->hierarchical_clustering_init_mode = hierarchical_clustering_init_mode;
    qbg_build_parameters->number_of_first_objects = number_of_first_objects;
    qbg_build_parameters->number_of_first_clusters = number_of_first_clusters;
    qbg_build_parameters->number_of_second_objects = number_of_second_objects;
    qbg_build_parameters->number_of_second_clusters = number_of_second_clusters;
    qbg_build_parameters->number_of_third_clusters = number_of_third_clusters;
    qbg_build_parameters->number_of_objects = number_of_objects;
    qbg_build_parameters->number_of_subvectors = number_of_subvectors;
    qbg_build_parameters->optimization_clustering_init_mode = optimization_clustering_init_mode;
    qbg_build_parameters->rotation_iteration = rotation_iteration;
    qbg_build_parameters->subvector_iteration = subvector_iteration;
    qbg_build_parameters->number_of_matrices = number_of_matrices;
    qbg_build_parameters->rotation = rotation;
    qbg_build_parameters->repositioning = repositioning;
}

void Property::set_hierarchical_clustering_init_mode(rust::i32 hierarchical_clustering_init_mode)
{
    qbg_build_parameters->hierarchical_clustering_init_mode = hierarchical_clustering_init_mode;
}

void Property::set_number_of_first_objects(rust::usize number_of_first_objects)
{
    qbg_build_parameters->number_of_first_objects = number_of_first_objects;
}

void Property::set_number_of_first_clusters(rust::usize number_of_first_clusters)
{
    qbg_build_parameters->number_of_first_clusters = number_of_first_clusters;
}

void Property::set_number_of_second_objects(rust::usize number_of_second_objects)
{
    qbg_build_parameters->number_of_second_objects = number_of_second_objects;
}

void Property::set_number_of_second_clusters(rust::usize number_of_second_clusters)
{
    qbg_build_parameters->number_of_second_clusters = number_of_second_clusters;
}

void Property::set_number_of_third_clusters(rust::usize number_of_third_clusters)
{
    qbg_build_parameters->number_of_third_clusters = number_of_third_clusters;
}

void Property::set_number_of_objects(rust::usize number_of_objects)
{
    qbg_build_parameters->number_of_objects = number_of_objects;
}

void Property::set_number_of_subvectors_for_bp(rust::usize number_of_subvectors)
{
    qbg_build_parameters->number_of_subvectors = number_of_subvectors;
}

void Property::set_optimization_clustering_init_mode(rust::i32 optimization_clustering_init_mode)
{
    qbg_build_parameters->optimization_clustering_init_mode = optimization_clustering_init_mode;
}

void Property::set_rotation_iteration(rust::usize rotation_iteration)
{
    qbg_build_parameters->rotation_iteration = rotation_iteration;
}

void Property::set_subvector_iteration(rust::usize subvector_iteration)
{
    qbg_build_parameters->subvector_iteration = subvector_iteration;
}

void Property::set_number_of_matrices(rust::usize number_of_matrices)
{
    qbg_build_parameters->number_of_matrices = number_of_matrices;
}

void Property::set_rotation(bool rotation)
{
    qbg_build_parameters->rotation = rotation;
}

void Property::set_repositioning(bool repositioning)
{
    qbg_build_parameters->repositioning = repositioning;
}

Index::Index(const rust::String &path, Property &p)
{
    NGTError err = ngt_create_error_object();
    std::string cpath(path);
    bool ok = qbg_create(cpath.c_str(), p.get_qbg_construction_parameters(), err);
    if (!ok)
    {
        string s = ngt_get_error_string(err);
        ngt_destroy_error_object(err);
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << s << std::endl;
        throw std::runtime_error(s);
    }
    open_index(cpath.c_str(), false);
    ngt_destroy_error_object(err);
}

Index::Index(const rust::String &path, bool prebuilt)
{
    std::string cpath(path);
    open_index(cpath.c_str(), prebuilt);
}

Index::~Index() {}

void Index::open_index(const rust::String &path, bool prebuilt)
{
    NGTError err = ngt_create_error_object();
    std::string cpath(path);
    index = qbg_open_index(cpath.c_str(), prebuilt, err);
    if (index == 0)
    {
        string s = ngt_get_error_string(err);
        ngt_destroy_error_object(err);
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << s << std::endl;
        throw std::runtime_error(s);
    }
    ngt_destroy_error_object(err);
}

void Index::build_index(const rust::String &path, Property &p)
{
    NGTError err = ngt_create_error_object();
    std::string cpath(path);
    bool ok = qbg_build_index(cpath.c_str(), p.get_qbg_build_parameters(), err);
    if (!ok)
    {
        string s = ngt_get_error_string(err);
        ngt_destroy_error_object(err);
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << s << std::endl;
        throw std::runtime_error(s);
    }
    ngt_destroy_error_object(err);
}

void Index::save_index()
{
    NGTError err = ngt_create_error_object();
    bool ok = qbg_save_index(index, err);
    if (!ok)
    {
        string s = ngt_get_error_string(err);
        ngt_destroy_error_object(err);
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << s << std::endl;
        throw std::runtime_error(s);
    }
    ngt_destroy_error_object(err);
}

void Index::close_index()
{
    qbg_close_index(index);
}

rust::i32 Index::append(rust::Slice<const rust::f32> v)
{
    NGTError err = ngt_create_error_object();
    std::vector<float> vec(v.begin(), v.end());
    unsigned int id = qbg_append_object(index, vec.data(), v.length(), err);
    if (id == 0)
    {
        string s = ngt_get_error_string(err);
        ngt_destroy_error_object(err);
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << s << std::endl;
        throw std::runtime_error(s);
    }
    ngt_destroy_error_object(err);
    return id;
}

rust::i32 Index::insert(rust::Slice<const rust::f32> v)
{
    NGTError err = ngt_create_error_object();
    std::vector<float> vec(v.begin(), v.end());
    unsigned int id = qbg_insert_object(index, vec.data(), v.length(), err);
    if (id == 0)
    {
        string s = ngt_get_error_string(err);
        ngt_destroy_error_object(err);
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << s << std::endl;
        throw std::runtime_error(s);
    }
    ngt_destroy_error_object(err);
    return id;
}

void Index::remove(rust::usize id)
{
    NGTError err = ngt_create_error_object();
    bool ok = qbg_remove_object(index, id, err);
    if (!ok)
    {
        string s = ngt_get_error_string(err);
        ngt_destroy_error_object(err);
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << s << std::endl;
        throw std::runtime_error(s);
    }
    ngt_destroy_error_object(err);
}

std::unique_ptr<std::vector<SearchResult>> Index::search(rust::Slice<const rust::f32> v, rust::usize k, rust::f32 radius, rust::f32 epsilon)
{
    QBGQuery query;
    qbg_initialize_query(&query);
    std::vector<float> vec(v.begin(), v.end());
    query.query = vec.data();
    query.number_of_results = k;
    query.radius = radius;
    query.epsilon = epsilon;

    NGTError err = ngt_create_error_object();
    NGTObjectDistances results = ngt_create_empty_results(err);
    bool ok = qbg_search_index(index, query, results, err);
    if (!ok)
    {
        string s = ngt_get_error_string(err);
        ngt_destroy_error_object(err);
        qbg_destroy_results(results);
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << s << std::endl;
        throw std::runtime_error(s);
    }

    size_t rsize = qbg_get_result_size(results, err);
    if (rsize == 0)
    {
        string s = ngt_get_error_string(err);
        ngt_destroy_error_object(err);
        qbg_destroy_results(results);
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << s << std::endl;
        throw std::runtime_error(s);
    }
    size_t limit = std::min(k, rsize);
    std::vector<SearchResult> searchResults;
    for (size_t i = 0; i < limit; i++)
    {
        NGTObjectDistance obj = qbg_get_result(results, i, err);
        if (obj.id == 0)
        {
            string s = ngt_get_error_string(err);
            ngt_destroy_error_object(err);
            qbg_destroy_results(results);
            std::cerr << "Error: " << __func__ << std::endl;
            std::cerr << s << std::endl;
            throw std::runtime_error(s);
        }
        searchResults.push_back(SearchResult{obj.id, obj.distance});
    }
    ngt_destroy_error_object(err);
    qbg_destroy_results(results);
    return std::make_unique<std::vector<SearchResult>>(searchResults);
}

rust::f32 *Index::get_object(rust::usize id)
{
    NGTError err = ngt_create_error_object();
    float *vec = qbg_get_object(index, id, err);
    if (vec == 0)
    {
        string s = ngt_get_error_string(err);
        ngt_destroy_error_object(err);
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << s << std::endl;
        throw std::runtime_error(s);
    }
    ngt_destroy_error_object(err);
    return vec;
}

rust::usize Index::get_dimension()
{
    NGTError err = ngt_create_error_object();
    size_t dim = qbg_get_dimension(index, err);
    if (dim == 0)
    {
        string s = ngt_get_error_string(err);
        ngt_destroy_error_object(err);
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << s << std::endl;
        throw std::runtime_error(s);
    }
    ngt_destroy_error_object(err);
    return dim;
}

std::unique_ptr<Property> new_property()
{
    return std::make_unique<Property>();
}

std::unique_ptr<Index> new_index(const rust::String &path, Property &p)
{
    return std::make_unique<Index>(path, p);
}

std::unique_ptr<Index> new_prebuilt_index(const rust::String &path, bool prebuilt)
{
    return std::make_unique<Index>(path, prebuilt);
}