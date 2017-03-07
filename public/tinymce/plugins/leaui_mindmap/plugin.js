/**
 * leaui mind map plugin
 * copyright leaui
 * leaui.com
 */
var LEAUI_MIND = {};
tinymce.PluginManager.add('leaui_mindmap', function(editor, url) {
	
	function showDialog() {
		var dom = editor.dom;

		var content = editor.selection.getContent();
		// get images and attrs
		var p = /<img.*?\/>/g;
		var images = content.match(p);
		var newNode = document.createElement("p");
		LEAUI_MIND = {};
		for(var i in images) {
			newNode.innerHTML = images[i];
			var imgElm = newNode.firstChild;
			if(imgElm && imgElm.nodeName == "IMG") {
				LEAUI_MIND.json = dom.getAttrib(imgElm, 'data-mind-json');
				break;
			}
		}

		function GetTheHtml(){
			var lang = editor.settings.language;
			var u = url + '/mindmap/index.html?i=1';
			var html = '<iframe id="leauiMindMapIfr" src="'+ u + '?' + new Date().getTime() + '&lang=' + lang + '" frameborder="0"></iframe>';
			return html;
		}

		var w = window.innerWidth - 10;
		var h = window.innerHeight - 150;

		win = editor.windowManager.open({
			title: "Mind Map",
			width : w,
			height : h,
			html: GetTheHtml(),
			buttons: [
				{
					text: 'Cancel',
					onclick: function() {
						this.parent().parent().close();
					}
				},
				{
				text: 'Insert',
				subtype: 'primary',
				onclick: function(e) {
					var me = this;
					var _iframe = document.getElementById('leauiMindMapIfr').contentWindow;
					var km = _iframe.km;
					// window.km= km;
					// return
					km.exportData('png').then(function(data) {
						var json = JSON.stringify(km.exportJson());
						json = json.replace(/'/g, "Ж");
						// console.log(json);
						var img = '<img src="' + data + '" data-mce-src="-" data-mind-json=\'' + json + '\'>';
						editor.insertContent(img);

						me.parent().parent().close();
					});
					return;
				}
				}]
		});
	}
	
	editor.addButton('leaui_mindmap', {
		// image: url + '/icon.png',
		icon: 'mind',
		tooltip: 'Insert/edit mind map',
		onclick: showDialog,
		stateSelector: 'img[data-mind-json]'
	});
});
