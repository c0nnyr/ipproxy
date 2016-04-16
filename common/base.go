package common

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

// mongo
type IPProxyItem struct {
	IP              string        `bson:"ip"`
	Country         string        `bson:"country"`
	Hide_Type       string        `bson:"hide_type"`
	Connection_Type string        `bson:"connection_type"`
	ID              bson.ObjectId `bson:"_id"`
}

const MGO_URL = "127.0.0.1:27017"
const MDB_NAME = "ipproxy_db"
const COLLECTION_ITEMS = "ip_item"

var mongo_basic_sesstion *mgo.Session = nil

func get_session() *mgo.Session {
	if mongo_basic_sesstion == nil {
		var err error
		mongo_basic_sesstion, err = mgo.Dial(MGO_URL)
		if err != nil {
			panic(err)
		}
	}
	return mongo_basic_sesstion.Clone()
}

func handle_with_collection(connection_name string, handler func(*mgo.Collection) error) error {
	session := get_session()
	defer session.Close()
	c := session.DB(MDB_NAME).C(connection_name)
	return handler(c)

}

func UpsertItemToDB(item IPProxyItem) {
	item_tmp := IPProxyItem{}
	err := handle_with_collection(COLLECTION_ITEMS, func(c *mgo.Collection) error {
		return c.Find(bson.M{"ip": item.IP}).One(&item_tmp)
	})
	if err != nil {
		//log.Println("Insert ip ", item.IP)
		err2 := add_item_to_db(item)
		if err2 != nil {
			//log.Println("Cannot insert ip", item.IP, err2)
			log.Fatal(err2)
		}
	} else {
		item.ID = item_tmp.ID
		//log.Println("Update ip", item.IP)
		err3 := update_item_to_db(item)
		if err3 != nil {
			//log.Println("Cannot update ip", item.IP, err3)
			log.Fatal(err3)
		}
	}
}

func update_item_to_db(item IPProxyItem) error {
	return handle_with_collection(COLLECTION_ITEMS, func(c *mgo.Collection) error {
		return c.UpdateId(item.ID, item)
	})
}
func add_item_to_db(item IPProxyItem) error {
	item.ID = bson.NewObjectId()
	return handle_with_collection(COLLECTION_ITEMS, func(c *mgo.Collection) error {
		return c.Insert(item)
	})
}
