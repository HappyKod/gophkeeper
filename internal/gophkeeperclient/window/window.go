package window

import (
	"context"
	"crypto/aes"
	"encoding/json"
	"errors"
	"strings"
	"time"

	clientmodels "yudinsv/gophkeeper/internal/gophkeeperclient/models"
	"yudinsv/gophkeeper/internal/gophkeeperclient/service"
	"yudinsv/gophkeeper/internal/keeperstorage"
	"yudinsv/gophkeeper/internal/models"
	"yudinsv/gophkeeper/internal/utils"

	"github.com/google/uuid"
	"github.com/pterm/pterm"
)

// RunWindow entry point to the app interface
func RunWindow(serviceClient service.ClientService, storage keeperstorage.KeeperStorage) {
	addSecret := "add secret"
	viewSecret := "view secrets"
	registration := "registration"
	authorization := "authorization"
	var options []string
	options = append(options, registration)
	options = append(options, authorization)

	selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()
	pterm.Info.Printfln("Selected: %s", pterm.Green(selectedOption))
	var user models.User
	var secretKey string
	if selectedOption == registration {
		user = registrationWindow()
		err := serviceClient.RegistryService.Register(user)
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		pterm.Info.Printfln("Registration successful")
		secretKey = uuid.New().String()
		pterm.Info.Printfln("Your private key to encrypt and decrypt data: %s", secretKey)

	} else if selectedOption == authorization {
		user = authorizationWindow()
		err := serviceClient.AuthService.Authorization(user)
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		secretKey, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Enter private key: ").WithMultiLine(false).Show()
		_, err = uuid.Parse(secretKey)
		if err != nil {
			pterm.Error.Println("Invalid private key")
			return
		}
		pterm.Info.Printfln("Authorization successful")
	} else {
		pterm.Error.Println("Invalid option")
	}
	go serviceClient.SyncService.StartSync(user.Login)
	for {
		var optionsMenu []string
		optionsMenu = append(optionsMenu, addSecret)
		optionsMenu = append(optionsMenu, viewSecret)
		selectedMenu, _ := pterm.DefaultInteractiveSelect.WithOptions(optionsMenu).Show()
		pterm.Info.Printfln("Selected: %s", pterm.Green(selectedMenu))
		if selectedMenu == addSecret {
			secret, err := addSecretWindow(secretKey)
			if err != nil {
				pterm.Error.Println(err)
				return
			}
			secret.OwnerID = user.Login
			err = storage.PutSecret(context.Background(), secret)
			if err != nil {
				pterm.Error.Println(err)
			}
		} else if selectedMenu == viewSecret {
			var secrets []models.Secret
			syncSecret, err := storage.SyncSecret(context.Background(), user.Login)
			if err != nil {
				return
			}
			for _, liteSecret := range syncSecret {
				secret, err := storage.GetSecret(context.Background(), liteSecret.ID)
				if err != nil {
					pterm.Error.Println(err)
					return
				}
				if !secret.IsDeleted {
					secrets = append(secrets, secret)
				}
			}
			viewSecretWindow(storage, secretKey, secrets)
		}
	}
}

// registrationWindow register user registration rendering
func registrationWindow() models.User {
	username, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter username").WithMultiLine(false).Show()
	password, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter password").WithMultiLine(false).Show()
	return models.User{Login: username, Password: password}
}

// register user authorization rendering
func authorizationWindow() models.User {
	username, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter username").WithMultiLine(false).Show()
	password, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter password").WithMultiLine(false).Show()
	return models.User{Login: username, Password: password}
}

// addSecretWindow adding new models.Secret rendering
func addSecretWindow(secretKey string) (models.Secret, error) {
	typeCards := "bank_cards"
	typePassword := "login_password"
	typeBinary := "binary"
	typeText := "text"

	var options []string
	options = append(options, typeCards)
	options = append(options, typePassword)
	options = append(options, typeBinary)
	options = append(options, typeText)

	selectedOption, _ := pterm.DefaultInteractiveSelect.WithDefaultText("Please secret type secret").WithOptions(options).Show()
	pterm.Info.Printfln("Selected: %s", pterm.Green(selectedOption))
	var data []byte
	var err error
	if selectedOption == typePassword {
		data, err = addLoginPasswordWindow()
		if err != nil {
			return models.Secret{}, err
		}
	} else if selectedOption == typeCards {
		data, err = addBankCardsWindow()
		if err != nil {
			return models.Secret{}, err
		}
	} else if selectedOption == typeText {
		data = addTextWindow()
	} else if selectedOption == typeBinary {
		data = addBinaryWindow()
	} else {
		return models.Secret{}, errors.New("invalid option")
	}
	description, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter Description:").WithMultiLine(false).Show()
	keyStr := []byte(secretKey)
	iv := keyStr[len(keyStr)-aes.BlockSize:]
	key := keyStr[:aes.BlockSize]
	decrypt, err := utils.Encrypt(data, key, iv)
	if err != nil {
		return models.Secret{}, err
	}
	return models.Secret{
		ID:          uuid.New(),
		Value:       decrypt,
		Description: description,
		Type:        selectedOption,
		Ver:         time.Now(),
	}, nil
}

// addLoginPasswordWindow adding new LoginPassword rendering
func addLoginPasswordWindow() ([]byte, error) {
	login, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter login").WithMultiLine(false).Show()
	password, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter password").WithMultiLine(false).Show()
	return json.Marshal(clientmodels.LoginPassword{Login: login, Password: password})
}

// addBankCardsWindow adding new BankCard rendering
func addBankCardsWindow() ([]byte, error) {
	cardnumber, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter card number").WithMultiLine(false).Show()
	expirym, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter expiry month").WithMultiLine(false).Show()
	expiryy, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter expiry year").WithMultiLine(false).Show()
	cvv, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter CVV").WithMultiLine(false).Show()
	cardholdername, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter cardHolder name").WithMultiLine(false).Show()
	return json.Marshal(clientmodels.BankCard{
		CardNumber:     cardnumber,
		ExpiryMonth:    expirym,
		ExpiryYear:     expiryy,
		CVV:            cvv,
		CardHolderName: cardholdername,
	})
}

// addTextWindow adding new Text rendering
func addTextWindow() []byte {
	text, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter text").WithMultiLine(false).Show()
	return []byte(text)
}

// addBinaryWindow adding new Binary rendering
func addBinaryWindow() []byte {
	text, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter text binary").WithMultiLine(false).Show()
	return []byte(text)
}

// getSecret get secret from secret store
func viewSecretWindow(storage keeperstorage.KeeperStorage, secretKey string, secrets []models.Secret) {
	var viewSecrets []string
	for _, v := range secrets {
		viewSecrets = append(viewSecrets, v.ID.String()+"\t"+v.Description)
	}

	selectedOption, _ := pterm.DefaultInteractiveSelect.WithDefaultText("Please select a secret").WithOptions(viewSecrets).Show()
	pterm.Info.Printfln("Selected: %s", pterm.Green(selectedOption))
	uuidStr := strings.Split(selectedOption, "\t")[0]
	for _, s := range secrets {
		if uuidStr == s.ID.String() {
			keyStr := []byte(secretKey)
			iv := keyStr[len(keyStr)-aes.BlockSize:]
			key := keyStr[:aes.BlockSize]
			data, err := utils.Decrypt(s.Value, key, iv)
			if err != nil {
				pterm.Error.Println(err)
			}
			pterm.Info.Printfln("Secret: %s", string(data))
			oneSecretWindow(storage, s)
			return
		}
	}
}

// oneSecretWindow selected option models.Secret
func oneSecretWindow(storage keeperstorage.KeeperStorage, secret models.Secret) {
	closeOp := "close"
	changeOp := "change description"
	deleteOp := "delete"
	var options []string
	options = append(options, closeOp)
	options = append(options, changeOp)
	options = append(options, deleteOp)
	selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()
	if selectedOption == closeOp {
		return
	} else if selectedOption == changeOp {
		description, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter new description: ").WithMultiLine(false).Show()
		secret.Description = description
		secret.Ver = time.Now()
		err := storage.PutSecret(context.Background(), secret)
		if err != nil {
			pterm.Error.Println(err)
		}
	} else if selectedOption == deleteOp {
		err := storage.DeleteSecret(context.Background(), secret.ID)
		if err != nil {
			pterm.Error.Println(err)
		}
	}
}
