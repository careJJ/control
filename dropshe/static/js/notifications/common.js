(function($) {
    //页面输出错误信息
    window.errors = function (msg,name){
        var data = sessionStorage.getItem('login');
        var inputs = $("input[name= '"+name+"']");
        var input_plus = $("input[name= '"+data+"']");
        var str = "<p id='"+name+"plus' style='color: #F08118;'>"+msg+"</p>"+
            "<input type='hidden' id='"+name+"' value='"+msg+"'>";

        if($("#"+name+"").val() == undefined){
            inputs.css("border","1px solid #F08118");
            inputs.after(str);
            if(data != name){
                $('#'+data+'').remove();
                $('#'+data+'plus').remove();
                input_plus.css("border","1px solid #EEF1F8");
                sessionStorage.setItem('login', ''+name+'');
            }
        }else{
            $('#'+data+'').remove();
            $('#'+data+'plus').remove();
            inputs.css("border","1px solid #F08118");
            inputs.after(str);
        }
    };
    //客户收集页面单独处理错误信息
    window.error_plus = function (msg,name){
        if(name === 'team_name'){
            var inputs = $("input[name= 'team_name']");
            var str = "<p id='team_name_plus' style='color: #F08118;'>Team name is required</p>";
            if($('#team_name_plus').val() == undefined){
                inputs.css("border","1px solid #F08118");
                inputs.after(str);
            }
            setTimeout(function () {
                $('#team_name_plus').remove();
                inputs.css("border","1px solid #EEF1F8");
            },5000);
        }else {
            var data = sessionStorage.getItem('test');
            console.log(data);
            var inputs = $("select[name= '"+name+"']");
            var input_plus = $("select[name= '"+data+"']");
            var str = "<p id='"+name+"plus' style='color: #F08118;'>"+msg+"</p>"+
                "<input type='hidden' id='"+name+"s' value='"+msg+"'>";

            if($("#"+name+"s").val() == undefined){
                inputs.css("border","1px solid #F08118");
                inputs.after(str);
                if(data != name){
                    $('#'+data+'s').remove();
                    $('#'+data+'plus').remove();
                    input_plus.css("border","1px solid #EEF1F8");
                    sessionStorage.setItem('test', ''+name+'');
                }
            }
        }
    };

})(jQuery);