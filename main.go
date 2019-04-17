package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	csvfilename := flag.String("c", "", "input csv filename")
	templatefilename := flag.String("t", "", "template filename")
	mode := flag.String("m", "stdio", "output mode. stdio or file")
	outputfilepattern := flag.String("o", "{0}.txt", "file output mode pattern. like {0}.txt")
	flag.Parse()

	if *csvfilename == "" {
		fmt.Println("csv filename is required")
		flag.PrintDefaults()
		return
	}

	if *mode == "file" && *outputfilepattern == "" {
		fmt.Println("output file pattern is required for file mode")
		flag.PrintDefaults()
		return
	}

	baseDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	absCsvFilename := getAbsFile(baseDir, *csvfilename)
	absTemplateFilename := getAbsFile(baseDir, *templatefilename)
	absOutputFilePattern := getAbsFile(baseDir, *outputfilepattern)

	lines, err := readCsv(absCsvFilename)
	if err != nil {
		log.Fatal(err)
	}

	template := readTemplate(absTemplateFilename)

	if *mode == "stdio" {
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

func readCsv(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	lines := []string{}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	return lines, nil
}

func readTemplate(filename string) (template string) {
	b, _ := ioutil.ReadFile(filename)
	template = string(b)
	return
}

func processLines(lines []string, template string) (fullReplace string) {
	for i, line := range lines {
		fullReplace = fullReplace + processLine(line, template, i)
	}

	return
}

func processLine(line string, template string, num int) string {
	s := strings.Split(line, ",")
	tmp := template
	for i := range s {
		index := fmt.Sprintf("{%v}", i)
		tmp = strings.Replace(tmp, index, s[i], -1)
		tmp = strings.Replace(tmp, "{index}", strconv.Itoa(num), -1)
	}

	return tmp
}

func processLinesFileOutput(lines []string, template string, outputfilenamepattern string) {
	for i, line := range lines {
		output := processLine(line, template, i)
		outputFilename := processLine(line, outputfilenamepattern, i)

		if err := ioutil.WriteFile(outputFilename, []byte(output), os.ModePerm); err != nil {
			fmt.Println(err)
		}
	}
}
