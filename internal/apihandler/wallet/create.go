package wallet

import (
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/wallet"
)

func NewCreate(walletStorage *sync.Map) operations.MgmtWalletCreateHandler {
	return &walletCreate{walletStorage: walletStorage}
}

type walletCreate struct {
	walletStorage *sync.Map
}

func (c *walletCreate) Handle(params operations.MgmtWalletCreateParams) middleware.Responder {
	if params.Body.Nickname == nil || len(*params.Body.Nickname) == 0 {
		e := errorCodeWalletCreateNoNickname
		msg := "Error: nickname field is mandatory."

		return operations.NewMgmtWalletCreateBadRequest().WithPayload(
			&models.Error{
				Code:    &e,
				Message: &msg,
			})
	}

	_, ok := c.walletStorage.Load(*params.Body.Nickname)
	if ok {
		e := errorCodeWalletAlreadyExists
		msg := "Error: a wallet with the same nickname already exists."

		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    &e,
				Message: &msg,
			})
	}

	if params.Body.Password == nil || len(*params.Body.Password) == 0 {
		e := errorCodeWalletCreateNoPassword
		msg := "Error: password field is mandatory."

		return operations.NewMgmtWalletCreateBadRequest().WithPayload(
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

	c.walletStorage.Store(newWallet.Nickname, newWallet)

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
