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
			NewQuote("What do you start from when initiating a project?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you work with the other people involved in the project?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you use the time of the creation?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What are the phases in the creation of a performance?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you consider the spectators’ experience when preparing a performance?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Which skills do you use?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Which conventions do you use?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Which modes of representation do you use?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you position yourself in relation to the object/topic/discourse of your work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("When does the discourse of your work become clear to you?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Which habits do you recognise in your way of working?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Which ones would you like to change?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What kind of theoretical concepts are you interested in for the future and how do they differ from what you have been busy with in the past?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What methods do you use to test your ideas/works?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you involve other people in your work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What do you think about the idea of a performance user/ test group?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Does this differ from a conventional showing/work in progress?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How Is your work fixed/finalized once it has become a performance?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What are your most recent strategies of working?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What is it that makes you work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What are the problems that interest you?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you decide on possible collaborations?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("When does fi nancial production come into play?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How does writing a proposal in the form of an application infl uence your work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Does it help you to formulate or do you feel restrained by the institutional frames produced by governments, theater, research and work spaces.?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you feel you can work outside of structures of support?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What does community mean to you, are you part of one/several and how does this infl uence or support you work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How would you define the nature of sharing you are taking part in?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("When working on a performance, how do you start?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What kind of preparation do you do, or rather when and how does an idea pass into actuality?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What role do other media play in your work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you relate theory to practice?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Does the one have supremacy over the other?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("When and how do you decide of the frame to expose your work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you articulate the working process?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you relate an idea, a sensation, a material, to the process of work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you choose the material?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What role does 'spectacularity' play in your work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you deal with problems that come up?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What is your reference frame when you meet problems?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What do you think is significant?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Which role does sensation/ intuition play?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("On which level would you place expression?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you expose your work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What is aesthetics for you?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What subjects are you really interested in?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Why is it interesting for you to occupy yourself with that?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Why is it interesting for you to work with this medium?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What have your performancesto do with you?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you have an opinion or point of view on the subject?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What is your attitude towards your medium?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Who looks over your shoulder when you are at work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What is your specific input in your performance?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Are you provocative in your work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Are you making a comment?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Are you expressing feelings?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What is your position as an artist in your work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you think people can learn something from you?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Why are you not a teacher?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you like to work with symbols and metaphors?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Are you interested in psychoanalysis?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What makes you want to use it in your work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Are you interested in politics and economics?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Why have you not become a journalist?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What is the relationship between your political interests and your interest in theatre?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is there a project that you have been eager to realize for a long time now?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Where do you get your ideas?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you judge the viability of your ideas?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is your idea one that you would want to spend a year or more on?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What does this idea have to do with you?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How does the idea become workable?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you start a new project – from a tabula rasa?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you create a space for thinking within the work process?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you organise the elements of your work, your ideas, your media, your criteria, the technology?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Imagine you are a painter: what is your attitude towards the canvas?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Imagine you are a painter: what do and don’t you want to have on the canvas?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How can you realise what is in your head?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you construct a performance?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you choose one solution or do you explore several possibilities?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What part does the audience play in your performance?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you want something to happen between the audience and the performance?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What do you expect of the audience?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you feel about communicating with your audience?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What do you want to communicateto your audience?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you want to give something to the audience?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you think it important to know what the audience thinks about you?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("And about your work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you think it important that the audience believes you?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Does the audience have to exert itself with your work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you want them to?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is it important that the audience discovers the link between the performance and your ideas that lie behind it?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you think your inner world is interesting for the audience?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What must happen to the audience during your piece?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you manipulate the audience?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you want to share your emotions with the audience?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you want the audience to have an emotional experience?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you want the audience to start thinking about your subject or take a position on it?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Why?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is to “realize” something for the audience the same as to “experience” something?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you set out to let the audience identify with the people in your performance, or with your problems?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What reaction from the audience would displease you?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is the audience informed about your subject before they enter the performance space?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What is the aim of this piece?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you propose to achieve that aim?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("When openness is your aim, how do you measure whether your piece is open enough?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How many objectives can you achieve with one piece?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Can you achieve your objectives with the media you have chosen to employ?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you use material that is close to you or that which is distant?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Are you reconstructing a situation in order to underline something?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Are you looking for a framework?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you want to limit or define your ideas?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is your work open or closed?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is your performance a statement?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("At the level of content or form or both?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What was your last performance about?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is there a connection between your last piece and this new one?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What is your next project about?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How does „content“ become theatre?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do content and form coincide in your work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Can one separate content and form?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Can one separate the „what“ and the „how“?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What comes first, content or form?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Can form be presented as content?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you want it to be?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Can content direct you towards form?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is there a connection or a dissonance between the form and content of what you create?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you think that aesthetics are a form of manipulation?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Have you ever considered leaving out any and all suggestion, illusion and reference?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Have you ever considered including any and all suggestion, illusion and reference?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What is specifi c about theatre for you?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How specifically is the medium of theatre used in your piece?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What do you mean by theatre and theatricality?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you accept certain principles specific to theatre such as begginning, end, time-line, and the shared space of performers and audience?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is theatre the right medium, as far as credibility and identifi cation are concerned?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Where do you think illusion works best, in theatre or in film?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is theatre a medium suitable for educating people?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Are you eager to do something new in the theatre?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is this the only way for you to make theatre or are there others?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("According to you, what is dance?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you translate a subject into a dance performance?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Does your choreography say what you want to say?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Why do you choose this art form to say what you have to say?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you use the medium in a way that is not possible with another medium?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How does the idea become workable?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you start a new project – from a tabula rasa?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Do you like to work with symbols and metaphors?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Are you interested in psychoanalysis?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What makes you want to use it in your work?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Are you interested in politics and economics?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Why have you not become a journalist?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What is the relationship between your political interests and your interest in theatre?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is there a project that you have been eager to realize for a long time now?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Where do you get your ideas?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you judge the viability of your ideas?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is your idea one that you would want to spend a year or more on?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What does this idea have to do with you?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("How do you arrive at your selection of artistic media, such as video and music?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What is the performance about?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What statement do you want to make with this performance?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Where lies the essence of the piece for you?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is it a conventional or unconventional piece?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is it a dance project, a theatre project or a visual arts project?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Why?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is there a political aspect to the piece?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Is there a story behind the piece?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Why do you wish to make this?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("And why exactly in this way?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("Why is this performance interesting to watch, to listen to, or experience?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What are the other questions?", "https://web.archive.org/web/20110611064738/http://www.everybodystoolbox.net/?q=node/14", "alice chauchat", []string{"EverybodysToolbox", "Question"}),
			NewQuote("What did you long for?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you have any sides or private or jugmental thoughts going on? When and what did you think of?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you have a private score / focus you were working on? Which one(s)?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Do you think it was a successful improvisation? Why?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("What was your favorite moment?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you see something which particularly irritated you? What?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did someone do something to you which particularly irritated you? What?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you experience something which absolutely did not work for you? What?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("What were your strategies to go through the improvisation when you felt stuck?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Which elements were you composing with?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Which compositional tools did you use?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Can you give us some examples of elements/situations/persons who moved/inspired you when you were inside the improvisation?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Can you give us some examples of elements/situations/persons who moved/inspired you when you were outside the improvisation?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("When you were out, what motivated you to go in?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("When you were out, what motivated you to stay out?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("When you were in, what motivated you to stay in?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("When you were in, what motivated you to go out?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did something/someone bring you out of your performance/improvisational state? What/who?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you have someone whom you really enjoyed watching? Who?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you enjoy how the group worked together? In what way?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Do you think the group has an intelligengie of it's own? ", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Could you see beauty in the group's way of composing? What kind?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you feel supported by the group? In what way?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did we just create an art piece together?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you sometimes not know what you were doing? Please describe…", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you sometimes not know what to do? What did you do then?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you get bored of others? When? What did you do then?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you get bored of yourself? When? What did you do then?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Do you like improvising for others? ", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Do you like improvising with no one watching from the outside?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Do you enjoy watching people improvise for you?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you have to be courageous at any point? When?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you have to be patient at any point? When?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you enjoy when people were using their voice / talking / singing? Why?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you feel like breaking the format? When?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you break it? When?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Were you judging what were you doing? What kind of judjements did you do?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Were you judging what the others where doing? What kind of judjements did you do?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you think about the future? When?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you think about the past? When?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you miss anything / anyone? When?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did you mostly have fun?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("Did the warm up help you?", "http://www.idocde.net/idocs/333", "Sandra Wieser", []string{"IDOCDE", "Question", "Improvisation", "Composition"}),
			NewQuote("What is your professional formation? Your technical, artistic and pedagogical formation? How does this inform your teaching? How do you “pay it forward”?", "http://www.idocde.net/idocs/797", "Defne Erdur", []string{"IDOCDE", "Question", "Teaching"}),
			NewQuote("What is your personal interest in teaching? Why do you teach?", "http://www.idocde.net/idocs/797", "Defne Erdur", []string{"IDOCDE", "Question", "Teaching"}),
			NewQuote("What do you transmit?", "http://www.idocde.net/idocs/797", "Defne Erdur", []string{"IDOCDE", "Question", "Teaching"}),
			NewQuote("What  do you believe that Contemporary Dance Training needs today?", "http://www.idocde.net/idocs/797", "Defne Erdur", []string{"IDOCDE", "Question", "Teaching"}),
			NewQuote("What is the aim/goal of your teaching?", "http://www.idocde.net/idocs/797", "Defne Erdur", []string{"IDOCDE", "Question", "Teaching"}),
			NewQuote("What are the tools and principles do you utilize? How do you teach what you teach?", "http://www.idocde.net/idocs/797", "Defne Erdur", []string{"IDOCDE", "Question", "Teaching"}),
			NewQuote("Where does your teaching content and format fit into? Which context  is better for you – institutional, festivals, workshop, master class, collaboration with other artist?", "http://www.idocde.net/idocs/797", "Defne Erdur", []string{"IDOCDE", "Question", "Teaching"}),
			NewQuote("Who do you teach? Who do you prefer to work with – Professionals, non-professionals; children, teenagers, adults, elders; dis-advantaged populations?", "http://www.idocde.net/idocs/797", "Defne Erdur", []string{"IDOCDE", "Question", "Teaching"}),
			NewQuote("Do you adjust your teaching as a whole for different populations? If so why and how?", "http://www.idocde.net/idocs/797", "Defne Erdur", []string{"IDOCDE", "Question", "Teaching"}),
			NewQuote("What is “Art Education” for you?", "http://www.idocde.net/idocs/797", "Defne Erdur", []string{"IDOCDE", "Question", "Teaching"}),
			NewQuote("What is the place of Creativity, Personal Development, Inter-disciplinarity in your teaching?", "http://www.idocde.net/idocs/797", "Defne Erdur", []string{"IDOCDE", "Question", "Teaching"}),
			NewQuote("Where do you place yourself in the Artist – Teacher spectrum?", "http://www.idocde.net/idocs/797", "Defne Erdur", []string{"IDOCDE", "Question", "Teaching"}),
			NewQuote("What are the pro’s and con’s of being an active Artist for being a Teacher?", "http://www.idocde.net/idocs/797", "Defne Erdur", []string{"IDOCDE", "Question", "Teaching"}),
			NewQuote("", "http://www.idocde.net/idocs/797", "Defne Erdur", []string{"IDOCDE", "Question", "Teaching"}),
			NewQuote("In relation to what you planned to do (your concept) where did you end up?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("Did anything change along the way? What? And how did you adapt to this?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("What was the most challenging part of the process?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("What was the most surprising? The most inspiring? The most confronting?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("Were there any pivotal moments in the process?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("What were your difficulties? What did you need/do to overcome them?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("What were your victories? How did they come about?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("Were there thing(s) for which you didn’t find a solution/answer?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("How was the process of choosing and integrating the technical aspects of your piece (light/sound/the theater space)?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("How was the experience of performing your work for an audience?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("Did the dialogue with your public provide any new/additional information about what you were aiming to accomplish or communicate?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("What kind of feedback was helpful to hear from the public?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("Along the way, what did you discover? About making dances? About yourself and your own way of being/working? Or about dance/theater as a medium?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("Are there things you’d like to continue investigating/ working on as a result?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("Are there things you’d consider doing differently next time you make own work?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("What do you take away with you from this process?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("How will you work with these realizations in the future?", "http://www.idocde.net/idocs/368", "Hillary Blake Firestone", []string{"IDOCDE", "Question", "Teaching", "Reflection"}),
			NewQuote("To what extent do visuals play an effective role in triggering movement improvisation?", "http://www.idocde.net/idocs/504", "aylin kalem", []string{"IDOCDE", "Question", "Improvisation", "Image"}),
			NewQuote("What is the relationship between movement and visual perception?", "http://www.idocde.net/idocs/504", "aylin kalem", []string{"IDOCDE", "Question", "Improvisation", "Image"}),
			NewQuote("How does a person get into interaction with visual images?", "http://www.idocde.net/idocs/504", "aylin kalem", []string{"IDOCDE", "Question", "Improvisation", "Image"}),
			NewQuote("How to generate a choreographic perspective in relation to images?", "http://www.idocde.net/idocs/504", "aylin kalem", []string{"IDOCDE", "Question", "Improvisation", "Image"}),
			NewQuote("Can images augment the expressive quality of gestures?", "http://www.idocde.net/idocs/504", "aylin kalem", []string{"IDOCDE", "Question", "Improvisation", "Image"}),
			NewQuote("what is my costume?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("what am i interested in today?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("how do i use my voice?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("what kind of language to i use?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("what kind of walls do i put up?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("how porous am i?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("how fearless can i be?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("do i believe in what i am doing?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("do i believe in what i am saying?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("am i faking it?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("am i for real?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("what am i most excited about sharing?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("how can i be surprised?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("how do i listen to the room?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("how can i direct the flow of the class?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("how can i follow the flow of the class?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("how do i entertain?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
			NewQuote("how do i leave space for others to fill?", "http://www.idocde.net/idocs/594", "alicia grant", []string{"IDOCDE", "Question", "Teaching", "Performance"}),
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
