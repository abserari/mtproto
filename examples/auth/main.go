package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/xelaj/go-dry"
	"github.com/xelaj/mtproto"
	"github.com/xelaj/mtproto/telegram"
	"golang.org/x/net/proxy"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("second argument must be phone number!")
		os.Exit(1)
	}
	phoneNumber := os.Args[1]

	// edit these params for you!
	client, err := telegram.NewClient(telegram.ClientConfig{
		// where to store session configuration. must be set
		SessionFile: "./session1.json",
		// host address of mtproto server. Actually, it can'be mtproxy, not only official
		ServerHost: "149.154.167.50:443",
		// public keys file is patrh to file with public keys, which you must get from https://my.telelgram.org
		PublicKeysFile: "./keys.pem",
		AppID:          94575,                              // app id, could be find at https://my.telegram.org
		AppHash:        "a3406de8d171bb422bb6ddf3bbd800e2", // app hash, could be find at https://my.telegram.org
		ProxyDialer: telegram.NewSOCK5ProxyDialer("tcp",
			"106.53.131.105:65431",
			&proxy.Auth{
				User:     "Kassulke8264",
				Password: "wFKo1z8xOr",
			}),
	})
	dry.PanicIfErr(err)

	setCode, err := client.AuthSendCode(
		phoneNumber, 94575, "a3406de8d171bb422bb6ddf3bbd800e2", &telegram.CodeSettings{},
	)
	dry.PanicIfErr(err)
	pp.Println(setCode)

	fmt.Print("Auth code:")
	code, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	code = strings.ReplaceAll(code, "\n", "")

	auth, err := client.AuthSignIn(
		phoneNumber,
		setCode.PhoneCodeHash,
		code,
	)
	if err == nil {
		pp.Println(auth)

		fmt.Println("Success! You've signed in!")
		return
	}

	// if you don't have password protection — THAT'S ALL! You're already logged in.
	// but if you have 2FA, you need to make few more steps:

	// could be some errors:

	mtError, ok := errors.Unwrap(err).(*mtproto.ErrResponseCode)
	// SESSION_PASSWORD_NEEDED says that your account has 2FA protection
	if !ok || mtError.Message != "SESSION_PASSWORD_NEEDED" {
		fmt.Println("SignIn failed:", err)
		println("\n\nMore info about error:")
		pp.Println(err)
		return
	}

	fmt.Print("Password:")
	password, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	password = strings.ReplaceAll(password, "\n", "")

	accountPassword, err := client.AccountGetPassword()
	dry.PanicIfErr(err)

	// GetInputCheckPassword is fast response object generator
	inputCheck, err := telegram.GetInputCheckPassword(password, accountPassword)
	dry.PanicIfErr(err)

	auth, err = client.AuthCheckPassword(inputCheck)
	dry.PanicIfErr(err)

	pp.Println(auth)
	fmt.Println("Success! You've signed in!")
}
