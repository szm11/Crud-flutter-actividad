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
    ID        string `bson:"id"`
    Capacidad int    `bson:"capacidad"`
    Estado    string `bson:"estado"`
    Kilometraje int  `bson:"kilometraje"`
    ChoferID  string `bson:"chofer_id"`
}

type Chofer struct {
    ID      string `bson:"id"`
    Nombre  string `bson:"nombre"`
    Licencia string `bson:"licencia"`
    Estado  string `bson:"estado"`
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

    // Operaciones CRUD
    realizarOperacionesCRUD()
}

func realizarOperacionesCRUD() {
    // Insertar un nuevo chofer
    chofer := Chofer{ID: "1", Nombre: "Juan Perez", Licencia: "ABC123", Estado: "Activo"}
    insertChofer(chofer)

    // Insertar una nueva ruta escolar
    ruta := RutaEscolar{ID: "1", Capacidad: 50, Estado: "Disponible", Kilometraje: 1000, ChoferID: "1"}
    insertRuta(ruta)

    // Consultar una ruta escolar por su ID
    consultarRuta("1")
}

func insertChofer(chofer Chofer) {
    collection := client.Database("pruebas").Collection("choferes")
    _, err := collection.InsertOne(context.Background(), chofer)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Chofer insertado correctamente")
}

func insertRuta(ruta RutaEscolar) {
    collection := client.Database("pruebas").Collection("rutas_escolares")
    _, err := collection.InsertOne(context.Background(), ruta)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Ruta escolar insertada correctamente")
}

func consultarRuta(id string) {
    collection := client.Database("pruebas").Collection("rutas_escolares")
    var ruta RutaEscolar
    err := collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&ruta)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Ruta escolar encontrada: %+v\n", ruta)
}
