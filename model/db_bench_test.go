package model

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func BenchmarkGetPreCommands(b *testing.B) {
	homeDir := os.Getenv("HOME")

	baseTimeFolder := strconv.Itoa(int(time.Now().Unix()))
	baseFolder := baseTimeFolder + "-withPre"
	InitFolder(baseFolder)

	// Create test directory structure
	err := os.MkdirAll(fmt.Sprintf("%s/%s", homeDir, COMMAND_STORAGE_FOLDER), 0755)
	if err != nil {
		b.Fatal("Failed to create test directory:", err)
	}

	// Create test pre-commands file with sample data
	preFilePath := fmt.Sprintf("%s/%s", homeDir, COMMAND_PRE_STORAGE_FILE)

	now := time.Now().UnixNano()
	nowStr := strconv.Itoa(int(now))
	testData := []string{
		`{"shell":"bash","sid":123,"cmd":"ls -la","main":"","hn":"localhost","un":"user1","t":"2024-12-13T21:57:02.204345+08:00","et":"0001-01-01T00:00:00Z","result":0,"phase":0}` + "\t" + nowStr,
		`{"shell":"zsh","sid":456,"cmd":"cd /tmp","main":"","hn":"localhost","un":"user2","t":"2024-12-13T21:57:02.204345+08:00","et":"0001-01-01T00:00:00Z","result":0,"phase":0}` + "\t" + nowStr,
		`{"shell":"bash","sid":789,"cmd":"grep pattern","main":"","hn":"localhost","un":"user3","t":"2024-12-13T21:57:02.204345+08:00","et":"0001-01-01T00:00:00Z","result":0,"phase":0}` + "\t" + nowStr,
		`{"shell":"fish","sid":101,"cmd":"cat file.txt","main":"","hn":"localhost","un":"user1","t":"2024-12-13T21:57:02.204345+08:00","et":"0001-01-01T00:00:00Z","result":0,"phase":0}` + "\t" + nowStr,
	}

	f, err := os.Create(preFilePath)
	if err != nil {
		b.Fatal("Failed to create test file:", err)
	}
	for _, line := range testData {
		_, err := f.WriteString(line + "\n")
		if err != nil {
			f.Close()
			b.Fatal("Failed to write test data:", err)
		}
	}
	f.Close()

	// Cleanup function
	b.Cleanup(func() {
		os.Remove(preFilePath)
		os.RemoveAll(fmt.Sprintf("%s/%s", homeDir, COMMAND_STORAGE_FOLDER))
	})

	// Run benchmark
	ctx := context.Background()
	b.ResetTimer() // Start timing from here

	for i := 0; i < b.N; i++ {
		commands, err := GetPreCommands(ctx)
		if err != nil {
			b.Fatal("Benchmark failed:", err)
		}
		if len(commands) == 0 {
			b.Fatal("Expected non-empty commands slice")
		}
	}
}

// Benchmark with different file sizes
func BenchmarkGetPreCommandsWithDifferentSizes(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			// Setup test data
			homeDir := os.Getenv("HOME")
			baseTimeFolder := strconv.Itoa(int(time.Now().Unix()))
			baseFolder := baseTimeFolder + "-withPre"
			InitFolder(baseFolder)

			// Create test directory structure
			err := os.MkdirAll(fmt.Sprintf("%s/%s", homeDir, COMMAND_STORAGE_FOLDER), 0755)
			if err != nil {
				b.Fatal("Failed to create test directory:", err)
			}

			// Create test pre-commands file with varying amounts of data
			preFilePath := fmt.Sprintf("%s/%s", homeDir, COMMAND_PRE_STORAGE_FILE)
			f, err := os.Create(preFilePath)
			if err != nil {
				b.Fatal("Failed to create test file:", err)
			}

			// Write test data
			for i := 0; i < size; i++ {
				now := time.Now()
				line := fmt.Sprintf(`{"shell":"bash","sid":%d,"cmd":"command%d","main":"","hn":"localhost","un":"user%d","t":"%s","et":"0001-01-01T00:00:00Z","result":0,"phase":0}`,
					i, i, i%10, now.Format(time.RFC3339))
				_, err := f.WriteString(line + "\t" + strconv.Itoa(int(now.UnixNano())) + "\n")
				if err != nil {
					f.Close()
					b.Fatal("Failed to write test data:", err)
				}
			}
			f.Close()

			// Cleanup function
			b.Cleanup(func() {
				os.Remove(preFilePath)
				os.RemoveAll(fmt.Sprintf("%s/%s", homeDir, COMMAND_STORAGE_FOLDER))
			})

			// Run benchmark
			ctx := context.Background()
			b.ResetTimer() // Start timing from here

			for i := 0; i < b.N; i++ {
				commands, err := GetPreCommands(ctx)
				if err != nil {
					b.Fatal("Benchmark failed:", err)
				}
				if len(commands) == 0 {
					b.Fatal("Expected non-empty commands slice")
				}
			}
		})
	}
}

func BenchmarkGetPreCommandsTree(b *testing.B) {
	homeDir := os.Getenv("HOME")

	baseTimeFolder := strconv.Itoa(int(time.Now().Unix()))
	baseFolder := baseTimeFolder + "-withPre"
	InitFolder(baseFolder)

	// Create test directory structure
	err := os.MkdirAll(fmt.Sprintf("%s/%s", homeDir, COMMAND_STORAGE_FOLDER), 0755)
	if err != nil {
		b.Fatal("Failed to create test directory:", err)
	}

	// Create test pre-commands file with sample data
	preFilePath := fmt.Sprintf("%s/%s", homeDir, COMMAND_PRE_STORAGE_FILE)

	now := time.Now().UnixNano()
	nowStr := strconv.Itoa(int(now))
	testData := []string{
		`{"shell":"bash","sid":123,"cmd":"ls -la","main":"","hn":"localhost","un":"user1","t":"2024-12-13T21:57:02.204345+08:00","et":"0001-01-01T00:00:00Z","result":0,"phase":0}` + "\t" + nowStr,
		`{"shell":"zsh","sid":456,"cmd":"cd /tmp","main":"","hn":"localhost","un":"user2","t":"2024-12-13T21:57:02.204345+08:00","et":"0001-01-01T00:00:00Z","result":0,"phase":0}` + "\t" + nowStr,
		`{"shell":"bash","sid":789,"cmd":"grep pattern","main":"","hn":"localhost","un":"user3","t":"2024-12-13T21:57:02.204345+08:00","et":"0001-01-01T00:00:00Z","result":0,"phase":0}` + "\t" + nowStr,
		`{"shell":"fish","sid":101,"cmd":"cat file.txt","main":"","hn":"localhost","un":"user1","t":"2024-12-13T21:57:02.204345+08:00","et":"0001-01-01T00:00:00Z","result":0,"phase":0}` + "\t" + nowStr,
	}

	f, err := os.Create(preFilePath)
	if err != nil {
		b.Fatal("Failed to create test file:", err)
	}
	for _, line := range testData {
		_, err := f.WriteString(line + "\n")
		if err != nil {
			f.Close()
			b.Fatal("Failed to write test data:", err)
		}
	}
	f.Close()

	// Cleanup function
	b.Cleanup(func() {
		os.Remove(preFilePath)
		os.RemoveAll(fmt.Sprintf("%s/%s", homeDir, COMMAND_STORAGE_FOLDER))
	})

	// Run benchmark
	ctx := context.Background()
	b.ResetTimer() // Start timing from here

	for i := 0; i < b.N; i++ {
		tree, err := GetPreCommandsTree(ctx)
		if err != nil {
			b.Fatal("Benchmark failed:", err)
		}
		if len(tree) == 0 {
			b.Fatal("Expected non-empty tree")
		}
	}
}

// Benchmark with different file sizes
func BenchmarkGetPreCommandsTreeWithDifferentSizes(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			// Setup test data
			homeDir := os.Getenv("HOME")
			baseTimeFolder := strconv.Itoa(int(time.Now().Unix()))
			baseFolder := baseTimeFolder + "-withPre"
			InitFolder(baseFolder)

			// Create test directory structure
			err := os.MkdirAll(fmt.Sprintf("%s/%s", homeDir, COMMAND_STORAGE_FOLDER), 0755)
			if err != nil {
				b.Fatal("Failed to create test directory:", err)
			}

			// Create test pre-commands file with varying amounts of data
			preFilePath := fmt.Sprintf("%s/%s", homeDir, COMMAND_PRE_STORAGE_FILE)
			f, err := os.Create(preFilePath)
			if err != nil {
				b.Fatal("Failed to create test file:", err)
			}

			// Write test data
			for i := 0; i < size; i++ {
				now := time.Now()
				line := fmt.Sprintf(`{"shell":"bash","sid":%d,"cmd":"command%d","main":"","hn":"localhost","un":"user%d","t":"%s","et":"0001-01-01T00:00:00Z","result":0,"phase":0}`,
					i, i, i%10, now.Format(time.RFC3339))
				_, err := f.WriteString(line + "\t" + strconv.Itoa(int(now.UnixNano())) + "\n")
				if err != nil {
					f.Close()
					b.Fatal("Failed to write test data:", err)
				}
			}
			f.Close()

			// Cleanup function
			b.Cleanup(func() {
				os.Remove(preFilePath)
				os.RemoveAll(fmt.Sprintf("%s/%s", homeDir, COMMAND_STORAGE_FOLDER))
			})

			// Run benchmark
			ctx := context.Background()
			b.ResetTimer() // Start timing from here

			for i := 0; i < b.N; i++ {
				tree, err := GetPreCommandsTree(ctx)
				if err != nil {
					b.Fatal("Benchmark failed:", err)
				}
				if len(tree) == 0 {
					b.Fatal("Expected non-empty tree")
				}
			}
		})
	}
}
