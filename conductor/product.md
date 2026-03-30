# Initial Concept
jxl-go is a Go library for reading JPEG XL (JXL) images, with a current focus on lossless decoding and plans for future expansion and performance optimization.

# Product Guide: jxl-go

## Overview
**jxl-go** is a high-performance Go library specifically designed for reading and decoding JPEG XL (JXL) images. This project aims to provide a Go-native implementation of the JPEG XL specification, catering to Go developers, performance-focused system integrators, and enthusiasts who want a reliable tool for handling JXL data.

## Target Audience
- **Go Developers:** Those who need to seamlessly integrate JPEG XL support into their Go-based applications.
- **Performance-Focused Developers:** Users seeking a native Go solution that minimizes the overhead of decoding and emphasizes efficient resource usage.
- **JXL Enthusiasts/Researchers:** Individuals interested in the JXL format and its Go implementation, potentially contributing to or learning from the codebase.

## Primary Goals
- **Spec Compliance/Correctness:** Ensuring that the decoder strictly adheres to the JPEG XL specification to provide accurate and predictable results.
- **Efficiency and Speed:** Implementing decoding processes that are both fast and resource-efficient, including future optimizations like parallelism and SIMD.
- **Ease of Use/Integration:** Providing a simple and intuitive API that allows developers to easily read and process JXL images with minimal boilerplate.

## Core Features
- **Lossless Decoding (Current):** A primary focus on correctly decoding lossless JXL images to maintain the highest quality possible.
- **Lossy Decoding (Future):** Expanding the library's capabilities to include support for lossy image decoding, as outlined in the project roadmap.
- **Performance Optimizations (Current):** Actively optimized for memory efficiency and parallelized decoding of LF groups and upsampling filters.

## Scope of Integration
- **Standard Image Library Integration:** Enabling easy conversion of decoded JXL data into Go's standard `image.Image` format for broad compatibility.
- **Advanced Color/Metadata Support:** Providing comprehensive handling for different color spaces and metadata within the JPEG XL format.
- **CLI Tools Development:** Creating robust command-line tools for image conversion and format verification, making the library useful for both developers and end-users.
