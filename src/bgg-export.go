/*
Fetch game collections from boardgamegeek.com for one or more users using the
boardgamegeek.com public XML API. Retreive the games owned by each user and
combine the lists, keeping a count of the duplicates.  Output the result to
stdout as a json file.

Usage:

	bgg-export <user> [<user> ...]
*/
package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type Reply struct {
	Message string `xml:",chardata"`
}

type Errors struct {
	Messages []string `xml:"error>message"`
}

type ItemStatus struct {
	Own          int    `xml:"own,attr"`
	PrevOwned    int    `xml:"prevowned,attr"`
	ForTrade     int    `xml:"fortrade,attr"`
	Want         int    `xml:"want,attr"`
	WantToPlay   int    `xml:"wanttoplay,attr"`
	WantToBuy    int    `xml:"wanttobuy,attr"`
	Wishlist     int    `xml:"wishlist,attr"`
	Preordered   int    `xml:"preordered,attr"`
	LastModified string `xml:"lastmodified,attr"`
}

type CollectionItem struct {
	Id            int        `xml:"objectid,attr"`
	Subtype       string     `xml:"subtype,attr"`
	CollectionId  string     `xml:"collid,attr"`
	Name          string     `xml:"name"`
	YearPublished int        `xml:"yearpublished"`
	Image         string     `xml:"image"`
	Thumbnail     string     `xml:"thumbnail"`
	NumPlays      int        `xml:"numplays"`
	Status        ItemStatus `xml:"status"`
	Comments      string     `xml:"comment"`
	ConditionText string     `xml:"conditiontext"`
}

type Collection struct {
	Owner string
	Items []CollectionItem `xml:"item"`
}

type Game struct {
	Name          string
	YearPublished int
	Url           string
	Copies        int
}

func fetchCollection(user string) (*Collection, error) {
	url := fmt.Sprintf("https://boardgamegeek.com/xmlapi/collection/%s?own=1", user)
	for attempts := 0; attempts < 10; attempts++ {
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("Get URL '%s' : %v", url, err)
		}

		if resp.StatusCode != 200 {
			time.Sleep(5 * time.Second)
			continue /* Server may be busy or trottling, retry. */
		}

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed to read body of url '%s': %v", url, err)
		}

		errors := Errors{}
		err = xml.Unmarshal(bytes, &errors)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal '%s': %v", url, err)
		}
		if len(errors.Messages) > 0 {
			messages := strings.Join(errors.Messages, ", ")
			return nil, fmt.Errorf("Errors from url '%s': %s", url, messages)
		}

		reply := Reply{}
		err = xml.Unmarshal(bytes, &reply)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal '%s': %v", url, err)
		}
		message := strings.TrimSpace(reply.Message)
		if strings.HasPrefix(message, "Your request for this collection has been accepted") {
			time.Sleep(5 * time.Second)
			continue
		}

		collection := Collection{Owner: user}
		err = xml.Unmarshal(bytes, &collection)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal '%s': %v", url, err)
		}
		return &collection, nil
	}
	return nil, fmt.Errorf("Retries exceeded")
}

type fetchResults struct {
	collection *Collection
	err        error
}

func asyncFetchCollection(user string, ch chan<- fetchResults) {
	collection, err := fetchCollection(user)
	results := fetchResults{collection, err}
	ch <- results
}

func main() {
	users := os.Args[1:]

	if len(users) == 0 {
		fmt.Println("usage: bgg-export <user> [<user> ...]")
		os.Exit(1)
	}

	ch := make(chan fetchResults, len(users))
	for _, user := range users {
		go asyncFetchCollection(user, ch)
	}

	/*
	 * Gather collections in a map of games, keeping a count of the number of
	 * duplicates found.
	 */
	games := make(map[int]*Game)
	for range users {
		results := <-ch
		if results.err != nil {
			fmt.Println("error", results.err)
			os.Exit(1)
		}
		collection := results.collection
		for _, item := range collection.Items {
			game, ok := games[item.Id]
			if !ok {
				url := fmt.Sprintf("https://boardgamegeek.com/boardgame/%d/", item.Id)
				games[item.Id] = &Game{item.Name, item.YearPublished, url, 1}
			} else {
				game.Copies++
			}
		}
	}

	/* Convert the map of games to a slice for output. */
	gameList := []*Game{}
	for _, game := range games {
		gameList = append(gameList, game)
	}
	sort.Slice(gameList, func(a, b int) bool {
		if strings.ToLower(gameList[a].Name) == strings.ToLower(gameList[b].Name) {
		    return gameList[a].YearPublished < gameList[b].YearPublished
        }
		return strings.ToLower(gameList[a].Name) < strings.ToLower(gameList[b].Name)
	})

	/* Print our game list as json. */
	output, err := json.MarshalIndent(gameList, "", "  ")
	if err != nil {
		fmt.Println("json marshal error:", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}
