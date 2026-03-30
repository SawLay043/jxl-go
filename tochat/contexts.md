# PNG Writer sRGB Parameter Cleanup

## Context
In `core\png_writer.go`, the `writeSRGB` method currently takes an unused first parameter of type `*JXLImage` (renamed to `_`).

```go
func (w *PNGWriter) writeSRGB(_ *JXLImage, output io.Writer) error {
```

## Maintainer Hint
The maintainer mentioned: *"If we're cleaning up, might as well remove the first param completely."*

## Analysis
- **Safety:** It is technically safe to remove since it's an unexported method only used within the `core` package (specifically in `WritePNG`).
- **Functionality:** Currently, `writeSRGB` writes a hardcoded rendering intent (`0x01`).
- **Future Use Case:** `JXLImage` contains a `ColourEncodingBundle` which includes a `RenderingIntent`. If `writeSRGB` is ever updated to use the actual intent from the source JXL, the `*JXLImage` parameter will be required.
- **Consistency:** Other `write*` methods in `PNGWriter` (e.g., `writeIHDR`, `writeIDAT`, `writeICCP`) all accept `*JXLImage` as their first parameter.


## Recommendation for Future
- **Option 1 (Minimalism):** Remove the parameter to satisfy the "cleanup" request.
- **Option 2 (Functional Improvement):** Keep the parameter but rename it back to `jxlImage` and use `jxlImage.imageHeader.ColourEncoding.RenderingIntent` to write the correct sRGB chunk.
- **Option 3 (Status Quo):** Leave it as `_` for now to maintain interface consistency with other chunk writers while acknowledging it's currently unused.

---

# BitWriter Duplication (frame/Quantizer_test.go vs entropy/SymbolDistribution_test.go)

## Context
PR asks: *"Can this be merged/replaced with the one in SymbolDistribution_test.go ?"* referring to the `BitWriter` struct in `frame\Quantizer_test.go`.

## Analysis
- **Duplication:** Both `frame\Quantizer_test.go` and `entropy\SymbolDistribution_test.go` contain an identical `BitWriter` struct (~20 lines of code) to craft bitstreams for testing.
- **Dependency:** They are currently independent (no dependency between `frame_test` and `entropy_test`).
- **Relation to Quantizer.go:** The `BitWriter` is a **test utility** only; it has no relation to the production logic of `Quantizer.go`.
- **Merging Challenges:** In Go, you cannot easily "merge" them without:
    1. Making one test package depend on another (e.g., `frame_test` importing `entropy_test`), which is bad practice.
    2. Exporting `BitWriter` from a shared utility package (like `testcommon`).

## Recommendation
- **No (Isolation):** Keep them separate to maintain test file isolation and naming consistency. A 20-line duplication is often preferable to an unnatural cross-package test dependency.
- **Yes (Shared Utility):** If the maintainer insists on zero duplication, the correct path is to move `BitWriter` to `testcommon/bitwriter.go` and export it as `BitWriter`. This would allow any test in the project to use it without adding inter-package dependencies.

