package impl

import (
	"fmt"
	"reflect"
)

// SortImplementation はソートアルゴリズムの基本実装を提供する
type SortImplementation struct {
    minRun int
}

// NewSortImplementation は新しいSortImplementationを作成する
func NewSortImplementation() *SortImplementation {
    return &SortImplementation{
        minRun: 32,
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

    // runのマージ
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

// findRuns は配列をrunに分割する
func (s *SortImplementation) findRuns(arr []interface{}, compare func(a, b interface{}) bool) [][]interface{} {
    runs := make([][]interface{}, 0)
    currentRun := make([]interface{}, 0)
    
    for i := 0; i < len(arr); i++ {
        currentRun = append(currentRun, arr[i])
        
        // runの終了条件をチェック
        if i == len(arr)-1 || !compare(arr[i], arr[i+1]) {
            if len(currentRun) < s.minRun {
                // 最小runサイズに満たない場合は挿入ソートで拡張
                s.insertionSort(currentRun, compare)
            }
            runs = append(runs, currentRun)
            currentRun = make([]interface{}, 0)
        }
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
