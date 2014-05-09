/**
 * plugin.js
 *
 * Copyright, Moxiecode Systems AB
 * Released under LGPL License.
 *
 * License: http://www.tinymce.com/license
 * Contributing: http://www.tinymce.com/contributing
 */

/*global tinymce:true */

// 与jbimages结合
// 使用了jbimages/dialog-v4.htm
// jbimages/js/dialog-v4.js
tinymce.PluginManager.add('image', function(editor, url) {
	// 弹框
	function showDialog() {
		// 与jbimages结合
		// Simple default dialog
		win = editor.windowManager.open({
			title: 'Insert/edit image',
			file : url + '/dialog.htm',
			width : 550,
			height: 345
		});
	}

	// 添加按钮
	editor.addButton('image', {
		icon: 'image',
		tooltip: 'Insert/edit image',
		onclick: showDialog,
		stateSelector: 'img:not([data-mce-object])'
	});
});