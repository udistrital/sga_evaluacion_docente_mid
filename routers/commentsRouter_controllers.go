package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Formulario_por_tipoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Formulario_por_tipoController"],
        beego.ControllerComments{
            Method: "GetFormularioTipo",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:MetricasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:MetricasController"],
        beego.ControllerComments{
            Method: "MetricasHeteroevaluacion",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:MetricasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:MetricasController"],
        beego.ControllerComments{
            Method: "MetricasAutoevaluacion",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:MetricasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:MetricasController"],
        beego.ControllerComments{
            Method: "MetricasCoevaluacion",
            Router: "/",
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
