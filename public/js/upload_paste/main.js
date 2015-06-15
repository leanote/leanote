// 上传, 粘贴图片
requirejs.config({
	paths: {
		// life
		'editor_drop_paste': 'js/upload_paste/editor_drop_paste',
    	'attachment_upload': 'js/upload_paste/attachment_upload',
    	'jquery.ui.widget': 'js/upload_paste/libs/jquery.ui.widget',
    	'fileupload': 'js/upload_paste/libs/jquery.fileupload',
    	'iframe-transport': 'js/upload_paste/libs/jquery.iframe-transport'
	},
	shim: {
		// life
    	'fileupload': {deps: ['jquery.ui.widget', 'iframe-transport']},
	}
});
require(["editor_drop_paste", "attachment_upload"], function() {
});