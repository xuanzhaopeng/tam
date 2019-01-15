package tests

import (
	"testing"
	"tam/account"
	"github.com/stretchr/testify/assert"
	"time"
)

var (
	defaultDuration = 1 * time.Second
)

func TestLoadAccountsSuccessfully(t *testing.T) {
	jsonStr := `[
		{
			"username": "a1",
			"password": "p1",
			"region": "US"
		},
		{
			"username": "a2",
			"password": "p2",
			"region": "NL"
		}
	]`
	accountKey := "username"
	accounts, err := account.BuildAccounts(accountKey, defaultDuration, []byte(jsonStr))

	assert.Nil(t, err, "The accounts load failed")
	assert.Len(t, accounts.Data, 2, "The accounts length should be 1")
}

func TestLoadAccountsFailedByIncorrectFormat(t *testing.T) {
	jsonStr := `[
		{
			abc: abc
		},
	]`
	_, err := account.BuildAccounts("any", defaultDuration, []byte(jsonStr))
	assert.Error(t, err, "It should raised err")
}

func TestLoadAccountsFailedByIncorrectAccountKey(t *testing.T) {
	jsonStr := `[
		{
			"username": "a1",
			"password": "p1",
			"region": "US"
		}
	]`
	_, err := account.BuildAccounts("any", defaultDuration, []byte(jsonStr))
	assert.Error(t, err, "It should raised err")
}

func TestLoadAccountsFailedByDuplicatedKey(t *testing.T) {
	jsonStr := `[
		{
			"username": "a1",
			"password": "p1",
			"region": "US"
		},
		{
			"username": "a1",
			"password": "p2",
			"region": "NL"
		}
	]`
	_, err := account.BuildAccounts("username", defaultDuration, []byte(jsonStr))
	assert.Error(t, err, "It should raised err")
}

func TestFetchAndReleaseAccountSuccessfully(t *testing.T) {
	jsonStr := `[
		{
			"username": "tester1",
			"password": "p1",
			"region": "US"
		},
		{
			"username": "tester2",
			"password": "p2",
			"region": "US"
		},
		{
			"username": "tester3",
			"password": "p2",
			"region": "NL"
		}
	]`
	filterStr := `{
		"region": "US"
	}`
	accountKey := "username"
	duration := 10 * time.Second
	accounts, _ := account.BuildAccounts(accountKey, duration, []byte(jsonStr))
	filter, _ := account.BuildFilter([]byte(filterStr))
	assert.Len(t, accounts.Data, 3, "load accounts failed")

	// US account1
	data1, _ := accounts.Fetch(filter)
	assert.True(t, len(data1) > 0)

	// US account2
	data2, _ := accounts.Fetch(filter)
	assert.True(t, len(data2) > 0)

	// No available US account
	_, err := accounts.Fetch(filter)
	assert.Error(t, err)

	// Release US account1
	accounts.Release(data1[accountKey].(string))

	// Get US account again
	data4, _ := accounts.Fetch(filter)
	assert.True(t, len(data4) > 0)
}

func TestFetchAndAutoReleaseAccount(t *testing.T) {
	jsonStr := `[
		{
			"username": "tester1",
			"password": "p1",
			"region": "US"
		},
		{
			"username": "tester3",
			"password": "p2",
			"region": "NL"
		}
	]`
	filterStr := `{
		"region": "NL"
	}`
	accountKey := "username"
	duration := 2 * time.Second
	accounts, _ := account.BuildAccounts(accountKey, duration, []byte(jsonStr))
	filter, _ := account.BuildFilter([]byte(filterStr))

	data1, _ := accounts.Fetch(filter)
	assert.True(t, len(data1) > 0)

	time.Sleep(3 * time.Second)

	data2, _ := accounts.Fetch(filter)
	assert.True(t, len(data2) > 0)
}

func TestFetchAnyAccount(t *testing.T) {
	jsonStr := `[
		{
			"username": "tester1",
			"password": "p1",
			"region": "US"
		}
	]`
	filterStr := `{}`

	accountKey := "username"
	duration := 10 * time.Second
	accounts, _ := account.BuildAccounts(accountKey, duration, []byte(jsonStr))
	filter, _ := account.BuildFilter([]byte(filterStr))

	data1, _ := accounts.Fetch(filter)
	assert.True(t, len(data1) > 0)

	_, err := accounts.Fetch(filter)
	assert.Error(t, err)
}