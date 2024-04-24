package storage

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://localhost:27017"
const collName = "quotes"

var client *mongo.Client
var coll *mongo.Collection

type Quote struct {
	Text              string   `json:"text"`
	Link              string   `json:"link"`
	Author            string   `json:"author"`
	Hashtags          []string `json:"hashtags"`
	HashtagsLowercase []string `json:"hashtags_lowercase"`
}

func NewQuote(text string, link string, author string, hashtags []string) Quote {
	var hashtagsLowercase []string
	for _, hashtag := range hashtags {
		hashtagsLowercase = append(hashtagsLowercase, strings.ToLower(hashtag))
	}
	return Quote{
		Text:              text,
		Link:              link,
		Author:            author,
		Hashtags:          hashtags,
		HashtagsLowercase: hashtagsLowercase,
	}
}

func AddQuote(quote Quote) {
	newQuote := bson.M{
		"text":               quote.Text,
		"link":               quote.Link,
		"author":             quote.Author,
		"hashtags":           quote.Hashtags,
		"hashtags_lowercase": quote.HashtagsLowercase,
	}
	res, err := coll.InsertOne(context.TODO(), newQuote)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", res)
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
	list, _ = client.Database("db").ListCollectionNames(context.TODO(), bson.D{{}})

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
				"Author": bson.M{
					"bsonType":    "string",
					"description": "author",
				},
				"Hashtags": bson.M{
					"bsonType": "array",
					"items": bson.M{
						"bsonType":    "string",
						"description": "associated hashtags",
					},
				},
				"HashtagsLowercase": bson.M{
					"bsonType": "array",
					"items": bson.M{
						"bsonType":    "string",
						"description": "associated hashtags in lowercase",
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
		docs := []Quote{
			NewQuote(
				"I violently severed the continuation of my thoughts\non the question of how and from where can the total fulfillment of this existence be fathomed within this alien form that is fate created by nature.",
				"https://ko-murobushi.com/eng/biblio_selves/p6839/",
				"Ko Murobushi",
				[]string{"KoMurobushi", "Butoh"},
			),
			NewQuote(
				"Standing on that severance,\nfrom the layering darkness the blood drips down, life bubbling up.\nI want to see it, I want to possess it.\nInclined life, labored breath, an inclined ship,\na single-winged flight, an omen of falling from the sky that is an omen of drowning,\nform decaying, the hidden coming to light.",
				"https://ko-murobushi.com/eng/biblio_selves/p6839/",
				"Ko Murobushi",
				[]string{"KoMurobushi", "Butoh"},
			),
			NewQuote(
				"Towards the edge, the corner, the margin, the frontier, the border of being;\nthe stench of death makes life stand out.\nThe thing that stands out desperately—\nthe spasm, cruelly and humorously seared onto the severance as an unmovable fate, is the origin of Butoh.\nIt is the original form of life.",
				"https://ko-murobushi.com/eng/biblio_selves/p6839/",
				"Ko Murobushi",
				[]string{"KoMurobushi", "Butoh"},
			),
			NewQuote(
				"I am a catastrophe at the center of every laceration, every crisis.\nWithout knowing to which world I am being arrested to,\nI, myself, am the den of the demon, arresting the world.\nI must be my medicine man,\na beaten corpse.",
				"https://ko-murobushi.com/eng/biblio_selves/p6839/",
				"Ko Murobushi",
				[]string{"KoMurobushi", "Butoh"},
			),
			NewQuote(
				"I want to be the thing that robs my existence of its totality.\nThe urgency towards the unknown is Eros beyond limits, a love.\nI stop being a god and become a king of madness, unifying myself with the immeasurable.\nIt is not fear nor anxiety that fills my existence or my non-existence.\nThe chaos of the primitive sea in my center,\nthe faraway land from which this sense of existence washes up;\nin front of the ocean of reality, I am a trembling bridge.",
				"https://ko-murobushi.com/eng/biblio_selves/p6839/",
				"Ko Murobushi",
				[]string{"KoMurobushi", "Butoh"},
			),
			NewQuote(
				"I am a nikutai (body), suspended in an escape.\nThe nikutai, masochistically filled with a sense of pain, is conjoined with the darkness of this country.\nThe urgency towards tension, spasms, pulverization, towards the absolute silence……\nThis pain and fascination is unbearable, leaving me lost.\nThe thing that lies on the dumbfounded boundary of the non-self is also the origin of Butoh.",
				"https://ko-murobushi.com/eng/biblio_selves/p6839/",
				"Ko Murobushi",
				[]string{"KoMurobushi", "Butoh"},
			),
			NewQuote(
				"I kill myself in the Butoh space.\nBut I am immortal; I will be ressurected.\nWhen the power that kills me as a whole and revives me as a whole appears from the hina (darkness), rushing through with whirling winds,\nthe hinagata (shape of constant darkness) will be harvested as awakened steel.",
				"https://ko-murobushi.com/eng/biblio_selves/p6839/",
				"Ko Murobushi",
				[]string{"KoMurobushi", "Butoh"},
			),
			NewQuote(
				"To be a ghost is having a conversation with your destination. To be a ghost is having a conversation with the air. (There is a village and there is no sound. You are standing at its gates.)",
				"https://butoh-kaden.com/en/worlds/abyss/ghosts/",
				"Ko Murobushi",
				[]string{"ButohKaden", "Butoh"},
			),
			NewQuote(
				"Ghosts are always transforming into other things at tremendous speeds. Ghosts sometimes imitate living people. Ghosts are also that ephemeral substance that melts into the surroundings. This ghost, unlike a person, has the ability to sense a thousand branches of a tree at the same time. And the ghost, unlike a person, can hear the sounds of these branches grow at the same time. The ghost does not have the form of a person.",
				"https://butoh-kaden.com/en/worlds/abyss/ghosts/",
				"Ko Murobushi",
				[]string{"ButohKaden", "Butoh"},
			),
			NewQuote(
				"The ghost dwells in a place without time and space where numerous white flowers are blooming. Or maybe the ghost hides behind trees and rocks in a Japanese garden. The ghost misses the time and space wher it once lived. Sometimes on the very fingertips, he remenbers the time when he was alive. The ghost is like the mist, the fog, always changing.",
				"https://butoh-kaden.com/en/worlds/abyss/ghosts/",
				"Ko Murobushi",
				[]string{"ButohKaden", "Butoh"},
			),
			NewQuote(
				"A stuffed bird is flying in the wind. (It is almost like the bird which used to fly when it was alive.) It is aware of wires that run inside the body. its cloudy eyes, and feathers which tremble at the slightest wind. There is cotton inside you instead of internal organs. Your feathers have become slightly yellow. You have a desire to fly. You are aware that young chicks found in a nest underneath a rock want flight more than birds which are already flying. You are also aware of the time when you used to fly, and the space in which you used to fly. You are aware of yorr weightlessness. you are also aware of the dry parts of your body after you have become a stuffed bird.",
				"https://butoh-kaden.com/en/worlds/abyss/stuffed-birds/",
				"Ko Murobushi",
				[]string{"ButohKaden", "Butoh"},
			),
			NewQuote(
				"You are also aware of the tiny space of contact with the pedestal. There may be bugs eating into you. The pedestal may slip out of position. The down on your chest may cover your face in the wind.",
				"https://butoh-kaden.com/en/worlds/abyss/stuffed-birds/",
				"Ko Murobushi",
				[]string{"ButohKaden", "Butoh"},
			),
			NewQuote(
				"Speed of these various expressions are important.",
				"https://butoh-kaden.com/en/worlds/abyss/stuffed-birds/",
				"Ko Murobushi",
				[]string{"ButohKaden", "Butoh"},
			),
			NewQuote(
				`10 STATEMENTS (or HOW TO)


				History and Objectives:
				The statements can be used to define a specific area of interest within performance, and to elaborate and develop thoughts on a certain topic. It relies on the form of manifesto where being precise to the point of excluding other possibilities is desirable. The statements do not need to have eternal value, but they should trigger you to think differently. The tool is about producing opinions and positions that can be productive within your work. The purpose of writing 10 statements is to clarify your own ideology and make it visible to others. It’s also about daring to take a stand, exposing yourself to critique and put some fire in the debate.
				
				Description:
				1. Choose a topic that you would like to work on, for instance “statements on how to work”, “statements on site-specific performance”, “statements on spectatorship” or “statements on what practice is”.
				
				2. Think of the format of writing and decide whether or not you want to use a formula. For instance super short and precise, long and descriptive or starting each sentence the same way, x is.../x must.../x is considered...
				
				3. Write the 10 statements on the topic. Try to be as specific as possible and write them in a manner that is coherent with its ideological content and don’t be afraid of being categorical.
				
				Example:
				
				7 statements on how/why 10 statements
				
				Writing ten statements on something is a quick and fun way to help you articulate postitions, strategies or ideologies.
				
				You can choose to write your statements from various positions, also ones you don't nescessarily agree with. Take the more difficult option.
				
				Try to be as precise as possible. This could either mean being super categorical or trying to completely capture whatever complexity exists in what you are writing about and anything inbetween.
				
				Let your statements take the form which is best suited to what you are writing about, the position you are writing from or how you want your statements to work.
				
				Challenge yourself to write all ten statements or more. Don't stop if you first think three is enough.
				
				Writing ten statements on something has the potential to be a versatile tool with lots of uses. It can function as a starting point for discourse, as ideology, as a guide on how to work on something or think around it, or as a departure point from which to develop stuff, like subjects to work on, practices and ways of working.
				
				Write your statements with a clear intention of what or who they are for, or simply begin writing and see what happens. Sometimes the use becomes clear afterwards. Writing ten statements can be good simply to clear your head.
				
				Collective Statement Game
				
				Write one statement on a piece of paper and pass it to your neighbor. He/she can choose to add, change or delete a statement and the paper is finished when it has 10 statements. Continue passing the remaining papers until all papers have 10 statements. Read and discuss the statements in order to understand what they propose in terms of cohesion, fragmentation, perspective etc.
								`,
				"https://web.archive.org/web/20130425011222/http://www.everybodystoolbox.net/?q=node/126",
				"kr66t",
				[]string{"EverybodysToolbox", "Game"},
			),
			NewQuote(
				`BETTER GAME
				History and objectives:
				The objective of this game is to do things in total agreement with a group without getting into group dynamics (no dynamics) or losing individuality.
				Description:
				Anyone in the group can state: IT IS BETTER TO … (Ex. Be in the centre). Everybody agrees and as a concequence does what was said according to one’s own understanding. (As though this is THE new trend) After the activity is established a next statement will be uttered which should not be related to the previous one. This game can go on forever, statements can relate to being, doing, having. If something is not possible or difficult to achieve it is important to genuinely attempt it or else accept to have lost the better game of.`,
				"https://web.archive.org/web/20130425010727/http://www.everybodystoolbox.net/?q=node/68",
				"kr66t",
				[]string{"EverybodysToolbox", "Game"},
			),
		}
		for _, q := range docs {
			AddQuote(q)
		}
	}

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

func SearchQuotes(hashtag string) []Quote {
	var quotes []Quote
	var filter bson.D
	filter = bson.D{{Key: "hashtags_lowercase", Value: strings.ToLower(hashtag)}}
	if hashtag == "" {
		filter = bson.D{{}}
	}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &quotes); err != nil {
		log.Fatal(err)
	}

	if len(quotes) == 0 && hashtag != "" {
		return SearchQuotes("")
	}

	return quotes
}
