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

### Release Process

#### Option 1: Tag-Based Release (Recommended)

Push a version tag to automatically create a release with the latest Jolt Physics:

1. **Update `download.go` line 14** with the new release tag:
   ```go
   releaseTag = "v0.1.0"  // Update this
   ```

2. **Commit and push** this change:
   ```bash
   git add download.go
   git commit -m "Bump version to v0.1.0"
   git push
   ```

3. **Create and push the version tag**:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

This will:
1. Build binaries for both platforms
2. Run tests on both platforms
3. Create a GitHub Release with tag `v0.1.0`
4. Use latest Jolt Physics version
5. Upload binaries to the release

#### Option 2: Manual Release (Ad-hoc)

Use this when you need to specify a specific Jolt Physics version or rebuild an existing release:

1. Go to Actions tab on GitHub
2. Select "Build Pre-built Binaries" workflow
3. Click "Run workflow"
4. Fill in:
   - **release_tag**: e.g., `v0.1.0` (required)
   - **jolt_version**: e.g., `v5.4.0`, commit hash, or `latest` (optional, defaults to latest)
5. Wait for workflow to complete
6. A new GitHub Release will be created with binaries attached

### Build Triggers

- **Version tags** (`v*`): Builds, tests, and creates release
- **Manual workflow dispatch**: Builds, tests, and creates release with specified versions
- **Changes to `wrapper/`, `scripts/`, or workflow file**: Builds and tests only (no release)

### Release Assets

After the workflow completes, these files are uploaded to the GitHub Release:
- `darwin_arm64_libJolt.a`
- `darwin_arm64_libjolt_wrapper.a`
- `linux_amd64_libJolt.a`
- `linux_amd64_libjolt_wrapper.a`

Users automatically download these binaries on first `go get`.

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
2. Verify the release tag matches the version in `download.go` (`releaseTag = "..."`)
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
