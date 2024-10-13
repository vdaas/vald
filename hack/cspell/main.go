// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"bufio"
	"cmp"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"

	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
)

type CSpellConfig struct {
	Version          string              `json:"version"`
	Language         string              `json:"language"`
	Import           []string            `json:"import"`
	IgnorePaths      []string            `json:"ignorePaths"`
	Patterns         []PatternEntry      `json:"patterns,omitempty"`
	IgnoreRegExpList []string            `json:"ignoreRegExpList,omitempty"`
	IgnoreWords      []string            `json:"ignoreWords,omitempty"`
	IgnoreWordsMap   map[string][]string `json:"ignoreWordsMap,omitempty"`
}

type PatternEntry struct {
	Name    string `json:"name"`
	Pattern string `json:"pattern"`
}

type Override struct {
	Files []string `json:"files"`
	Words []string `json:"words"`
}

var (
	mandatoryImports = []string{
		"@cspell/dict-cpp/cspell-ext.json",
		"@cspell/dict-docker/cspell-ext.json",
		"@cspell/dict-en_us/cspell-ext.json",
		"@cspell/dict-fullstack/cspell-ext.json",
		"@cspell/dict-git/cspell-ext.json",
		"@cspell/dict-golang/cspell-ext.json",
		"@cspell/dict-k8s/cspell-ext.json",
		"@cspell/dict-makefile/cspell-ext.json",
		"@cspell/dict-markdown/cspell-ext.json",
		"@cspell/dict-npm/cspell-ext.json",
		"@cspell/dict-public-licenses/cspell-ext.json",
		"@cspell/dict-rust/cspell-ext.json",
		"@cspell/dict-shell/cspell-ext.json",
	}

	mandatoryIgnorePaths = []string{
		"**/*.ai",
		"**/*.drawio",
		"**/*.hdf5",
		"**/*.key",
		"**/*.lock",
		"**/*.log",
		"**/*.md5",
		"**/*.pack",
		"**/*.pdf",
		"**/*.pem",
		"**/*.png",
		"**/*.sum",
		"**/*.svg",
		"**/.cspell.json",
		"**/.git/objects/**",
		"**/cmd/agent/core/faiss/faiss",
		"**/cmd/agent/core/ngt/ngt",
		"**/cmd/agent/sidecar/sidecar",
		"**/cmd/discoverer/k8s/discoverer",
		"**/cmd/gateway/filter/filter",
		"**/cmd/gateway/lb/lb",
		"**/cmd/gateway/mirror/mirror",
		"**/cmd/index/job/correction/index-correction",
		"**/cmd/index/job/creation/index-creation",
		"**/cmd/index/job/creation/index-deletion",
		"**/cmd/index/job/readreplica/rotate/readreplica-rotate",
		"**/cmd/index/job/save/index-save",
		"**/cmd/index/operator/index-operator",
		"**/cmd/manager/index/index",
		"**/cmd/tools/benchmark/job/job",
		"**/cmd/tools/benchmark/operator/operator",
		"**/cmd/tools/cli/loadtest/loadtest",
		"**/hack/cspell/**",
		"**/internal/core/algorithm/ngt/assets/index",
		"**/internal/test/data/agent/ngt/validIndex",
	}
	suffixes = []string{
		"addr",
		"addrs",
		"buf",
		"cancel",
		"cfg",
		"ch",
		"cnt",
		"conf",
		"conn",
		"ctx",
		"dim",
		"dur",
		"env",
		"err",
		"error",
		"errors",
		"errs",
		"idx",
		"len",
		"mu",
		"opt",
		"opts",
		"pool",
		"req",
		"res",
		"size",
		"vec",
	}

	sufReg = regexp.MustCompile(fmt.Sprintf("(%s)$", strings.Join(suffixes, "|")))

	prexp = regexp.MustCompile(`Unknown word \((.*?)\) Suggestions`)
)

func extractLine(line string) (filePath, word string, ok bool) {
	filePath, line, ok = strings.Cut(line, ":")
	if !ok || len(filePath) == 0 {
		return "", "", false
	}
	_, s, ok := strings.Cut(line, " - ")
	if ok {
		line = s
	}
	matches := prexp.FindStringSubmatch(line)
	if len(matches) > 1 {
		return filePath, matches[1], true
	}
	return "", "", false
}

func parseCspellResult(filePath string, th int) (map[string][]string, map[string]bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	wordsByFile := make(map[string][]string)
	filesByWord := make(map[string][]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Extract the unknown word
			if path, word, ok := extractLine(line); ok {
				if sufReg.MatchString(word) {
					return
				}
				lword := strings.ToLower(word)
				mu.Lock()
				w, ok := wordsByFile[path]
				if !ok || w == nil {
					w = make([]string, 0, 2)
				}
				wordsByFile[path] = append(w, word)

				f, ok := filesByWord[word]
				if !ok || f == nil {
					f = make([]string, 0, 2)
				}
				filesByWord[word] = append(f, path)

				if word != lword {
					f, ok = filesByWord[lword]
					if !ok || f == nil {
						f = make([]string, 0, 2)
					}
					filesByWord[lword] = append(f, path)
				}
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	globalWords := make(map[string]bool)
	for word, files := range filesByWord {
		lword := strings.ToLower(word)
		if word != lword {
			lfiles, ok := filesByWord[lword]
			if ok {
				files = append(files, lfiles...)
				slices.Sort(files)
				files = slices.Compact(files)
				if len(files) >= th {
					globalWords[lword] = true
					globalWords[word] = true
				}
			}
		} else if len(files) >= th {
			globalWords[word] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading file: %w", err)
	}

	return wordsByFile, globalWords, nil
}

func loadConfig(path string) (config *CSpellConfig, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config = new(CSpellConfig)
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func saveConfig(path string, config *CSpellConfig) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

func addPatterns(config *CSpellConfig) {
	if config.Patterns == nil {
		config.Patterns = make([]PatternEntry, 0, len(suffixes))
	}
	if config.IgnoreRegExpList == nil {
		config.IgnoreRegExpList = make([]string, 0, len(suffixes))
	}
	for _, suffix := range suffixes {
		pat := fmt.Sprintf("/\\b\\w*%s\\b/", suffix)
		name := fmt.Sprintf("Ignore_%s_suffix", suffix)
		config.Patterns = append(config.Patterns, PatternEntry{
			Name:    name,
			Pattern: pat,
		})
		config.IgnoreRegExpList = append(config.IgnoreRegExpList, name)
	}
	slices.SortFunc(config.Patterns, func(left, right PatternEntry) int {
		return cmp.Compare(left.Name, right.Name)
	})
	config.Patterns = slices.CompactFunc(config.Patterns, func(left, right PatternEntry) bool {
		return left.Name == right.Name
	})
	slices.Sort(config.IgnoreRegExpList)
	config.IgnoreRegExpList = slices.Compact(config.IgnoreRegExpList)
}

func main() {
	configPath := flag.String("config", ".cspell.json", "Path to the existing .cspell.json file")
	outputPath := flag.String("output", "", "Path to the cspell output log")
	threshold := flag.Int("threshold", 5, "Threshold for declaring global words")

	flag.Parse()

	if *outputPath == "" {
		fmt.Println("Error: output path is required")
		os.Exit(1)
	}

	config, err := loadConfig(*configPath)
	if err != nil || config == nil {
		config = new(CSpellConfig)
	}
	config.Import = mandatoryImports
	config.IgnorePaths = mandatoryIgnorePaths
	config.Version = "0.2"
	config.Language = "en"

	addPatterns(config)

	wordsByFile, globalWords, err := parseCspellResult(*outputPath, *threshold)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if config.IgnoreWords == nil {
		config.IgnoreWords = make([]string, 0, len(globalWords))
	}
	for word := range globalWords {
		config.IgnoreWords = append(config.IgnoreWords, word)
		fmt.Println(config.IgnoreWords)
	}
	slices.Sort(config.IgnoreWords)
	for _, word := range config.IgnoreWords {
		globalWords[word] = true
	}
	if config.IgnoreWordsMap == nil {
		config.IgnoreWordsMap = make(map[string][]string, len(wordsByFile))
	}
	for filePath, words := range wordsByFile {
		words = slices.DeleteFunc(words, func(word string) bool {
			return globalWords[strings.ToLower(word)]
		})
		if len(words) > 0 {
			im, ok := config.IgnoreWordsMap[filePath]
			if !ok || im == nil {
				slices.Sort(words)
				config.IgnoreWordsMap[filePath] = words
			} else {
				words = append(im, words...)
				slices.Sort(words)
				config.IgnoreWordsMap[filePath] = slices.Compact(words)
			}
		}
	}

	if err := saveConfig(*configPath, config); err != nil {
		fmt.Println("Error: output path is required")
		os.Exit(1)
	}
}
