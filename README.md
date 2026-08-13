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

This library is an ongoing conversion of the OSCAT `v3.3.5` libraries. The port focuses on creating idiomatic Go code that is both performant and easy to use, while maintaining the logical integrity of the original functions.

Key packages include:
*   **`basic`**: A wide array of fundamental utilities for:
    *   `time_date`: Advanced date and time calculations, including astronomical functions like sunrise/sunset, holiday calculations, and more.
    *   `math`: Mathematical functions beyond the standard library.
    *   `buffer`: Utilities for byte slice manipulation.
    *   And many more...
*   **`building`**: Functions specific to building automation tasks.

## Example Usage

Here is a simple example of how to use a function from the `time_date` package to calculate the date of Easter for a given year:

```go
package main

import (
	"fmt"
	"github.com/apiarytech/beebread/basic"
	"time"
)

func main() {
	// Calculate the date of Easter for the year 2025
	easterDate := basic.Easter(2025)
	fmt.Printf("Easter Sunday in 2025 is on: %s\n", easterDate.Format("January 2"))
	// Output: Easter Sunday in 2025 is on: April 20
}
```
