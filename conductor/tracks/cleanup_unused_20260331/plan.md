# Implementation Plan: Clean up PNG writer sRGB parameter and functionalize frame debug functions

## Phase 1: Preparation and Testing Setup
- [x] Task: Create reproduction/verification tests for `writeSRGB` in `core/png_writer_test.go`
    - [x] Add a test case that checks if `writeSRGB` correctly handles different rendering intents.
    - [x] Verify that the current hardcoded implementation always produces the same output regardless of input.
- [~] Task: Create reproduction/verification tests for `displayBuffer` and `displayBuffers` in `frame/Frame_test.go`
    - [ ] Ensure tests cover the execution of these functions (already partially done in `frame/Frame_test.go`).
    - [ ] Prepare to verify that `log.Debugf` is called with the correct label and sum.

## Phase 2: PNG Writer Cleanup (`core/png_writer.go`)
- [x] Task: Refactor `writeSRGB` to use `*JXLImage` and dynamic CRC (3e26e93)
    - [x] Rename parameter `_` to `jxlImage`.
    - [x] Implement rendering intent retrieval from `jxlImage.imageHeader.ColourEncoding.RenderingIntent`.
    - [x] Implement dynamic CRC32 calculation for the `sRGB` chunk.
    - [x] Standardize the chunk writing process (Length, Type, Data, CRC).
- [x] Task: Verify `writeSRGB` changes (3e26e93)
    - [x] Run tests in `core/` and ensure success for `writeSRGB`.
    - [x] Check coverage for `core/png_writer.go`.
- [x] Task: Conductor - User Manual Verification 'PNG Writer Cleanup' (Protocol in workflow.md)

## Phase 3: Frame Debug Functions Functionalization (`frame/Frame.go`)
- [ ] Task: Update `displayBuffer` and `displayBuffers` to use all parameters and return sum
    - [ ] Rename `_` to `label`.
    - [ ] Add `log.Debugf` calls to output the `label` and calculated `total` sum.
    - [ ] Update functions to return the calculated `total` sum (float64) for verification.
    - [ ] Remove any unnecessary `// nolint` markers if the IDE no longer complains.
- [ ] Task: Verify Frame debug changes
    - [ ] Run tests in `frame/` and ensure they pass.
    - [ ] Verify that debug logs are produced.
- [ ] Task: Conductor - User Manual Verification 'Frame Debug Functionalization' (Protocol in workflow.md)

## Phase 4: Final Validation and Cleanup
- [ ] Task: Run full project test suite
- [ ] Task: Final code review and quality gate check
- [ ] Task: Conductor - User Manual Verification 'Final Validation' (Protocol in workflow.md)
