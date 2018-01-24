Tyk/LDAP Go plugin
==================

Overview
--------

This is a simple Go pluging for Tyk providing connection to an LDAP server.
This allows to retreive meta-data from an LDAP server and add it to a Header in the query.
Of course, you need an already configured LDAP server.

Build
-----

Build the server:

    go build

Build the Tyk bundle:

    tyk-cli bundle build -y -o ldap-bundle.zip

Configuration
-------------

This plugin uses the _per API_ `config_data` parameter in order to be configured and environment variables for sensible informations.

Example of `config_data` section in a Tyk API definition:

    "custom_middleware_bundle": "ldap-bundle.zip",
    "config_data": {
        "headername": "X-Ldap-Result",
        "basedn": "dc=exampl.dc=com",
        "filter": "(cn=$ldapparam)",
        "attributes": ["meta"]
    }

* headername: Name of the header to write the result to
* basedn: Start point of the search
* filter: Filter to apply. `$ldapparam` is replaced by the content of the `X-Ldap-Param` HTTP header. For example it can be used to give the `cn` of the object we are looking for
* attributes: List of attributes to retreive from the LDAP response

Environment variables:

* TYKLDAPHOST: Hostname or IP address of the LDAP server
* TYKLDAPPORT: TCP port of the LDAP server
* TYKLDAPBINDDN, TYKLDABINDPW: user/password to connect to the LDAP server

Example
-------

    TYKLDAPHOST=127.0.0.1 TYKLDAPPORT=389 TYKLDAPBINDDN="cn=admin,dc=example,dc=com" TYKLDAPBINDPW=rootpw ./tyk-plugin-ldap-go

    curl -H 'X-Ldap-Param: athing' -H 'Content-Type: application/json' -d '{"h":"v"}' -s -v http://localhost:8080/post|jq .
