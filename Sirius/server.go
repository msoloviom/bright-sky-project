package main

import (
	"crypto/x509"
	"log"
	"math/big"
	"sirius/ca"
	"time"

	"github.com/labstack/echo"
)

type saveable interface {
	Save() error
}

type Supplier struct {
	id   int
	name string
	cert string
}

func (s *Supplier) Save() error {
	return nil
}

type Customer struct {
	id   int
	name string
	cert string
}

func (c *Customer) Save() error {
	return nil
}

type Contract struct {
	id       int
	supplier *Supplier
	customer *Customer
	stage    int
	created  time.Time

	title       string
	description string
	amount      int
	mustBeDone  time.Time

	supplierSignature []byte
	customerSignature []byte
}

func (c *Contract) Save() error {
	return nil
}

type ECDSASignature struct {
	r big.Int
	s big.Int
}

func ListContracts(c echo.Context) error {
	return nil
}

func GetContract(c echo.Context) error {
	return nil
}

func CreateContract(c echo.Context) error {
	return nil
}

func UpdateContract(c echo.Context) error {
	return nil
}

func DeleteContract(c echo.Context) error {
	return nil
}

/*
	GET contracts/ - list of contracts, filterable params - supplier_id, customer_id, title
	GET contracts/{id} - retrieve contract with specific id
	POST contracts/ - create contract
	PATCH contracts/{id} - update contract
	DELETE contracts/{id} - delete contract with specific id
*/
func main() {

	log.Println("Generating an ECDSA P-256 Private Key")
	ECKey := ca.GenerateECKey()
	ca.GenerateCert(&ECKey.PublicKey, ECKey, "ECDSA Sirius Root Authority", x509.KeyUsageCertSign|x509.KeyUsageCRLSign, "sirius.pem")

	e := echo.New()
	e.GET("/contracts", ListContracts)
	e.GET("/contracts/:id", GetContract)
	e.POST("/contracts", CreateContract)
	e.PATCH("/contracts/:id", UpdateContract)
	e.DELETE("/contracts/:id", DeleteContract)

	e.Logger.Fatal(e.Start(":1323"))
}
