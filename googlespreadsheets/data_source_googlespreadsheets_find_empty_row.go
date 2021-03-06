package googlespreadsheets

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func dataSourceGooglespreadsheetsFindEmptyRow() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGooglespreadsheetsFindEmptyRowRead,

		Schema: map[string]*schema.Schema{
			"spreadsheet_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"range": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"position": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGooglespreadsheetsFindEmptyRowRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	srv := config.ClientSheets
	vrange := d.Get("range").(string)

	mysheet, err := resourceRowsBuild(d, meta)
	if err != nil {
		return errors.Wrap(err, "failed to get sheet")
	}

	resp, err := srv.Spreadsheets.Values.Get(mysheet.SpreadsheetId, vrange).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve data from sheet. %v", err)
	}

	i := 0
	for _, row := range resp.Values {
		i++
		if len(row) == 0 {
			break
		}
	}
	//No empty lines found in response, add new line
	//TODO check if the new row is still included on range
	if i == len(resp.Values) {
		i++
	}

	d.SetId(mysheet.SpreadsheetId + "/" + vrange)
	d.Set("spreadsheet_id", mysheet.SpreadsheetId)
	d.Set("position", strconv.Itoa(i))
	return nil
}
