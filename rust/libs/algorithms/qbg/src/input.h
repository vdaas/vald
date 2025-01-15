//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
#pragma once

#include <memory>
#include "NGT/NGTQ/Capi.h"
#include "NGT/NGTQ/QuantizedGraph.h"
#include "rust/cxx.h"

struct SearchResult
{
    rust::u32 id;
    rust::f32 distance;

    rust::u32 get_id() { return id; }
    rust::f32 get_distance() { return distance; }
};

class Property
{
    QBGConstructionParameters *qbg_construction_parameters;
    QBGBuildParameters *qbg_build_parameters;

public:
    Property();
    ~Property();
    QBGConstructionParameters *get_qbg_construction_parameters();
    void init_qbg_construction_parameters();
    void set_qbg_construction_parameters(
        rust::usize,
        rust::usize,
        rust::usize,
        rust::usize,
        rust::i32,
        rust::i32,
        rust::i32);
    void set_extended_dimension(rust::usize);
    void set_dimension(rust::usize);
    void set_number_of_subvectors(rust::usize);
    void set_number_of_blobs(rust::usize);
    void set_internal_data_type(rust::i32);
    void set_data_type(rust::i32);
    void set_distance_type(rust::i32);
    QBGBuildParameters *get_qbg_build_parameters();
    void init_qbg_build_parameters();
    void set_qbg_build_parameters(
        // hierarchical kmeans
        rust::i32,
        rust::usize,
        rust::usize,
        rust::usize,
        rust::usize,
        rust::usize,
        // optimization
        rust::usize,
        rust::usize,
        rust::i32,
        rust::usize,
        rust::usize,
        rust::usize,
        bool,
        bool);
    void set_hierarchical_clustering_init_mode(rust::i32);
    void set_number_of_first_objects(rust::usize);
    void set_number_of_first_clusters(rust::usize);
    void set_number_of_second_objects(rust::usize);
    void set_number_of_second_clusters(rust::usize);
    void set_number_of_third_clusters(rust::usize);
    void set_number_of_objects(rust::usize);
    void set_number_of_subvectors_for_bp(rust::usize);
    void set_optimization_clustering_init_mode(rust::i32);
    void set_rotation_iteration(rust::usize);
    void set_subvector_iteration(rust::usize);
    void set_number_of_matrices(rust::usize);
    void set_rotation(bool);
    void set_repositioning(bool);
};

class Index
{
    void *index;

public:
    Index(
        const rust::String &,
        Property &);
    Index(
        const rust::String &,
        bool);
    ~Index();
    void open_index(const rust::String &, bool);
    void build_index(const rust::String &, Property &);
    void save_index();
    void close_index();
    rust::i32 append(rust::Slice<const rust::f32>);
    rust::i32 insert(rust::Slice<const rust::f32>);
    void remove(rust::usize);
    std::unique_ptr<std::vector<SearchResult>> search(rust::Slice<const rust::f32>, rust::usize, rust::f32, rust::f32);
    rust::f32 *get_object(rust::usize);
    rust::usize get_dimension();
};

std::unique_ptr<Property> new_property();
std::unique_ptr<Index> new_index(const rust::String &, Property &);
std::unique_ptr<Index> new_prebuilt_index(const rust::String &, bool);
