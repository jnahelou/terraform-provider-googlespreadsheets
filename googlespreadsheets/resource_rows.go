package googlespreadsheets

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	sheets "google.golang.org/api/sheets/v4"
)

func resourceRows() *schema.Resource {
	return &schema.Resource{
		Create: resourceRowsCreate,
		Read:   resourceRowsRead,
		Update: resourceRowsUpdate,
		Delete: resourceRowsDelete,

		Schema: map[string]*schema.Schema{
			"spreadsheet_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"range": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"rows": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"values": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required: true,
						},
					},
				},
			},
		},
	}
}

// resourceRowsCreate creates a new row via the API.
func resourceRowsCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	srv := config.ClientSheets

	sheet, err := resourceRowsBuild(d, meta)
	if err != nil {
		return errors.Wrap(err, "failed to build row")
	}

	vrange := d.Get("range").(string)
	var values [][]interface{}
	var one_row []interface{}

	rows := d.Get("rows").([]interface{})
	log.Printf("[DEBUG] row %v\n", rows)
	for i, r := range rows {
		values = append(values, one_row)

		rc := r.(map[string]interface{})
		rvs := rc["values"].([]interface{})
		for _, rv := range rvs {
			log.Printf("[DEBUG] Row values: %v\n", rv)
			values[i] = append(values[i], rv)
		}
	}

	rb := &sheets.ValueRange{
		Values: values,
	}

	valueInputOption := "USER_ENTERED"
	ctx := context.Background()
	_, err = srv.Spreadsheets.Values.Update(sheet.SpreadsheetId, vrange, rb).ValueInputOption(valueInputOption).Context(ctx).Do()
	if err != nil {
		panic(fmt.Errorf("Error update. %v", err))
	}

	return resourceRowsRead(d, meta)
}

// resourceRowsRead reads information about the row from the API.
func resourceRowsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	srv := config.ClientSheets
	vrange := d.Get("range").(string)

	mysheet, err := resourceRowsBuild(d, meta)
	if err != nil {
		return errors.Wrap(err, "failed to get sheet")
	}

	resp, err := srv.Spreadsheets.Values.Get(mysheet.SpreadsheetId, vrange).Do()
	if err != nil {
		panic(fmt.Errorf("unable to retrieve data from sheet. %v", err))
	}

	var result []interface{}
	if len(resp.Values) > 0 {
		for _, row := range resp.Values {
			m := make(map[string]interface{})

			s := make([]string, len(row))
			for i, v := range row {
				s[i] = v.(string)
			}
			m["values"] = s
			result = append(result, m)
		}
	}
	log.Printf("[DEBUG]Rows : %v\n", result)
	d.SetId(mysheet.SpreadsheetId + "/" + vrange)
	d.Set("spreadsheet_id", mysheet.SpreadsheetId)
	d.Set("rows", result)

	return nil
}

// resourceRowsUpdate updates an row via the API.
func resourceRowsUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceRowsCreate(d, meta)
}

// resourceRowsDelete deletes an row via the API.
func resourceRowsDelete(d *schema.ResourceData, meta interface{}) error {
	//TODO
	//Code here

	d.SetId("")

	return nil
}

func resourceRowsBuild(d *schema.ResourceData, meta interface{}) (*sheets.Spreadsheet, error) {
	config := meta.(*Config)

	srv := config.ClientSheets
	spreadsheet_id := d.Get("spreadsheet_id").(string)

	mysheet, err := srv.Spreadsheets.Get(spreadsheet_id).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	return mysheet, nil
}
