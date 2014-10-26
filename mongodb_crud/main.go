package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Person struct {
	Name  string `json:"name" bson:"name"`
	Phone string `json:"phone" bson:"phone"`
}

func (p *Person) readOneFromPointer(c *mgo.Collection) {
	err := c.Find(bson.D{{"name", p.Name}}).One(p)
	if err != nil {
		if err == mgo.ErrNotFound {
			log.Fatal("Not found")
		} else {
			log.Fatal("Internal error")
		}
	}
	fmt.Println("Phone:", p.Phone)
}

func main() {
	session, err := mgo.Dial("172.17.8.101")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	//db.people.update({name: "Ale"}, { $set: {phone: "+55 53 8116 9639"}}, {upsert: true})
	_, err = c.Upsert(bson.M{"name": "Ale"}, &Person{Name: "Ale", Phone: "+55 53 8116 9639"})
	//db.people.update({name: "Cla"}, { $set: {phone: "55 53 8402 8510"}}, {upsert: true})
	_, err = c.Upsert(bson.M{"name": "Cla"}, &Person{Name: "Cla", Phone: "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	p := &Person{Name: "Ale"}
	p.readOneFromPointer(c)
}
