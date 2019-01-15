package account

import (
	"sync"
	"time"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"log"
)


type Account struct {
	Key string
	Data map[string]interface{}
	duration time.Duration
	mutex  *sync.Mutex
	isLocked bool
}

func NewAccount(key string, data map[string]interface{}, duration time.Duration) *Account {
	return &Account{Key:key, Data:data, duration:duration, mutex: &sync.Mutex{}, isLocked: false}
}

func (account *Account) Lock() (map[string]interface{}, error) {
	if account.IsLocked() {
		return nil, errors.New(fmt.Sprintf("the account %s has been already locked", account.Key))
	}

	account.mutex.Lock()
	defer account.mutex.Unlock()
	account.isLocked = true
	log.Printf("[Lock] %s\n", account.Key)
	go func(account *Account) {
		timer := time.NewTimer(account.duration)
		<- timer.C
		if account.isLocked {
			account.Unlock()
		}
	}(account)
	return account.Data, nil
}

func (account *Account) Unlock() {
	account.mutex.Lock()
	account.isLocked = false
	log.Printf("[UnLock] %s\n", account.Key)
	account.mutex.Unlock()
}

func (account *Account) IsLocked() bool {
	account.mutex.Lock()
	defer account.mutex.Unlock()
	return account.isLocked
}

func (account *Account) MatchFilter(filter Filter) (bool,error) {
	return account.matchFilter(account.Data, filter)
}

func (account *Account) matchFilter(data map[string]interface{}, filter Filter) (bool, error) {
	if len(filter) <= 0 {
		return true, nil
	}

	for k, v := range filter {
		_data := data[k]
		log.Printf("[Matching] compare data: {%s}, filter: {%s}", _data, v)
		if _data == nil {
			return false, nil
		}
		dataType := reflect.ValueOf(_data).Kind()
		filterType := reflect.ValueOf(v).Kind()
		if dataType == reflect.Map && filterType == reflect.Map {
			result, err := account.matchFilter(_data.(map[string]interface{}), v.(map[string]interface{}))
			if err != nil || !result {
				return false, err
			}
		} else if (dataType == reflect.Map && filterType != reflect.Map) || (dataType != reflect.Map && filterType == reflect.Map) {
			return false, nil
		} else if dataType == filterType {
			switch dataType {
			case reflect.String:
				dataValue := _data.(string)
				filterValue := v.(string)
				r, _ := regexp.Compile(filterValue)
				if !r.MatchString(dataValue) {
					log.Printf("[Match] compare failed: data: {%s}, filter: {%s}", dataValue, filterValue)
					return false, nil
				}
			case reflect.Int:
				dataValue := _data.(int)
				filterValue := v.(int)
				if dataValue != filterValue {
					log.Printf("[Match] compare failed: data: {%s}, filter: {%s}", dataValue, filterValue)
					return false, nil
				}
			case reflect.Float32:
				dataValue := _data.(float32)
				filterValue := v.(float32)
				if dataValue != filterValue {
					log.Printf("[Match] compare failed: data: {%s}, filter: {%s}", dataValue, filterValue)
					return false, nil
				}
			case reflect.Float64:
				dataValue := _data.(float64)
				filterValue := v.(float64)
				if dataValue != filterValue {
					log.Printf("[Match] compare failed: data: {%s}, filter: {%s}", dataValue, filterValue)
					return false, nil
				}
			default:
				return false, errors.New("cannot identify the type of values, it only support string, int, float32, float64")
			}
		} else {
			log.Printf("[Match] compare failed: types are different")
			return false, nil
		}
	}

	log.Printf("[Match] matched!")
	return true, nil
}
