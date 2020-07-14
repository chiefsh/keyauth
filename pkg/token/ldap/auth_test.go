package ldap_test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/infraboard/keyauth/pkg/token/ldap"
)

func TestAuthenticate(t *testing.T) {
	if testConfig.Server == "" {
		t.Skip("LDAP_SERVER not set")
		return
	}

	if testConfig.BaseDN == "" {
		t.Skip("LDAP_BASEDN not set")
		return
	}

	config := &ldap.Config{Server: testConfig.Server, Port: testConfig.Port, Security: testConfig.BindSecurity, BaseDN: testConfig.BaseDN}

	fmt.Println(config)
	status, err := ldap.Authenticate(config, "ladp-auth", "invalid password")
	if err != nil {
		t.Fatal("Invalid credentials: Expected err to be nil but got:", err)
	}
	if status {
		t.Error("Invalid credentials: Expected authenticate status to be false")
	}

	badConfig := &ldap.Config{Server: testConfig.Server, Port: testConfig.Port, Security: testConfig.BindSecurity, BaseDN: "Bad BaseDN"}
	if _, err = ldap.Authenticate(badConfig, "ladp-auth", "invalid password"); !strings.Contains(err.Error(), "invalid BaseDN") {
		t.Error("Invalid configuration: Expected invalid BaseDN error but got:", err)
	}

	badConfig = &ldap.Config{Server: "127.0.0.1", Port: 1, Security: testConfig.BindSecurity, BaseDN: testConfig.BaseDN}
	if _, err = ldap.Authenticate(badConfig, "ladp-auth", "invalid password"); !strings.Contains(err.Error(), "Connection error") {
		t.Error("Connect error: Expected connection error but got:", err)
	}

	if testConfig.BindUPN == "" || testConfig.BindPass == "" {
		t.Skip("LDAP_BIND_UPN or LDAP_BIND_PASS not set")
		return
	}

	status, err = ldap.Authenticate(config, testConfig.BindUPN, testConfig.BindPass)
	if err != nil {
		t.Fatal("Valid UPN: Expected err to be nil but got:", err)
	}
	if !status {
		t.Error("Valid UPN: Expected authenticate status to be true")
	}

	var username string

	if splits := strings.Split(testConfig.BindUPN, "@"); len(splits) == 2 {
		username = splits[0]
	} else {
		t.Fatalf("Expected BIND_UPN (%s) to be splittable", testConfig.BindUPN)
	}

	status, err = ldap.Authenticate(config, username, testConfig.BindPass)
	if err != nil {
		t.Fatal("Valid username: Expected err to be nil but got:", err)
	}
	if !status {
		t.Error("Valid username: Expected authenticate status to be true")
	}
}

func TestAuthenticateExtended(t *testing.T) {
	if testConfig.Server == "" {
		t.Skip("LDAP_SERVER not set")
		return
	}

	if testConfig.BaseDN == "" {
		t.Skip("LDAP_BASEDN not set")
		return
	}

	config := &ldap.Config{Server: testConfig.Server, Port: testConfig.Port, Security: testConfig.BindSecurity, BaseDN: testConfig.BaseDN}

	status, _, _, err := ldap.AuthenticateExtended(config, "ladp-auth", "invalid password", []string{""}, nil)
	if err != nil {
		t.Fatal("Invalid credentials: Expected err to be nil but got:", err)
	}
	if status {
		t.Error("Invalid credentials: Expected authenticate status to be false")
	}

	badConfig := &ldap.Config{Server: testConfig.Server, Port: testConfig.Port, Security: testConfig.BindSecurity, BaseDN: "Bad BaseDN"}
	if _, _, _, err = ldap.AuthenticateExtended(badConfig, "ladp-auth", "invalid password", []string{""}, nil); !strings.Contains(err.Error(), "invalid BaseDN") {
		t.Error("Invalid configuration: Expected invalid BaseDN error but got:", err)
	}

	badConfig = &ldap.Config{Server: "127.0.0.1", Port: 1, Security: testConfig.BindSecurity, BaseDN: testConfig.BaseDN}
	if _, _, _, err = ldap.AuthenticateExtended(badConfig, "ladp-auth", "invalid password", []string{""}, nil); !strings.Contains(err.Error(), "Connection error") {
		t.Error("Connect error: Expected connection error but got:", err)
	}

	if testConfig.BindUPN == "" || testConfig.BindPass == "" {
		t.Skip("LDAP_BIND_UPN or LDAP_BIND_PASS not set")
		return
	}

	status, _, _, err = ldap.AuthenticateExtended(config, testConfig.BindUPN, testConfig.BindPass, []string{""}, nil)
	if err != nil {
		t.Fatal("Valid UPN: Expected err to be nil but got:", err)
	}
	if !status {
		t.Error("Valid UPN: Expected authenticate status to be true")
	}

	var username string

	if splits := strings.Split(testConfig.BindUPN, "@"); len(splits) == 2 {
		username = splits[0]
	} else {
		t.Fatalf("Expected BIND_UPN (%s) to be splittable", testConfig.BindUPN)
	}

	status, _, _, err = ldap.AuthenticateExtended(config, username, testConfig.BindPass, []string{""}, nil)
	if err != nil {
		t.Fatal("Valid username: Expected err to be nil but got:", err)
	}
	if !status {
		t.Error("Valid username: Expected authenticate status to be true")
	}

	status, entry, _, err := ldap.AuthenticateExtended(config, testConfig.BindUPN, testConfig.BindPass, []string{"memberOf"}, nil)
	if err != nil {
		t.Fatal("memberOf attrs: Expected err to be nil but got:", err)
	}
	if !status {
		t.Error("memberOf attrs: Expected authenticate status to be true")
	}

	//use dn for even groups and cn for odd groups
	dnGroups := entry.GetAttributeValues("memberOf")
	var checkGroups []string
	for i, group := range dnGroups {
		if i%2 == 0 {
			checkGroups = append(checkGroups, group)
		} else {
			cn := dnToCN(group)
			if cn != "" {
				checkGroups = append(checkGroups, cn)
			}
		}
	}

	status, entry, userGroups, err := ldap.AuthenticateExtended(config, testConfig.BindUPN, testConfig.BindPass, []string{"sAMAccountName"}, checkGroups)
	if err != nil {
		t.Fatal("memberOf attrs: Expected err to be nil but got:", err)
	}
	if !status {
		t.Error("memberOf attrs: Expected authenticate status to be true")
	}

	sort.Strings(checkGroups)
	sort.Strings(userGroups)

	if len(checkGroups) != len(userGroups) {
		t.Fatalf("Expected returned group count (%d) to be equal to searched group count (%d)", len(userGroups), len(checkGroups))
	}

	for i := range checkGroups {
		if checkGroups[i] != userGroups[i] {
			t.Fatalf("Expected returned group (%s) to be equal to searched group (%s):", userGroups[i], checkGroups[i])
		}
	}

	if entry.GetAttributeValue("sAMAccountName") != username {
		t.Fatalf("Expected sAMAccountName (%s) to be equal to username (%s)", entry.GetAttributeValue("sAMAccountName"), username)
	}
}
