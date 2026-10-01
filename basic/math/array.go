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

package math

import (
	"slices"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

// The array functions take an array of REAL and its size in bytes, as OSCAT
// calls them with ADR(array) and SIZEOF(array). A REAL is 4 bytes.

// elements returns the elements of pt that a size in bytes covers.
func elements(pt []iec.REAL, size iec.UINT) []iec.REAL {
	n := int(size >> 2)
	if n > len(pt) {
		n = len(pt)
	}
	return pt[:n]
}

// ARRAY_ABS_ replaces each element of the array with its absolute value.
func ARRAY_ABS_(pt []iec.REAL, size iec.UINT) iec.BOOL {
	a := elements(pt, size)
	for i := range a {
		a[i] = ABS(a[i])
	}
	return true
}

// ARRAY_ADD_ adds x to each element of the array.
func ARRAY_ADD_(pt []iec.REAL, size iec.UINT, x iec.REAL) iec.BOOL {
	a := elements(pt, size)
	for i := range a {
		a[i] = a[i] + x
	}
	return true
}

// ARRAY_INIT_ sets each element of the array to init.
func ARRAY_INIT_(pt []iec.REAL, size iec.UINT, init iec.REAL) iec.BOOL {
	a := elements(pt, size)
	for i := range a {
		a[i] = init
	}
	return true
}

// ARRAY_MEDIAN_ returns the median of the array. The array is sorted and
// stays sorted: for [12,0,4,7,1] the median is 4 and the array afterwards is
// [0,1,4,7,12].
func ARRAY_MEDIAN_(pt []iec.REAL, size iec.UINT) iec.REAL {
	ARRAY_SORT_(pt, size)
	a := elements(pt, size)
	if len(a) == 0 {
		return 0.0
	}
	stop := len(a) - 1
	if stop%2 == 0 {
		return a[stop>>1]
	}
	i := stop >> 1
	return (a[i] + a[i+1]) * 0.5
}

// ARRAY_MUL_ multiplies each element of the array with x.
func ARRAY_MUL_(pt []iec.REAL, size iec.UINT, x iec.REAL) iec.BOOL {
	a := elements(pt, size)
	for i := range a {
		a[i] = a[i] * x
	}
	return true
}

// ARRAY_SHUFFLE_ randomly shuffles the elements of the array.
func ARRAY_SHUFFLE_(pt []iec.REAL, size iec.UINT) iec.BOOL {
	a := elements(pt, size)
	stop := iec.INT(len(a) - 1)
	var pos iec.INT
	for i := iec.INT(0); i <= stop; i++ {
		pos = RDM2(i+pos, 0, stop)
		a[i], a[pos] = a[pos], a[i]
	}
	return true
}

// ARRAY_SORT_ sorts the array in ascending order.
func ARRAY_SORT_(pt []iec.REAL, size iec.UINT) iec.BOOL {
	slices.Sort(elements(pt, size))
	return true
}

// ARRAY_AVG returns the average of the array.
func ARRAY_AVG(pt []iec.REAL, size iec.UINT) iec.REAL {
	a := elements(pt, size)
	if len(a) == 0 {
		return 0.0
	}
	return ARRAY_SUM(pt, size) / iec.REAL(len(a))
}

// ARRAY_GAV returns the geometric average of the array. It is 0 if an
// element is not above 0.
func ARRAY_GAV(pt []iec.REAL, size iec.UINT) iec.REAL {
	a := elements(pt, size)
	var gav iec.REAL = 1.0
	for _, v := range a {
		if v > 0.0 {
			gav = gav * v
		} else {
			return 0.0
		}
	}
	return SQRTN(gav, iec.INT(len(a)))
}

// ARRAY_HAV returns the harmonic average of the array. It is 0 if an
// element is 0.
func ARRAY_HAV(pt []iec.REAL, size iec.UINT) iec.REAL {
	a := elements(pt, size)
	if len(a) == 0 {
		return 0.0
	}
	var hav iec.REAL
	for _, v := range a {
		if v != 0.0 {
			hav = hav + 1.0/v
		} else {
			return 0.0
		}
	}
	return iec.REAL(len(a)) / hav
}

// ARRAY_MAX returns the maximum of the array.
func ARRAY_MAX(pt []iec.REAL, size iec.UINT) iec.REAL {
	a := elements(pt, size)
	if len(a) == 0 {
		return 0.0
	}
	mx := a[0]
	for _, v := range a[1:] {
		if v > mx {
			mx = v
		}
	}
	return mx
}

// ARRAY_MIN returns the minimum of the array.
func ARRAY_MIN(pt []iec.REAL, size iec.UINT) iec.REAL {
	a := elements(pt, size)
	if len(a) == 0 {
		return 0.0
	}
	mn := a[0]
	for _, v := range a[1:] {
		if v < mn {
			mn = v
		}
	}
	return mn
}

// ARRAY_SDV returns the standard deviation of the array.
func ARRAY_SDV(pt []iec.REAL, size iec.UINT) iec.REAL {
	return SQRT(ARRAY_VAR(pt, size))
}

// ARRAY_SPR returns the spread of the array, its maximum less its minimum.
func ARRAY_SPR(pt []iec.REAL, size iec.UINT) iec.REAL {
	return ARRAY_MAX(pt, size) - ARRAY_MIN(pt, size)
}

// ARRAY_SUM returns the sum of the array.
func ARRAY_SUM(pt []iec.REAL, size iec.UINT) iec.REAL {
	var sum iec.REAL
	for _, v := range elements(pt, size) {
		sum = sum + v
	}
	return sum
}

// ARRAY_TREND returns the trend of the array, the average of its second half
// less the average of its first half: for [0,1,4,5,3,4,6,3] it is
// 4 - 2.5 = 1.5. The middle element of an array with an odd number of
// elements counts in both halves.
func ARRAY_TREND(pt []iec.REAL, size iec.UINT) iec.REAL {
	a := elements(pt, size)
	if len(a) == 0 {
		return 0.0
	}
	stop := len(a) - 1
	stop2 := stop >> 1
	var x iec.REAL
	for i := 0; i <= stop2; i++ {
		x = x - a[i]
	}
	start := stop2 + 1
	if stop%2 == 0 {
		start = stop2
	}
	for i := start; i <= stop; i++ {
		x = x + a[i]
	}
	return x / iec.REAL(stop2+1)
}

// ARRAY_VAR returns the variance of the array.
func ARRAY_VAR(pt []iec.REAL, size iec.UINT) iec.REAL {
	a := elements(pt, size)
	if len(a) == 0 {
		return 0.0
	}
	avg := ARRAY_AVG(pt, size)
	var sum iec.REAL
	for _, v := range a {
		sum = sum + (v-avg)*(v-avg)
	}
	return sum / iec.REAL(len(a)-1)
}

// IS_SORTED reports whether the array is sorted in ascending order.
func IS_SORTED(pt []iec.REAL, size iec.UINT) iec.BOOL {
	a := elements(pt, size)
	for i := 0; i+1 < len(a); i++ {
		if a[i] > a[i+1] {
			return false
		}
	}
	return true
}
