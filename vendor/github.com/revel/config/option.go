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

import "errors"

// AddOption adds a new option and value to the configuration.
//
// If the section is nil then uses the section by default; if it does not exist,
// it is created in advance.
//
// It returns true if the option and value were inserted, and false if the value
// was overwritten.
func (c *Config) AddOption(section string, option string, value string) bool {
	c.AddSection(section) // Make sure section exists

	if section == "" {
		section = DefaultSection
	}

	_, ok := c.data[section][option]

	c.data[section][option] = &tValue{c.lastIDOption[section], value}
	c.lastIDOption[section]++

	return !ok
}

// RemoveOption removes a option and value from the configuration.
// It returns true if the option and value were removed, and false otherwise,
// including if the section did not exist.
func (c *Config) RemoveOption(section string, option string) bool {
	if _, ok := c.data[section]; !ok {
		return false
	}

	_, ok := c.data[section][option]
	delete(c.data[section], option)

	return ok
}

// HasOption checks if the configuration has the given option in the section.
// It returns false if either the option or section do not exist.
func (c *Config) HasOption(section string, option string) bool {
	if _, ok := c.data[section]; !ok {
		return false
	}

	_, okd := c.data[DefaultSection][option]
	_, oknd := c.data[section][option]

	return okd || oknd
}

// Options returns the list of options available in the given section.
// It returns an error if the section does not exist and an empty list if the
// section is empty. Options within the default section are also included.
func (c *Config) Options(section string) (options []string, err error) {
	if _, ok := c.data[section]; !ok {
		return nil, errors.New(SectionError(section).Error())
	}

	// Keep a map of option names we've seen to deduplicate.
	optionMap := make(map[string]struct{},
		len(c.data[DefaultSection])+len(c.data[section]))
	for s := range c.data[DefaultSection] {
		optionMap[s] = struct{}{}
	}
	for s := range c.data[section] {
		optionMap[s] = struct{}{}
	}

	// Get the keys.
	i := 0
	options = make([]string, len(optionMap))
	for k := range optionMap {
		options[i] = k
		i++
	}

	return options, nil
}

// SectionOptions returns only the list of options available in the given section.
// Unlike Options, SectionOptions doesn't return options in default section.
// It returns an error if the section doesn't exist.
func (c *Config) SectionOptions(section string) (options []string, err error) {
	if _, ok := c.data[section]; !ok {
		return nil, errors.New(SectionError(section).Error())
	}

	options = make([]string, len(c.data[section]))
	i := 0
	for s := range c.data[section] {
		options[i] = s
		i++
	}

	return options, nil
}
