# BeeBread: Industrial-Grade Automation Logic in Go

`beebread` is a comprehensive Go library that ports the highly-regarded OSCAT (Open Source Community for Automation Technology) `basic`, `building` and `network` libraries from their original IEC 61131-3 Structured Text (ST) to idiomatic Go.

## Purpose

The OSCAT libraries are a de-facto standard in the world of industrial automation and Programmable Logic Controllers (PLCs), providing thousands of robust, well-tested functions for a wide range of tasks. The primary goal of the `beebread` project is to make this invaluable collection of industrial-grade logic available to the modern Go developer.

By porting these libraries, `beebread` aims to:

*   **Bridge Worlds**: Provide a familiar and powerful toolset for PLC programmers and automation engineers transitioning to Go.
*   **Enable Modern Applications**: Allow developers to leverage proven industrial algorithms in contemporary applications, such as IoT gateways, edge computing devices, backend services for industrial monitoring, and digital twin simulations.
*   **Ensure Robustness**: Offer a collection of functions for mathematics, time/date calculations, string manipulation, and control loop logic that has been tested and refined over many years in real-world industrial environments.

Whether you are building a home automation system, a factory monitoring dashboard, or any application that can benefit from time-tested utility functions, `beebread` provides a solid foundation.

## Features

All of OSCAT BASIC 3.35, 559 functions, function blocks and types, is ported. The packages under `basic` follow the subjects of the OSCAT source:

| Package | OSCAT subject |
|---|---|
| `basic` | types, global constants, and the IEC standard conversions and string functions the port relies on |
| `basic/buffer` | buffer management |
| `basic/engineering` | automation, control, conversion, measurements, sensors, signal generators and signal processing |
| `basic/list` | list processing |
| `basic/logic` | gate logic, flip-flops, generators and memory |
| `basic/math` | mathematics, arrays, complex numbers, double precision, geometry and vectors |
| `basic/other` | event, status and error reports |
| `basic/string` | strings |
| `basic/time_date` | time and date, sun, holidays and clocks |

The port follows the conventions of [royaljelly](https://github.com/apiarytech/royaljelly), the runtime the [beedance](https://github.com/apiarytech/beedance) transpiler targets, so transpiled IEC 61131-3 programs can call it:

*   Names are OSCAT's. A name starting with an underscore, which Go does not export, has it at the end instead: `_BUFFER_CLEAR` is `BUFFER_CLEAR_`.
*   Types are royaljelly's `iec` types: `REAL` is `iec.REAL`, a 32 bit float, as on a PLC.
*   A function takes its inputs in OSCAT's order. A function block is a struct with its inputs and outputs as fields and the methods `INIT()`, which sets the initial values, and `Execute(now time.Time)`, which runs it at the scan time `now`. Blocks never read the wall clock, so they can run on simulated time.

See the package documentation of `basic` for the details.

The OSCAT sources, cleaned for beedance, are in `doc/`: `beedance_basic.st`, `beedance_building.st` and `beedance_network.st`. Pointers are written as beedance's references (`REF_TO`, `REF()`); a POU beedance cannot run, such as one that reads the bytes of a STRING through a pointer, is commented out with the reason, and still has a Go port.

`go run ./tools/oscatcov` reports what of `doc/beedance_basic.st` is ported, by kind; its test fails if anything is not.

All of OSCAT BUILDING 1.00, 57 functions and function blocks, is ported too, built on the BASIC port. Its source, cleaned for IEC 61131-3 by `go run ./tools/stclean` (see its documentation), is `doc/beedance_building.st`, which beedance reads; the Go port follows the cleaned source:

| Package | OSCAT subject |
|---|---|
| `building` | the library version |
| `building/actuators` | valves, coils, pumps and up/down drives |
| `building/electrical` | switches, push buttons, dimmers, lamps and timers |
| `building/hlk` | heating, ventilation and air conditioning: the physics of air and water, boilers, burners, heat meters and tanks |
| `building/jalousie` | blinds and roller shutters |

The test of `building` fails if a POU of the cleaned source has no port.

All of OSCAT NETWORK 1.35 for TwinCAT, 127 functions and function blocks and 34 types, is ported too, built on the BASIC port. Its source was the binary CODESYS 2.3 library `.lib`, not an ST export; `tools/stclean` read its records and cleaned them into `doc/beedance_network.st`. A cleaned source can be cleaned again, for instance to write its pointers as references:

```
go run ./tools/stclean -refs -exclude tools/stclean/testdata/network.exclude -uses doc/beedance_basic.st -in doc/beedance_network.st -out doc/beedance_network.st
```

| Package | OSCAT subject |
|---|---|
| `network` | the library version, the types and the global variables |
| `network/tcpip` | TwinCAT's `Tc2_TcpIp` sockets, emulated with Go's `net` |
| `network/ads` | TwinCAT's ADS file service `FW_AdsRdWrt`, emulated with Go's `os` |
| `network/ip` | `IP_CONTROL`, `IP_CONTROL2` and `IP_FIFO`: shared TCP and UDP connections |
| `network/inet` | DNS, HTTP, SNTP, syslog, telnet log, FTP, SMTP and MySQL clients and servers |
| `network/modbus` | Modbus TCP and UDP client and servers |
| `network/netvar` | network variables between PLCs |
| `network/file` | the file server and reading files by block |
| `network/parser` | CSV, INI and XML parsers |
| `network/encoding` | Base64, HTML, URL, IP4 and element strings |
| `network/crypto` | MD5, SHA1, RC4 and CRAM-MD5 |
| `network/logging` | the log and its viewport |
| `network/dlog` | the data logger: values to CSV, HTML and XML files, MySQL and RRD, files on to FTP and SMTP, and a cron table |
| `network/telnet` | a telnet user interface: screen, input elements and menus |
| `network/irtrans` | IRTrans infrared transceivers |
| `network/weather` | weather services and the phase of the moon |

The blocks run on a real network: their tests talk to loopback servers. SSL is not supported, as TwinCAT's `IP_CONTROL` does not. The test of `network` fails if a POU of the cleaned source has no port.

beebread requires royaljelly v0.3.0 or later, the first release series with the `iec` package it is written against.

## Community Contributions

Contributions from the community are highly encouraged and welcome! This project is a significant undertaking, and your help is invaluable in making it a complete and robust port of the OSCAT libraries.

You can contribute in several ways:
*   **Writing Tests**: The most critical need is to achieve 100% test coverage to ensure the ported logic is bug-free and behaves identically to the original.
*   **Fixing Bugs**: If you find a discrepancy between the Go implementation and the original ST code, please open an issue or submit a pull request with a fix.
*   **Porting more libraries**: OSCAT BASIC, BUILDING and NETWORK are ported; others of the archive are not.
*   **Improving Documentation**: Enhancing the documentation helps everyone.

When contributing, please strive to write clean, idiomatic Go code.

## Original OSCAT Library

The original OSCAT Basic and Building library source files in IEC 61131-3 Structured Text can be found at the official OSCAT libs archive on GitHub:

https://github.com/eclipse-oscat/oscat-libs-archive

## Example Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/apiarytech/beebread/basic/engineering"
	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/royaljelly/iec"
)

func main() {
	// A function: the date of easter sunday.
	easter := td.EASTER(2025)
	fmt.Println(time.Time(easter).Format("January 2")) // April 20

	// A function block: a low pass filter, run once per scan.
	var filter engineering.FT_PT1
	filter.INIT()
	filter.T = iec.TIME(time.Second)
	now := time.Now()
	filter.Execute(now) // the first scan starts the filter at IN = 0
	for i := 0; i < 1000; i++ {
		now = now.Add(time.Millisecond)
		filter.IN = 1
		filter.Execute(now)
	}
	fmt.Printf("%.2f\n", filter.OUT) // 0.63
}
```
