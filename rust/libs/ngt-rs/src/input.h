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
    void set_object_type(const ObjectType);
    void set_distance_type(const DistanceType);
};

class Index {
    NGT::Index* index;
public:
    Index(Property&);
    rust::u32 insert(const rust::Vec<rust::f32>& object);
    void create_index(rust::u32);
    void remove(rust::u32);
};

std::unique_ptr<Property> new_property();
std::unique_ptr<Index> new_index(Property& p);