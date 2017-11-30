###### Notice

*This file documents the changes in **config** versions that are listed below.*

*Items should be added to this file as:*

	### YYYY-MM-DD  Release

	+ Additional changes.

	+ More changes.

* * *

### 2011-??-??  v0.9.6

+ Changed to line comments.


### 2010-09-15  v0.9.5

+ Sections, options and values are all case-sensitive.

+ Changed API:

  Type *File* -> *Config*  
  *NewFile()* -> *NewDefault*  
  *ReadFile()* -> *ReadDefault*

+ Added functions, *New()*, *Read()*, which allow to choose the character of
comment and separator, and the spaces around separator.

+ Better error handling.

+ Both sections and options are showed by its input order.


### 2010-08-22  v0.9

+ The files has been splitted, formatted via *gomft*.

+ Methods use *self* to refer to its own type.

+ *Get* has been removed from the functions names.

+ Fixed some errors. All tests are passed.

+ At write the header in configuration file, it is added the comment character
after of each new line.

+ Better documentation.

