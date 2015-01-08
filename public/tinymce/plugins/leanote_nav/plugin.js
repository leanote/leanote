/**
 * leanote nav
 * 
 */
tinymce.PluginManager.add('leanote_nav', function(editor) {
	var self = this;
	var preBody = "";
    function genNav() {
    	var body = editor.getBody();
    	var $con = $(body);
    	var html = $con.html();
    	if(preBody == html) {
    		return;
    	}
    	preBody = html;
		// 构造一棵树
		// {"h1-title":{h2-title:{}}}
		var tree = [];//[{title: "xx", children:[{}]}, {title:"xx2"}];
		var hs = $con.find("h1,h2,h3,h4,h5,h6").toArray();
		var titles = '<ul>';
		for(var i = 0; i < hs.length; ++i) {
			var text = $(hs[i]).text(); 
			var tagName = hs[i].tagName.toLowerCase();
			// scrollTo在page.js中定义
			titles += '<li class="nav-' + tagName + '"><a data-a="' + tagName + '-' + encodeURI(text)+'" onclick="scrollTo(this, \'' + tagName + '\', \'' + text + '\')">' + text + '</a></li>';
		}
		titles += "</ul>";
		$("#leanoteNavContent").html(titles).height("auto"); // auto
		if(!hs.length) {
			$("#leanoteNavContent").html("&nbsp; Nothing...");
		}
		
		// 这里, resize Height
		var curH = $("#leanoteNavContent").height();
		var pH = $("#editorContent").height()-29;
		if(curH > pH) {	
			$("#leanoteNavContent").height(pH);
		}
	}

    var lastResizeTime = null;
    editor.on('init', function() {
     	window.setTimeout(function() {
  			// 最开始时加载
  			genNav();
  			editor.on('setcontent beforeadd undo paste', genNav);
  			
  			// 这里, 以前是keydown, 太频繁
	        editor.on('ExecCommand', function(e) {
	             genNav();
	        });
	        
	        // 为了把下拉菜单关闭
	        editor.on("click", function(e) {
	        	genNav();
	        	$("body").trigger("click");
	        });

  		}, 0);
     	
	});
});