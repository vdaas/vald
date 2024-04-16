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
#include <vector>
#include "ngt-rs/src/input.h"
#include "ngt-rs/src/lib.rs.h"

Property::Property(): p() {}

void Property::set_dimension(rust::i32 dimension) {
    p.dimension = dimension;
}

rust::i32 Property::get_dimension() {
    return p.dimension;
}

void Property::set_distance_type(const DistanceType t) {
    p.distanceType = static_cast<NGT::Property::DistanceType>(t);
}

void Property::set_object_type(const ObjectType t) {
    p.objectType = static_cast<NGT::Property::ObjectType>(t);
}

void Property::set_edge_size_for_creation(rust::i16 edge_size) {
    p.edgeSizeForCreation = edge_size;
}

void Property::set_edge_size_for_search(rust::i16 edge_size) {
    p.edgeSizeForSearch = edge_size;
}

Index::Index(const rust::String& path, bool read_only): index(new NGT::Index(std::string(path), read_only)) {}

Index::Index(const rust::String& path, Property& p) {
    std::string cpath(path);
    NGT::Index::createGraphAndTree(cpath, p.p, true);
    index = new NGT::Index(cpath);
}

Index::Index(Property& p): index(new NGT::GraphAndTreeIndex(p.p)) {}

void Index::search(rust::Slice<const rust::f32> v, rust::i32 size, rust::f32 epsilon, rust::f32 radius, rust::i32 edge_size, rust::i32 *ids, rust::f32 *distances) {
    if (radius < 0.0) {
        radius = FLT_MAX;
    }
    std::vector<float> vquery(v.begin(), v.end());
    NGT::Object* query = index->allocateObject(vquery);
    NGT::SearchContainer sc(*query);
    NGT::ObjectDistances ngtresults;
    sc.setResults(&ngtresults);
    sc.setSize(size);
    sc.setRadius(radius);
    sc.setEpsilon(epsilon);
    if (edge_size != INT_MIN) {
        sc.setEdgeSize(edge_size);
    }

    index->search(sc);
    index->deleteObject(query);

    for (int i = 0; i < size; i++) {
        ids[i] = ngtresults[i].id;
        distances[i] = ngtresults[i].distance;
    }
}

void Index::linear_search(rust::Slice<const rust::f32> v, rust::i32 size, rust::i32 edge_size, rust::i32 *ids, rust::f32 *distances) {
    std::vector<float> vquery(v.begin(), v.end());
    auto query = index->allocateObject(vquery);
    NGT::SearchContainer sc(*query);
    NGT::ObjectDistances ngtresults;
    sc.setResults(&ngtresults);
    sc.setSize(size);
    if (edge_size != INT_MIN) {
        sc.setEdgeSize(edge_size);
    }

    index->linearSearch(sc);
    index->deleteObject(query);

    for (int i = 0; i < size; i++) {
        ids[i] = ngtresults[i].id;
        distances[i] = ngtresults[i].distance;
    }
}

rust::u32 Index::insert(rust::Slice<const rust::f32> v) {
    return index->insert(std::vector<float>(v.begin(), v.end()));
}

void Index::create_index(rust::u32 pool_size) {
    index->createIndex(pool_size);
}

void Index::remove(rust::u32 id) {
    index->remove(id);
}

void Index::save(const rust::String& path) {
    index->saveIndex(std::string(path));
}

rust::Slice<const rust::f32> Index::get_vector(rust::u32 id) {
    NGT::ObjectSpace& os = index->getObjectSpace();
    rust::f32* v = static_cast<rust::f32*>(os.getObject(id));
    return rust::Slice<const rust::f32>(v, os.getDimension());
}

std::unique_ptr<Property> new_property() {
    return std::make_unique<Property>();
}

std::unique_ptr<Index> open_index(const rust::String& path, bool read_only) {
    return std::make_unique<Index>(path, read_only);
}

std::unique_ptr<Index> new_index(const rust::String& path, Property& p) {
    return std::make_unique<Index>(path, p);
}

std::unique_ptr<Index> new_index_in_memory(Property& p) {
    return std::make_unique<Index>(p);
}
