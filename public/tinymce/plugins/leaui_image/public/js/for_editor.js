// for editor.
// drag image to editor
// Copyright leaui

define('leaui_image', ['jquery.ui.widget', 'fileupload'], function(){
	var editor = tinymce.activeEditor;
	var dom = editor.dom;
	
	// 当url改变时, 得到图片的大小
	function getImageSize(url, callback) {
		var img = document.createElement('img');
	
		function done(width, height) {
			img.parentNode.removeChild(img);
			callback({width: width, height: height});
		}
	
		img.onload = function() {
			done(img.clientWidth, img.clientHeight);
		};
	
		img.onerror = function() {
			done();
		};
	
		img.src = url;
	
		var style = img.style;
		style.visibility = 'hidden';
		style.position = 'fixed';
		style.bottom = style.left = 0;
		style.width = style.height = 'auto';
	
		document.body.appendChild(img);
	}
	
	var i = 1;
	function insertImage(data2) {
		// 这里, 如果图片宽度过大, 这里设置成500px
		var d = {};
		var imgElm;
		// 先显示loading...
		d.id = '__mcenew' + (i++);
		d.src = "http://leanote.com/images/loading-24.gif";
		imgElm = dom.createHTML('img', d);
		editor.insertContent(imgElm);
		imgElm = dom.get(d.id);
		
		function callback (wh) {
			if(wh && wh.width) {
				if(wh.width > 600) {
					wh.width = 600;
				}
				data2.width = wh.width;
			}
			dom.setAttrib(imgElm, 'src', data2.src);
			dom.setAttrib(imgElm, 'width', data2.width);
			if(data2.title) {
				dom.setAttrib(imgElm, 'title', data2.title);
			}
			
			dom.setAttrib(imgElm, 'id', null);
		};
		getImageSize(data2.src, callback);
	}
	
	var initUploader =  function() {
		var ul = $('#upload ul');
	
	    $('#drop a').click(function() {
	        // trigger to show file select
	        $(this).parent().find('input').click();
	    });
	
	    // Initialize the jQuery File Upload plugin
	    $('#upload').fileupload({
	        dataType: 'json',
	        acceptFileTypes: /(\.|\/)(gif|jpg|jpeg|png|jpe)$/i,
	        maxFileSize: 210000,
	
	        // This element will accept file drag/drop uploading
	        dropZone: $('#drop'),
	        formData: function(form) {
	        	return [{name: 'albumId', value: ""}]
	        },
	        // This function is called when a file is added to the queue;
	        // either via the browse button, or via drag/drop:
	        add: function(e, data) {
	            var tpl = $('<li><div class="alert alert-info"><img class="loader" src="/tinymce/plugins/leaui_image/public/images/ajax-loader.gif"> <a class="close" data-dismiss="alert">×</a></div></li>');
	
	            // Append the file name and file size
	            tpl.find('div').append(data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>]</small>');
	
	            // Add the HTML to the UL element
	            data.context = tpl.appendTo(ul);
	
	            // data.form[0].action += "&album_id=" + $("#albumsForUpload").val();
	
	            // Automatically upload the file once it is added to the queue
	            var jqXHR = data.submit();
	        },
	
	        done: function(e, data) {
	            if (data.result.Ok == true) {
	                data.context.remove();
	                var data2 = {src: data.result.Id}
	                insertImage(data2);
	            } else {
	                data.context.empty();
	                var tpl = $('<li><div class="alert alert-danger"><a class="close" data-dismiss="alert">×</a></div></li>');
	                tpl.find('div').append('<b>Error:</b> ' + data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>]</small> ' + data.result.Msg);
	                data.context.append(tpl);
	                setTimeout((function(tpl) {
	                	return function() {
		                	tpl.remove();
	                	}
	                })(tpl), 2000);
	            }
	            $("#uploadMsg").scrollTop(1000);
	        },
	        fail: function(e, data) {
	            data.context.empty();
	            var tpl = $('<li><div class="alert alert-danger"><a class="close" data-dismiss="alert">×</a></div></li>');
	            tpl.find('div').append('<b>Error:</b> ' + data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>]</small> ' + data.errorThrown);
	            data.context.append(tpl);
	            setTimeout((function(tpl) {
	                	return function() {
		                	tpl.remove();
	                	}
	                })(tpl), 2000);
	
	            $("#uploadMsg").scrollTop(1000);
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
	    
	    function showUpload() {
	    	$("#upload").css("z-index", 12);
	    	var top = +$("#mceToolbar").css("height").slice(0, -2); // px
	    	$("#upload").css("top", top - 8);
	    	$("#upload").show();
	    }
	    
	    function hideUpload() { 
	    	$("#upload").css("z-index", 0).css("top", "auto");
	    }
	    
	    // drag css
		$(document).bind('dragover', function (e) {
		    var dropZone = $('#drop'),
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