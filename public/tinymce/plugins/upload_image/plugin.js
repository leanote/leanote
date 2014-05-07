tinymce.PluginManager.add('upload_image', function(editor) {
	var self = this;

    function genNav() {
    	var $con = $(editor.getContent({format: 'raw'}));
      // 构造一棵树
      // {"h1-title":{h2-title:{}}}
      var tree = [];//[{title: "xx", children:[{}]}, {title:"xx2"}];
      var hs = $con.find("h1,h2,h3,h4,h5,h6").toArray();
      var titles = '<ul>';
      for(var i = 0; i < hs.length; ++i) {
        var text = $(hs[i]).text(); 
        var tagName = hs[i].tagName.toLowerCase();
        titles += '<li class="nav-' + tagName + '"><a class="scrollTo-a" onclick="scrollTo(\'' + tagName + '\', \'' + text + '\')">' + text + '</a></li>';
      }
      titles += "</ul>";

      $("#leanoteNav").html(titles);

      resizeIfrmHeight();
    }

    var lastResizeTime = null;

    function resizeIfrmHeight() {
      return;
      // if(!lastResizeTime) {
      //   lastResizeTime = new Date().getTime();
      // }
      // var n = new Date().getTime();
      // if(n - lastResizeTime < 5000) {
      //   return;
      // }
      // lastResizeTime = n;
      log("resize....")
      $("#content_ifr").parent().css("position", "relative");
      // 在这里设置iframe高度
      // 之前这段在上面, 每次输入都会重新计算高度, 调整iframe高度
      var height = $("#content_ifr").contents().find("html").height(); // 之前是body的高度, 不出入
      $("#content_ifr").height(height);
      $("#content_ifr").css("border", "1px solid #ccc")
      // $("#noteTitleDiv").hide();
      var t = $("#content_ifr").parent().scrollTop();
      // $("#content_ifr").parent().scrollTop(t + 30)
      
      var rng = editor.selection.getRng(true).cloneRange()
      var eOffsetTop = rng.endContainer.offsetTop;
      // var eOffsetTop = rng.startContainer.offsetTop;
      log(rng);
      log("firstElementChild: "+ $(rng.endContainer.firstElementChild).offset().top);
      if(!eOffsetTop) {
        var html = $(rng.endContainer.parentNode).html();
        eOffsetTop = $(rng.endContainer.parentNode).offset().top

      }
      
      var vStart = $("#content_ifr").parent().scrollTop();
      var vEnd = $("#content_ifr").parent().height() + vStart;
      log("vStart: " + vStart + "; vEnd" + vEnd);

      var realEoffsetTop = eOffsetTop; // + $("#noteTitle").height();    
      log("realEoffsetTop " + realEoffsetTop);  
      // 如果不在可视范围内
      // vStart realEoffsetTop vEnd
      if(vStart < realEoffsetTop && realEoffsetTop < vEnd) {
      } else {
        var top = realEoffsetTop - vEnd + 30;
        log("resize.............." + top)
        // $("#content_ifr").parent().scrollTop(top);
      }
    }


    editor.on('init', function() {
      // resizeIfrmHeight();
      // $("#content_ifr").attr("scrolling", "no");

     	window.setTimeout(function() {
  			// 最开始时加载
  			genNav();
  			editor.on('setcontent beforeadd undo paste', genNav);
  			editor.on('keyup', function(e) {
  				// if (e.keyCode == 32) {
  					// genNav();
            // setTimeout(resizeIfrmHeight, 1000);
  				// }
  			});
        editor.on('keydown', function(e) {
          // if (e.keyCode == 32) {
            genNav();
            // setTimeout(resizeIfrmHeight, 1000);
          // }
        });
        // 为了把下拉菜单关闭
        editor.on("click", function(e) {
          $("body").trigger("click");
          // log($("#content_ifr").contents().find("body").height());
        });

  		}, 0);
	});
});