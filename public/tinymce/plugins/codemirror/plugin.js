/**
 * plugin.js
 *
 * Copyright 2013 Web Power, www.webpower.nl
 * @author Arjan Haverkamp
 */

/*jshint unused:false */
/*global tinymce:true */

tinymce.PluginManager.requireLangPack('codemirror');

tinymce.PluginManager.add('codemirror', function(editor, url) {

	function showSourceEditor() {
		// Insert caret marker
		editor.focus();
		editor.selection.collapse(true);
		editor.selection.setContent('<span class="CmCaReT" style="display:none">&#0;</span>');

		// Open editor window
		var height = $(document).height();
		var win = editor.windowManager.open({
			title: 'HTML source code',
			url: url + '/source.html',
			width: 800,
			height: height-150,
			resizable : true,
			maximizable : true,
			buttons: [
				{ text: 'Ok', subtype: 'primary', onclick: function(){
					var doc = document.querySelectorAll('.mce-container-body>iframe')[0];
					doc.contentWindow.submit();
					win.close();
				}},
				{ text: 'Cancel', onclick: 'close' }
			]
		});
	};

	// Add a button to the button bar
	editor.addButton('code', {
		title: 'Source code',
		image: url + '/img/file-html.png',
		onclick: showSourceEditor
	});

	// Add a menu item to the tools menu
	editor.addMenuItem('code', {
		icon: 'code',
		text: 'Source code',
		context: 'tools',
		onclick: showSourceEditor
	});
});
