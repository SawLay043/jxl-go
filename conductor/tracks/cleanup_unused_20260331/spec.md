# Track Specification: Clean up PNG writer sRGB parameter and functionalize frame debug functions

## Overview
This track focuses on cleaning up unused parameters in `core/png_writer.go` and `frame/Frame.go`, while simultaneously improving functionality and code quality. In `writeSRGB`, the unused `*JXLImage` parameter will be used to dynamically determine the rendering intent from the JXL metadata, and the hardcoded CRC will be replaced with a dynamic calculation. In `frame/Frame.go`, the unused `label` and `total` parameters in `displayBuffer` and `displayBuffers` will be put to use in debug log statements.

## Functional Requirements
- **Dynamic `writeSRGB` in `core/png_writer.go`:**
    - Rename the first parameter from `_` to `jxlImage`.
    - Use `jxlImage.imageHeader.ColourEncoding.RenderingIntent` (if available) to set the rendering intent byte.
    - Dynamically calculate the CRC for the `sRGB` chunk using the `crc32` package.
    - Ensure the chunk length, type, and data are written consistently with other `write*` methods in `PNGWriter`.
- **Functional Debug in `frame/Frame.go`:**
    - Rename the first parameter from `_` to `label` in `displayBuffer` and `displayBuffers`.
    - Calculate the `total` of the buffer values (already partially implemented).
    - Use `log.Debugf` to output the `label` and `total` sum for both functions.
    - Remove `// nolint` (if applicable and no longer needed).

## Non-Functional Requirements
- **Consistency:** Align the structure of `writeSRGB` with existing methods like `writeICCP` and `writeIHDR`.
- **Maintainability:** Eliminate unused variables and hardcoded constants (magic numbers) related to PNG chunks.

## Acceptance Criteria
- `core/png_writer.go`: `writeSRGB` successfully uses the `jxlImage` parameter and writes a correct `sRGB` chunk with a dynamic CRC.
- `frame/Frame.go`: `displayBuffer` and `displayBuffers` use both their `label` and `frameBuffer` parameters and output the result via `log.Debugf`.
- The code compiles and passes existing tests.

## Out of Scope
- Major refactoring of the `PNGWriter` class beyond the `writeSRGB` method.
- Performance optimization of the `displayBuffer` sum calculation.
