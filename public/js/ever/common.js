/**
 * 命名空间与全局方法
 * @author life
 */

// 命名空间
var dk = 
	order = 
	shipment = 
	gis = 
	party = 
	contract = 
	network = 
	equip =
	manual = 
	reason =
	payment = 
	auth = 
	cron = 
	{};

// 全局配置
G_CONFIG = {
	defaultCountryId: "402895853e3bd7a5013e3bd7b2c20000",
	split: "/"
};

// cache
dk.cache = {
	countries: {},
	regions: {},
	cities: {},
	towns: {}
};

// 公用方法
// 除公用方法外, 其它方法必须在命名空间下

//-----------
// 国际化
//-----------

/**
 * 得到国际化
 * @param key String 键名
 * @param data Array 替换的值
 * @param [source] Object 数据源
 * @returns String 国际化后的值
 */
function getText(key, data, source) {
	// 判断key有无'.', 比如 msg.title.info
	// 得到最终的value
	var keyArr = key.split('.');
	var value = source ? source : i18n;
	for(var i = 0; i < keyArr.length; ++i) {
		value = value[keyArr[i]];
	}
	
	if(!data) return value;

	// 替换之
	// 防止{1}/{1}, 里面的数字是有意义的!
	var i = 0;
	var str = value.replace(/{[0-9]+}/g, function(each) {
		// 从{11}中得到11
		var index = each.substr(1, each.length-2);
		return data[index];
	});
	
	return str;
}

/**
 * 
 * @param key
 * @param data
 * @returns
 */
function getMsgText(key, data) {
	return getText(key, data, i18n.msg);
}
function getMsgTitleText(key, data) {
	return getText(key, data, i18n.msg.title);
}

function getWinText(key, data) {
	return getText(key, data, i18n.win);
}
function getWinTitleText(key, data) {
	return getText(key, data, i18n.win.title);
}

//---------------
// dataGrid 表格
//---------------

/**
 * 检查是否选中了
 * 
 * @return Object row
 */
function checkSelectOne(gridId) {
	var datagrid = $('#' + gridId);
	if(!datagrid) return false;
	
	var rows = datagrid.datagrid('getChecked');
	if(rows.length < 1) {
		msgAlert('info', 'selectOne', 'info');
		
		return false;
	}
	if(rows.length > 1) {
		msgAlert('info', 'selectOnlyOne', 'info');
		return false;
	}
	
	return rows[0];
}

/**
 * 可以选中多行
 */
function checkSelects(gridId) {
	var datagrid = $('#' + gridId);
	if(!datagrid) return false;
	
	var rows = datagrid.datagrid('getChecked');
	if(rows.length < 1) {
		msgAlert('info', 'select', 'info');
		return false;
	}
	
	return rows;
}

/**
 * 表格移除多条
 * @param gridId String
 * @param rows Array[Object]
 */
function deleteRows(gridId, rows) {
	for(var i = 0; i < rows.length; ++i) {
		deleteRow(gridId, rows[i]);
	}
}

/**
 * 表格移除一条
 * @param gridId String
 * @param row Object
 */
function deleteRow(gridId, row) {
	var index = $('#' + gridId).datagrid('getRowIndex', row);
	$('#' + gridId).datagrid('deleteRow', index);
}

/**
 * 表格添加多条
 * @param gridId String
 * @param rows Array[Object]
 */
function appendRows(gridId, rows) {
	for(var i = 0; i < rows.length; ++i) {
		appendRow(gridId, rows[i]);
	}
}

/**
 * 表格添加一条
 * @param gridId String
 * @param row Object
 */
function appendRow(gridId, row) {
	$("#" + gridId).datagrid('appendRow', row);
}

/**
 * 得到多行记录的主键列表
 * @param rows
 * @param fieldId 主键字段名
 * @returns {Array}
 */
function getRowsIds(rows, fieldId) {
	var ids = [];
	for(var i in rows) {
		ids.push(rows[i][fieldId]);
	}
	return ids;
}

/**
 * 刷新dataGrid
 * @param gridId
 */
function reloadGrid(gridId) {
	$('#' + gridId).datagrid('reload');
}

/**
 * 为grid加载数据
 * @param gridId
 * @param queryParams
 * @param [url]
 */
function loadGrid(gridId, queryParams, url) {
	var gridObj = $('#' + gridId);
	if(url) {
		gridObj.datagrid('options').url = url;
	}
	gridObj.datagrid('options').queryParams = queryParams;
    gridObj.datagrid('load');
}

/**
 * 表格添加标题
 * 显示数据用
 * @param title
 * @returns {String}
 */
function addTitle(title) {
	return '<span title="' + title + '">' + title + '</span>';
}

/**
 * 将true/false转成是/否
 * @param bool true or false
 * @returns
 */
function getYesOrNo(bool) {
	return bool ? getText('yes') : getText('no');
}

//------------
// 弹框
//------------

// 弹框默认配置
dk.winDefaultConfig = {
	minimizable: false, // 不可最小化
    resizable: true,
	collapsible: false, // 不能折叠
    modal: true, // 遮罩
    closed: false,
    onClose : function() {
		$(this).window('destroy');
	}
};

// 得到grid 宽度
// 传过来的winWidth是win的宽度, 根据win宽度调整grid宽度
// min是否有下拉滚动条, 有会减少grid宽度
function getGridWidth(winWidth, min) {
	var borderWidth = 16;
	var borderWidthm = 36;
	if(min) {
		return winWidth - borderWidthm;
	}
	return winWidth - borderWidth;
}

// 得到win的宽度
// size是规格大小, selfWidth是自定义的大小, 如果有该值, 表示最大规格是size
// selfWidth超过它将用size
function getWinWidth(size, selfWidth) {
	var width = $('body').width() * 0.1 * size;
	if(typeof selfWidth != 'undefined') {
		if(width > selfWidth) return selfWidth;
	}
	return width;
}

// 得到win的高度
function getWinHeight(size, selfHeight) {
	// 这里之前是$('body').height(), 有bug, 当最大化win后, 该值为0
	var height = $(window).height() * 0.1 * size;
	if(typeof selfHeight != 'undefined') {
		if(height > selfHeight) return selfHeight;
	}
	return height;
}

/**
 * 打开win弹框
 * @param winVar bool 是否自定义win, 如果为false, 则用g_custom_win, 否则用winVar
 * @param winId string [不要使用, 将会弃用]
 */
dk.win = {};
function openWin(config, winVar, winId) {
	config = $.extend({}, dk.winDefaultConfig, config);
	if(!winId) winId = 'newWin'; // + (new Date().getTime());
	
	// win的标题
	config.title = getWinTitleText(config.title);
	// 调整宽度和高度
	config.width = getWinWidth(9, config.width);
	config.height = getWinHeight(9, config.height);
		
	var win = $('<div id="'+ winId+ '" class="newWin"/>').window(config);
	
	if(!winVar) { // 用全局的
		dk.win['g_custom_win'] = win;
	} else {
		dk.win[winVar] = win; // 避免不是由自己的子来关闭的情况, 其它可以控制关闭
		// window[winVar] = win; 
	}
}

// 关闭弹框
// 空 or target Object或创建时定义的全局变量名
function closeWin(target) {
	if(target == undefined) {
		var theEvent = window.event || arguments.callee.caller.arguments[0];
		var theObj = theEvent.target || theEvent.srcElement;	
		closeWin(theObj);
	} else if(typeof target == "object") {
		// 向上找class="newWin"的div
		var obj = $(target).parents("div.newWin");
		if(obj) {
			$(obj).window('destroy');
		}
	} else if(typeof target == "string") {
		if(dk.win[target]){
			dk.win[target].window('destroy');
		}
	}
}

/**
 * 关闭win, 更新grid
 * @param winVar
 * @param gridId
 */
function closeWinAndReloadGrid(winVar, gridId) {
	closeWin(winVar);
	reloadGrid(gridId);
}

//---------
// 表单
//---------

/**
 * 提交数据
 * @param formId string <form id="formId"> 会通过该ID取该form的serialize数据
 * @param successFunc Function 成功处理函数
 * @param failureFunc Function 失败处理函数
 */
function formSubmit(formId, formDataFunc, successFunc, failureFunc, hasProgress) {
	hasProgress = hasProgress == undefined ? true : hasProgress;
	if(hasProgress) {
		showProgress('info', 'processing');
	}
	
	formId = $('#' + formId);
	var isValid = $(formId).form('validate'); // 验证数据
	if(!isValid) {
		hideProgress();
		msgAlert("info", "inputError", "info");
		return;
	}
	
	var formData = $(formId).serialize();
	if(typeof formDataFunc == 'function') {
		formData = formDataFunc(formData);
	}
	
	ajaxPost($(formId).attr('action'), formData, function(ret) {
		_ajaxCallback(ret, successFunc, failureFunc);
		hideProgress();
	});
}

/**
 * form的数据通过Json形式提交到后台
 * @param formId
 * @param formDataFunc 数据处理方法
 * @param successFunc
 * @param failureFunc
 * @author Life
 * @datetime 2013/4/20 18:00
 */
function formSubmitJson(formId, formDataFunc, successFunc, failureFunc, hasProgress) {
	hasProgress = hasProgress == undefined ? true : hasProgress;
	if(hasProgress) {
		showProgress('info', 'processing');
	}
	var formId2 = $('#' + formId);
	var isValid = $(formId2).form('validate'); // 验证数据
	
	if(!isValid) {
		hideProgress();
		msgAlert("info", "inputError", "info");
		return;
	}
	
	// 得到form Json数据
	var formData = getFormJsonData(formId);
	if(typeof formDataFunc == 'function') {
		formData = formDataFunc(formData);
	}
	
	$.ajax({
	    url : $(formId2).attr('action'),
	    type : "POST",
	    contentType: "application/json; charset=utf-8",
	    datatype: "json",
	    data : JSON.stringify(formData),
	    success : function(data, stats) {
	    	_ajaxCallback(data, successFunc, failureFunc);
	    	if(hasProgress) {
	    		hideProgress();
	    	}
	    },
		error: function(ret) {
			_ajaxCallback(ret, successFunc, failureFunc);
			if(hasProgress) {
				hideProgress();
			}
		}
	});
}

// 得到form的数据
// 返回json
// 会处理date, autocomplete值
function getFormJsonData(formId, needFilter) {
	if(needFilter == undefined) {
		needFilter = true;
	}
	var data = formArrDataToJson($('#' + formId).serializeArray());
	
	// date->dbType，如果是date类型，dbType是timestamp类型，要转
	$("#" + formId + " input.Wdate").each(function() {
		var name = $(this).attr("name");
		if(name && data[name]) {
			var viewType = $(this).attr("viewType");
			var dbType = $(this).attr("dbType");
			if(viewType == dbType) {
				return;
			}
			if(viewType == "date" && dbType == "datetime") {
				data[name] = data[name] + " 00:00:00";
			}
		}
	});
	
	/*if(needFilter) {
		// 如果是autocomplete，那么不需要该值, 删除它
		$("#" + formId + " input.ac_input").each(function() {
			var name = $(this).attr("name");
			if(!name) {
				return;
			}
			if(data[name]) {
				delete data[name];
			}
		});
	}*/
	
	//mxm 2013.11.29
	//如果是autocomplete，删除名字，那对应的Id也不应该保留
	$("#" + formId + " input.ac_input").each(function() {
		var name = $(this).attr("name");
		if(!name) {
			return;
		}
		if(data[name]==null || data[name]=="" || data[name]==undefined) {
			var id = name.substring(0,(name.length-4))+'Id';
			delete data[id];
		}
	});
	
	return data;
}

// $('#form').serializeArray()的数据[{name: a, value: b}, {name: "c[]", value: d}]
// 转成{a:b}
function formArrDataToJson(arrData) {
	var datas = {};
	var arrObj= {}; // {a:[1, 2], b:[2, 3]};
	for(var i in arrData) {
		var attr = arrData[i].name;
		var value = arrData[i].value;
		// 判断是否是a[]形式
		if(attr.substring(attr.length-2, attr.length) == '[]') {
			attr = attr.substring(0, attr.length-2);
			if(arrObj[attr] == undefined) {
				arrObj[attr] = [value];
			} else {
				arrObj[attr].push(value);
			}
			continue;
		}
		
		datas[attr] = value;
	}
	
	return $.extend(datas, arrObj);
}

// 将serialize的的form值转成json
function formSerializeDataToJson(formSerializeData) {
	var arr = formSerializeData.split("&");
	var datas = {};
	var arrObj= {}; // {a:[1, 2], b:[2, 3]};
	for(var i = 0; i < arr.length; ++i) {
		var each = arr[i].split("=");
		var attr = decodeURI(each[0]);
		var value = decodeURI(each[1]);
		// 判断是否是a[]形式
		if(attr.substring(attr.length-2, attr.length) == '[]') {
			attr = attr.substring(0, attr.length-2);
			if(arrObj[attr] == undefined) {
				arrObj[attr] = [value];
			} else {
				arrObj[attr].push(value);
			}
			continue;
		}
		datas[attr] = value;
	}
	
	return $.extend(datas, arrObj);
}

/**
 * 选择radio
 * TODO 有问题, 只能处理一批只有两个的情况
 * @param names Array [] radio名
 * @param source Object {} 值
 */
function setRadios(names, source, parentId) {
	for(var i = 0; i < names.length; ++i) {
		setRadio(names[i], source[names[i]], parentId);
	}
}
// jquery 选择radio
// $("#rdo1").attr("checked","checked");
// $("#rdo1").removeAttr("checked");
function setRadio(name, value, parentId) {
	var index = value ? 0 : 1; // 是在前, 否在后
	id = 'input[name="' + name + '"]';
	if(parentId) {
		id = '#' + parentId + ' ' + id;
	}
	$(id).eq(index).attr("checked", "checked");
}

/**
 * 设置input, textarea值
 * @param names
 * @param source
 */
function setInputTextareas(ids, values, parentId) {
	for(var i = 0; i < ids.length; ++i) {
		setInputTextarea(ids[i], values[ids[i]], parentId);
	}
}
function setInputTextarea(id, value, parentId) {
	var id = '#' + id;
	if(parentId) {
		id = '#' + parentId + ' ' + id;
	}
	$(id).val(value);
}

/**
 * 设置comboboxs值, 选中.
 * @param names
 * @param values
 */
function setComboboxs(names, values, parentId) {
	for(var i = 0; i < names.length; ++i) {
		setCombobox(names[i], values[names[i]], parentId);
	}
}
function setCombobox(name, value, parentId) {
	if(value == null || value == undefined) {
		return;
	}
	
	var id = '#' + name;
	if(parentId) {
		id = '#' + parentId + ' ' + id;
	}
	$(id).val(value);
	// 等待easyui渲染完闭
	setTimeout(function() {
		$(id).combobox();
		$(id).combobox('setValue', value);	
	}, 0);
}

function getCombobox(name, parentId) {
	if(name == null || name == undefined) {
		return;
	}
	var id = '#' + name;
	if(parentId) {
		id = '#' + parentId + ' ' + id;
	}
	var value=$(id).combobox('getValue');
	if(value){
		return value;
	}else{
		return $(id).val();
	}
}

/**
 * 在父parentId下所有值为value的checkbox选中
 * @param parentId String 父ID
 * @param values Array [vlaue1, value2]
 */
function setCheckboxs(values, parentId) {
	if(!values) return;
	for(var i = 0; i < values.length; ++i) {
		setCheckbox(values[i], parentId);
	}
}
function setCheckbox(value, parentId) {
	id = 'input[value="' + value + '"]';
	if(parentId) {
		id = '#' + parentId + ' ' + id;
	}
	$(id).attr('checked', 'checked');
}

/**
 * 展示信息
 * @Param ids <span id="id1"></span> id集合
 * @param sourceData 数据源, sourceData[id]
 * @Param parent 父对象, for精确查询#id
 */
function showInfos(ids, sourceData, parentId) {
	if(!ids || !sourceData) return;
	for(var i = 0; i < ids.length; ++i) {
		showInfo(ids[i], sourceData, parentId);
	}
}
function showInfo(id, sourceData, parentId) {
	var id2 = '#' + id;
	if(parentId) {
		id2 = '#' + parentId + ' ' + id2;
	}
	$(id2).html(sourceData[id]);
}

//------------
// messager 
//------------

/**
 * 提示框
 */
function msgAlert(title, msg, icon, confirmFunc) {
	if(icon == undefined || icon == "") icon = "info";
	title = getMsgTitleText(title);
	msg = getMsgText(msg);
	$.messager.alert(title, msg, icon, confirmFunc);
}

/**
 * 提示, 自己输入信息, 不拿国际化
 */
function msgAlertSelf(title, msg, icon, confirmFunc) {
	if(icon == undefined || icon == "") icon = "info";
	$.messager.alert(getMsgTitleText(title), msg, icon, confirmFunc);
}

/**
 * 确认?
 * @param title
 * @param msg
 * @param yesFunc
 * @param noFunc
 */
function msgConfirm(title, msg, yesFunc, noFunc) {
	title = getMsgTitleText(title);
	msg = getMsgText(msg);
	
	$.messager.confirm(title, msg, function(r) {
		if(r) {
			if(typeof yesFunc == "function") yesFunc();
		} else {
			if(typeof noFunc == "function") noFunc();
		}
	});
}


//-----------
// 工具
//-----------

// 判断后台处理的结果是否为true
function _checkRetStatus(ret) {
	if(!ret) {
		//此处为地图专门处理
		return true;
	}
	
	if(typeof ret == "object" && (ret.success == true || typeof ret.success == 'undefined')) {
		return true;
	}
	
	return false;
}

// 显示进度条
function showProgress(title, msg) {
	var param = {};
	if(title != undefined) {
		param.title = getMsgTitleText(title);
	}
	if(msg != undefined) {
		param.msg = getMsgText(msg)
	}
	$.messager.progress(param);
}
// 关闭进度条
function hideProgress() {
	$.messager.progress('close'); // 关闭进度条
}

/**
 * 是否无权限, 登录(session过期)
 * @param ret
 * @returns {Boolean}
 */
function _noAuth(ret) {
	if(typeof ret == "object") {
		// 没有权限
		if(ret.objectId == "noAuth") {
			msgAlert("info", "noAccess", "info", function() {
			});
			return true;
		// 没有登录, session过时
		} else if(ret.objectId == "sessionTimeOut" || ret.objectId == "otherOneLogined") {
			msgAlert("info", ret.objectId, "info", function() {
				location.href = G_PATH + "/app/login.jsp";
			});
			// 清空定时
			if(dk.activeUserCountInterval) {
				clearInterval(dk.activeUserCountInterval);
			}
			if(dk.notReadMsgInterval) {
				clearInterval(dk.notReadMsgInterval);
			}
			return true;
		}
	}
	return false;
}

/**
 * 得到异常信息
 * @param ret
 * @returns
 */
function _getExceptionMsgs(ret) {
	var msg = "";
	//var msg = getMsgText('opFailure');
	if(ret.exceptionList) {
		for(var i in ret.exceptionList) {
			if(ret.exceptionList[i]) {
				msg += ret.exceptionList[i] + "<br />";
			}
		}
	}
	if(!msg) {
		msg = getMsgText('opFailure');
	}
	return msg;
}

/**
 * 得到异常信息
 * @param ret
 * @returns
 */
function getExceptionMsg(ret) {
	var msg = "";
	if(ret.exceptionList) {
		for(var i in ret.exceptionList) {
			if(ret.exceptionList[i]) {
				msg += ret.exceptionList[i] + "<br />";
			}
		}
	}
	return msg;
}

/**
 * 异常错误! 供datagrid中调用
 * @param ret
 * @returns bool 有错误,异常返回true
 */
function exceptionError(ret) {
	// no auth
	if(_noAuth(ret)) {
		//...
		return true;
	} else {
		// processStatus有错误, 输出异常信息
		if(typeof ret == "object" && typeof ret.success != "undefined" && !ret.success) {
			msgAlertSelf('info', _getExceptionMsgs(ret), 'info');	
			return true;
		} else if(typeof ret != "object" ) {
			msgAlert("info", "opFailure", "info");
			return true;
		}
	}
	
	return false;
}

// ajax请求返回结果后的操作
// 用于ajaxGet(), ajaxPost(), formSubmit()
function _ajaxCallback(ret, successFunc, failureFunc) {
	if(_checkRetStatus(ret)) {
		if(typeof successFunc == 'function') successFunc(ret);
		else msgAlert('info', 'opSuccess', 'info');
	} else {
		// no auth
		if(_noAuth(ret)) {
			//...
		} else {
			if(typeof failureFunc == 'function') {
				failureFunc(ret);
			} else {
				// processStatus有错误, 输出异常信息
				if(typeof ret == "object") {
					msgAlertSelf('info', _getExceptionMsgs(ret), 'info');	
				} else {
					msgAlert("info", "opFailure", "info");
				}
			}
		}
	}
}

function _ajax(type, url, param, successFunc, failureFunc, hasProgress, async) {
	if(typeof async == "undefined") {
		async = true;
	} else {
		async = false;
	}
	$.ajax({
		type: type,
		url: url,
		data: param,
		async: async, // 是否异步
		success: function(ret) {
			_ajaxCallback(ret, successFunc, failureFunc);
			if(hasProgress) hideProgress();
		},
		error: function(ret) {
			_ajaxCallback(ret, successFunc, failureFunc);
			if(hasProgress) hideProgress();
		}
	});
}

/**
 * 发送ajax get请求
 * @param url
 * @param param
 * @param successFunc
 * @param failureFunc
 * @param hasProgress
 * @param async 是否异步
 * @returns
 */
function ajaxGet(url, param, successFunc, failureFunc, hasProgress, async) {
	if(hasProgress == undefined) hasProgress = true;
	if(hasProgress) showProgress('info', 'processing');
	_ajax("GET", url, param, successFunc, failureFunc, hasProgress, async);
}

/**
 * 发送post请求
 * @param url
 * @param param
 * @param successFunc
 * @param failureFunc
 * @param hasProgress
 * @param async 是否异步, 默认为true
 * @returns
 */
function ajaxPost(url, param, successFunc, failureFunc, hasProgress, async) {
	if(hasProgress == undefined) hasProgress = true;
	if(hasProgress) showProgress('info', 'processing');
	_ajax("POST", url, param, successFunc, failureFunc, hasProgress, async);
}

/**
 * 传送Json数据
 * @param url
 * @param param
 * @param successFunc
 * @param failureFunc
 * @param hasProgress
 * @param async
 * @returns
 */
function ajaxPostJson(url, param, successFunc, failureFunc, hasProgress, async) {
	if(hasProgress == undefined) hasProgress = true;
	if(hasProgress) showProgress('info', 'processing');
	// 默认是异步的
	if(typeof async == "undefined") {
		async = true;
	} else {
		async = false;
	}
	$.ajax({
	    url : url,
	    type : "POST",
	    contentType: "application/json; charset=utf-8",
	    datatype: "json",
	    async: async,
	    data : JSON.stringify(param),
	    success : function(data, stats) {
	    	_ajaxCallback(data, successFunc, failureFunc);
	    	if(hasProgress) hideProgress();
	    },
		error: function(ret) {
			_ajaxCallback(ret, successFunc, failureFunc);
			if(hasProgress) hideProgress();
		}
	});
}

/**
 * 隐藏按钮
 * @param btns Object | string 按钮IDs
 */
function hideBtn(btns) {
	if(typeof btns == "object") {
		for(var i in btns) {
			$('#' + btns[i]).css("display", "none");
		}
	} else {
		$('#' + btns).css("display", "none");
	}
}

/**
 * 显示按钮
 * @param btns Object | string 按钮IDs
 */
function showBtn(btns) {
	if(typeof btns == "object") {
		for(var i in btns) {
			$('#' + btns[i]).css("display", "inline-block");
		}
	} else {
		$('#' + btns).css("display", "inline-block");
	}
}

// 关闭按钮
function disableBtn(btns) {
	if(typeof btns == "object") {
		for(var i in btns) {
			$('#' + btns[i]).linkbutton({disabled: true});
		}
	} else {
		$('#' + btns).linkbutton({disabled: true});
	}
}

function enableBtn(btns) {
	if(typeof btns == "object") {
		for(var i in btns) {
			$('#' + btns[i]).linkbutton({disabled: false});
		}
	} else {
		$('#' + btns).linkbutton({disabled: false});
	}
}

/**
 * 
 * country, region, city, town多级联动
 * 
 */
// function areaCascade(countryWebId, regionWebId, cityWebId, townWebId, defaultCountryId, defaultRegionId, defaultCityId, defaultTownId) {
function areaCascade(config) {
	var areaData = {regions: {}, cities: {}, towns: {}}; // 缓存
	
	// 配置处理
	
	// 默认中国
	if(config.defaultCountryId == undefined) {
		config.defaultCountryId = G_CONFIG['defaultCountryId'];
	}
	var countryWebId = (!config.countryWebId) ? false : '#' + config.countryWebId;
	var regionWebId = (!config.regionWebId) ? false : '#' + config.regionWebId;
	var cityWebId = (!config.cityWebId) ? false : '#' + config.cityWebId;
	var townWebId = (!config.townWebId) ? false : '#' + config.townWebId;
	
	function clearRegion() {
		if(regionWebId && $(regionWebId)) {
			$(regionWebId).combobox('loadData', false);
			$(regionWebId).combobox('clear');
		}
	}

	function clearCity() {
		if(cityWebId && $(cityWebId)) {
			$(cityWebId).combobox('loadData', false);
			$(cityWebId).combobox('clear');
		}
	}

	function clearTown() {
		if(townWebId && $(townWebId)) {
			$(townWebId).combobox({data: []});
			$(townWebId).combobox('clear');
		}
	}
	
	// 保存缓存
	// 以后地址有用
	// type: countries, regions, cities, towns
	function saveCache(rows, type) {
		if(!rows) {
			return;
		}
		for(var i = 0; i < rows.length; ++i) {
			dk.cache[type][rows[i].typeId] = rows[i].shortDesc;
		}
	}
	
	// 国家
	if(countryWebId && $(countryWebId)) {
		$(countryWebId).combobox({
			valueField: 'typeId',
			textField: 'shortDesc',
			required: false,
			url: G_PATH + '/geo/getCountries', // 首先载入国家, 通过url, onLoadSuccess有问题
			onLoadSuccess: function(ret) {
				// $(this).combobox('select', config.defaultCountryId); // 选中中国
				// 选中第一行记录
				if(ret && ret.length > 0) {
					$(this).combobox('select', ret[0].typeId);
				}
				
				saveCache(ret, "countries");
			},
			onSelect: function() {
				var countryId = $(this).combobox('getValue');
				
				// 载入省
				clearRegion();
				clearCity();
				clearTown();
				if(regionWebId && $(regionWebId)) {
					if(areaData['regions'][countryId]) {
						$(regionWebId).combobox('loadData', areaData['regions'][countryId]);
					} else {
						$.get(G_PATH + '/geo/getRegions', {countryId: countryId}, function(ret) {
							if(ret.length) {
								areaData['regions'][countryId] = ret;
								$(regionWebId).combobox('loadData', ret);
								
								saveCache(ret, "regions");
							}
							
							// 第一次初始化时选中默认的省. 第一次肯定没有默认的, 所以, 必须要从数据库中, 所以
							// 肯定会执行!
							if(typeof config.defaultRegionId != 'undefined') {
								// alert("region " + ret);
								$(regionWebId).combobox('select', config.defaultRegionId);
								config.defaultRegionId = null;
							}
						});
					}
				}
			}
		})
	}
	
	// 选中中国
	// _selectArea(config.countryWebId, config.defaultCountryId);
	
	// 省
	if(regionWebId && $(regionWebId)) {
		$(regionWebId).combobox({
			valueField: 'typeId',
			textField: 'shortDesc',
			required: false,
			onSelect: function() {
				var regionId = $(this).combobox('getValue');
				
				// 载入市
				clearCity();
				clearTown();
				if(cityWebId && $(cityWebId)) {
					if(areaData['cities'][regionId]) {
						$(cityWebId).combobox('loadData', areaData['cities'][regionId]);
						// 如果市只有一行记录, 则默认选中
						if(areaData['cities'][regionId].length == 1) {
							$(cityWebId).combobox('select', areaData['cities'][regionId][0].typeId);
						}
					} else {
						$.get(G_PATH + '/geo/getCities', {regionId: regionId}, function(ret) {
							if(ret.length) {
								areaData['cities'][regionId] = ret;
								$(cityWebId).combobox('loadData', ret);
								
								// 如果市只有一行记录, 则默认选中
								if(ret && typeof ret == "object" && ret.length == 1) {
									$(cityWebId).combobox('select', ret[0].typeId);
								}
								
								saveCache(ret, "cities");
							}
							
							// 第一次初始化时选中默认的市
							if(typeof config.defaultCityId != 'undefined') {
								$(cityWebId).combobox('select', config.defaultCityId);
								config.defaultCityId = null;
							}
						});
					}
				}
			}
		});
	}
	
	// 市
	if(cityWebId && $(cityWebId)) {
		$(cityWebId).combobox({
			valueField: 'typeId',
			textField: 'shortDesc',
			required: false,
			onSelect: function() {
				var cityId = $(this).combobox('getValue');
				// 载入县
				clearTown();
				if(townWebId && $(townWebId)) {
					if(areaData['towns'][cityId]) {
						$(townWebId).combobox('loadData', areaData['towns'][cityId]);
					} else {
						$.get(G_PATH + '/geo/getTowns?cityId=' + cityId, {}, function(ret) {
							if(ret.length) {
								areaData['towns'][cityId] = ret;
								$(townWebId).combobox('loadData', ret);
							}
							
							// 第一次初始化时选中默认的镇
							if(typeof config.defaultTownId != 'undefined') {
								$(townWebId).combobox('select', config.defaultTownId);
								config.defaultTownId = null;
							}
							
							saveCache(ret, "towns");
						})
					}
				}
			}
		});
	}

	// 县
	if(townWebId && $(townWebId)) {
		$(townWebId).combobox({
			valueField: 'typeId',
			textField: 'shortDesc',
			required: false,
			onSelect: function() {
			}
		});
	}
}

// 选择联动
// config与areaCascade的config一至!// 选择区域
// 因为ajax的原因, 所以要定时循环判断
function _selectArea(areaId, value) {
	if(areaId != undefined && $('#' + areaId)) {
		areaObj = $('#' + areaId);
		var loadedData = areaObj.combobox('getData');
		if(loadedData != "") {
			// 判断是否有该值, 若没该值, 则不选择, 不然会显示其它值
			for(var i in loadedData) {
				if(loadedData[i].typeId == value) {
					areaObj.combobox('select', value);
					return;
				}
			}
			// 执行到此, 没有值
			return;
		}
		setTimeout(function() {_selectArea(areaId, value);}, 10);
	}
}
function selectArea(config) {
	if(config.defaultCountryId != undefined) {
		_selectArea(config.countryWebId, config.defaultCountryId);
		if(config.defaultRegionId != undefined) {
			_selectArea(config.regionWebId, config.defaultRegionId);
			if(config.defaultCityId != undefined) {
				_selectArea(config.cityWebId, config.defaultCityId);
				if(config.defaultTownId != undefined) {
					_selectArea(config.townWebId, config.defaultTownId);
				}
			}
		}
	}
}

/*
function a() {
	<th>国家</th>
	<td><input id="areaCountry" class="easyui-combobox combobox-width80" name="country"></td>
	<th><fmt:message key="geo.region" bundle="${AdminStrings}"/></th>
	<td><input id="areaRegion" class="easyui-combobox combobox-width80" name="region"></td>
	<th><fmt:message key="geo.city" bundle="${AdminStrings}"/></th>
	<td><input id="areaCity" class="easyui-combobox combobox-width80" name="city"></td>
	<th><fmt:message key="geo.town" bundle="${AdminStrings}"/></th>
	<td><input id="areaTown" class="easyui-combobox combobox-width80" name="town"></td>
}
*/

// 类型获取

// type全局变量
// 缓存
dk.types = {
	//id: "desc"	
};
dk.catTypes = {
	//'catName' : {id: "desc"}
}; 

/**
 * 是否有types
 * @param type
 */
function hasTypes(type) {
	if(dk.types[type] != undefined) {
		return true;
	}
	return false;
}

/**
 * 通过type获取类型, type可取 weightUoms
 * @Param type String 类型名
 * @Param callback function 异步获取types的回调方法
 * @return Object
 */
function getTypes(type, callback) {
	// 有回调方法，异步加载
	if(typeof callback == "function") {
		if(dk.types[type] != undefined) {
			callback(dk.types[type]);
			return;
		}
		ajaxGet(G_PATH + '/type/getTypes', {type: type}, function(types) {
			setTypes(type, types);
			callback(types);
		}, '', false, true);
		
	} else {
		if(dk.types[type] != undefined) return dk.types[type];
		
		jQuery.ajax({
	         type: 'GET',
	         url: G_PATH + '/type/getTypes',
	         data: {type: type},
	         async: false, // false代表只有在等待ajax执行完毕后才执行
	         success: function(ret) {
	        	 thisTypes = ret;
	         }
	     });
		
		setTypes(type, thisTypes);
		
		return thisTypes;
	}
}

/**
 * 得到type的字段描述
 * @Param type String 类型名
 * @return String
 */
function getTypeText(type, key) {
	var types = getTypes(type);
	if(types) return types[key];
	return '';
}

/**
 * 为某类设置types
 * @param type
 * @param types
 */
function setTypes(type, types) {
	dk.types[type] = types;
}


/**
 * 得到复杂类型下某分类的types
 * @param cat catName
 * @param catId 所选择的catId
 * @returns
 */
function getCatTypes(cat, catId) {
	if(dk.catTypes[cat] != undefined && dk.catTypes[cat][catId] != undefined) {
		return dk.catTypes[cat][catId];
	}
	var thisTypes;
	jQuery.ajax({
         type: 'GET',
         url: G_PATH + '/type/getCatTypes',
         data: {cat: cat, catId: catId},
         async: false, // false代表只有在等待ajax执行完毕后才执行
         success: function(ret) { // {id: "", id: ""}
        	 thisTypes = ret;
         }
    });
	
	// 缓存
	if(dk.catTypes[cat] == undefined) {
		dk.catTypes[cat] = {};
	}
	dk.catTypes[cat][catId] = thisTypes;
	
	return thisTypes;
}

/**
 * type的32位字符串是否等于 desc
 * @param type
 * @param key
 * @param desc
 * @returns {Boolean}
 */
function typeEqualTo(type, desc, targetDesc) {
	if(!type || !desc || !targetDesc) return false;
	if(!G_TYPES || !G_TYPES[type] || !G_TYPES[type][desc]) return false;
	if(G_TYPES[type][desc] == targetDesc) return true;
	return false;
}

/**
 * 先得到combobox Text，再判断是否相等
 * @param type
 * @param comboboxId
 * @param desc
 * @returns {Boolean}
 */
function comboboxTypeEqualTo(type, comboboxId, targetDesc) {
	if(!type || !comboboxId || !targetDesc) return false;
	var desc = $('#' + comboboxId).combobox('getText');

	if(!desc || !G_TYPES || !G_TYPES[type] || !G_TYPES[type][desc]) return false;
	if(G_TYPES[type][desc] == targetDesc) return true;
	return false;
}

/**
 * 得到某type，某key的desc
 * @param type
 * @param key
 * @returns
 */
function getTypeDesc(type, desc) {
	if(!type || !desc || !G_TYPES || !G_TYPES[type]) return false;
	return G_TYPES['type'][desc];
}

//----------------
// get url
//----------------

function getPlatformPath(url) {
	return G_PATH + '/getJsp?jsp=platform/' + url;
}

function getOrgPath(url) {
	return G_PATH + '/getJsp?jsp=org/' + url;
}

function getUserPath(url) {
	return G_PATH + '/getJsp?jsp=user/' + url;
}

function getCommonPath(url) {
	return G_PATH + '/getJsp?jsp=common/' + url;
}

//.. web重构, 加载jsp页面时会自动加载其相应的js(如果有)

function getPlatformJspJsPath(url) {
	return G_PATH + '/getJspJs?jsp=platform/' + url;
}

function getOrgJspJsPath(url) {
	return G_PATH + '/getJspJs?jsp=org/' + url;
}

function getUserJspJsPath(url) {
	return G_PATH + '/getJspJs?jsp=user/' + url;
}

function getCommonJspJsPath(url) {
	return G_PATH + '/getJspJs?jsp=common/' + url;
}

/**
easyui validate
*/
$.extend($.fn.validatebox.defaults.rules, {  
    CHS: {
        validator: function (value, param) {  
            return /^[\u0391-\uFFE5]+$/.test(value);  
        },  
        message: '请输入汉字'  
    },
    ZIP: {  
        validator: function (value, param) {  
            return /^[1-9]\d{5}$/.test(value);  
        }, 
        message: '邮政编码不存在'  
    },
    QQ: {  
        validator: function (value, param) {  
            return /^[1-9]\d{4,10}$/.test(value);  
        },  
        message: 'QQ号码不正确'  
    },  
    mobile: {
        validator: function (value, param) {  
            return /^13\d{9}$/g.test(value) || /^14\d{9}$/g.test(value) || (/^15[0-35-9]\d{8}$/g.test(value)) || (/^18[05-9]\d{8}$/g.test(value));
        },
        message: '手机号码不正确'  
    },
    mobileOrPhone: {
        validator: function (value, param) {  
            return /^([0-9]|\-)*$/.test(value);
        },  
        message: '号码不正确'  
    },  
    loginName: {  
        validator: function (value, param) {  
            return /^[\u0391-\uFFE5\w]+$/.test(value);  
        },
        message: '登录名称只允许汉字、英文字母、数字及下划线。'  
    },  
    safePassword: {  
        validator: function (value, param) {  
        	if(value.length < 6) {
        		return false;
        	}
        	if(/[a-zA-Z]+/.test(value) && /[0-9]+/.test(value)/* && /\W+\D+/.test(string)*/) {
        		return true;
        	}
        },
        message: '密码由字母和数字组成，至少6位'  
    },  
    equalTo: {
        validator: function (value, param) {  
            return value == $(param[0]).val();  
        },
        message: '两次输入的字符不一至'  
    },  
    number: {
        validator: function (value, param) { 
            return /^\d*\.?\d*$/.test(value);  
        },
        message: '请输入数字'  
    },
    integer: {
    	validator: function (value, param) { 
            return /^\d+$/.test(value);  
        },
        message: '请输入整数'  
    },
    time: {
        validator: function (value, param) { 
            if(/^\d+\:\d+$/.test(value)) {
            	var a = value.split(":");
            	if(a[0] > 23 || a[1] > 59) {
            		return false;
            	}
            	return true;
			}
            return false;
        },
        message: '请输入时间'  
    }, 
    email: {  
        validator: function (value, param) { 
        	var myreg = /^([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+\.[a-zA-Z]{2,3}$/;
        	return myreg.test(value);
        },
        message: '请输入正确邮箱'  
    },
    postalCode: {
        validator: function (value, param) {  
            if(/^\d+$/.test(value)) {
            	if(value.length == 6) return true;
            }
            return false;
        },
        message: '请输入6位数字'  
    },
    idcard: {  
        validator: function (value, param) {  
            return idCard(value);  
        },  
        message:'请输入正确的身份证号码'  
    }, 
    maxLength: {
        validator: function(value, param){
            return value.length <= param[0];
        },
        message: '最多输入{0}个字符'
    }
});  

// var myDate = new Date();
// var date = myDate.format('yyyy-MM-dd hh:mm:ss'); 
Date.prototype.format = function(format)
{
    var o =
    {
        "M+" : this.getMonth()+1, //month
        "d+" : this.getDate(),    //day
        "h+" : this.getHours(),   //hour
        "m+" : this.getMinutes(), //minute
        "s+" : this.getSeconds(), //second
        "q+" : Math.floor((this.getMonth()+3)/3),  //quarter
        "S" : this.getMilliseconds() //millisecond
    }
    if(/(y+)/.test(format))
    format=format.replace(RegExp.$1,(this.getFullYear()+"").substr(4 - RegExp.$1.length));
    for(var k in o)
    if(new RegExp("("+ k +")").test(format))
    format = format.replace(RegExp.$1,RegExp.$1.length==1 ? o[k] : ("00"+ o[k]).substr((""+ o[k]).length));
    return format;
}

function getCurrentTimestamp() {
	var myDate = new Date();
	return myDate.format('yyyy-MM-dd hh:mm:ss');
}

function getCurrentDate() {
	var myDate = new Date();
	return myDate.format('yyyy-MM-dd');
}

function getTomorrowDate() {
	var myDate = new Date();
	myDate.setDate(myDate.getDate()+1);
	return myDate.format('yyyy-MM-dd');
}

function getEndTimestamp() {
	return "2100-12-31 12:59:59";
}

function getEndDate() {
	return "2100-12-31";
}

/**
 * 是否为空
 * 可判断任意类型，string array
 */
function isEmpty(obj) {
	if(!obj) {
		return true;
	}
	
	if(isArray(obj)) {
		if(obj.length == 0) {
			return true;
		}
	}
	
	return false;
}

/**
 * Enter键触发点击事件
 * @param btnId
 */
function enterPressSubmit(btnId) {
	var theEvent = window.event || arguments.callee.caller.arguments[0];
	if(theEvent.keyCode == 13||theEvent.keyCode == 108) {
		$('#' + btnId).trigger('click');
	}
}

/**
 * Enter键触发点击事件
 * @param btnId
 */
function enterPressSearch(func) {
	var theEvent = window.event || arguments.callee.caller.arguments[0];
	if(theEvent.keyCode == 13||theEvent.keyCode == 108) {
		$('#' + btnId).click(func);
	}
}

// 其它 在哪用到了?

/**
 * 设置单个或者多个下拉列表的默认值，id和value的个数必须一致，即一一对应
 * @param selectIds Array 下拉列表Id的数组，比如['numId','numId2']
 * @param values Array 要设置的值的数组，比如['numId','numId2']
 */
function setSelectValues(selectIds,values){
	if(selectIds.length == 1 && values.length == 1){
		$('#' + selectIds).combobox('select', values);
		
	}else if(selectIds.length>1 && values.length>1 && selectIds.length==values.length){
		for(var i=0;i<selectIds.length;i++){
			$('#' + selectIds[i]).combobox('select', values[i]);
		}
	}
}
/**
 * 验证数字格式是否正确，正确，返回true，反之false
 * @param numberIds Array 输入框Id的数组，比如['numId','numId2']
 */
function regIsNums(numberIds){
	if(numberIds.length == 1){
		var s = $("#" + numberIds).val();
		if(!isNum(s)){
			$("#" + numberIds).css("boreder","10px solid red");
			msgAlert('info','numberNotCorrect','info');		
			return false;
		}
	}else if(numberIds.length > 1){
		for(var i=0;i<numberIds.length;i++){
			var s = $("#" + numberIds[i]).val();
			if(!isNum(s)){
				msgAlert('info','numberNotCorrect','info');
				return false;
			}
		}
	}
	return true;
}

function isNum(num) {
    reg = /^\d+$/;
    reg2 = /^\d+.\d+$/;
    return reg.test(num) || reg2.test(num);
}
/**
 * 验证输入框是否是否为空，不空，返回true，反之false
 * @param inputIds Array 输入框Id的数组，比如['rateTypeId','rateTypeId2']
 */
function checkInputIsNull(inputIds){
	if(inputIds.length == 1){
		var s = $("#" + inputIds).val();
		if(s==null || s==""){
			msgAlert('info','pageNotCorrect','info');		
			return false;
		}
	}else if(inputIds.length > 1){
		for(var i=0;i<inputIds.length;i++){
			var s = $("#" + inputIds[i]).val();
			if(s==null || s==""){
				msgAlert('info','pageNotCorrect','info');
				return false;
			}
		}
	}
	return true;
}
/**
 * 验证select下拉框是否选中
 * @param selectIds Array 下拉框Id的数组，比如['rateTypeId','rateTypeId2']
 */
function checkIsSelected(selectIds){
	if(selectIds.length == 1){
		var s = $("#" + selectIds).combobox('getValue');
		if(s==null || s==""){
			msgAlert('info','pageNotCorrect','info');		
			return false;
		}
	}else if(selectIds.length > 1){
		for(var i=0;i<selectIds.length;i++){
			var s = $("#" + selectIds[i]).combobox('getValue');
			if(s==null || s==""){
				msgAlert('info','pageNotCorrect','info');
				return false;
			}
		}
	}
	return true;
}

// 与具体应用相关

/**
 * 查询站点
 * @param inputId 要显示站点名称的input
 * @param valueId 要存储站点ID的input
 * 
 * 用户添加/修改, 订单添加/修改
 */
function searchLocation(inputId, valueId) {
	areaInputId = inputId;
	areaValueId = valueId;
	var areaName = $('#' + inputId).val();
	
	var config = {
		title: 'locationSearch',
	    width: 820,
	    height: 450,
	    href: G_PATH + '/getJsp?jsp=common/commom_window_area_name_search',
	    onLoad: function() {
	    	initGrid('common_locationSearch', 'main', 'areaSearchGrid', {
	    		fit: true,
	    		singleSelect: true,
				url: G_PATH + '/site/checkLocation',
				queryParams: {areaName: areaName}
	    	});
	    	$("#areaSearchForm #areaName").val(areaName);
	    }
	};

	openWin(config, 'commonSearchLocationWin');
}


/**
 * 查询设备组
 * @param inputId 要显示设备组的input
 * @param valueId 要存储设备组ID的input
 * 
 * 用户添加/修改, 订单添加/修改
 */
var equipTeamInputId;
var equipTeamValueId;
function searchEquipTeam(inputId, valueId) {
	equipTeamInputId = inputId;
	equipTeamValueId = valueId;
	var equipTeamName = $('#' + inputId).val();
	
	var config = {
		title: 'equipTeamSearch',
	    width: 820,
	    height: 450,
	    href: G_PATH + '/getJsp?jsp=common/equip/equip_list_window_search',
	    onLoad: function() {

	    	$("#equipListSearchForm #equipListName").val(equipTeamName);
	    	$('#searchEquipListBtn').click();
	    }
	};

	openWin(config, 'equipTeamSearchWin');
}
/**
 * 合同线路 查询  出发地 目的地  
 */
var areaAndLocationInputId;
var areaAndLocationValueId;
function searchAreaAndLocation(inputId, valueId) {
	areaAndLocationInputId = inputId;
	areaAndLocationValueId = valueId;
	var areaName = $('#' + inputId).val();
	
	var config = {
		title: 'areaSearch',
	    width: 880,
	    height: 550,
	    href: G_PATH + '/getJsp?jsp=common/area/commom_window_area_search',
	    onLoad: function() {
	    	
	    }
	};

	openWin(config, 'areaSearchWin');
}
/**
 * 承运商（外部） 组织 查询
 */
var carrierOrgInputId;
var carrierOrgValueId;
function searchCarrierOrgCommon(inputId, valueId) {
	carrierOrgInputId = inputId;
	carrierOrgValueId = valueId;
	var orgName = $('#' + inputId).val();
	var config = {
		title: 'carrierOrgSearch',
	    width: 900,
	    height: 487,
	    href: G_PATH + '/orgCommon/searchCarrierOrgJsp'
	};

	openWin(config, 'searchCarrierOrgWin');
}

/**
 * 客户 组织 查询
 */
var customerOrgInputId;
var customerOrgValueId;
function searchCustomerOrgCommon(inputId, valueId) {
	customerOrgInputId = inputId;
	customerOrgValueId = valueId;
	var orgName = $('#' + inputId).val();
	var config = {
		title: 'customerOrgSearch',
	    width: 900,
	    height: 487,
	    href: G_PATH + '/orgCommon/searchCustomerOrgJsp'
	};

	openWin(config, 'searchCustomerOrgWin');
}

/**
 * 承运商(内部) 组织 查询
 */
var carrierInnerOrgInputId;
var carrierInnerOrgValueId;
function searchCarrierInnerOrgCommon(inputId, valueId) {
	carrierInnerOrgInputId = inputId;
	carrierInnerOrgValueId = valueId;
	var config = {
		title: 'carrierOrgSearch',
	    width: 900,
	    height: 487,
	    href: G_PATH + '/orgCommon/searchCarrierInnerOrgJsp'
	};

	openWin(config, 'searchCarrierInnerOrgWin');
}

/**
 * 组织分组 查询
 */
var memberListInputId;
var memberListValueId;
function searchMemberListCommon(inputId, valueId) {
	memberListInputId = inputId;
	memberListValueId = valueId;
	var orgName = $('#' + inputId).val();
	
	var config = {
		title: 'searchMemberList',
	    width: 900,
	    height: 487,
	    href:G_PATH + '/getJsp?jsp=common/memberlist/org_layout_center_searchMemberlist'
	};

	openWin(config, 'searchMemberListWin');
}
var orgInputId;
var orgValueId;
function searchOrg(inputId, valueId) {
	orgInputId = inputId;
	orgValueId = valueId;
	var orgName = $('#' + inputId).val();
	
	var config = {
		title: 'orgSearch',
	    width: 900,
	    height: 450,
	    href: G_PATH + '/getJsp?jsp=common/org/org_window_search',
	    onLoad: function() {
	    	initGrid('orgSearch', 'main', 'orgSearchGrid', {
	    		fit: true,
	    		singleSelect: true,
				url: G_PATH + '/orgCommon/searchOrg'
	    	});
	    	//$("#areaSearchForm #areaName").val(areaName);
	    }
	};

	openWin(config, 'orgSearchWin');
}

var fcOrgInputId;
var fcOrgValueId;
function searchFCOrg(inputId, valueId) {
	fcOrgInputId = inputId;
	fcOrgValueId = valueId;
	var orgName = $('#' + inputId).val();
	
	var config = {
		title: 'orgSearch',
	    width: 900,
	    height: 450,
	    href: G_PATH + '/getJsp?jsp=common/org/org_window_searchFcOrg',
	    onLoad: function() {
	    	initGrid('orgSearch', 'main', 'orgSearchGrid', {
	    		fit: true,
	    		singleSelect: true,
				url: G_PATH + '/orgCommon/searchFCOrg'
	    	});
	    	//$("#areaSearchForm #areaName").val(areaName);
	    }
	};

	openWin(config, 'fcOrgSearchWin');
}

/**
 * 通用添加分段费率
 */
var selectedShipmentAttribTypeId;
var selectedUomId;
var breakTemplateCommonId;
var breakTemplateCommonName;
var showBreakTemplate;
function searchBreakTemplateCommon(shipmentAttribTypeId, uomId,breakTemplateName,breakTemplateId,showFunc) {
	selectedShipmentAttribTypeId = shipmentAttribTypeId;
	selectedUomId = uomId;
	breakTemplateCommonName=breakTemplateName;
	breakTemplateCommonId=breakTemplateId;
	showBreakTemplate=showFunc;
	var config = {
		title: 'createBreakTemplate',
	    width: 990,
	    height: 555,
	    href: G_PATH + '/getJsp?jsp=org/breakTemplate/common_breakTemplate_search',
    	onLoad:function(){
	    	initGrid('contract', 'addBreakTemplate', 'breakCostInfo', {
	    		fit:true,
	    		afterLoad: function() {
	    		}
	    	});	
	    }   
	};

	openWin(config, 'searchBreakTemplateCommonWin');
}

var simpleOrgInputId;
var simpleOrgValueId;
function searchOriginSimpleOrg(inputId,valueId) {
	simpleOrgInputId = inputId;
	simpleOrgValueId = valueId;
	
	var winWidth = getWinWidth(9, 950);
	var config = {
		title: 'orgNews',
	    width: winWidth,
	    height: 450,
	    href: G_PATH + '/getJsp?jsp=common/org/user_window_origin_org_simpleInfo_search'
	};
	
	
	openWin(config, 'originOrgSimpleInfoWin');
}

function searchDestSimpleOrg(inputId,valueId) {
	simpleOrgInputId = inputId;
	simpleOrgValueId = valueId;
	
	var winWidth = getWinWidth(9, 850);
	var config = {
		title: 'orgNews',
	    width: winWidth,
	    height: 450,
	    href: G_PATH + '/getJsp?jsp=common/org/user_window_dest_org_simpleInfo_search'
	};
	
	
	openWin(config, 'destOrgSimpleInfoWin');
}

/**
 * 得到状态描述
 * @param activeFlag
 * @returns
 */
function getActiveFlagDesc(activeFlag) {
	return activeFlag ? getText('activeFlagOn') : getText('activeFlagOff');
}

/**
 * 判断浏览器
 */
function judgeNav() {
	var nav = {};
	var userAgent = navigator.userAgent; //取得浏览器的userAgent字符串
	nav.isOpera = userAgent.indexOf("Opera") > -1; //判断是否Opera
	nav.isMaxthon = userAgent.indexOf("Maxthon") > -1 ; //判断是否傲游3.0
	nav.isIE = userAgent.indexOf("compatible") > -1 && userAgent.indexOf("MSIE") > -1 && !nav.isOpera ; //判断是否IE 
	nav.isFF = userAgent.indexOf("Firefox") > -1 ; //判断是否Firefox
	nav.isSafari = userAgent.indexOf("Safari") > -1 && userAgent.indexOf("Chrome") < 1 ; //判断是否Safari
	nav.isChrome = userAgent.indexOf("Chrome") > -1 ; //判断是否Chrome
	
	return nav;
}

/**
 * 得到rows的ids
 * @param rows
 * @param pk
 * @returns {Array}
 */
function getIds(rows, pk) {
	var ids = new Array();
	for( var i = 0; i < rows.length; ++i) {
		ids.push(rows[i][pk]);
	}
	return ids;
}

/**
 * 查看帮助
 * @Depreciated
 */
function viewManual(tag) {
	if(!tag) {
		return;
	}
	var config = {
		title : 'viewManual',
		width : 800,
		height : 600,
		href : G_PATH + "/manual/viewManual?tag=" + tag
	};
	
	openWin(config, 'viewManual');
}

/**
 * 搜索并刷新grid
 * @param formId
 * @param tableId
 * @param url
 */
function searchReloadGrid(formId, tableId, url, extData) {
	if(!tableId) {
		tableId = "mainGrid";
	}
	var data = getFormJsonData(formId);
	
	$.extend(data, extData);
	loadGrid(tableId, data, url);
}



function clearHiddenId(name,id) {
	var value=$("#"+name).val();
	if(!value){
		log(value);
		$("#"+id).val("");
	}
}

// 重构 -----------------

// 是否是数组
function isArray(obj) {  
	return Object.prototype.toString.call(obj) === '[object Array]';   
}

/**
 * 得到深层次值
 * @param row
 * @param name a.b.c
 * @returns
 */
function getDeepValue(row, name) {
	if(!name) {
		return null;
	}
	var nameArr = name.split('.');
	if(!nameArr || nameArr.length == 0) {
		return null;
	}
	/*
	if(nameArr.length == 1) {
		return row[name];
	}
	return _getDeepValue(row, nameArr);
	*/
	
	// 循环版
	var curData = row;
	for(var i = 0; i < nameArr.length; ++i) {
		var key = nameArr[i];
		if(i == nameArr.length - 1) {
			return curData[key];
		}
		if(!curData[key]) {
			return null;
		}
		curData = curData[key];
	}
}

// 深层次解析
function _getDeepValue(row, fieldArr, index) {
	if(!index) {
		index = 0;
	}
	if(row[fieldArr[index]]) {
		if(index == fieldArr.length - 1) {
			// 到最后了
			return row[fieldArr[index]];
		} else {
			// 递归下一位
			return _getDeepValue(row[fieldArr[index]], fieldArr, index+1);
		}
	}
	return null;
}

/**
 * 设置表单值
 * @param formId
 * @param row
 */
function setForm(formId, row) {
	if(!formId || !row) {
		return;
	}
	var parent = "#" + formId + " ";
	// input[type='text']
	$(parent + "input[type='text'] ," + parent + "input[type='hidden']").each(function() {
		var name = $(this).attr("name");
		if(name) {
			$(this).val(getDeepValue(row, name));
		}
	});
	// input[type="radio"]
	var hasSetRaioNames = {}; // 已经设置过的 radio name
	$(parent + "input[type='radio']").each(function() {
		var name = $(this).attr("name");
		if(name && !hasSetRaioNames[name]) {
			var val = getDeepValue(row, name);
			if(typeof val == "boolean") {
				val = val ? 1 : 0;
			}
			$(parent + " input[name='" + name + "'][value='" + val + "']").attr("checked", "checked");
			hasSetRaioNames[name] = true;
		}
	});
	// input[type="checkbox"]
	$(parent + "input[type='checkbox']").each(function() {
		var name = $(this).attr("name");
		if(name) {
			var val = getDeepValue(row, name);
			if(isArray(val)) {
				for(var i = 0; i < val.length; ++i) {
					$(parent + " input[name='" + name + "'][value='" + val[i] + "']").attr("checked", "checked");
				}
			} else {
				$(parent + " input[name='" + name + "'][value='" + val + "']").attr("checked", "checked");
			}
		}
	});
	
	// 以下datebox, datetimebox以后会弃用，全部采用my97
	
	// datebox
	$(parent + "input.easyui-datebox").each(function() {
		var name = $(this).attr("comboname") || $(this).attr("name");
		if(name) {
			var val = getDeepValue(row, name);
			if(val) {
				val = getDate(val);
				var id = $(this).attr("id");
				$("#" + id).datebox({value: val});
			}
		}
	});
	// datetimebox
	$(parent + "input.easyui-datetimebox").each(function() {
		var name = $(this).attr("comboname") || $(this).attr("name");
		if(name) {
			var val = getDeepValue(row, name);
			if(val) {
				var id = $(this).attr("id");
				$("#" + id).datetimebox({value: val});
			}
		}
	});
	
	// my97日期时间控件
	// <dk:date name="a"/>
	$(parent + "input.Wdate").each(function() {
		var name = $(this).attr("name");
		if(name) {
			var val = getDeepValue(row, name);
			if(!val) {
				return;
			}
			var viewType = $(this).attr("viewType");
			if(viewType == "date") {
				val = transeToDate(val);
			} else if(viewType == "time") {
				val = transeToTime(val);
			}
			
			if(val) {
				$(this).val(val);
			}
		}
	});
	
	// textarea
	$(parent + "textarea").each(function() {
		var name = $(this).attr("name");
		if(name) {
			$(this).val(getDeepValue(row, name));
		}
	});
	
	// combobox
	$(parent + "select.easyui-combobox").each(function() {
		var name = $(this).attr("comboname") || $(this).attr("name");
		if(name) {
			setCombobox($(this).attr("id"), getDeepValue(row, name), formId);
		}
	});
	
	// timespinner
	$(parent + "input.easyui-timespinner").each(function() {
		var name = $(this).attr("comboname") || $(this).attr("name");
		if(name) {
			var val = getDeepValue(row, name);
			if(val) {
				var id = $(this).attr("id");
				$("#" + id).timespinner({value: val});
			}
		}
	});
	
/*	// input[type='hidden']
	$(parent + "input[type='hidden'] ") .each(function() {
		var name = $(this).attr("name");
		if(name) {
			$(this).val(getDeepValue(row, name));
		}
	});*/
}

/**
 * 更深层次解析
 * 如果是date，查看属性dbType
 * 
 * needValidateForm 需要验证表单
 */
function getFormDeepJsonData(formId, needValidateForm) {
	if(needValidateForm) {
		if(!$("#" + formId).valid()){
			return false;
		}
	}
	var data = getFormJsonData(formId);
	if(!data) {
		return null;
	}
	return transToDeepJsonData(data);
}
/**
 * 将普通的没有层次的json数据转成有层次的数据
 * data = {'a.b.c': "xx", 'a.b.d': 'kk', 'b.c': 'll'};
 * => data = {a: {b:{c:"xx"}}}
 */
function transToDeepJsonData(data) {
	if(!data) {
		return null;
	}
	function extendData(data, keyArr, val, index) {
		if(!index) {
			index = 0;
		}
		var curKey = keyArr[index];
		if(index == keyArr.length - 1) {
			data[curKey] = val;
			return;
		}
		
		if(!data[curKey]) {
			data[curKey] = {};
		}
		
		extendData(data[curKey], keyArr, val, index + 1);
	}
	
	// data = {'a.b.c': "life", 'a.b.d': "xx"}
	var data2 = {};
	for(var key in data) {
		var keyArr = key.split('.');
		if(!keyArr || keyArr.length <= 0) {
			continue;
		}
		if(keyArr.length == 1) {
			data2[key] = data[key];
		}
		
		extendData(data2, keyArr, data[key], 0);
	}
	return data2;
}

/**
 * 得到date值, date, datetime
 * 会判断dbType转化
 * @param id
 * @returns
 */
function getDateValue(id, parent) {
	var obj;
	if(parent) {
		obj = $("#" + parent + " #" + id);
	} else {
		obj = $("#" + id);
	}
	var val = $(obj).val();
	if(!val) {
		return null;
	}
	var viewType = $(obj).attr("viewType");
	var dbType = $(obj).attr("dbType");
	if(viewType == dbType) {
		return val;
	}
	if(viewType == "date" && dbType == "datetime") {
		return val + " 00:00:00";
	}
}

// 设置日期
function setDateValue(id, val) {
	if(typeof id != "object") {
		id = $("#" + id);
	}
	var viewType = $(id).attr("viewType");
	if(viewType == "date") {
		val = transeToDate(val);
	} else if(viewType == "time") {
		val = transeToTime(val);
	}
	
	if(val) {
		$(id).val(val);
	}
}

// 将2012-12-12 12:12:12分隔出2012-12-12
function transeToDate(dateTime) {
	var dateTimeArr = dateTime.split(' ');
	if(dateTimeArr && dateTimeArr.length > 0) {
		return dateTimeArr[0];
	}
	return dateTime;
}

// 将2012-12-12 12:12:12分隔出12:12:12
function transeToTime(dateTime) {
	// 如果是datetime格式的
	var dateTimeArr = dateTime.split(' ');
	if(dateTimeArr && dateTimeArr.length > 1) {
		return dateTimeArr[1];
	}
	
	// 如果是12:12:12格式的
	return dateTime;
}

// 设置信息 idDesc没有controller还未测
// <span id="a.b" class="info" data-options="type:'date', source:'a.b'"/></span>
// <span id="a.dddb" class="info" data-options="type:'constant', source:'d', typeValue:'snsRoles'"/></span>
// <span id="c" class="info" data-options="type:'datetime', source:'modelName'"/></span>
// <span id="d" class="info" data-options="type:'idDesc', typeValue:'org/getOrgNames'"/></span>
// <span id="resizable" class="info" data-options="type:'boolean', typeValue:'yes,no'"/></span>
// <span id="f" class="info" data-options="type:'combine', typeValue:'{modelName}-{modelType}-{d}'"/></span>
// <span id="f" class="info" data-options="type:'address', typeValue:'{country}-{region}-{city}', source:'addressPartyInfo'"/></span>
function setInfo(parentId, row) {
	if(!parentId) {
		return;
	}
	
	// 找到所有的
	$("#" + parentId + " .info").each(function() {
		setValueOrInfo(this, row);
	});
}

// 设置input, span信息
function setValueOrInfo(idOrgObjOrArr, row) {
	var obj = idOrgObjOrArr;
	if(isArray(obj)) {
		for(var i = 0; i < obj.length; ++i) {
			setValueOrInfo(obj[i], row);
		}
		return;
	} else if(typeof obj == "object") {
		obj = $(obj);
	} else {
		// <span id=""> <input id="" />
		obj = $("#" + obj);
	}
	
	if($(obj).is("input") || $(obj).is("textarea")) {
		var setValue = function(value) {
			obj.val(value);
		}
	} else {
		var setValue = function(value) {
			obj.html(value);
		}
	}
	
	var id = obj.attr("id");
	
	var dataOptions = obj.attr("data-options");
	if(!dataOptions) {
		dataOptions = {};
	} else {
		try {
			dataOptions = eval("({" + dataOptions + "})");
		} catch(e) {
			alert(dataOptions + " 有误！无法解析成Json, 少了单引号?");
			return;
		}
	}
	
	// 数据源
	var source = dataOptions['source'] || obj.attr("name") || id;
	if(!source) {
		return;
	}
	
	// 组合描述
	function combineDescs(types, ids) {
		if(!types || !ids) {
			return;
		}
		var descs = "";
		for(var i = 0; i < ids.length; ++i) {
			descs += types[ids[i]];
			if(i != ids.length - 1) {
				descs += G_CONFIG['split'];
			}
		}
		return descs;
	}
	
	var type = dataOptions.type;
	var typeValue = dataOptions.typeValue;
	
	// 只有conbine类型不需要数据源
	var val = getDeepValue(row, source);
	if(val == undefined && type != "combine") {
		return;
	}
	if(!type) {
		setValue(val);
	} else if(type == 'datetime') {
		setValue(val);
	} else if(type == 'date') {
		setValue(transeToDate(val));
	} else if(type == 'time') {
		setValue(transeToTime(val));
	} else if(type == 'boolean') {
		var options = typeValue.split(","); // yes,no
		if(!options || options.length <= 1) {
			return;
		}
		setValue(val ? getText($.trim(options[0])) : getText($.trim(options[1])));
	} else if(type == 'constant') {
		if(!typeValue) {
			return;
		}
		// 常量，需要异步
		// 异步加载
		getTypes(typeValue, function(types) {
			if(isArray(val)) {
				setValue(combineDescs(types, val));
			} else {
				setValue(types[val]);
			}
		});
	} else if(type == 'idDesc') {
		if(!typeValue) {
			return;
		}
		
		// id转decs, 需要异步，先把所有的controller收集
		// {controller: "", target: that, ids: []}
		if(dk.cache[typeValue] && dk.cache[typeValue][val]) {
			setValue(dk.cache[typeValue][val]);
		} else {
			ajaxPost(G_PATH + '/' + typeValue, {ids: [val]}, function(ret) {
				// 数据缓存起来
				if(!dk.cache[typeValue]) {
					dk.cache[typeValue] = {}
				}
				$.extend(dk.cache[typeValue], ret); 
				setValue(dk.cache[typeValue][val]);
			}, '', false);
		}
	} else if(type == 'combine') {
		if(!typeValue) {
			return;
		}
		
		// 组合类型 {a.b}-{c}-{d}
		var fieldArr = typeValue.match(/\{.+?\}/g);
		if(!fieldArr) {
			return;
		}
		
		var fieldArr2 = [];
		for(var i = 0; i < fieldArr.length; ++i) {
			if(fieldArr[i].length < 3) { // 至少{a}
				continue;
			}
			fieldArr2.push(fieldArr[i].substr(1, fieldArr[i].length - 2));
		}
		if(fieldArr2.length == 0) {
			return;
		}
		
		function contactValue(row, typeValue, fieldArr2) {
			for(var i = 0; i < fieldArr2.length; ++i) {
				var str = getDeepValue(row, fieldArr2[i]);
				if(!str) {
					str = "";
				}
				// 替换
				typeValue = typeValue.replace("{" + fieldArr2[i] + "}", str);
			}
			return typeValue;
		}
		
		setValue(contactValue(row, typeValue, fieldArr2));
		
	} else if(type = "address") {
		// 地址组合
		if(!typeValue) {
			return;
		}
		var addressInfo = getDeepValue(row, source);
		if(!addressInfo) {
			return;
		}
		var address = "";
		var ajaxData = {};
		var needAjax = false;
		var names = {
			country: '',
			region: '',
			city: '',
			town: ''
		}
		// 收集要发送请求的id
		if(addressInfo.countryId && typeValue.indexOf("country")) {
			if(dk.cache.countries[addressInfo.countryId]) {
				names['country'] = dk.cache.countries[addressInfo.countryId];
			} else {
				ajaxData['countryId'] = addressInfo.countryId;
				needAjax = true;
			}
		}
		if(addressInfo.regionId && typeValue.indexOf("region")) {
			if(dk.cache.regions[addressInfo.regionId]) {
				names['region'] = dk.cache.regions[addressInfo.regionId];
			} else {
				ajaxData['regionId'] = addressInfo.regionId;
				needAjax = true;
			}
		}
		if(addressInfo.cityId && typeValue.indexOf("city")) {
			if(dk.cache.cities[addressInfo.cityId]) {
				names['city'] = dk.cache.cities[addressInfo.cityId];
			} else {
				ajaxData['cityId'] = addressInfo.cityId;
				needAjax = true;
			}
		}
		if(addressInfo.townId && typeValue.indexOf("town")) {
			if(dk.cache.towns[addressInfo.townId]) {
				names['town'] = dk.cache.towns[addressInfo.townId];
			} else {
				ajaxData['townId'] = addressInfo.townId;
				needAjax = true;
			}
		}
		
		// ajax取
		if(needAjax) {
			ajaxGet(G_PATH + "/desc/getAddress", ajaxData, function(ret) {
				$.extend(names, ret);
				// 保存到cache中
				if(ajaxData['countryId']) {
					dk.cache.countries[ajaxData['countryId']] = names['country'];
				}
				if(ajaxData['regionId']) {
					dk.cache.regions[ajaxData['regionId']] = names['region'];
				}
				if(ajaxData['cityId']) {
					dk.cache.cities[ajaxData['cityId']] = names['city'];
				}
				if(ajaxData['townId']) {
					dk.cache.towns[ajaxData['townId']] = names['town'];
				}
				renderAddress();
			}, '', false);
		} else {
			renderAddress();
		}
		
		function renderAddress() {
			if(typeValue.indexOf("country")) {
				typeValue = typeValue.replace("{country}", names.country || "");
			}
			if(typeValue.indexOf("region")) {
				typeValue = typeValue.replace("{region}",  names.region || "");
			}
			if(typeValue.indexOf("city")) {
				typeValue = typeValue.replace("{city}",  names.city || "");
			}
			if(typeValue.indexOf("town")) {
				typeValue = typeValue.replace("{town}",  names.town || "");
			}
			for(var j = 1; j <= 4; ++j) {
				var street = "street" + j;
				if(typeValue.indexOf(street)) {
					var desc = addressInfo[street] || "";
					typeValue = typeValue.replace("{" + street + "}", desc);
				}
			}
			
			setValue(typeValue);
		}
	}
}


/**
 * 自动将表单的数据解析成有层次的Json数据提交
 * @param formId
 * @param formDataFunc
 * @param successFunc
 * @param failureFunc
 * @param hasProgress
 */
function formSubmitDeepJson(formId, url, formDataFunc, successFunc, failureFunc, hasProgress) {
	hasProgress = hasProgress == undefined ? true : hasProgress;
	if(hasProgress) {
		showProgress('info', 'processing');
	}
	var formId2 = $('#' + formId);
	initValidator(formId);
	var isValid = valid(formId); // 验证数据, 用jQuery validation
	
	if(!isValid) {
		hideProgress();
		// msgAlert("info", "inputError", "info");
		return;
	}
	
	// 得到form Json数据
	var formData = getFormDeepJsonData(formId); // 这里这之前是getFormJsonData();
	if(typeof formDataFunc == 'function') {
		formData = formDataFunc(formData);
	}
	
	$.ajax({
	    url : url,
	    type : "POST",
	    contentType: "application/json; charset=utf-8",
	    datatype: "json",
	    data : JSON.stringify(formData),
	    success : function(data, stats) {
	    	_ajaxCallback(data, successFunc, failureFunc);
	    	if(hasProgress) {
	    		hideProgress();
	    	}
	    },
		error: function(ret) {
			_ajaxCallback(ret, successFunc, failureFunc);
			if(hasProgress) {
				hideProgress();
			}
		}
	});
}

/**
 * 与formSubmit()不同的是验证采用jQuery validation
 * 取的是json数据
 * 显示传递url
 * @param formId
 * @param formDataFunc
 * @param successFunc
 * @param failureFunc
 * @param hasProgress
 */
function formSubmit2(formId, url, formDataFunc, successFunc, failureFunc, hasProgress) {
	hasProgress = hasProgress == undefined ? true : hasProgress;
	if(hasProgress) {
		showProgress('info', 'processing');
	}
	var formId2 = $('#' + formId);
	// initValidator(formId);
	var isValid = valid(formId); // 验证数据, 用jQuery validation
	if(!isValid) {
		hideProgress();
		// msgAlert("info", "inputError", "info");
		return;
	}
	
	var formData = getFormJsonData(formId);
	if(typeof formDataFunc == 'function') {
		formData = formDataFunc(formData);
	}
	
	ajaxPost(url, formData, function(ret) {
		_ajaxCallback(ret, successFunc, failureFunc);
		hideProgress();
	});
}

// jquery validate
/* 支持:	
	required – Makes the element required.
	remote – Requests a resource to check the element for validity.
	minlength – Makes the element require a given minimum length.
	maxlength – Makes the element require a given maxmimum length.
	rangelength – Makes the element require a given value range.
	min – Makes the element require a given minimum.
	max – Makes the element require a given maximum.
	range – Makes the element require a given value range.
	email – Makes the element require a valid email
	url – Makes the element require a valid url
	date – Makes the element require a date.
	dateISO – Makes the element require a ISO date.
	number – Makes the element require a decimal number.
	digits – Makes the element require digits only.
	creditcard – Makes the element require a creditcard number.
	equalTo – Requires the element to be the same as another one
*/
if($.validator) {
	$.validator.setDefaults({
		submitHandler: function() { alert("submitted!"); },
		errorPlacement: function(error, element) {
			var parent = element.parent();
			if(parent.is(".each-chk")) { // dk:checkbox，checkbox与 label是一对
				parent = parent.parent();
			} else if(element.is(".combo-value")) {
				parent = parent.parent();
			}
			error.appendTo(parent);  
		},
		debug: true,
		ignore: "" // 隐藏表单也要验证
	});

	$.validator.addMethod("comboboxRequired", function(value, element) {
		var value = $(element).val();
		if(!value) {
			return false;
		}
		return true;
	}, "必选项");

	$.validator.addMethod("minSelect", function(value, element, param) {
		var length = $.isArray( value ) ? value.length : this.getLength($.trim(value), element);
		return this.optional(element) || length >= param;
	}, $.validator.format("至少选择 {0} 项"));
	
	
	//Double 数据库小数点后最多三位
	$.validator.addMethod("double", function(value, element, param) {
		
		if(value!=null && value!=""){
			var array=value.split(".");
			if(array.length>2){
				return false;
			}else if(array.length==2) {
				if(array[1].length>param){
					return false;
				}
			}
		}
		
		return true;
	}, $.validator.format("小数点后最多 {0} 位"));
}

/**
 * 初始化验证器
 * @param formId
 * @param isId 默认为true
 */
function initValidator(formId, isId) {
	if(isId == undefined) {
		isId = true;
	}
	if(isId) {
		$("#" + formId).validate();
	} else {
		$(formId).validate();
	}
}
// 验证
function valid(formId, isId) {
	if(isId == undefined) {
		isId = true;
	}
	if(isId) {
		return $("#" + formId).valid();
	} else {
		return $(formId).valid();
	}
}

// autocomplete

/**
 原生用法
 $("#auto").autocomplete(G_PATH + "/tips/searchOrg", {
	minChars: 2,
	width: 230,
	matchContains: false,
	autoFill: false,
    dataType: 'json', // 返回的数据类型为JSON类型 
	parse: function(data) { // 解释返回的数据，把其存在数组里 
		var rows = data.rows;
		var parsed = []; 
		for (var i = 0; i < rows.length; i++) { 
			var row = rows[i];
			parsed[i] = {
				data: row, // 数据行
				value: row.orgId, // 存储的结果
				result: row.orgName + "(" + row.account +")" // 显示到表单
			};
		}
		return parsed; 
	},
	// 显示到选项列表
    formatItem: function(row, i, max) {
		return row.orgName + "(" + row.account + ")";
	}
}).result(function(event, row, formatted) {
	// 把结果存在某个地方
    log(row.orgId);
});

// 示例
initAutocomplete("auto", G_PATH + "/tips/searchOrg", {
	valueKey: "orgId", // 要得到的值的key
	targetInputId: "orgIdHidden", // 存储选中值的inputId <input type="hidden" id="xx" /> 也可以不要，自动给你生成一个
	formatResult: function(row) { // 用户选中一条数据后，显示在表单里
		return row.orgName + "(" + row.account +")";
	},
	formatItem: function(row) { // 列表显示
		return row.orgName + "(" + row.account + ")";
	}
});

或：
initAutocomplete("auto", G_PATH + "/tips/searchOrg", {
	valueKey: "orgId",
	controller: 'desc/get',
	itemKeys: ["orgName", "account"],
	resultKeys: ["orgName", "account"]
});
*/
dk.autocomplete = {}; // 定义id与targerId映射
function initAutocomplete(id, url, userConfig) {
	// userConfig['formatResult'] = userConfig['formatResult'] || userConfig['formatItem'];
	// userConfig['formatItem'] = userConfig['formatItem'] || userConfig['formatResult'];
	var config = {
		id: id,
		minChars: 2,
		width: 230,
		matchContains: false,
		autoFill: false,
	    dataType: 'json', // 返回的数据类型为JSON类型
	    delimiter: ' ', // 分隔符
	    itemKeys: [], // 列表显示的字段
		resultKeys: [] // 选中一条数据后，在input里要显示的字段
	};
	
	$.extend(config, userConfig);
	
	// 是否已定义把值放在的input
	// 没有，则插入一个
	if(!config['targetInputId']) {
		var hiddenInputId = id + "_hidden_input";
		config['targetInputId'] = hiddenInputId;
		var name = $("#" + id).attr('name') || id;
		$("#" + id).attr("name", name + "_raw_input");
		$("#" + id).parent().append('<input type="hidden" name="' + name + '" id="' + hiddenInputId + '"/>');
	}

	
	// 是否需要通过listKeys, resultKeys来建立
	if(typeof config['formatResult'] != "function") {
		if(!isArray(config['resultKeys']) || config['resultKeys'].length == 0) {
			log("config['resultKeys'] 未定义，或formatResult方法未定义");
			return;
		}
		config['formatResult'] = function(row) {
			var str = "";
			for(var i = 0; i < config['resultKeys'].length; ++i) {
				var key = config['resultKeys'][i];
				if(!row[key]) {
					continue;
				}
				str += row[key];
				if(i != config['resultKeys'].length - 1) {
					str += config['delimiter'];
				}
			}
			return str;
		};
	}
	if(typeof config['formatItem'] != "function") {
		if(!isArray(config['itemKeys']) || config['itemKeys'].length == 0) {
			log("config['itemKeys'] 未定义，或formatItem方法未定义");
			return;
		}
		config['formatItem'] = function(row) {
			var str = "";
			for(var i = 0; i < config['itemKeys'].length; ++i) {
				var key = config['itemKeys'][i];
				str += row[key];
				if(i != config['itemKeys'].length - 1) {
					str += config['delimiter'];
				}
			}
			return str;
		};
	}
	
	// parse
	var valueKey = config['valueKey'];
	if(!config['parse']) {
		if(!valueKey) {
			return;
		}
		config.parse = function(rows) { // 解释返回的数据，把其存在数组里 
			if(!rows || !isArray(rows)) {
				return null;
			}
			var parsed = []; 
			for (var i = 0; i < rows.length; i++) { 
				var row = rows[i];
				parsed[i] = {
					data: row, // 数据行
					value: row[valueKey], // 存储的结果
					result:  config['formatResult'](row)// 显示到表单
				};
			}
			return parsed; 
		};
	}
	
	// 建立映射，可以得到值
	dk.autocomplete[id] = config;
	
	$("#" + id).autocomplete(url, config).result(function(event, row, formatted) {
		$("#" + config.targetInputId).val(row[valueKey]);
	});
}

/**
 * 得到autocomplete值
 * @param id autocomplete ID，不是targerInputId
 * @returns
 */
function getAutocompleteValue(id) {
	if(!dk.autocomplete[id] || !dk.autocomplete[id]['targetInputId']) {
		return null;
	}
	return $("#" + dk.autocomplete[id]['targetInputId']).val();
}

/**
 * 设置值
 * @param row
 * @returns
 */
function setAutocompleteValue(id, row) {
	if(!dk.autocomplete[id] || !dk.autocomplete[id]['targetInputId']) {
		return null;
	}
	
	var config = dk.autocomplete[id];
	
	if(!config) {
		return;
	}
	
	var result = config['formatResult'](row);
	var value = row[config['valueKey']];

	log(result==undefined);
	if(!result||result==undefined) {
		//log(result);
		// 需要通过controller来取?
		if(config.controller) {
			ajaxGet(G_PATH + "/" + config.controller, {ids: [value]}, function(ret) {
				$("#" + id).val(ret[value]);
			}, "", false);
		}
	}
	
	$("#" + id).val(result);
	$("#" + config.targetInputId).val(value);
	
	return true;
}

// log 日志
function log(o) {
	if(window.console) {
		console.log(o);
	}
}

// 刷新单元格
// rowNum从0开始
function refreshCell(tableId, field, rowNum, val) {
	var target = $('#' + tableId).parent().children().find('tr[datagrid-row-index="' + rowNum + '"] td[field="' + field + '"] div');
	if(target) {
		target.html(val);
	}
}

//searchFlag 为指定的 结算类型查询controller 方法的 的url
var commonPaymentItemTypeName;
var commonPaymentItemTypeId;
var commonPaymentItemTypeFlag;
var searchPaymentItemFunc;
function searchPaymentItemType(paymentItemTypeName,paymentItemTypeId,searchFlag,func){
	commonPaymentItemTypeName=paymentItemTypeName;
	commonPaymentItemTypeId=paymentItemTypeId;
	commonPaymentItemTypeFlag=searchFlag;
	searchPaymentItemFunc=func;
	var config = {
			title: 'searchPaymentItemType',
		    width: 900,
		    height: 487,
		    href:getCommonPath('paymentItemType/common_paymentItemType_search')
		};
		openWin(config, 'searchPaymentItemTypeWin');
}

// theme
function getThemeTypeId(themeName) {
	// themeName是css文件名, typeId是存数据库中的
	for(var themeTypeId in G_USER['themeTypeId2ThemeName']) {
		if(themeName2ThemeTypeId[themeTypeId] == themeName) {
			return themeTypeId;
		}
	}
	return "";
}
// 根据themeTypeId得到themeName
function getThemeName(themeTypeId) {
	return G_USER['themeTypeId2ThemeName'][themeTypeId];
}