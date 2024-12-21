package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Estructuras para las tablas
type RutaEscolar struct {
	ID         string `bson:"id"`
	Capacidad  int    `bson:"capacidad"`
	Estado     string `bson:"estado"`
	Kilometraje int   `bson:"kilometraje"`
}

type Chofer struct {
	ID       string `bson:"id"`
	Nombre   string `bson:"nombre"`
	Licencia string `bson:"licencia"`
	Estado   string `bson:"estado"`
}

type Asistente struct {
	ID     string `bson:"id"`
	Nombre string `bson:"nombre"`
	Estado string `bson:"estado"`
}

type Apoderado struct {
	ID        string `bson:"id"`
	Nombre    string `bson:"nombre"`
	Direccion string `bson:"direccion"`
	Telefono  string `bson:"telefono"`
}

type Nino struct {
	ID       string `bson:"id"`
	Nombre   string `bson:"nombre"`
	Direccion string `bson:"direccion"`
	Colegio  string `bson:"colegio"`
}

type Contrato struct {
	ID          string  `bson:"id"`
	RutaID      string  `bson:"ruta_id"`
	FechaInicio string  `bson:"fecha_inicio"`
	FechaFin    string  `bson:"fecha_fin"`
	Modalidad   string  `bson:"modalidad"`
	Tarifa      float64 `bson:"tarifa"`
}

type Recorrido struct {
	ID          string `bson:"id"`
	Barrio      string `bson:"barrio"`
	Colegio     string `bson:"colegio"`
	Jornada     string `bson:"jornada"`
	HoraSalida  string `bson:"hora_salida"`
	HoraLlegada string `bson:"hora_llegada"`
}

type Evento struct {
	ID          string `bson:"id"`
	Tipo        string `bson:"tipo"`
	Descripcion string `bson:"descripcion"`
	Fecha       string `bson:"fecha"`
}

// Definiciones de las estructuras
type RutaManager struct{}

func (r *RutaManager) createRuta(ruta RutaEscolar) {
	collection := client.Database("pruebas").Collection("rutas_escolares")
	_, err := collection.InsertOne(context.Background(), ruta)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Ruta escolar creada correctamente")
}

func (r *RutaManager) listRutas() {
	collection := client.Database("pruebas").Collection("rutas_escolares")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var rutas []RutaEscolar
	if err = cursor.All(context.Background(), &rutas); err != nil {
		log.Fatal(err)
	}
	
	if len(rutas) == 0 {
		fmt.Println("No hay rutas escolares disponibles.")
		return
	}

	fmt.Println("Rutas escolares listadas:")
	for _, ruta := range rutas {
		fmt.Printf("ID: %s, Capacidad: %d, Estado: %s, Kilometraje: %d\n", ruta.ID, ruta.Capacidad, ruta.Estado, ruta.Kilometraje)
	}
}

type ContratoManager struct {
	rutaManager RutaManager
}

func (c *ContratoManager) createContrato(contrato Contrato) {
	collection := client.Database("pruebas").Collection("contratos")
	_, err := collection.InsertOne(context.Background(), contrato)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Contrato creado correctamente")
}

func (c *ContratoManager) readContrato(id string) {
	collection := client.Database("pruebas").Collection("contratos")
	var contrato Contrato
	err := collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&contrato)
	if err != nil {
		log.Printf("Contrato con ID %s no encontrado.\n", id)
		return
	}
	fmt.Printf("Contrato encontrado: %+v\n", contrato)
}

func (c *ContratoManager) updateContrato(id string, nuevaTarifa float64, nuevaModalidad, nuevaFechaInicio, nuevaFechaFin string) {
	collection := client.Database("pruebas").Collection("contratos")
	update := bson.M{
		"$set": bson.M{
			"tarifa":       nuevaTarifa,
			"modalidad":    nuevaModalidad,
			"fecha_inicio": nuevaFechaInicio,
			"fecha_fin":    nuevaFechaFin,
		},
	}
	_, err := collection.UpdateOne(context.Background(), bson.M{"id": id}, update)
	if err != nil {
		log.Printf("Error al actualizar el contrato con ID %s: %v\n", id, err)
		return
	}
	fmt.Println("Contrato actualizado correctamente")
}

func (c *ContratoManager) listContratos() {
	collection := client.Database("pruebas").Collection("contratos")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var contratos []Contrato
	if err = cursor.All(context.Background(), &contratos); err != nil {
		log.Fatal(err)
	}
	
	if len(contratos) == 0 {
		fmt.Println("No hay contratos disponibles.")
		return
	}

	fmt.Println("Contratos listados:")
	for _, contrato := range contratos {
		fmt.Printf("ID: %s, RutaID: %s, Fecha Inicio: %s, Fecha Fin: %s, Modalidad: %s, Tarifa: %.2f\n", contrato.ID, contrato.RutaID, contrato.FechaInicio, contrato.FechaFin, contrato.Modalidad, contrato.Tarifa)
	}
}

func (c *ContratoManager) searchContratos(id string) {
	collection := client.Database("pruebas").Collection("contratos")
	var contrato Contrato
	err := collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&contrato)
	if err != nil {
		log.Printf("Contrato con ID %s no encontrado.\n", id)
		return
	}
	fmt.Printf("Contrato encontrado: %+v\n", contrato)
}

func (c *ContratoManager) deleteContrato(id string) {
	collection := client.Database("pruebas").Collection("contratos")
	result, err := collection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		log.Printf("Error al eliminar el contrato con ID %s: %v\n", id, err)
		return
	}
	if result.DeletedCount == 0 {
		fmt.Printf("No se encontró ningún contrato con ID %s para eliminar.\n", id)
		return
	}
	fmt.Println("Contrato eliminado correctamente")
}

func main() {
	// Establecer la conexión a MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://Sergio-All:OwSA0rjmhJKtyE7M@cluster0.ztbdwxd.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	// Verificar la conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Conexión exitosa a MongoDB")

	// Inicializar los managers
	rutaManager := RutaManager{}
	contratoManager := ContratoManager{rutaManager: rutaManager}

	// Menú interactivo
	for {
		fmt.Println("\nSeleccione una opción:")
		fmt.Println("1. Crear Ruta Escolar")
		fmt.Println("2. Listar Rutas Escolares")
		fmt.Println("3. Crear Contrato")
		fmt.Println("4. Leer Contrato")
		fmt.Println("5. Actualizar Contrato")
		fmt.Println("6. Listar Contratos")
		fmt.Println("7. Buscar Contrato")
		fmt.Println("8. Eliminar Contrato")
		fmt.Println("9. Salir")

		var opcion int
		fmt.Scan(&opcion)

		switch opcion {
		case 1:
			var ruta RutaEscolar
			fmt.Print("Ingrese ID de la ruta: ")
			fmt.Scan(&ruta.ID)
			fmt.Print("Ingrese capacidad de la ruta: ")
			fmt.Scan(&ruta.Capacidad)
			fmt.Print("Ingrese estado de la ruta: ")
			fmt.Scan(&ruta.Estado)
			fmt.Print("Ingrese kilometraje de la ruta: ")
			fmt.Scan(&ruta.Kilometraje)
			rutaManager.createRuta(ruta)

		case 2:
			rutaManager.listRutas()

		case 3:
			var contrato Contrato
			fmt.Print("Ingrese ID del contrato: ")
			fmt.Scan(&contrato.ID)
			fmt.Print("Ingrese ID de la ruta asociada: ")
			fmt.Scan(&contrato.RutaID)
			fmt.Print("Ingrese fecha de inicio del contrato: ")
			fmt.Scan(&contrato.FechaInicio)
			fmt.Print("Ingrese fecha de fin del contrato: ")
			fmt.Scan(&contrato.FechaFin)
			fmt.Print("Ingrese modalidad del contrato: ")
			fmt.Scan(&contrato.Modalidad)
			fmt.Print("Ingrese tarifa del contrato: ")
			fmt.Scan(&contrato.Tarifa)
			contratoManager.createContrato(contrato)

		case 4:
			var id string
			fmt.Print("Ingrese ID del contrato a leer: ")
			fmt.Scan(&id)
			contratoManager.readContrato(id)

		case 5:
			var id string
			var nuevaTarifa float64
			var nuevaModalidad, nuevaFechaInicio, nuevaFechaFin string
			fmt.Print("Ingrese ID del contrato a actualizar: ")
			fmt.Scan(&id)
			fmt.Print("Ingrese nueva tarifa del contrato: ")
			fmt.Scan(&nuevaTarifa)
			fmt.Print("Ingrese nueva modalidad del contrato: ")
			fmt.Scan(&nuevaModalidad)
			fmt.Print("Ingrese nueva fecha de inicio del contrato: ")
			fmt.Scan(&nuevaFechaInicio)
			fmt.Print("Ingrese nueva fecha de fin del contrato: ")
			fmt.Scan(&nuevaFechaFin)
			contratoManager.updateContrato(id, nuevaTarifa, nuevaModalidad, nuevaFechaInicio, nuevaFechaFin)

		case 6:
			contratoManager.listContratos()

		case 7:
			var id string
			fmt.Print("Ingrese ID del contrato a buscar: ")
			fmt.Scan(&id)
			contratoManager.searchContratos(id)

		case 8:
			var id string
			fmt.Print("Ingrese ID del contrato a eliminar: ")
			fmt.Scan(&id)
			contratoManager.deleteContrato(id)

		case 9:
			fmt.Println("Saliendo...")
			return

		default:
			fmt.Println("Opción no válida. Intente de nuevo.")
		}
	}
}
  