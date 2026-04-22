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
	Destination             Destination       `json:"Destination"`
	ReplacementTemplateData EmailTemplateData `json:"ReplacementTemplateData"`
	Attachments             []EmailAttachment `json:"Attachments"`
}

type EnvioEmailPayload struct {
	Source              string             `json:"Source"`
	Template            string             `json:"Template"`
	Destinations        []EmailDestination `json:"Destinations"`
	DefaultTemplateData EmailTemplateData  `json:"DefaultTemplateData"`
}

type DatosEmail struct {
	Docente            string `json:"docente,omitempty"`
	Estudiante         string `json:"estudiante,omitempty"`
	ConsejoCurricular  string `json:"nombreConsejoCurricular,omitempty"`
	EspacioAcademico   string `json:"espacioAcademico,omitempty"`
	Grupo              string `json:"grupo,omitempty"`
	ProyectoCurricular string `json:"proyectoCurricular,omitempty"`
}

type DatosEmailRequest struct {
	NombreEvaluacion string     `json:"nombreEvaluacion"`
	Periodo          int        `json:"periodo"`
	Correo           string     `json:"correo"`
	Documento        string     `json:"documento"`
	Fecha            string     `json:"fecha"`
	Hora             string     `json:"hora"`
	DatosEmail       DatosEmail `json:"datosEvaluacion"`
}
