package routes

import (
	"fmt"
	"net/http"
	"projet-groupie/controllers"
)

func Init() {
	//Permet de réduperai tout les fichier static type css,img,js,music.
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/accueil", controllers.HomeControler)
	http.HandleFunc("/recherche", controllers.SearchControler)

	fmt.Println("Le serveur est opérationel : http://localhost:8000")
	http.ListenAndServe("localhost:8000", nil)
}
