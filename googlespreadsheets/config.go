/**
*
* Fork from Google provider in case of merge between GCP and GSuite
*
 */
package googlespreadsheets

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/pathorcontents"
	"github.com/hashicorp/terraform/version"
	sheets "google.golang.org/api/sheets/v4"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// Config is the configuration structure used to instantiate the Google
// provider.
type Config struct {
	Credentials string
	Project     string
	Region      string
	Zone        string

	client    *http.Client
	userAgent string

	tokenSource oauth2.TokenSource

	ClientSheets *sheets.Service
}

func (c *Config) loadAndValidate() error {
	var account accountFile
	clientScopes := []string{
		"https://www.googleapis.com/auth/spreadsheets",
	}

	var client *http.Client
	var tokenSource oauth2.TokenSource

	if c.Credentials != "" {
		contents, _, err := pathorcontents.Read(c.Credentials)
		if err != nil {
			return fmt.Errorf("Error loading credentials: %s", err)
		}

		// Assume account_file is a JSON string
		if err := parseJSON(&account, contents); err != nil {
			return fmt.Errorf("Error parsing credentials '%s': %s", contents, err)
		}

		// Get the token for use in our requests
		log.Printf("[INFO] Requesting Google token...")
		log.Printf("[INFO]   -- Email: %s", account.ClientEmail)
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		log.Printf("[INFO]   -- Private Key Length: %d", len(account.PrivateKey))

		conf := jwt.Config{
			Email:      account.ClientEmail,
			PrivateKey: []byte(account.PrivateKey),
			Scopes:     clientScopes,
			TokenURL:   "https://accounts.google.com/o/oauth2/token",
		}

		// Initiate an http.Client. The following GET request will be
		// authorized and authenticated on the behalf of
		// your service account.
		client = conf.Client(context.Background())

		tokenSource = conf.TokenSource(context.Background())
	} else {
		log.Printf("[INFO] Authenticating using DefaultClient")
		err := error(nil)
		client, err = google.DefaultClient(context.Background(), clientScopes...)
		if err != nil {
			return err
		}

		tokenSource, err = google.DefaultTokenSource(context.Background(), clientScopes...)
		if err != nil {
			return err
		}
	}

	c.tokenSource = tokenSource

	client.Transport = logging.NewTransport("Google", client.Transport)

	projectURL := "https://www.terraform.io"
	userAgent := fmt.Sprintf("Terraform/%s (+%s)",
		version.String(), projectURL)

	c.client = client
	c.userAgent = userAgent

	var err error

	log.Printf("[INFO] Instantiating Spreadsheet client...")
	c.ClientSheets, err = sheets.New(client)
	if err != nil {
		return err
	}
	c.ClientSheets.UserAgent = userAgent

	return nil
}

// accountFile represents the structure of the account file JSON file.
type accountFile struct {
	PrivateKeyId string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
	ClientId     string `json:"client_id"`
}

func parseJSON(result interface{}, contents string) error {
	r := strings.NewReader(contents)
	dec := json.NewDecoder(r)

	return dec.Decode(result)
}
