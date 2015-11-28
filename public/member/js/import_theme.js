define('import_theme', ['fileupload'], function(){
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
    
    setDropStyle("#dropAvatar", "#uploadAvatar");
    
	var initUploader = function() {
	    $('.dropzone .btn-choose-file').click(function() {
	        $(this).parent().find('input').click();
	    });
	
	    var $msg2 = $('#avatarUploadMsg');
	    $('#uploadAvatar').fileupload({
	        dataType: 'json',
	        dropZone: $('#dropAvatar'),
	        add: function(e, data) {
	            var tpl = $('<div class="alert alert-info"><img class="loader" src="/images/ajax-loader.gif"> <a class="close" data-dismiss="alert">×</a></div>');
	
	            // Append the file name and file size
	            tpl.append(data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>]</small>');
	
	            // Add the HTML to the UL element
	            $msg2.html(tpl);
	            data.context = $msg2;
	            
	            // 检查文件大小
	            var size = data.files[0].size;
	            if(typeof size == 'number' && size > 10 * 1024 * 1024) {
	            	tpl.find("img").remove();
	            	tpl.removeClass("alert-info").addClass("alert-danger");
	            	tpl.append("Warning: File size is bigger than 10M");
	            	return;
	            }
	            
	            // Automatically upload the file once it is added to the queue
	            var jqXHR;
	            setTimeout(function() {
		            jqXHR = data.submit();
	            }, 10);
	        },
	        done: function(e, data) {
	            if (data.result.Ok == true) {
	                data.context.html("");
	                var re = data.result;
	            	art.tips("Success");
	            	setTimeout(function() {
		            	location.reload();
	            	}, 1000);
	            } else {
	                var re = data.result;
	                data.context.html("");
	                var tpl = $('<div class="alert alert-danger"><a class="close" data-dismiss="alert">×</a></div>');
	                tpl.append('<b>Error:</b> ' + data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>]</small> ' + data.result.Msg);
	                data.context.html(tpl);
	            }
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
	        }
	    });
	}
	
	initUploader();
});