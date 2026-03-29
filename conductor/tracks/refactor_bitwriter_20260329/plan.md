# Implementation Plan: Refactor BitWriter to Shared Test Utility

## Phase 1: Preparation & Setup
- [ ] Task: Create the `testcommon` directory if it doesn't exist.
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Preparation & Setup' (Protocol in workflow.md)

## Phase 2: Create Shared Utility
- [ ] Task: Create `testcommon/bitwriter.go` and implement the `BitWriter` struct with exported methods: `WriteBit`, `WriteBits`, `WriteU8`, and `Bytes`.
- [ ] Task: Write unit tests for the shared `BitWriter` in `testcommon/bitwriter_test.go`.
- [ ] Task: Verify shared `BitWriter` tests pass.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Create Shared Utility' (Protocol in workflow.md)

## Phase 3: Refactor frame Package
- [ ] Task: Modify `frame/Quantizer_test.go` to import `testcommon` and use `testcommon.BitWriter`.
- [ ] Task: Remove the local `BitWriter` definition from `frame/Quantizer_test.go`.
- [ ] Task: Run `go test ./frame/...` and ensure all tests pass.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Refactor frame Package' (Protocol in workflow.md)

## Phase 4: Refactor entropy Package
- [ ] Task: Modify `entropy/SymbolDistribution_test.go` to import `testcommon` and use `testcommon.BitWriter`.
- [ ] Task: Remove the local `BitWriter` definition from `entropy/SymbolDistribution_test.go`.
- [ ] Task: Run `go test ./entropy/...` and ensure all tests pass.
- [ ] Task: Conductor - User Manual Verification 'Phase 4: Refactor entropy Package' (Protocol in workflow.md)

## Phase 5: Final Cleanup & Verification
- [ ] Task: Perform a final search for any remaining local `BitWriter` implementations in test files.
- [ ] Task: Ensure all project tests pass: `go test ./...`.
- [ ] Task: Conductor - User Manual Verification 'Phase 5: Final Cleanup & Verification' (Protocol in workflow.md)
