package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:CargaAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:CargaAcademicaController"],
        beego.ControllerComments{
            Method: "PostCargaAcademica",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:EspacioAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:EspacioAcademicoController"],
        beego.ControllerComments{
            Method: "PostEspacioAcademico",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Formulario_por_tipoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Formulario_por_tipoController"],
        beego.ControllerComments{
            Method: "GetFormularioTipo",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Formulario_por_tipoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Formulario_por_tipoController"],
        beego.ControllerComments{
            Method: "PostFormularioTipo",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:MetricasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:MetricasController"],
        beego.ControllerComments{
            Method: "MetricasAutoevaluacion",
            Router: "/Autoevaluacion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:MetricasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:MetricasController"],
        beego.ControllerComments{
            Method: "MetricasCoevaluacion",
            Router: "/Coevaluacion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:MetricasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:MetricasController"],
        beego.ControllerComments{
            Method: "MetricasHeteroevaluacion",
            Router: "/Heteroevaluacion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Respuesta_formularioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Respuesta_formularioController"],
        beego.ControllerComments{
            Method: "PostRespuestaFormulario",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
