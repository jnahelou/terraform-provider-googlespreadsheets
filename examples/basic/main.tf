variable "sheet_id" {
  type        = "string"
  description = "Sheet ID"
}

data "googlespreadsheets_find_empty_row" "empty" {
  spreadsheet_id = "${var.sheet_id}"
  range          = "'Feuille 1'!A1:A20"
}

resource "googlespreadsheets_rows" "foo" {
  spreadsheet_id = "${var.sheet_id}"
  range          = "'Feuille 1'!A${data.googlespreadsheets_find_empty_row.empty.position}:B${data.googlespreadsheets_find_empty_row.empty.position + 1}"

  rows {
    values = ["fooA", "fooB"]
  }

  //First empty change every time, ignore it
  lifecycle {
    ignore_changes = ["range"]
  }
}
