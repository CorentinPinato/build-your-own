package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

// Data
// id, name, content

func main() {
	documents := map[int]string{
		0:"At Scale You Will Hit Every Performance Issue I used to think I knew a bit about performance scalability and how to keep things trucking when you hit large amounts of data Truth is I know diddly squat on the subject since the most I have ever done is read about how its done To understand how I came about realising this you need some background",
		1:"Richard Stallman to visit Australia Im not usually one to promote events and the like unless I feel there is a genuine benefit to be had by attending but this is one stands out Richard M Stallman the guru of Free Software is coming Down Under to hold a talk You can read about him here Open Source Celebrity to visit Australia",
		2:"MySQL Backups Done Easily One thing that comes up a lot on sites like Stackoverflow and the like is how to backup MySQL databases The first answer is usually use mysqldump This is all fine and good till you start to want to dump multiple databases You can do this all in one like using the all databases option however this makes restoring a single database an issue since you have to parse out the parts you want which can be a pain",
		3:"Why You Shouldnt roll your own CAPTCHA At a TechEd I attended a few years ago I was watching a presentation about Security presented by Rocky Heckman read his blog its quite good In it he was talking about security algorithms The part that really stuck with me went like this",
		4:"The Great Benefit of Test Driven Development Nobody Talks About The feeling of productivity because you are writing lots of code Think about that for a moment Ask any developer who wants to develop why they became a developer One of the first things that comes up is I enjoy writing code This is one of the things that I personally enjoy doing Writing code any code especially when its solving my current problem makes me feel productive It makes me feel like Im getting somewhere Its empowering",
		5:"Setting up GIT to use a Subversion SVN style workflow Moving from Subversion SVN to GIT can be a little confusing at first I think the biggest thing I noticed was that GIT doesnt have a specific workflow you have to pick your own Personally I wanted to stick to my Subversion like work-flow with a central server which all my machines would pull and push too Since it took a while to set up I thought I would throw up a blog post on how to do it",
		6:"Why CAPTCHA Never Use Numbers 0 1 5 7 Interestingly this sort of question pops up a lot in my referring search term stats Why CAPTCHAs never use the numbers 0 1 5 7 Its a relativity simple question with a reasonably simple answer Its because each of the above numbers are easy to confuse with a letter See the below",
	}

	index := map[int]Concordance{
		0:concordance(transform(documents[0])),
		1:concordance(transform(documents[1])),
		2:concordance(transform(documents[2])),
		3:concordance(transform(documents[3])),
		4:concordance(transform(documents[4])),
		5:concordance(transform(documents[5])),
		6:concordance(transform(documents[6])),
	}


	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter search term: ")
	searchTerm, _ := reader.ReadString('\n')
	searchTerm = searchTerm[:len(searchTerm) - 1]

	matches := make([]Match, 0)

	con := concordance(transform(searchTerm))
	for docID,idx := range index {
		relation := relation(con, idx)
		if relation != 0 {
			matches = append(matches, Match{relation, documents[docID][:100]})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].score > matches[j].score
	})

	// Filter out lower scores
	//
	// temp := matches[:0]
	// for _,match := range matches {
	// 	if match.score >= 0.2 {
	// 		temp = append(temp, match)
	// 	}
	// }
	// matches = temp

	for _,match := range matches {
		fmt.Printf("Score: %f Doc: %s\n", match.score, match.document)
	}
}

type Concordance map[string]int
type Match struct {
	score float64
	document string
}

func concordance(document string) map[string]int {
	con := map[string]int{}

	for _,word := range strings.Split(document, " ") {
		_,ok := con[word]
		if ok {
			con[word] += 1
		} else {
			con[word] = 1
		}
	}

	return con
}

func magnitude(concordance map[string]int) float64 {
	var total float64 = 0

	for _,count := range concordance {
		total += math.Pow(float64(count), 2)
	}
	
	return math.Sqrt(total)
}

func relation(concordance1, concordance2 Concordance) float64 {
	var topValue int = 0

	for word,count := range concordance1 {
		_,ok := concordance2[word]
		if ok {
			topValue += count * concordance2[word]
		}
	}

	mag := magnitude(concordance1) * magnitude(concordance2)
	if mag != 0 {
		return float64(topValue) / mag
	} else {
		return 0
	}
}

// qwerty => qwe, wer, ert, rty
//
func trigram(term string) string {
	size := len(term)
	result := ""

	if size <= 3 { return " " + term }

	for i := 0; i + 3 <= size; i += 1 {
		result += " " + term[i:i + 3]
	}

	return result
}

func transform(document string) string {
	document = strings.ToLower(document)
	result := ""

	for _,term := range strings.Split(document, " ") {
		result += trigram(term)
	}

	return result
}

func printConcordance(con Concordance) {
	for k,v := range con {
		fmt.Printf("%s => %d\n", k, v)
	}
}
