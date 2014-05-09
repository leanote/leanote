/**
 * plugin.js
 *
 * Copyright, Moxiecode Systems AB
 * Released under LGPL License.
 *
 * License: http://www.tinymce.com/license
 * Contributing: http://www.tinymce.com/contributing
 */

/*global tinymce:true */

// 与jbimages结合
// 使用了jbimages/dialog-v4.htm
// jbimages/js/dialog-v4.js
tinymce.PluginManager.add('image', function(editor, url) {
	function getImageSize(url, callback) {
		var img = document.createElement('img');

		function done(width, height) {
			img.parentNode.removeChild(img);
			callback({width: width, height: height});
		}

		img.onload = function() {
			done(img.clientWidth, img.clientHeight);
		};

		img.onerror = function() {
			done();
		};

		img.src = url;

		var style = img.style;
		style.visibility = 'hidden';
		style.position = 'fixed';
		style.bottom = style.left = 0;
		style.width = style.height = 'auto';

		document.body.appendChild(img);
	}

	//??
	function createImageList(callback) {
		return function() {
			var imageList = editor.settings.image_list;

			if (typeof(imageList) == "string") {
				tinymce.util.XHR.send({
					url: imageList,
					success: function(text) {
						callback(tinymce.util.JSON.parse(text));
					}
				});
			} else {
				callback(imageList);
			}
		};
	}

	// 弹框
	function showDialog(imageList) {
		var win, data = {}, dom = editor.dom, imgElm = editor.selection.getNode();
		var width, height, imageListCtrl;

		function buildImageList() {
			var imageListItems = [{text: 'None', value: ''}];

			tinymce.each(imageList, function(image) {
				imageListItems.push({
					text: image.text || image.title,
					value: editor.convertURL(image.value || image.url, 'src'),
					menu: image.menu
				});
			});

			return imageListItems;
		}

		function recalcSize(e) {
			var widthCtrl, heightCtrl, newWidth, newHeight;

			widthCtrl = win.find('#width')[0];
			heightCtrl = win.find('#height')[0];

			newWidth = widthCtrl.value();
			newHeight = heightCtrl.value();

			if (win.find('#constrain')[0].checked() && width && height && newWidth && newHeight) {
				if (e.control == widthCtrl) {
					newHeight = Math.round((newWidth / width) * newHeight);
					heightCtrl.value(newHeight);
				} else {
					newWidth = Math.round((newHeight / height) * newWidth);
					widthCtrl.value(newWidth);
				}
			}

			width = newWidth;
			height = newHeight;
		}

		function onSubmitForm() {
			function waitLoad(imgElm) {
				function selectImage() {
					imgElm.onload = imgElm.onerror = null;
					editor.selection.select(imgElm);
					editor.nodeChanged();
				}

				// 如果没有设置width, height, 就用图片本身的大小
				imgElm.onload = function() {
					if (!data.width && !data.height) {
						dom.setAttribs(imgElm, {
							width: imgElm.clientWidth,
							height: imgElm.clientHeight
						});
					}

					selectImage();
				};

				imgElm.onerror = selectImage;
			}

			// 这是通过url插入图片
			// iframe里得到...
			var data = {};
			var ifr = $("iframe[src$='v4.htm']").contents().find("body"); 
			data.width = $(ifr).find("#imageWidth").val();
			data.height = $(ifr).find("#imageHeight").val();
			data.src = $(ifr).find("#imageSrc").val();
			data.alt = "";

			if (data.width === '') {
				data.width = null;
			}

			if (data.height === '') {
				data.height = null;
			}

			if (data.style === '') {
				data.style = null;
			}

			data = {
				src: data.src,
				alt: data.alt,
				width: data.width,
				height: data.height,
				style: data.style
			};

			editor.undoManager.transact(function() {
				// 删除图片
				if (!data.src) {
					if (imgElm) {
						dom.remove(imgElm);
						editor.nodeChanged();
					}

					return;
				}

				if (!imgElm) {
					data.id = '__mcenew';
					editor.selection.setContent(dom.createHTML('img', data));
					imgElm = dom.get('__mcenew');
					dom.setAttrib(imgElm, 'id', null);
				} else {
					dom.setAttribs(imgElm, data);
				}

				waitLoad(imgElm);
			});

			tinymce.EditorManager.activeEditor.windowManager.close(window);
		}

		function removePixelSuffix(value) {
			if (value) {
				value = value.replace(/px$/, '');
			}

			return value;
		}

		// 当url改变时, 得到图片的大小
		function srcChange() {
			if (imageListCtrl) {
				imageListCtrl.value(editor.convertURL(this.value(), 'src'));
			}

			getImageSize(this.value(), function(data) {
				if (data.width && data.height) {
					width = data.width;
					height = data.height;

					win.find('#width').value(width);
					win.find('#height').value(height);
				}
			});
		}

		width = dom.getAttrib(imgElm, 'width');
		height = dom.getAttrib(imgElm, 'height');

		if (imgElm.nodeName == 'IMG' && !imgElm.getAttribute('data-mce-object')) {
			data = {
				src: dom.getAttrib(imgElm, 'src'),
				alt: dom.getAttrib(imgElm, 'alt'),
				width: width,
				height: height
			};
		} else {
			imgElm = null;
		}

		if (imageList) {
			imageListCtrl = {
				type: 'listbox',
				label: 'Image list',
				values: buildImageList(),
				value: data.src && editor.convertURL(data.src, 'src'),
				onselect: function(e) {
					var altCtrl = win.find('#alt');

					if (!altCtrl.value() || (e.lastControl && altCtrl.value() == e.lastControl.text())) {
						altCtrl.value(e.control.text());
					}

					win.find('#src').value(e.control.value());
				},
				onPostRender: function() {
					imageListCtrl = this;
				}
			};
		}

		// General settings shared between simple and advanced dialogs
		var generalFormItems = [
			{name: 'src', type: 'filepicker', filetype: 'image', label: 'Source', autofocus: true, onchange: srcChange},
			imageListCtrl,
			{name: 'alt', type: 'textbox', label: 'Image description'},
			{
				type: 'container',
				label: 'Dimensions',
				layout: 'flex',
				direction: 'row',
				align: 'center',
				spacing: 5,
				items: [
					{name: 'width', type: 'textbox', maxLength: 3, size: 3, onchange: recalcSize},
					{type: 'label', text: 'x'},
					{name: 'height', type: 'textbox', maxLength: 3, size: 3, onchange: recalcSize},
					{name: 'constrain', type: 'checkbox', checked: true, text: 'Constrain proportions'}
				]
			}
		];

		function updateStyle() {
			function addPixelSuffix(value) {
				if (value.length > 0 && /^[0-9]+$/.test(value)) {
					value += 'px';
				}

				return value;
			}

			var data = win.toJSON();
			var css = dom.parseStyle(data.style);

			dom.setAttrib(imgElm, 'style', '');

			delete css.margin;
			css['margin-top'] = css['margin-bottom'] = addPixelSuffix(data.vspace);
			css['margin-left'] = css['margin-right'] = addPixelSuffix(data.hspace);
			css['border-width'] = addPixelSuffix(data.border);

			win.find('#style').value(dom.serializeStyle(dom.parseStyle(dom.serializeStyle(css))));
		}

		// 与jbimages结合
		// Simple default dialog
		win = editor.windowManager.open({
			title: 'Insert/edit image',
			data: data,
			file : url + '/../jbimages/dialog-v4.htm',
			// body: generalFormItems,
			onSubmit: onSubmitForm,
			width : 550,
			height: 305,
			buttons: [{
				text: '确定',
				classes: 'widget btn primary first abs-layout-item',
				// disabled: false,
				onclick: onSubmitForm
				},
				{
					text: 'Close',
					onclick: 'close'
				}]
		});
	}

	// 添加按钮
	editor.addButton('image', {
		icon: 'image',
		tooltip: 'Insert/edit image',
		onclick: createImageList(showDialog),
		stateSelector: 'img:not([data-mce-object])'
	});
});