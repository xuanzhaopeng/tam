package tests

import (
	"testing"
	"encoding/json"
	"tam/account"
	"time"
	"github.com/stretchr/testify/assert"
	"log"
)

func TestAccountCouldBeAutoUnlocked(t *testing.T) {
	testAccount := getAccount(2 * time.Second)
	testAccount.Lock()
	assert.True(t, testAccount.IsLocked(), "account should be locked")
	time.Sleep(1 * time.Second)
	assert.True(t, testAccount.IsLocked(), "account should be unlocked")
	time.Sleep(3 * time.Second)
	assert.False(t, testAccount.IsLocked(), "account should be unlocked")
}

func TestAccountCannotBeLockedWhenAlreadyLocked(t *testing.T) {
	testAccount := getAccount(10 * time.Second)
	_, err := testAccount.Lock()
	assert.Nil(t, err, "account locked failed")
	data, err := testAccount.Lock()
	assert.Nil(t, data, "account data is nil")
	assert.Error(t, err, "account should be locked failed")
}

func TestAccountCouldBeUnlockedMaully(t *testing.T) {
	testAccount := getAccount(10 * time.Second)
	testAccount.Lock()
	testAccount.Unlock()
	assert.False(t, testAccount.IsLocked())
	testAccount.Lock()
	assert.True(t, testAccount.IsLocked())
}

func TestAccountCouldBeMatchedBySingleFilter(t *testing.T) {
	jsonData := `{
		"username": "test1",
		"password": "password"
    }`
    filterData := `{
		"username": "^test\\d+$"
	}`
    testAccount := getCusAccount(jsonData, 5 * time.Second)
    filter := getFilter(filterData)

    result, err := testAccount.MatchFilter(filter)
    assert.Nil(t, err, "match filter failed")
    assert.True(t, result, "match failed")
}

func TestAccountCouldBeMatchedByMultipleFilters(t *testing.T) {
	jsonData := `{
		"username": "test1",
		"password": "hello"
    }`
	filterData := `{
		"username": "^test\\d+$",
		"password": "^he.*"
	}`
	testAccount := getCusAccount(jsonData, 5 * time.Second)
	filter := getFilter(filterData)

	result, err := testAccount.MatchFilter(filter)
	assert.Nil(t, err, "match filter failed")
	assert.True(t, result, "match failed")
}

func TestAccountCouldBeMatchedInDeepLevel(t *testing.T) {
	jsonData := `{
		"username": "test1",
		"password": "hello",
		"data": {
			"a1": 1,
			"a2": "deep1"
		}
    }`
    filterData := `{
		"username": "^test\\d+",
		"data": {
			"a1": 1
		}
	}`
	testAccount := getCusAccount(jsonData, 5 * time.Second)
	filter := getFilter(filterData)
	result, err := testAccount.MatchFilter(filter)
	assert.Nil(t, err, "match filter failed")
	assert.True(t, result, "match failed")
}

func TestAccountCouldBeMatchAny(t *testing.T) {
	filterData := `{}`
	testAccount := getAccount(5 * time.Second)
	filter := getFilter(filterData)
	result, err := testAccount.MatchFilter(filter)
	assert.Nil(t, err, "match filter failed")
	assert.True(t, result, "match failed")
}


func getAccount(duration time.Duration) *account.Account {
	var accountData map[string]interface{}
	jsonData := `{
		"username": "abc",
		"password": "pass"
	}`
	json.Unmarshal([]byte(jsonData), &accountData)
	return account.NewAccount("username", accountData, duration)
}

func getCusAccount(jsonData string, duration time.Duration) *account.Account {
	var accountData map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &accountData)
	if err != nil {
		log.Fatal(err)
	}
	return account.NewAccount("username", accountData, duration)
}

func getFilter(jsonData string) account.Filter {
	filter, err := account.BuildFilter([]byte(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	return filter
}
