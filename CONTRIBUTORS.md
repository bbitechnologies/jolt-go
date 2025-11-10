# Contributor Guide

This guide is for contributors who need to rebuild the pre-built binaries for jolt-go.

## Overview

jolt-go ships with pre-built static libraries for:
- `darwin/arm64` (macOS Apple Silicon)
- `linux/amd64` (Linux x86-64)

These binaries are committed to the repository in the `lib/` directory, ensuring the package works immediately with `go get`.

**Current Jolt Physics Version:** v5.4.0

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
- `lib/{platform}/libJolt.a`
- `lib/{platform}/libjolt_wrapper.a`

## Testing Binaries

After building, test that everything works:

```bash
go run example/main.go
```

Expected output:
```
Player controller demo - character walks in a circle
==========================================================
[0.0s] Position: X=  0.05 Y=  5.00 Z=  0.00 | State: InAir
[0.5s] Position: X=  1.45 Y=  3.65 Z=  0.47 | State: InAir
[1.0s] Position: X=  2.30 Y=  1.90 Z=  1.67 | State: OnGround
[1.5s] Position: X=  2.29 Y=  1.90 Z=  3.15 | State: OnGround
[2.0s] Position: X=  1.41 Y=  1.90 Z=  4.33 | State: OnGround
[2.5s] Position: X=  0.00 Y=  1.90 Z=  4.77 | State: OnGround
[3.0s] Position: X= -1.40 Y=  1.90 Z=  4.30 | State: OnGround
[3.5s] Position: X= -2.25 Y=  1.90 Z=  3.10 | State: OnGround
[4.0s] Position: X= -2.24 Y=  1.90 Z=  1.63 | State: OnGround
[4.5s] Position: X= -1.36 Y=  1.90 Z=  0.44 | State: OnGround
==========================================================
Demo complete!
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

4. **Commit and tag the binaries:**
   ```bash
   git add lib/
   git commit -m "Update binaries to Jolt Physics v5.5.0"
   git tag v0.2.0
   git push origin main --tags
   ```

### Release Process

To create a new release after updating binaries:

   ```bash
   git tag v0.2.0
   git push origin main --tags
   ```
