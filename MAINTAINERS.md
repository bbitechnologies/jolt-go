# Maintainer Guide

This guide is for maintainers who need to rebuild the pre-built binaries for jolt-go.

## Overview

jolt-go ships with pre-built static libraries for:
- `darwin/arm64` (macOS Apple Silicon)
- `linux/amd64` (Linux x86-64)

These binaries are hosted on GitHub Releases and automatically downloaded when users import the package.

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

5. **Create a GitHub Release:**
   - Use the GitHub Actions workflow (recommended - see CI/CD section below)
   - Or manually create a release and upload the binaries

   The binaries will be automatically downloaded by users on first import.

## CI/CD

### GitHub Actions

The repository has a GitHub Actions workflow (`.github/workflows/build-binaries.yml`) that automatically builds binaries and creates GitHub releases.

**Manual Trigger (Recommended for Releases):**
1. Go to Actions tab on GitHub
2. Select "Build Pre-built Binaries" workflow
3. Click "Run workflow"
4. Optionally specify a Jolt version/commit
5. Wait for workflow to complete
6. A new GitHub Release will be created with binaries attached

**Automatic Trigger:**
- Workflow runs automatically when files in `wrapper/` or `scripts/` are modified (builds only, no release)

### Release Process

After the workflow completes:
1. **Binaries are built and tested** across both platforms
2. **GitHub Release is created** (if triggered manually via workflow_dispatch)
3. **Binaries are uploaded** as release assets with platform-specific names:
   - `darwin_arm64_libJolt.a`
   - `darwin_arm64_libjolt_wrapper.a`
   - `linux_amd64_libJolt.a`
   - `linux_amd64_libjolt_wrapper.a`
4. **Users automatically download** these binaries on first `go get`

**Note:** Update the version in `.github/workflows/build-binaries.yml` and `download.go` when creating a new release.

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

### Binary download fails for users

If users report download failures:
1. Check that the GitHub Release exists and has all 4 binary files
2. Verify the release tag matches the version in `download.go` (currently `v0.2.0`)
3. Check that file names match the expected format: `{platform}_{filename}`
4. Users can manually download binaries from GitHub Releases and place them in `lib/{platform}/`

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
