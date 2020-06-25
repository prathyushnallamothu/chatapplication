package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)
type Data struct{
	Username string `json:"username"`
	Message string `json:"message"`
}
var upgrader=websocket.Upgrader{
	ReadBufferSize:1024,
	WriteBufferSize:1024,
	CheckOrigin:func(r *http.Request)bool{
		return true
	},
}
var msgtype int
var msgchan = make(chan Data)
var clients = make(map[*websocket.Conn]bool)
var tpl *template.Template
func init(){
	tpl=template.Must(template.ParseGlob("templates/*.html"))
}
func main(){
m:=mux.NewRouter().StrictSlash(true)
m.HandleFunc("/",homehandler)
m.HandleFunc("/login",loginhandler)
m.HandleFunc("/chat",chathandler)
m.HandleFunc("/ws",websockethandler)
go writer()
http.ListenAndServe(":8080",m)
}
func homehandler(w http.ResponseWriter,r *http.Request){
	tpl.ExecuteTemplate(w,"login.html",nil)
}
func loginhandler(w http.ResponseWriter,r *http.Request){
	q:=r.FormValue("name")
	if q==""{
		fmt.Fprintf(w,"Go back and enter a valid name!!!")
	}else{
	http.Redirect(w,r,"/chat?q="+q,307)
	}
}
func chathandler(w http.ResponseWriter,r *http.Request){
	q:=r.FormValue("q")
	tpl.ExecuteTemplate(w,"chat.html",q)
}
func websockethandler(w http.ResponseWriter,r *http.Request){
	ws,err:=upgrader.Upgrade(w,r,nil)
	if err!=nil{
		log.Fatal(err)
	}
	clients[ws]=true
	for{
		var jsonmsg Data
		err := ws.ReadJSON(&jsonmsg)
		if err != nil {
			log.Println(err)
			return
		}
	msgchan<-jsonmsg
	//msgtype=messagetype

	}

}
func writer(){
	for {
		messgs := <-msgchan
		for client := range clients {
	err1:=client.WriteJSON(messgs)
	if err1!=nil{
		log.Fatal(err1)
	}
}
}
}
