package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func cargarHtml(a string) string {
	html, _ := ioutil.ReadFile(a)

	return string(html)
}

func form(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("form.html"),
	)
}

var materias = make(map[string]map[string]float64)
var alumnos = make(map[string]map[string]float64)

var alum []string
var mat []string

func existeAlumno(nombre string) int64 {
	for _, name := range alum {
		if name == nombre {
			return 1
		}
	}
	return 0
}

func existeMateria(nombre string) int64 {
	for _, name := range mat {
		if name == nombre {
			return 1
		}
	}
	return 0
}

func existeCalificacion(nombre string, mate string) int64 {
	if alumnos[nombre][mate] == 0 {
		return 0
	}
	return 1
}

var strPromediogeneral string

func PromedioGeneral() string {
	var promedio float64
	var numMaterias float64

	var numAlumnos float64
	numAlumnos = 0
	var promediogeneral float64
	promediogeneral = 0

	var html string

	for _, name := range alum {
		promedio = 0
		numMaterias = 0
		for materia, calificacion := range alumnos[name] {
			calS := fmt.Sprintf("%f", calificacion)
			html += "<tr>" +
				"<td>" + name + "</td>" +
				"<td>" + materia + "</td>" +
				"<td>" + calS + "</td>" +
				"</tr>"
			promedio = promedio + calificacion
			numMaterias = numMaterias + 1
		}
		promedio = promedio / numMaterias
		promediogeneral = promediogeneral + promedio
		numAlumnos = numAlumnos + 1
	}
	promediogeneral = promediogeneral / numAlumnos
	proS := fmt.Sprintf("%.2f", promediogeneral)
	strPromediogeneral = proS
	return html
}

func regCalificacion(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		nombre := req.FormValue("alumnoCali")
		mate := req.FormValue("materiaCali")
		calificacionString := req.FormValue("calificacion")
		calificacion, _ := strconv.ParseFloat(calificacionString, 64)

		if existeCalificacion(nombre, mate) == 0 {
			bandA := 1
			bandM := 1
			if existeAlumno(nombre) == 0 {
				bandA = 0
			}
			if existeMateria(mate) == 0 {
				bandM = 0
			}

			if bandM == 0 {
				mat = append(mat, mate)
				var alumno = make(map[string]float64)
				alumno[nombre] = calificacion
				materias[mate] = alumno
			} else {
				materias[mate][nombre] = calificacion
			}

			if bandA == 0 {
				alum = append(alum, nombre)
				var materia = make(map[string]float64)
				materia[mate] = calificacion
				alumnos[nombre] = materia
			} else {
				alumnos[nombre][mate] = calificacion
			}

			fmt.Println(alumnos)
			fmt.Println(materias)

			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				cargarHtml("respuestaRP.html"),
				nombre,
				mate,
			)
		} else {
			fmt.Println("El alumno ya tiene calificación")

			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				cargarHtml("respuestaRN.html"),
				nombre,
				mate,
			)
		}

	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("tablaProG.html"),
			PromedioGeneral(),
			strPromediogeneral,
		)
	}
}

func promedioAlumno(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		nombre := req.FormValue("promeAlum")

		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("tablaProG.html"),
			PromedioAlumnoP(nombre),
			strPromediogeneral,
		)
	}
}

func promedioMateria(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		nombre := req.FormValue("promeMate")

		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("tablaProG.html"),
			PromedioMateriaP(nombre),
			strPromediogeneral,
		)
	}
}

func PromedioAlumnoP(nombre string) string {
	var promedio float64
	var numMaterias float64
	promedio = 0
	numMaterias = 0

	var html string

	for materia, calificacion := range alumnos[nombre] {
		calS := fmt.Sprintf("%f", calificacion)
		html += "<tr>" +
			"<td>" + nombre + "</td>" +
			"<td>" + materia + "</td>" +
			"<td>" + calS + "</td>" +
			"</tr>"
		promedio = promedio + calificacion
		numMaterias = numMaterias + 1
	}

	promedio = promedio / numMaterias
	proS := fmt.Sprintf("%.2f", promedio)
	strPromediogeneral = proS
	return html
}

func PromedioMateriaP(mate string) string {
	var promedio float64
	var numAlumnos float64
	promedio = 0
	numAlumnos = 0

	var html string

	for alumno, calificacion := range materias[mate] {
		calS := fmt.Sprintf("%f", calificacion)
		html += "<tr>" +
			"<td>" + alumno + "</td>" +
			"<td>" + mate + "</td>" +
			"<td>" + calS + "</td>" +
			"</tr>"
		promedio = promedio + calificacion
		numAlumnos = numAlumnos + 1
	}
	promedio = promedio / numAlumnos
	proS := fmt.Sprintf("%.2f", promedio)
	strPromediogeneral = proS
	return html
}

// http://localhost:9000/form
//Promedio General: http://localhost:9000/regCalificacion
func main() {
	http.HandleFunc("/form", form)
	http.HandleFunc("/regCalificacion", regCalificacion)
	http.HandleFunc("/promedioAlumno", promedioAlumno)
	http.HandleFunc("/promedioMateria", promedioMateria)
	fmt.Println("Corriendo servirdor de tareas...")
	http.ListenAndServe(":9000", nil)
}
