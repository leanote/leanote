// for editor.
// drag image to editor
define('editor_drop_paste', ['fileupload'], function() {

	// 在toggle成pre或ace时
	// 最后没有元素, 或者元素不是p, 则在最后插入之
	function insertPIfNotExists() {
		var children = $('#editorContent').children();
		var lastChild = children && children.length > 0 ? children[children.length - 1] : null;
		if (!lastChild || lastChild.tagName != 'P') {
			$('#editorContent').append('<p><br data-mce-bogus="1"></p>');
		}
	}

	// 粘贴图片的进度控制
	function Process(editor) {
		var id = '__mcenew' + (new Date()).getTime();
		var str = '<div contenteditable="false" id="' + id + '" class="leanote-image-container">' + 
			'<img class="loader" src="/images/ajax-loader.gif">' + 
				'<div class="progress">' + 
					'<div class="progress-bar progress-bar-success progress-bar-striped" role="progressbar" aria-valuenow="2" aria-valuemin="0" aria-valuemax="100" style="width: 0%;">' + 
						'0%' + 
					'</div>' + 
				'</div>' + 
			'</div>';
		this.containerStr = str;
		
		editor.insertContent(str);
		insertPIfNotExists();
		
		var container = $('#' + id);
		this.container = container;
		this.id = id;
		this.processBar = container.find('.progress-bar');
	}
	Process.prototype.update = function(process) {
		var me = this;
		// 98%, 不要小数
		process = Math.ceil(process * 100);
		if(process >= 100) {
			process = 99;
		}
		process += "%";
		$('#' + me.id + ' .progress-bar').html(process).css('width', process);
	}
	Process.prototype.replace = function(src) {
		var me = this;
		getImageSize(src, function() {
			$('#' + me.id).replaceWith('<img src="' + src + '" />');
		});
	}
	Process.prototype.remove = function() {
		var me = this;
		$('#' + me.id).remove();
	}

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
	function insertImage(data) {
		var editor = tinymce.activeEditor;
		var dom = editor.dom;
	
		var renderImage = function(data2) {
			// 这里, 如果图片宽度过大, 这里设置成500px
			var d = {};
			var imgElm;
			// 先显示loading...
			d.id = '__mcenew' + (i++);
			d.src = "/images/loading-24.gif";
			imgElm = dom.createHTML('img', d);
			tinymce.activeEditor.insertContent(imgElm);
			imgElm = dom.get(d.id);
			
			function callback (wh) {
				dom.setAttrib(imgElm, 'src', data2.src);
				// dom.setAttrib(imgElm, 'width', data2.width);
				if(data2.title) {
					dom.setAttrib(imgElm, 'title', data2.title);
				}
				
				dom.setAttrib(imgElm, 'id', null);
			};
			getImageSize(data.src, callback);
		}
		
		//-------------
		// outputImage?fileId=123232323
		var fileId = "";
		fileIds = data.src.split("fileId=")
		if(fileIds.length == 2 && fileIds[1].length == "53aecf8a8a039a43c8036282".length) {
			fileId = fileIds[1];
		}
		if(fileId) {
			// 得到fileId, 如果这个笔记不是我的, 那么肯定是协作的笔记, 那么需要将图片copy给原note owner
			var curNote = Note.getCurNote();
			if(curNote && curNote.UserId != UserInfo.UserId) {
				(function(data) {
					ajaxPost("/file/copyImage", {userId: UserInfo.UserId, fileId: fileId, toUserId: curNote.UserId}, function(re) {
						if(reIsOk(re) && re.Id) {
							data.src = "/api/file/getImage?fileId=" + re.Id;
						}
						renderImage(data);
					});
				})(data);
			} else {
				renderImage(data);
			}
		} else {
			renderImage(data);
		}
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
	        pasteZone: '', // 不允许paste
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
	            var tpl = $('<li><div class="alert alert-info"><img class="loader" src="/images/ajax-loader.gif"> <a class="close" data-dismiss="alert">×</a></div></li>');
	
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
	                // life
	                var data2 = {src: "/api/file/getImage?fileId=" + data.result.Id}
	                insertImage(data2);
	            } else {
	                data.context.empty();
	                var tpl = $('<li><div class="alert alert-danger"><a class="close" data-dismiss="alert">×</a></div></li>');
	                tpl.find('div').append('<b>' + getMsg('Error') + ':</b> ' + data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>]</small> ' + data.result.Msg);
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
	            tpl.find('div').append('<b>' + getMsg('Error') + ':</b> ' + data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>]</small> ' + data.errorThrown);
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
	    	$("#upload").css("z-index", 0).css("top", "auto").hide();
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
		        
		        // 如果在只读状态, 转换之
			    if (LEA.readOnly) {
			    	LEA.toggleWriteable();
			    }
			    
		    } else {
		        dropZone.removeClass('hover');
		    }
		    
		    window.dropZoneTimeout = setTimeout(function () {
		        window.dropZoneTimeout = null;
		        dropZone.removeClass('in hover');
		        hideUpload();
		    }, 500);
		});
	};

	var lastTime = 0;

	// pasteImage
	var pasteImageInit =  function() {
		var curNote;
	    // Initialize the jQuery File Upload plugin
	    var dom, editor;
	    // 2015/4/17 添加wmd-input markdown paste image
	    $('#editorContent, #left-column').fileupload({
	        dataType: 'json',
	        pasteZone: $('#editorContent, #left-column'),
	        dropZone: '', // 只允许paste
	        maxFileSize: 210000,
	        url: "/file/pasteImage",
	        paramName: 'file',
	        formData: function(form) {
	        	return [{name: 'from', value: 'pasteImage'}, {name: 'noteId', value: Note.curNoteId}]
	        },
	        /*
	        paste: function(e, data) {
	        	var jqXHR = data.submit();
	        },
	        */
	        progress: function(e, data) {
	        	if(curNote && !curNote.IsMarkdown) {
		        	data.process.update(data.loaded / data.total);
	        	}
	        },

	        // 调用了两次
	        // 不知道为什么会触发两次
	        add: function(e, data) {
	        	// 防止两次
        		// console.trace(e);
	        	var now = (new Date()).getTime();
	        	if (now - lastTime < 500) {
	        		return;
	        	}
	        	// console.log('nono');
	        	lastTime = now;

	        	var note = Note.getCurNote();
	        	curNote = note;
	        	if(!note || note.IsNew) {
	        		// alert(getMsg("Please save note firstly!"));
	        		// return;
	        	}
	        	
	        	// LEA.removePasteBin();
	        	// 为什么要延迟? 为了让paste plugin先执行, 删除掉paste bin
	        	setTimeout(function () {
		        	// 先显示loading...
					editor = tinymce.EditorManager.activeEditor; 
					if(!note.IsMarkdown) {
						var process = new Process(editor);
					}
					data.process = process;
		            var jqXHR = data.submit();
	        	}, 20);
	        },
	
	        done: function(e, data) {
	            if (data.result.Ok == true) {
		    		// 这里, 如果图片宽度过大, 这里设置成500px
		    		var re = data.result;
					var src = "/api/file/getImage?fileId=" + re.Id;

					if(curNote && !curNote.IsMarkdown) {
						data.process.replace(src);
					} else {
						MD && MD.insertLink(src, 'title', true);
					}
				
	            } else {
					data.process.remove();
	            }
	        },
	        fail: function(e, data) {
	        	if(curNote && !curNote.IsMarkdown) {
					data.process.remove();
				}
	        }
	    });
	};
	
	initUploader();
	pasteImageInit();
});
