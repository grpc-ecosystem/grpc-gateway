package main

import (
	"fmt"
	"encoding/base64"
	"encoding/hex"
)


func main()  {

/*
	{
		"account_name": "Zion",
		"address": "a0b4750e2cd76e19dca331bf5d089b71c3c2798548",
		"balance": "15000000000000000"
	},
*/

/*
	{
		"account_name": "Wmlvbg==",
		"address": "oLR1DizXbhncozG/XQibccPCeYVI",
		"balance": "15000000000000000"
	},
*/

	fmt.Println(hexstr2base64("a0b4750e2cd76e19dca331bf5d089b71c3c2798548"))
	fmt.Println(base64Tohexstr("oLR1DizXbhncozG/XQibccPCeYVI"))

	fmt.Println(strToBase64("Zion"))
	fmt.Println(base64ToStr("Wmlvbg=="))

}

func hexstr2base64(s string) (string){
	addressBytes, _ := hex.DecodeString(s)
	return base64.StdEncoding.EncodeToString(addressBytes)
}

func base64Tohexstr(s string) (string) {
	unBase64Bytes, _ := base64.StdEncoding.DecodeString(s)
	return hex.EncodeToString(unBase64Bytes)
}

func strToBase64(s string) (string) {
	bytes := []byte(s)
	return base64.StdEncoding.EncodeToString(bytes)
}

func base64ToStr(s string) (string) {
	bytes, _ := base64.StdEncoding.DecodeString(s)
	return string(bytes)
}