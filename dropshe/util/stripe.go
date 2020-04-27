package util

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/client"
	"github.com/stripe/stripe-go/paymentintent"
)

//提出测试API请求
//要检查您的集成是否正常运行，请使用您的测试密钥创建一个Payment API来进行测试API请求。

//我们已使用您的测试秘密API密钥预先填充了此代码示例-只有您可以看到此值。

const (
	stripecarekey="pk_test_TJbZ48ANjDAXvyRiQnCsgAyY00OBMrsdL2"
	stripecarekey2="sk_test_TyAe7P109fmUiQBHTf0xGn9G00EaocYpqU"
)

func StripePay(){
	//设置你的密匙。记得在生产中切换到您的live密钥!
	//点击这里查看你的密钥:https://dashboard.stripe.com/account/apikeys
	stripe.Key = stripecarekey

	params := &stripe.PaymentIntentParams{
		Amount: stripe.Int64(1000),
		Currency: stripe.String(string(stripe.CurrencyHKD)),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),

	}
	paymentintent.New(params)



}
func Stripesk(){
	stripe.Key=stripecarekey2
	sc:=&client.API{}
	sc.Init("sk_test_TyAe7P109fmUiQBHTf0xGn9G00EaocYpqU", nil)
	params := &stripe.ChargeParams{}
	sc.Charges.Get("ch_1GVv6zEVqDizVrmQBqkWmzGG", params)
	//账户id,一般以acct_开头
	params.SetStripeAccount("acct_1GJYX4EVqDizVrmQ")
	ch, err := charge.Get("ch_1GVv6zEVqDizVrmQBqkWmzGG", params)
beego.Alert(ch.ID)
	//错误处理
	//_, err :=
		// Go library call

	if err != nil {
		// Try to safely cast a generic error to a stripe.Error so that we can get at
		// some additional Stripe-specific information about what went wrong.
		if stripeErr, ok := err.(*stripe.Error); ok {
			// The Code field will contain a basic identifier for the failure.
			switch stripeErr.Code {
			case stripe.ErrorCodeCardDeclined:
			case stripe.ErrorCodeExpiredCard:
			case stripe.ErrorCodeIncorrectCVC:
			case stripe.ErrorCodeIncorrectZip:
				// etc.
			}

			// The Err field can be coerced to a more specific error type with a type
			// assertion. This technique can be used to get more specialized
			// information for certain errors.
			if cardErr, ok := stripeErr.Err.(*stripe.CardError); ok {
				fmt.Printf("Card was declined with code: %v\n", cardErr.DeclineCode)
			} else {
				fmt.Printf("Other Stripe error occurred: %v\n", stripeErr.Error())
			}
		} else {
			fmt.Printf("Other error occurred: %v\n", err.Error())
		}
	}

}


func StripeTest(i int64)*stripe.PaymentIntent{
	stripe.Key = "sk_test_TyAe7P109fmUiQBHTf0xGn9G00EaocYpqU"

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(i),
		Currency: stripe.String(string(stripe.CurrencyHKD)),

	}
	// Verify your integration in this guide by including this parameter
	params.AddMetadata("integration_check", "accept_a_payment")

	pi, _ := paymentintent.New(params)
	return pi
}