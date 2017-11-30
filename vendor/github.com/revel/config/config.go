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
	"regexp"
	"strings"
)

// config constants
const (
	// Default section name.
	DefaultSection = "DEFAULT"
	// Maximum allowed depth when recursively substituing variable names.
	DepthValues = 200

	DefaultComment       = "# "
	AlternativeComment   = "; "
	DefaultSeparator     = ":"
	AlternativeSeparator = "="
)

var (
	// Strings accepted as boolean.
	boolString = map[string]bool{
		"t":     true,
		"true":  true,
		"y":     true,
		"yes":   true,
		"on":    true,
		"1":     true,
		"f":     false,
		"false": false,
		"n":     false,
		"no":    false,
		"off":   false,
		"0":     false,
	}

	varRegExp    = regexp.MustCompile(`%\(([a-zA-Z0-9_.\-]+)\)s`) // %(variable)s
	envVarRegExp = regexp.MustCompile(`\${([a-zA-Z0-9_.\-]+)}`)   // ${envvar}
)

// Config is the representation of configuration settings.
type Config struct {
	comment   string
	separator string

	// Sections order
	lastIDSection int            // Last section identifier
	idSection     map[string]int // Section : position

	// The last option identifier used for each section.
	lastIDOption map[string]int // Section : last identifier

	// Section -> option : value
	data map[string]map[string]*tValue
}

// tValue holds the input position for a value.
type tValue struct {
	position int    // Option order
	v        string // value
}

// New creates an empty configuration representation.
// This representation can be filled with AddSection and AddOption and then
// saved to a file using WriteFile.
//
// == Arguments
//
// comment: has to be `DefaultComment` or `AlternativeComment`
// separator: has to be `DefaultSeparator` or `AlternativeSeparator`
// preSpace: indicate if is inserted a space before of the separator
// postSpace: indicate if is added a space after of the separator
func New(comment, separator string, preSpace, postSpace bool) *Config {
	if comment != DefaultComment && comment != AlternativeComment {
		panic("comment character not valid")
	}

	if separator != DefaultSeparator && separator != AlternativeSeparator {
		panic("separator character not valid")
	}

	// == Get spaces around separator
	if preSpace {
		separator = " " + separator
	}

	if postSpace {
		separator += " "
	}
	//==

	c := new(Config)

	c.comment = comment
	c.separator = separator
	c.idSection = make(map[string]int)
	c.lastIDOption = make(map[string]int)
	c.data = make(map[string]map[string]*tValue)

	c.AddSection(DefaultSection) // Default section always exists.

	return c
}

// NewDefault creates a configuration representation with values by default.
func NewDefault() *Config {
	return New(DefaultComment, DefaultSeparator, false, true)
}

// Merge merges the given configuration "source" with this one ("target").
//
// Merging means that any option (under any section) from source that is not in
// target will be copied into target. When the target already has an option with
// the same name and section then it is overwritten (i.o.w. the source wins).
func (target *Config) Merge(source *Config) {
	if source == nil || source.data == nil || len(source.data) == 0 {
		return
	}

	for section, option := range source.data {
		for optionName, optionValue := range option {
			target.AddOption(section, optionName, optionValue.v)
		}
	}
}

// == Utility

func stripComments(l string) string {
	// Comments are preceded by space or TAB
	for _, c := range []string{" ;", "\t;", " #", "\t#"} {
		if i := strings.Index(l, c); i != -1 {
			l = l[0:i]
		}
	}
	return l
}
