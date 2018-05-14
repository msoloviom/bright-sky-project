package main

import (
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"./ca"
	"github.com/huandu/go-sqlbuilder"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

type Supplier struct {
	ID   int
	Name sql.NullString
	Cert sql.NullString
}

type Investor struct {
	ID   int
	Name sql.NullString
	Cert sql.NullString
}

type Contract struct {
	ID       int
	Supplier *Supplier
	Investor *Investor
	Stage    sql.NullInt64
	Created  sql.NullInt64

	Title       sql.NullString
	Description sql.NullString
	Amount      sql.NullInt64
	MustBeDone  sql.NullInt64

	SupplierSignature []byte
	InvestorSignature []byte
}

type Offer struct {
	id                int
	contract          *Contract
	supplier          *Supplier
	supplierSignature []byte
	comment           string
}

type saveable interface {
	Save() error
}

func (c *Contract) Save() error {
	return nil
}

func (o *Offer) Save() error {
	return nil
}

type ECDSASignature struct {
	r big.Int
	s big.Int
}

const DbDriver = "sqlite3"
const DbName = "./contracts.sqlite3"

func ListContracts(c echo.Context) error {
	supplierID := c.QueryParam("supplier_id")
	investorID := c.QueryParam("investor_id")
	title := c.QueryParam("title")

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*")
	sb.From("contracts")
	if supplierID != "" {
		a, err := strconv.Atoi(supplierID)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		sb.Where(sb.Equal("supplierID", a))
	}
	if investorID != "" {
		a, err := strconv.Atoi(investorID)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		sb.Where(sb.Equal("investorID", a))
	}
	if title != "" {
		sb.Where(sb.Like("title", title))
	}
	q, args := sb.Build()

	db, err := sql.Open(DbDriver, DbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(q, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var contracts []Contract

	for rows.Next() {
		supplier := Supplier{}
		investor := Investor{}
		contract := Contract{Investor: &investor, Supplier: &supplier}

		err = rows.Scan(&contract.ID, &contract.Supplier.ID, &contract.Investor.ID, &contract.Stage,
			&contract.Created, &contract.Title, &contract.Description, &contract.Amount, &contract.MustBeDone,
			&contract.SupplierSignature, &contract.InvestorSignature)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(contract.Title)
		contracts = append(contracts, contract)
	}

	fmt.Println(contracts)
	return c.JSON(http.StatusOK, contracts)
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

func GetOffer(c echo.Context) error {
	return nil
}

func ListOffers(c echo.Context) error {
	return nil
}

func CreateOffer(c echo.Context) error {
	return nil
}

func DeleteOffer(c echo.Context) error {
	return nil
}

/*
	GET contracts/ - list of contracts, filterable params - supplierID, investorID, title
	GET contracts/{id} - retrieve contract with specific id
	POST contracts/ - create contract
	PATCH contracts/{id} - update contract(accept offer)
	DELETE contracts/{id} - delete contract with specific id

	GET offers/ - list of offers, fliterable params - supplierID, contract_id
	GET offers/{id} - retrieve offer with specific id
	POST offers/ - create offer
	DELETE offers/{id} - delete offer with specific id
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
