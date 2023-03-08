package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	// 读取语料库文件 https://norvig.com/big.txt
	filename := "big.txt"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	wordRegex := regexp.MustCompile("^\\w+$")
	// 读取文件中的所有单词并计算所有 3-gram 的频率
	scanner := bufio.NewScanner(file)
	ngrams := map[string]int{}
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, " ")
		for _, word := range words {
			word := strings.ToLower(word)
			if !wordRegex.Match([]byte(word)) {
				continue
			}
			for _, ngram := range generateNGrams(word, 3) {
				ngrams[ngram]++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	output, err := os.OpenFile("output.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer output.Close()
	writer := bufio.NewWriter(output)
	defer writer.Flush()
	// 生成指定数量和长度的单词
	numWords := 100000
	wordLength := 5
	suffix := ""
	for i := 0; i < numWords; i++ {
		word := generateWord(ngrams, wordLength)
		//fmt.Println(word)
		_, err := writer.WriteString(fmt.Sprintf("%s%s\n", word, suffix))
		if err != nil {
			panic(err)
		}
	}

}

// 生成指定长度的单词
func generateWord(ngrams map[string]int, length int) string {
	// 从概率模型中选择一个 3-gram 作为单词的前缀
	prefix := ""
	for !isValidPrefix(ngrams, prefix) {
		keys := make([]string, 0, len(ngrams))
		for k := range ngrams {
			if len(k) == 3 {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)

		var sum int
		for _, k := range keys {
			sum += ngrams[k]
		}

		pos := rand.Intn(sum)
		for _, k := range keys {
			pos -= ngrams[k]
			if pos < 0 {
				prefix = k
				break
			}
		}
	}

	word := prefix

	// 继续选择概率模型中符合条件的3-gram
	for len(word) < length {
		keys := make([]string, 0, len(ngrams))
		for k := range ngrams {
			if strings.HasPrefix(k, word[len(word)-2:]) && len(k) == 3 {
				keys = append(keys, k)
			}
		}

		if len(keys) == 0 {
			break
		}

		sort.Strings(keys)

		var sum int
		for _, k := range keys {
			sum += ngrams[k]
		}

		pos := rand.Intn(sum)
		for _, k := range keys {
			pos -= ngrams[k]
			if pos < 0 {
				word += k[2:]
				break
			}
		}
	}

	return word
}

// 判断指定的 3-gram 是否是合法的前缀
func isValidPrefix(ngrams map[string]int, prefix string) bool {
	if len(prefix) < 3 {
		return false
	}

	for k := range ngrams {
		if strings.HasPrefix(k, prefix[len(prefix)-2:]) {
			return true
		}
	}

	return false
}

// 生成指定长度的 n-gram 序列
func generateNGrams(s string, n int) []string {
	runes := []rune(s)

	if len(runes) < n {
		return nil
	}

	ngrams := make([]string, len(runes)-n+1)
	for i := 0; i < len(runes)-n+1; i++ {
		ngrams[i] = string(runes[i : i+n])
	}

	return ngrams
}
