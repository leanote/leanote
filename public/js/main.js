require.config({
	baseUrl: '/public',
    paths: {
    	// 'jquery': 'js/jquery-1.9.0.min',
    	// base editor
    	'tinymce': 'tinymce/tinymce',
    	'jquery.slimscroll': 'js/jQuery-slimScroll-1.3.0/jquery.slimscroll',
    	'contextmenu': 'js/contextmenu/jquery.contextmenu-min',
    	'jquery.cookie': 'js/jquery-cookie',
    	'page': 'js/app/page-min',
    	'note': 'js/app/note-min',
    	'notebook': 'js/app/notebook-min',
    	'tag': 'js/app/tag-min',
    	'share': 'js/app/share-min',
    	'objectId': 'js/object_id-min',
    	'ZeroClipboard': 'js/ZeroClipboard/ZeroClipboard-min',
    	'bootstrap': 'js/bootstrap-min',
    	'leanote': 'js/main',
    	
    	// ajax upload image/attach
    	'editor_drop_paste': 'js/app/editor_drop_paste-min',
    	'attachment_upload': 'js/app/attachment_upload-min',
    	'jquery.ui.widget': 'tinymce/plugins/leaui_image/public/js/jquery.ui.widget',
    	'fileupload': '/tinymce/plugins/leaui_image/public/js/jquery.fileupload',
    	'iframe-transport': '/tinymce/plugins/leaui_image/public/js/jquery.iframe-transport',
    	
    	// mdeditor
    	'Markdown.Converter': 'mdeditor/editor/pagedown/Markdown.Converter-min',
    	'Markdown.Sanitizer': 'mdeditor/editor/pagedown/Markdown.Sanitizer-min',
    	'Markdown.Editor': 'mdeditor/editor/pagedown/Markdown.Editor',
    	'Markdown.zh': 'mdeditor/editor/pagedown/local/Markdown.local.zh-min',
    	'Markdown.en': 'mdeditor/editor/pagedown/local/Markdown.local.en-min',
    	'Markdown.Extra': 'mdeditor/editor/Markdown.Extra-min',
    	'underscore': 'mdeditor/editor/underscore-min',
    	'scrollLink': 'mdeditor/editor/scrollLink-min',
    	'mathJax': 'mdeditor/editor/mathJax-min',
    	'jquery.waitforimages': 'mdeditor/editor/jquery.waitforimages-min',
    	'pretty': 'mdeditor/editor/google-code-prettify/prettify',
    	'mdeditor': 'mdeditor/editor/mdeditor',
    	
    	'jquery.mobile': 'js/jquery.mobile-1.4.4.min',
    	'fastclick': 'js/fastclick'
    },
    shim: {
    	'page': {deps: ['tinymce']},
    	'fileupload': {deps: ['jquery.ui.widget', 'iframe-transport']},
    	'Markdown.Sanitizer': {deps: ['Markdown.Converter']},
    	'Markdown.Editor': {deps: ['Markdown.Converter']},
    	'Markdown.Extra': {deps: ['Markdown.Editor']},
    	'Markdown.zh': {deps: ['Markdown.Editor']},
    	'Markdown.en': {deps: ['Markdown.Editor']}
    }
});

/*
// leanote, 这里使用requireJs很慢, 不用
define('leanote', ['tinymce', 'page'], function(){
});

require(['jquery.slimscroll', 'contextmenu', 'jquery.cookie', 
'note', 'notebook', 'share', 'tag', 'objectId', 'ZeroClipboard', 'bootstrap'], function() {
	// 没有执行
	Notebook.renderNotebooks(notebooks);
	Share.renderShareNotebooks(sharedUserInfos, shareNotebooks);
	
	Note.renderNotes(notes);
	if(!isEmpty(notes)) {
		Note.changeNote(notes[0].NoteId);
	}
	
	Note.setNoteCache(noteContentJson);
	Note.renderNoteContent(noteContentJson)
	
	Tag.renderTagNav(tagsJson);
	
	// init notebook后才调用
	require(['page'], function() {
		initSlimScroll();
	});	
});
*/

// require(['mdeditor'], function(mdeditor) {});
require(['editor_drop_paste'], function(leaui_image) {});
require(['attachment_upload'], function(attachment_upload) {});