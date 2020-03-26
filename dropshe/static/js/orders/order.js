let order_data = [];
$(document).keyup(function(e){

    let curKey = e.which;

    if(curKey==13){
        let searchContent = $(".serach-input").val();

        if (searchContent == '') {
            return false;
        }

        $(".serach-btn").trigger('click');
        //调用查询方法
    }
});

$(document).ready(function () {
    $("#paymentMethodEditModal .close").on("click", function () {
        $("#paymentMethodEditModal").hide();
    });
    $(".express-select").select2();

    $(".express-select").on("change", function () {
        let express_id = $(this).closest('tr').find('.express-select option:selected').attr('data-express-id');
        console.log(express_id);
        let from_id = $(this).attr('data-from-country-id');
        let to_id = $(this).attr('data-to-country-id');
        let price = $(this).attr('data-product-cost-price');
        let num = $(this).attr('data-product-num');
        let express_price = $(this).val();
        let th = $(this);

        let pay_fee = price * num + parseFloat(express_price);
        th.closest('tr').find('.express-price').text(express_price);
        th.closest('tr').find('.substituting-fees').text(pay_fee);
        th.closest('tr').find('.pay-btn').attr('data-price', pay_fee);
        let pay_hidden_input = th.closest('tr').find('.pay-hidden-input');
        pay_hidden_input.attr('data-price', pay_fee);
        pay_hidden_input.attr('data-express-id', express_id);

        let total_price = 0;
        th.closest('table').find('.pay-hidden-input').each(function () {
            total_price += parseFloat($(this).attr('data-price'));
        });
        total_price = total_price.toFixed(2);
        th.closest('table').find('.order-all-fee').text(total_price);
    });

    $(".btn-pay-now").on('click', function() {
        let pay_type = $("input[name='paymentMethod']:checked").val();
        if (pay_type != 1 && pay_type != 2) {
            alert('Please choose right pay type');
            return false;
        }

        let th = $(this);
        $(this).addClass('btn-loading');

        let product_id = $("#pay-now-product-id").val();
        let order_id = $("#pay-now-order-id").val();
        let express_id = $("#pay-now-express-id").val();

        console.log(order_data);

        $.ajax({
            url : '/order/drop/pay',
            type: 'get',
            data: {
                'pay_type' : pay_type,
                'order_data': order_data,
            },
            success: function(data) {
                if (data.code != 0) {
                    th.removeClass('btn-loading');
                    alert(data.msg);
                    return false;
                }
                alert('Pay success');
                location.reload();
            }
        });
    });

    $(".pay-all-btn").on("click", function () {
        $(this).closest('table');

        let total_price = 0;

        let index = 0;
        $(this).closest('table').find('.pay-hidden-input').each(function () {
            let product_price = $(this).attr('data-price');
            product_price = parseFloat(product_price);
            total_price += product_price;
            order_data[index] = {'product_id' : $(this).attr('data-product-id'), 'order_id' : $(this).attr('data-order-id'), 'express_id' : $(this).attr('data-express-id')};
            console.log(order_data[index]);
            index++;
        });
        console.log(order_data);
        updatePayData(total_price);
    });

    $(".pay-btn").on("click", function() {
        let order_id = $(this).attr('data-order-id');
        let product_id = $(this).attr('data-product-id');
        let pay_fee = $(this).attr('data-price');
        let pay_express_id = $(this).attr('data-express-id');
        $("#pay-now-order-id").val(order_id);
        $("#pay-now-product-id").val(product_id);
        $("#pay-now-express-id").val(pay_express_id);

        order_data = [{'order_id' : order_id, 'product_id' : product_id, 'express_id' : pay_express_id}];

        updatePayData(pay_fee);
    });

    $(".btn-cancel").on("click", function () {
        $("#paymentMethodEditModal").hide();
    });

    $(".sync-btn").on("click", function (event) {
        sync_order();
    });

    $(".serach-btn").on("click", function () {
        let search_input = $(".serach-input").val();
        let start_time = $("#start_time").val();
        let end_time = $("#end_time").val();

        let result = check();
        if (!result) {
            return false;
        }
        if (search_input == '') {
            return false;
        }

        location.href = '/order/drop?serach-input=' + search_input + '&start_time=' + start_time + '&end_time=' + end_time;
    });

    $('.js-example-basic-single').select2({
        minimumResultsForSearch: -1
    });

    $('.datepicker-input').fdatepicker({
        format: 'yyyy-mm-dd hh:ii',
        pickTime: true,
        disableDblClickSelection: true
    });

    //替换时间格式
    for (var el of $('.tolocaltime')) {
        var e = $(el);
        var text = e.attr('data');
        if (text == '-' || text == '') {
            continue;
        }
        try {
            var time = parseInt(text) * 1000;
        } catch (error) {
            //do nothing
            return;
        }
        var d = moment(time).format('YYYY-MM-DD hh:mm');
        if (e.is('input')) {
            e.val(d);
        } else {
            e.text(d);
        }

    }
});

/**
 * 更新数据
 *
 * @param total_price
 * @returns {boolean}
 */
function updatePayData(total_price)
{
    if (total_price <= 0) {
        return false;
    }
    total_price = parseFloat(total_price);
    total_price = total_price.toFixed(2);

    $("#total-pay-fee").text(total_price);

    $.ajax({
        url : '/order/drop/getPayAccount',
        type: 'get',
        data: {
        },
        success: function(data) {
            let pay_account = data.data['pay_account'];
            let balance = data.data['balance'];
            if (pay_account == '') {
                let html = '<a href="/accout/billing" target="_blank"><button type="button" class="btn btn-bg-main btn-big" style="width: 200px;margin-left:200px;background-color: ">Add Credit Card</button></a>';
                $(".payment-method-item:first").empty().html(html);
            } else {
                let html = '<div class="radio1"><input type="radio" name="paymentMethod" value="1"/><span class="icon"></span></div>\n' +
                    '                            <img class="payment-icon" src="/static/home/img/icon-pme1.png" />\n' +
                    '                            <div class="payment-text">\n' +
                    '                                <div class="title">Credit card</div>\n' +
                    '                                <p class="subtitle">end of the number（' + pay_account['last4'] + '）</p>\n' +
                    '                            </div>';
                $(".payment-method-item:first").empty().html(html);
            }

            if (balance != '' && parseFloat(balance) >= parseFloat(total_price)) {
                let balance_html = '<div class="radio1"><input type="radio" name="paymentMethod" value="2" /><span class="icon"></span></div>\n' +
                    '                            <img class="payment-icon" src="/static/home/img/icon-pme2.png" />\n' +
                    '                            <div class="payment-text">\n' +
                    '                                <div class="title">Account balances</div>\n' +
                    '                                <p class="subtitle">account balances：$<span id="user-fee">' + balance + '</span> left.</p>\n' +
                    '                            </div>';
                $(".payment-method-item:last").empty().html(balance_html);
            } else {
                let balance_html = '<span>Account balance is empty or insufficient!</span>';
                $(".payment-method-item:last").empty().html(balance_html);
            }
            if (data.data['pay_account'] == '' && (balance == '' || parseFloat(balance) < parseFloat(total_price))) {
                $(".btn-pay-now").attr('disabled', true);
            } else {
                $(".btn-pay-now").attr('disabled', false);
            }

            $("#paymentMethodEditModal").show();
        }
    });
}

function datepickerToTimestamp(orgin, to_id) {
    var date = $(orgin);
    var to = $('#' + to_id);

    var timeStamp = new Date(date.val()).getTime() / 1000;

    to.val(timeStamp);
}

function check() {
    if ($('#end_time').val() < $('#start_time').val()) {
        alert('time is wrong');
        return false;
    }
    return true;
}

function sync_order() {
    $(".sync-btn").addClass("btn-loading");
    $(".sync-btn").attr('disabled',true);

    $.ajax({
        type:"post",
        url: "/order/drop/syncOrder",
        data: {
        },
        //回调函数
        success: function(data) {
            setTimeout(function () {
                $(".sync-btn").removeClass("btn-loading");
                $(".sync-btn").attr('disabled',false);
                location.reload();
            }, 3000);
        }
    });
    return false;
}

function alertMsg(data)
{
    $.amaran({
        'message'   : data.msg,
        'position'  :'bottom right'
    });
}