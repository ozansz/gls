package local

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	replacements = []struct {
		from string
		to   string
	}{
		{from: `.`, to: `\.`},
		{from: `*`, to: `.*`},
	}
)

type ignoreRule struct {
	isDirRule bool
	re        *regexp.Regexp
}

type IgnoreChecker struct {
	ruleFiles map[string]struct{}
	rules     []ignoreRule
}

type IgnoreCheckerOption func(*IgnoreChecker)

func NewIgnoreChecker(opts ...IgnoreCheckerOption) (*IgnoreChecker, error) {
	c := &IgnoreChecker{
		ruleFiles: make(map[string]struct{}),
		rules:     nil,
	}
	c.ruleFiles[defaultIgnoreFile] = struct{}{}
	for _, o := range opts {
		o(c)
	}
	if err := c.generateRules(); err != nil {
		return nil, err
	}
	return c, nil
}

func WithRuleFile(path string) IgnoreCheckerOption {
	return func(ic *IgnoreChecker) {
		if _, ok := ic.ruleFiles[path]; !ok {
			ic.ruleFiles[path] = struct{}{}
		}
	}
}

func (ic *IgnoreChecker) ShouldIgnore(path string, isDir bool) bool {
	for _, r := range ic.rules {
		if ((r.isDirRule && isDir) || (!r.isDirRule && !isDir)) && r.re.MatchString(path) {
			return true
		}
	}
	return false
}

func (ic *IgnoreChecker) Dump() string {
	s := ""
	for i, r := range ic.rules {
		s += fmt.Sprintf("rule #%d: '%s', isDirRule: %t\n", i, r.re.String(), r.isDirRule)
	}
	return s
}

func (ic *IgnoreChecker) generateRules() error {
	ic.rules = make([]ignoreRule, 0)
	for path := range ic.ruleFiles {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			if path == defaultIgnoreFile && os.IsNotExist(err) {
				// Default ignore file does not exist is not an error case
				continue
			}
			return err
		}
		r := bufio.NewReader(bytes.NewReader(b))
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				if err != io.EOF {
					return fmt.Errorf("error while reading line from %q: %v", path, err)
				}
				break // io.EOF means we've read all of this file.
			}
			reStr, isDirRule := convertIgnorePatternToRegex(string(line))
			re, err := regexp.Compile(reStr)
			if err != nil {
				return fmt.Errorf("error while compiling regex rule %q from ignore file %q: %v", reStr, path, err)
			}
			ic.rules = append(ic.rules, ignoreRule{
				re:        re,
				isDirRule: isDirRule,
			})
		}
	}
	return nil
}

func convertIgnorePatternToRegex(pattern string) (s string, isDirRule bool) {
	isDirRule = false
	s = pattern
	for _, r := range replacements {
		s = strings.ReplaceAll(s, r.from, r.to)
	}
	if s[len(s)-1] == '/' {
		isDirRule = true
		s = s[:len(s)-1]
	}
	s += "$"
	return
}
