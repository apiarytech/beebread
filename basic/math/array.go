/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - GPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

package math

import (
	"math"
	"math/rand"
	"sort"
)

// ARRAY_ABS calculates the absolute value of every element in an array of float64.
// The operation is done in-place.
func ARRAY_ABS(arr []float64) {
	for i := range arr {
		arr[i] = math.Abs(arr[i])
	}
}

// ARRAY_ADD adds a value to every element in an array of float64.
// The operation is done in-place.
func ARRAY_ADD(arr []float64, val float64) {
	for i := range arr {
		arr[i] += val
	}
}

// ARRAY_INIT initializes every element in an array of float64 with a given value.
// The operation is done in-place.
func ARRAY_INIT(arr []float64, val float64) {
	for i := range arr {
		arr[i] = val
	}
}

// ARRAY_MEDIAN calculates the median of an array of float64.
// Note: This function sorts the array in-place before calculating the median.
func ARRAY_MEDIAN(arr []float64) float64 {
	size := len(arr)
	if size == 0 {
		return 0.0
	}

	// Sort the array to find the median
	sort.Float64s(arr)

	mid := size / 2
	if size%2 == 1 {
		// Odd number of elements
		return arr[mid]
	}

	// Even number of elements
	return (arr[mid-1] + arr[mid]) / 2.0
}

// ARRAY_MUL multiplies every element in an array of float64 with a value.
// The operation is done in-place.
func ARRAY_MUL(arr []float64, val float64) {
	for i := range arr {
		arr[i] *= val
	}
}

// ARRAY_SHUFFLE shuffles every element in an array of float64 randomly.
// The operation is done in-place.
func ARRAY_SHUFFLE(arr []float64) {
	// Using Go's standard library rand.Shuffle is the idiomatic way.
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}

// ARRAY_SORT sorts an array of float64 either ascending or descending.
// The operation is done in-place.
func ARRAY_SORT(arr []float64, desc bool) {
	if desc {
		sort.Sort(sort.Reverse(sort.Float64Slice(arr)))
	} else {
		sort.Float64s(arr)
	}
}

// ARRAY_AVG calculates the average of a given array of float64.
func ARRAY_AVG(arr []float64) float64 {
	if len(arr) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, v := range arr {
		sum += v
	}
	return sum / float64(len(arr))
}

// ARRAY_MAX finds the maximum value in an array of float64.
func ARRAY_MAX(arr []float64) float64 {
	if len(arr) == 0 {
		return 0.0
	}
	max := arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] > max {
			max = arr[i]
		}
	}
	return max
}

// ARRAY_MIN finds the minimum value in an array of float64.
func ARRAY_MIN(arr []float64) float64 {
	if len(arr) == 0 {
		return 0.0
	}
	min := arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] < min {
			min = arr[i]
		}
	}
	return min
}

// ARRAY_SUM calculates the sum of all elements in an array of float64.
func ARRAY_SUM(arr []float64) float64 {
	sum := 0.0
	for _, v := range arr {
		sum += v
	}
	return sum
}

// ARRAY_VAR calculates the sample variance of a given array of float64.
func ARRAY_VAR(arr []float64) float64 {
	n := len(arr)
	if n < 2 { // Variance requires at least 2 data points.
		return 0.0
	}

	avg := ARRAY_AVG(arr)
	variance := 0.0
	for _, v := range arr {
		variance += (v - avg) * (v - avg)
	}
	return variance / float64(n-1)
}

// ARRAY_SDV calculates the standard deviation of a given array of float64.
func ARRAY_SDV(arr []float64) float64 {
	return math.Sqrt(ARRAY_VAR(arr))
}

// ARRAY_SPR calculates the spread (range) of a given array of float64.
func ARRAY_SPR(arr []float64) float64 {
	if len(arr) == 0 {
		return 0.0
	}
	return ARRAY_MAX(arr) - ARRAY_MIN(arr)
}

// ARRAY_TREND calculates the trend of a given array.
// It's the average of the second half minus the average of the first half.
func ARRAY_TREND(arr []float64) float64 {
	n := len(arr)
	if n < 2 {
		return 0.0
	}

	mid := n / 2
	firstHalf := arr[:mid]
	secondHalf := arr[mid:]

	// The original ST code has a slight difference for even/odd lengths,
	// but this approach is cleaner and captures the intent.
	return ARRAY_AVG(secondHalf) - ARRAY_AVG(firstHalf)
}

// ARRAY_GAV calculates the geometric average of a given array of float64.
func ARRAY_GAV(arr []float64) float64 {
	n := len(arr)
	if n == 0 {
		return 0.0
	}

	sumOfLogs := 0.0
	for _, v := range arr {
		if v <= 0 {
			return 0.0 // Geometric mean is undefined for non-positive numbers.
		}
		sumOfLogs += math.Log(v)
	}
	return math.Exp(sumOfLogs / float64(n))
}

// ARRAY_HAV calculates the harmonic average of a given array of float64.
func ARRAY_HAV(arr []float64) float64 {
	n := len(arr)
	if n == 0 {
		return 0.0
	}

	sumOfInverses := 0.0
	for _, v := range arr {
		if v == 0 {
			return 0.0 // Harmonic mean is undefined if any element is zero.
		}
		sumOfInverses += 1.0 / v
	}
	return float64(n) / sumOfInverses
}

// IS_SORTED checks if an array of float64 is sorted in ascending order.
func IS_SORTED(arr []float64) bool {
	// Go's standard library provides a convenient function for this.
	return sort.Float64sAreSorted(arr)
}
