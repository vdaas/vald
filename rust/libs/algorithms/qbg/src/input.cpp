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

void Property::set_hierarchical_clustering_init_mode(rust::i16 hierarchical_clustering_init_mode)
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

void Property::set_number_of_second_objects(rust::u32 number_of_second_objects)
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
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << ngt_get_error_string(err) << std::endl;
        ngt_destroy_error_object(err);
        return;
    }
    open_index(cpath.c_str(), false);
    ngt_destroy_error_object(err);
}

Index::Index(const rust::String &path, bool prebuilt)
{
    std::string cpath(path);
    open_index(cpath.c_str(), prebuilt);
}

Index::~Index()
{
    // close_index();
}

void Index::open_index(const rust::String &path, bool prebuilt)
{
    NGTError err = ngt_create_error_object();
    std::string cpath(path);
    index = qbg_open_index(cpath.c_str(), prebuilt, err);
    if (index == 0)
    {
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << ngt_get_error_string(err) << std::endl;
        ngt_destroy_error_object(err);
        return;
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
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << ngt_get_error_string(err) << std::endl;
        ngt_destroy_error_object(err);
        return;
    }
    ngt_destroy_error_object(err);
}

void Index::save_index()
{
    NGTError err = ngt_create_error_object();
    bool ok = qbg_save_index(index, err);
    if (!ok)
    {
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << ngt_get_error_string(err) << std::endl;
        ngt_destroy_error_object(err);
        return;
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
        std::cerr << ngt_get_error_string(err) << std::endl;
        ngt_destroy_error_object(err);
        return 0;
    }
    ngt_destroy_error_object(err);
    return id;
}

void Index::search(
    rust::Slice<const rust::f32> v,
    rust::usize k,
    rust::i32 *ids,
    rust::f32 *distances)
{
    QBGQuery query;
    qbg_initialize_query(&query);
    std::vector<float> vec(v.begin(), v.end());
    query.query = vec.data();

    NGTError err = ngt_create_error_object();
    NGTObjectDistances results = ngt_create_empty_results(err);
    bool ok = qbg_search_index(index, query, results, err);
    if (!ok)
    {
        std::cerr << "Error: " << __func__ << std::endl;
        std::cerr << ngt_get_error_string(err) << std::endl;
        qbg_destroy_results(results);
        ngt_destroy_error_object(err);
        return;
    }

    size_t rsize = qbg_get_result_size(results, err);
    size_t limit = std::min(k, rsize);
    for (size_t i = 0; i < limit; i++)
    {
        NGTObjectDistance obj = qbg_get_result(results, i, err);
        ids[i] = obj.id;
        distances[i] = obj.distance;
    }

    qbg_destroy_results(results);
    ngt_destroy_error_object(err);
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