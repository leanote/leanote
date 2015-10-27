LEA.cmroot = 1;
(function($) {
    function returnfalse() { return false; };
    $.fn.contextmenu = function(option) {
    	var cmroot = "contextmenu" + LEA.cmroot; //Math.floor((Math.random() + Math.random()) * 10000);
    	LEA.cmroot++;
        option = $.extend({ alias: cmroot, width: 150 }, option);
        var ruleName = null, target = null,
	    groups = {}, mitems = {}, actions = {}, showGroups = [],
        itemTpl = '<div class="b-m-$[type]" unselectable="on"><div class="clearfix cm-item"><div class="b-m-icon pull-left"><i class="fa $[faIcon]"></i>$[imgIcon]</div><div class="pull-left cm-text"><span class="c-text" unselectable="on">$[text]</span></div></div></div>';
		itemNoIconTpl = "<div class='b-m-$[type]' unselectable=on><nobr unselectable=on><span align='absmiddle'></span><span class='c-text' unselectable=on>$[text]</span></nobr></div>";
        var gTemplet = $("<div/>").addClass("b-m-mpanel").attr("unselectable", "on").css("display", "none");
        var iTemplet = $("<div/>").addClass("b-m-item").attr("unselectable", "on");
        var sTemplet = $("<div/>").addClass("b-m-split");
        var $body = $("body");
        
        var itemsCache = {}; // idx => items
        
        //build group item, which has sub items
        var buildGroup = function(obj) { // this = $("")对象, obj=item
            groups[obj.alias] = this;
            this.gidx = obj.alias;
            this.id = obj.alias;
            if (obj.disable) {
                this.disable = obj.disable;
                this.className = "b-m-idisable";
            }
            $(this).width(obj.width).click(function(){}).mousedown(returnfalse).appendTo($body);
            
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
            var len = items.length;
            for (var i = 0; i < len; i++) {
            	var item = items[i];
                if (item.type == "splitLine") {
                    tmp = sTemplet.clone()[0];
                } else {
                	// life, alias可以不需要, 从text取, 但必须唯一
                	if(!item.alias) {
                		if(parentAlias) {
                			item.alias = parentAlias + "." + item.text; // 移动.Hadoop 
                		} else {
                			item.alias = item.text;
                		}
                		// log(item.alias);
                	}
                    item.gidx = gidx;
                    if (item.type == "group" && !item.action) {
                        //group 
                        buildGroup.apply(gTemplet.clone()[0], [item]);
                        itemsCache[item.alias] = item.items;
                        // 递归调用, 可以动态生成?
                        // arguments.callee(item.alias, item.items, item.alias); // life 传上级的alias, 避免重复
                        item.type = "arrow";
                        tmp = buildItem.apply(iTemplet.clone()[0], [item]);
                    } else {
                    	// 如果group有action还是可以点击的 life
                    	if(item.type == "group") {
	                    	//group 
	                        buildGroup.apply(gTemplet.clone()[0], [item]);
	                        itemsCache[item.alias] = item.items;
	                        // 递归调用
	                        // arguments.callee(item.alias, item.items, item.alias); // life 传上级的alias, 避免重复
	                        item.type = "arrow";
	                        tmp = buildItem.apply(iTemplet.clone()[0], [item]);
                    	} else {
	                        //normal item
	                        item.type = "ibody";
	                        tmp = buildItem.apply(iTemplet.clone()[0], [item]);
                        }
						// 
                        var thisItem = item;
                        
                        // 点击item
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
                tmp = item = item.items = null;
            } //end for
            gidx = items = null;
        };
        // hover
        var overItem = function(e) {
            //menu item is disabled          
            if (this.disable)
                return false;
            hideMenuPane.call(groups[this.gidx]);
            //has sub items
            if (this.group) {
                var pos = $(this).offset();
                var width = $(this).outerWidth();
                showMenuGroup.apply(groups[this.idx], [pos, width, this]);
            }
            this.className = "b-m-ifocus";
            return false;
        };
        // hover out
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
        
        // 显示group, 这里可以动态生成
        // show menu group at specified position
        var showMenuGroup = function(pos, width, t) {
        	var $this = $(this); // dom 对象
        	// 没有东西, 那么生成之, 动态生成 life [ok]
        	if($this.html() == "") {
        		addItems(t.idx, itemsCache[t.idx], t.idx);
        	}
            var bwidth = $body.width();
            // var bheight = $body.height();
            var bheight = document.documentElement.clientHeight-10;
            bheight = bheight < 0 ? 100 : bheight;
            var mwidth = $(this).outerWidth();
            var mheight = $(this).outerHeight()-10;
            mheight = mheight < 0 ? 100 : mheight;
            var mwidth = $(this).outerWidth();
            pos.left = (pos.left + width + mwidth > bwidth) ? (pos.left - mwidth < 0 ? 0 : pos.left - mwidth) : pos.left + width;
            pos.top = (pos.top + mheight > bheight) ? (pos.top - mheight + (width > 0 ? 25 : 0) < 0 ? 0 : pos.top - mheight + (width > 0 ? 25 : 0)) : pos.top;
            $(this).css(pos).show().css("max-height", bheight);

            showGroups.push(this.gidx);
        };

        //to hide menu
        var hideMenuPane = function() {
            var alias = null;

            // console.log('showGroups: ' + showGroups.length)
            for (var i = showGroups.length - 1; i >= 0; i--) {
                if (showGroups[i] == this.gidx)
                    break;
                alias = showGroups.pop();
                groups[alias].style.display = "none";
                mitems[alias] && (mitems[alias].className = "b-m-item");
            }
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
            // 先隐藏之前的
            hideMenuPane();
            removeContextmenuClass();

            target = menutarget;
            showMenuGroup.call(groups[cmroot], { left: e.pageX, top: e.pageY }, 0);

            // 在该target上添加contextmenu-hover
            if(target && !$(target).hasClass("item-active")) {
                $(target).addClass("contextmenu-hover");
            }
        }
        
        // 初始化
        var $root = $("#" + option.alias);
        var root = null;
        if ($root.length == 0) {
            root = buildGroup.apply(gTemplet.clone()[0], [option]);
            root.applyrule = applyRule;
            root.showMenu = showMenu;
        	// 这里很费时
            addItems(option.alias, option.items);
        }
        else {
            root = $root[0];
        }
        
        function onShowMenu(e) {
            var bShowContext = (option.onContextMenu && $.isFunction(option.onContextMenu)) ? option.onContextMenu.call(this, e) : true;
            if (bShowContext) {
                if (option.onShow && $.isFunction(option.onShow)) {
                    option.onShow.call(this, root);
                }
                root.showMenu(e, this);
            }
            // 阻止冒泡, 默认事件
            if(e) {
	            e.preventDefault();
            }
            return false;
        }
        // bind event
        var me = $(option.parent).on('contextmenu', option.children, function(e){ 
        	onShowMenu.call(this, e);
        });

        function removeContextmenuClass() {
            Note.$itemList.find('li').removeClass('contextmenu-hover');
        }

        // life , 之前是document, 绑定document即使stopPro也会执行到这里
        $('body').on('click', function(e) {
            removeContextmenuClass();
            hideMenuPane();
        });

        //to apply rule
        if (option.rule) {
            applyRule(option.rule);
        }
        /*
        gTemplet = iTemplet = sTemplet = itemTpl = buildGroup = buildItem = null;
        addItems = overItem = outItem = null;
        */
        //CollectGarbage();
        
        var out = {
        	destroy: function() {
        		me.unbind("contextmenu");
        	},
        	showMenu: function(e, target) {
        		onShowMenu.call(target, e);
        	}
        }
        return out;
    }
})(jQuery);