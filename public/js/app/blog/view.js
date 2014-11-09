function scrollTo(self, tagName, text) {
	var iframe = $("#content");
	var target = iframe.find(tagName + ":contains(" + text + ")");
	
	// 找到是第几个
	// 在nav是第几个
	var navs = $('#blogNavContent [data-a="' + tagName + '-' + encodeURI(text) + '"]');
	var len = navs.size();
	for(var i = 0; i < len; ++i) {
		if(navs[i] == self) {
			break;
		}
	}
	
	if (target.size() >= i+1) {
		target = target.eq(i);
		// 之前插入, 防止多行定位不准
		var top = target.offset().top;
		if(LEA.isMobile) {
			top -= 50;
		}
		var nowTop = $(document).scrollTop();
		// 用$("body").scrllTop(10)没反应 firefox下
		$('html,body').animate({scrollTop: top}, 200);
		return;
	}
}
function genNav() {
	var $con = $("#content");
	var html = $con.html();
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
	$("#blogNavContent").html(titles);
	if(!hs.length) {
		$("#blogNavContent").html(getMsg("none"));
		return false;
	}
	return true;
}

function initNav() {
	var hasNav = genNav();
	if(!hasNav) {
		return;
	}
	
	var $title = $(".title");
	var titlePos = $title.offset();
	var top = titlePos.top + 10;// - $title.height();
	// 手机下不要与标题在同一高度
	if(LEA.isMobile){ 
		top += 30;
	}
	if(top < 0) {
		top = 10;
	}

	var left = $title.width() + titlePos.left - 100;
	$("#blogNav").css("top", top).css("left", left);
	$("#blogNav").show();
	
	$("#blogNavNav").click(function() {
		var $o = $("#blogNavContent");
		if($o.is(":hidden")) {
			$o.show();
		} else {
			$o.hide();
		}
	});
	
	var $d = $(document);
	function reNav() {
	    var vtop = $d.scrollTop();
	    if(vtop <= top) {
			$("#blogNav").css("top", top-vtop);
	    } else {
	    	// 差距很磊了
	    	if(LEA.isMobile) {
				$("#blogNav").css("top", 50);
			} else {
				$("#blogNav").css("top", 10);
			}
	    }
	}
	reNav();
	$(window).scroll(reNav);
}

var C = {
	info: null,
	noteId: noteId,
	preLikeNum: preLikeNum,
	commentNum: commentNum,
	likeBtnO: $("#likeBtn"),
	likeNumO: $("#likeNum"),
	tLikersO: $("#tLikers"),
	likersO: $("#likers"),
	tCommentsO: $("#tComments"),
	commentsO: $("#comments"),
	
	commentBtnO: $("#commentBtn"),
	
	commentsLoadingO: $(".comments-loading"),
	commentsMoreO: $(".comments-more"),
	
	commentBoxO: $(".comment-box"),
	init: function() {
		var self = this;
		if(UserBlogInfo.CanComment && UserBlogInfo.CommentType != "disqus") {
			self.initLikeAndComments();
		} else {
			self.initLike();
		}
		self.initEvent();
		self.incReadNum();
	},
	incReadNum: function() {
		var self = this;
		if(!$.cookie(self.noteId)) {
			$.cookie(self.noteId, 1);
			ajaxGet(staticUrl + "/blog/incReadNum", {noteId: self.noteId});
		}
	},
	initLike: function() {
		var self = this;
		ajaxGet(staticUrl + "/blog/getLike", {noteId: self.noteId}, function(ret) {
			self.info = ret;
			self.toggleLikeBtnActive();
			self.renderLikers();
		});
	},
	initLikeAndComments: function() {
		var self = this;
		ajaxGet(staticUrl + "/blog/getLikeAndComments", {noteId: self.noteId}, function(ret) {
			self.info = ret;
			self.toggleLikeBtnActive();
			self.renderLikers();
			// 是否需要renderComments?
			self.info.commentUserInfo = self.info.commentUserInfo || {};
			// 为了防止第一条评论找不到用户信息情况
			if(visitUserInfo.UserId) {
				self.info.commentUserInfo[visitUserInfo.UserId] = visitUserInfo;
			}
			self.renderComments();
			
			self.commentBoxO.removeClass("hide");
			self.commentsLoadingO.addClass("hide");
			if(self.info.pageInfo.TotalPage > self.info.pageInfo.CurPage) {
				self.commentsMoreO.removeClass("hide");
				self.initMoreComments();
			}
		});
	},
	initMoreComments: function() {
		var self = this;
		self.commentsMoreO.find("a").click(function(){ 
			if(self.info.pageInfo.TotalPage > self.info.pageInfo.CurPage) {
				self.commentsMoreO.addClass("hide");
				self.commentsLoadingO.removeClass("hide");
				ajaxGet(staticUrl + "/blog/listComments", {noteId: self.noteId, page: self.info.pageInfo.CurPage+1}, function(ret) {
					var pageInfo = ret.pageInfo;
					var comments = ret.comments;
					var commentUserInfo = ret.commentUserInfo;
					
					$.extend(self.info.commentUserInfo, commentUserInfo);
					
					// 渲染之
					for(var i in comments) {
						var comment = comments[i];
						comment = self.parseComment(comment);
					}
					var html = self.tCommentsO.render({comments: comments, visitUserInfo: visitUserInfo});
					self.commentsO.append(html);
					
					self.info.pageInfo = pageInfo;
					
					if(self.info.pageInfo.TotalPage > self.info.pageInfo.CurPage) {
						self.commentsMoreO.removeClass("hide");
					} else {
						self.commentsMoreO.addClass("hide");
					}
					
					self.commentsLoadingO.addClass("hide");
				});
			}
		});
	},
	addCommentRender: function(comment){
		var self = this;
		comment = self.parseComment(comment);
		var html = self.tCommentsO.render({blogUrl: blogUrl, comments: [comment], visitUserInfo: visitUserInfo});
		self.commentsO.prepend(html);
		var li = self.commentsO.find("li").eq(0);
		li.hide();
		li.show(500);
		li.addClass("item-highlight");
		setTimeout(function() {
			li.removeClass("item-highlight");
		}, 2000);
	},
	parseComment: function(comment) {
		var self = this;
		var authorUserId = UserInfo.UserId;
		commentUserInfo = self.info.commentUserInfo;
		comment.UserInfo = commentUserInfo[comment.UserId];
		// 是作者自己
		if(visitUserInfo.UserId == UserInfo.UserId) {
			comment.IsMyNote = true;
		}
		if(comment.UserId == authorUserId) {
			comment.IsAuthorComment = true;
		}
		if(comment.UserId == visitUserInfo.UserId) {
			comment.IsMyComment = true;
		}
		// 不是回复自己
		if(comment.ToUserId && comment.ToUserId != comment.UserId) { 
			comment.ToUserInfo = commentUserInfo[comment.ToUserId];
			if(comment.ToUserInfo.UserId == UserInfo.UserId) {
				comment.ToUserIsAuthor = true;
			}
		}
		comment.PublishDate = getDateDiff(Date.parse(goNowToDatetime(comment.CreatedTime)));
		return comment;
	},
	// 渲染评论
	renderComments: function() {
		var self = this;
		var comments = self.info.comments || [];
		if(comments.length == 0) {
			return;
		}
		
		// 整理数据
		// 回复谁, 是否是作者?
		// 回复日期, 几天前, 刚刚
		for(var i in comments) {
			var comment = comments[i];
			comment = self.parseComment(comment);
		}
		var html = self.tCommentsO.render({blogUrl: blogUrl, comments: comments, visitUserInfo: visitUserInfo});
		self.commentsO.html(html);
	},
	
	// 重新渲染likers
	reRenderLikers: function(addMe) {
		var self = this;
		var likedUsers = self.info.likedUsers || [];
		for(var i = 0; i < likedUsers.length; ++i) {
			var user = likedUsers[i];
			if(user.UserId == visitUserInfo.UserId) {
				likedUsers.splice(i, 1);
				break;
			}
		}
		if(addMe) {
			likedUsers = [visitUserInfo].concat(likedUsers);
			self.info.likedUsers = likedUsers;
		}
		self.renderLikers();
	},
	renderLikers: function() {
		var self = this;
		var users = self.info.likedUsers || [];
		var html = self.tLikersO.render({blogUrl: blogUrl, users: users});
		self.likersO.html(html);
	},
	toggleLikeBtnActive: function() {
		var self = this;
		if(self.info.isILikeIt) {
			self.likeBtnO.addClass("active");
		} else {
			self.likeBtnO.removeClass("active");
		}
	},
	commentNumO: $("#commentNum"),
	bindCommentNum: function(fix) {
		var self = this;
		self.commentNum += fix;
		self.commentNumO.text(self.commentNum);
	},
	initEvent: function() {
		var self = this;
		
		// like or not
		self.likeBtnO.click(function() {
			if(!visitUserInfo.UserId) {
				needLogin();
				return;
			}
			ajaxPost(staticUrl + "/blog/likeBlog", {noteId: self.noteId}, function(ret) {
				if(ret.Ok) {
					// like
					if(ret.Item) {
						var num = self.preLikeNum+1;
					} else {
						var num = self.preLikeNum-1;
					}
					self.preLikeNum = num >= 0 ? num : 0;
					self.likeNumO.text(self.preLikeNum);
					self.info.isILikeIt = ret.Item;
					self.toggleLikeBtnActive();
					
					// 重新render likers
					// 我是否在列表中
					self.reRenderLikers(ret.Item);
				}
			});
		});
		
		// 显示回复回复
		$("#comments").on("click", ".comment-reply", function() {
			var form = $(this).closest("li").find("form");
			if(form.is(":hidden")) {
				form.show();
				form.find("textarea").focus();
			} else {
				form.hide();
			}
		});
		$("#comments").on("click", ".reply-cancel", function() {
			$(this).closest("form").hide();
		});
		
		// 回复
		$(".comment-box").on("click", ".reply-comment-btn", function(e) {
			e.preventDefault();
			var commentId = $(this).data("comment-id");
			var $form = $(this).closest("form");
			var $content = $form.find("textarea");
			var content = $.trim($content.val());
			if(!content) {
				$content.focus();
				return;
			}
			var t = $(this);
			t.button("loading");
			var data = {noteId: self.noteId, toCommentId: commentId, content: content};
			ajaxPost(staticUrl + "/blog/comment", data, function(ret) {
				t.button("reset");
				$content.val("");
				self.bindCommentNum(1);
				if(commentId) {
					$form.hide();
				}
				
				if(commentId) {
					scrollToTarget("#comments", -200);
				}
				
				// 添加一个
				self.addCommentRender(ret.Item);
			});
		});
		
		// 删除
		$(".comment-box").on("click", ".comment-trash", function(e) {
			var commentId = $(this).parent().data("comment-id");
			var t = this;
			BootstrapDialog.confirm(getMsg("confirmDeleteComment"), function(yes) {
				if(yes) {
					ajaxPost(staticUrl + "/blog/deleteComment", {noteId: self.noteId, commentId: commentId}, function(ret) {
						if(ret.Ok) {
							var li = $(t).closest("li");
							li.hide(500); // remove();
							setTimeout(function() {
								li.remove();
							}, 300);
							
							self.bindCommentNum(-1);
						}
					});
				}
			});
		});
		
		// 点zan
		$(".comment-box").on("click", ".comment-like", function(e) {
			var commentId = $(this).parent().data("comment-id");
			var t = this;
		
			ajaxPost(staticUrl + "/blog/likeComment", {commentId: commentId}, function(re) {
				if(re.Ok) {
					var ret = re.Item;
					if(ret.Num <= 0) {
						$(t).parent().find(".like-num").addClass("hide");
					} else {
						$(t).parent().find(".like-num").removeClass("hide");
						$(t).parent().find(".like-num-i").text(ret.Num)
					}
					if(ret.IsILikeIt) {
						$(t).find(".like-text").text(getMsg("unlike"));
					} else {
						$(t).find(".like-text").text(getMsg('like'));
					}
				}
			});
		});
		
		// 举报
		function report(commentId, noteId, title) {
			var form = $("#reportMsg").html();
			var body;
	        var input;
	        var isOver = false;
			var modal = BootstrapDialog.show({
	            title: title,
	            message: form,
	            nl2br: false,
	            buttons: [{
                    label: getMsg("cancel"),
                    action: function(dialog) {
                        dialog.close();
                    }
                }, {
                    label: getMsg("confirm"),
                    cssClass: 'btn-primary',
                    action: function(dialog) {
                    	if(isOver) {
                    		dialog.close();
                    	}
                    	var val = body.find("input[type='radio']:checked").val();
                    	if(!val) {
                    		var val = body.find(".input-container input").val();
                    	}
                    	if(!val) {
                    		body.find(".footnote").html(getMsg("chooseReason"));
                    		return;
                    	}
                    	ajaxPost(staticUrl + "/blog/report", {commentId: commentId, noteId: noteId, reason: val}, function(re) {
                    		isOver = true;
                    		if(reIsOk(re)) {
		                        body.html(getMsg("reportSuccess"));
                    		} else {
		                        body.html(getMsg("error"));
                    		}
                			setTimeout(function() {
		                        dialog.close();
	                        }, 3000);
                    	});
                    }
                }]
	        });
	        body = modal.getModalBody();
	        input = body.find(".input-container");
	        body.find("input[type='radio']").click(function(){ 
	        	if(!$(this).val()) {
	        		input.show();
	        		input.find("input").focus();
	        	} else {
	        		input.hide();
	        	}
	        });
		}
		$(".comment-box").on("click", ".comment-report", function() {
			if(needLogin()) {
				return;
			}
			var commentId = $(this).parent().data("comment-id");
			report(commentId, self.noteId, getMsg("reportComment?"));
		});
		$("#reportBtn").click(function() {
			if(needLogin()) {
				return;
			}
			report("", self.noteId, getMsg("reportBlog?"));
		});
		
		self.initShare();
	},
	weixinQRCodeO: $("#weixinQRCode"),
	initShare: function() {
		var self = this;
		$(".btn-weixin").click(function() {
			if(!self.weixinQRCodeO.html()) {
				self.weixinQRCodeO.qrcode(viewUrl + "/" + self.noteId);
			}
			BootstrapDialog.show({
	            title: getMsg('scanQRCode'),
	            message: self.weixinQRCodeO
	        });
		});
		
		$(".btn-share").click(function() {
			var $this = $(this);
			var map = {"btn-weibo": shareSinaWeibo, "tencent-weibo": shareTencentWeibo, "qq": shareQQ, "renren": shareRenRen};
			for(var i in map) {
				if($this.hasClass(i)) {
					map[i](self.noteId, document.title);
					break;
				}	
			}
		});
	}
}

$(function() {
	C.init();
});