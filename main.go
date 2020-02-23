package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dhowden/tag"
	"github.com/fatih/color"
)

func main() {
	args := os.Args[1:]
	var dir string

	if len(args) > 0 {
		winPath := filepath.Dir(args[0])
		basePath := filepath.Base(args[0])
		dir = filepath.Join(winPath, basePath)
	} else {
		fakedir, err := os.Getwd()
		dir = fakedir
		if err != nil {
			log.Fatal("Failed to ascertain working directory :", err)
		}
	}

	files, err := os.Open(dir)

	if err != nil {
		fmt.Println("Failed to open directory: ", err)
		return
	}
	defer files.Close()
	list, err := files.Readdirnames(0)

	if err != nil {
		fmt.Println("Failed to read dirnames: ", err)
		return
	}

	audioObjs := make(AudioFiles, 0, len(list))
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	megenta := color.New(color.FgMagenta)
	red := color.New(color.FgRed)

	for _, name := range list {
		ext := filepath.Ext(name)
		if isAudioFile(ext) {
			tags := parseAudioFile(name, dir)
			a := AudioFile{
				Title:        format(tags.Title()),
				Artist:       format(tags.Artist()),
				OriginalName: name,
				Ext:          strings.ToLower(ext),
			}
			if a.Title == "" || a.Artist == "" {
				fmt.Printf("[")
				red.Printf("-")
				fmt.Printf("] ")
				megenta.Printf(name)
				fmt.Printf(" does not have enough meta data\n")
				continue
			}
			newName := a.Title + " - " + a.Artist + ext
			if a.OriginalName == newName {
				fmt.Printf("[")
				red.Printf("-")
				fmt.Printf("] ")
				megenta.Printf(name)
				fmt.Printf(" has already a good name\n")
				continue
			}
			fmt.Printf("[")
			green.Printf("+")
			fmt.Printf("] ")
			megenta.Printf(name)
			fmt.Printf(" will be renamed to ")
			green.Printf(newName)
			fmt.Printf("\n")

			audioObjs = append(audioObjs, a)
		}
	}
	if len(audioObjs) == 0 {
		fmt.Println("Not enough files to rename")
		return
	}
	if !prompt("Do you want to continue? (y/N)") {
		return
	}
	for _, obj := range audioObjs {
		if fileExists(obj.OriginalName, dir) {
			if err := os.Rename(filepath.Join(dir, obj.OriginalName), filepath.Join(dir, obj.Title+" - "+obj.Artist+obj.Ext)); err != nil {
				log.Fatal("Failed to rename", obj.OriginalName)
			}
			fmt.Printf(obj.OriginalName)
			yellow.Printf(" => ")
			fmt.Printf(obj.Title + " - " + obj.Artist + obj.Ext + "\n")
		} else {
			fmt.Println(obj.OriginalName, "does not exist")
		}
	}
	if !prompt("Do you want to undo all previous changes? (y/N)") {
		return
	}
	for _, obj := range audioObjs {
		if fileExists(obj.Title+" - "+obj.Artist+obj.Ext, dir) {
			if err := os.Rename(filepath.Join(dir, obj.Title+" - "+obj.Artist+obj.Ext), filepath.Join(dir, obj.OriginalName)); err != nil {
				log.Fatal("Failed to rename", obj.OriginalName)
			}
			fmt.Printf(obj.Title + " - " + obj.Artist + obj.Ext)
			yellow.Printf(" => ")
			fmt.Printf(obj.OriginalName + "\n")
		} else {
			fmt.Println(obj.OriginalName, "does not exist")
		}
	}
}

// Parse meta data of an audio file
func parseAudioFile(str string, wd string) tag.Metadata {
	file, err := os.Open(filepath.Join(wd, str))
	if err != nil {
		log.Fatal("Failed to open audio file ", str, err)
	}
	defer file.Close()
	m, err := tag.ReadFrom(file)
	if err != nil {
		log.Fatal("Failed to read from file ", str, err)
	}
	return m
}

// Returns true when audio file
func isAudioFile(ext string) bool {
	if ext == ".flac" || ext == ".mp3" || ext == ".wav" || ext == ".m4a" || ext == ".ogg" || ext == ".acc" || ext == ".alac" {
		return true
	}
	return false
}

// Attempt to format
func format(input string) string {
	str := strings.Title(strings.ToLower(input))
	re := regexp.MustCompile(`'[A-Za-z]( |\z)`)
	str = re.ReplaceAllStringFunc(str, strings.ToLower)

	// re := regexp.MustCompile(`(?m)\b\w\B`)
	// str := strings.ToLower(input)
	// str = re.ReplaceAllStringFunc(str, strings.ToUpper)

	return str
}

// Returns the answer of a question as bool
func prompt(question string) bool {
	var s string
	fmt.Printf(question)
	if _, err := fmt.Scan(&s); err != nil {
		log.Fatal("Failed to scan: ", err)
	}
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "y" || s == "yes" {
		return true
	}
	return false
}

// Checks if the file exists
func fileExists(filename string, wd string) bool {
	info, err := os.Stat(filepath.Join(wd, filename))
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// AudioFiles is an array of Audiofile
type AudioFiles []AudioFile

// AudioFile is a struct
type AudioFile struct {
	Title        string
	Artist       string
	OriginalName string
	Ext          string
}
