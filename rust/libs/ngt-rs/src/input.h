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
