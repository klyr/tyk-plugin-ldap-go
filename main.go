package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	coprocess "github.com/TykTechnologies/tyk-protobuf/bindings/go"
)

func getEnv(varname string) string {
	value, ok := os.LookupEnv(varname)
	if !ok {
		log.Fatalf("Environment variable '%s' must be set\n", varname)
	}
	return value
}

func main() {
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("Failed to listen: %w", err)
	}

	ldapHost := getEnv("TYKLDAPHOST")
	ldapPort := getEnv("TYKLDAPPORT")
	ldapBindDn := getEnv("TYKLDAPBINDDN")
	ldapBindPw := getEnv("TYKLDAPBINDPW")

	ldapBind(ldapHost, ldapPort, ldapBindDn, ldapBindPw)

	log.Println("Listening ...")
	s := grpc.NewServer()
	coprocess.RegisterDispatcherServer(s, &Dispatcher{})
	s.Serve(lis)
}
