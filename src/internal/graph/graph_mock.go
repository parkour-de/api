package graph

import (
	"fmt"
	"log"
	"math/rand"
	"pkv/api/src/domain"
	"pkv/api/src/internal/dpv"
	"time"
)

var author domain.User

func CreateMultiple[T any](db *Db, count int, createFunc func(*Db, int) (T, error)) ([]T, error) {
	var entities []T
	for i := 1; i <= count; i++ {
		entity, err := createFunc(db, i)
		if err != nil {
			return nil, fmt.Errorf("could not create multiple entities: %w", err)
		}
		entities = append(entities, entity)
	}
	return entities, nil
}

func CreateDescriptions(text string) domain.Descriptions {
	german := CreateDescription("de-DE", text, false)
	english := CreateDescription("en-GB", text, true)
	return domain.Descriptions{
		"de-DE": german,
		"en-GB": english,
	}
}

func CreateDescription(language string, text string, translated bool) domain.Description {
	switch language {
	case "de-DE":
		return domain.Description{
			text + " - Eine tolle Sache",
			"Das wird euch sicher ganz gut gefallen",
			translated,
		}
		break
	case "en-GB":
		return domain.Description{
			text + " - A cool thing",
			"That will be super cool",
			translated,
		}
		break
	}
	return domain.Description{}
}

func CreateComments(_db *Db, count int) ([]domain.Comment, error) {
	createFunc := func(db *Db, i int) (domain.Comment, error) {
		return CreateComment(), nil
	}
	return CreateMultiple(_db, count, createFunc)
}

func CreateComment() domain.Comment {
	return domain.Comment{
		[]string{"Geil", "Super", "Klasse", "Wahnsinn"}[rand.Intn(4)],
		[]string{"Da muss ich unbedingt hin", "Immer wieder schön hier", "Ich liebe es einfach", "Kann ich nicht genug von"}[rand.Intn(4)],
		"admin",
		time.Now(),
	}
}

func CreateUser(db *Db, i int) (domain.User, error) {
	user := domain.User{}
	user.Name = fmt.Sprintf("Test User %d", i)
	user.Type = "person"
	user.Information = map[string]string{"email": "john.doe@example.com", "twitter": "johndoe"}
	var err error
	user.Photos, err = CreateMultiple(db, 5, CreatePhoto)
	if err != nil {
		return user, err
	}
	user.Comments, err = CreateComments(db, rand.Intn(5))
	if err != nil {
		return user, err
	}
	err = db.Users.Create(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func CreateTraining(db *Db, i int) (domain.Training, error) {
	training := domain.Training{}
	training.Descriptions = CreateDescriptions(fmt.Sprintf("Test Training %d", i))
	var err error
	training.Comments, err = CreateComments(db, rand.Intn(5))
	if err != nil {
		return training, err
	}
	training.Photos, err = CreateMultiple(db, 5, CreatePhoto)
	if err != nil {
		return training, err
	}
	training.Cycles, err = CreateMultiple(db, 2, CreateCycle)
	err = db.Trainings.Create(&training)
	if err != nil {
		return training, err
	}
	location, err := CreateLocation(db, i)
	if err != nil {
		return training, err
	}
	err = db.ConnectTrainingLocation(&training, &location)
	if err != nil {
		return training, err
	}
	return training, nil
}

func CreateLocation(db *Db, i int) (domain.Location, error) {
	location := domain.Location{}
	var err error
	location.Descriptions = CreateDescriptions(fmt.Sprintf("Test Location %d", i))
	location.Comments, err = CreateComments(db, rand.Intn(5))
	if err != nil {
		return location, err
	}
	if i < 40 {
		location.City = "Hamburg"
		location.Lat = 53.55
		location.Lng = 9.99
	} else {
		location.City = "Berlin"
		location.Lat = 52.52
		location.Lng = 13.40
	}
	location.Photos, err = CreateMultiple(db, 5, CreatePhoto)
	if err != nil {
		return location, err
	}
	err = db.Locations.Create(&location)
	if err != nil {
		return location, err
	}
	return location, nil
}

func CreatePhoto(db *Db, i int) (domain.Photo, error) {
	photo := domain.Photo{}
	photo.Src = fmt.Sprintf("photo_%d.jpg", i)
	photo.W = 640
	photo.H = 480
	return photo, nil
}

func CreateCycle(db *Db, i int) (domain.Cycle, error) {
	cycle := domain.Cycle{}
	cycle.Weekday = rand.Intn(7) + 1
	return cycle, nil
}

func NewTestDB() *Db {
	config, err := dpv.NewConfig("../../config.yml")
	if err != nil {
		log.Fatal(err)
	}
	c, err := Connect(config, true)
	if err != nil {
		log.Fatal(err)
	}
	database, err := GetOrCreateDatabase(c, "dpv", config)
	if err != nil {
		log.Fatal(err)
	}

	db, err := NewDB(database)

	/*
		author, err = CreateUser(db, 0)
		if err != nil {
			log.Fatal(err)
		}

		x1 := time.Now()
		users, err := CreateMultiple(db, 1000, CreateUser)
		if err != nil {
			log.Fatal(fmt.Errorf("create multiple users failed: %w", err))
		}
		x2 := time.Now()
		trainings, err := CreateMultiple(db, 1000, CreateTraining)
		if err != nil {
			log.Fatal(fmt.Errorf("create multiple trainings failed: %w", err))
		}
		x3 := time.Now()
		for _, v := range trainings {
			err = db.ConnectUserTraining(users[rand.Intn(1000)], v)
			if err != nil {
				log.Fatal(fmt.Errorf("connect multiple users to trainings: %w", err))
			}
		}
		x4 := time.Now()
		fmt.Printf("Create 1000 users: %d ms\n", x2.Sub(x1).Milliseconds())
		fmt.Printf("Create 1000 trainings: %d ms\n", x3.Sub(x2).Milliseconds())
		fmt.Printf("Create 1000 organisers: %d ms\n", x4.Sub(x3).Milliseconds())
	*/
	return db
}
