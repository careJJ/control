package util

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

//提出测试API请求
//要检查您的集成是否正常运行，请使用您的测试密钥创建一个Payment API来进行测试API请求。

//我们已使用您的测试秘密API密钥预先填充了此代码示例-只有您可以看到此值。

const (
	stripecarekey="pk_test_TJbZ48ANjDAXvyRiQnCsgAyY00OBMrsdL2"

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
		ReceiptEmail: stripe.String("709412810@me.com"),
	}
	paymentintent.New(params)



}

