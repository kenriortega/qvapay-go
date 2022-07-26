/*
	Package qvapay is the no-official qvapay apps SDK for Go.
	Use it to interact with the Qvapay API.


	func main() {
		api := qvapay.NewPaymentAppClient(
			os.Getenv("APP_ID"),     // app_id
			os.Getenv("APP_SECRET"), // secret_id
			qvapay.BaseURL,          // constants url base https://qvapay.com/api
			false,                   // skip verificationSSL
			nil,                     // custom http.Client
			nil,                     // debug io.Writter (os.Stdout)
		)

		info, err := api.GetInfo(context.Background())
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Println(info)
	}


	Examples can be found at
	https://github.com/kenriortega/qvapay-go/tree/main/examples
	If you find an issue with the SDK, please report through
	https://github.com/kenriortega/qvapay-go/issues/new
*/
package qvapay
