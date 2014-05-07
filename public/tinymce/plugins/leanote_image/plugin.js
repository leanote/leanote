/**
 * 图片上传插件
 * 结合了原tinymce image插件
 * 添加了input上传和拖拽上传
 * http://leanote.com
 */

tinymce.PluginManager.add('leanote_image', function(editor, url) {
	// 弹框
	function showDialog() {
		/*
		win = editor.windowManager.open({
			title: 'Insert/edit image',
			file : url + '/dialog.html',
			width : 550,
			height: 345
		});
		*/
		showDialog2("#imageDialog", {postShow:function() {
			$("#imageDialog iframe").attr("src", url + '/dialog.html');
		}});
	}

	// 添加按钮
	editor.addButton('leanote_image', {
		icon: 'image',
		tooltip: 'ctrl+shift+i 插入/修改图片',
		onclick: showDialog,
		stateSelector: 'img:not([data-mce-object])'
	});
	
	editor.addShortcut('ctrl+shift+i', '', showDialog);
});