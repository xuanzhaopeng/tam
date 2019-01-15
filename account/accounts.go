package account

import (
	"encoding/json"
	"reflect"
	"errors"
	"time"
	"fmt"
	"log"
	"tam/utils"
)

type Accounts struct {
	Data map[string]*Account
}

func BuildAccounts(accountKey string, autoReleaseDuration time.Duration, accountJsonData []byte) (accounts *Accounts, err error) {
	var accountsData []map[string]interface{}
	err = json.Unmarshal(accountJsonData, &accountsData)
	if err != nil {
		return accounts, err
	}

	accounts = &Accounts{Data: make(map[string]*Account)}

	for _, accountData := range accountsData {
		if key, found := accountData[accountKey]; found && reflect.ValueOf(key).Kind() == reflect.String && accounts.Data[key.(string)] == nil {
			account := NewAccount(key.(string), accountData, autoReleaseDuration)
			accounts.Data[key.(string)] = account
		} else {
			return accounts, errors.New("invalid accounts data")
		}
	}

	return accounts, nil
}

func (accounts *Accounts) Fetch(filter Filter) (fetchedAccount map[string]interface{}, err error) {
	log.Printf("[Fetching] By Filter: %s\n", utils.MapToJson(filter))
	for _, account := range accounts.Data {
		if isMatched, err := account.MatchFilter(filter); isMatched {
			if fetchedAccount, err = account.Lock(); err == nil {
				log.Printf("[Fetched] %s\n", utils.MapToJson(fetchedAccount))
				return fetchedAccount, nil
			}
		}
	}
	log.Printf("[Fetched] no available account")
	return nil, errors.New("no available account")
}

func (accounts *Accounts) Release(accountKey string) (err error) {
	if account, found := accounts.Data[accountKey]; found {
		account.Unlock()
		return nil
	} else {
		return errors.New(fmt.Sprintf("cannot find the account %s", accountKey))
	}
}
