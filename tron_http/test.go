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




	address := "a0b4750e2cd76e19dca331bf5d089b71c3c2798548"
	addressBytes, _ := hex.DecodeString(address)
	hex.EncodeToString(addressBytes)
	encodeString := base64.StdEncoding.EncodeToString(addressBytes)
	fmt.Printf("address->base64: %s %s\n", address, encodeString)

	base64Address := "oLR1DizXbhncozG/XQibccPCeYVI"
	unBase64Bytes, _ := base64.StdEncoding.DecodeString(base64Address)
	toAddress := hex.EncodeToString(unBase64Bytes)
	fmt.Printf("base64->address: %s %s\n", base64Address, toAddress)



}