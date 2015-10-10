/* /index 首页 */

$(function() {
    // 平滑滚动
    $(".smooth-scroll").click(function(e) {
        e.preventDefault();
        var t = $(this).data("target");
        var targetOffset = $(t).offset().top - 80;
        $('html,body').animate({scrollTop: targetOffset}, 300);
    });
    function slider(webImgs, descT) {
        var webImgsLen = webImgs.length;
        var curIndex = 0;
        setInterval(function() {
            webImgs.eq(curIndex).stop().animate({opacity: '0'}, 1000);
            curIndex = (curIndex+1)%webImgsLen;
            var curImg = webImgs.eq(curIndex);
            curImg.stop().animate({opacity: '1'}, 1000);
            descT.text(curImg.data("text"));
        }, 5000);
    }
    slider($(".web-slider"), $("#webText"));
    slider($(".mobile-slider"), $("#mobileText"));
    
    function setCookie(name, value) {
        var Days = 10*365;  
        var exp  = new Date();  
        exp.setTime(exp.getTime() + Days*24*60*60*1000);  
        document.cookie = name + "="+ escape (value) + ";expires=" + exp.toGMTString();  
    }
    
    $('#lang a').click(function() {
        var lang = $(this).data('lang');
        setCookie('LEANOTE_LANG', lang);
        location.reload();
    });
});
