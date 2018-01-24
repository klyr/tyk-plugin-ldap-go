package main

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	coprocess "github.com/TykTechnologies/tyk-protobuf/bindings/go"
)

type LdapConfigData struct {
	HeaderName string   `json:headername`
	BaseDn     string   `json:basedn`
	Filter     string   `json:filter`
	Attributes []string `json:attributes`
}

func HookAddMeta(object *coprocess.Object) (*coprocess.Object, error) {
	log.Println("Hook 'HookAddMeta' called ...")

	ldapParam := object.Request.Headers["X-Ldap-Param"]
	configData, _ := object.Spec["config_data"]
	log.Println("ConfigData", configData)

	var config LdapConfigData

	err := json.Unmarshal([]byte(configData), &config)
	if err != nil {
		return nil, errors.New("Error while reading config data, doing nothing ...")
	}

	log.Printf(`ConfigData\n
============\n
basedn: %s
filter: %s
attributes: %s
header: %s
`, config.BaseDn, config.Filter, config.Attributes, config.HeaderName)

	var filter string
	if len(ldapParam) > 0 {
		filter = strings.Replace(config.Filter, "$ldapparam", ldapParam, -1)
	} else {
		filter = config.Filter
	}
	result, err := ldapSearch(LdapConn, config.BaseDn, filter, config.Attributes)

	object.Request.SetHeaders = map[string]string{
		config.HeaderName: strings.Join(result, ","),
	}

	return object, nil
}
