# Specification: Refactor BitWriter to Shared Test Utility

## Overview
Currently, `frame/Quantizer_test.go` and `entropy/SymbolDistribution_test.go` contain duplicate implementations of a `BitWriter` struct used for crafting mock bitstreams in unit tests. This track aims to deduplicate this code by moving it to a shared `testcommon` package.

## Functional Requirements
- **Create Shared Utility:** Move the `BitWriter` struct and its associated methods (`WriteBit`, `WriteBits`, `WriteU8`, and `Bytes`) to a new file `testcommon/bitwriter.go`.
- **Export BitWriter:** Ensure the struct and its methods are exported (e.g., `BitWriter`, `WriteBit`, etc.) to be accessible by other packages' tests.
- **Update Frame Tests:** Modify `frame/Quantizer_test.go` to import and use `testcommon.BitWriter` instead of its local implementation.
- **Update Entropy Tests:** Modify `entropy/SymbolDistribution_test.go` to import and use `testcommon.BitWriter` instead of its local implementation.

## Non-Functional Requirements
- **Maintainability:** Reduce code duplication and improve test suite consistency.
- **Independence:** Avoid creating direct dependencies between the `frame` and `entropy` packages themselves; use `testcommon` as the bridge for test helpers.

## Acceptance Criteria
- `go test ./frame/...` passes with the refactored `Quantizer_test.go`.
- `go test ./entropy/...` passes with the refactored `SymbolDistribution_test.go`.
- No local `BitWriter` implementation remains in `frame/Quantizer_test.go` or `entropy/SymbolDistribution_test.go`.

## Out of Scope
- Refactoring the production `jxlio.BitReader`.
- Adding any non-test related functionality to `testcommon`.
