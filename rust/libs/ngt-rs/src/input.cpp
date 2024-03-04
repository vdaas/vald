#include "ngt-rs/src/input.h"
#include "ngt-rs/src/lib.rs.h"

Property::Property(): p() {}

void Property::set_dimension(rust::i32 dimension) {
    p.dimension = dimension;
}

void Property::set_distance_type(const DistanceType t) {
    p.distanceType = static_cast<NGT::Property::DistanceType>(t);
}

void Property::set_object_type(const ObjectType t) {
    p.objectType = static_cast<NGT::Property::ObjectType>(t);
}

Index::Index(Property& p) {
    index = new NGT::GraphAndTreeIndex(p.p);
}

rust::u32 Index::insert(const rust::Vec<rust::f32>& object) {
    std::vector<float> v;
    std::copy(object.begin(), object.end(), std::back_inserter(v));
    return index->insert(v);
}

void Index::create_index(rust::u32 pool_size) {
    index->createIndex(pool_size);
}

void Index::remove(rust::u32 id) {
    index->remove(id);
}

std::unique_ptr<Property> new_property() {
    return std::make_unique<Property>();
}

std::unique_ptr<Index> new_index(Property& p) {
    return std::make_unique<Index>(p);
}