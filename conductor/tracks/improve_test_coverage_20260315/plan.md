# Track: Improve Test Coverage and Performance Optimization - Implementation Plan

## Phase 1: Test Coverage Baseline and Benchmarking [checkpoint: 074ba6d]
- [x] Task: Establish Test Coverage Baseline
    - [x] Generate current test coverage report.
    - [x] Identify critical areas with low coverage.
    Baseline Coverage:
    - bundle: 60.9%
    - colour: 69.5%
    - core: 41.0%
    - entropy: 28.6%
    - frame: 65.9%
    - image: 65.9%
    - jxlio: 64.2%
    - util: 88.3%
    - overall core decoding packages (~55%)
- [x] Task: Create Performance Benchmarks
    - [x] Develop benchmark tests for decoding representative JXL images.
    - [x] Record baseline performance metrics.
    Baseline Performance (on i5-12600KF):
    - BenchmarkDecodeUnittest: ~644 ms/op, ~1006 MB/op, ~745k allocs/op
    - BenchmarkDecodeTiny2: ~50 ms/op, ~7 MB/op, ~103k allocs/op
    - BenchmarkDecodeLossless: ~640 ms/op, ~1006 MB/op, ~745k allocs/op
    - BenchmarkDecodeGrayscale: ~59 ms/op, ~16 MB/op, ~143k allocs/op
- [x] Task: Conductor - User Manual Verification 'Phase 1: Test Coverage Baseline and Benchmarking' (Protocol in workflow.md)

## Phase 2: Targeted Testing Improvements
- [ ] Task: Increase Coverage for `entropy` Package
    - [ ] Write unit tests for ANS and Prefix symbol distributions.
- [ ] Task: Increase Coverage for `frame` Package
    - [ ] Write unit tests for frame headers and block context.
- [ ] Task: Implement Integration Tests
    - [ ] Add tests that decode various JXL images and verify output against reference data.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Targeted Testing Improvements' (Protocol in workflow.md)

## Phase 3: Performance Analysis and Optimization
- [ ] Task: Perform Profiling
    - [ ] Use profiling tools to identify performance bottlenecks during decoding.
- [ ] Task: Optimize Memory Allocation
    - [ ] Implement optimizations for frequent memory allocations in hot paths.
- [ ] Task: Explore Parallelism/SIMD
    - [ ] Investigate and prototype parallelism for independent decoding tasks.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Performance Analysis and Optimization' (Protocol in workflow.md)
