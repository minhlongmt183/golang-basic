package main

// import all necessary packages we get in our previous go get -u
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var db, _ = gorm.Open("mysql", "root:root@/todolist?charset=utf8&parseTime=True&loc=Local")

// defining the TodoItem model as what we described earlier
type TodoItemModel struct {
	Id          int `gorm:"primary_key"`
	Description string
	Completed   bool
}

// create Healthz function that will respond {"alive": true} to client every time it's invoked
func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Healh is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

// set our function to setup our logurs logger setting.
func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	defer db.Close()

	// we are running automigration against our MySQL databse immediately after
	// starting our API server.
	db.Debug().DropTableIfExists(&TodoItemModel{})
	db.Debug().AutoMigrate(&TodoItemModel{})

	log.Info("Starting Todolist API server")

	// init our gorilla/mx HTTP router with a walrus operator. We route `/healthz` HTTP GET requests to Health() function.
	// The router will listen to port 8000
	router := mux.NewRouter()
	router.HandleFunc("/c", Healthz).Methods("GET")
	http.ListenAndServe(":8000", router)
}
