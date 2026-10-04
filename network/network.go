/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - EPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

// Package network is the port of the OSCAT NETWORK 1.35 library for
// Beckhoff TwinCAT, built on the port of OSCAT BASIC, from
// doc/beedance_network.st, which tools/stclean read from the library
// beckoff_network_135.st. This package has the library's types, global
// variables and NETWORK_VERSION; its POUs are in a package for each
// subject.
//
// The ports follow beebread's conventions: royaljelly's iec types, OSCAT's
// names in upper case, a function block as a struct with INIT and
// Execute(now), and a VAR_IN_OUT as a pointer field. An array OSCAT
// declares from 1 is a Go array from 0: the comment of its field gives its
// bounds. A POINTER TO a buffer, with ADR, is a slice.
package network

import (
	. "github.com/apiarytech/beebread/basic"
	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/royaljelly/iec"
)

// NETWORK_VERSION returns the version of the library, 135 for 1.35, or its
// release date as DATE_TO_DWORD gives it if IN is true.
func NETWORK_VERSION(IN iec.BOOL) iec.DWORD {
	if IN {
		return DATE_TO_DWORD(td.SET_DATE(2016, 7, 1))
	}
	return 135
}

// The library's constants.
const (
	// NETWORK_BUFFER_LONG_SIZE is the last index of a NETWORK_BUFFER.
	NETWORK_BUFFER_LONG_SIZE iec.UINT = 4095
	// NETWORK_BUFFER_SHORT_SIZE is the last index of a NETWORK_BUFFER_SHORT.
	NETWORK_BUFFER_SHORT_SIZE iec.UINT = 1407
	// LOG_MAX is the last index of the messages of a LOG_CONTROL.
	LOG_MAX iec.INT = 40
	// LOG_SIZE is the length of a log message.
	LOG_SIZE iec.INT = 80
	// ELEMENT_LENGTH is the length of an element of a list.
	ELEMENT_LENGTH iec.INT = 250
)

// The library's global variables.
var (
	// TCP_SERVER_RESET is for the library's own use (TwinCAT only).
	TCP_SERVER_RESET iec.BYTE
	// SSRVNETID is the network address of the TwinCAT TCP/IP connection
	// server, '' for the local computer. OSCAT names it sSrvNetId.
	SSRVNETID iec.STRING
	// SLOCALHOST is the local IP address of a UDP socket, '' for the
	// default network adapter. OSCAT names it sLocalHost.
	SLOCALHOST iec.STRING
	// SYSLIBSOCKETS_OPTION configures CODESYS's SysLibSocket.lib; TwinCAT
	// does not use it.
	SYSLIBSOCKETS_OPTION iec.BYTE
	// LOG_CL is the log LOG_MSG writes to.
	LOG_CL = NewLOG_CONTROL()
)
