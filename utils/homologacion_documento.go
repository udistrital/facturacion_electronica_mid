package utils

// MapeoTipoDocumentoTerceroASGA contiene el mapeo de IDs de tipo documento de la tabla terceros
// a los códigos de variables (TDO_CODVAR) utilizados en SGA
var MapeoTipoDocumentoTerceroASGA = map[int]string{
	1:  "R", // REGISTRO CIVIL DE NACIMIENTO → REGISTRO CIVIL
	2:  "T", // TARJETA DE IDENTIDAD → TARJETA DE IDENTIDAD
	3:  "C", // CÉDULA DE CIUDADANÍA → CÉDULA DE CIUDADANÍA
	4:  "O", // CERTIFICADO REGISTRADURÍA SIN IDENTIFICACIÓN → OTRO
	5:  "O", // TARJETA DE EXTRANJERÍA → OTRO
	6:  "E", // CÉDULA DE EXTRANJERÍA → CÉDULA DE EXTRANJERÍA
	7:  "N", // NIT → NIT
	8:  "O", // IDENTIFICACIÓN DE EXTRANJEROS DIFERENTE AL NIT ASIGNADO DIAN → OTRO
	9:  "P", // PASAPORTE → PASAPORTE
	10: "O", // DOCUMENTO DE IDENTIFICACIÓN EXTRANJERO → OTRO
	11: "O", // SIN IDENTIFICACIÓN DEL EXTERIOR → OTRO
	12: "O", // DOCUMENTO DE IDENTIFICACIÓN EXTRANJERO PERSONA JURÍDICA → OTRO
	13: "O", // CARNÉ DIPLOMÁTICO → OTRO
	14: "O", // CARNÉ ESTUDIANTE → OTRO
}

// ObtenerTipoDocumentoSGA convierte un ID de tipo documento de la tabla terceros
// al código de variable (TDO_CODVAR) utilizado en SGA
// Retorna el código SGA y un booleano indicando si se encontró el mapeo
func ObtenerTipoDocumentoSGA(idTipoDocumentoTercero int) (string, bool) {
	codigo, existe := MapeoTipoDocumentoTerceroASGA[idTipoDocumentoTercero]
	return codigo, existe
}
