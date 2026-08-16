# BeeBread: Industrial-Grade Automation Logic in Go

`beebread` is a comprehensive Go library that ports the highly-regarded OSCAT (Open Source Community for Automation Technology) `basic` and `building` libraries from their original IEC 61131-3 Structured Text (ST) to idiomatic Go.

## Purpose

The OSCAT libraries are a de-facto standard in the world of industrial automation and Programmable Logic Controllers (PLCs), providing thousands of robust, well-tested functions for a wide range of tasks. The primary goal of the `beebread` project is to make this invaluable collection of industrial-grade logic available to the modern Go developer.

By porting these libraries, `beebread` aims to:

*   **Bridge Worlds**: Provide a familiar and powerful toolset for PLC programmers and automation engineers transitioning to Go.
*   **Enable Modern Applications**: Allow developers to leverage proven industrial algorithms in contemporary applications, such as IoT gateways, edge computing devices, backend services for industrial monitoring, and digital twin simulations.
*   **Ensure Robustness**: Offer a collection of functions for mathematics, time/date calculations, string manipulation, and control loop logic that has been tested and refined over many years in real-world industrial environments.

Whether you are building a home automation system, a factory monitoring dashboard, or any application that can benefit from time-tested utility functions, `beebread` provides a solid foundation.

## Features

This library is an ongoing conversion of the OSCAT libraries. The port focuses on creating idiomatic Go code that is both performant and easy to use, while maintaining the logical integrity of the original functions.

Key packages include:
*   **`basic`**: A wide array of fundamental utilities for:
    *   `time_date`: Advanced date and time calculations, including astronomical functions like sunrise/sunset, holiday calculations, and more.
    *   `math`: Mathematical functions beyond the standard library.
    *   `buffer`: Utilities for byte slice manipulation.
    *   And many more...
*   **`building`**: Functions specific to building automation tasks.

## Community Contributions

Contributions from the community are highly encouraged and welcome! This project is a significant undertaking, and your help is invaluable in making it a complete and robust port of the OSCAT libraries.

You can contribute in several ways:
*   **Writing Tests**: The most critical need is to achieve 100% test coverage to ensure the ported logic is bug-free and behaves identically to the original.
*   **Fixing Bugs**: If you find a discrepancy between the Go implementation and the original ST code, please open an issue or submit a pull request with a fix.
*   **Porting New Functions**: There are still many functions in the OSCAT libraries waiting to be ported. Feel free to pick one and submit it.
*   **Improving Documentation**: Enhancing the documentation helps everyone.

When contributing, please strive to write clean, idiomatic Go code.

## Example Usage

## Original OSCAT Library

The original OSCAT Basic and Building library source files in IEC 61131-3 Structured Text can be found at the official OSCAT libs archive on GitHub:

https://github.com/eclipse-oscat/oscat-libs-archive

Here is a simple example of how to use a function from the `time_date` package to calculate the date of Easter for a given year:

```go
package main

import (
	"fmt"
	"github.com/apiarytech/beebread/basic/time_date"
	"time"
)

func main() {
	// Calculate the date of Easter for the year 2025
	easterDate := time_date.Easter(2025)
	fmt.Printf("Easter Sunday in 2025 is on: %s\n", easterDate.Format("January 2"))
	// Output: Easter Sunday in 2025 is on: April 20
}
```
