package tests

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/alejandrogonzalr/load-test-cli/cmd"
)

type testCase map[string]struct {
	args     []string
	expected []string
}

func TestRootCmd(t *testing.T) {
	testUrl := "http://google.com"
	tests := testCase{
		"url param":        {args: []string{testUrl}, expected: []string{testUrl, "1"}},
		"concurrency flag": {args: []string{"--concurrency", "5", testUrl}, expected: []string{testUrl, "5"}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := execTestCmd(test.args)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Fatalf("Expected \"%s\" got \"%s\"", test.expected, result)
			}
		})
	}

}

func execTestCmd(args []string) ([]string, error) {
	cmd := cmd.RootCmd
	buf := bytes.NewBufferString("")
	cmd.SetOut(buf)
	cmd.SetArgs(args)
	cmd.Execute()
	out, err := ioutil.ReadAll(buf)
	outString := string(out)
	outString = strings.TrimSuffix(outString, "\n")
	c, _ := cmd.Flags().GetInt("concurrency")
	result := []string{
		outString,
		strconv.Itoa(c),
	}

	return result, err
}
