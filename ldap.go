package main

import (
	"fmt"
	"log"

	ldap "gopkg.in/ldap.v2"
)

var LdapConn *ldap.Conn

func ldapBind(ldapHost, ldapPort, ldapBindDn, ldapBindPw string) {
	var err error

	LdapConn, err = ldap.Dial("tcp", fmt.Sprintf("%s:%s", ldapHost, ldapPort))
	if err != nil {
		log.Fatal(err)
	}

	err = LdapConn.Bind(ldapBindDn, ldapBindPw)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("=== Connected to LDAP server ===")
}

func ldapSearch(conn *ldap.Conn, baseDn string, filter string, attributes []string) ([]string, error) {
	searchRequest := ldap.NewSearchRequest(
		baseDn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,     // The filter to apply
		attributes, // A list attributes to retrieve
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entry := range sr.Entries {
		result = append(result, entry.GetAttributeValue(attributes[0]))
	}

	return result, nil
}
