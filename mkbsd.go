/*
Created on 01 Oct 2024
@author Sai Sumanth
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	const downloadUrl = "https://storage.googleapis.com/panels-api/data/20240916/media-1a-i-p~s"

	// get the JSON data
	jsonData, err := makeNetworkRequest(downloadUrl)
	if err != nil {
		log.Fatalf("error while fetching data: %v", err)
	}
	// Extract HD urls
	hdUrls := extractDownloadUrls(jsonData)

	log.Printf("Found %d images, started downloading.....", len(hdUrls))
	downloadImages(hdUrls)
	fmt.Println("Time taken:", time.Since(start))
}

func downloadImages(urlsMap map[string]string) {
	// create a directory "downloads" at cwd
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get current working directory details")
	}
	downloadsDirectory := filepath.Join(currentWorkingDir, "downloads")
	if err := os.MkdirAll(downloadsDirectory, os.ModePerm); err != nil {
		log.Fatal("failed to create downloads directory:", err)
	}
	wg := &sync.WaitGroup{}
	for k, v := range urlsMap {
		if fileExtension, err := getExtension(v); err == nil {
			imagePath := filepath.Join(downloadsDirectory, k+fileExtension)
			// log.Println(imagePath)
			wg.Add(1)
			go func(imageDownloadUrl, path string) {
				defer wg.Done()
				downloadImage(imageDownloadUrl, path)
			}(v, imagePath)
		}
	}
	wg.Wait()
}

// returns the image extension obtained from given image download url
func getExtension(imageUrl string) (string, error) {
	parsedUrl, err := url.Parse(imageUrl)
	if err != nil {
		return "", err
	}
	return path.Ext(parsedUrl.Path), nil
}

// extract HD download urls from
func extractDownloadUrls(jsonData map[string]interface{}) map[string]string {
	downloadUrls := make(map[string]string)
	for k, v := range jsonData {
		eachWallpaperData, ok := v.(map[string]interface{})
		if !ok || len(eachWallpaperData) == 0 {
			continue
		}

		if hd, exists := eachWallpaperData["dhd"]; exists {
			downloadUrls[k] = hd.(string)
		}
	}
	return downloadUrls
}

// makes a network request to fetch JSON Data
func makeNetworkRequest(url string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	log.Println("Making network request to url:", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("status code is not 200")
	}

	var response map[string]interface{}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse body: %w", err)
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to extract the data from json response")
	}

	return data, nil
}

// downlods the image file from given url and saves it locally at specified path
func downloadImage(imageUrl string, downloadPath string) {
	imageResponse, err := http.Get(imageUrl)
	if err != nil {
		log.Println("error downloading image", imageUrl)
		return
	}
	if imageResponse.StatusCode != 200 {
		return
	}

	// create file
	file, err := os.Create(downloadPath)
	if err != nil {
		log.Println("failed to create file")
		return
	}
	defer file.Close()

	_, err = io.Copy(file, imageResponse.Body)
	if err != nil {
		log.Println("failed to save image at given path")
		return
	}

	log.Println("Image downloaded successfully")
}
