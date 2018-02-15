Tyk/LDAP Go plugin
==================

Overview
--------

This is a simple Go plugin for Tyk providing connection to an LDAP server.
This allows to retreive meta-data from an LDAP server and add it to a Header in the query.
Of course, you need an already configured LDAP server.

The Tyk documentation about the same kind of plugin, in Python, is available at: https://tyk.io/docs/customise-tyk/plugins/rich-plugins/python/custom-auth-python-tutorial/

    +---------+       +-----+                           +-------------+                                                          +------------+ +----------------+
    | Client  |       | Tyk |                           | gRPCPlugin  |                                                          | LDAPServer | | UpstreamServer |
    +---------+       +-----+                           +-------------+                                                          +------------+ +----------------+
         |               |                                     |                                                                        |                |
         | HTTP query    |                                     |                                                                        |                |
         |-------------->|                                     |                                                                        |                |
         |               |                                     |                                                                        |                |
         |               | Original query object passed        |                                                                        |                |
         |               |------------------------------------>|                                                                        |                |
         |               |                                     | ----------------------------------------------------------------\      |                |
         |               |                                     |-| get the LDAP cn from the configured header (ex: X-Ldap-Param) |      |                |
         |               |                                     | |---------------------------------------------------------------|      |                |
         |               |                                     |                                                                        |                |
         |               |                                     | LDAP query to fetch the LDAP object                                    |                |
         |               |                                     |----------------------------------------------------------------------->|                |
         |               |                                     |                                                                        |                |
         |               |                                     |                                                                        |                |
         |               |                                     |<-----------------------------------------------------------------------|                |
         |               |                                     | --------------------------------------------------------------\        |                |
         |               |                                     |-| update the original query HTTP headers with the LDAP result |        |                |
         |               |                                     | |-------------------------------------------------------------|        |                |
         |               |                                     |                                                                        |                |
         |               |                                     |                                                                        |                |
         |               |<------------------------------------|                                                                        |                |
         |               |                                     |                                                                        |                |
         |               |                                     |                                                                        |                |
         |               |------------------------------------------------------------------------------------------------------------------------------>|
         |               |                                     |                                                                        |                |
         |               |                                     |                                                                        |                |
         |<--------------|                                     |                                                                        |                |
         |               |                                     |                                                                        |                |


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
        "basedn": "dc=example.dc=com",
        "filter": "(cn=$ldapparam)",
        "attributes": ["meta"]
    }

* headername: Name of the header to write the result to
* basedn: Start point of the search
* filter: Filter to apply. `$ldapparam` is replaced by the content of the `X-Ldap-Param` HTTP header. For example it can be used to give the `cn` of the object we are looking for
* attributes: List of attributes to retreive from the LDAP response


Update the main Tyk configuration (tyk.conf) to enable the gRPC server:

    "coprocess_options": {
        "enable_coprocess": true,
        "coprocess_grpc_server": "tcp://127.0.0.1:5000"
	},
	"enable_bundle_downloader": true,
	"bundle_base_url": "http://localhost/bundles/",


Environment variables of the gRPC server:

* TYKLDAPHOST: Hostname or IP address of the LDAP server
* TYKLDAPPORT: TCP port of the LDAP server
* TYKLDAPBINDDN, TYKLDABINDPW: user/password to connect to the LDAP server

Example
-------

Let's imagine you want to add meta data to an HTTP query, in the headers. This meta data is fetched from the LDAP server with the following schema:

    objectidentifier cisco 1.3.6.1.4.1.9
    objectidentifier iotgwSchema cisco:42
    objectidentifier iotgwAttrs iotgwSchema:3
    objectidentifier iotgwOCs iotgwSchema:4
    
    attributeType ( iotgwAttrs:1
        NAME 'deviceid'
        DESC 'Device unique identifier'
        SUP name )
    
    attributeType ( iotgwAttrs:2
        NAME 'meta'
        DESC 'Opaque metadata string'
        SYNTAX 1.3.6.1.4.1.1466.115.121.1.40 )
    
    objectClass ( iotgwOCs:1
        NAME 'thing'
        DESC 'Describe an IOT thing'
        SUP ( top ) AUXILIARY
        MUST ( deviceid $ meta ) )

With this example object:

    dn: cn=athing,ou=things,dc=cisco,dc=com
    objectClass: top
    objectClass: device
    objectClass: thing
    ou: things
    cn: athing
    deviceid: athing
    meta: This is a meta data for the 'athing' thing

You then can launch the gRPC server which will connect to the LDAP server

    TYKLDAPHOST=127.0.0.1 TYKLDAPPORT=389 TYKLDAPBINDDN="cn=admin,dc=cisco,dc=com" TYKLDAPBINDPW=rootpw ./tyk-plugin-ldap-go

    curl -H 'X-Ldap-Param: athing' -H 'Content-Type: application/json' -d '{"h":"v"}' -s -v http://localhost:8080/post|jq .
