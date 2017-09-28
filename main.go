package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func exit(err error) {
	fmt.Printf("error:%v\n", err)
	os.Exit(1)
}

func Average(r io.Reader) (float64, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, fmt.Errorf("could not read coverage data:%v", err)
	}

	lines := strings.Split(string(b), "\n")

	re, err := regexp.Compile("coverage:\\s\\d+\\.\\d+")
	if err != nil {
		return 0, fmt.Errorf("could not compile regexp:%v", err)
	}

	var totalCov float64
	var pkgCount float64

	for _, line := range lines {
		match := re.FindString(line)
		if match == "" {
			continue
		}

		cov, err := strconv.ParseFloat(strings.Split(match, " ")[1], 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse coverage:%v", err)
		}

		totalCov += cov
		pkgCount = pkgCount + 1
	}

	return totalCov / pkgCount, nil
}

func main() {
	var covFile string
	var threshold float64
	flag.StringVar(&covFile, "file", "", "file to go test -cover output")
	flag.Float64Var(&threshold, "threshold", 0.0, "the minimum coverage")
	flag.Parse()

	fi, err := os.Stdin.Stat()
	if err != nil {
		exit(fmt.Errorf("failed to stat stdin:%v", err))
	}

	var avg float64
	if fi.Size() > 0 {
		avg, err = Average(os.Stdin)
	} else {
		if covFile == "" {
			fmt.Printf("error:file is required")
			flag.Usage()
			os.Exit(1)
		}

		f, err := os.Open(covFile)
		if err != nil {
			exit(err)
		}
		defer f.Close()

		avg, err = Average(f)
	}

	fmt.Printf("average coverage:%f\n", avg)

	if threshold == 0 {
		os.Exit(0)
	}

	if avg < threshold {
		fmt.Println("threshold not met")
		os.Exit(1)
	}

	fmt.Println("threshold met!")

}
