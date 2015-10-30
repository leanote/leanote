function log(o) {
	// console.log(o);
}
function retIsOk(ret) {
	if(ret && typeof ret == "object" && ret.Ok == 1) {
		return true;
	}
	return false;
}

var urlPrefix = UrlPrefix;
var getMsg = parent.getMsg;
if (!getMsg) {
	getMsg = function(msg) {
		return msg;
	};
}

// load image
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

var o = {
	maxSelected: G.maxSelected, 
	selectedZoneO:$("#preview"),
	previewO: $("#preview"),
	selectedImages:[], // selected urls
	imageAttrs: {}, // src => {width, height, title}

	pageNum: 1,

	pagination: function(count) {
		var self = this;
		$(".pagination").pagination(count, {
			items_per_page: G.perPageItems,
			prev_text: getMsg('Prev'),
			next_text: getMsg('Next'),
			callback: function(pageNum) {
				self.pageNum = pageNum+1;
				self.renderImages($("#albumsForList").val(), self.pageNum, false);
			}
	    });
    },

    showMsg: function(msg) {
    	$("#msg").html(msg).css("display", "inline");
    	setTimeout(function() {
    		$("#msg").fadeOut();
    	}, 2000);
    },

    pageAddAlbum: function(ret) {
    	var html = '<option value="' + ret.AlbumId + '">' + ret.Name + '</option>';
    	$("#albumsForUpload").append(html).val(ret.AlbumId);
    	$("#albumsForList").append(html);
    },
    pageUpdateAlbum: function(albumId, albumName) {
    	$('option[value="' + albumId + '"]').html(albumName);
    },
    // add, delete, update album
    processAlbum: function() {
	    var self = this;
		var isAddAlbum = true;
		var curAlbum = "";
		function toggleAddAlbum() {
			if($("#addOrUpdateAlbumForm").is(":hidden")) {
				$("#addOrUpdateAlbumForm").show();
				$("#albumSelect").hide();
			} else {
				$("#addOrUpdateAlbumForm").hide();
				$("#albumSelect").show();
			}
		}
		// rename
		$("#renameAlbumBtn").click(function(){ 
			curAlbum = $("#albumsForUpload").val();
			if(!curAlbum) {
				alert(getMsg("Cannot rename default album"));
				return;
			}
			toggleAddAlbum();
			$("#addOrUpdateAlbumBtn").html(getMsg("Rename Album"));
			$("#albumName").val($("#albumsForUpload option:selected").html()).focus();
			isAddAlbum = false;
		});
		// add album
	    $("#addAlbumBtn").click(function() {
			toggleAddAlbum();
			$("#addOrUpdateAlbumBtn").html(getMsg("Add Album"));
    		$("#albumName").val("").focus();
			isAddAlbum = true;
		});
		$("#cancelAlbumBtn").click(function() {
			toggleAddAlbum();
		});
		// add or update album
		$("#addOrUpdateAlbumBtn").click(function() {
	    	var albumName = $("#albumName").val();
	    	if(!albumName) {
	    		$("#albumName").focus();
	    		return;
	    	}
	    	if(isAddAlbum) {
		    	$.get("/album/addAlbum", {name: albumName}, function(ret) {
		    		if(typeof ret == "object" && ret.AlbumId != "") {
		    			$("#albumName").val("");
		    			self.showMsg(getMsg("Add Success!"));
		    			self.pageAddAlbum(ret);
		    			
						setTimeout(function() {
							toggleAddAlbum();
						}, 200);
		    		} else {
		    			alert(getMsg("error"));
		    		}
		    	});
	    	} else {
		    	$.get("/album/updateAlbum", {albumId: curAlbum, name: albumName}, function(ret) {
		    		if(typeof ret == "boolean" && ret) {
		    			$("#albumName").val("");
		    			self.showMsg(getMsg("Rename Success!"));
		    			self.pageUpdateAlbum(curAlbum, albumName);
		    			
						setTimeout(function() {
							toggleAddAlbum();
						}, 200);
		    		} else {
		    			alert(getMsg("error!"));
		    		}
		    	});
	    	}
	    })
	    // delete album
	    $("#deleteAlbumBtn").click(function() {
	    	var albumId = $("#albumsForUpload").val();
	    	if(!albumId) {
	    		alert(getMsg("Cannot delete default album"));
	    		return;
	    	}
	    	$.get("/album/deleteAlbum", {albumId: albumId}, function(ret) {
	    		if(typeof ret == "object" && ret.Ok == true) {
	    			self.showMsg(getMsg("Delete Success!"));
	    			// delete this album from select
	    			$("#albumsForUpload option[value='" + albumId + "']").remove();

	    			// if the albumsForList has select this album, it must refresh list after delete it;
	    			if($("#albumsForList").val() == albumId) {
		    			self.needRefresh = true;
	    			}

	    			$("#albumsForList option[value='" + albumId + "']").remove();

	    		} else {
	    			alert(getMsg("This album has images, please delete it's images at first."));
	    		}
	    	});

	    });
	},

    renderAlbums: function() {
    	var self = this;
    	$.get("/album/getAlbums", function(ret) {
    		if(!ret) return;
    		var html = "";
    		for(var i in ret) {
    			var each = ret[i];
		    	var option = '<option value="' + each.AlbumId + '">' + each.Name+ '</option>';
		    	html += option;
    		}

	    	$("#albumsForUpload").append(html);
	    	$("#albumsForList").append(html);

	    	var albumId = $("#albumsForList").val();
		    self.renderImages(albumId, 1, true);
    	});
    },

    imageMaskO: $("#imageMask"),
    noImagesO: $("#noImages"),
   	loadingO: $("#loading"),

    loadingStart: function() {
    	if(this.imageMaskO.is(":hidden")) {
	    	this.imageMaskO.css("opacity", 0.8).show();
	    }
    	this.noImagesO.hide();
    	this.loadingO.show();
    },
    loadingEnd: function() {
    	this.imageMaskO.hide();
    },

    noImages: function () {
    	this.imageMaskO.show().css("opacity", 1);
    	this.noImagesO.show();
    	this.loadingO.hide();
    },


	search: function() {
		var self = this;
		var t1 = 1;
		$("#key").on("keyup", function() {
			var t2 = ++t1;
			var key = $(this).val();
			var albumId = $("#albumsForList").val();

			self.renderImages(albumId, 1, true, key, function(){
				return t1 == t2;
			});
		});
	},
    
	renderImages: function(albumId, page, needRenderPagination, key, needRender) {
		var self = this;

    	if(!page) {
    		page = 1;
    	}
    	self.loadingStart();
    	$.get("/file/getImages", {albumId: albumId, page: page, key: key}, function(ret) {
    		if(!ret || !ret.Count) {
    			self.noImages();
    			return;
    		}
    		self.loadingEnd();
    		var datas = ret.List;
    		var selectedMap = {};

			for(var i in self.selectedImages) {
				var src = self.selectedImages[i]; // src include G.imageSrcPrefix
				selectedMap[src] = true;
			}
			// log(self.selectedImages);

			var html = "";
			for(var i in datas){
				var each = datas[i];
				var classes = "";
				// life edit
				// 之前的
				if(each.Path != "" && each.Path[0] == "/") {
					each.Path = each.Path.substr(1);
				}
				if(each.Path != "" && each.Path.substr(0, 7) == "upload/") {
					var src = urlPrefix + "/" + each.Path;
				} else {
					var src = urlPrefix + "/api/file/getImage?fileId=" + each.FileId;
				}
				// log(src);
				if(selectedMap[src]) {
					classes = 'class="selected"';
				}
				html += '<li ' + classes + '>';
				html += '<a title="" href="javascript:;" class="a-img"><img  alt="" src="' + src + '" data-original="' + src + '" ></a>';
				// html += '<div class="tools"><a href="javascript:;" class="del" data-id="' + each.FileId + '"><span class="fa fa-trash"></span></a></div>';
				html += '<div class="tools clearfix" data-id="' + each.FileId + '"><div class="file-title pull-left">' + each.Title + '</div><div class="pull-right"><a href="javascript:;" class="del" data-id="' + each.FileId + '"><span class="fa fa-trash"></span></a></div></div>';
				html += "</li>";
			}
				
    		// var html = $("#tImage").render(datas);
    		$("#imageList").html(html);

    		if(needRenderPagination) {
	    		self.pagination(ret.Count);
    		}

    		// $("#imageList img").lazyload({effect : "fadeIn"});
    		// $("#imageList img").lazyload();
    	});
    },

	// 初始化已选择的图片区域
	initSelectedZones: function() {
		var self = this;
		num = this.maxSelected;
		self.previewO.html("");
		for(var i = 1; i <= num; ++i) {
			// self.previewO.append("<li>" + i + "</li>");
			self.previewO.append("<li>?</li>");
		}
	},
		
	reRenderSelectedImages: function(isRemove, addSrc) {
		var self = this;

		var lis = this.selectedZoneO.find("li");
		var size = this.selectedImages.length-1;
		for(var i = 0; i < this.maxSelected; ++i) {
			var target = lis.eq(i);
			if(i > size) {
				target.html('?');
			} else {
				src = this.selectedImages[i];

				var data = self.imageAttrs[src];
				var attrs = "";
				if(data) {
					if(data.width) attrs += ' data-width="' + data.width + '"';
					if(data.height) attrs += ' data-height="' + data.height + '"';
					if(data.title) attrs += ' data-title="' + data.title + '"';
				}

				target.html('<img ' + attrs + ' src="' + src + '" width="60"/><div class="tools"><a title="' + getMsg('click to remove this image') + '" href="javascript:;" class="del"><span class="fa fa-trash"></span></a></div>');
			}

			// remove selected
			if(isRemove) {
				target.removeClass("selected");
			} else {
				// is add
				// trigger click and set attrs
				if(addSrc == src) {
					target.click();
				}
			}
		}
	},
	removeSelectedImage: function($li) {
		var self = this;

		var src = $li.find("img").attr('src');
		for(var i in this.selectedImages) {
			if(this.selectedImages[i] == src) {
				this.selectedImages.splice(i, 1)
			}
		}
		this.reRenderSelectedImages(true);

		// clear attrs and disable it
		self.clearAttrs();
	},
	addSelectedImage: function($li) {
		if(this.maxSelected > 1 && this.maxSelected <= this.selectedImages.length) {
			return false;
		}
		
		// life 为了图片安全
		if(typeof $li == "object") {
			var src = $li.find("img").attr('src');
		} else {
			// 也有可能来自url
			if($li.indexOf("http://") != -1 || $li.indexOf("https://") != -1) {
				src = $li;
			} else {
				// 来自内部
				src = urlPrefix + "/api/file/getImage?fileId=" + $li;
			}
		}
		
		// 如果只允许选1个
		if(this.maxSelected == 1) {
			// 先把其它的去掉
			$("#imageList li").removeClass("selected");
			this.selectedImages = [src];
		} else {
			this.selectedImages.push(src);
		}
		
		this.reRenderSelectedImages(false, src);
	
		return true;
	},
	initDataFromTinymce: function() {
		var self = this;
		var datas = top.LEAUI_DATAS;
		var lastSrc = "";
		if(datas && datas.length > 0) {
			for(var i in datas) {
				var data = datas[i];
				data.constrain = true;
				lastSrc = data.src;
				self.selectedImages.push(data.src);
				self.imageAttrs[data.src] = data;
			}

			self.reRenderSelectedImages(false, lastSrc);
		}
	},

	init: function() {	
		var self = this;
		
		self.processAlbum();

		$("#albumsForList").change(function() {
			var albumId = $(this).val();
			self.renderImages(albumId, 1, true);
		});
		
		$("#imageList").on("click", 'li', function() {
			if($(this).hasClass("selected")) {
				$(this).removeClass("selected");
				self.removeSelectedImage($(this));
			} else {
				if(self.addSelectedImage($(this))){ 
					$(this).addClass("selected");
				}
			}	
		});


		// delete file
		$("#imageList").on("click", '.del', function(e) {
			var t = this;
			e.stopPropagation();

			if(confirm(getMsg("Are you sure to delete this image ?"))) {
				var fileId = $(this).data('id');
				$.get("/file/deleteImage", {fileId: fileId}, function(ret) {
					if(ret) {
						var $li = $(t).closest('li');
						if($li.hasClass("selected")) {
							self.removeSelectedImage($li);
						}
						$(t).closest('li').remove();
					}
				});
			}	
		});
		// edit file title
		$("#imageList").on("click", '.file-title', function(e) {
			var p = this;
			var t = this;
			e.stopPropagation();
			if($(this).children().eq(0).is("input")){
				return;
			}
			var fileId = $(p).parent().data('id');
			var fileTitle = $(this).text();
			$(this).html('<input type="text" value="' + fileTitle + '" />');

			var $input = $(this).find("input");
			$input.focus();
			$input.keydown(function(e){
				if(e.keyCode==13){
					$(this).trigger("blur");
				}
			});
			$input.blur(function() {
				var title = $(this).val();
				if(!title) {
					title = fileTitle;
				} else {
					$.post("/file/updateImageTitle", {fileId: fileId, title: title});
				}
				$(p).html(title);
			});
		});

		// remove preview
		$("#preview").on("click", '.del', function(e) {
			e.stopPropagation();
			var $li = $(this).closest("li");
			var src = $li.find("img").attr("src");
			self.removeSelectedImage($li);

			// 在当前的imagesList下看是否有
			$("#imageList img").each(function() {
				var src2 = $(this).attr('src');
				if(src == src2) {
					$(this).parent().parent().removeClass("selected");
				}
			});
		});

		// 
		$("#goAddImageBtn").click(function() {
			$("#albumsForUpload").val($("#albumsForList").val());
			$('#myTab li:eq(1) a').tab('show');
		});

		// toggle tab
		// refresh 
		$('#myTab a').on('shown.bs.tab', function(e) {
			e.preventDefault()
			$(this).tab('show');
			var href = $(this).attr("href");

			if(self.needRefresh && href == "#images") {
				setTimeout(function(){
					var albumId = $("#albumsForList").val();
					var key = $("#key").val();
					self.renderImages(albumId, self.pageNum, true, key);
				}, 200);
				self.needRefresh = false;
			}

			if(href == "#url") {
				$("#imageUrl").focus();
			}
		});
		$("#refresh").click(function() {
			var albumId = $("#albumsForList").val();
			var key = $("#key").val();
			self.renderImages(albumId, self.pageNum, false, key);
		});

		// add url
		$("#addImageUrlBtn").click(function(e) {
			e.preventDefault();
			var url = $.trim($("#imageUrl").val());
			if(!url) {
				$("#imageUrl").focus();
				return;
			}

			getImageSize(url, function(ret) {
				if(!ret.width || !ret.height){
					$("#msgForUrl").show();
					return;
				}
				$("#msgForUrl").hide();
				$("#imageUrl").val("");
				self.addSelectedImage(url);
			});
		});

		// 设置属性
		$("#preview").on("click", 'li', function() {
			if($(this).hasClass("selected")) {
				// $(this).removeClass("selected");
			} else {
				if($(this).find("img").length){ 
					$("#preview li").removeClass("selected");
					$(this).addClass("selected");

					self.initAttr($(this));
				}
			}	
		});

		$("#attrTitle, #attrWidth, #attrHeight").on("keyup", function(){
			self.modifyAttr($(this));
		});
		$("#attrConstrain").on("click", function(){
			self.modifyAttr($(this));
		});

		// file search
		self.search();

		self.initSelectedZones();
		
		self.initDataFromTinymce();

		self.renderAlbums();
		//...
		self.initUploader();
	},

	// 设置
	curSrc: "",
	curLi: null,
	attrTitleO: $("#attrTitle"),
	attrWidthO: $("#attrWidth"),
	attrHeightO: $("#attrHeight"),
	attrConstrainO: $("#attrConstrain"),
	// clear attrs and disable it
	clearAttrs: function() {
		var self = this;
		self.attrTitleO.val("").attr("disabled", true);
		self.attrHeightO.val("").attr("disabled", true);
		self.attrWidthO.val("").attr("disabled", true);
		self.attrConstrainO.prop("checked", false).attr("disabled", true);
	},
	scale: function(isWidth) {
		var self = this;
		var autoScale = self.attrConstrainO.is(":checked");
		var width = +self.attrWidthO.val();
		var height = +self.attrHeightO.val();
		if(isNaN(width) || isNaN(height)) {
			return;
		}

		var curAttrs = self.getCurAttrs();
		var preWidth = curAttrs.preWidth || curAttrs.width;
		var preHeight = curAttrs.preHeight || curAttrs.height;
		
		if(autoScale && preWidth && preHeight) {
			if(isWidth) {
				height = parseInt((width/preWidth) * preHeight);
				self.attrHeightO.val(height);
			} else {
				width = parseInt((height/preHeight) * preWidth);
				self.attrWidthO.val(width);
			}
		}
		var ret = {width: width, height: height};
		return ret;
	},
	getCurAttrs: function() {
		var self = this;
		return self.imageAttrs[self.curSrc];
	},
	setCurDataAttrs: function(attrs) {
		var self = this;
		var img = self.curLi.find("img");
		img.attr("data-width", attrs.width);
		img.attr("data-height", attrs.height);
		img.attr("data-title", attrs.title);

		self.imageAttrs[self.curSrc] = attrs;
	},
	// 修改属性
	modifyAttr: function($target){
		var self = this;
		var id = $target.attr("id");
		var val = $target.val();
		var curAttrs = self.getCurAttrs();
		if(!curAttrs) {
			return;
		}
		switch(id) {
			case "attrConstrain":
				val = 0;
				if($target.is(":checked")) {
					val = 1;
				}
				curAttrs['constrain'] = val;
				break;
			case "attrTitle":
				curAttrs['title'] = val;
				break;
			case "attrWidth":
				$.extend(curAttrs, self.scale(true));
				break;
			case "attrHeight":
				$.extend(curAttrs, self.scale(false));
				break;
		}

		self.setCurDataAttrs(curAttrs);
	},
	// when click preview li
	initAttr: function($li) {
		var self = this;
		if(typeof $li != "object") {
			$li = $("#preview").find('img[src="' + $li + '"]').parent();
		};

		var src = $li.find("img").attr("src");
		self.curSrc = src;
		self.curLi = $li;

		var attrs = self.imageAttrs[src];
		function setAttr(attrs) {
			attrs = attrs || {};
			self.attrTitleO.val(attrs.title).attr("disabled", false);
			self.attrWidthO.val(attrs.width).attr("disabled", false);
			self.attrHeightO.val(attrs.height).attr("disabled", false);
			self.attrConstrainO.attr("disabled", false);

			if(attrs.constrain) {
				self.attrConstrainO.prop('checked', true);
			} else {
				self.attrConstrainO.prop('checked', false);
			}

			self.setCurDataAttrs(attrs);
		}
		attrs = attrs || {};
		if(!attrs || !attrs.width || !attrs.height) {
			getImageSize(src, function(ret) {
				ret.title = attrs.title || "";
				ret.constrain = 1;
				ret.preWidth = ret.width;
				ret.preHeight = ret.height;
				if(src != self.curSrc) {
					self.imageAttrs[src] = ret;
					// in case user click fast
					self.setCurDataAttrs(attrs);
					return;
				}
				// set attrs
				setAttr(ret);
				return;
			});
		} else {
			setAttr(attrs);
		}
	},

	needRefresh: false,

	uploadRefreshImageList: function() {
		var self = this;
		// check albumId
		var albumId = $("#albumsForList").val();
		if(albumId == $("#albumsForUpload").val()){
			self.needRefresh = true;
		}
	},

	initUploader: function() {
		var self = this;
		var ul = $('#upload ul');

	    $('#drop a').click(function() {
	        // trigger to show file select
	        $(this).parent().find('input').click();
	    });
	    // Initialize the jQuery File Upload plugin
	    $('#upload').fileupload({
	        dataType: 'json',
	        pasteZone: '',
	        acceptFileTypes: /(\.|\/)(gif|jpg|jpeg|png|jpe)$/i,
	        // maxFileSize: 210000,

	        // This element will accept file drag/drop uploading
	        dropZone: $('#drop'),
	        formData: function(form) {
	        	return [{name: 'albumId', value: $("#albumsForUpload").val()}]
	        },
	        /*
	        urlFunc: function() {
	        	return 'server/index.php?action=file:uploadImage&album_id=' + $("#albumsForUpload").val();
	        },
	        */

	        // This function is called when a file is added to the queue;
	        // either via the browse button, or via drag/drop:
	        add: function(e, data) {
	        	// 文件大小限制
				var size = data.files[0].size;
	            var maxFileSize = +parent.GlobalConfigs["uploadImageSize"] || 100;
	            if(typeof size == 'number' && size > 1024 * 1024 * maxFileSize) {
	                var tpl = $('<li><div class="alert alert-danger"><a class="close" data-dismiss="alert">×</a></div></li>');
	                tpl.find('div').append('<b>Warning:</b> ' + data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>] is bigger than ' + maxFileSize + 'M</small> ');
	                tpl.appendTo(ul);
	            	return;
	            }
	            
	            var tpl = $('<li><div class="alert alert-info"><img class="loader" src="/public/album/images/ajax-loader.gif"> <a class="close" data-dismiss="alert">×</a></div></li>');
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

	               	// add image to preview
	                self.addSelectedImage(data.result.Id);
	                // reresh image list
	                self.uploadRefreshImageList();
	            } else {
	                data.context.empty();
	                var tpl = $('<li><div class="alert alert-danger"><a class="close" data-dismiss="alert">×</a></div></li>');
	                tpl.find('div').append('<b>' + getMsg('Error') + ':</b> ' + data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>]</small> ' + data.result.Msg);
	                data.context.append(tpl);
	                setTimeout((function(tpl) {
	                	return function() {
		                	tpl.remove();
	                	}
	                })(tpl), 3000);
	            }
	            $("#upload-msg").scrollTop(1000);
	        },
	        fail: function(e, data) {
	            data.context.empty();
	            var tpl = $('<li><div class="alert alert-danger"><a class="close" data-dismiss="alert">×</a></div></li>');
	            tpl.find('div').append('<b>Error:</b> ' + data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>]</small> ' + data.errorThrown);
	            data.context.append(tpl);

	            $("#upload-msg").scrollTop(1000);
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
		    var dropZone = $('#drop'),
		        timeout = window.dropZoneTimeout;
		    if (!timeout) {
		        dropZone.addClass('in');
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
		    }, 100);
		});
	}
}

$(function() {
	o.init();
});

// 为md得到图片链接
function mdGetImgSrc() {
	if(o.selectedImages && o.selectedImages.length) {
		return o.selectedImages[0];
	}
	return "";
}
