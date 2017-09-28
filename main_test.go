package main

import (
	"math"
	"strings"
	"testing"
)

func TestAverage(t *testing.T) {
	var testCases = []struct {
		coverage    string
		expectedAvg float64
	}{
		{
			coverage: `
			ok      bitbucket.org/alkira/email/config       0.010s  coverage: 90.0% of statements
			ok      bitbucket.org/alkira/email/email        0.102s  coverage: 0.0% of statements
			ok      bitbucket.org/alkira/email/message      0.331s  coverage: 67.6% of statements
			ok      bitbucket.org/alkira/email/model        0.155s  coverage: 85.9% of statements
			`,
			expectedAvg: 60.875,
		}, {
			coverage: `
			ok      bitbucket.org/alkira/email/config       0.010s  coverage: 90.0% of statements
			ok      bitbucket.org/alkira/email/message      0.331s  coverage: 67.6% of statements
			ok      bitbucket.org/alkira/email/model        0.155s  coverage: 85.9% of statements
			`,
			expectedAvg: 81.166667,
		},
	}
	for _, tC := range testCases {
		t.Run("", func(t *testing.T) {
			avg, err := Average(strings.NewReader(tC.coverage))
			if err != nil {
				t.Fatalf("average:%v", err)
			}

			avg = math.Ceil(avg)
			expectedAvg := math.Ceil(tC.expectedAvg)

			if avg != expectedAvg {
				t.Fatalf("got average %f want average %f", avg, expectedAvg)
			}
		})
	}
}
