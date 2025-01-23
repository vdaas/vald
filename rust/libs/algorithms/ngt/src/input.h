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

#include "NGT/Index.h"
#include <memory>
#include "rust/cxx.h"

enum class DistanceType;
enum class ObjectType;

class Property {
public:
    NGT::Property p;
    Property();
    void set_dimension(rust::i32);
    rust::i32 get_dimension();
    void set_object_type(const ObjectType);
    void set_distance_type(const DistanceType);
    void set_edge_size_for_creation(rust::i16);
    void set_edge_size_for_search(rust::i16);
};

class Index {
    NGT::Index* index;
public:
    Index(const rust::String&, bool);
    Index(const rust::String&, Property&);
    Index(Property&);
    void search(rust::Slice<const rust::f32>, rust::i32, rust::f32, rust::f32, rust::i32, rust::i32*, rust::f32*);
    void linear_search(rust::Slice<const rust::f32>, rust::i32, rust::i32, rust::i32*, rust::f32*);
    rust::u32 insert(rust::Slice<const rust::f32>);
    void create_index(rust::u32);
    void remove(rust::u32);
    void save(const rust::String&);
    rust::Slice<const rust::f32> get_vector(rust::u32);
};

std::unique_ptr<Property> new_property();
std::unique_ptr<Index> open_index(const rust::String&, bool);
std::unique_ptr<Index> new_index(const rust::String&, Property&);
std::unique_ptr<Index> new_index_in_memory(Property&);
