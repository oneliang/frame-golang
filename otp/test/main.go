package main

import (
	"bufio"
	"fmt"
	myotp "github.com/oneliang/frame-golang/otp"
	"github.com/pquerna/otp"
	"os"
)

func display(key *otp.Key, data []byte) {
	fmt.Printf("Issuer:       %s\n", key.Issuer())
	fmt.Printf("Account Name: %s\n", key.AccountName())
	fmt.Printf("Secret:       %s\n", key.Secret())
	fmt.Println("Writing PNG to qr-code.png....")
	os.WriteFile("qr-code.png", data, 0644)
	fmt.Println("")
	fmt.Println("Please add your TOTP to your OTP Application now!")
	fmt.Println("")
}

func promptForPasscode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Passcode: ")
	text, _ := reader.ReadString('\n')
	return text
}

func main() {
	//key, _ := myotp.GenerateKey("Example.com", "alice@example.com")
	//var buf bytes.Buffer
	//img, err := key.Image(200, 200)
	//if err != nil {
	//	panic(err)
	//}
	//png.Encode(&buf, img)
	//os.WriteFile("qr-code.png", buf.Bytes(), 0644)
	var secret = "Y65R5LXKBWUHUSKXZ2FCCSNSZCJ3EJIH"
	passcode, _ := myotp.GeneratePasscode(secret)
	fmt.Println("passcode:" + passcode)
	fmt.Println("secret:" + secret)
	inputPasscode := promptForPasscode()
	valid := myotp.Validate(inputPasscode, secret)
	if valid {
		println("Valid passcode!")
		os.Exit(0)
	} else {
		println("Invalid passcode!")
		os.Exit(1)
	}
	return
	//key, err := totp.Generate(totp.GenerateOpts{
	//	Issuer:      "Example.com",
	//	AccountName: "alice@example.com",
	//})
	//if err != nil {
	//	panic(err)
	//}
	//// Convert TOTP key into a PNG
	//var buf bytes.Buffer
	//img, err := key.Image(200, 200)
	//if err != nil {
	//	panic(err)
	//}
	//png.Encode(&buf, img)
	//
	//// display the QR code to the user.
	//
	//generate_passcode, _ := totp.GenerateCode(key.Secret(), time.Now())
	//fmt.Println(generate_passcode)
	//
	//display(key, buf.Bytes())
	//
	//// Now Validate that the user's successfully added the passcode.
	//fmt.Println("Validating TOTP...")
	////passcode := promptForPasscode()
	//valid := totp.Validate(generate_passcode, key.Secret())
	//if valid {
	//	println("Valid passcode!")
	//	os.Exit(0)
	//} else {
	//	println("Invalid passcode!")
	//	os.Exit(1)
	//}
}
