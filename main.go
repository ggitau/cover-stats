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
	"time"
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

	covRe, err := regexp.Compile("coverage:\\s\\d+\\.\\d+")
	if err != nil {
		return 0, fmt.Errorf("compile coverage regex:%v", err)
	}

	failRe, err := regexp.Compile("^FAIL")
	if err != nil {
		return 0, fmt.Errorf("compile fail regex:%v", err)
	}

	var totalCov float64
	var pkgCount float64
	for _, line := range lines {
		failMatch := failRe.FindString(line)
		if failMatch != "" {
			return 0, fmt.Errorf("tests failed")
		}

		match := covRe.FindString(line)
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

	if pkgCount == 0 {
		return 0, nil
	}
	return totalCov / pkgCount, nil
}

func main() {
	var covFile string
	var threshold float64
	var waitPeriod time.Duration
	flag.StringVar(&covFile, "file", "", "file to go test -cover output")
	flag.Float64Var(&threshold, "threshold", 0.0, "the minimum coverage")
	flag.DurationVar(&waitPeriod, "wait-period", 5*time.Second, "how long to wait for tests to run if tests results are being piped in")
	flag.Parse()


	var avg float64
	if covFile != "" {
		f, err := os.Open(covFile)
		if err != nil {
			exit(err)
		}
		defer f.Close()

		avg, err = Average(f)
	} else {

		time.Sleep(waitPeriod)

		fi, err := os.Stdin.Stat()
		if err != nil {
			exit(fmt.Errorf("failed to stat stdin:%v", err))
		}

		if fi.Size() == 0 {
			fmt.Printf("error:please provide an input file or pipe in the coverage results.")
			flag.Usage()
			os.Exit(1)
		} else {
			avg, err = Average(os.Stdin)
			if err != nil {
				exit(err)
			}
		}
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
