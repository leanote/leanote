function bindMathJaxHooks(converter) {

    var msie = /msie/.test(navigator.userAgent.toLowerCase());

	var inline = "$"; // the inline math delimiter
	//
	// The pattern for math delimiters and special symbols
	// needed for searching for math in the page.
	//
	var SPLIT = /(\$\$?|\\(?:begin|end)\{[a-z]*\*?\}|\\[\\{}$]|[{}]|(?:\n\s*)+|@@\d+@@)/i;

	//
	// The math is in blocks i through j, so
	// collect it into one block and clear the others.
	// Replace &, <, and > by named entities.
	// For IE, put <br> at the ends of comments since IE removes \n.
	// Clear the current math positions and store the index of the
	// math, then push the math string onto the storage array.
	//
	function processMath(i, j) {
		var block = blocks.slice(i, j + 1).join("").replace(/&/g, "&amp;") // use
																			// HTML
																			// entity
																			// for
																			// &
		.replace(/</g, "&lt;") // use HTML entity for <
		.replace(/>/g, "&gt;") // use HTML entity for >
		;
		if (msie) {
			block = block.replace(/(%[^\n]*)\n/g, "$1<br/>\n")
		}
		while (j > i) {
			blocks[j] = "";
			j--
		}
		blocks[i] = "@@" + math.length + "@@";
		math.push(block);
		start = end = last = null;
	}

    //
	// Break up the text into its component parts and search
	// through them for math delimiters, braces, linebreaks, etc.
	// Math delimiters must match and braces must balance.
	// Don't allow math to pass through a double linebreak
	// (which will be a paragraph).
	//
	function removeMath(text) {
		start = end = last = null; // for tracking math delimiters
		math = []; // stores math strings for latter

		blocks = text.replace(/\r\n?/g, "\n").split(SPLIT);
		for ( var i = 1, m = blocks.length; i < m; i += 2) {
			var block = blocks[i];
			if (block.charAt(0) === "@") {
				//
				// Things that look like our math markers will get
				// stored and then retrieved along with the math.
				//
				blocks[i] = "@@" + math.length + "@@";
				math.push(block);
			} else if (start) {
				//
				// If we are in math, look for the end delimiter,
				// but don't go past double line breaks, and
				// and balance braces within the math.
				//
				if (block === end) {
					if (braces) {
						last = i
					} else {
						processMath(start, i)
					}
				} else if (block.match(/\n.*\n/)) {
					if (last) {
						i = last;
						processMath(start, i)
					}
					start = end = last = null;
					braces = 0;
				} else if (block === "{") {
					braces++
				} else if (block === "}" && braces) {
					braces--
				}
			} else {
				//
				// Look for math start delimiters and when
				// found, set up the end delimiter.
                //
				// jiawzhang: to avoid handle $LaTex expression$, since sometimes, it conflicts. Handle $$LaTeX expression$$ only.
				// if (block === inline || block === "$$") {
				if (block === "$$") {
					start = i;
					end = block;
					braces = 0;
				} else if (block.substr(1, 5) === "begin") {
					start = i;
					end = "\\end" + block.substr(6);
					braces = 0;
				}
			}
		}
		if (last) {
			processMath(start, last)
		}
		return blocks.join("");
	}

	//
	// Put back the math strings that were saved,
	// and clear the math array (no need to keep it around).
	//  
	function replaceMath(text) {
		text = text.replace(/@@(\d+)@@/g, function(match, n) {
			return math[n]
		});
		math = null;
		return text;
	}

    converter.hooks.chain("preConversion", removeMath);
    converter.hooks.chain("postConversion", replaceMath);
}
