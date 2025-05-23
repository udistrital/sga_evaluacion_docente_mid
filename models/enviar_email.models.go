package models

type EmailAttachment struct {
	ContentType string `json:"ContentType"`
	FileName    string `json:"FileName"`
	Base64File  string `json:"Base64File"`
}

type EmailTemplateData struct {
	NombreUsuario     string `json:"nombre_usuario"`
	DocumentoUsuario  string `json:"documento_usuario"`
	NombreEvaluacion  string `json:"nombre_evaluacion"`
	NumeroPeriodo     string `json:"numero_periodo"`
	FechaEvaRealizada string `json:"fecha_eva_realizada"`
	HoraEvaRealizada  string `json:"hora_eva_realizada"`
}

type Destination struct {
	ToAddresses []string `json:"ToAddresses"`
}

type EmailDestination struct {
	Destination             Destination        `json:"Destination"`
	ReplacementTemplateData EmailTemplateData  `json:"ReplacementTemplateData"`
	Attachments             []EmailAttachment  `json:"Attachments"`
}

type EnvioEmailPayload struct {
	Source              string            `json:"Source"`
	Template            string            `json:"Template"`
	Destinations        []EmailDestination `json:"Destinations"`
	DefaultTemplateData EmailTemplateData  `json:"DefaultTemplateData"`
}
