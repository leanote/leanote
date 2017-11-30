// Copyright 2009  The "config" Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"bufio"
	"os"
	"reflect"
	"strings"
	"testing"
)

const (
	tmpFilename    = "testdata/__test.go"
	sourceFilename = "testdata/source.cfg"
	targetFilename = "testdata/target.cfg"
)

func testGet(t *testing.T, c *Config, section string, option string,
	expected interface{}) {
	ok := false
	switch expected.(type) {
	case string:
		v, _ := c.String(section, option)
		if v == expected.(string) {
			ok = true
		}
	case int:
		v, _ := c.Int(section, option)
		if v == expected.(int) {
			ok = true
		}
	case bool:
		v, _ := c.Bool(section, option)
		if v == expected.(bool) {
			ok = true
		}
	default:
		t.Fatalf("Bad test case")
	}
	if !ok {
		v, _ := c.String(section, option)
		t.Errorf("Get failure: expected different value for %s %s (expected: [%#v] got: [%#v])", section, option, expected, v)
	}
}

// TestInMemory creates configuration representation and run multiple tests in-memory.
func TestInMemory(t *testing.T) {
	c := NewDefault()

	// == Test empty structure

	// should be empty
	if len(c.Sections()) != 1 {
		t.Errorf("Sections failure: invalid length")
	}

	// test presence of missing section
	if c.HasSection("no-section") {
		t.Errorf("HasSection failure: invalid section")
	}

	// get options for missing section
	_, err := c.Options("no-section")
	if err == nil {
		t.Errorf("Options failure: invalid section")
	}

	// test presence of option for missing section
	if c.HasOption("no-section", "no-option") {
		t.Errorf("HasSection failure: invalid/section/option")
	}

	// get value from missing section/option
	_, err = c.String("no-section", "no-option")
	if err == nil {
		t.Errorf("String failure: got value for missing section/option")
	}

	// get value from missing section/option
	_, err = c.Int("no-section", "no-option")
	if err == nil {
		t.Errorf("Int failure: got value for missing section/option")
	}

	// remove missing section
	if c.RemoveSection("no-section") {
		t.Errorf("RemoveSection failure: removed missing section")
	}

	// remove missing section/option
	if c.RemoveOption("no-section", "no-option") {
		t.Errorf("RemoveOption failure: removed missing section/option")
	}

	// == Fill up structure

	// add section
	if !c.AddSection("section1") {
		t.Errorf("AddSection failure: false on first insert")
	}

	// re-add same section
	if c.AddSection("section1") {
		t.Errorf("AddSection failure: true on second insert")
	}

	// default section always exists
	if c.AddSection(DefaultSection) {
		t.Errorf("AddSection failure: true on default section insert")
	}

	// add option/value
	if !c.AddOption("section1", "option1", "value1") {
		t.Errorf("AddOption failure: false on first insert")
	}
	testGet(t, c, "section1", "option1", "value1") // read it back

	// overwrite value
	if c.AddOption("section1", "option1", "value2") {
		t.Errorf("AddOption failure: true on second insert")
	}
	testGet(t, c, "section1", "option1", "value2") // read it back again

	// remove option/value
	if !c.RemoveOption("section1", "option1") {
		t.Errorf("RemoveOption failure: false on first remove")
	}

	// remove again
	if c.RemoveOption("section1", "option1") {
		t.Errorf("RemoveOption failure: true on second remove")
	}

	// read it back again
	_, err = c.String("section1", "option1")
	if err == nil {
		t.Errorf("String failure: got value for removed section/option")
	}

	// remove existing section
	if !c.RemoveSection("section1") {
		t.Errorf("RemoveSection failure: false on first remove")
	}

	// remove again
	if c.RemoveSection("section1") {
		t.Errorf("RemoveSection failure: true on second remove")
	}

	// == Test types

	// add section
	if !c.AddSection("section2") {
		t.Errorf("AddSection failure: false on first insert")
	}

	// add number
	if !c.AddOption("section2", "test-number", "666") {
		t.Errorf("AddOption failure: false on first insert")
	}
	testGet(t, c, "section2", "test-number", 666) // read it back

	// add 'yes' (bool)
	if !c.AddOption("section2", "test-yes", "yes") {
		t.Errorf("AddOption failure: false on first insert")
	}
	testGet(t, c, "section2", "test-yes", true) // read it back

	// add 'false' (bool)
	if !c.AddOption("section2", "test-false", "false") {
		t.Errorf("AddOption failure: false on first insert")
	}
	testGet(t, c, "section2", "test-false", false) // read it back

	// == Test cycle

	c.AddOption(DefaultSection, "opt1", "%(opt2)s")
	c.AddOption(DefaultSection, "opt2", "%(opt1)s")

	_, err = c.String(DefaultSection, "opt1")
	if err == nil {
		t.Errorf("String failure: no error for cycle")
	} else if !strings.Contains(err.Error(), "cycle") {
		t.Errorf("String failure: incorrect error for cycle")
	}
}

// TestReadFile creates a 'tough' configuration file and test (read) parsing.
func TestReadFile(t *testing.T) {
	file, err := os.Create(tmpFilename)
	if err != nil {
		t.Fatal("Test cannot run because cannot write temporary file: " + tmpFilename)
	}

	err = os.Setenv("GO_CONFIGFILE_TEST_ENV_VAR", "configvalue12345")
	if err != nil {
		t.Fatalf("Test cannot run because cannot set environment variable GO_CONFIGFILE_TEST_ENV_VAR: %#v", err)
	}

	buf := bufio.NewWriter(file)
	buf.WriteString("optionInDefaultSection=true\n")
	buf.WriteString("[section-1]\n")
	buf.WriteString("option1=value1 ; This is a comment\n")
	buf.WriteString("option2 : 2#Not a comment\t#Now this is a comment after a TAB\n")
	buf.WriteString("  # Let me put another comment\n")
	buf.WriteString("option3= line1\n line2: \n\tline3=v # Comment multiline with := in value\n")
	buf.WriteString("; Another comment\n")
	buf.WriteString("[" + DefaultSection + "]\n")
	buf.WriteString("variable1=small\n")
	buf.WriteString("variable2=a_part_of_a_%(variable1)s_test\n")
	buf.WriteString("[secTION-2]\n")
	buf.WriteString("IS-flag-TRUE=Yes\n")
	buf.WriteString("[section-1] # comment on section header\n") // continue again [section-1]
	buf.WriteString("option4=this_is_%(variable2)s.\n")
	buf.WriteString("envoption1=this_uses_${GO_CONFIGFILE_TEST_ENV_VAR}_env\n")
	buf.WriteString("optionInDefaultSection=false")
	buf.Flush()
	file.Close()

	c, err := ReadDefault(tmpFilename)
	if err != nil {
		t.Fatalf("ReadDefault failure: %s", err)
	}

	// check number of sections
	if len(c.Sections()) != 3 {
		t.Errorf("Sections failure: wrong number of sections")
	}

	// check number of options 6 of [section-1] plus 2 of [default]
	opts, _ := c.Options("section-1")
	if len(opts) != 8 {
		t.Errorf("Options failure: wrong number of options: %d", len(opts))
	}

	testGet(t, c, "section-1", "option1", "value1")
	testGet(t, c, "section-1", "option2", "2#Not a comment")
	testGet(t, c, "section-1", "option3", "line1\nline2:\nline3=v")
	testGet(t, c, "section-1", "option4", "this_is_a_part_of_a_small_test.")
	testGet(t, c, "section-1", "envoption1", "this_uses_configvalue12345_env")
	testGet(t, c, "section-1", "optionInDefaultSection", false)
	testGet(t, c, "section-2", "optionInDefaultSection", true)
	testGet(t, c, "secTION-2", "IS-flag-TRUE", true) // case-sensitive
}

// TestWriteReadFile tests writing and reading back a configuration file.
func TestWriteReadFile(t *testing.T) {
	cw := NewDefault()

	// write file; will test only read later on
	cw.AddSection("First-Section")
	cw.AddOption("First-Section", "option1", "value option1")
	cw.AddOption("First-Section", "option2", "2")

	cw.AddOption("", "host", "www.example.com")
	cw.AddOption(DefaultSection, "protocol", "https://")
	cw.AddOption(DefaultSection, "base-url", "%(protocol)s%(host)s")

	cw.AddOption("Another-Section", "useHTTPS", "y")
	cw.AddOption("Another-Section", "url", "%(base-url)s/some/path")

	cw.WriteFile(tmpFilename, 0644, "Test file for test-case")

	// read back file and test
	cr, err := ReadDefault(tmpFilename)
	if err != nil {
		t.Fatalf("ReadDefault failure: %s", err)
	}

	testGet(t, cr, "First-Section", "option1", "value option1")
	testGet(t, cr, "First-Section", "option2", 2)
	testGet(t, cr, "Another-Section", "useHTTPS", true)
	testGet(t, cr, "Another-Section", "url", "https://www.example.com/some/path")

	defer os.Remove(tmpFilename)
}

// TestSectionOptions tests read options in a section without default options.
func TestSectionOptions(t *testing.T) {
	cw := NewDefault()

	// write file; will test only read later on
	cw.AddSection("First-Section")
	cw.AddOption("First-Section", "option1", "value option1")
	cw.AddOption("First-Section", "option2", "2")

	cw.AddOption("", "host", "www.example.com")
	cw.AddOption(DefaultSection, "protocol", "https://")
	cw.AddOption(DefaultSection, "base-url", "%(protocol)s%(host)s")

	cw.AddOption("Another-Section", "useHTTPS", "y")
	cw.AddOption("Another-Section", "url", "%(base-url)s/some/path")

	cw.WriteFile(tmpFilename, 0644, "Test file for test-case")

	// read back file and test
	cr, err := ReadDefault(tmpFilename)
	if err != nil {
		t.Fatalf("ReadDefault failure: %s", err)
	}

	options, err := cr.SectionOptions("First-Section")

	if err != nil {
		t.Fatalf("SectionOptions failure: %s", err)
	}

	if len(options) != 2 {
		t.Fatalf("SectionOptions reads wrong data: %v", options)
	}

	expected := map[string]bool{
		"option1": true,
		"option2": true,
	}
	actual := map[string]bool{}

	for _, v := range options {
		actual[v] = true
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("SectionOptions reads wrong data: %v", options)
	}

	options, err = cr.SectionOptions(DefaultSection)

	if err != nil {
		t.Fatalf("SectionOptions failure: %s", err)
	}

	expected = map[string]bool{
		"host":     true,
		"protocol": true,
		"base-url": true,
	}
	actual = map[string]bool{}

	for _, v := range options {
		actual[v] = true
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("SectionOptions reads wrong data: %v", options)
	}

	defer os.Remove(tmpFilename)
}

// TestMerge tests merging 2 configurations.
func TestMerge(t *testing.T) {
	target, error := ReadDefault(targetFilename)
	if error != nil {
		t.Fatalf("Unable to read target config file '%s'", targetFilename)
	}

	source, error := ReadDefault(sourceFilename)
	if error != nil {
		t.Fatalf("Unable to read source config file '%s'", sourceFilename)
	}

	target.Merge(source)

	// Assert whether a regular option was merged from source -> target
	if result, _ := target.String(DefaultSection, "one"); result != "source1" {
		t.Errorf("Expected 'one' to be '1' but instead it was '%s'", result)
	}
	// Assert that a non-existent option in source was not overwritten
	if result, _ := target.String(DefaultSection, "five"); result != "5" {
		t.Errorf("Expected 'five' to be '5' but instead it was '%s'", result)
	}
	// Assert that a folded option was correctly unfolded
	if result, _ := target.String(DefaultSection, "two_+_three"); result != "source2 + source3" {
		t.Errorf("Expected 'two_+_three' to be 'source2 + source3' but instead it was '%s'", result)
	}
	if result, _ := target.String(DefaultSection, "four"); result != "4" {
		t.Errorf("Expected 'four' to be '4' but instead it was '%s'", result)
	}

	// Assert that a section option has been merged
	if result, _ := target.String("X", "x.one"); result != "sourcex1" {
		t.Errorf("Expected '[X] x.one' to be 'sourcex1' but instead it was '%s'", result)
	}
	if result, _ := target.String("X", "x.four"); result != "x4" {
		t.Errorf("Expected '[X] x.four' to be 'x4' but instead it was '%s'", result)
	}
}

func TestLoadContextOneConf(t *testing.T) {
	ctx, err := LoadContext("app.conf", []string{"testdata/conf-path1"})
	if err != nil {
		t.Errorf("Error: %v", err)
		t.FailNow()
	}

	ctx.SetSection("X")
	result, found := ctx.String("x.three")
	if !strings.EqualFold("conf1-sourcex3", result) {
		t.Errorf("Expected '[X] x.three' to be 'conf1-sourcex3' but instead it was '%s'", result)
	}

	_, found = ctx.String("x.notexists")
	if found {
		t.Error("Config 'x.notexists' shouldn't found")
	}

	ctx.SetSection("Y")
	result, found = ctx.String("y.one")
	if !strings.EqualFold("conf1-sourcey1", result) {
		t.Errorf("Expected '[Y] y.one' to be 'conf1-sourcey1' but instead it was '%s'", result)
	}

	_, found = ctx.String("y.notexists")
	if found {
		t.Error("Config 'y.notexists' shouldn't found")
	}
}

func TestLoadContextMultipleConfWithPriority(t *testing.T) {
	ctx, err := LoadContext("app.conf", []string{"testdata/conf-path1", "testdata/conf-path2"})
	if err != nil {
		t.Errorf("Error: %v", err)
		t.FailNow()
	}

	ctx.SetSection("X")
	result, found := ctx.String("x.two")
	if !strings.EqualFold("override-conf2-sourcex2", result) {
		t.Errorf("Expected '[X] x.two' to be 'override-conf2-sourcex2' but instead it was '%s'", result)
	}

	_, found = ctx.String("x.notexists")
	if found {
		t.Error("Config 'x.notexists' shouldn't be found")
	}

	ctx.SetSection("Y")
	result, found = ctx.String("y.three")
	if !strings.EqualFold("override-conf2-sourcey3", result) {
		t.Errorf("Expected '[Y] y.three' to be 'override-conf2-sourcey3' but instead it was '%s'", result)
	}

	_, found = ctx.String("y.notexists")
	if found {
		t.Error("Config 'y.notexists' shouldn't be found")
	}
}

func TestLoadContextConfNotFound(t *testing.T) {
	_, err := LoadContext("notfound.conf", []string{"testdata/conf-path1"})
	if err != nil && !strings.EqualFold("open testdata/conf-path1/notfound.conf: no such file or directory", err.Error()) {
		t.Errorf("This is not expected error: %v", err)
	}
}

func TestLoadContextInvalidConf(t *testing.T) {
	_, err := LoadContext("app-invalid.conf", []string{"testdata"})
	if err != nil && !strings.EqualFold("testdata/app-invalid.conf: could not parse line #7: %(two)s + %(four)s", err.Error()) {
		t.Errorf("This is not expected error: %v", err)
	}
}
