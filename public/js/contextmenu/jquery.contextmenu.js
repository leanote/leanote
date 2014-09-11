LEA.cmroot = 1;
(function($) {
    function returnfalse() { return false; };
    $.fn.contextmenu = function(option) {
    	var cmroot = "contextmenu" + LEA.cmroot; //Math.floor((Math.random() + Math.random()) * 10000);
    	LEA.cmroot++;
        option = $.extend({ alias: cmroot, width: 150 }, option);
        var ruleName = null, target = null,
	    groups = {}, mitems = {}, actions = {}, showGroups = [],
        itemTpl = '<div class="b-m-$[type]" unselectable="on"><div class="clearfix"><div class="b-m-icon pull-left"><i class="fa $[faIcon]"></i>$[imgIcon]</div><div class="pull-left"><span class="c-text" unselectable="on">$[text]</span></div></div></div>';
		itemNoIconTpl = "<div class='b-m-$[type]' unselectable=on><nobr unselectable=on><span align='absmiddle'></span><span class='c-text' unselectable=on>$[text]</span></nobr></div>";
        var gTemplet = $("<div/>").addClass("b-m-mpanel").attr("unselectable", "on").css("display", "none");
        var iTemplet = $("<div/>").addClass("b-m-item").attr("unselectable", "on");
        var sTemplet = $("<div/>").addClass("b-m-split");
        //build group item, which has sub items
        var buildGroup = function(obj) {
            groups[obj.alias] = this;
            this.gidx = obj.alias;
            this.id = obj.alias;
            if (obj.disable) {
                this.disable = obj.disable;
                this.className = "b-m-idisable";
            }
            $(this).width(obj.width).click(returnfalse).mousedown(returnfalse).appendTo($("body"));
            obj = null;
            return this;
        };
        var buildItem = function(obj) {
            var T = this;
            T.title = obj.text;
            T.idx = obj.alias;
            T.gidx = obj.gidx;
            T.data = obj;
            var imgIcon = "";
            if(obj.icon) {
            	imgIcon = '<img src="' + obj.icon + '"/>';
            }
            obj.imgIcon = imgIcon;
			if(obj.icon) {
				T.innerHTML = itemTpl.replace(/\$\[([^\]]+)\]/g, function() {
					return obj[arguments[1]];
				});
			} else {
				T.innerHTML = itemTpl.replace(/\$\[([^\]]+)\]/g, function() {
					return obj[arguments[1]];
				});
			}
            if (obj.disable) {
                T.disable = obj.disable;
                T.className = "b-m-idisable";
            }
            obj.items && (T.group = true);
            obj.action && (actions[obj.alias] = obj.action);
            mitems[obj.alias] = T;
            T = obj = null;
            return this;
        };
        //add new items
        var addItems = function(gidx, items, parentAlias) {
            var tmp = null;
            for (var i = 0; i < items.length; i++) {
                if (items[i].type == "splitLine") {
                    //split line
                    tmp = sTemplet.clone()[0];
                } else {
                	// life, alias可以不需要, 从text取, 但必须唯一
                	if(!items[i].alias) {
                		if(parentAlias) {
                			items[i].alias = parentAlias + "." + items[i].text; // 移动.Hadoop 
                		} else {
                			items[i].alias = items[i].text;
                		}
                		// log(items[i].alias);
                	}
                    items[i].gidx = gidx;
                    if (items[i].type == "group") {
                        //group 
                        buildGroup.apply(gTemplet.clone()[0], [items[i]]);
                        arguments.callee(items[i].alias, items[i].items, items[i].alias); // life 传上级的alias, 避免重复
                        items[i].type = "arrow";
                        tmp = buildItem.apply(iTemplet.clone()[0], [items[i]]);
                    } else {
                        //normal item
                        items[i].type = "ibody";
                        tmp = buildItem.apply(iTemplet.clone()[0], [items[i]]);
						// 
                        var thisItem = items[i];
                        
                        // 用闭包来存储变量
                        (function(thisItem, tmp) {
	                        $(tmp).click(function(e) {
	                            if (!this.disable) {
									// console.log(target);
	                            	// 调用...
	                                if ($.isFunction(actions[this.idx])) {
	                                    actions[this.idx].call(this, target, thisItem);
	                                }
	                                hideMenuPane();
	                                
					            	// life
						            $(target).removeClass("contextmenu-hover");
	                            }
	                            return false;
	                        });
	                        	
                        }(thisItem, tmp));
                        

                    } //end if
                    $(tmp).bind("contextmenu", returnfalse).hover(overItem, outItem);
                }
                groups[gidx].appendChild(tmp);
                tmp = items[i] = items[i].items = null;
            } //end for
            gidx = items = null;
        };
        var overItem = function(e) {
            //menu item is disabled          
            if (this.disable)
                return false;
            hideMenuPane.call(groups[this.gidx]);
            //has sub items
            if (this.group) {
                var pos = $(this).offset();
                var width = $(this).outerWidth();
                showMenuGroup.apply(groups[this.idx], [pos, width]);
            }
            this.className = "b-m-ifocus";
            return false;
        };
        //menu loses focus
        var outItem = function(e) {
            //disabled item
            if (this.disable )
                return false;
            if (!this.group) {
                //normal item
                this.className = "b-m-item";
            } //Endif
            return false;
        };
        //show menu group at specified position
        var showMenuGroup = function(pos, width) {
            var bwidth = $("body").width();
            var bheight = document.documentElement.clientHeight;
            var mwidth = $(this).outerWidth();
            var mheight = $(this).outerHeight();
            pos.left = (pos.left + width + mwidth > bwidth) ? (pos.left - mwidth < 0 ? 0 : pos.left - mwidth) : pos.left + width;
            pos.top = (pos.top + mheight > bheight) ? (pos.top - mheight + (width > 0 ? 25 : 0) < 0 ? 0 : pos.top - mheight + (width > 0 ? 25 : 0)) : pos.top;
            $(this).css(pos).show();
            showGroups.push(this.gidx);
        };
        //to hide menu
        var hideMenuPane = function() {
            var alias = null;
            for (var i = showGroups.length - 1; i >= 0; i--) {
                if (showGroups[i] == this.gidx)
                    break;
                alias = showGroups.pop();
                groups[alias].style.display = "none";
                mitems[alias] && (mitems[alias].className = "b-m-item");
                
            } //Endfor
            //CollectGarbage();
        };
        function applyRule(rule) {
        	/*
            if (ruleName && ruleName == rule.name)
                return false;
            */
            for (var i in mitems)
                disable(i, !rule.disable);
            for (var i = 0; i < rule.items.length; i++)
                disable(rule.items[i], rule.disable);
            ruleName = rule.name;
        };
        function disable(alias, disabled) {
            var item = mitems[alias];
            if(!item || !item.lastChild) {
            	return;
            }
            item.className = (item.disable = item.lastChild.disabled = disabled) ? "b-m-idisable" : "b-m-item";
        };

        /* to show menu  */
        function showMenu(e, menutarget) {
            target = menutarget;
            showMenuGroup.call(groups[cmroot], { left: e.pageX, top: e.pageY }, 0);
            // life
            if(!$(target).hasClass("item-active")) {
	            $(target).addClass("contextmenu-hover");
            }
            
            $(document).one('mousedown', function() {
            	hideMenuPane();
            	// life
	            $(target).removeClass("contextmenu-hover");
            })
        }
        
        // 初始化
        var $root = $("#" + option.alias);
        var root = null;
        if ($root.length == 0) {
            root = buildGroup.apply(gTemplet.clone()[0], [option]);
            root.applyrule = applyRule;
            root.showMenu = showMenu;
            addItems(option.alias, option.items);
        }
        else {
            root = $root[0];
        }
        
        /*
        var me = $(this).each(function() {
            return $(this).bind('contextmenu', function(e) {
                var bShowContext = (option.onContextMenu && $.isFunction(option.onContextMenu)) ? option.onContextMenu.call(this, e) : true;
                if (bShowContext) {
                    if (option.onShow && $.isFunction(option.onShow)) {
                        option.onShow.call(this, root);
                    }
                    root.showMenu(e, this);
                }
                return false;
            });
        });
        */
        
        // bind event
        var me = $(option.parent).on('contextmenu', option.children, function(e) {
            var bShowContext = (option.onContextMenu && $.isFunction(option.onContextMenu)) ? option.onContextMenu.call(this, e) : true;
            if (bShowContext) {
                if (option.onShow && $.isFunction(option.onShow)) {
                    option.onShow.call(this, root);
                }
                root.showMenu(e, this);
            }
            // 阻止冒泡, 默认事件
            e.preventDefault();
            return false;
        });
        
        //to apply rule
        if (option.rule) {
            applyRule(option.rule);
        }
        gTemplet = iTemplet = sTemplet = itemTpl = buildGroup = buildItem = null;
        addItems = overItem = outItem = null;
        //CollectGarbage();
        
        return me;
    }
})(jQuery);