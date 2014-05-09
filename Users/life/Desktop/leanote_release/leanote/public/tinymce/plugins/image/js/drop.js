$(function() {
	var dropbox = $('#dropbox'),
		message = $('.message', dropbox);
	
	dropbox.filedrop({
		// The name of the $_FILES entry:
		paramname: 'file',
		
		maxfiles: 1,
    	maxfilesize: 3,
    	url: '/file/uploadImageJson',
		
		uploadFinished: function(i, file, re) {
			jbImagesDialog.uploadFinish({
				filename: re.Id,
				result: re.Msg,
				resultCode: re.Code + ''
			});
			
			$.data(file).addClass('done');
		},
		
    	error: function(err, file) {
			switch(err) {
				case 'BrowserNotSupported':
					showMessage('您的浏览器不支持拖拉上传文件, 请使用chrome或firefox浏览器');
					break;
				case 'TooManyFiles':
					alert('文件太多了');
					break;
				case 'FileTooLarge':
					alert(file.name+'文件太大, 最大支持3M图片');
					break;
				default:
					alert("出错了");
					break;
			}
		},
		
		// Called before each upload is started
		beforeEach: function(file){
			if(!file.type.match(/^image\//)){
				alert('只能上传图片!');
				return false;
			}
		},
		
		uploadStarted:function(i, file, len){
			//createImage(file);
		},
		
		progressUpdated: function(i, file, progress) {
			$('.progress-bar').css("width", progress + "%");
		}
    	 
	});
	
	var template = '<div class="preview">'+
						'<span class="imageHolder">'+
							'<img />'+
							'<span class="uploaded"></span>'+
						'</span>'+
						'<div class="progressHolder">'+
							'<div class="progress"></div>'+
						'</div>'+
					'</div>'; 
	
	
	function createImage(file){

		var preview = $(template), 
			image = $('img', preview);
			
		var reader = new FileReader();
		
		image.width = 100;
		image.height = 100;
		
		reader.onload = function(e){
			
			// e.target.result holds the DataURL which
			// can be used as a source of the image:
			
			image.attr('src',e.target.result);
			image.css("max-width", "200px")
		};
		
		// Reading the file as a DataURL. When finished,
		// this will trigger the onload function above:
		reader.readAsDataURL(file);
		
		message.hide();
		preview.appendTo(dropbox);
		
		// Associating a preview container
		// with the file, using jQuery's $.data():
		
		$.data(file,preview);
	}

	function showMessage(msg){
		message.html(msg);
	}

});