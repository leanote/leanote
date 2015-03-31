// leanote blog share & comment
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
		if(blogInfo.OpenComment && blogInfo.CommentType != "disqus") {
			self.initLikesAndComments();
		} else {
			self.initLikes();
		}
		self.initEvent();
		self.incReadNum();
	},
	
	// 博客的统计信息
	getPostStat: function() {
	},
	// 增加阅读量
	incReadNum: function() {
		incReadNum(self.noteId);
	},
	initLikes: function() {
		var self = this;
		getLikes(self.noteId, function(ret) {
			if(reIsOk(ret)) {
				ret = ret.Item;
				self.info = ret;
				self.toggleLikeBtnActive();
				self.renderLikers();
				self.visitUserInfo = self.info.visitUserInfo || {};
			}
		});
	},
	initLikesAndComments: function() {
		var self = this;
		getLikesAndComments(self.noteId, function(ret) {
			if(reIsOk(ret)) {
				ret = ret.Item;
				self.info = ret;
				self.toggleLikeBtnActive();
				self.renderLikers();
				// 是否需要renderComments?
				self.info.commentUserInfo = self.info.commentUserInfo || {};
				self.visitUserInfo = self.info.visitUserInfo || {};
				// 为了防止第一条评论找不到用户信息情况
				if(self.visitUserInfo.UserId) {
					self.info.commentUserInfo[self.visitUserInfo.UserId] = self.visitUserInfo;
				}
				self.renderComments();
				
				// 之前是隐藏的, 取消之
				self.commentBoxO.removeClass("hide");
				self.commentsLoadingO.addClass("hide");
				if(self.info.pageInfo.TotalPage > self.info.pageInfo.CurPage) {
					self.commentsMoreO.removeClass("hide");
					self.initMoreComments();
				}
				
				// 是否已经登录?
				if(self.visitUserInfo.UserId) {
					$("#commentForm").removeClass("hide");
					$("#visitUserLogo").attr("src", self.visitUserInfo.Logo);
				} else {
					$("#noLoginContainer").removeClass("hide");
				}
			}
		});
	},
	initMoreComments: function() {
		var self = this;
		self.commentsMoreO.find("a").click(function(){ 
			if(self.info.pageInfo.TotalPage > self.info.pageInfo.CurPage) {
				self.commentsMoreO.addClass("hide");
				self.commentsLoadingO.removeClass("hide");
				getComments(self.noteId, self.info.pageInfo.CurPage+1, function(ret) {
					var pageInfo = ret.pageInfo;
					var comments = ret.comments;
					var commentUserInfo = ret.commentUserInfo;
					
					$.extend(self.info.commentUserInfo, commentUserInfo);
					
					// 渲染之
					for(var i in comments) {
						var comment = comments[i];
						comment = self.parseComment(comment);
					}
					var html = self.tCommentsO.render({comments: comments, visitUserInfo: self.visitUserInfo});
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
		var html = self.tCommentsO.render({comments: [comment], visitUserInfo: self.visitUserInfo});
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
		var authorUserId = blogInfo.UserId;
		commentUserInfo = self.info.commentUserInfo;
		comment.UserInfo = commentUserInfo[comment.UserId];
		// 是作者自己
		if(self.visitUserInfo.UserId == authorUserId) {
			comment.IsMyNote = true;
		}
		if(comment.UserId == authorUserId) {
			comment.IsAuthorComment = true;
		}
		if(comment.UserId == self.visitUserInfo.UserId) {
			comment.IsMyComment = true;
		}
		// 不是回复自己
		if(comment.ToUserId && comment.ToUserId != comment.UserId) { 
			comment.ToUserInfo = commentUserInfo[comment.ToUserId];
			if(comment.ToUserInfo.UserId == authorUserId) {
				comment.ToUserIsAuthor = true;
			}
		}
		comment.PublishDate = getDateDiff(goNowToDate(comment.CreatedTime));
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
		var html = self.tCommentsO.render({comments: comments, visitUserInfo: self.visitUserInfo});
		self.commentsO.html(html);
	},
	
	// 重新渲染likers
	reRenderLikers: function(addMe) {
		var self = this;
		var likedUsers = self.info.likedUsers || [];
		for(var i = 0; i < likedUsers.length; ++i) {
			var user = likedUsers[i];
			if(user.UserId == self.visitUserInfo.UserId) {
				likedUsers.splice(i, 1);
				break;
			}
		}
		if(addMe) {
			likedUsers = [self.visitUserInfo].concat(likedUsers);
			self.info.likedUsers = likedUsers;
		}
		self.renderLikers();
	},
	renderLikers: function() {
		var self = this;
		var users = self.info.likedUsers || [];
		var html = self.tLikersO.render({users: users});
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
			if(!self.visitUserInfo.UserId) {
				needLogin();
				return;
			}
			likePost(noteId, function(ret) {
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
			commentPost(self.noteId, commentId, content, function(ret) {
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
			try {
				BootstrapDialog.confirm("Are you sure?", function(yes) {
					if(yes) {
						deleteComment(noteId, commentId, function(ret) {
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
			} catch(e) {}
		});
		
		// 点zan
		$(".comment-box").on("click", ".comment-like", function(e) {
			var commentId = $(this).parent().data("comment-id");
			var t = this;
			
			likeComment(commentId, function(re) {
				if(re.Ok) {
					var ret = re.Item;
					if(ret.Num <= 0) {
						$(t).parent().find(".like-num").addClass("hide");
					} else {
						$(t).parent().find(".like-num").removeClass("hide");
						$(t).parent().find(".like-num-i").text(ret.Num)
					}
					if(ret.IsILikeIt) {
						var ever = $(t).find(".like-text").text();
						if(ever == "赞") {
							$(t).find(".like-text").text("取消赞");
						} else {
							$(t).find(".like-text").text("Unlike");
						}
					} else {
						var ever = $(t).find(".like-text").text();
						if(ever == "取消赞") {
							$(t).find(".like-text").text("赞");
						} else {
							$(t).find(".like-text").text("Like");
						}
					}
				}
			});
		});
		self.initShare();
	},
	weixinQRCodeO: $("#weixinQRCode"),
	initShare: function() {
		var self = this;
		$(".btn-weixin").click(function() {
			if(!self.weixinQRCodeO.html()) {
				self.weixinQRCodeO.qrcode(location.href);
			}
			BootstrapDialog.show({
	            title: "Open Wechat to scan the code",
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