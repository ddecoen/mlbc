package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Define the teamNameToNumber map
var teamNameToNumber = map[string]string{
	"Richmond":        "0",
	"Chattanooga":     "1",
	"Bangor":          "2",
	"Hartford":        "3",
	"Peninsula":       "4",
	"Midland":         "5",
	"Lynchburg":       "6",
	"Amarillo":        "7",
	"Peoria":          "8",
	"Sugar Land":      "9",
	"Omaha":           "10",
	"San Jose":        "11",
	"Lake Elsinore":   "12",
	"Oklahoma City":   "13",
	"Las Vegas":       "14",
	"Maine":           "15",
	"Mahoning Valley": "16",
	"Scranton":        "17",
	"Lousiville":      "18",
	"Bowie":           "19",
	"El Paso":         "20",
	"Rocket City":     "21",
	"Edmonton":        "22",
	"Altoona":         "23",
	"Wisconsin":       "24",
	"Idaho Falls":     "25",
	"Corpus Christi":  "26",
	"Portland":        "27",
	"Casper":          "28",
	"Billings":        "29",
}

// Function to add headers to each box score file
func addHeader(filePath string, outputFolder string) {
	// Read the content of the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) < 4 {
		fmt.Println("Invalid file format: Not enough lines.")
		return
	}

	var teamsLineIndex int
	for i, line := range lines {
		if strings.Contains(line, "FINAL") {
			teamsLineIndex = i + 2
			break
		}
	}

	if teamsLineIndex >= len(lines) {
		fmt.Println("Could not find teams line in file.")
		return
	}

	filename := filepath.Base(filePath)
	day := strings.Split(filename, "_")[0]

	teamXLine := strings.Fields(lines[teamsLineIndex])
	teamYLine := strings.Fields(lines[teamsLineIndex+1])

	if len(teamXLine) < 1 || len(teamYLine) < 1 {
		fmt.Println("Invalid file format: Not enough fields.")
		return
	}

	// Find the index of the first digit in the line
	findFirstDigit := func(str string) int {
		for i, c := range str {
			if c >= '0' && c <= '9' {
				return i
			}
		}
		return -1
	}

	xIndex := findFirstDigit(lines[teamsLineIndex])
	yIndex := findFirstDigit(lines[teamsLineIndex+1])

	if xIndex == -1 || yIndex == -1 {
		fmt.Println("Could not find team names in file.")
		return
	}

	teamXName := strings.TrimSpace(lines[teamsLineIndex][:xIndex])
	teamYName := strings.TrimSpace(lines[teamsLineIndex+1][:yIndex])

	// Use the teamNameToNumber map to fetch the corresponding numbers
	teamXNum, exists := teamNameToNumber[teamXName]
	if !exists {
		teamXNum = "UNKNOWN"
	}

	teamYNum, exists := teamNameToNumber[teamYName]
	if !exists {
		teamYNum = "UNKNOWN"
	}

	header := fmt.Sprintf("Team%s vs Team%s (Day %s)", teamXNum, teamYNum, day)
	newContent := header + "\n" + string(content)

	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		err = os.MkdirAll(outputFolder, 0755)
		if err != nil {
			panic(err)
		}
	}

	newFilePath := filepath.Join(outputFolder, filename)
	err = ioutil.WriteFile(newFilePath, []byte(newContent), 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Loop through the "boxes" directory
	err := filepath.Walk("boxes", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") && !strings.Contains(info.Name(), "team_numbers") {
			addHeader(path, "new_boxes")
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
	}
}
