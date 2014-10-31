package main

import "fmt"

func getTestOffsets1() IntervalOffsetArray {
	return IntervalOffsetArray{
		IntervalOffset{serverIndex: 0, offset: 10, offsetType: BEGIN},
		IntervalOffset{serverIndex: 0, offset: 20, offsetType: END},
		IntervalOffset{serverIndex: 1, offset: 30, offsetType: BEGIN},
		IntervalOffset{serverIndex: 1, offset: 40, offsetType: END},
		IntervalOffset{serverIndex: 2, offset: 32, offsetType: BEGIN},
		IntervalOffset{serverIndex: 2, offset: 40, offsetType: END},
		IntervalOffset{serverIndex: 3, offset: 32, offsetType: BEGIN},
		IntervalOffset{serverIndex: 3, offset: 40, offsetType: END},
		IntervalOffset{serverIndex: 4, offset: 35, offsetType: BEGIN},
		IntervalOffset{serverIndex: 4, offset: 47, offsetType: END},
	}
}

func getTestOffsets2() IntervalOffsetArray {
	return IntervalOffsetArray{
		IntervalOffset{serverIndex: 0, offset: 10, offsetType: BEGIN},
		IntervalOffset{serverIndex: 0, offset: 20, offsetType: END},
		IntervalOffset{serverIndex: 1, offset: 30, offsetType: BEGIN},
		IntervalOffset{serverIndex: 1, offset: 40, offsetType: END},
		IntervalOffset{serverIndex: 2, offset: 50, offsetType: BEGIN},
		IntervalOffset{serverIndex: 2, offset: 60, offsetType: END},
		IntervalOffset{serverIndex: 3, offset: 70, offsetType: BEGIN},
		IntervalOffset{serverIndex: 3, offset: 80, offsetType: END},
		IntervalOffset{serverIndex: 4, offset: 90, offsetType: BEGIN},
		IntervalOffset{serverIndex: 4, offset: 100, offsetType: END},
	}
}

func getTestOffsets3() IntervalOffsetArray {
	return IntervalOffsetArray{
		IntervalOffset{serverIndex: 0, offset: 0, offsetType: BEGIN},
		IntervalOffset{serverIndex: 0, offset: 20, offsetType: END},
		IntervalOffset{serverIndex: 1, offset: 15, offsetType: BEGIN},
		IntervalOffset{serverIndex: 1, offset: 30, offsetType: END},
		IntervalOffset{serverIndex: 2, offset: 40, offsetType: BEGIN},
		IntervalOffset{serverIndex: 2, offset: 60, offsetType: END},
		IntervalOffset{serverIndex: 3, offset: 70, offsetType: BEGIN},
		IntervalOffset{serverIndex: 3, offset: 85, offsetType: END},
		IntervalOffset{serverIndex: 4, offset: 75, offsetType: BEGIN},
		IntervalOffset{serverIndex: 4, offset: 100, offsetType: END},
	}
}

func TestGetMaxIntersection() {
	testOffsets1 := getTestOffsets1()
  maxIntersections, maxIntersectionStart, maxIntersectionEnd := GetMaxIntersection(testOffsets1)
  if maxIntersections == 4 && maxIntersectionStart.offset == 35 && maxIntersectionEnd.offset == 40 {
	 fmt.Println("Passed")
  }

	testOffsets2 := getTestOffsets2()
	maxIntersections, maxIntersectionStart, maxIntersectionEnd = GetMaxIntersection(testOffsets2)
  if maxIntersections == 1 && maxIntersectionStart.offset == 50 && maxIntersectionEnd.offset == 60 {
   fmt.Println("Passed")
  }

	testOffsets3 := getTestOffsets3()
	maxIntersections, maxIntersectionStart, maxIntersectionEnd = GetMaxIntersection(testOffsets3)
  if maxIntersections == 1 && maxIntersectionStart.offset == 75 && maxIntersectionEnd.offset == 85 {
   fmt.Println("Passed")
  }
}

func main() {
  TestGetMaxIntersection()
}
