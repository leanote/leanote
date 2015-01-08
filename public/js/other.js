require.config({
    paths: {
    	'editor_drop_paste': '/public/js/app/editor_drop_paste-min',
    	'attachment_upload': '/public/js/app/attachment_upload-min',
    	'jquery.ui.widget': '/public/tinymce/plugins/leaui_image/public/js/jquery.ui.widget',
    	'fileupload': '/public/tinymce/plugins/leaui_image/public/js/jquery.fileupload',
    	'iframe-transport': '/public/tinymce/plugins/leaui_image/public/js/jquery.iframe-transport',
    },
    shim: {
    	'fileupload': {deps: ['jquery.ui.widget', 'iframe-transport']},
    }
});
require(['editor_drop_paste'], function(leaui_image) {});
require(['attachment_upload'], function(attachment_upload) {});