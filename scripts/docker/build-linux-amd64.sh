#!/bin/bash
set -e

echo 'Building Jolt Physics for linux/amd64...'

# Build Jolt Physics
mkdir -p /build/JoltPhysics/Build/linux_amd64_release
cd /build/JoltPhysics/Build/linux_amd64_release

cmake .. \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_CXX_COMPILER=g++ \
    -DCMAKE_CXX_FLAGS="-fno-lto" \
    -DCMAKE_INTERPROCEDURAL_OPTIMIZATION=OFF \
    -DDISABLE_CUSTOM_ALLOCATOR=ON \
    -DTARGET_UNIT_TESTS=OFF \
    -DTARGET_HELLO_WORLD=OFF \
    -DTARGET_PERFORMANCE_TEST=OFF \
    -DTARGET_SAMPLES=OFF \
    -DTARGET_VIEWER=OFF

cmake --build . -j$(nproc)

echo 'Building Jolt wrapper for linux/amd64...'

# Build wrapper
cd /build/wrapper

g++ -std=c++17 \
    -I/build/JoltPhysics \
    -DJPH_DISABLE_CUSTOM_ALLOCATOR \
    -DJPH_OBJECT_STREAM \
    -fPIC \
    -fno-lto \
    -c jolt_wrapper.cpp \
    -o jolt_wrapper.o

ar rcs libjolt_wrapper.a jolt_wrapper.o

echo 'Copying binaries to output...'

# Copy binaries
mkdir -p /build/output
cp libjolt_wrapper.a /build/output/
cp /build/JoltPhysics/Build/linux_amd64_release/libJolt.a /build/output/

echo '✅ Linux amd64 binaries built successfully!'
ls -lh /build/output/
