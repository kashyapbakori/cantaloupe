package services

import(
	"encoding/json"
	"time"
)

type AllInvoiceDetails struct {
	Number string `json:"number"`
	Description string `json:"description"`
	Date string `json:"date"`
	InvoicerDetails struct {
		Name string `json:"name"`
		Address string `json:"address"`
		Contact string `json:"contact"`
	} `json:"invoicer_details"`
	ComponentDetails struct {
		Name []string `json:"name"`
		Description []string `json:"description"`
		WarrantyTill []string `json:"warranty_till"`
		SerialNo []string `json:"serial_no"`
		Category []int `json:"category"`
	} `json:"component_details"`

}

type InvoiceDetails struct {
	InvoiceNumber string
	InvoicerName string
	Invoicer_add string
	Invoicer_contact string
	Description string
	InvoiceDate string
	// InvoiceDate time.Time
}

type ComponentsInfo struct {
	Name string
	SerialNo string
	Description string
	WarrantyTill string
	InvoiceId int64
	CategoryId int
	Active bool
}

func AddInvoice(details string) {
	var m AllInvoiceDetails
	err := json.Unmarshal([]byte(details), &m)//converting JSON Object to GO structure ...
	CheckErr(err)
	sess := SetupDB()
	i := InvoiceDetails{}
	i.InvoiceNumber = m.Number
	i.InvoicerName = m.InvoicerDetails.Name
	i.Invoicer_add = m.InvoicerDetails.Address
	i.Invoicer_contact = m.InvoicerDetails.Contact
	i.Description = m.Description
	i.InvoiceDate = m.Date
	_, err1 := sess.InsertInto("invoices").
		Columns("invoice_number", "invoicer_name", "invoicer_add", "invoicer_contact", "description", "invoice_date").
		Record(i).
		Exec()
	CheckErr(err1)

	recentInsertedId, err2 := sess.Select("MAX(id)").
		From("invoices").
		ReturnInt64()
	CheckErr(err2)

	c := ComponentsInfo{}
	c.InvoiceId = recentInsertedId
	c.Active = false
	for i := 0; i < len(m.ComponentDetails.Name); i++ {
		c.Name = m.ComponentDetails.Name[i]
		c.SerialNo = m.ComponentDetails.SerialNo[i]
		c.CategoryId = m.ComponentDetails.Category[i]
		c.Description = m.ComponentDetails.Description[i]
		c.WarrantyTill = m.ComponentDetails.WarrantyTill[i]
		_, err3 := sess.InsertInto("components").
			Columns("invoice_id", "serial_no", "name", "category_id", "warranty_till", "description", "active").
			Record(c).
			Exec()
		CheckErr(err3)
	}
}

type DisplayInvoice struct{
	Id int
	Invoice_number string
	Invoicer_name string
	Invoice_timestamp time.Time
	Invoice_date string

	Components []DisplayAllComponents
}

func DisplayInvoices() []byte {
	sess := SetupDB()
	invoicesDetails := []DisplayInvoice{}
	sess.Select("id, invoice_number, invoicer_name, invoice_date as Invoice_timestamp").
		From("invoices").
		LoadStruct(&invoicesDetails)
	//extract only date from timestamp========
	for i := 0; i < len(invoicesDetails); i++ {
		t := invoicesDetails[i].Invoice_timestamp
		invoicesDetails[i].Invoice_date = t.Format("2006-01-02")
	}
	//================================
	b, err := json.Marshal(invoicesDetails)
	CheckErr(err)

	return b
}

func DisplayOneInvoice(invoiceId int) []byte {
	sess := SetupDB()
	//===============Perticuller Invoice Details =================================
	invoiceDetails := DisplayInvoice{}
	sess.Select("id, invoice_number, invoicer_name, invoice_date as Invoice_timestamp").
		From("invoices").
		Where("id = ?", invoiceId ).
		LoadStruct(&invoiceDetails)
		t := invoiceDetails.Invoice_timestamp
		invoiceDetails.Invoice_date = t.Format("2006-01-02")
	//============================================================================

	//==============components details related to perticuller invoice ============

	components := []DisplayAllComponents{}
	sess.Select("c.id, c.invoice_id, c.serial_no, c.name, c.description, c.warranty_till as Warranty_timestamp, machines.name as Machine").
		From("components c").
    LeftJoin("machine_components", "c.id = machine_components.component_id").
    LeftJoin("machines", "machines.id = machine_components.machine_id").
		Where("c.invoice_id = ?", invoiceId).
		LoadStruct(&components)


	//extract only date from timestamp========
	for i := 0; i < len(components); i++ {
		t := components[i].Warranty_timestamp
		components[i].Warranty_till = t.Format("2006-01-02")
	}
	//=========================================

	invoiceDetails.Components = components

	//============================================================================

	b, err := json.Marshal(invoiceDetails)
	CheckErr(err)

	return b
}
