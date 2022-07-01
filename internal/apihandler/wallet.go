package apihandler

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/wallet"
)

func WalletCreate(params operations.MgmtWalletCreateParams) middleware.Responder {
	if params.Body.Nickname == nil || len(*params.Body.Nickname) == 0 {
		e := errorCodeWalletCreateNoNickname
		msg := "Error: nickname field is mandatory"

		operations.NewMgmtWalletCreateBadRequest().WithPayload(
			&models.Error{
				Code:    &e,
				Message: &msg,
			})
	}

	if params.Body.Password == nil || len(*params.Body.Password) == 0 {
		e := errorCodeWalletCreateNoPassword
		msg := "Error: password field is mandatory"

		operations.NewMgmtWalletCreateBadRequest().WithPayload(
			&models.Error{
				Code:    &e,
				Message: &msg,
			})
	}

	newWallet, err := wallet.New(*params.Body.Nickname)
	if err != nil {
		e := errorCodeWalletCreateNew
		msg := err.Error()

		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    &e,
				Message: &msg,
			})
	}

	err = newWallet.Protect(*params.Body.Password, 0)
	if err != nil {
		e := errorCodeWalletCreateNew
		msg := err.Error()

		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    &e,
				Message: &msg,
			})
	}

	privK := base58.CheckEncode(newWallet.KeyPairs[0].PrivateKey)
	pubK := base58.CheckEncode(newWallet.KeyPairs[0].PublicKey)
	salt := base58.CheckEncode(newWallet.KeyPairs[0].Salt[:])
	nonce := base58.CheckEncode(newWallet.KeyPairs[0].Nonce[:])

	return operations.NewMgmtWalletCreateOK().WithPayload(
		&models.Wallet{
			Nickname: &newWallet.Nickname,
			Address:  &newWallet.Address,
			KeyPairs: []*models.WalletKeyPairsItems0{{
				PrivateKey: &privK,
				PublicKey:  &pubK,
				Salt:       &salt,
				Nonce:      &nonce,
			}},
		})
}

func WalletImport(params operations.MgmtWalletImportParams) middleware.Responder {
	var err error

	keyPairs := make([]wallet.KeyPair, len(params.Body.KeyPairs))
	for i := 0; i < len(params.Body.KeyPairs); i++ {
		keyPairs[i].PrivateKey, err = base58.CheckDecode(*params.Body.KeyPairs[i].PrivateKey)
		if err != nil {
			return operations.NewMgmtWalletCreateUnprocessableEntity()
		}

		keyPairs[i].PublicKey, err = base58.CheckDecode(*params.Body.KeyPairs[i].PublicKey)
		if err != nil {
			return operations.NewMgmtWalletCreateUnprocessableEntity()
		}

		salt, err := base58.CheckDecode(*params.Body.KeyPairs[i].Salt)
		if err != nil {
			return operations.NewMgmtWalletCreateUnprocessableEntity()
		}

		copy(keyPairs[i].Salt[:], salt)

		nonce, err := base58.CheckDecode(*params.Body.KeyPairs[i].Nonce)
		if err != nil {
			return operations.NewMgmtWalletCreateUnprocessableEntity()
		}

		copy(keyPairs[i].Nonce[:], nonce)

		keyPairs[i].Protected = true
	}

	w := wallet.Wallet{
		Version:  0,
		Nickname: *params.Body.Nickname,
		Address:  *params.Body.Address,
		KeyPairs: keyPairs,
	}

	fmt.Println(w)

	return operations.NewMgmtWalletImportNoContent()
}
