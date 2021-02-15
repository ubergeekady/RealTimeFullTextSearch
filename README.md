Fast indexing and searching

## Indexing

1. The news titles and body is loaded from database or CSV file (currently hardcoded inside a function)

2. `documents` hashmap contains a list of all the news item structs indexed by numbers 1,2,3

   {1:{title,body}, 2:{title,body}, 3:{title,body}.....}

3. `titleIndex` and `bodyIndex` are two other inverted index hashmaps which map words to documentIDs

   {"word1":[3,4] , "word2":[7,8,9]}

4. A suffix array is created which holds all the words in all the articles



## Querying

1. On GET request, we first lookup the suffix array for matching words. So. "tr" will return "trump", "trust", "travel" as present in the article titles and bodies.
2. First lookup the titleIndex for the word and get the matching documents.
3. 