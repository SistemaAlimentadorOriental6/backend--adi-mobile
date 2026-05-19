package domain

type Empleado struct {
	IdEmpleado             string  `json:"IdEmpleado"`
	Cedula                 string  `json:"Cedula"`
	Nombres                string  `json:"Nombres"`
	Direccion              string  `json:"Direccion"`
	Telefono               string  `json:"Telefono"`
	Celular                string  `json:"Celular"`
	CodigoOperador         string  `json:"CodigoOperador"`
	IdEmpresa              string  `json:"IdEmpresa"`
	Estado                 *string `json:"Estado"`
	Puntaje                string  `json:"Puntaje"`
	FechaEntrada           string  `json:"FechaEntrada"`
	Foto                   string  `json:"Foto"`
	IdEstado               string  `json:"IdEstado"`
	Email                  string  `json:"Email"`
	NumeroInterno          string  `json:"NumeroInterno"`
	FechaContrato          string  `json:"FechaContrato"`
	NumeroPase             string  `json:"NumeroPase"`
	CategoriaLicencia      string  `json:"CategoriaLicencia"`
	TarjetaEnteGestor      string  `json:"TarjetaEnteGestor"`
	TarjetaEnteGestorVence string  `json:"TarjetaEnteGestorVence"`
	IdCargo                string  `json:"IdCargo"`
	Cargo                  string  `json:"Cargo"`
}

type EmpleadoApiResponse struct {
	Status bool       `json:"status"`
	Count  int        `json:"count"`
	Data   []Empleado `json:"data"`
}

type ProxyRepository interface {
	GetEmpleados() (*EmpleadoApiResponse, error)
}
