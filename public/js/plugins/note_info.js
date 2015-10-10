/**
 * @file 笔记信息
 * @author  life
 * 
 */
define('note_info', [], function() {
    var tpl = ['<table>',
        '<tr><th>' + getMsg('Create Time') + '</th><td id="noteInfoCreatedTime"></td></tr>',
        '<tr><th>' + getMsg('Update Time') + '</th><td id="noteInfoUpdatedTime"></td></tr>',
        '<tr class="post-url-tr">',
            '<th>' +  getMsg('Post Url') + '</th>',
            '<td>',
                '<div class="post-url-wrap">',
                        '<span class="post-url-base">http://blog.leanote.com/life/post/</span><span><span class="post-url-text">life-life-life-a-leanote</span>',
                        '<input type="text" class="form-control">',
                        '</span>',
                        ' <a class="post-url-pencil" title="' + getMsg('update') + '"><i class="fa fa-pencil"></i></a>',
                    '</div>',
                '</td>',
            '</tr>',
        '</table>'].join('');
    var $tpl = $(tpl);

    var $noteInfoCreatedTime = $tpl.find('#noteInfoCreatedTime');
    var $noteInfoUpdatedTime = $tpl.find('#noteInfoUpdatedTime');
    var $noteInfoPostUrl = $tpl.find('#noteInfoPostUrl');

    var $noteInfoPostUrlTr = $tpl.find('.post-url-tr');
    var $postUrlWrap = $tpl.find('.post-url-wrap'); 
    var $input = $tpl.find('input');

    var $postUrlBase = $tpl.find('.post-url-base');
    var $postUrlText = $tpl.find('.post-url-text');

    var view = {

        $noteInfo: $('#noteInfo'),

        note: null,

        bind: function () {
            var me = this;
            $('#noteInfoDropdown').click(function () {
                me.render();
            });

            $tpl.find('.post-url-pencil').click(function () {
                $postUrlWrap.addClass('post-url-edit');
                $input.val(decodeURI(me.note.UrlTitle));
                $input.focus();
            });
            $input.keydown(function (e) {
                if(e.keyCode === 13) {
                    $input.blur();
                }
            });
            $input.blur(function () {
                $postUrlWrap.removeClass('post-url-edit');

                var val = $input.val();
                if (!val) {
                    return;
                }

                ajaxPost("/member/blog/updateBlogUrlTitle", {noteId: me.note.NoteId, urlTitle: val}, function(re) {
                    if(reIsOk(re)) {
                        var encodedUrl = encodeURI(re.Item);
                        me.note.UrlTitle = encodedUrl;
                        $postUrlText.text(decodeURI(me.note.UrlTitle));
                    } else {
                        alert(re.Msg || "error");
                    }
                });
            });

            // 当笔记Change时, 重新render
            LEA.on('noteChanged', function (note) {
                me.render(note);
            });
        },

        getPostUrl: function (note) {
            return '';
        },

        rendered: false,
        render: function (note) {
            var me = this;
            if (!note) {
                note = Note.getCurNote();
            }
            if (!note) {
                return;
            }
            me.note = note;

            $noteInfoCreatedTime.html(goNowToDatetime(note.CreatedTime));
            $noteInfoUpdatedTime.html(goNowToDatetime(note.UpdatedTime));

            if (!note.IsBlog) {
                $noteInfoPostUrlTr.addClass('hide');
            }
            else {
                $noteInfoPostUrlTr.removeClass('hide');

                // post-url
                $postUrlBase.text(UserInfo.PostUrl + '/');
                $postUrlText.text(decodeURI(note.UrlTitle));
            }

            if (!me.rendered) {
                me.$noteInfo.html($tpl);
                me.rendered = true;
            }
        },

        init: function () {
            this.bind();
        }
    };

    view.init();
});