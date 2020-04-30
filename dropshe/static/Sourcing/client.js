// A reference to Stripe.js initialized with a fake API key.
//Sign in to see examples pre-filled with your key.
var stripe = Stripe("pk_test_TJbZ48ANjDAXvyRiQnCsgAyY00OBMrsdL2");

// Disable the button until we have Stripe set up on the page
document.getElementById("button-text").disabled = true;
function Check() {
  var sourceidList = []
  var Check = $("table input[type=checkbox]:checked");//在table中找input下类型为checkbox属性为选中状态的数据
  Check.each(function () {//遍历
    var row = $(this).parent("td").parent("tr");//获取选中行
    var ida = row.find(".sourcing_id").text();//获取name='skuid'的值
    sourceidList.push(eval(ida))


  })

  $.ajax({
    url:"/sourcing/pay",
    data:{sourceidList},
    type:"post",
    contentType:'application/x-www-form-urlencoded',
    success:function(data){
      // 传sourceidList给后台获取clientSecret
      if( data.errno!=10){
        alert(data.errmsg)
        location.reload();

                       }
      $('#payment-form').css(
        'display','block'
      )
      var elements = stripe.elements(
          {        //设置默认显示语种   en 英文 cn 中文 auto 自动获取语种
            locale: 'en'
          }
      );
      var style = {
        base: {
          color: "#32325d",
          fontFamily: 'Arial, sans-serif',
          fontSmoothing: "antialiased",
          fontSize: "16px",
          "::placeholder": {
            color: "#32325d"
          }
        },
        invalid: {
          fontFamily: 'Arial, sans-serif',
          color: "#fa755a",
          iconColor: "#fa755a"
        }
      };
      var card = elements.create("card", { style: style });
      // Stripe injects an iframe into the DOM
      card.mount("#card-element");

      $('#button-text').html('pay  $ '+data.payMomeny)
      card.addEventListener('change', ({error}) => {
        const displayError = document.getElementById('card-errors');
        if (error) {
          displayError.textContent = error.message;
        } else {
          displayError.textContent = '';
        }
      });
      card.on("change", function (event) {

        // Disable the Pay button if there are no card details in the Element
        document.querySelector("button").disabled = event.empty;
        document.querySelector("#card-errors").textContent = event.error ? event.error.message : "";
      });

      setTimeout(function(){
        var form = document.getElementById("payment-form");
        form.addEventListener("submit", function(event) {
          event.preventDefault();
          // Complete payment when the submit button is clicked
          payWithCard(stripe, card, data.clientSecret,sourceidList);
        });

      },2000)
    }

  })

}


// Calls stripe.confirmCardPayment
// If the card requires authentication Stripe shows a pop-up modal to
// prompt the user to enter authentication details without leaving your page.
var payWithCard = function(stripe, card, clientSecret,sourceidList) {
  loading(true);
  stripe
    .confirmCardPayment(clientSecret, {
      payment_method: {
        card: card
      }
    })
    .then(function(result) {
      if (result.error) {
        //支付失败
       alert('Payment failed, please enter the correct information')
        location.reload();
        // Show error to your customer
        showError(result.error.message);
      } else {
        //支付成功
        // The payment succeeded!
        console.log(sourceidList,'sourceidList')

        $.post('/pay/after',{sourceidList},function(res){})
        alert(' The payment succeeded')
        orderComplete(result.paymentIntent.id);
        $('#payment-form').css(
            'display','none'
        )
        location.reload();
      }
    });
};

/* ------- UI helpers ------- */

// Shows a success message when the payment is complete
var orderComplete = function(paymentIntentId) {
  loading(false);
  document
    .querySelector(".result-message a")
    .setAttribute(
      "href",
      "https://dashboard.stripe.com/test/payments/" + paymentIntentId
    );
  document.querySelector(".result-message").classList.remove("hidden");
  document.querySelector("button").disabled = true;
};

// Show the customer the error from Stripe if their card fails to charge
var showError = function(errorMsgText) {
  loading(false);
  var errorMsg = document.querySelector("#card-errors");
  errorMsg.textContent = errorMsgText;
  setTimeout(function() {
    errorMsg.textContent = "";
  }, 4000);
};

// Show a spinner on payment submission
var loading = function(isLoading) {
  if (isLoading) {
    // Disable the button and show a spinner
    document.querySelector("button").disabled = true;
    document.querySelector("#spinner").classList.remove("hidden");
    document.querySelector("#button-text").classList.add("hidden");
  } else {
    document.querySelector("button").disabled = false;
    document.querySelector("#spinner").classList.add("hidden");
    document.querySelector("#button-text").classList.remove("hidden");
  }
};
