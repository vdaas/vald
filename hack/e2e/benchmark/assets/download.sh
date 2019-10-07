#!/bin/sh

if [[ ! -e $1.hdf5 ]]; then
    curl -fsSLO http://vectors.erikbern.com/$1.hdf5
fi
md5sum -c $1.md5
