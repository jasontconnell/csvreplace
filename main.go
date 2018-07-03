package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var csvfilename string
var templatefilename string
var mode string
var outputfilepattern string

func init() {
	flag.StringVar(&csvfilename, "c", "", "Please provide the csv filename")
	flag.StringVar(&templatefilename, "t", "", "Please provide the template filename")
	flag.StringVar(&mode, "m", "stdio", "Mode for the output. Either stdio or file")
	flag.StringVar(&outputfilepattern, "o", "{0}.txt", "For file output mode, the filename pattern for files")
}

func main() {
	flag.Parse()

	baseDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	absCsvFilename := getAbsFile(baseDir, csvfilename)
	absTemplateFilename := getAbsFile(baseDir, templatefilename)
	absOutputFilePattern := getAbsFile(baseDir, outputfilepattern)

	lines := readCsv(absCsvFilename)
	template := readTemplate(absTemplateFilename)

	if mode == "stdio" {
		updated := processLines(lines, template)
		fmt.Println(updated)
	} else {
		processLinesFileOutput(lines, template, absOutputFilePattern)
	}

}

func getAbsFile(baseDir, path string) string {
	abspath := path
	if !filepath.IsAbs(abspath) {
		abspath = filepath.Join(baseDir, abspath)
	}
	return abspath
}

func readCsv(filename string) (lines []string) {
	if f, err := os.Open(filename); err == nil {
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			var txt = scanner.Text()
			lines = append(lines, txt)
		}
	}
	return
}

func readTemplate(filename string) (template string) {
	b, _ := ioutil.ReadFile(filename)
	template = string(b)
	return
}

func processLines(lines []string, template string) (fullReplace string) {
	for _, line := range lines {
		fullReplace = fullReplace + processLine(line, template)
	}

	return
}

func processLine(line string, template string) string {
	s := strings.Split(line, ",")
	tmp := template
	for i, _ := range s {
		index := fmt.Sprintf("{%v}", i)
		tmp = strings.Replace(tmp, index, s[i], -1)
	}

	return tmp
}

func processLinesFileOutput(lines []string, template string, outputfilenamepattern string) {
	for _, line := range lines {
		output := processLine(line, template)
		outputFilename := processLine(line, outputfilenamepattern)

		if err := ioutil.WriteFile(outputFilename, []byte(output), os.ModePerm); err != nil {
			fmt.Println(err)
		}
	}
}
