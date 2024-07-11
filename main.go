package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	moviesDir = "/home/eric/Downloads/movies/"
	apiKey    = "3cb8847a0a44b51f501e9b813c14b52b"
)

func getMovieInfo(movieName string) {

	baseURL := "https://api.themoviedb.org/3/search/movie"
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	params := url.Values{}
	params.Add("api_key", apiKey)
	params.Add("query", movieName)
	reqURL.RawQuery = params.Encode()

	resp, err := http.Get(reqURL.String())
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Received non-200 response code", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var movieResponse Movie
	err = json.Unmarshal(body, &movieResponse)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return
	}

	for _, result := range movieResponse.Results {
		jsonData, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			fmt.Println("Error marshaling to JSON:", err)
			return
		}
		fmt.Println(string(jsonData))
	}
}

func main() {

	files, err := os.ReadDir(moviesDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	examples := make([]string, 0, len(files))

	for _, file := range files {
		if !file.IsDir() {
			examples = append(examples, file.Name())
		}
	}

	re := regexp.MustCompile(`^([^()0-9]*)(\d{1,2})? ?(\(?\d{4}\)?)?.*`)
	formattedExamples := make([]string, len(examples))

	for i, example := range examples {
		example = strings.TrimSuffix(example, filepath.Ext(example))

		if strings.Contains(example, ":") || (len(example) > 0 && !strings.ContainsAny(example, "0123456789")) {
			formattedExamples[i] = example
			continue
		}

		matches := re.FindStringSubmatch(example)
		if matches != nil {
			formattedTitle := matches[1]
			part := matches[2]
			year := matches[3]

			formattedTitle = strings.ReplaceAll(formattedTitle, ".", " ")
			formattedTitle = strings.ReplaceAll(formattedTitle, "/", " ")

			formattedTitle = strings.TrimSpace(formattedTitle)

			if part != "" {
				formattedTitle = fmt.Sprintf("%s %s", formattedTitle, part)
			}

			if year != "" {
				formattedTitle = fmt.Sprintf("%s %s", formattedTitle, year)
			}

			formattedTitle = strings.ReplaceAll(formattedTitle, "(", "")
			formattedTitle = strings.ReplaceAll(formattedTitle, ")", "")

			formattedTitle = strings.TrimSpace(strings.ReplaceAll(formattedTitle, "\u00A0", " "))

			formattedExamples[i] = formattedTitle
		} else {
			fmt.Printf("Warning: Couldn't format title: %s\n", example)
			formattedExamples[i] = example
		}
	}

	for _, movieName := range formattedExamples {
		getMovieInfo(movieName)
	}

}
