package graph

import (
	"fmt"
	"log"
	"math/rand"
	"pkv/api/src/domain"
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
	german := CreateDescription("de", text, false)
	english := CreateDescription("en", text, true)
	return domain.Descriptions{
		"de": german,
		"en": english,
	}
}

func CreateDescription(language string, text string, translated bool) domain.Description {
	switch language {
	case "de":
		return domain.Description{
			text + " - Eine tolle Sache",
			"Das wird euch sicher ganz gut gefallen",
			translated,
		}
		break
	case "en":
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
	err = db.Users.Create(&user, nil)
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
	err = db.Trainings.Create(&training, nil)
	if err != nil {
		return training, err
	}
	location, err := CreateLocation(db, i)
	if err != nil {
		return training, err
	}
	err = db.TrainingHappensAtLocation(&training, &location, nil)
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
	err = db.Locations.Create(&location, nil)
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

func SampleData(db *Db) {
	admin := domain.User{
		Key:  "admin",
		Name: "Admin",
		Type: "administrator",
	}
	if err := db.Users.Create(&admin, nil); err != nil {
		log.Fatal(err)
	}
	dpv := domain.User{
		Key:  "dpv",
		Name: "Deutscher Parkour Verband",
		Type: "association",
		Descriptions: map[string]domain.Description{
			"de": {
				Title: "Deutscher Parkour Verband",
				Text:  "Der Deutsche Parkour Verband e.V. (DPV) ist der Dachverband für Parkour und Freerunning in Deutschland. Er wurde 2024 gegründet und vertritt die Interessen der Parkour- und Freerunning-Szene in Deutschland. Der DPV ist zudem Mitglied im Deutschen Olympischen Sportbund (DOSB).",
			},
			"en": {
				Title: "German Parkour Association",
				Text:  "The German Parkour Association (DPV) is the umbrella organization for parkour and freerunning in Germany. It was founded in 2024 and represents the interests of the parkour and freerunning scene in Germany. The DPV is also a member of the German Olympic Sports Confederation (DOSB).",
			},
		},
		Comments: []domain.Comment{
			{
				"Endlich!",
				"Das hat lange gedauert",
				"admin",
				time.Now(),
			},
		},
	}
	if err := db.Users.Create(&dpv, nil); err != nil {
		log.Fatal(err)
	}
	if err := db.UserAdministersUser(admin, dpv, nil); err != nil {
		log.Fatal(err)
	}
	berlin := domain.Location{
		Key:  "berlin",
		City: "Berlin",
		Lat:  52.52,
		Lng:  13.40,
		Type: "office",
	}
	if err := db.Locations.Create(&berlin, nil); err != nil {
		log.Fatal(err)
	}
	meeting := domain.Training{
		Key:  "berlin-meeting-november-2023",
		Type: "meeting",
		Descriptions: map[string]domain.Description{
			"de": {
				Title: "Berlin Meeting November 2023",
				Text:  "Das Berlin Meeting ist ein monatliches Treffen der Parkour- und Freerunning-Szene in Deutschland. Es findet jeden Monat statt und wird vom DPV organisiert.",
			},
			"en": {
				Title: "Berlin Meeting November 2023",
				Text:  "The Berlin Meeting is a monthly meeting of the parkour and freerunning scene in Germany. It takes place every month and is organized by the DPV.",
			},
		},
	}
	if err := db.Trainings.Create(&meeting, nil); err != nil {
		log.Fatal(err)
	}
	if err := db.TrainingHappensAtLocation(&meeting, &berlin, nil); err != nil {
		log.Fatal(err)
	}
	if err := db.UserOrganisesTraining(dpv, meeting, nil); err != nil {
		log.Fatal(err)
	}
	page := domain.Page{
		Key: "satzung",
		Descriptions: map[string]domain.Description{
			"de": {
				Title: "Satzung",
				Text:  "Die Satzung des Deutschen Parkour Verbandes e.V.",
			},
			"en": {
				Title: "Articles of Association",
				Text:  "The articles of association of the German Parkour Association e.V.",
			},
		},
	}
	if err := db.Pages.Create(&page, nil); err != nil {
		log.Fatal(err)
	}
	if err := db.UserOwnsPage(dpv, page, nil); err != nil {
		log.Fatal(err)
	}
}

func NewTestDB(db *Db) {
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
}
