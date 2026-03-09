package models

import (
	"encoding/json"
	"strings"
	"time"
)

// CustomTime handles multiple date formats
type CustomTime struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler interface
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)

	if str == "null" || str == "" {
		ct.Time = time.Time{}
		return nil
	}

	// List of possible time formats
	formats := []string{
		"2006-01-02T15:04:05Z07:00",              // RFC3339
		"2006-01-02T15:04:05.000Z",               // RFC3339 with milliseconds
		"2006-01-02T15:04:05",                    // ISO format without timezone
		"2006-01-02 15:04:05.999999 -0700 MST",   // Custom format with timezone
		"2006-01-02 15:04:05.999999 +0000 +0000", // Your specific format
		"2006-01-02 15:04:05",                    // Simple datetime
		"2006-01-02",                             // Date only
	}

	var err error
	for _, format := range formats {
		ct.Time, err = time.Parse(format, str)
		if err == nil {
			return nil
		}
	}

	return err
}

// MarshalJSON implements json.Marshaler interface
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(ct.Time.Format(time.RFC3339))
}

// TipoContribuyenteId model
type TipoContribuyenteId struct {
	Activo            bool       `json:"Activo"`
	CodigoAbreviacion string     `json:"CodigoAbreviacion"`
	Descripcion       string     `json:"Descripcion"`
	FechaCreacion     CustomTime `json:"FechaCreacion"`
	FechaModificacion CustomTime `json:"FechaModificacion"`
	Id                int        `json:"Id"`
	Nombre            string     `json:"Nombre"`
}

// TerceroId model
type TerceroId struct {
	Activo              bool                `json:"Activo"`
	FechaCreacion       CustomTime          `json:"FechaCreacion"`
	FechaModificacion   CustomTime          `json:"FechaModificacion"`
	FechaNacimiento     CustomTime          `json:"FechaNacimiento"`
	Id                  int                 `json:"Id"`
	LugarOrigen         int                 `json:"LugarOrigen"`
	NombreCompleto      string              `json:"NombreCompleto"`
	PrimerApellido      string              `json:"PrimerApellido"`
	PrimerNombre        string              `json:"PrimerNombre"`
	SegundoApellido     string              `json:"SegundoApellido"`
	SegundoNombre       string              `json:"SegundoNombre"`
	TipoContribuyenteId TipoContribuyenteId `json:"TipoContribuyenteId"`
	UsuarioWSO2         string              `json:"UsuarioWSO2"`
}

// TipoDocumentoId model
type TipoDocumentoId struct {
	Activo            bool       `json:"Activo"`
	CodigoAbreviacion string     `json:"CodigoAbreviacion"`
	Descripcion       string     `json:"Descripcion"`
	FechaCreacion     CustomTime `json:"FechaCreacion"`
	FechaModificacion CustomTime `json:"FechaModificacion"`
	Id                int        `json:"Id"`
	Nombre            string     `json:"Nombre"`
	NumeroOrden       int        `json:"NumeroOrden"`
}

// DatosIdentificacion model
type DatosIdentificacion struct {
	Activo             bool            `json:"Activo"`
	CiudadExpedicion   int             `json:"CiudadExpedicion"`
	DigitoVerificacion int             `json:"DigitoVerificacion"`
	DocumentoSoporte   int             `json:"DocumentoSoporte"`
	FechaCreacion      CustomTime      `json:"FechaCreacion"`
	FechaExpedicion    CustomTime      `json:"FechaExpedicion"`
	FechaModificacion  CustomTime      `json:"FechaModificacion"`
	Id                 int             `json:"Id"`
	Numero             string          `json:"Numero"`
	TerceroId          TerceroId       `json:"TerceroId"`
	TipoDocumentoId    TipoDocumentoId `json:"TipoDocumentoId"`
}
