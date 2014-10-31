package main

import (
  "sort"
  "math"
)

// Consts

const (
  BEGIN           = -1
  END             = 1
  INVALID_OFFSET  = 0.0
)

// Types

type IntervalOffset struct {
  serverIndex int
  offset      float64
  offsetType  int
}

type IntervalOffsetArray []IntervalOffset

// Functions

func (arr IntervalOffsetArray) Len() int      { return len(arr) }
func (arr IntervalOffsetArray) Swap(i, j int) { arr[i], arr[j] = arr[j], arr[i] }

func (arr IntervalOffsetArray) Less(i, j int) bool {
  if arr[i].offset == arr[j].offset {
    if arr[i].offsetType == arr[j].offsetType {
      return i < j
    }
    return arr[i].offsetType < arr[j].offsetType
  }
  return arr[i].offset < arr[j].offset
}

func GetMaxIntersection(offsets IntervalOffsetArray) (maxIntersections int,
  maxIntersectionStart IntervalOffset,
  maxIntersectionEnd IntervalOffset) {

  var intersectionStart, intersectionEnd IntervalOffset
  var isMaxIntersection bool = false
  var numIntersections int = 0
  var maxIntersectionMidpoint float64 = 0

  sort.Sort(offsets)
  var globalMidpoint = (offsets[0].offset + offsets[len(offsets)-1].offset) / 2

  for _, offset := range offsets {
    if offset.offset == INVALID_OFFSET {
      continue
    }

    if offset.offsetType == BEGIN {
      intersectionStart = offset
    } else {
      intersectionEnd = offset
    }

    numIntersections -= offset.offsetType
    if numIntersections > maxIntersections {
      maxIntersections = numIntersections
      maxIntersectionStart = intersectionStart
      isMaxIntersection = true
    }

    if numIntersections == maxIntersections-1 && !isMaxIntersection && offset.offsetType == END {
      var midpoint = (intersectionStart.offset + offset.offset) / 2

      if math.Abs(midpoint-globalMidpoint) < math.Abs(midpoint-maxIntersectionMidpoint) {
        maxIntersectionStart = intersectionStart
        isMaxIntersection = true
      }
    }

    if offset.offsetType == END && isMaxIntersection && numIntersections == maxIntersections-1 {
      maxIntersectionEnd = intersectionEnd
      maxIntersectionMidpoint = (maxIntersectionStart.offset + maxIntersectionEnd.offset) / 2
      isMaxIntersection = false
    }
  }

  return maxIntersections, maxIntersectionStart, maxIntersectionEnd
}