package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileEntry struct {
	Path     string
	Contents string
}

var allDirs = make(map[string]bool)

func main() {
	root := "."
	outputFile := "project_summary.txt"
	var entries []FileEntry
	var structure strings.Builder

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(root, path)
		if d.IsDir() && (d.Name() == ".git" || d.Name() == "vendor") {
			return filepath.SkipDir
		}
		if d.IsDir() {
			allDirs[relPath] = true
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") {
			if strings.HasPrefix(relPath, "docs/") {
				return nil
			}
			if strings.HasPrefix(relPath, "test/") {
				return nil
			}
			//if strings.HasPrefix(relPath, "modules/") {
			//	if !(strings.HasPrefix(relPath, "modules/auth/") || strings.HasPrefix(relPath, "modules/cart/") || strings.HasPrefix(relPath, "modules/order/") || strings.HasPrefix(relPath, "modules/product/")) {
			//		return nil
			//	}
			//}
			entries = append(entries, FileEntry{
				Path:     relPath,
				Contents: readFile(path),
			})
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	printFullStructure(&structure)

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	file, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("📂 PROJECT STRUCTURE (.go only):\n\n")
	writer.WriteString(structure.String())
	writer.WriteString("\n\n📄 GO FILES CONTENT:\n\n")

	for _, entry := range entries {
		writer.WriteString(fmt.Sprintf("%s\n", entry.Path))
		writer.WriteString(strings.Repeat("-", 80) + "\n")
		writer.WriteString(entry.Contents + "\n\n")
	}
	writer.Flush()

	fmt.Println("✅ Project summary written to", outputFile)

	splitFileIntoThree(outputFile, "part1.txt", "part2.txt", "part3.txt")
}

func readFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("// ERROR reading %s: %v", path, err)
	}
	return string(content)
}

func printFullStructure(sb *strings.Builder) {
	dirs := make([]string, 0, len(allDirs))
	for dir := range allDirs {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)
	for _, dir := range dirs {
		levels := strings.Split(dir, string(os.PathSeparator))
		indent := ""
		for i, level := range levels {
			if level == "." {
				continue
			}
			indent = strings.Repeat("│   ", i)
			sb.WriteString(fmt.Sprintf("%s├── /%s\n", indent, level))
		}
	}
}

func splitFileIntoThree(input, output1, output2, output3 string) {
	data, err := os.ReadFile(input)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	total := len(lines)
	third := total / 3

	// Первая треть
	err = os.WriteFile(output1, []byte(strings.Join(lines[:third], "\n")), 0644)
	if err != nil {
		panic(err)
	}

	// Вторая треть
	err = os.WriteFile(output2, []byte(strings.Join(lines[third:2*third], "\n")), 0644)
	if err != nil {
		panic(err)
	}

	// Оставшиеся строки (третья часть)
	err = os.WriteFile(output3, []byte(strings.Join(lines[2*third:], "\n")), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("📝 File split into:", output1, ",", output2, "and", output3)
}
