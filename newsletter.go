package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/mitchellh/goamz/s3"
	uuid "github.com/satori/go.uuid"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Newsletter struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	Date        string        `bson:"date"`
	Image       string        `bson:"image"`
	BackImage   string        `bson:"backimage"`
	LetterNo    string        `bson:"letterno"`
	Type        string        `bson:"type"`
}

type NewsletterRepo struct {
	coll *mgo.Collection
}

func (r *NewsletterRepo) UploadLetter(nl Newsletter) error {
	err := r.coll.Insert(nl)

	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func (r *NewsletterRepo) GetLetters() ([]Newsletter, error) {
	data := []Newsletter{}
	err := r.coll.Find(bson.M{}).Sort("-date").All(&data)
	if err != nil {
		log.Println(err)
		return data, err
	}
	return data, nil

}

/*func (r *NewsletterRepo) GetPackages() ([]Newsletter, error) {
	data := []Newsletter{}
	err := r.coll.Find(bson.M{}).All(&data)
	if err != nil {
		log.Println(err)
		return data, err
	}
	return data, nil

}*/

/*func (c *Config) GetPackagesHandler(w http.ResponseWriter, r *http.Request) {
	u := NewsletterRepo{c.MongoSession.DB(c.MONGODB).C("packages")}
	data, err := u.GetPackages()
	if err != nil {
		log.Println(err)
	}
	res, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}*/

func (c *Config) GetLetterHandler(w http.ResponseWriter, r *http.Request) {
	u := NewsletterRepo{c.MongoSession.DB(c.MONGODB).C("newsletter")}
	data, err := u.GetLetters()
	if err != nil {
		log.Println(err)
	}
	res, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (c *Config) UploadHandler(w http.ResponseWriter, r *http.Request) {
	data := Newsletter{}
	u := NewsletterRepo{c.MongoSession.DB(c.MONGODB).C("newsletter")}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}
	if data.Image != "" {
		bucket := c.S3Bucket
		byt, err := base64.StdEncoding.DecodeString(strings.Split(data.Image, "base64,")[1])
		if err != nil {
			log.Println(err)
		}
		meta := strings.Split(data.Image, "base64,")[0]
		newmeta := strings.Replace(strings.Replace(meta, "data:", "", -1), ";", "", -1)
		name := uuid.NewV4().String() + data.LetterNo
		err = bucket.Put(name, byt, newmeta, s3.PublicReadWrite)
		if err != nil {
			log.Println(err)
		}
		data.Image = bucket.URL(name)
		if data.BackImage != "" {
			bucket := c.S3Bucket
			byt, err = base64.StdEncoding.DecodeString(strings.Split(data.BackImage, "base64,")[1])
			if err != nil {
				log.Println(err)
			}
			meta := strings.Split(data.BackImage, "base64,")[0]
			newmeta := strings.Replace(strings.Replace(meta, "data:", "", -1), ";", "", -1)
			name := uuid.NewV4().String() + data.LetterNo + "2"
			err = bucket.Put(name, byt, newmeta, s3.PublicReadWrite)
			if err != nil {
				log.Println(err)
			}
			data.BackImage = bucket.URL(name)
		}
		if err != nil {
			log.Println(err)
		}
	}
	err = u.UploadLetter(data)

}
