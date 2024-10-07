package main

import (
    "bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
    "regexp"
)

func getFileNames() []string{
    dir := "./dataset"
    filenames := []string{}
    // Read all files from the directory
    files, err := os.ReadDir(dir)
    if err != nil {
        log.Fatalf("failed to read directory: %s", err)
    }
    
    for _, file := range files {
        if !file.IsDir() {
            // Get the file extension
            if filepath.Ext(file.Name()) == ".txt" {
                // fmt.Println(file.Name())
                filenames = append(filenames, file.Name())
            }
        }
    }
    return filenames
}

func fileReader(filename string) string{
    var builder strings.Builder
    file, err := os.Open(fmt.Sprintf("dataset/%s", filename))
    if err != nil {
        log.Fatalf("failed to open file: %s", err)
    }
    defer file.Close()
	scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        builder.WriteString(scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        log.Fatalf("failed to read file: %s", err)
    }
    return builder.String()
}

func stateGenerator(states *map[string]nextstate, content string) []string{
    tokenList := strings.Split(content, " ")
    cleanTokens := []string{}
    for _, val := range tokenList{
        re := regexp.MustCompile(`[^a-zA-Z]+`)
        val = re.ReplaceAllString(val, "")
        if val != ""{
            cleanTokens = append(cleanTokens, strings.ToLower(val))
        }
    }

    groupedTokens := []string{}

    for i := 0; i < len(cleanTokens)-1;  i = i + 2{
        grouped:= fmt.Sprintf("%s %s", cleanTokens[i], cleanTokens[i+1])
        groupedTokens = append(groupedTokens, grouped)
    }

    for i := 1; i < len(cleanTokens)-1;  i = i + 2{
        grouped:= fmt.Sprintf("%s %s", cleanTokens[i], cleanTokens[i+1])
        groupedTokens = append(groupedTokens, grouped)
    }

    for i := 0; i <len(groupedTokens)-1; i++{
        rightTokens, foundLeft := (*states)[groupedTokens[i]]
        if !foundLeft{
            rightTokens := make(nextstate)
            rightTokens[groupedTokens[i+1]] = 1
            (*states)[groupedTokens[i]] = rightTokens
        }else{
            _, foundRight := rightTokens[groupedTokens[i+1]]
            if !foundRight{
                rightTokens[groupedTokens[i+1]] = 1
            }else{
                rightTokens[groupedTokens[i+1]]++
            }
        }
    }
    return groupedTokens
}

func findMaxKey(myMap map[string]int) (string, int) {
	var maxKey string

	maxValue := -1
	for key, value := range myMap {
		if value > maxValue {
			maxValue = value
			maxKey = key
		}
	}
	return maxKey, maxValue
}

func storyGenerator(states *map[string]nextstate, startingToken string, maxLen int) string{
    var builder strings.Builder
    builder.WriteString(startingToken)
    currToken := startingToken 
    for i := 0; i < maxLen; i++{
        nextToken, val := findMaxKey((*states)[currToken])
        if val == -1{
            // builder.WriteString(".")
            return builder.String()
        }else{
            builder.WriteString(" ")
            builder.WriteString(nextToken)
            currToken = nextToken
        }
    }
    // builder.WriteString(".")
    return builder.String()
}

type nextstate map[string]int

func main(){
    fmt.Print("\033[H\033[2J")
    fmt.Println("Please wait! Getting States Ready.....")

    states := make(map[string]nextstate)
    
    filenames := getFileNames()
    groupedTokens := []string{}

    for _, values := range filenames{
        content := fileReader(values)
        groupedTokens = stateGenerator(&states, content)
    }

    for {
        var maxLength int
        var startingToken string
        var startingPoints string

        fmt.Print("Do You want some Staring Suggestions (yes/no): ")
        _, err3 := fmt.Scanf("%s", &startingPoints)

        if err3 != nil {
            fmt.Println("Error reading input:", err3)
            return
        }

        if startingPoints == "yes"{
            for _ , val :=  range groupedTokens{
                fmt.Println(val)
            }
            fmt.Println("\n+++++++++++++++++++++++++++++")
            fmt.Println("")
        }

        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Enter starting Words: ")
        startingToken, err := reader.ReadString('\n')

        if err != nil {
            fmt.Println("Error reading input:", err)
            return
        }

        fmt.Print("Enter Max Length: ")
        _, err1 := fmt.Scanf("%d", &maxLength)

        if err1 != nil {
            fmt.Println("Error reading input:", err1)
            return
        }

        fmt.Println(storyGenerator(&states, strings.TrimSpace(startingToken), maxLength))

        var cont string 
        fmt.Println("\n+++++++++++++++++++++++++++++")
        fmt.Println("")
        fmt.Print("Want to Continue (yes/no): ")
        _, err2 := fmt.Scanf("%s", &cont)

        if err2 != nil {
            fmt.Println("Error reading input:", err2)
            return
        }

        if cont != "yes" {
            break
        }

        fmt.Println("\n+++++++++++++++++++++++++++++")
        fmt.Println("")
    }
}