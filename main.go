package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// 用「大小 + sha256」两阶段找重复，先按大小分组省掉大量哈希计算。
func findDuplicates(root string) ([][]string, error) {
	sizeGroups := map[int64][]string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		sizeGroups[info.Size()] = append(sizeGroups[info.Size()], path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	hashGroups := map[string][]string{}
	for _, paths := range sizeGroups {
		// 同大小的才有可能重复，逐个算哈希
		if len(paths) < 2 {
			continue
		}
		for _, p := range paths {
			h, err := hashFile(p)
			if err != nil {
				continue
			}
			hashGroups[h] = append(hashGroups[h], p)
		}
	}

	var dups [][]string
	for _, group := range hashGroups {
		if len(group) > 1 {
			sort.Strings(group)
			dups = append(dups, group)
		}
	}
	sort.Slice(dups, func(i, j int) bool { return dups[i][0] < dups[j][0] })
	return dups, nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("用法: go-findup <目录>")
		return
	}
	if args[0] == "-h" || args[0] == "--help" {
		fmt.Println("go-findup 找出目录下内容重复的文件")
		return
	}
	dups, err := findDuplicates(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(dups) == 0 {
		fmt.Println("没找到重复文件")
		return
	}
	for i, group := range dups {
		fmt.Printf("=== 重复组 %d ===\n", i+1)
		for _, p := range group {
			fmt.Println("  " + p)
		}
	}
}
