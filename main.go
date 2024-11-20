package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	source := "../afero"
	destination := "./"
	destination = filepath.Clean(destination + "/" + filepath.Base(source))

	err := copyFolderWithBufio(source, destination, 2*1024*1024)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Folder copied successfully!")
	}
}

func copyFolderWithBufio(src string, dst string, maxSize int64) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, relativePath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		if fileInfo.Size() > maxSize {
			fmt.Printf("Skipping file: %s (size: %d bytes)\n", path, fileInfo.Size())
			return nil
		}

		return copyFileWithBufio(path, destPath)
	})
}

func copyFileWithBufio(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	reader := bufio.NewReader(srcFile)
	writer := bufio.NewWriter(destFile)

	buffer := make([]byte, 4096)
	for {
		n, err := reader.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			return err
		}
		if n == 0 {
			break
		}
		if _, writeErr := writer.Write(buffer[:n]); writeErr != nil {
			return writeErr
		}
	}

	return writer.Flush()
}
