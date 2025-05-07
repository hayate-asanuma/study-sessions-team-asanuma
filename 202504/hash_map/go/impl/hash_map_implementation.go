package impl

import "fmt"

// HashMapImplementation はHashMapの基本実装を提供する
type HashMapImplementation struct {
	buckets    [][]keyValue
	capacity   int
	size       int
	loadFactor float64
}

const defaultLoadFactor = 0.75

// NewHashMap は新しいHashMapを作成する
func NewHashMap(bucketSize int) *HashMapImplementation {
	if bucketSize <= 0 {
		bucketSize = 1
	}
	return &HashMapImplementation{
		buckets:    make([][]keyValue, bucketSize),
		capacity:   bucketSize,
		loadFactor: defaultLoadFactor,
	}
}

// hashKey はキーのハッシュ値を計算する
func (h *HashMapImplementation) hashKey(key string) int {
	hash := 0
	a := 256
	for i := 0; i < len(key); i++ {
		hash = (hash*a + int(key[i])) % h.capacity
	}
	if hash < 0 {
		hash = -hash
	}
	return hash
}

type keyValue struct {
	key   string
	value interface{}
}


// Put はキーと値のペアを格納する
func (h *HashMapImplementation) Put(key, value interface{}) {
	keyStr := fmt.Sprintf("%v", key)
	idx := h.hashKey(keyStr)
	bucket := h.buckets[idx]

	for i, kv := range bucket {
		if kv.key == keyStr {
			h.buckets[idx][i].value = value
			return
		}
	}
	h.buckets[idx] = append(bucket, keyValue{key: keyStr, value: value})
	h.size++

	// resize
	if float64(h.size)/float64(h.capacity) > h.loadFactor {
		h.resize()
	}
}

// Get はキーに対応する値を取得する
func (h *HashMapImplementation) Get(key interface{}) (interface{}, bool) {
	keyStr := fmt.Sprintf("%v", key)
	idx := h.hashKey(keyStr)
	for _, kv := range h.buckets[idx] {
		if kv.key == keyStr {
			return kv.value, true
		}
	}
	return nil, false
}

// Remove はキーに対応するエントリを削除する
func (h *HashMapImplementation) Remove(key interface{}) bool {
	keyStr := fmt.Sprintf("%v", key)
	idx := h.hashKey(keyStr)
	bucket := h.buckets[idx]
	for i, kv := range bucket {
		if kv.key == keyStr {
			h.buckets[idx] = append(bucket[:i], bucket[i+1:]...)
			h.size--
			return true
		}
	}
	return false
}

// resize はバケットサイズを拡張する
func (h *HashMapImplementation) resize() {
	newCapacity := h.capacity * 2 + 1
	newBuckets := make([][]keyValue, newCapacity)

	for _, bucket := range h.buckets {
		for _, kv := range bucket {
			hash := 0
			a := 256
			for i := 0; i < len(kv.key); i++ {
				hash = (hash*a + int(kv.key[i])) % newCapacity
			}
			if hash < 0 {
				hash = -hash
			}
			newBuckets[hash] = append(newBuckets[hash], kv)
		}
	}

	h.buckets = newBuckets
	h.capacity = newCapacity
}


// Size は現在の要素数を取得する
func (h *HashMapImplementation) Size() int {
	return h.size
}

// GetAllEntries は全てのエントリを取得する（テスト用）
func (h *HashMapImplementation) GetAllEntries() map[string]interface{} {
	result := make(map[string]interface{})
	for _, bucket := range h.buckets {
		for _, kv := range bucket {
			result[kv.key] = kv.value
		}
	}
	return result
}
