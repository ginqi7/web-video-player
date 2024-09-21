package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"net/url"
)

type Navigation struct {
	Path string
	Name string
}

type Directory struct {
	RelativePath string
	Name         string
}

type File struct {
	Name string
}

func parseNavigation(url string) []Navigation {

	names := strings.Split(url, "/")
	navigations := []Navigation{}
	path := "/listing"
	for _, name := range names {
		navigation := Navigation{
			Path: path + name,
			Name: name,
		}
		path += name + "/"
		navigations = append(navigations, navigation)
	}

	return navigations
}

func getDirectories(url string) []Directory {
	currentDir := getBasePath() + url
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	directories := []Directory{}
	for _, entry := range entries {
		if entry.IsDir() {
			directory := Directory{
				RelativePath: url,
				Name:         entry.Name(),
			}
			directories = append(directories, directory)
		}
	}
	return directories
}

func isVideoFile(filePath string) bool {
	videoExtensions := []string{
		".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm", ".mpeg", ".mpg"}
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, videoExt := range videoExtensions {
		if ext == videoExt {
			return true
		}
	}
	return false
}

func getFiles(url string) []File {
	currentDir := getBasePath() + url

	entries, err := os.ReadDir(currentDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	files := []File{}
	for _, entry := range entries {
		if isVideoFile(entry.Name()) {
			file := File{
				Name: entry.Name(),
			}
			files = append(files, file)
		}
	}
	return files
}

func getBasePath() string {
	varName := "WEB_VIDEO_PLAYER_BASE_PATH"
	return os.Getenv(varName)
}

func openFile(path string) (*os.File, error) {
	full_path := getBasePath() + "/" + path
	full_path, err := url.QueryUnescape(full_path)
	if err != nil {
		fmt.Println("Error decoding URL:", err)
		return nil, err
	}
	return os.Open(full_path)
}
