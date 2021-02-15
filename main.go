package main

import (
    "fmt"
    "index/suffixarray"
    "regexp"
    "strings"
    "unicode"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

type NewsItem struct{
    title   string
    body    string
}

var documentsIndex = make(map[int]NewsItem)
var titleindex = make(map[string][]int)
var bodyindex = make(map[string][]int)
var sa = suffixarray.New([]byte{})
var joinedStrings = ""

var stopwords = map[string]struct{}{
    "a": {}, "and": {}, "be": {}, "have": {}, "i": {},
    "in": {}, "of": {}, "that": {}, "the": {}, "to": {},
}

func tokenize(text string) []string {
    return strings.FieldsFunc(text, func(r rune) bool {
        return !unicode.IsLetter(r) && !unicode.IsNumber(r)
    })
}

func lowercaseFilter(tokens []string) []string {
    r := make([]string, len(tokens))
    for i, token := range tokens {
        r[i] = strings.ToLower(token)
    }
    return r
}

func stopwordFilter(tokens []string) []string {
    r := make([]string, 0, len(tokens))
    for _, token := range tokens {
        if _, ok := stopwords[token]; !ok {
            r = append(r, token)
        }
    }
    return r
}

func analyze(text string) []string {
    tokens := tokenize(text)
    tokens = lowercaseFilter(tokens)
    tokens = stopwordFilter(tokens)
    return tokens
}

func addDocument(doc NewsItem) {
    keys := make([]int, 0, len(documentsIndex))
    for k := range documentsIndex {
        keys = append(keys, k)
    }
    max := maxSlice(keys)
    index := max+1
    documentsIndex[index] = doc
}

func maxSlice(keys []int) int{
    max := 0
    for _, item := range keys{
        if item > max {
            max = item
        }
    }
    return max
}

func buildIndexes() {
    keys := make([]int, 0, len(documentsIndex))
    for k := range documentsIndex {
        keys = append(keys, k)
    }

    for _ , k := range keys{
        item := documentsIndex[k]
        title := item.title
        titleWords := analyze(title)
        for _,word := range titleWords{
            if val, ok := titleindex[word]; ok {
                titleindex[word] = append(val, k)
            } else {
                titleindex[word] = []int{k}
            }       
        }
    }

    for _ , k := range keys{
        item := documentsIndex[k]
        body := item.body
        bodyWords := analyze(body)
        for _,word := range bodyWords{
            if val, ok := bodyindex[word]; ok {
                bodyindex[word] = append(val, k)
            } else {
                bodyindex[word] = []int{k}
            }       
        }
    }
}

func buildSuffixArray(){
    var words []string
    for key, _ := range titleindex{
        words = append(words, key)
    }
    for key, _ := range bodyindex{
        words = append(words, key)
    }
    fmt.Println(words)
    joinedStrings = "\x00" + strings.Join(words, "\x00")
    sa = suffixarray.New([]byte(joinedStrings))
}

func home(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "GET":
        query := r.URL.Query().Get("query")

        // User has typed in "he"
        match, err := regexp.Compile("\x00"+query+"[^\x00]*")
        if err != nil {
            panic(err)
        }
        ms := sa.FindAllIndex(match, -1)

        var matchedWords []string
        for _, m := range ms {
            start, end := m[0], m[1]
            matchWord := joinedStrings[start+1:end]
            matchedWords = append(matchedWords, matchWord)
        }

        fmt.Println(matchedWords)

        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": `+query+`}`))
    default:
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"message": "invalid"}`))
    }
}


//Can be loaded from CSV or Database
func buildDocumentIndex() {
    n1 := NewsItem{
        title: "Covid Cases Surge in Maharashtra Again, State Records Over 4,000 Cases in 24 Hrs, Mumbai More than 600",
        body: "The state last recorded 4,000-plus cases (4,382) on January 6 and the city recorded (607) daily cases on January 14, exactly a month ago.",
    }

    n2 := NewsItem{
        title: "Activist Arrested for Greta man Thunberg 'Toolkit' Was Working With Pro-Khalistani Group: Delhi Police",
        body: "According to officials, Ravi is 21 years old and lives in Bengaluru. She was active in allegedly disseminating the toolkit, which Thunberg had referred to in her post for the farmers and attached a Google document with details.",
    }

    n3 := NewsItem{
        title: "PM Modi's Photo, Bhagwad Gita & Names of toolkit 25,000 Citizens: Pvt Satellite to be Launched by Feb-End",
        body: "The nanosatellite according has been many developed by SpaceKidz India, an organisation dedicated to promoting space science among students.",
    }

    addDocument(n1)
    addDocument(n2)
    addDocument(n3)
}

func main() {
    buildDocumentIndex()
    buildIndexes()
    buildSuffixArray()
    r := mux.NewRouter()
    r.HandleFunc("/", home)
    log.Fatal(http.ListenAndServe(":8080", r))    
}