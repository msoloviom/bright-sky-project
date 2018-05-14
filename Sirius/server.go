package main

import (
	"crypto/ecdsa"
	"crypto/sha512"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"

	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"./ca"
	"github.com/huandu/go-sqlbuilder"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

type Supplier struct {
	ID   sql.NullInt64
	Name sql.NullString
	Cert sql.NullString
}

type Investor struct {
	ID   int64
	Name sql.NullString
	Cert sql.NullString
}

type ContractBody struct {
	Title       string
	Description string
	Amount      int64
	MustBeDone  string
}

type Contract struct {
	ID       int64
	Supplier *Supplier
	Investor *Investor
	Stage    int64
	Created  string

	ContractBody ContractBody

	SupplierSignature []byte
	InvestorSignature []byte
}

func (c *Contract) GetEncoded() []byte {
	r, _ := json.Marshal(c.ContractBody)
	return r
}

type Timestamp time.Time

type ContractQuery struct {
	InvestorID  int64
	Title       string
	Description string
	Amount      int64
	MustBeDone  Timestamp
}

type Signature []byte

type OfferAcceptionQuery struct {
	OfferID           int64
	ContractID        int64
	InvestorSignature Signature
}

func (s *Signature) UnmarshalParam(src string) error {
	b, err := base64.StdEncoding.DecodeString(src)
	*s = Signature(b)
	return err
}

func (t *Timestamp) UnmarshalParam(src string) error {
	ts, err := time.Parse(time.RFC3339, src)
	*t = Timestamp(ts)
	return err
}

type Offer struct {
	ID                int
	Created           string
	Contract          *Contract
	Supplier          *Supplier
	SupplierSignature []byte
	Comment           string
}

type ECDSASignature struct {
	r big.Int
	s big.Int
}

const dbDriver = "sqlite3"
const dbName = "./contracts.sqlite3"

// ListContracts - api controller for getting list of available contracts
func ListContracts(c echo.Context) error {
	supplierID := c.QueryParam("SupplierID")
	investorID := c.QueryParam("InvestorID")
	title := c.QueryParam("Title")

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*")
	sb.From("contracts")
	if supplierID != "" {
		a, err := strconv.Atoi(supplierID)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		sb.Where(sb.Equal("supplier_id", a))
	}
	if investorID != "" {
		a, err := strconv.Atoi(investorID)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		sb.Where(sb.Equal("investor_id", a))
	}
	if title != "" {
		sb.Where(sb.Like("title", title))
	}
	q, args := sb.Build()

	db, err := sql.Open(dbDriver, dbName)
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
			&contract.Created, &contract.ContractBody.Title, &contract.ContractBody.Description, &contract.ContractBody.Amount, &contract.ContractBody.MustBeDone,
			&contract.SupplierSignature, &contract.InvestorSignature)
		if err != nil {
			log.Fatal(err)
		}
		contracts = append(contracts, contract)
	}

	fmt.Printf("%s", contracts[0].GetEncoded())
	return c.JSON(http.StatusOK, contracts)
}

// GetContract - api controller for retrieving contract by ID
func GetContract(c echo.Context) error {
	ID := c.Param("id")

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*")
	sb.From("contracts")
	if ID != "" {
		a, err := strconv.Atoi(ID)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		sb.Where(sb.Equal("id", a))
	}
	q, args := sb.Build()

	db, err := sql.Open(dbDriver, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	supplier := Supplier{}
	investor := Investor{}
	contract := Contract{Investor: &investor, Supplier: &supplier}

	err = db.QueryRow(q, args...).Scan(&contract.ID, &contract.Supplier.ID, &contract.Investor.ID,
		&contract.Stage, &contract.Created, &contract.ContractBody.Title, &contract.ContractBody.Description,
		&contract.ContractBody.Amount, &contract.ContractBody.MustBeDone, &contract.SupplierSignature, &contract.InvestorSignature)

	if err == sql.ErrNoRows {
		return c.String(http.StatusNotFound, "Contract not found")
	} else if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, contract)
}

// CreateContract - api controller for creating new contract
func CreateContract(c echo.Context) error {
	contractQuery := new(ContractQuery)
	err := c.Bind(contractQuery)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto("contracts")
	ib.Cols("investor_id", "title", "created", "title", "description", "amount", "must_be_done")
	ib.Values(contractQuery.InvestorID, time.Now().Format(time.RFC3339), contractQuery.Title, contractQuery.Description,
		contractQuery.Amount, time.Time(contractQuery.MustBeDone).Format(time.RFC3339))
	q, args := ib.Build()

	db, err := sql.Open(dbDriver, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	res, err := db.Exec(q, args...)
	if err != nil {
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusCreated, struct{ id int64 }{id: id})
}

// UpdateContract - api controller for accepting an offer and finalizing contract creation
func UpdateContract(c echo.Context) error {
	offerAcceptionQuery := new(OfferAcceptionQuery)
	err := c.Bind(offerAcceptionQuery)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	db, err := sql.Open(dbDriver, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*")
	sb.From("contracts")
	sb.Where(sb.Equal("id", offerAcceptionQuery.ContractID))
	q, args := sb.Build()

	supplier := Supplier{}
	investor := Investor{}
	contract := Contract{Investor: &investor, Supplier: &supplier}

	err = db.QueryRow(q, args...).Scan(&contract.ID, &contract.Supplier.ID, &contract.Investor.ID,
		&contract.Stage, &contract.Created, &contract.ContractBody.Title, &contract.ContractBody.Description,
		&contract.ContractBody.Amount, &contract.ContractBody.MustBeDone, &contract.SupplierSignature, &contract.InvestorSignature)

	if err == sql.ErrNoRows {
		return c.String(http.StatusNotFound, "Contract not found")
	} else if err != nil {
		log.Fatal(err)
	}

	sb.Select("supplier_signature")
	sb.From("offers")
	sb.Where(sb.Equal("id", offerAcceptionQuery.OfferID), sb.Equal("contract_id", contract.ID))
	q, args = sb.Build()
	var supplierSignature []byte
	err = db.QueryRow(q, args...).Scan(&supplierSignature)
	if err == sql.ErrNoRows {
		return c.String(http.StatusNotFound, "Offer not found")
	} else if err != nil {
		log.Fatal(err)
	}

	hash := sha512.Sum512(contract.GetEncoded())
	certDer, err := base64.StdEncoding.DecodeString(contract.Investor.Cert.String)
	certObj, err := x509.ParseCertificate(certDer)
	pubKey := certObj.PublicKey.(ecdsa.PublicKey)
	ecdsa.Verify(pubKey, hash)
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
