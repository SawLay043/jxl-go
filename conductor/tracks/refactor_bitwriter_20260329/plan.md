# Implementation Plan: Refactor BitWriter to Shared Test Utility

## Phase 1: Preparation & Setup [checkpoint: e68ee0b]
- [x] Task: Create the `testcommon` directory if it doesn't exist.
- [x] Task: Conductor - User Manual Verification 'Phase 1: Preparation & Setup' (Protocol in workflow.md)

## Phase 2: Create Shared Utility [checkpoint: 1ec1ac9]
- [x] Task: Create `testcommon/bitwriter.go` and implement the `BitWriter` struct with exported methods: `WriteBit`, `WriteBits`, `WriteU8`, and `Bytes`.
- [x] Task: Write unit tests for the shared `BitWriter` in `testcommon/bitwriter_test.go`.
- [x] Task: Verify shared `BitWriter` tests pass.
- [x] Task: Conductor - User Manual Verification 'Phase 2: Create Shared Utility' (Protocol in workflow.md)

## Phase 3: Refactor frame Package [checkpoint: a39d52e]
- [x] Task: Modify `frame/Quantizer_test.go` to import `testcommon` and use `testcommon.BitWriter`.
- [x] Task: Remove the local `BitWriter` definition from `frame/Quantizer_test.go`.
- [x] Task: Run `go test ./frame/...` and ensure all tests pass.
- [x] Task: Conductor - User Manual Verification 'Phase 3: Refactor frame Package' (Protocol in workflow.md)

## Phase 4: Refactor entropy Package [checkpoint: 09fda5d]
- [x] Task: Modify `entropy/SymbolDistribution_test.go` to import `testcommon` and use `testcommon.BitWriter`.
- [x] Task: Remove the local `BitWriter` definition from `entropy/SymbolDistribution_test.go`.
- [x] Task: Run `go test ./entropy/...` and ensure all tests pass.
- [x] Task: Conductor - User Manual Verification 'Phase 4: Refactor entropy Package' (Protocol in workflow.md)

## Phase 5: Final Cleanup & Verification
- [ ] Task: Perform a final search for any remaining local `BitWriter` implementations in test files.
- [ ] Task: Ensure all project tests pass: `go test ./...`.
- [ ] Task: Conductor - User Manual Verification 'Phase 5: Final Cleanup & Verification' (Protocol in workflow.md)
