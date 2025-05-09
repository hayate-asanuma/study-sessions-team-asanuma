package impl

import (
	"fmt"
	"reflect"
	"sync"
)

// SortImplementation はソートアルゴリズムの基本実装を提供する
type SortImplementation struct {
	minRun            int
	parallelThreshold int
}

// NewSortImplementation は新しいSortImplementationを作成する
func NewSortImplementation() *SortImplementation {
	return &SortImplementation{
		minRun:            32,
		parallelThreshold: 10000, // 並列化するかの閾値
	}
}

// Sort はTimSortアルゴリズムを使用してソートを行う
func (s *SortImplementation) Sort(data []interface{}) []interface{} {
	if len(data) == 0 {
		return data
	}

	// 元スライスをコピー
	newArr := make([]interface{}, len(data))
	copy(newArr, data)

	// 型チェックと比較関数の取得
	compare, err := s.getCompareFunction(newArr)
	if err != nil {
		return newArr
	}

	// runの検出と挿入ソート
	runs := s.findRuns(newArr, compare)

	if len(runs) == 0 {
		return newArr
	}
	if len(runs) == 1 {
		return runs[0]
	}

	if len(data) >= s.parallelThreshold {
		return s.parallelMergeRuns(runs, compare)
	}
	return s.mergeRuns(runs, compare)
}

// getCompareFunction は配列の型に応じた比較関数を返す
func (s *SortImplementation) getCompareFunction(arr []interface{}) (func(a, b interface{}) bool, error) {
	if len(arr) == 0 {
		return nil, nil
	}

	// 最初の要素の型を取得
	firstType := reflect.TypeOf(arr[0])

	// 配列内のすべての要素が同じ型かチェック
	for _, v := range arr {
		if reflect.TypeOf(v) != firstType {
			return nil, fmt.Errorf("mixed types in array")
		}
	}

	// 型に応じた比較関数を返す
	switch firstType {
	case reflect.TypeOf(""):
		return func(a, b interface{}) bool {
			return a.(string) < b.(string)
		}, nil
	case reflect.TypeOf(0):
		return func(a, b interface{}) bool {
			return a.(int) < b.(int)
		}, nil
	case reflect.TypeOf(0.0):
		return func(a, b interface{}) bool {
			return a.(float64) < b.(float64)
		}, nil
	default:
		return nil, fmt.Errorf("unsupported type: %v", firstType)
	}
}

// スライス arr の指定範囲の要素を反転
func (s *SortImplementation) reverseRange(arr []interface{}, start, end int) {
	for start < end {
		arr[start], arr[end] = arr[end], arr[start]
		start++
		end--
	}
}

// findRuns は配列 'arr' を TimSort の run に分割
// 処理済 run のコピーのスライスを返します。
func (s *SortImplementation) findRuns(arr []interface{}, compare func(a, b interface{}) bool) [][]interface{} {
	n := len(arr)
	runs := make([][]interface{}, 0)
	idx := 0

	for idx < n {
		runLo := idx
		runHi := idx

		if idx == n-1 {
			runHi = idx
		} else {
			descending := compare(arr[idx+1], arr[idx])

			runHi = idx + 1
			for runHi < n-1 {
				if descending {
					if !compare(arr[runHi+1], arr[runHi]) {
						break
					}
				} else {
					if compare(arr[runHi+1], arr[runHi]) {
						break
					}
				}
				runHi++
			}

			if descending {
				s.reverseRange(arr, runLo, runHi)
			}
		}

		currentNaturalRunLen := runHi - runLo + 1
		if currentNaturalRunLen < s.minRun {
			forcedHi := runLo + s.minRun - 1
			if forcedHi >= n {
				forcedHi = n - 1
			}

			if forcedHi > runLo {
				s.insertionSort(arr[runLo:forcedHi+1], compare)
			}
			runHi = forcedHi
		}

		runLen := runHi - runLo + 1
		runCopy := make([]interface{}, runLen)
		copy(runCopy, arr[runLo:runHi+1])
		runs = append(runs, runCopy)

		idx = runHi + 1
	}
	return runs
}

// insertionSort は挿入ソートを実行する
func (s *SortImplementation) insertionSort(run []interface{}, compare func(a, b interface{}) bool) {
	for i := 1; i < len(run); i++ {
		key := run[i]
		j := i - 1

		for j >= 0 && compare(key, run[j]) {
			run[j+1] = run[j]
			j--
		}
		run[j+1] = key
	}
}

// mergeRuns はrunをマージする
func (s *SortImplementation) mergeRuns(runs [][]interface{}, compare func(a, b interface{}) bool) []interface{} {
	if len(runs) == 0 {
		return nil
	}
	if len(runs) == 1 {
		return runs[0]
	}

	// 2つのrunをマージ
	merged := s.mergeTwoRuns(runs[0], runs[1], compare)

	// 残りのrunをマージ
	for i := 2; i < len(runs); i++ {
		merged = s.mergeTwoRuns(merged, runs[i], compare)
	}

	return merged
}

// mergeTwoRuns は2つのrunをマージする
func (s *SortImplementation) mergeTwoRuns(run1, run2 []interface{}, compare func(a, b interface{}) bool) []interface{} {
	result := make([]interface{}, 0, len(run1)+len(run2))
	i, j := 0, 0

	for i < len(run1) && j < len(run2) {
		if compare(run1[i], run2[j]) {
			result = append(result, run1[i])
			i++
		} else {
			result = append(result, run2[j])
			j++
		}
	}

	// 残りの要素を追加
	result = append(result, run1[i:]...)
	result = append(result, run2[j:]...)

	return result
}

// parallelMergeRuns は run を並列にマージする
func (s *SortImplementation) parallelMergeRuns(runs [][]interface{}, compare func(a, b interface{}) bool) []interface{} {
	for len(runs) > 1 {
		var wg sync.WaitGroup
		mergedCount := (len(runs) + 1) / 2
		merged := make([][]interface{}, mergedCount)

		for i := 0; i < len(runs); i += 2 {
			i := i // goroutine 内で使うためコピー
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				if index+1 < len(runs) {
					merged[index/2] = s.mergeTwoRuns(runs[index], runs[index+1], compare)
				} else {
					merged[index/2] = runs[index]
				}
			}(i)
		}

		wg.Wait()
		runs = merged
	}
	return runs[0]
}
