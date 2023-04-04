package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

const MAX_API_COUNT_IMAGES = 50
const API_URL = "https://dog.ceo/api/breeds/image/random"
const PATH_PREFIX = "images_dogs/"

var urls = []string{}
var myClient = &http.Client{}

var wg sync.WaitGroup

func getJson(url string, target *RandomDogAPIResponse) error {
	response, err := myClient.Get(url)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	json.Unmarshal(data, &target)

	return nil
}

type RandomDogAPIResponse struct {
	Message []string `json:"message"`
	Status  string   `json:"status"`
}

func partional_download_dog_images(count_images int) {
	response := RandomDogAPIResponse{}
	getJson(fmt.Sprintf("%s/%d", API_URL, count_images), &response)

	for _, url := range response.Message {
		urls = append(urls, url)
	}
}

func errorHandler(message string) error {
	println("ERROR: " + message)
	defer wg.Done()
	return errors.New(message)
}

func load_image(url string, _id int) error {
	response, err := http.Get(url)
	if err != nil {
		return errorHandler("Url invalide")
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errorHandler("Received non 200 response code")
	}
	fileName := PATH_PREFIX + uuid.New().String() + ".jpg"
	file, err := os.Create(fileName)
	if err != nil {
		return errorHandler("File not created")
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return errorHandler("File not saved")
	}
	println(_id)
	defer wg.Done()
	return nil
}

func main() {
	var count_images = 100
	for count_images > 0 {
		partional_count := count_images
		if partional_count > MAX_API_COUNT_IMAGES {
			partional_count = MAX_API_COUNT_IMAGES
		}
		count_images -= partional_count
		partional_download_dog_images(partional_count)
	}
	start := time.Now()

	wg.Add(len(urls))
	println(len(urls))
	for i, url := range urls {
		go load_image(url, i)
	}
	wg.Wait()
	fmt.Printf("took %v\n", time.Since(start))
}
