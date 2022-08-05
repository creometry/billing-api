package utils

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
    m "billing-api/models"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"fmt"
)

//pdf requestpdf struct
type RequestPdf struct {
	body string
}

//new request to pdf function
func NewRequestPdf(body string) *RequestPdf {
	return &RequestPdf{
		body: body,
	}
}

//parsing template function
func (r *RequestPdf) ParseTemplate(templateFileName string, data m.BillingAccount) error {

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

//generate pdf function
func (r *RequestPdf) GeneratePDF(pdfPath string) (bool, error) {
	t := time.Now().Unix()
	// write whole the body
	err1 := ioutil.WriteFile("storage/"+strconv.FormatInt(int64(t), 10)+".html", []byte(r.body), 0644)
	if err1 != nil {
		panic(err1)
	}

	f, err := os.Open("storage/" + strconv.FormatInt(int64(t), 10) + ".html")
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		os.Remove("storage/" + strconv.FormatInt(int64(t), 10) + ".html")
		log.Fatal(err)
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(f))

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		log.Fatal(err)
	}
	os.Remove("storage/" + strconv.FormatInt(int64(t), 10) + ".html")

	return true, nil
}

func generatePDF(outputPath string,billingAccount m.BillingAccount) {

	r := NewRequestPdf("")

	//html template path
	templatePath := "templates/sample.html"

	//html template data
	//templateData := m.BillingAccount{
	//	BillingAdmins: []m.AdminDetails{{Email: "exmaleadmin@email.com", Phone_number: "21452012", Name: "mohsen"},{Email: "exmaleadmin@email.com", Phone_number: "21452012", Name: "mohsen"}},
	//	Company:     m.Company  {IsCompany : true, TaxId: "222-555", Name: "mohsenlacharikalahou"},
	//	Projects: []m.Project{{
	//		ProjectId:          "85c0194d-488c-48d0-8b01-0917ea578def",
	//		ClusterId:          "7f432018-3f06-4c43-8dda-1578e5c61f29",
	//		State:              "active",
	//		BillingAccountRefer: "b8369c1b-00e9-417b-b145-0a5f76d40550"}}}

	if err := r.ParseTemplate(templatePath, billingAccount); err == nil {
		ok, _ := r.GeneratePDF(outputPath)
		fmt.Println(ok, "pdf generated successfully")
	} else {
		fmt.Println(err)
	}
}