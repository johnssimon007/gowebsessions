package main
import(
  "log"
  "net/http"
  "html/template"
  "github.com/gorilla/sessions"
  "math/rand"

)
var (
	key = []byte("silverhawk007")
	store = sessions.NewCookieStore(key)

)
type Employee struct{
  Name string
  Company string
  Staffid int
  Email string
  Status bool
  Address *Location
}
type Location struct{
   City string
   State string
   Country string
   Status bool

}
type Error struct{
  Err string
  Status bool
}

func middleware(f http.HandlerFunc) http.HandlerFunc{
  return func(w http.ResponseWriter,r *http.Request){
  if r.Method!="GET"{
        tmpl:=template.Must(template.ParseFiles("./assets/index.html"))
        err:=Error{Status:false,Err:"unsupported Method Type"}
        parse:=tmpl.Execute(w,err)
        log.Println(parse)
  }
  session, _ := store.Get(r, "datr")
  if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

      tmpl:=template.Must(template.ParseFiles("./assets/Error.html"))
      err:=Error{Status:false,Err:"you are not authenticated,generate a key first"}
      parse:=tmpl.Execute(w,err)
      log.Println(parse)
    return
}else{
  f(w,r)
}
  }
}
func users(w http.ResponseWriter,r *http.Request){
  tmpl:=template.Must(template.ParseFiles("./assets/index.html"))
  Employee:=Employee{Status:true,Name:"johns",Company:"NetObjex",Staffid:rand.Intn(100000),Email:"johnssimon@gmail.com",Address:&Location{Status:true,City:"Phnom penh",State:"Cambodia",Country:"Cambodia"}}
  errs:=tmpl.Execute(w,Employee)
  log.Println(errs)

}
func Key(w http.ResponseWriter,r *http.Request){
  c := make(chan bool)
  go func() {
  session,_:=store.Get(r,"datr")
  session.Values["authenticated"] = true
  err:=session.Save(r,w)
  if err!=nil{
    c<-false
  }else{
  c <-true
}
  }()
  select {
    case res := <-c:
       if res==true{
        w.Write([]byte("created session"))
        log.Println(res)
      }else{

         w.Write([]byte("error while creating session"))
         log.Println(res)
      }
    }

}
func main(){
  http.HandleFunc("/",middleware(users))
  http.HandleFunc("/GenerateKey",Key)
  http.ListenAndServe(":8080",nil)
}

