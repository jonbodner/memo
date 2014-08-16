package memo

import (
	_ "fmt"
	"testing"
)

var sumCount int

//use side-effects to make sure that it's only called once per identical params
func sum(a, b int) int {
	sumCount++
	return a + b
}

var multCount int

func mult(a, b int) int {
	multCount++
	return a * b
}

func TestLongMemo(t *testing.T) {
	localSumCount := 0

	localSum := func(a, b, c, d, e, f, g, h int) int {
		localSumCount++
		return a + b + c + d + e + f + g + h
	}
	Memoize(&localSum)
	result := localSum(2, 1, 1, 1, 1, 1, 1, 1)
	if result != 9 {
		t.Fatalf("expected sum of 9, got %v\n", result)
	}
	result2 := localSum(2, 1, 1, 1, 1, 1, 1, 1)
	if result2 != 9 {
		t.Fatalf("expected sum of 9, got %v\n", result2)
	}
	if localSumCount != 1 {
		t.Fatalf("expected sum to only be called once, was %v\n", sumCount)
	}
	result3 := localSum(2, 3, 1, 1, 1, 1, 1, 1)
	if result3 != 11 {
		t.Fatalf("expected sum of 11, got %v\n", result3)
	}
	result4 := localSum(2, 3, 1, 1, 1, 1, 1, 1)
	if result4 != 11 {
		t.Fatalf("expected sum of 11, got %v\n", result4)
	}
	if localSumCount != 2 {
		t.Fatalf("expected sum to be called twice, was %v\n", sumCount)
	}
}

func TestMemo(t *testing.T) {
	//reset before running tests
	sumCount = 0
	multCount = 0

	localSum := func(a, b int) int {
		return sum(a, b)
	}
	Memoize(&localSum)
	localMult := func(a, b int) int {
		return mult(a, b)
	}
	Memoize(&localMult)
	result := localSum(2, 1)
	if result != 3 {
		t.Fatalf("expected sum of 3, got %v\n", result)
	}
	result2 := localSum(2, 1)
	if result2 != 3 {
		t.Fatalf("expected sum of 3, got %v\n", result2)
	}
	if sumCount != 1 {
		t.Fatalf("expected sum to only be called once, was %v\n", sumCount)
	}
	result3 := localSum(2, 3)
	if result3 != 5 {
		t.Fatalf("expected sum of 5, got %v\n", result3)
	}
	result4 := localSum(2, 3)
	if result4 != 5 {
		t.Fatalf("expected sum of 5, got %v\n", result4)
	}
	if sumCount != 2 {
		t.Fatalf("expected sum to be called twice, was %v\n", sumCount)
	}
	result5 := localMult(2, 1)
	if result5 != 2 {
		t.FailNow()
	}
	result6 := localMult(2, 1)
	if result6 != 2 {
		t.FailNow()
	}
	if multCount != 1 {
		t.FailNow()
	}
	result7 := localMult(2, 3)
	if result7 != 6 {
		t.FailNow()
	}
	result8 := localMult(2, 3)
	if result8 != 6 {
		t.FailNow()
	}
	if multCount != 2 {
		t.FailNow()
	}
}

func TestRecursive(t *testing.T) {
	fibCount := 0
	var localFib func(int) int
	localFib = func(n int) int {
		fibCount++
		if n == 0 {
			return 0
		}
		if n == 1 {
			return 1
		}
		return localFib(n-2) + localFib(n-1)
	}

	Memoize(&localFib)
	r := localFib(1)
	if r != 1 {
		t.FailNow()
	}
	if fibCount != 1 {
		t.FailNow()
	}
	r2 := localFib(1)
	if r2 != 1 {
		t.FailNow()
	}
	if fibCount != 1 {
		t.FailNow()
	}
	//should result in 10 more calls
	r3 := localFib(10)
	// 10th fib == 0 1 1 2 3 5 8 13 21 34 55
	if r3 != 55 {
		t.FailNow()
	}
	if fibCount != 11 {
		t.FailNow()
	}
}

var squareCount int = 0

var square func(int) int = func(x int) int {
	squareCount++
	return x * x
}

func TestPackageVar(t *testing.T) {
	Memoize(&square)
	s := square(2)
	if s != 4 {
		t.FailNow()
	}
	s2 := square(2)
	if s2 != 4 {
		t.FailNow()
	}
	if squareCount != 1 {
		t.FailNow()
	}
}

func BenchmarkMemo(b *testing.B) {
	localSum := func(a, b int) int {
		return sum(a, b)
	}
	Memoize(&localSum)
	for i := 0; i < b.N; i++ {
		localSum(2, 1)
	}
}

func BenchmarkSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sum(2, 1)
	}
}

func fib(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return fib(n-2) + fib(n-1)
}

func BenchmarkMemoFib(b *testing.B) {
	var localFib func(int) int
	localFib = func(n int) int {
		if n == 0 {
			return 0
		}
		if n == 1 {
			return 1
		}
		return localFib(n-2) + localFib(n-1)
	}
	Memoize(&localFib)
	for i := 0; i < b.N; i++ {
		localFib(20)
	}
}

func BenchmarkFib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fib(20)
	}
}

func BenchmarkBig(b *testing.B) {
	localSum := func(a, b, c, d, e, f, g, h int) int {
		return a + b + c + d + e + f + g + h
	}
	Memoize(&localSum)
	for i := 0; i < b.N; i++ {
		localSum(1, 2, 3, 4, 5, 6, 7, 8)
	}
}
