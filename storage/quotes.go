package storage

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://localhost:27017"
const collName = "quotes"

var client *mongo.Client
var coll *mongo.Collection
var quotes []Quote

type Quote struct {
	Text     string   `json:"text"`
	Link     string   `json:"link"`
	Hashtags []string `json:"hashtags"`
}

func NewQuote(text string, link string, hashtags []string) Quote {
	return Quote{Text: text, Link: link, Hashtags: hashtags}
}

func AddQuote(text string, link string, hashtags []string) {
	newQuote := bson.M{"text": text, "link": link, "hashtags": hashtags}
	res, err := coll.InsertOne(context.TODO(), newQuote)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", res)
	reloadQuotes()
}

func init() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	var err error
	client, err = mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	///

	var list []string
	list, err = client.Database("db").ListCollectionNames(context.TODO(), bson.D{{}})

	fmt.Printf("%v\n", list)

	collFound := false
	for _, name := range list {
		if name == collName {
			collFound = true
		}
	}
	if !collFound {
		fmt.Printf("Creating and initializing the collection %s\n", collName)
		jsonSchema := bson.M{
			"bsonType": "object",
			"required": []string{"text"},
			"properties": bson.M{
				"Text": bson.M{
					"bsonType":    "string",
					"description": "content",
				},
				"Link": bson.M{
					"bsonType":    "string",
					"description": "source url",
				},
				"Hashtags": bson.M{
					"bsonType": "array",
					"items": bson.M{
						"bsonType":    "string",
						"description": "associated hashtags",
					},
				},
			},
		}
		validator := bson.M{
			"$jsonSchema": jsonSchema,
		}
		coll_opts := options.CreateCollection().SetValidator(validator)
		err = client.Database("db").CreateCollection(context.TODO(), collName, coll_opts)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Collection %s found\n", collName)
	}

	coll = client.Database("db").Collection(collName)

	if !collFound {
		docs := []interface{}{
			Quote{
				Text:     "I violently severed the continuation of my thoughts\non the question of how and from where can the total fulfillment of this existence be fathomed within this alien form that is fate created by nature.",
				Link:     "https://ko-murobushi.com/eng/biblio_selves/p6839/",
				Hashtags: []string{"KoMurobushi", "butoh"},
			},
			Quote{
				Text:     "Standing on that severance,\nfrom the layering darkness the blood drips down, life bubbling up.\nI want to see it, I want to possess it.\nInclined life, labored breath, an inclined ship,\na single-winged flight, an omen of falling from the sky that is an omen of drowning,\nform decaying, the hidden coming to light.",
				Link:     "https://ko-murobushi.com/eng/biblio_selves/p6839/",
				Hashtags: []string{"KoMurobushi", "Butoh"},
			},
			Quote{
				Text:     "Towards the edge, the corner, the margin, the frontier, the border of being;\nthe stench of death makes life stand out.\nThe thing that stands out desperately—\nthe spasm, cruelly and humorously seared onto the severance as an unmovable fate, is the origin of Butoh.\nIt is the original form of life.",
				Link:     "https://ko-murobushi.com/eng/biblio_selves/p6839/",
				Hashtags: []string{"KoMurobushi", "Butoh"},
			},
			Quote{
				Text:     "I am a catastrophe at the center of every laceration, every crisis.\nWithout knowing to which world I am being arrested to,\nI, myself, am the den of the demon, arresting the world.\nI must be my medicine man,\na beaten corpse.",
				Link:     "https://ko-murobushi.com/eng/biblio_selves/p6839/",
				Hashtags: []string{"KoMurobushi", "Butoh"},
			},
			Quote{
				Text:     "I want to be the thing that robs my existence of its totality.\nThe urgency towards the unknown is Eros beyond limits, a love.\nI stop being a god and become a king of madness, unifying myself with the immeasurable.\nIt is not fear nor anxiety that fills my existence or my non-existence.\nThe chaos of the primitive sea in my center,\nthe faraway land from which this sense of existence washes up;\nin front of the ocean of reality, I am a trembling bridge.",
				Link:     "https://ko-murobushi.com/eng/biblio_selves/p6839/",
				Hashtags: []string{"KoMurobushi", "Butoh"},
			},
			Quote{
				Text:     "I am a nikutai (body), suspended in an escape.\nThe nikutai, masochistically filled with a sense of pain, is conjoined with the darkness of this country.\nThe urgency towards tension, spasms, pulverization, towards the absolute silence……\nThis pain and fascination is unbearable, leaving me lost.\nThe thing that lies on the dumbfounded boundary of the non-self is also the origin of Butoh.",
				Link:     "https://ko-murobushi.com/eng/biblio_selves/p6839/",
				Hashtags: []string{"KoMurobushi", "Butoh"},
			},
			Quote{
				Text:     "I kill myself in the Butoh space.\nBut I am immortal; I will be ressurected.\nWhen the power that kills me as a whole and revives me as a whole appears from the hina (darkness), rushing through with whirling winds,\nthe hinagata (shape of constant darkness) will be harvested as awakened steel.",
				Link:     "https://ko-murobushi.com/eng/biblio_selves/p6839/",
				Hashtags: []string{"KoMurobushi", "Butoh"},
			},
			Quote{
				Text:     "To be a ghost is having a conversation with your destination. To be a ghost is having a conversation with the air. (There is a village and there is no sound. You are standing at its gates.)",
				Link:     "https://butoh-kaden.com/en/worlds/abyss/ghosts/",
				Hashtags: []string{"ButohKaden", "Butoh"},
			},
			Quote{
				Text:     "Ghosts are always transforming into other things at tremendous speeds. Ghosts sometimes imitate living people. Ghosts are also that ephemeral substance that melts into the surroundings. This ghost, unlike a person, has the ability to sense a thousand branches of a tree at the same time. And the ghost, unlike a person, can hear the sounds of these branches grow at the same time. The ghost does not have the form of a person.",
				Link:     "https://butoh-kaden.com/en/worlds/abyss/ghosts/",
				Hashtags: []string{"ButohKaden", "Butoh"},
			},
			Quote{
				Text:     "The ghost dwells in a place without time and space where numerous white flowers are blooming. Or maybe the ghost hides behind trees and rocks in a Japanese garden. The ghost misses the time and space wher it once lived. Sometimes on the very fingertips, he remenbers the time when he was alive. The ghost is like the mist, the fog, always changing.",
				Link:     "https://butoh-kaden.com/en/worlds/abyss/ghosts/",
				Hashtags: []string{"ButohKaden", "Butoh"},
			},
			Quote{
				Text:     "A stuffed bird is flying in the wind. (It is almost like the bird which used to fly when it was alive.) It is aware of wires that run inside the body. its cloudy eyes, and feathers which tremble at the slightest wind. There is cotton inside you instead of internal organs. Your feathers have become slightly yellow. You have a desire to fly. You are aware that young chicks found in a nest underneath a rock want flight more than birds which are already flying. You are also aware of the time when you used to fly, and the space in which you used to fly. You are aware of yorr weightlessness. you are also aware of the dry parts of your body after you have become a stuffed bird.",
				Link:     "https://butoh-kaden.com/en/worlds/abyss/stuffed-birds/",
				Hashtags: []string{"ButohKaden", "Butoh"},
			},
			Quote{
				Text:     "You are also aware of the tiny space of contact with the pedestal. There may be bugs eating into you. The pedestal may slip out of position. The down on your chest may cover your face in the wind.",
				Link:     "https://butoh-kaden.com/en/worlds/abyss/stuffed-birds/",
				Hashtags: []string{"ButohKaden", "Butoh"},
			},
			Quote{
				Text:     "Speed of these various expressions are important.",
				Link:     "https://butoh-kaden.com/en/worlds/abyss/stuffed-birds/",
				Hashtags: []string{"ButohKaden", "Butoh"},
			},
		}

		res, err := coll.InsertMany(context.TODO(), docs)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v", res)
	}

	reloadQuotes()

	// rootPath, _ := os.Getwd()
	// jsonData, err := os.ReadFile(rootPath + "/data.json")
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return
	// }

	// err = json.Unmarshal(jsonData, &quotes)
	// if err != nil {
	// 	fmt.Printf("could not unmarshal json: %s\n", err)
	// 	return
	// }
	// fmt.Printf("json map: %v\n", quotes)
}

func GetQuotes() []Quote {
	return quotes
}

func reloadQuotes() {
	filter := bson.D{{}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	// var results []Quote
	if err = cursor.All(context.TODO(), &quotes); err != nil {
		log.Fatal(err)
	}
}
