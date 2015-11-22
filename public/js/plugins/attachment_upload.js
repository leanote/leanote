// upload attachment
// 依赖note
var urlPrefix = UrlPrefix;
define('attachment_upload', ['fileupload'], function(){
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
    
    function setDropStyle(dropzoneId, formId) {
	    // drag css
	    var dropZone = $(dropzoneId);
		$(formId).bind('dragover', function (e) {
			e.preventDefault();
		    var timeout = window.dropZoneTimeoutAttach;
		    if(timeout) {
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
		    window.dropZoneTimeoutAttach = setTimeout(function () {
		        window.dropZoneTimeoutAttach = null;
		        dropZone.removeClass('in hover');
		    }, 100);
		});
    }
    
    setDropStyle("#dropAttach", "#uploadAttach");
    
	var initUploader = function() {
	    $('.dropzone .btn-choose-file').click(function() {
	        $(this).parent().find('input').click();
	    });

		var $msg = $('#attachUploadMsg');
	    // Initialize the jQuery File Upload plugin
	    $('#uploadAttach').fileupload({
	        dataType: 'json',
	        pasteZone: '', // 不能通过paste来上传图片
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

	            var tpl = $('<div class="alert alert-info"><img class="loader" src="/images/ajax-loader.gif"> <a class="close" data-dismiss="alert">×</a></div>');
	
	            // Append the file name and file size
	            tpl.append(data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>]</small>');
	
	            // Add the HTML to the UL element
	            $msg.html(tpl);
	            data.context = $msg;
	            
	            // 检查文件大小
	            var size = data.files[0].size;
	            var maxFileSize = +GlobalConfigs["uploadAttachSize"] || 100;
	            if(typeof size == 'number' && size > 1024 * 1024 * maxFileSize) {
	            	tpl.find("img").remove();
	            	tpl.removeClass("alert-info").addClass("alert-danger");
	            	tpl.append(" Warning: File size is bigger than " + maxFileSize + "M");
	            	setTimeout((function(tpl) {
	                	return function() {
		                	tpl.remove();
	                	}
	                })(tpl), 3000);
	            	return;
	            }
	            
	            // Automatically upload the file once it is added to the queue
	            var jqXHR;
	            setTimeout(function() {
		            jqXHR = data.submit();
	            }, 10);
	        },
	        /*
	        progress: function (e, data) {
	        },
	        */
	        done: function(e, data) {
	            if (data.result.Ok == true) {
	                data.context.html("");
	                Attach.addAttach(data.result.Item);
	            } else {
	                var re = data.result;
	                data.context.html("");
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
                data.context.html("");
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
	};
	
	initUploader();
});
