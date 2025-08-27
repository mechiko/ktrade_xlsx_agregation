package domain

type Utilisation struct {
	Id               int64  `db:"id,omitempty"`
	IdOrderMarkCodes int64  `db:"id_order_mark_codes"`
	CreateDate       string `db:"create_date"`
	ProductionDate   string `db:"production_date"`
	ExpirationDate   string `db:"expiration_date"`
	UsageType        string `db:"usage_type"`
	Inn              string `db:"inn"`
	Kpp              string `db:"kpp"`
	Version          string `db:"version"`
	State            string `db:"state"`
	Status           string `db:"status"`
	ReportId         string `db:"report_id"`
	Archive          int    `db:"archive"`
	Json             string `db:"json"`
	Quantity         string `db:"quantity"`
	PrimaryDocNumber string `db:"primary_doc_number"`
	PrimaryDocDate   string `db:"primary_doc_date"`
	AlcVolume        string `db:"alc_volume"`
}

type UtilisationCodes struct {
	Id                     int64  `db:"id,omitempty"`
	IdOrderMarkUtilisation int64  `db:"id_order_mark_utilisation"`
	SerialNumber           string `db:"serial_number"`
	Code                   string `db:"code"`
	Status                 string `db:"status"`
}

type Order struct {
	Id                  int64  `db:"id"`
	CreateDate          string `db:"create_date"`
	Gtin                string `db:"gtin"`
	Quantity            int    `db:"quantity"`
	SerialNumberType    string `db:"serial_number_type"`
	TemplateId          int    `db:"template_id"`
	CisType             string `db:"cis_type"`
	ContactPerson       string `db:"contact_person"`
	ReleaseMethodType   string `db:"release_method_type"`
	CreateMethodType    string `db:"create_method_type"`
	PaymentType         string `db:"payment_type"`
	ProductionOrderId   string `db:"production_order_id"`
	ProductName         string `db:"product_name"`
	ProductCapacity     string `db:"product_capacity"`
	ProductShelfLife    string `db:"product_shelf_life"`
	ProductTemplate     string `db:"product_template"`
	Comment             string `db:"comment"`
	Version             string `db:"version"`
	State               string `db:"state"`
	Status              string `db:"status"`
	OrderId             string `db:"order_id"`
	Archive             int    `db:"archive"`
	Json                string `db:"json"`
	ServiceProviderId   string `db:"service_provider_id"`
	ServiceProviderName string `db:"service_provider_name"`
	ServiceProviderRole string `db:"service_provider_role"`
}

type OrderMarkCodesSerialNumbers struct {
	Id               int64  `db:"id"`
	IdOrderMarkCodes int64  `db:"id_order_mark_codes"`
	Gtin             string `db:"gtin"`
	SerialNumber     string `db:"serial_number"`
	Code             string `db:"code"`
	BlockId          string `db:"block_id"`
	Status           string `db:"status"`
}
