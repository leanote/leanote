// upload attachment
// 依赖note
var urlPrefix = window.location.protocol + "//" + window.location.host;
define('attachment_upload', ['jquery.ui.widget', 'fileupload'], function(){
	// var editor = tinymce.activeEditor;
	// var dom = editor.dom;
	
	var initUploader =  function() {
		var $msg = $('#attachUploadMsg');
	
	    $('#dropAttach .btn-choose-file').click(function() {
	        // trigger to show file select
	        $(this).parent().find('input').click();
	    });
	
	    // Initialize the jQuery File Upload plugin
	    $('#uploadAttach').fileupload({
	        dataType: 'json',
	        maxFileSize: 210000,
	        // This element will accept file drag/drop uploading
	        dropZone: $('#dropAttach'),
	        formData: function(form) {
	        	return [{name: 'noteId', value: Note.curNoteId}] // 传递笔记本过去
	        },
	        // This function is called when a file is added to the queue;
	        // either via the browse button, or via drag/drop:
	        add: function(e, data) {
	        	var note = Note.getCurNote();
	        	if(!note || note.IsNew) {
	        		alert("This note hasn't saved, please save it firstly!")
	        		return;
	        	}
	            var tpl = $('<div class="alert alert-info"><img class="loader" src="/tinymce/plugins/leaui_image/public/images/ajax-loader.gif"> <a class="close" data-dismiss="alert">×</a></div>');
	
	            // Append the file name and file size
	            tpl.append(data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>]</small>');
	
	            // Add the HTML to the UL element
	            tpl.appendTo($msg);
	            data.context = $msg;
	            
	            // Automatically upload the file once it is added to the queue
	            var jqXHR;
	            setTimeout(function() {
		            jqXHR = data.submit();
	            }, 0);
	        },
	
	        done: function(e, data) {
	            if (data.result.Ok == true) {
	                data.context.remove();
	                Attach.addAttach(data.result.Item);
	            } else {
	                var re = data.result;
	                data.context.empty();
	                var tpl = $('<div class="alert alert-danger"><a class="close" data-dismiss="alert">×</a></div>');
	                tpl.append('<b>Error:</b> ' + data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>]</small> ' + data.result.Msg);
	                data.context.html(tpl);
	                setTimeout((function(tpl) {
	                	return function() {
		                	tpl.remove();
	                	}
	                })(tpl), 3000);
	            }
	            $("#uploadAttachMsg").scrollTop(1000);
	        },
	        fail: function(e, data) {
	            data.context.empty();
	            var tpl = $('<div class="alert alert-danger"><a class="close" data-dismiss="alert">×</a></div>');
	            tpl.append('<b>Error:</b> ' + data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>]</small> ' + data.errorThrown);
	            data.context.html(tpl);
	            setTimeout((function(tpl) {
	                	return function() {
		                	tpl.remove();
	                	}
	                })(tpl), 3000);
	
	            $("#uploadAttachMsg").scrollTop(1000);
	        }
	    });
	
	    // Prevent the default action when a file is dropped on the window
	    $(document).on('drop dragover', function(e) {
	        e.preventDefault();
	    });
	
	    // Helper function that formats the file sizes
	    function formatFileSize(bytes) {
	        if (typeof bytes !== 'number') {
	            return '';
	        }
	        if (bytes >= 1000000000) {
	            return (bytes / 1000000000).toFixed(2) + ' GB';
	        }
	        if (bytes >= 1000000) {
	            return (bytes / 1000000).toFixed(2) + ' MB';
	        }
	        return (bytes / 1000).toFixed(2) + ' KB';
	    }
	    
	    // drag css
		$(document).bind('dragover', function (e) {
		    var dropZone = $('#dropAttach'),
		        timeout = window.dropZoneTimeout;
		    if (!timeout) {
		        dropZone.addClass('in');
		        showUpload();
		    } else {
		        clearTimeout(timeout);
		    }
		    var found = false,
		        node = e.target;
		    do {
		        if (node === dropZone[0]) {
		            found = true;
		            break;
		        }
		        node = node.parentNode;
		    } while (node != null);
		    if (found) {
		        dropZone.addClass('hover');
		    } else {
		        dropZone.removeClass('hover');
		    }
		    window.dropZoneTimeout = setTimeout(function () {
		        window.dropZoneTimeout = null;
		        dropZone.removeClass('in hover');
		        hideUpload();
		    }, 100);
		});
	}
	
	initUploader();
});