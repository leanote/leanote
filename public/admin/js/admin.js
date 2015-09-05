function log(o) {
	console.log(o);
}

var ARTDIALOG = {stack:[], id: 1};
ARTDIALOG.defaultConfig = {title:"", draggable: false, padding: 0, fixed: false, lock: false, opacity: 0.3};
function openDialog(config) {
	config = config || {};
	if(!config.id) {
		config.id = ARTDIALOG.id++;
	}
	if(config.content) {
		// $("#life")
		if(typeof config.content == "object") {
			try {
				config.content = config.content.get(0);
			} catch(e) {
				// 不是jquery对象, 而是dom对象, document.getElementById("xx")
				// config.content = config.content;
			}
		}
	}
	config = $.extend({}, ARTDIALOG.defaultConfig, config);
	// var content = "<div id='sys-ibank'>" + $("#sys-ibank").html() + "</div>";
	var d = art.dialog(config);
	
	if(config.url) {
		$.get(config.url, {}, function(ret) {
			d.content(ret);
		});
	}
	
	ARTDIALOG.stack.push(config.id);
	return d;
	// $( "#sys-ibank" ).dialog({title:"插入图片", width: 800,draggable: false});
}
function closeDialog(){
	var list = art.dialog.list;
	if(!list) {
		return;
	}
	while(true) {
		var d = ARTDIALOG.stack.pop();
		if(d) {
			if(list[d]) {
				list[d].close();
				return;
			}
		} else {
			return;
		}
	}
}
// 关闭最近内容为loading的dialog
function closeLatestLoadingDialog() {
	var list = art.dialog.list;
	if(!list) {
		return;
	}
	while(true) {
		var d = ARTDIALOG.stack.pop();
		if(d) {
			if(list[d]) {
				var dd = list[d];
				if($(dd.content()).text() == "loading..") {
					dd.close();
				}
				return;
			}
		} else {
			return;
		}
	}
}

if(typeof art != "undefined") {
	/**
	 * 警告
	 * @param	{String}	消息内容
	 */
	art.alert = function (content, callback) {
	    return artDialog({
	        id: 'Alert',
	        icon: 'warning',
	        fixed: true,
	        lock: true,
	        content: content,
	        ok: true,
	        opacity: 0.3,
	        close: callback
	    });
	};
	
	/**
	 * 确认
	 * @param	{String}	消息内容
	 * @param	{Function}	确定按钮回调函数
	 * @param	{Function}	取消按钮回调函数
	 */
	art.confirm = function (content, yes, no) {
	    return artDialog({
	        id: 'Confirm',
	        icon: 'question',
	        fixed: true,
	        lock: true,
	        opacity: .3,
	        content: content,
	        ok: function (here) {
	            return yes.call(this, here);
	        },
	        cancel: function (here) {
	            return no && no.call(this, here);
	        }
	    });
	};
	
	
	/**
	 * 提问
	 * @param	{String}	提问内容
	 * @param	{Function}	回调函数. 接收参数：输入值
	 * @param	{String}	默认值
	 */
	art.prompt = function (content, yes, value) {
	    value = value || '';
	    var input;
	    
	    return artDialog({
	        id: 'Prompt',
	        icon: 'question',
	        fixed: true,
	        lock: true,
	        opacity: .3,
	        content: [
	            '<div style="margin-bottom:5px;font-size:12px">',
	                content,
	            '</div>',
	            '<div>',
	                '<input value="',
	                    value,
	                '" style="width:18em;padding:6px 4px" />',
	            '</div>'
	            ].join(''),
	        init: function () {
	            input = this.DOM.content.find('input')[0];
	            input.select();
	            input.focus();
	        },
	        ok: function (here) {
	            return yes && yes.call(this, input.value, here);
	        },
	        cancel: true
	    });
	};
	
	
	/**
	 * 短暂提示
	 * @param	{String}	提示内容
	 * @param	{Number}	显示时间 (默认1.5秒)
	 */
	art.tips = function (content, time) {
	    return artDialog({
	        id: 'Tips',
	        title: false,
	        cancel: false,
	        fixed: true,
	        top: 0,
	        // lock: true,
	        opacity: 0.3
	    })
	    .content('<div style="padding: 0 1em;">' + content + '</div>')
	    .time(time || 1);
	};
	
	// art dialog bind 
	// <a href="javascript:;" id="agree_btn" class="button art-dialog" data-url="index.php?app=seller_refund&amp;action=confirm_refund&amp;order_id=55" data-title="确认退款">同意退款</a>
	$(function() {
		$(".art-dialog").click(function(){
			var title = $(this).data('title');
			var url = $(this).data("url");
			var lock = +$(this).data('lock');
			var width = $(this).data('width');
			var config = {url: url, title: title, lock: lock};
			if(width) {
				config.width = width;
			}
			openDialog(config);
		});
	});
}

// 删除确认
function drop_confirm(msg, url) {
	if(art) {
		art.confirm(msg,function(){
	        // window.location = url;
	        var self = this;
	        self.content("正在处理...");
	        ajaxGet(url, {}, function(ret) {
	        	if(ret.done) {
	        		self.content("操作成功, 正在刷新...");
	        		location.reload();
	        	} else {
	        		art.alert(ret.msg);
	        	}
        		self.close();
	        });
	        return false;
		});
	} else {
	    if(confirm(msg)){
	        window.location = url;
	    }
    }
}
function init_validator(target, rules, messages) {
	var config = {
	        errorElement : 'div',
	        errorClass : 'help-block alert alert-warning',
	        focusInvalid : false,
	        ignore: ".ignore",
	        highlight : function(element) {
	        	var $p = $(element).closest('.control-group');
	            $p.removeClass("success").addClass('error');
	        },
	        success : function(label) {
	        	var $p = label.closest('.control-group');
	            $p.removeClass('error');
	            $p.addClass("success");
	            
	            $p.find(".help-block").hide();
	            $(label).hide();
	        },
	        errorPlacement : function(error, element) {
	        	var $p = element.parent('div');
	            element.parent('div').append(error);
	            log(element);
	            log($p);
	        },
	        submitHandler : function(form) {
	            form.submit();
	        }
	    };
	if(rules) {
		config.rules = rules;
	}
	if(messages) {
		config.messages = messages;
	}
	return $(target).validate(config);
}


function enter_submit(btnId) {
	var theEvent = window.event || arguments.callee.caller.arguments[0];
	if(theEvent.keyCode == 13||theEvent.keyCode == 108) {
		$(btnId).trigger('click');
	}
}

// send email dialog
function openSendEmailDialog(emails) {
	openDialog({width: 500,  url: "/adminEmail/sendEmailDialog?emails=" + emails, title: "Send Email"});
}

function goNowToDatetime(goNow) {
	if(!goNow) {
		return "";
	}
	return goNow.substr(0, 10) + " " + goNow.substr(11, 8);
}

!function ($) {
  $(function(){
  	
  	//navigation
	!function ($) {
		/**
		* description
		* every href must begin with "?t="
		*/
		//initial
		function getParameterByName(name, queryString) {
		    name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]"); //address array [] condiion
		    var regex = new RegExp("[\\?&]" + name + "=([^&#]*)"),
		        results = regex.exec(queryString);
		    return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
		}

		var pathname = window.location.pathname; // admin/t
		var search = window.location.search; // ?t=xxx
		var fullPath = pathname;

		//test case
		//http://localhost:9000/admin/t?p=0&t=email/sendToUsers
		//t=email/sendToUsers  and t=email/send
		var paramId = getParameterByName("t",window.location.search)
		if(paramId !== "") {
			var fullPath = pathname + "?t=" + paramId; // /admin/t?t=xxx
		}
		
		$("#nav > li").removeClass("active");

		var $thisLi = $('#nav a[href="' + fullPath + '"]').parent();
		// 自己
		$thisLi.addClass("active");
		// 父也active
		$thisLi.parent().parent().addClass('active');

		// event binding
	  	$(".nav li > a").click(function(e) {
	  		$p = $(this).closest("ul");
	  		var $li = $(this).closest("li");
	  		if($li.find("ul").length == 0) {
	  			return true;
	  		}
	  		e.preventDefault();
	  		var hasClass = $li.hasClass("active");
	  		$p.find("li").removeClass("active");
	  		if(hasClass) {
	  		} else {
	  			$li.addClass("active");
	  		}
	  	});
	}($);
  	
  	// sort
  	$(".th-sortable").click(function() {
  		var up = $(this).hasClass("th-sort-up");
  		var down = $(this).hasClass("th-sort-down");
  		var url = $(this).data("url");
  		var sorter = $(this).data("sorter");
  		var t = "th-sort-up";
  		if(up) {
  			t = "th-sort-down";
  			var sUrl = "sorter=" + sorter + "-down";
  		} else {
  			var sUrl = "sorter=" + sorter + "-up";
  		}
  		
  		if(url.indexOf("?") > 0) {
  			location.href = url + "&" + sUrl;
  		} else {
  			location.href = url + "?" + sUrl;
  		}
  		$(this).removeClass("th-sort-up th-sort-down").addClass(t);
  	});
  	
  	// search 
  	$(".search-group input").keyup(function(e){
  		enter_submit(".search-group button");
  	});
  	$(".search-group button").click(function(e){
  		var url = $(this).data("url");
  		$input = $(this).closest(".search-group").find("input");
  		var keywords = $input.val();
  		/*
  		if(!keywords) {
  			$input.focus();
  			return;
  		}
  		*/
  		if(url.indexOf("?") > 0) {
	  		location.href = url + "&keywords=" + keywords;
  		} else {
	  		location.href = url + "?keywords=" + keywords;
  		}
  	});
  	
  	
  	//--------------------------
 	
	// sparkline
	var sr, sparkline = function($re){
		$(".sparkline").each(function(){
			var $data = $(this).data();
			if($re && !$data.resize) return;
			($data.type == 'pie') && $data.sliceColors && ($data.sliceColors = eval($data.sliceColors));
			($data.type == 'bar') && $data.stackedBarColor && ($data.stackedBarColor = eval($data.stackedBarColor));
			$data.valueSpots = {'0:': $data.spotColor};
			$(this).sparkline('html', $data);
		});
	};
	$(window).resize(function(e) {
		clearTimeout(sr);
		sr = setTimeout(function(){sparkline(true)}, 500);
	});
	sparkline(false);


	// easypie
    $('.easypiechart').each(function(){
    	var $this = $(this), 
    	$data = $this.data(), 
    	$step = $this.find('.step'), 
    	$target_value = parseInt($($data.target).text()),
    	$value = 0;
    	$data.barColor || ( $data.barColor = function($percent) {
            $percent /= 100;
            return "rgb(" + Math.round(200 * $percent) + ", 200, " + Math.round(200 * (1 - $percent)) + ")";
        });
    	$data.onStep =  function(value){
    		$value = value;
    		$step.text(parseInt(value));
    		$data.target && $($data.target).text(parseInt(value) + $target_value);
    	}
    	$data.onStop =  function(){
    		$target_value = parseInt($($data.target).text());
    		$data.update && setTimeout(function() {
		        $this.data('easyPieChart').update(100 - $value);
		    }, $data.update);
    	}
		$(this).easyPieChart($data);
	});

  	// combodate
	$(".combodate").each(function(){ 
		$(this).combodate();
		$(this).next('.combodate').find('select').addClass('form-control');
	});

	// datepicker
	$(".datepicker-input").each(function(){ $(this).datepicker();});

	// dropfile
	$('.dropfile').each(function(){
		var $dropbox = $(this);
		if (typeof window.FileReader === 'undefined') {
		  $('small',this).html('File API & FileReader API not supported').addClass('text-danger');
		  return;
		}

		this.ondragover = function () {$dropbox.addClass('hover'); return false; };
		this.ondragend = function () {$dropbox.removeClass('hover'); return false; };
		this.ondrop = function (e) {
		  e.preventDefault();
		  $dropbox.removeClass('hover').html('');
		  var file = e.dataTransfer.files[0],
		      reader = new FileReader();
		  reader.onload = function (event) {
		  	$dropbox.append($('<img>').attr('src', event.target.result));
		  };
		  reader.readAsDataURL(file);
		  return false;
		};
	});

	// fuelux pillbox
	var addPill = function($input){
		var $text = $input.val(), $pills = $input.closest('.pillbox'), $repeat = false, $repeatPill;
		if($text == "") return;
		$("li", $pills).text(function(i,v){
	        if(v == $text){
	        	$repeatPill = $(this);
	        	$repeat = true;
	        }
	    });
	    if($repeat) {
	    	$repeatPill.fadeOut().fadeIn();
	    	return;
	    };
	    $item = $('<li class="label bg-dark">'+$text+'</li> ');
		$item.insertBefore($input);
		$input.val('');
		$pills.trigger('change', $item);
	};

	$('.pillbox input').on('blur', function() {
		addPill($(this));
	});

	$('.pillbox input').on('keypress', function(e) {
	    if(e.which == 13) {
	        e.preventDefault();
	        addPill($(this));
	    }
	});

	// slider
	$('.slider').each(function(){
		$(this).slider();
	});

	// wizard
  $(document).on('change', '.wizard', function (e, data) {
    if(data.direction !== 'next' ) return;
    var item = $(this).wizard('selectedItem');
    var $step = $(this).find('.step-pane:eq(' + (item.step-1) + ')');
    var validated = true;
    $('[data-required="true"]', $step).each(function(){
      return (validated = $(this).parsley( 'validate' ));
    });
    if(!validated) return e.preventDefault();
  });

	// sortable
	if ($.fn.sortable) {
	  $('.sortable').sortable();
	}

	// slim-scroll
	$('.no-touch .slim-scroll').each(function(){
		var $self = $(this), $data = $self.data(), $slimResize;
		$self.slimScroll($data);
		$(window).resize(function(e) {
			clearTimeout($slimResize);
			$slimResize = setTimeout(function(){$self.slimScroll($data);}, 500);
		});

    $(document).on('updateNav', function(){
      $self.slimScroll($data);
    });
	});

	// pjax
	if ($.support.pjax) {
	  $(document).on('click', 'a[data-pjax]', function(event) {
	  	event.preventDefault();
	    var container = $($(this).data('target'));
	    $.pjax.click(event, {container: container});
	  })
	};

	// portlet
	$('.portlet').each(function(){
		$(".portlet").sortable({
	        connectWith: '.portlet',
            iframeFix: false,
            items: '.portlet-item',
            opacity: 0.8,
            helper: 'original',
            revert: true,
            forceHelperSize: true,
            placeholder: 'sortable-box-placeholder round-all',
            forcePlaceholderSize: true,
            tolerance: 'pointer'
	    });
    });

	// docs
    $('#docs pre code').each(function(){
	    var $this = $(this);
	    var t = $this.html();
	    $this.html(t.replace(/</g, '&lt;').replace(/>/g, '&gt;'));
	});

	// fontawesome
	$(document).on('click', '.fontawesome-icon-list a', function(e){
		e && e.preventDefault();
	});

	// table select/deselect all
	$(document).on('change', 'table thead [type="checkbox"]', function(e){
		e && e.preventDefault();
		var $table = $(e.target).closest('table'), $checked = $(e.target).is(':checked');
		$('tbody [type="checkbox"]',$table).prop('checked', $checked);
	});

	// random progress
	$(document).on('click', '[data-toggle^="progress"]', function(e){
		e && e.preventDefault();

		$el = $(e.target);
		$target = $($el.data('target'));
		$('.progress', $target).each(
			function(){
				var $max = 50, $data, $ps = $('.progress-bar',this).last();
				($(this).hasClass('progress-xs') || $(this).hasClass('progress-sm')) && ($max = 100);
				$data = Math.floor(Math.random()*$max)+'%';
				$ps.css('width', $data).attr('data-original-title', $data);
			}
		);
	});
	
	// add notes
	function addMsg($msg){
		var $el = $('.nav-user'), $n = $('.count:first', $el), $v = parseInt($n.text());
		$('.count', $el).fadeOut().fadeIn().text($v+1);
		$($msg).hide().prependTo($el.find('.list-group')).slideDown().css('display','block');
	}
	var $msg = '<a href="#" class="media list-group-item">'+
                  '<span class="pull-left thumb-sm text-center">'+
                    '<i class="fa fa-envelope-o fa-2x text-success"></i>'+
                  '</span>'+
                  '<span class="media-body block m-b-none">'+
                    'Sophi sent you a email<br>'+
                    '<small class="text-muted">1 minutes ago</small>'+
                  '</span>'+
                '</a>';	
  setTimeout(function(){addMsg($msg);}, 1500);

	// select2 
 	if ($.fn.select2) {
      $("#select2-option").select2();
      $("#select2-tags").select2({
        tags:["red", "green", "blue"],
        tokenSeparators: [",", " "]}
      );
  	}


  });
}(window.jQuery);
