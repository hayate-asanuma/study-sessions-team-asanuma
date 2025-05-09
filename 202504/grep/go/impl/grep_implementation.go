package impl

import (
	"bufio"
	"log"
	"os"
)

// GrepImplementation はGrepの基本実装を提供する
type GrepImplementation struct{}

func (g *GrepImplementation) Search(filePath, pattern string) []string {
	var result []string
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("failed to open file %s: %v", filePath, err)
		return result
	}
	defer f.Close()

	lps := buildLPS(pattern) 
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if kmpMatch(line, pattern, lps) { 
			result = append(result, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("failed to read file %s: %v", filePath, err)
	}
	return result
}

func buildLPS(pat string) []int {
	lps := make([]int, len(pat))
	for i, lenPrev := 1, 0; i < len(pat); {
		if pat[i] == pat[lenPrev] {
			lenPrev++
			lps[i] = lenPrev
			i++
		} else if lenPrev != 0 {
			lenPrev = lps[lenPrev-1]
		} else {
			lps[i] = 0
			i++
		}
	}
	return lps
}

func kmpMatch(txt, pat string, lps []int) bool {
	i, j := 0, 0
	for i < len(txt) {
		if txt[i] == pat[j] {
			i++
			j++
			if j == len(pat) {
				return true
			}
		} else if j != 0 {
			j = lps[j-1]
		} else {
			i++
		}
	}
	return false
}