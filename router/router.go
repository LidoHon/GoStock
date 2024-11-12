package router

import (
	controller "github.com/LidoHon/GoStock/Controller"
	"github.com/gorilla/mux"
)



func Router() *mux.Router {
router := mux.NewRouter()
router.HandleFunc("/api/stock/{id}", controller.GetStock).Methods("GET","OPTIONS")
router.HandleFunc("/api/newstock", controller.CreateStock).Methods("POST", "OPTIONS")
router.HandleFunc("/api/stock", controller.GetAllStock).Methods("GET", "OPTIONS")
router.HandleFunc("/api/stock/{id}", controller.DeleteStock).Methods("DELETE", "OPTIONS")

router.HandleFunc("/api/stock/{id}", controller.UpdateStock).Methods("PUT", "OPTIONS")

return router
}