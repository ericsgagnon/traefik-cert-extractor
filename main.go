package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Acme is the struct for the traefik json certificate file
type Acme struct {
	Default Default `json:"default"`
}

// Default acme provider certs
type Default struct {
	Certificates []Certificate `json:"Certificates"`
}

// Domain has domains
type Domain struct {
	Main string
}

// Certificate holds tls certificate from acme.json
type Certificate struct {
	Domain      Domain `json:"domain"`
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
}

func main() {

	// open acme.json
	jsonFile, err := os.Open("acme.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	// extract all of acme.json to byte slices
	rrsBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	// unmarshall json byte slices to structs
	var acme Acme
	err = json.Unmarshal(rrsBytes, &acme)
	if err != nil {
		log.Fatal(err)
	}

	// extract domains, certs, keys, decode, and write to file
	certs := acme.Default.Certificates
	for i := range certs {

		// extract domain name
		domainName := certs[i].Domain.Main

		// make file name stub from domain name
		fileNameStub := strings.ReplaceAll(domainName, ".", "-")

		// extract key
		key, err := base64.StdEncoding.DecodeString(certs[i].Key)
		if err != nil {
			log.Fatalln(err)
		}

		// write key to file
		keyFileName := fmt.Sprint(fileNameStub, ".key")
		keyFile, err := os.Create(keyFileName)
		if err != nil {
			log.Fatalln(err)
		}
		defer keyFile.Close()
		_, err = io.WriteString(keyFile, string(key))

		//cert := certs[i].Certificate
		cert, err := base64.StdEncoding.DecodeString(certs[i].Certificate)
		certFileName := fmt.Sprint(fileNameStub, ".cert")
		certFile, err := os.Create(certFileName)
		if err != nil {
			log.Fatalln(err)
		}
		defer keyFile.Close()
		_, err = io.WriteString(certFile, string(cert))
	}
}
