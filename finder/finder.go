package finder

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
)

//Dir represents the root directory for which the search is made
type Dir struct {
	Name   string
	Path   string
	SubDir []*Dir
	Files  []*File
}

//File represents the physical file within each directory and has matches of the searched text
type File struct {
	Name    string
	Path    string
	Matches []*Match
}

//Match represents each match of the searched text
type Match struct {
	Line int
	Text string
}

//Config represents the configuration parameters which drives the application
type Config struct {
	ExcludeDirectories []string `json:"excludeDirectories"`
	ExcludeFiles       []string `json:"excludeFiles"`
	AllowedExtensions  []string `json:"allowedExtensions"`
	MatchCase          bool     `json:"matchCase"`
	MatchFullWord      bool     `json:"matchFullWord"`
}

type scan struct {
	text string
	line int
}

var config Config
var matcher Matcher

const root = 1

//NewDir creates a new directory with name and path
func NewDir(name string, path string) *Dir {
	return &Dir{Name: name, Path: path, Files: make([]*File, 0), SubDir: (make([]*Dir, 0))}
}

//Find searches the directory and the nested directories for the search string
func (dir *Dir) Find(search string) {

	var wg sync.WaitGroup

	wg.Add(1)

	go dir.Scan(search, &wg)

	wg.Wait()

	dir.print()

}

func (dir *Dir) getSubDirs() []*Dir {
	files, err := ioutil.ReadDir(dir.Path)

	if err != nil {
		log.Fatal(err)
	}

	sub := make([]*Dir, 0)
	for _, fileInfo := range files {
		if fileInfo.IsDir() && !contains(config.ExcludeDirectories, fileInfo.Name()) {
			subDir := NewDir(fileInfo.Name(), filepath.Join(dir.Path, fileInfo.Name()))
			sub = append(sub, subDir)
		}
	}

	return sub
}

//Scan scans the sub directories
func (dir *Dir) Scan(search string, wg *sync.WaitGroup) {

	defer wg.Done()

	if contains(config.ExcludeDirectories, dir.Name) {
		return
	}

	files, err := ioutil.ReadDir(dir.Path)

	if err != nil {
		log.Fatal(err)
	}

	var wg1 sync.WaitGroup

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			wg.Add(1)
			go dir.scanSubDir(search, wg, fileInfo)

		} else {
			if contains(config.AllowedExtensions, filepath.Ext(fileInfo.Name())) {
				wg1.Add(1)
				go dir.scanFile(search, &wg1, fileInfo)
			}
		}
	}

	wg1.Wait()

}

func (dir *Dir) scanSubDir(search string, wg *sync.WaitGroup, dirInfo os.FileInfo) {
	subDir := NewDir(dirInfo.Name(), filepath.Join(dir.Path, dirInfo.Name()))
	subDir.Scan(search, wg)

	if len(subDir.Files) > 0 {
		dir.SubDir = append(dir.SubDir, subDir)

	}
}

func (dir *Dir) scanFile(search string, wg *sync.WaitGroup, fileInfo os.FileInfo) {
	file := newFile(fileInfo.Name(), filepath.Join(dir.Path, fileInfo.Name()))
	file.find(search, wg)

	if len(file.Matches) > 0 {
		dir.Files = append(dir.Files, file)
	}
}

//Print prints all the 'matches' found inside the root directory to the console
func (dir *Dir) print() {
	matches := dir.getTotalMatches()

	if matches > 0 {
		color.HiCyan("Total matches : " + strconv.Itoa(matches))
		dir.printIndent(1)
	} else {
		color.Red("No matches found")
	}

}

func (dir *Dir) printIndent(level int) {

	color.Yellow(getSpaces(level) + dir.Name)

	if len(dir.SubDir) > 0 {
		for _, dir := range dir.SubDir {
			dir.printIndent(level + 1)
		}
	}

	if len(dir.Files) > 0 {
		for _, file := range dir.Files {
			color.Green(getSpaces(level) + " " + file.Name)

			for _, match := range file.Matches {
				fmt.Println(getSpaces(level) + " " + match.Text + " " + strconv.Itoa(match.Line))
			}

			fmt.Println(" ")
		}
	}
}

func (dir *Dir) getTotalMatches() int {
	total := 0

	for _, subDir := range dir.SubDir {
		total += subDir.getTotalMatches()
	}
	for _, file := range dir.Files {
		total += len(file.Matches)
	}
	return total
}

func newFile(name string, path string) *File {
	return &File{Name: name, Path: path, Matches: make([]*Match, 0)}
}

func (file *File) find(search string, wg *sync.WaitGroup) {

	defer wg.Done()
	if contains(config.ExcludeFiles, file.Name) {
		return
	}

	f, err := os.Open(file.Path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	file.scan(f, search)
}

func (file *File) scan(r io.Reader, search string) {
	workQueue := make(chan scan)

	var wg1 sync.WaitGroup
	line := 1
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			workQueue <- scan{line: line, text: scanner.Text()}
			line++
		}
		close(workQueue)
	}()

	for i := 0; i < 100; i++ {
		wg1.Add(1)
		go file.process(workQueue, search, &wg1)
	}

	wg1.Wait()
}

func (file *File) process(queue chan scan, search string, wg *sync.WaitGroup) {

	defer wg.Done()
	for line := range queue {
		for _, match := range matcher.Match(line.text, search) {
			file.Matches = append(file.Matches, &Match{Line: line.line, Text: match})
		}
	}
}

//Init must be called to initialize the app with the config file
func Init(cfgFile string) {
	setConfig(cfgFile)
	setMatcher()
}

func setConfig(fileName string) {

	configFile, err := os.Open(fileName)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
}

func setMatcher() {
	if config.MatchCase && config.MatchFullWord {
		matcher = FullMatcher{}
	}
	if !config.MatchCase && config.MatchFullWord {
		matcher = FullMatcherIgnoreCase{}
	}
	if config.MatchCase && !config.MatchFullWord {
		matcher = PartialMatcher{}
	}
	if !config.MatchCase && !config.MatchFullWord {
		matcher = PartialMatcherIgnoreCase{}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getSpaces(level int) string {
	var sb strings.Builder
	for i := 0; i < level; i++ {
		sb.WriteString("   ")
	}

	return sb.String()
}

func batch(dir []*Dir, batch int) [][]*Dir {
	var divided [][]*Dir

	for i := 0; i < len(dir); i += batch {
		end := i + batch

		if end > len(dir) {
			end = len(dir)
		}

		divided = append(divided, dir[i:end])
	}
	return divided
}
