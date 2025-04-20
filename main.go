package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var folder = "/home/brunopoiano/Documents/Pessoal/subtitle-extract/videos"

var videoExtensions = map[string]bool{
	".mp4":  true,
	".mkv":  true,
	".avi":  true,
	".mov":  true,
	".wmv":  true,
	".flv":  true,
	".webm": true,
}

func main() {

	getItemsFromFolder(folder)

}

func getItemsFromFolder(location string) {

	items, err := os.ReadDir(location)
	if err != nil {
		println("Error going to folder")
	}

	for _, item := range items {

		if !item.IsDir() {
			ext := strings.ToLower(filepath.Ext(item.Name()))

			if videoExtensions[ext] {
				runningEmbedSubtitleCheck(location, item)
			}
		} else {

			new_location := fmt.Sprintf("%s/%s", location, item.Name())
			println(new_location)
			getItemsFromFolder(new_location)

		}

	}
}

func runningEmbedSubtitleCheck(folder string, item os.DirEntry) {

	name := item.Name()
	command := fmt.Sprintf("ffmpeg -i '%s/%s' 2>&1 | grep Stream.*Subtitle", folder, name)

	println("command:", command)

	cmd := exec.Command("bash", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}

	items := strings.Split(string(output), "\n")

	for _, item := range items {

		re := regexp.MustCompile(`\d+:\d+\([a-zA-Z0-9]{3}\)`)
		matches := re.FindAllString(item, -1)
		if len(matches) == 0 {
			continue
		}
		runningExtractSubtitle(folder, name, matches[0])

	}

}

func runningExtractSubtitle(folder, name, sub_language string) {

	re := regexp.MustCompile(`[()]`)

	sub_div := strings.Split(sub_language, "(")
	sub_num := string(sub_div[0])
	sub_len := re.ReplaceAllString(sub_div[1], "")

	println(string(sub_num))
	println(string(sub_len))

	full_path := fmt.Sprintf("%s/%s", folder, name)

	command := fmt.Sprintf("ffmpeg -i '%s' -map %s -c:s srt '%s.%s.srt'", full_path, sub_num, full_path, sub_len)

	cmd := exec.Command("bash", "-c", command)

	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}

}
