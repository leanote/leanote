config
======

This package implements a basic configuration file parser language which
provides a structure similar to what you would find on Microsoft Windows INI
files.

The configuration file consists of sections, led by a "*[section]*" header and
followed by "*name: value*" entries; "*name=value*" is also accepted. Note that
leading whitespace is removed from values. The optional values can contain
format strings which refer to other values in the same section, or values in a
special *DEFAULT* section. Additional defaults can be provided on initialization
and retrieval. Comments are indicated by ";" or "#"; a comment may begin
anywhere on a line, including on the same line after parameters or section
declarations.

For example:

	[My Section]
	foodir: %(dir)s/whatever
	dir=foo

would resolve the "*%(dir)s*" to the value of "*dir*" (*foo* in this case). All
reference expansions are done on demand.

The functionality and workflow is loosely based on the *configparser* package of
the Python Standard Library.

## Installation

	go get github.com/robfig/config

## Operating instructions

Given a sample configuration file:

	[DEFAULT]
	host: www.example.com
	protocol: http://
	base-url: %(protocol)s%(host)s

	[service-1]
	url: %(base-url)s/some/path
	delegation: on
	maxclients: 200 # do not set this higher
	comments: This is a multi-line
		entry	# And this is a comment

To read this configuration file, do:

	c, _ := config.ReadDefault("config.cfg")

	c.String("service-1", "url")
	// result is string "http://www.example.com/some/path"

	c.Int("service-1", "maxclients")
	// result is int 200

	c.Bool("service-1", "delegation")
	// result is bool true

	c.String("service-1", "comments")
	// result is string "This is a multi-line\nentry"

Note the support for unfolding variables (such as *%(base-url)s*), which are read
from the special (reserved) section name *[DEFAULT]*.

A new configuration file can also be created with:

	c := config.NewDefault()
	c.AddSection("Section")
	c.AddOption("Section", "option", "value")
	c.WriteFile("config.cfg", 0644, "A header for this file")

This results in the file:

	# A header for this file

	[Section]
	option: value

Note that sections, options and values are all case-sensitive.

## License

The source files are distributed under the [Mozilla Public License, version 2.0](http://mozilla.org/MPL/2.0/),
unless otherwise noted.  
Please read the [FAQ](http://www.mozilla.org/MPL/2.0/FAQ.html)
if you have further questions regarding the license.

