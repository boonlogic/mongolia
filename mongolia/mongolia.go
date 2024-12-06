package mongolia

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ODM struct {
	URI, DB                 string
	username, password      string
	certificateFile, caFile string
	client                  *mongo.Client
	database                *mongo.Database
	colls                   map[string]Collection
	timeout                 time.Duration
}

func NewODM() *ODM {
	return &ODM{
		URI:             "mongodb://localhost:27017",
		DB:              "mongolia-local",
		username:        "",
		password:        "",
		certificateFile: "",
		caFile:          "",
		timeout:         10 * time.Second,
		colls:           make(map[string]Collection),
	}
}

func (odm *ODM) SetURI(uri string) *ODM {
	odm.URI = uri
	return odm
}

func (odm *ODM) SetUsername(username string) *ODM {
	odm.username = username
	return odm
}

func (odm *ODM) SetPassword(password string) *ODM {
	odm.password = password
	return odm
}

func (odm *ODM) SetCertificateFile(certficateFile string) *ODM {
	odm.certificateFile = certficateFile
	return odm
}

func (odm *ODM) SetCAFile(caFile string) *ODM {
	odm.caFile = caFile
	return odm
}

func (odm *ODM) SetDBName(db string) *ODM {
	odm.DB = db
	return odm
}

func (odm *ODM) SetTimeout(timeout time.Duration) *ODM {
	odm.timeout = timeout
	return odm
}

func (odm *ODM) Connect() *Error {
	ctx, _ := context.WithTimeout(context.Background(), odm.timeout)
	var err error

	//set uri
	clientOptions := options.Client().ApplyURI(odm.URI)

	// Use authentication
	if odm.username != "" && odm.password != "" {
		credentials := options.Credential{
			Username: odm.username,
			Password: odm.password,
		}

		clientOptions.SetAuth(credentials)
	}

	//Add on a certificiate
	if odm.certificateFile != "" {
		// Load client certificate and private key from the same PEM file
		cert, err := tls.LoadX509KeyPair(odm.certificateFile, odm.certificateFile)
		if err != nil {
			errorString := fmt.Sprintf("Error: unable to load certificate file: %v (%v)", odm.certificateFile, err)
			return NewErrorString(500, errorString)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		// Optionally set the CA certificate file
		if odm.caFile != "" {
			caCert, err := os.ReadFile(odm.caFile)
			if err != nil {
				errorString := fmt.Sprintf("Error: unable to load CA file: %v (%v)", odm.caFile, err)
				return NewErrorString(500, errorString)
			}
			caCertPool := x509.NewCertPool()
			if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
				errorString := fmt.Sprintf("Error: CA file must be in PEM format: %v", odm.caFile)
				return NewErrorString(500, errorString)
			}
			tlsConfig.RootCAs = caCertPool
			tlsConfig.InsecureSkipVerify = false
		} else {
			// If we don't have a CA, we'll skip authenticating the server
			tlsConfig.InsecureSkipVerify = true
		}

		// Set the TLS option for mongo options
		clientOptions.SetTLSConfig(tlsConfig)
	}

	// create mongo connection
	odm.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return NewError(500, err)
	}

	// test connection
	err = odm.client.Ping(ctx, nil)
	if err != nil {
		errorString := fmt.Sprintf("Error unable to connect to Mongo: %v", err)
		return NewErrorString(500, errorString)
	}

	odm.database = odm.client.Database(odm.DB)
	return nil
}

func (odm *ODM) GetCollection(name string) *Collection {
	coll, ok := odm.colls[name]
	if !ok {
		return odm.CreateCollection(name, nil)
	}
	return &coll
}

func (odm *ODM) CreateCollection(name string, indexes interface{}) *Collection {
	coll := odm.database.Collection(name)
	c := &Collection{
		name:    name,
		coll:    coll,
		timeout: odm.timeout,
	}
	odm.colls[name] = *c
	if indexes != nil {
		err := c.CreateIndexes(indexes)
		if err != nil {
			log.Printf("Error Creating Indexes: %v\n", err.ToString())
		}
	}
	return c
}

func (odm *ODM) CreateTimeSeriesCollection(name string, opts *options.TimeSeriesOptions, indexes interface{}) *Collection {
	//Specify a timeseries collection
	col_opts := options.CreateCollection().SetTimeSeriesOptions(opts)
	ctx, _ := context.WithTimeout(context.Background(), odm.timeout)
	err := odm.database.CreateCollection(ctx, name, col_opts)
	if err != nil {
		switch e := err.(type) {
		case mongo.CommandError: // raises a specific CommandError if collection already exists
			if e.Name != "NamespaceExists" {
				log.Printf("Error Creating TimeSeries %v\n", err.Error())
				return nil
			}
		default:
			log.Printf("Error Creating TimeSeries %v\n", err.Error())
			return nil
		}
	}

	c := &Collection{
		name:    name,
		coll:    odm.database.Collection(name),
		timeout: odm.timeout,
	}
	odm.colls[name] = *c
	if indexes != nil {
		err := c.CreateIndexes(indexes)
		if err != nil {
			log.Printf("Error Creating Indexes: %v\n", err.ToString())
		}
	}
	return c
}

func (odm *ODM) Disconnect() {
	odm.client.Disconnect(context.Background())
}

// Drop deletes all ODM data.
// It fails if ODM is not ephemeral.
func (odm *ODM) Drop() {
	odm.database.Drop(context.Background())
}
