# Maintainer Guide

This guide is for maintainers who need to rebuild the pre-built binaries for jolt-go.

## Overview

jolt-go ships with pre-built static libraries for:
- `darwin/arm64` (macOS Apple Silicon)
- `linux/amd64` (Linux x86-64)

These binaries are committed to the repository in the `lib/` directory, ensuring the package works immediately with `go get`.

**Current Jolt Physics Version:** v5.4.0 (commit: be9036cf)

## Prerequisites

### Required Tools
- **Go 1.20+**
- **Git**
- **CMake 3.15+**
- **Docker** (for Linux builds)
- **macOS with Apple Silicon** (for darwin/arm64 builds)

## Setup

### 1. Clone/Update Jolt Physics

Clone Jolt Physics to a location of your choice:

```bash
git clone https://github.com/jrouwe/JoltPhysics.git ~/JoltPhysics
cd ~/JoltPhysics

# Optionally check out a specific version
git checkout v5.4.0
# Or use latest main
git pull origin main
```

### 2. Set JOLT_SRC Environment Variable

```bash
export JOLT_SRC=~/JoltPhysics  # or wherever you cloned it
```

**Tip:** Add this to your `.bashrc`, `.zshrc`, or `.profile` to make it permanent:

```bash
echo 'export JOLT_SRC=~/JoltPhysics' >> ~/.zshrc
```

## Rebuilding Binaries

### Build All Platforms

```bash
cd /path/to/jolt-go
export JOLT_SRC=/path/to/JoltPhysics
./scripts/build-libs.sh all
```

### Build Specific Platform

```bash
# Build darwin/arm64 only (requires macOS ARM64)
./scripts/build-libs.sh darwin_arm64

# Build linux/amd64 only (uses Docker)
./scripts/build-libs.sh linux_amd64
```

### Build Output

Binaries will be placed in:
- `lib/darwin_arm64/libJolt.a` (~68 MB)
- `lib/darwin_arm64/libjolt_wrapper.a` (~184 KB)
- `lib/linux_amd64/libJolt.a` (~76 MB)
- `lib/linux_amd64/libjolt_wrapper.a` (~179 KB)

## Testing Binaries

After building, test that everything works:

```bash
go run example/main.go
```

Expected output:
```
Simulating falling sphere...
Frame   0: Y = 20.00
Frame  30: Y = 18.66
Frame  60: Y = 14.94
Frame  90: Y = 8.88
Frame 120: Y = 0.56
Frame 150: Y = -9.98
Done!
```

## Updating Jolt Physics Version

When updating to a new Jolt Physics version:

1. **Update Jolt source:**
   ```bash
   cd $JOLT_SRC
   git pull origin main
   # Or checkout specific tag: git checkout v5.5.0
   ```

2. **Rebuild all binaries:**
   ```bash
   cd /path/to/jolt-go
   ./scripts/build-libs.sh all
   ```

3. **Test:**
   ```bash
   go run example/main.go
   go test ./...  # if you have tests
   ```

4. **Update this file** with the new version number at the top

5. **Commit and tag the binaries:**
   ```bash
   git add lib/
   git commit -m "Update binaries to Jolt Physics v5.5.0"
   git tag v0.2.0
   git push origin main --tags
   ```

## CI/CD

### GitHub Actions

The repository has a GitHub Actions workflow (`.github/workflows/build-binaries.yml`) that automatically builds and tests binaries on both platforms.

### Release Process

1. **Rebuild binaries** (if updating Jolt Physics version):
   ```bash
   ./scripts/build-libs.sh all
   go run example/main.go  # Test locally
   ```

2. **Commit the binaries:**
   ```bash
   git add lib/
   git commit -m "Update to Jolt Physics v5.5.0"
   ```

3. **Create and push a version tag:**
   ```bash
   git tag v0.2.0
   git push origin main --tags
   ```

The GitHub Actions workflow will automatically build and test on both platforms when you push a tag.

## Troubleshooting

### "JOLT_SRC is not set"

Make sure you've exported the `JOLT_SRC` variable:
```bash
export JOLT_SRC=/path/to/JoltPhysics
```

### "CMakeLists.txt not found"

Make sure your `JOLT_SRC` points to the root of the JoltPhysics repository (not a subdirectory).

### "Docker not found"

Install Docker Desktop (macOS) or Docker Engine (Linux):
```bash
# macOS
brew install --cask docker

# Ubuntu
apt-get install docker.io
```

### Linux build fails with "permission denied"

Make sure Docker is running and you have permissions:
```bash
# Add yourself to docker group (Linux)
sudo usermod -aG docker $USER
# Then logout and login again
```


## Architecture Notes

### Build System

- **darwin/arm64**: Built natively on macOS using Clang
- **linux/amd64**: Built in Docker container using GCC

### Compiler Flags

**Critical:** These flags MUST match between Jolt build and wrapper build:
- `-DJPH_DISABLE_CUSTOM_ALLOCATOR` (uses standard allocator)
- `-DJPH_PROFILE_ENABLED` (profiling support)
- `-DJPH_DEBUG_RENDERER` (debug rendering)
- `-DJPH_OBJECT_STREAM` (serialization support)

Mismatch will cause `RegisterTypes()` to abort with version errors.

### Wrapper Layer

The C wrapper (`wrapper/jolt_wrapper.cpp`) provides a C interface to Jolt's C++ API:
- Uses `extern "C"` to prevent name mangling
- Uses opaque pointers (`void*`) to hide C++ types from Go
- Compiled into `libjolt_wrapper.a` static library

## Support

Questions or issues with the build system? Open an issue on GitHub.
