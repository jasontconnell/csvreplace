package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func main() {
	csvfilename := flag.String("c", "", "input csv filename")
	templatefilename := flag.String("t", "", "template filename")
	mode := flag.String("m", "stdout", "output mode. stdout or file")
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

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	baseDir, _ := filepath.Abs(wd)

	absCsvFilename := getAbsFile(baseDir, *csvfilename)
	absTemplateFilename := getAbsFile(baseDir, *templatefilename)
	absOutputFilePattern := getAbsFile(baseDir, *outputfilepattern)

	lines, err := readCsv(absCsvFilename)
	if err != nil {
		log.Fatal(err)
	}

	template, err := readTemplate(absTemplateFilename)
	if err != nil {
		log.Fatal(err)
	}

	if *mode == "stdout" {
		updated := processLines(lines, template)
		fmt.Println(updated)
	} else {
		err = processLinesFileOutput(lines, template, absOutputFilePattern)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func getAbsFile(baseDir, path string) string {
	abspath := path
	if !filepath.IsAbs(abspath) {
		abspath = filepath.Join(baseDir, abspath)
	}
	return abspath
}

func readCsv(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rdr := csv.NewReader(f)

	return rdr.ReadAll()
}

func readTemplate(filename string) (string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	template := string(b)
	return template, nil
}

func processLines(lines [][]string, template string) (fullReplace string) {
	for i, line := range lines {
		fullReplace = fullReplace + processLine(line, template, i)
	}

	return
}

func processLine(s []string, template string, num int) string {
	tmp := template
	for i := range s {
		index := fmt.Sprintf("{%v}", i)
		tmp = strings.Replace(tmp, index, strings.Trim(s[i], " "), -1)
		tmp = strings.Replace(tmp, "{index}", strconv.Itoa(num), -1)
	}

	tmp = strings.Replace(tmp, "{newguid}", uuid.New().String(), -1)

	return tmp
}

func processLinesFileOutput(lines [][]string, template string, outputfilenamepattern string) error {
	for i, line := range lines {
		output := processLine(line, template, i)
		outputFilename := processLine(line, outputfilenamepattern, i)

		err := os.WriteFile(outputFilename, []byte(output), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
