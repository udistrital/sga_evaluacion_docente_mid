package services

import (
	"fmt"

	"github.com/udistrital/utils_oas/requestresponse"
)

type FormularioService struct {
	// Aquí puedes definir los campos y métodos necesarios para tu servicio
}

func ConsultaFormulario(id_tipo_formulario string, id_periodo string, id_tercero string, id_espacio string) (APIResponseDTO requestresponse.APIResponse) {
	fmt.Println("Inicializando el servicio de formulario...")
	return requestresponse.APIResponseDTO(false, 200, "", "")
}

// Aquí puedes agregar los métodos y funcionalidades de tu servicio
