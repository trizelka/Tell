
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Tell/julienschmidt/httprouter"
	"Tell/trizelka/Messaging/models"
	"Tell/trizelka/Messaging/gcm"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// UserController represents the controller for operating on the User resource
	UserController struct {
		session *mgo.Session
	}
)

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

func sendGCMMessages(code int, from string, to string, messages string){
        data := map[string]interface{}{"code": code, "snumber":from,  "dnumber":to, "messages":messages}

        client := gcm.New("AIzaSyB5Nnu1fM0y_nDIJzVt7PHrlel2mGvMa_s")
        load := gcm.NewMessage("1021501583544")
        load.AddRecipient(to)
        load.SetPayload("data",data)
        load.CollapseKey = "trizelka"
        load.DelayWhileIdle = false
        load.TimeToLive = 10

        resp, err := client.Send(load)

        fmt.Printf("id: %+v\n", resp)
        fmt.Println("err:", err)
        fmt.Println("err index:", resp.ErrorIndexes())
        fmt.Println("reg index:", resp.RefreshIndexes())
}


// GetUser retrieves an individual user resource
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab phone
	phone := p.ByName("phone")

	// Stub user
	u := models.User{}

	// Fetch user
	if err := uc.session.DB("hereiam").C("users").Find(bson.M{"phone":phone}).One(&u); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
	fmt.Printf("%s", uj)
}

// send Messages user retrieves an individual user resource
func (uc UserController) SendMessages(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
        messages := p.ByName("messages")
	from := p.ByName("from")
	
	// Stub user
        u := models.User{}

	// Populate the user data
        json.NewDecoder(r.Body).Decode(&u)

        // Fetch user
        if err := uc.session.DB("hereiam").C("users").Find(bson.M{"phone":u.Phone}).One(&u); err != nil {
                w.WriteHeader(404)
                return
        }

        sendGCMMessages(1,from,u.Regid,messages);

        // Marshal provided interface into JSON structure
        uj, _ := json.Marshal(u)

        // Write content-type, statuscode, payload
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        fmt.Fprintf(w, "%s", uj)
        fmt.Printf("%s", uj)
}


// CreateUser creates a new user resource
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab http params
        mth := p.ByName("methode")

	// Stub an user to be populated from the body
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.Id = bson.NewObjectId()

	if mth == "add_user" {uc.session.DB("hereiam").C("users").Insert(u)}
	if mth == "update_regid" {
		if err := uc.session.DB("hereiam").C("users").Find(bson.M{"phone":u.Phone}).One(&u); err != nil {
                w.WriteHeader(404)
                return
	        }
 
		query := bson.M{"name":u.Name,"email":u.Email,"phone":u.Phone}
        	update := bson.M{"regid":u.Regid}
 		// Update Register id user
        	if err := uc.session.DB("hereiam").C("users").Update(query,update); err != nil {
                	w.WriteHeader(404)
                	return
        	}
	}


	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
	fmt.Printf("%s", uj)
}

// RemoveUser removes an existing user resource
func (uc UserController) RemoveUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove user
	if err := uc.session.DB("hereiam").C("users").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}

