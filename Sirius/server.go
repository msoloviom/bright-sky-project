package main

import (
	"crypto/ecdsa"
	"crypto/sha512"
	"crypto/x509"
	"database/sql"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

type UserAbstract struct {
	ID   sql.NullInt64
	Name sql.NullString
	Cert sql.NullString
}

type userExternal struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
	Cert string `json:"cert"`
}

type loadable interface {
	Load(cache map[int64]UserAbstract) error
}

type Supplier struct {
	loadable
	UserAbstract
}

type Investor struct {
	loadable
	UserAbstract
}

const GatewayURL = "http://172.20.10.6:8000"

// Load Supplier from the external api
func (s *Supplier) Load(cache map[int64]UserAbstract) (err error) {
	suppliersAPI := GatewayURL + "/clients/%d"
	if !s.ID.Valid {
		return nil
	}
	v, ok := cache[s.ID.Int64]
	if ok {
		s.UserAbstract = v

	} else {
		res, err := http.Get(fmt.Sprintf(suppliersAPI, s.ID.Int64))

		if err != nil {
			log.Print(err)
			return err
		}
		defer res.Body.Close()
		u := userExternal{}
		err = json.NewDecoder(res.Body).Decode(&u)
		//fmt.Print(u)
		if err != nil {
			log.Print(err)
			return err
		}
		s.Name.Valid = true
		s.Name.String = u.Name
		s.Cert.Valid = true
		s.Cert.String = u.Cert

		cache[s.ID.Int64] = s.UserAbstract
	}
	return nil
}

// Load Investor from the external api
func (s *Investor) Load(cache map[int64]UserAbstract) error {
	investorsAPI := GatewayURL + "/investors/%d"
	if !s.ID.Valid {
		return nil
	}
	v, ok := cache[s.ID.Int64]
	if ok {
		s.UserAbstract = v

	} else {
		res, err := http.Get(fmt.Sprintf(investorsAPI, s.ID.Int64))

		if err != nil {
			log.Print(err)
			return err
		}
		defer res.Body.Close()
		u := userExternal{}
		err = json.NewDecoder(res.Body).Decode(&u)
		//fmt.Print(u)
		if err != nil {
			log.Print(err)
			return err
		}
		s.Name.Valid = true
		s.Name.String = u.Name
		s.Cert.Valid = true
		s.Cert.String = u.Cert

		cache[s.ID.Int64] = s.UserAbstract
	}
	return nil
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

	SupplierSignature sql.NullString
	InvestorSignature sql.NullString
}

func (c *Contract) GetEncoded() []byte {
	r, _ := json.Marshal(c.ContractBody)
	return r
}

type Timestamp time.Time

func (t *Timestamp) UnmarshalParam(src string) error {
	ts, err := time.Parse(time.RFC3339, src)
	*t = Timestamp(ts)
	return err
}

type ContractQuery struct {
	Title       string
	Description string
	Amount      int64
	MustBeDone  *Timestamp
}

type Signature []byte

type OfferAcceptionQuery struct {
	OfferID           int64
	ContractID        int64
	InvestorSignature string
}

type Offer struct {
	ID                int64
	Created           string
	ContractID        int64
	Supplier          *Supplier
	SupplierSignature sql.NullString
	Comment           sql.NullString
}

type OfferQuery struct {
	ContractID        int64
	Comment           string
	SupplierSignature string
}

type ECDSASignature struct {
	r big.Int
	s big.Int
}

type InvestorContext struct {
	echo.Context
	InvestorID sql.NullInt64
}

type SupplierContext struct {
	echo.Context
	SupplierID sql.NullInt64
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
	investorsCache := make(map[int64]UserAbstract)
	suppliersCache := make(map[int64]UserAbstract)

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
		contract.Investor.Load(investorsCache)
		contract.Supplier.Load(suppliersCache)
		contracts = append(contracts, contract)
	}

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
	investorsCache := make(map[int64]UserAbstract)
	suppliersCache := make(map[int64]UserAbstract)

	contract.Investor.Load(investorsCache)
	contract.Supplier.Load(suppliersCache)

	return c.JSON(http.StatusOK, contract)
}

// CreateContract - api controller for creating new contract
func CreateContract(c echo.Context) error {
	ic := c.(InvestorContext)

	contractQuery := new(ContractQuery)
	err := c.Bind(contractQuery)
	if err != nil {
		log.Print(err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto("contracts")
	ib.Cols("investor_id", "title", "created", "description", "amount", "must_be_done")
	ib.Values(ic.InvestorID, contractQuery.Title, time.Now().Format(time.RFC3339), contractQuery.Description,
		contractQuery.Amount, time.Time(*contractQuery.MustBeDone).Format(time.RFC3339))
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

// VerifySignature verifies ecdsa signature, algorithm - ECDSA with curve P-384 and hash - SHA-512-384
func VerifySignature(b64signature, pemcert string, data []byte) bool {
	derSignature, err := base64.StdEncoding.DecodeString(b64signature)
	if err != nil {
		return false
	}
	sig := ECDSASignature{}
	_, err = asn1.Unmarshal(derSignature, &sig)
	if err != nil {
		return false
	}
	hash := sha512.Sum384(data)
	certBlock, rest := pem.Decode([]byte(pemcert))
	if len(rest) > 0 {
		return false
	}
	certObj, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return false
	}
	pubKey := certObj.PublicKey.(ecdsa.PublicKey)
	return ecdsa.Verify(&pubKey, hash[:], &sig.r, &sig.s)
}

// UpdateContract - api controller for accepting an offer and finalizing contract creation
func UpdateContract(c echo.Context) error {
	ic := c.(InvestorContext)

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
	sb.Where(sb.Equal("id", offerAcceptionQuery.ContractID), sb.Equal("investor_id", ic.InvestorID))
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
	var supplierSignature string
	err = db.QueryRow(q, args...).Scan(&supplierSignature)

	if err == sql.ErrNoRows {
		return c.String(http.StatusNotFound, "Offer not found")
	} else if err != nil {
		log.Fatal(err)
	}

	contractToBeSigned := contract.GetEncoded()
	//supplierVerified := VerifySignature(supplierSignature, contract.Supplier.Cert.String, contractToBeSigned)
	investorVerified := VerifySignature(offerAcceptionQuery.InvestorSignature, contract.Investor.Cert.String, contractToBeSigned)

	if investorVerified {
		ub := sqlbuilder.NewUpdateBuilder()
		ub.Update("contracts")
		ub.Where(ub.Equal("id", contract.ID))
		ub.Set(ub.Assign("supplier_signature", supplierSignature), ub.Assign("investor_signature", offerAcceptionQuery.InvestorSignature), ub.Assign("stage", 1))
		q, args := ub.Build()
		res, err := db.Exec(q, args...)
		if err != nil {
			log.Fatal(err)
		}
		affected, err := res.RowsAffected()
		if err != nil || affected != 1 {
			log.Fatal(err, affected)
		}
		return c.String(http.StatusOK, "")

	} else {
		return c.String(http.StatusBadRequest, "Signature not verified")
	}
}

// DeleteContract - api controller for removing contract
func DeleteContract(c echo.Context) error {
	ic := c.(InvestorContext)

	dB, err := sql.Open(dbDriver, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer dB.Close()

	id := c.Param("id")
	db := sqlbuilder.NewDeleteBuilder()
	db.DeleteFrom("contracts")
	db.Where(db.Equal("id", id), db.Equal("investor_id", ic.InvestorID))
	q, args := db.Build()

	res, err := dB.Exec(q, args...)
	if err != nil {
		log.Fatal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if affected != 1 {
		return c.String(http.StatusNotFound, "Contract not found")
	}

	return c.String(http.StatusOK, "")
}

// GetOffer - api controller for retrieving an offer by ID
func GetOffer(c echo.Context) error {
	ID := c.Param("id")

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*")
	sb.From("offers")
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

	offer := Offer{Supplier: &supplier}

	err = db.QueryRow(q, args...).Scan(offer.ID, offer.ContractID, offer.Supplier.ID, offer.SupplierSignature, offer.Comment, offer.Created)
	if err == sql.ErrNoRows {
		return c.String(http.StatusNotFound, "Contract not found")
	} else if err != nil {
		log.Fatal(err)
	}

	suppliersCache := make(map[int64]UserAbstract)
	offer.Supplier.Load(suppliersCache)
	return c.JSON(http.StatusOK, offer)
}

// ListOffers - api controller for obtaining list of offers
func ListOffers(c echo.Context) error {
	supplierID := c.QueryParam("SupplierID")
	contractID := c.QueryParam("ContractID")

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*")
	sb.From("offers")
	if supplierID != "" {
		a, err := strconv.Atoi(supplierID)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		sb.Where(sb.Equal("supplier_id", a))
	}
	if contractID != "" {
		a, err := strconv.Atoi(contractID)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		sb.Where(sb.Equal("contract_id", a))
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

	var offers []Offer
	suppliersCache := make(map[int64]UserAbstract)

	for rows.Next() {
		supplier := Supplier{}
		offer := Offer{Supplier: &supplier}

		err = rows.Scan(&offer.ID, &offer.ContractID, &offer.Supplier.ID, &offer.SupplierSignature, &offer.Comment, &offer.Created)
		if err != nil {
			log.Fatal(err)
		}
		offer.Supplier.Load(suppliersCache)
		offers = append(offers, offer)
	}

	return c.JSON(http.StatusOK, offers)
}

// CreateOffer - api controller for creation of an offer
func CreateOffer(c echo.Context) error {
	sc := c.(SupplierContext)

	offerQuery := new(OfferQuery)

	supplier := Supplier{UserAbstract: UserAbstract{ID: sc.SupplierID}}

	m := make(map[int64]UserAbstract)
	err := supplier.Load(m)

	if err != nil {
		log.Print(err)
		return c.String(http.StatusBadGateway, "Supplier's certificate could not be loaded")
	}

	err = c.Bind(offerQuery)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	db, err := sql.Open(dbDriver, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("title", "description", "amount", "must_be_done")
	sb.From("contracts")
	sb.Where(sb.Equal("id", offerQuery.ContractID))
	q, args := sb.Build()

	contractBody := ContractBody{}

	err = db.QueryRow(q, args...).Scan(&contractBody.Title, &contractBody.Description, &contractBody.Amount, &contractBody.MustBeDone)

	if err == sql.ErrNoRows {
		return c.String(http.StatusNotFound, "Contract not found")
	} else if err != nil {
		log.Fatal(err)
	}

	contractEncoded, err := json.Marshal(contractBody)
	supplierVerified := VerifySignature(offerQuery.SupplierSignature, supplier.Cert.String, contractEncoded)

	if !supplierVerified {
		return c.String(http.StatusBadRequest, "Bad Signature")
	}

	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto("offers")
	ib.Cols("contract_id", "supplier_id", "supplier_signature", "comment", "created")
	ib.Values(offerQuery.ContractID, sc.SupplierID, offerQuery.SupplierSignature, offerQuery.Comment, time.Now().Format(time.RFC3339))
	q, args = ib.Build()

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

// DeleteOffer - api controller for removing offers by ID
func DeleteOffer(c echo.Context) error {
	sc := c.(SupplierContext)

	dB, err := sql.Open(dbDriver, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer dB.Close()

	id := c.Param("id")
	db := sqlbuilder.NewDeleteBuilder()
	db.DeleteFrom("offers")
	db.Where(db.Equal("id", id), db.Equal("supplier_id", sc.SupplierID))
	q, args := db.Build()

	res, err := dB.Exec(q, args...)
	if err != nil {
		log.Fatal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if affected != 1 {
		return c.String(http.StatusNotFound, "Offer not found")
	}

	return c.String(http.StatusOK, "")
}

type (
	SupplierAuthorizationToken string
	InvestorAuthorizationToken string
)

func (t SupplierAuthorizationToken) Authorize() (int64, error) {
	tokenSuppliersApi := GatewayURL + "/clients/token/"
	res, err := http.Get(tokenSuppliersApi + string(t))

	if err != nil {
		log.Print(err)
		return 0, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		idstr, err := ioutil.ReadAll(res.Body)
		a, err := strconv.Atoi(string(idstr))
		return int64(a), err
	} else {
		return 0, errors.New("Unauthorized")
	}
}

func (t InvestorAuthorizationToken) Authorize() (int64, error) {
	tokenInvestorsApi := GatewayURL + "/investors/token/"
	res, err := http.Get(tokenInvestorsApi + string(t))

	if err != nil {
		log.Print(err)
		return 0, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		idstr, err := ioutil.ReadAll(res.Body)
		a, err := strconv.Atoi(string(idstr))
		return int64(a), err
	} else {
		return 0, errors.New("Unauthorized")
	}
}

func ResponseHeaderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Sirius Contracts Service 1.0")
		return next(c)
	}
}

func SupplierAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var token SupplierAuthorizationToken
		var s string
		fmt.Sscanf(c.Request().Header.Get("Authorization"), "Token %s", &s)
		token = SupplierAuthorizationToken(s)
		id, err := token.Authorize()
		if err != nil {
			return echo.ErrUnauthorized
		}
		ID := sql.NullInt64{Int64: id, Valid: true}
		sc := SupplierContext{c, ID}
		return next(sc)
	}
}

func InvestorAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var token SupplierAuthorizationToken
		var s string
		fmt.Sscanf(c.Request().Header.Get("Authorization"), "Token %s", &s)
		token = SupplierAuthorizationToken(s)
		id, err := token.Authorize()
		if err != nil {
			return echo.ErrUnauthorized
		}
		ID := sql.NullInt64{Int64: id, Valid: true}
		sc := InvestorContext{c, ID}
		return next(sc)
	}
}

/*
	GET contracts/ - list of contracts, filterable params - supplierID, investorID, title
	GET contracts/{id} - retrieve contract with specific id
	POST contracts/ - create contract
	PATCH contracts/{id} - update contract(accept offer)
	DELETE contracts/{id} - delete contract with specific id

	GET offers/ - list of offers, fliterable params - supplierID, contractID
	GET offers/{id} - retrieve offer with specific id
	POST offers/ - create offer
	DELETE offers/{id} - delete offer with specific id
*/
func main() {
	e := echo.New()
	e.Use(ResponseHeaderMiddleware)
	e.GET("/contracts", ListContracts)
	e.GET("/contracts/:id", GetContract)
	e.POST("/contracts", CreateContract, InvestorAuthMiddleware)
	e.PATCH("/contracts/:id", UpdateContract, InvestorAuthMiddleware)
	e.DELETE("/contracts/:id", DeleteContract, InvestorAuthMiddleware)

	e.GET("/offers", ListOffers)
	e.GET("/offers/:id", GetOffer)
	e.POST("/offers", CreateOffer, SupplierAuthMiddleware)
	e.DELETE("/offers/:id", DeleteOffer, SupplierAuthMiddleware)

	e.Logger.Fatal(e.Start(":1323"))
}
