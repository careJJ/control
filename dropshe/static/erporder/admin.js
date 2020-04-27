jQuery(document).ready(function($) {

    window.toast = function(msg,time){

        if(time == undefined){
            time = 2000;
        }

        new $.zui.Messager(msg, {
            icon: 'bell', // 定义消息图标
            time: time,
        }).show();
    }

    // 更新网页地址
    window.update_url = function(url, key) {
        var key = (key || 't') + '='; //默认是"t"
        var reg = new RegExp(key + '\\d+'); //正则：t=1472286066028
        var timestamp = +new Date();
        if (url.indexOf(key) > -1) { //有时间戳，直接更新
            return url.replace(reg, key + timestamp);
        } else { //没有时间戳，加上时间戳
            if (url.indexOf('\?') > -1) {
                var urlArr = url.split('\?');
                if (urlArr[1]) {
                    return urlArr[0] + '?' + key + timestamp + '&' + urlArr[1];
                } else {
                    return urlArr[0] + '?' + key + timestamp;
                }
            } else {
                if (url.indexOf('#') > -1) {
                    return url.split('#')[0] + '?' + key + timestamp + location.hash;
                } else {
                    return url + '?' + key + timestamp;
                }
            }
        }
    }

    // 检测字符串是否符合json数据格式
    window.is_json = function(json){

        if (json.match("^\{(.+:.+,*){1,}\}$")){
            // 符合
            return true;
        }

        return false;
    }

    function check_img_url_ext(img_url){
        // 获取最后一个.的位置
        var index = img_url.lastIndexOf('.');
        // 获取后缀名
        var ext = img_url.substr(index + 1);
        // 判断后缀名是否在指定范围内
        var exts = ['png', 'jpg', 'jpeg', 'gif'];
        return exts.indexOf(ext.toLowerCase()) !== -1;
    }

    // 表格页开关
    $(document).on('change', '.adm-list-switch', function(event) {
        event.preventDefault();
        /* Act on the event */
        var $this = $(this);

        // 获取最新值
        var value = $this.prop('checked') ? 1 : 0 ;

        // 构建提交数据
        var data = {};
        data['token'] = token;
        data['id'] = $this.data('id');
        data['field'] = $this.data('name');
        data['value'] = value;

        $.ajax({
            url: $this.data('url'),
            type: 'POST',
            cache: false,
            data: data,
        }).done(function(res) {

            //判断返回值不是 json 格式
            if(is_json(res) == false){
                toast('无法连接服务器，请刷新重试');
                return false;
            }

            // 将返回的图片地址添加到更新到页面上。
            var data = jQuery.parseJSON(res);

            toast(data['msg']);

            if(data['status'] != 1){
                // 修改失败，将值恢复到原来的样子
                setTimeout(function(){
                    $this.prop('checked',!value);
                },800);
            }

        }).fail(function(res) {
            toast("网络出错");
        });

    });

    // 表格页开关
    $(document).on('change', '.adm-form-switch', function(event) {
        event.preventDefault();
        /* Act on the event */
        var $this = $(this);

        // 获取最新值
        var value = $this.prop('checked') ? 1 : 0 ;

        $this.val(value);

    });
    // 颜色选择器
    $(".adm-colorpicker .control-input").each(function(index, el) {
        var $el = $(el);
        $el.colorpicker({
            "color": ($el.data('color') != undefined ? $el.data('color') : '#e5e5e5')
        });
    });

    // 日期选择器
    $(".adm-datetimepicker .control-input").each(function(index, el) {

        var $el = $(el);

        $el.datetimepicker({
            language:  "zh-CN",
            weekStart: ($el.data('weekstart') != undefined ? $el.data('weekstart') : true),
            todayBtn:  ($el.data('todaybtn') != undefined ? $el.data('todaybtn') : true),
            autoclose: ($el.data('autoclose') != undefined ? $el.data('autoclose') : true),
            todayHighlight: ($el.data('todayhighlight') != undefined ? $el.data('todayhighlight') : true),
            startView: ($el.data('startview') != undefined ? $el.data('startview') : 2),
            maxView: ($el.data('maxview') != undefined ? $el.data('maxview') : 4),
            minView: ($el.data('minview') != undefined ? $el.data('minview') : 0),
            forceParse: ($el.data('forceparse') != undefined ? $el.data('forceparse') : 0),
            format: ($el.data('format') != undefined ? $el.data('format') : 'yyyy-mm-dd hh:ii:ss'),
            startDate: ($el.data('startdate') != undefined ? $el.data('startdate') : null),
            endDate: ($el.data('enddate') != undefined ? $el.data('enddate') : null),
            minuteStep: ($el.data('minutestep') != undefined ? $el.data('minutestep') : null),
            initialDate: ($el.data('initialdate') != undefined ? $el.data('initialdate') : null),
        });
    });

    // 富文本编辑器
    $(".wysiwyg-editor textarea").each(function(index, el) {
        KindEditor.create(this,{
            uploadJson : base_url + '/index.php/'+backend_module+'/attachment/ke_img_upload/',
            fileManagerJson : base_url + '/index.php/'+backend_module+'/attachment/ke_file_manager/',
            allowFileManager : true,
            minWidth:320,
            minHeight:240,
        });
    });

    // 右侧主要内容高度控制
    var right_con_height = 0;
    function window_resize(){
        sider_con_height = $(window).height() - 50;
        right_con_height = $(window).height() - 96;
        $(".adm-sider").css('height', sider_con_height + 'px');
        $(".adm-right").css('min-height', right_con_height + 'px');
        $(".adm-frame").height(right_con_height);
    }

    // 浏览器窗口大小变化时更新高度
    $(window).on('resize',function () {
        window_resize();
    });

    // 页面第一次打开时确定高度
    window_resize();

    // 左侧菜单隐藏显示
    $(document).on('click', '.item-header', function(event) {
        event.preventDefault();
        var $this = $(this);
        // 隐藏所有菜单
        $this.parents('.sider-menu').find('.item-list').hide();
        // 显示当前菜单
        $this.parent().find('.item-list').toggle();
    });

    // 左侧菜单点击
    $(document).on('click', '.item-list a', function(event) {
        // event.preventDefault();
        var $this = $(this);

        $this.parents(".sider-menu").find('.item-list a').removeClass('active');

        $this.addClass('active');

        $(".breadcrumb_pos_m").text($this.parents(".menu-item").find('.item-header a').text());

        $(".breadcrumb_pos_c").text($this.text());

    });

    // 右上角卸载功能
    $(document).on('click', '#uninstall', function(event) {
        if (confirm("确定卸载重装系统？")) {
            return true;
        } else {
            return false;
        }
    });

    // 应用安装页卸载应用提示
    $(document).on('click', '.app-uninstall', function(event) {
        if (confirm("确定卸载该应用？")) {
            return true;
        } else {
            return false;
        }
    });

    // 列表页删除单项
    $(document).on('click', '.del-item', function(event) {

        var $this = $(this);

        // 防止重复提交
        if($this.data('clk') == 'clicked'){
            toast('数据已提交，请耐心等待');
            return false;
        }

        // 询问是否要删除数据
        if(!confirm("确定删除所选数据？")){
            return false;
        }

        $this.data('clk','clicked');

        // 通过ajax提交删除请求
        var url = $this.attr('href');
        // 提交数据
        $.ajax({
            type: 'POST',
            url: url,
            data: {
                token:token
            },
            cache: false,
            success: function(json, textStatus, jqXHR){
                // 允许重新点击按钮
                $this.data('clk','reclick');
                //判断返回值不是 json 格式
                if(is_json(json) == false){
                    toast('无法连接服务器，请刷新重试');
                    return false;
                }
                //将字符串转换为对象
                var remsg = jQuery.parseJSON(json);
                if(remsg.status == 1){
                    // 成功，执行代码
                    toast(remsg.msg);
                    // 返回上一页
                    setTimeout(function(){
                        window.location.reload();
                    },1000);
                    return true;
                }
                // 失败，弹出错误信息
                toast(remsg.msg);
            },
            error: function(jqXHR, textStatus, errorThrown){
                // 失败，弹出错误信息
                toast(remsg.msg);
            }
        });

        return false;
    });

    // 列表页全选
    $(document).on('click', '#idscb', function(event) {
        var $this = $(this);
        if ($this.prop('checked') == true) {
            $(".idscb").prop("checked", true);
        } else {
            $(".idscb").prop("checked", false);
        }
    });

    // 列表页预览图片
    $(document).on('click', '.adm-list-image', function(event) {
        
        var $this = $(this);

        $('.adm-modal-img').attr('src', $this.attr('src'));

        $('.adm-modal').modal('show');

    });

    // 多列文本
    $(document).on('click', '.multiple-text-remove', function(event) {
        event.preventDefault();
        $(this).parents('.control-input').remove();
    });

    function multiple_text_add($obj){

        var $add_button = $obj.find('.multiple-text-add');

        var name = $add_button.data('name');

        var options = $add_button.data('options');


        var style = $add_button.data('style');

        if(!style){
            style = '';
        }

        var add_html = '';

        var timeid = Date.now();

        if(options){

            // 先遍历一次获取对象的个数
            var count=3;
            for(var i in options){ count++; }
            
            add_html = add_html + '<div class="control-input input-group mb15">';

            for(var i in options){
                add_html = add_html + '<input type="text" name="data['+ name +']['+timeid+']['+i+']" class="form-control" value="" placeholder="'+options[i]+'" style="display: inline-block;width:'+(100/count)+'%;'+style+'">';
            }

            add_html = add_html + '<button type="button" class="btn btn-action multiple-text-remove" >移除</button></div>';
        }

        $(add_html).insertBefore($add_button.parents('.control-input'));
    }

    // 添加多列文本
    $(document).on('click', '.multiple-text-add', function(event) {
        event.preventDefault();
        var $this = $(this);
        var $parent = $this.parents('.multiple-text');
        multiple_text_add($parent);
    });

    // 单图上传文本框
    $(document).on('change', '.single-img-upload .single-img-upload-input', function(event) {
        event.preventDefault();
        /* Act on the event */

        var $this = $(this);

        var img_url = $this.val().toLowerCase();

        // 判断url是否包含http或https
        if(img_url.indexOf('http') === 0){
            $this.parents('.control-input').find('.preview-img img').attr('src', $this.val());
        }

    });

    // 单图上传按钮
    $(document).on('change', '.single-img-upload .upload-input', function(event) {
        event.preventDefault();
        /* Act on the event */

        var $this = $(this);

        var $parent = $this.parents('.control-input');

        var formData = new FormData();

        formData.append('img', $this[0].files[0]);

        // 单图上传地址
        var url = base_url + '/index.php/'+backend_module+'/attachment/img_upload/';

        // 通过ajax使用post方式上传
        $.ajax({
            url: url,
            type: 'POST',
            cache: false,
            data: formData,
            processData: false,
            contentType: false
        }).done(function(res) {
            // 将返回的图片地址添加到更新到页面上。
            var data = jQuery.parseJSON(res);
            if(data['status'] == 1){

                $parent.find('.preview-img img').attr('src', data['image']);
                $parent.find('.preview-img .single-img-upload-input').val(data['image']);
                $parent.find('.preview-img,.preview-img .remove-btn').show();
                $parent.find('.img-upload-box').hide();

            }else{
                toast(data['msg']);
            }
        }).fail(function(res) {
            toast("上传失败，请重新上传文件！");
        });
    });

    // 单图上传-输入网址
    $(document).on('click', '.single-img-upload .add-text-btn', function(event) {
        event.preventDefault();
        var $this = $(this);
        var $parent = $this.parents('.control-input');

        bootbox.prompt({ 
            size: "large",
            title: "请输入图片地址：", 
            callback: function(result){
                if(result === null || result === ''){
                    return true;
                }
                // 分析后缀名
                if(check_img_url_ext(result) == false){
                    return true;
                }
                $parent.find('.preview-img img').attr('src', result);
                $parent.find('.preview-img .single-img-upload-input').val(result);
                $parent.find('.preview-img,.preview-img .remove-btn').show();
                $parent.find('.img-upload-box').hide();
            }
        });
    });

    // 单图上传-删除图片
    $(document).on('click', '.single-img-upload .remove-btn', function(event) {
        event.preventDefault();
        var $this = $(this);
        var $parent = $this.parents('.control-input');

        $parent.find('.preview-img .single-img-upload-input').val('');

        $parent.find('.preview-img,.preview-img .remove-btn').hide();
        $parent.find('.img-upload-box').show();
    });

    // 多图上传
    $(document).on('change', '.multiple-img-upload .multiple-img-upload-input', function(event) {
        event.preventDefault();
        /* Act on the event */

        var $this = $(this);

        var img_url = $this.val().toLowerCase();

        // 判断url是否包含http或https
        if(img_url.indexOf('http') === 0){
            $this.parents('.control-input').find('.preview-img img').attr('src', $this.val());
        }

    });

    // 预览图片
    $(document).on('click', '.preview-img img', function(event) {
        event.preventDefault();

        var $this = $(this);

        console.log($this.attr('src'));
        
        $('.adm-modal-img').attr('src', $this.attr('src'));

        $('.adm-modal').modal('show');
    });

    // 多图片上传，增加图片
    function multiple_img_upload_add($obj,imgs){

        var $add_box = $obj.find('.img-upload-box');

        var name = $add_box.data('name');

        var style = $add_box.data('style');

        if(!style){
            style = '';
        }
        var add_html = '';

        if(imgs){

            var imgs_len = imgs.length;

            for (var i = 0; i < imgs_len; i++) {
                add_html = add_html + '<div class="preview-img">'+
                        '<input type="text" hidden name="data['+ name +'][]" class="form-control multiple-img-upload-input adm-hide" value="'+imgs[i]+'" style="'+style+'">' +
                        '<img src="'+imgs[i]+'">'+
                        '<i class="icon icon-remove-circle remove-btn"></i>'+
                    '</div>';
            }
        }else{
            add_html = '<div class="preview-img">'+
                        '<input type="text" hidden name="data['+ name +'][]" class="form-control multiple-img-upload-input adm-hide" value="" style="'+style+'">' +
                        '<img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAALEgAACxIB0t1+/AAAABx0RVh0U29mdHdhcmUAQWRvYmUgRmlyZXdvcmtzIENTNui8sowAAAAWdEVYdENyZWF0aW9uIFRpbWUAMDcvMDUvMTj7ZJCvAAAADUlEQVQImWP4////fwAJ+wP9CNHoHgAAAABJRU5ErkJggg==">'+
                    '</div>';
        }

        $(add_html).insertBefore($add_box);
    }

    $(document).on('change', '.multiple-img-upload .upload-input', function(event) {
        event.preventDefault();
        /* Act on the event */
        var $this = $(this);

        var $parent = $this.parents('.multiple-img-upload');

        var formData = new FormData();

        var files_len = $this[0].files.length;

        for (var i = 0; i < files_len; i++) {
            formData.append('imgs'+i,$this[0].files[i]);
        }

        // 多图上传地址
        var url = base_url + '/index.php/'+backend_module+'/attachment/imgs_upload/';
        // 通过ajax使用post方式上传
        $.ajax({
            url: url,
            type: 'POST',
            cache: false,
            data: formData,
            processData: false,
            contentType: false
        }).done(function(res) {
            // 将返回的图片地址添加到更新到页面上。
            var data = jQuery.parseJSON(res);
            if(data['status'] == 1){
                // 上传成功
                multiple_img_upload_add($parent,data['images']);
            }else{
                toast(data['msg']);
            }
        }).fail(function(res) {
            toast("上传失败，请重新上传文件！");
        });
    });

    // 多图上传-删除单个图片
    $(document).on('click', '.multiple-img-upload .remove-btn', function(event) {
        event.preventDefault();
        var $this = $(this);
        $this.parents('.preview-img').remove();
    });

    // 多图上传-输入网址
    $(document).on('click', '.multiple-img-upload .add-text-btn', function(event) {
        event.preventDefault();
        var $this = $(this);
        var $parent = $this.parents('.multiple-img-upload');

        bootbox.prompt({
            size: 'large',
            title: "请输入图片地址，每行一个：",
            inputType: 'textarea',
            callback: function (result) {

                if(result === null || result === ''){
                    return true;
                }

                var all_imgs = result.split(/[\s\n]/);
                var final_imgs = new Array();

                var len = all_imgs.length;

                if(len <= 0){
                    return true;
                }

                for (var i = 0; i < len; i++) {

                    var temp_img = all_imgs[i];
                    // 分析后缀名
                    if(check_img_url_ext(temp_img) == true){
                        final_imgs.push(temp_img);
                    }
                }

                if(final_imgs.length > 0){
                    multiple_img_upload_add($parent,final_imgs);
                }
                
            }
        });
    });

    // 单文件上传
    $(document).on('change', '.single-file-upload .upload-input', function(event) {
        event.preventDefault();
        /* Act on the event */
        var $this = $(this);

        var $parent = $this.parents('.control-input');

        var formData = new FormData();

        formData.append('file', $this[0].files[0]);
        // 单文件上传地址
        var url = base_url + '/index.php/'+backend_module+'/attachment/file_upload/';
        // 通过ajax使用post方式上传
        $.ajax({
            url: url,
            type: 'POST',
            cache: false,
            data: formData,
            processData: false,
            contentType: false
        }).done(function(res) {
            // 将返回的图片地址添加到更新到页面上。
            var data = jQuery.parseJSON(res);
            if(data['status'] == 1){

                // 上传成功
                $parent.find('.single-file-upload-input').val(data['file']);

            }else{
                toast(data['msg']);
            }
        }).fail(function(res) {
            toast("上传失败，请重新上传文件！");
        });  
    });

    // 多文件上传
    $(document).on('click', '.multiple-file-upload-remove', function(event) {
        event.preventDefault();
        $(this).parents('.control-input').remove();
    });

    function multiple_file_upload_add($obj,files){

        var $add_button = $obj.find('.multiple-file-upload-add');

        var name = $add_button.data('name');

        var style = $add_button.data('style');

        if(!style){
            style = '';
        }

        var add_html = '';

        if(files){
            console.log(files);
            var files_len = files.length;

            for (var i = 0; i < files_len; i++) {
                add_html = add_html + '<div class="control-input mb15">' +
                    '<input type="text" name="data['+ name +'][]" class="form-control multiple-file-upload-input" value="'+files[i]+'" style="'+style+'">' +
                    '<button type="button" class="btn btn-action multiple-file-upload-remove" >移除</button>'+
                '</div>';
            }
        }else{
            add_html = '<div class="control-input mb15">' +
                    '<input type="text" name="data['+ name +'][]" class="form-control multiple-file-upload-input" value="" style="'+style+'">' +
                    '<button type="button" class="btn btn-action multiple-file-upload-remove" >移除</button>'+
                '</div>';
        }


        $(add_html).insertBefore($add_button.parents('.control-input'));
    }

    $(document).on('click', '.multiple-file-upload-add', function(event) {
        event.preventDefault();
        /* Act on the event */

        var $this = $(this);
        
        var $parent = $this.parents('.multiple-file-upload');

        multiple_file_upload_add($parent);
    });

    $(document).on('change', '.multiple-file-upload .upload-input', function(event) {
        event.preventDefault();
        /* Act on the event */
        var $this = $(this);

        var $parent = $this.parents('.multiple-file-upload');

        var formData = new FormData();

        var files_len = $this[0].files.length;
        
        for (var i = 0; i < files_len; i++) {
            formData.append('files'+i,$this[0].files[i]);
        }

        // 多文件上传地址
        var url = base_url + '/index.php/'+backend_module+'/attachment/files_upload/';
        // 通过ajax使用post方式上传
        $.ajax({
            url: url,
            type: 'POST',
            cache: false,
            data: formData,
            processData: false,
            contentType: false
        }).done(function(res) {
            // 将返回的图片地址添加到更新到页面上。
            var data = jQuery.parseJSON(res);
            if(data['status'] == 1){
                // 上传成功
                multiple_file_upload_add($parent,data['files']);
            }else{
                toast(data['msg']);
            }
        }).fail(function(res) {
            toast("上传失败，请重新上传文件！");
        });
    });

    function isInteger(num) {
        return typeof num === 'number' && num%1 === 0;
    }

    // 输入页码直接跳转
    $(document).on('click', '.btn-jump', function(event) {
        event.preventDefault();
        /* Act on the event */

        var $this = $(this);
        var page = parseInt($this.parents('.adm-list-jump').find('#adm-jump-page').val());

        if(isInteger(page)){
            // 跳转到新页面
            // 判断页码是否大于最大页码
            var last_page = $this.data('last');

            if(page >last_page){
                page = last_page;
            }

            if(page < 1){
                page = 1;
            }
            // 开始跳转
            var url = $this.data('url');
            if(url.charAt(url.length-1) != '/'){
                url = url + '/';
            }
            window.location.href = url + 'page/'+page;
        }else{
            toast("请输入正确的页码！");
        }
    });

    function role_checkbox(pid,status){

        $(".adm-role-pid-" + pid).each(function(index, el) {
            var $el = $(el);

            var id = $el.data('id');

            var parent = $el.data('parent');

            $el.prop('checked', status);

            if(parent == 1){
                // 还有下一级菜单
                role_checkbox(id,status);
            }

        });
    }


    // 角色管理——权限
    $(document).on('click', '.adm-role-inp', function(event) {
        /* Act on the event */
        var $this = $(this);

        var id = $this.data('id');

        var parent = $this.data('parent');

        var status = false;

        if ($this.prop('checked') == true) {
            status = true;
        } else {
            status = false;
        }

        if(parent == 1){
            role_checkbox(id,status);
        }else{
            
            $this.closest('tr').find('.checkbox-inline input').prop('checked', status);
        }

    });

    // 通用表单，提交数据按钮btn-submit
    $(document).on('submit', '.adm-form form,.adm-list form', function(event) {
        event.preventDefault();
        /* Act on the event */

        var $this = $(this);

        // 防止重复提交
        if($this.data('clk') == 'clicked'){
            toast('数据已提交，请耐心等待');
            return false;
        }

        $this.data('clk','clicked');

        var url = window.location.href;

        // 开始提交数据
        $.ajax({
            type: 'POST',
            url: url,
            data: $this.serialize(),
            cache: false,
            success: function(json, textStatus, jqXHR){

                // 允许重新点击按钮
                $this.data('clk','reclick');

                //判断返回值不是 json 格式
                if(is_json(json) == false){
                    toast('无法连接服务器，请刷新重试');
                    return false;
                }

                //将字符串转换为对象
                var remsg = jQuery.parseJSON(json);

                // 弹出提示信息
                toast(remsg.msg);

                if(remsg.status == 1){
                    // 成功，执行代码
                    // 返回上一页
                    setTimeout(function(){
                        window.location.href = update_url(document.referrer);
                    },350);

                    return true;
                }

                // 失败，不处理
            },
            error: function(jqXHR, textStatus, errorThrown){
                // 失败，弹出错误信息
                toast(remsg.msg);
            }
        });

        return false;
    });

    // 批量提交数据
    $(document).on('click', '.btns-bot .btn', function(event) {
        event.preventDefault();
        /* Act on the event */

        var $this = $(this);

        var action = $this.data('action');

        var $form = $('.adm-list form');

        // 防止重复提交
        if($form.data('clk') == 'clicked'){
            toast('数据已提交，请耐心等待');
            return false;
        }

        // 询问是否要删除数据
        if(action == 'delete' && !confirm("确定删除所选数据？")){
            return false;
        }

        // 询问是否要删除全部数据
        if(action == 'delall' && !confirm("确定删除表单全部数据？")){
            return false;
        }

        // 避免重复提交
        $form.data('clk','clicked');

        // 设置批量操作类型
        $form.find('.action').val(action);

        // 获取提交地址
        var url = window.location.href;

        // 开始提交数据
        $.ajax({
            type: 'POST',
            url: url,
            data: $form.serialize(),
            cache: false,
            success: function(json, textStatus, jqXHR){

                // 允许重新点击按钮
                $form.data('clk','reclick');

                //判断返回值不是 json 格式
                if(is_json(json) == false){
                    toast('无法连接服务器，请刷新重试');
                    return false;
                }

                //将字符串转换为对象
                var remsg = jQuery.parseJSON(json);

                // 弹出提示信息
                toast(remsg.msg);

                if(remsg.status == 1){
                    // 成功，执行代码
                    // 返回上一页
                    setTimeout(function(){
                        window.location.reload();
                    },350);

                    return true;
                }

                // 失败，不处理
                
            },
            error: function(jqXHR, textStatus, errorThrown){
                // 失败，弹出错误信息
                toast(remsg.msg);
            }
        });

        return false;
    });

    // 通过ajax提交数据
    $(document).on('click', 'a.adm-ajax, button.adm-ajax', function(event) {
        event.preventDefault();
        /* Act on the event */

        var $this = $(this);

        var url = $this.attr('href');

        // 防止重复提交
        if($this.data('clk') == 'clicked'){
            toast('数据已提交，请耐心等待');
            return false;
        }

        $this.data('clk','clicked');

        // 开始提交数据
        $.ajax({
            type: 'POST',
            url: url,
            data: {
                token : token
            },
            cache: false,
            success: function(json, textStatus, jqXHR){

                // 允许重新点击按钮
                $this.data('clk','reclick');

                //判断返回值不是 json 格式
                if(is_json(json) == false){
                    toast('无法连接服务器，请刷新重试');
                    return false;
                }

                //将字符串转换为对象
                var remsg = jQuery.parseJSON(json);

                // 弹出提示信息
                toast(remsg.msg);

                if(remsg.status == 1){
                    // 成功，执行代码
                    // 返回上一页
                    setTimeout(function(){
                        window.location.href = update_url(window.location.href);
                    },350);

                    return true;
                }

                // 失败，不处理
            },
            error: function(jqXHR, textStatus, errorThrown){
                // 失败，弹出错误信息
                toast(remsg.msg);
            }
        });

        return false;
    });
});