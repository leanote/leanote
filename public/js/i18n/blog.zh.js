var MSG = {"a":"a","aboutMe":"关于我","author":"作者","baseInfoSet":"基本设置","blog":"博客","blogClass":"分类","blogDesc":"博客描述","blogLogo":"博客Logo","blogLogoTips":"上传logo将显示logo(替代博客标题)","blogName":"博客标题","blogNav":"导航","blogNavs":"导航","blogSet":"博客设置","cancel":"取消","chooseReason":"请选择举报理由","comment":"评论","commentSet":"评论设置","commentSys":"leanote 使用 \u003ca href=\"http://disqus.com\" target=\"_blank\"\u003eDisqus\u003c/a\u003e 作为评论系统","comments":"条评论","community":"社区","confirm":"确认","confirmDeleteComment":"确定删除该评论?","createdTime":"创建","daysAgo":"天前","delete":"删除","disqusHelp":"请填写您申请的Disqus唯一url前缀. 建议您申请Disqus帐号, 这样可以自己管理评论. 或使用leanote的默认Disqus Id. ","elegant":"大气","error":"错误","fullBlog":"全文","home":"主页","hoursAgo":"个小时前","justNow":"刚刚","latestPosts":"最近发表","like":"赞","minutesAgo":"分钟前","monthsAgo":"个月前","moreShare":"更多分享","navFixed":"导航左侧固定","needHelp":"需要帮助?","noBlog":"无博客","noTag":"无","none":"无","openComment":"开启评论?","other":"其它","qqZone":"QQ空间","quickLinks":"快速链接","renren":"人人网","reply":"回复","report":"举报","reportBlog?":"举报该博客?","reportComment?":"举报该评论?","reportReason":"举报理由","reportReason1":"不友善内容","reportReason2":"广告等垃圾信息","reportReason3":"违法违规内容","reportReason4":"不宜公开讨论的政治内容","reportSuccess":"举报成功, 我们处理后会通知作者, 感谢您的监督","saveSuccess":"保存成功","scanQRCode":"打开微信扫一扫二维码","signIn":"登录","signUp":"注册","sinaWeibo":"新浪微博","submitComment":"发表评论","tencentWeibo":"腾讯微博","theme":"主题","themeSet":"主题设置","unlike":"取消赞","updatedTime":"更新","viewers":"人读过","weeksAgo":"周前","weixin":"微信"};
function getMsg(key, data) {
	var msg = MSG[key]
	if(msg) {
		if(data) {
			if(!isArray(data)) {
				data = [data];
			}
			for(var i = 0; i < data.length; ++i) {
				msg = msg.replace("%s", data[i]);
			}
		}
		return msg;
	}
	return key;
}