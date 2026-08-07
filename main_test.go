package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindDuplicates(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "sub"), 0755)
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("same"), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("same"), 0644)
	os.WriteFile(filepath.Join(dir, "c.txt"), []byte("different"), 0644)
	os.WriteFile(filepath.Join(dir, "sub", "d.txt"), []byte("same"), 0755)

	dups, err := findDuplicates(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(dups) != 1 {
		t.Fatalf("期望找到 1 组重复, 得到 %d", len(dups))
	}
	if len(dups[0]) != 3 {
		t.Errorf("重复组应有 3 个文件, 得到 %d", len(dups[0]))
	}
}

func TestFindDuplicatesNone(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("y"), 0644)
	dups, err := findDuplicates(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(dups) != 0 {
		t.Errorf("不应有重复, 得到 %v", dups)
	}
}

func TestHashFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "f")
	os.WriteFile(p, []byte("dup"), 0644)
	h1, _ := hashFile(p)
	h2, _ := hashFile(p)
	if h1 != h2 {
		t.Error("同一文件哈希应一致")
	}
}
