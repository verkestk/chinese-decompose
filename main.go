package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	compositionTypeGraphicalPrimitive = "一" // non composition (second character is always a deformed version of another character) - only has leftComponent
	compositionTypeHorizontal         = "吅" // horizontal composition (when repetition, the second character is deformed)
	compositionTypeVertical           = "吕" // horizontal composition (when repetition, the second character is deformed)
	compositionTypeInclusion          = "回" // inclusion of the second character inside the first (门, 囗, 匚...)
	compositionTypeVerticalRepetition = "咒" // vertical composition, the top part being a repetition.
	compositionTypeHorizontalOfThree  = "弼" // horizontal composition of three, the third being the repetition of the first.
	compositionTypeRepetitionOfThree  = "品" // repetition of three - only has leftComponent
	compositionTypeRepetitionOfFour   = "叕" // repetition of four - only has leftComponent
	compositionTypeVerticalSeparated  = "冖" // vertical composition, separated by "冖".
	compositionTypeSuperposition      = "+" // graphical superposition or addition

	positionPrimitive      = iota // only applies to compositionTypeGraphicalPrimitive
	positionLeft                  // applies to leftComponent of compositionTypeHorizontal and compositionTypeHorizontalOfThree
	positionRight                 // applies to rightComponent of compositionTypeHorizontal and leftComponent of compositionTypeHorizontalOfThree
	positionTop                   // applies to leftComponent of compositionTypeVertical and compositionTypeRepetitionOfThree and compositionTypeVerticalSeparated
	positionBottom                // applies to rightComponet of compositionTypeVertical and compositionTypeVerticalRepetition and compositionTypeVerticalSeparated
	positionOuter                 // applies to leftComponent of compositionTypeInclusion
	positionInner                 // applies to rightComponent of compositionTypeInclusion
	positionMiddle                // applies to rightComponent of compositionTypeHorizontalOfThree
	positionTopRepeated           // applies to leftComponent of compositionTypeVerticalRepetition and compositionTypeRepetitionOfFour
	positionBottomRepeated        // applies to rightComponent of compositionTypeRepetitionOfThree and compositionTypeRepetitionOfFour
	positionSuperPrimary          // applies to leftComponent of compositionTypeSuperposition
	positionSuperSecondary        // applies to rightComponent of compositionTypeSuperposition
)

var (
	positionDescriptions = map[int]string{
		positionPrimitive:      "primitive",
		positionLeft:           "left",
		positionRight:          "right",
		positionTop:            "top",
		positionBottom:         "bottom",
		positionOuter:          "outer",
		positionInner:          "inner",
		positionMiddle:         "middle",
		positionTopRepeated:    "top repeated",
		positionBottomRepeated: "bottom repeated",
		positionSuperPrimary:   "primary",
		positionSuperSecondary: "secondary",
	}
)

type vocabulary struct {
	term         string
	pinyin       string
	partOfSpeech string
	translation  string
}

type character struct {
	character             string
	strokes               uint
	compositionType       string
	leftComponent         string
	leftComponentStrokes  uint
	rightComponent        string
	rightComponentStrokes uint
	signature             string
	notes                 string
	section               string

	vocabulary []*vocabulary
}

type cluster struct {
	component  string
	position   int
	characters []*character
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("must pass vocabulary source file and decomposition database file")
	}
	if len(os.Args) > 3 {
		log.Fatal("too many args")
	}

	vocabularySourcePath := os.Args[1]
	decompositionDBPath := os.Args[2]

	charactersToVocabulary, err := loadSourceVocabulary(vocabularySourcePath)
	if err != nil {
		log.Fatal(err)
	}

	characterDB, err := loadDecompositionDB(decompositionDBPath)
	if err != nil {
		log.Fatal(err)
	}

	clusters, _ := makeClusters(charactersToVocabulary, characterDB)
	sortClusters(clusters)

	for _, c := range clusters {
		fmt.Printf("# Cluster for %s in position %s\n", c.component, positionDescriptions[c.position])
		fmt.Printf("## ")

		for _, char := range c.characters {
			fmt.Print(char.character)
		}
		fmt.Print("\n")

		for _, char := range c.characters {
			fmt.Printf("### %s vocabulary\n", char.character)
			fmt.Printf("| Term | Pinyin | PoS | Translation |\n| --- | --- | --- | --- |\n")

			for _, vocab := range char.vocabulary {
				fmt.Printf("| %s | %s | %s | %s |\n", vocab.term, vocab.pinyin, vocab.partOfSpeech, vocab.translation)
			}
		}
	}
}

func loadSourceVocabulary(path string) (map[string][]*vocabulary, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)

	values, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	inputVocabulary := make([]*vocabulary, 0)
	charactersToVocabulary := map[string][]*vocabulary{}
	for _, row := range values {
		if len(row) > 0 {
			vocab := &vocabulary{term: row[0]}
			vocab.term = row[0]
			if len(row) > 1 {
				vocab.pinyin = row[1]
				if len(row) > 2 {
					vocab.partOfSpeech = row[2]
					if len(row) > 3 {
						vocab.translation = row[3]
					}
				}
			}
			inputVocabulary = append(inputVocabulary, vocab)

			characters := splitCharacters(row[0])
			for _, char := range characters {
				if _, ok := charactersToVocabulary[char]; !ok {
					charactersToVocabulary[char] = make([]*vocabulary, 0)
				}

				charactersToVocabulary[char] = append(charactersToVocabulary[char], vocab)
			}
		}
	}

	return charactersToVocabulary, nil
}

func loadDecompositionDB(path string) (map[string]*character, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	reader.Comma = '\t'

	values, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	characters := make(map[string]*character)
	for i, row := range values {
		if len(row) != 10 {
			return nil, fmt.Errorf("unexpected number of values on row %d: %d", i, len(row))
		}

		characterStr := row[0]
		strokesStr := row[1]
		compositionTypeStr := row[2]
		leftComponentStr := row[3]
		leftComponentStrokesStr := row[4]
		rightComponentStr := row[5]
		rightComponentStrokesStr := row[6]
		signatureStr := row[7]
		notesStr := row[8]
		sectionStr := row[9]

		strokesInt, err := strconv.Atoi(strokesStr)
		if err != nil {
			return nil, err
		}

		leftComponentStrokesInt, err := strconv.Atoi(leftComponentStrokesStr)
		if err != nil {
			return nil, err
		}

		rightComponentStrokesInt, err := strconv.Atoi(rightComponentStrokesStr)
		if err != nil {
			return nil, err
		}

		characters[characterStr] = &character{
			character:             characterStr,
			strokes:               uint(strokesInt),
			compositionType:       compositionTypeStr,
			leftComponent:         leftComponentStr,
			leftComponentStrokes:  uint(leftComponentStrokesInt),
			rightComponent:        rightComponentStr,
			rightComponentStrokes: uint(rightComponentStrokesInt),
			signature:             signatureStr,
			notes:                 notesStr,
			section:               sectionStr,
		}
	}

	return characters, nil
}

func makeClusters(charactersToVocabulary map[string][]*vocabulary, characterDB map[string]*character) (clusters []*cluster, isolatedCharacters []*character) {
	knownCharacterDB := make(map[string]*character)
	for charStr, char := range characterDB {
		if charVocab, ok := charactersToVocabulary[charStr]; ok {
			knownCharacterDB[charStr] = char
			char.vocabulary = charVocab
		}
	}

	componentToPositionToCharacterMap := make(map[string]map[int][]*character)
	for _, char := range knownCharacterDB {
		// get the components and positions from each char and group by component/position match
		for component, positions := range getComponentsAndPositions(char) {
			if _, ok := componentToPositionToCharacterMap[component]; !ok {
				componentToPositionToCharacterMap[component] = make(map[int][]*character)
			}

			for _, position := range positions {
				if _, ok := componentToPositionToCharacterMap[component][position]; !ok {
					componentToPositionToCharacterMap[component][position] = make([]*character, 0)
				}

				componentToPositionToCharacterMap[component][position] = append(componentToPositionToCharacterMap[component][position], char)
			}
		}
	}

	clusters = make([]*cluster, 0)
	isolatedCharacters = make([]*character, 0)

	for component, positionCharacters := range componentToPositionToCharacterMap {
		for position, characters := range positionCharacters {
			if len(characters) < 2 {
				isolatedCharacters = append(isolatedCharacters, characters[0])
			} else {
				clusters = append(clusters, &cluster{component: component, position: position, characters: characters})
			}
		}
	}

	return clusters, isolatedCharacters
}

func getComponentsAndPositions(char *character) map[string][]int {
	components := make(map[string][]int)

	switch char.compositionType {
	case compositionTypeGraphicalPrimitive:
		components[char.leftComponent] = []int{positionPrimitive}
	case compositionTypeHorizontal:
		components[char.leftComponent] = []int{positionLeft}
	case compositionTypeVertical:
		components[char.leftComponent] = []int{positionTop}
	case compositionTypeInclusion:
		components[char.leftComponent] = []int{positionOuter}
	case compositionTypeVerticalRepetition:
		components[char.leftComponent] = []int{positionTopRepeated}
	case compositionTypeHorizontalOfThree:
		components[char.leftComponent] = []int{positionLeft, positionRight}
	case compositionTypeRepetitionOfThree:
		components[char.leftComponent] = []int{positionTop}
	case compositionTypeRepetitionOfFour:
		components[char.leftComponent] = []int{positionTopRepeated}
	case compositionTypeVerticalSeparated:
		components[char.leftComponent] = []int{positionTop}
	case compositionTypeSuperposition:
		components[char.leftComponent] = []int{positionSuperPrimary}
	}

	if char.rightComponent == "*" {
		return components
	}

	switch char.compositionType {
	case compositionTypeHorizontal:
		components[char.rightComponent] = []int{positionRight}
	case compositionTypeVertical:
		components[char.rightComponent] = []int{positionBottom}
	case compositionTypeInclusion:
		components[char.rightComponent] = []int{positionInner}
	case compositionTypeVerticalRepetition:
		components[char.rightComponent] = []int{positionBottom}
	case compositionTypeHorizontalOfThree:
		components[char.rightComponent] = []int{positionMiddle}
	case compositionTypeVerticalSeparated:
		components[char.rightComponent] = []int{positionBottom}
	case compositionTypeSuperposition:
		components[char.rightComponent] = []int{positionSuperSecondary}
	}

	return components
}

func sortClusters(clusters []*cluster) {
	// sort clusters by size
	sort.Slice(clusters, func(i, j int) bool {
		return len(clusters[i].characters) > len(clusters[j].characters)
	})

	// sort characters by stroke count
	for _, c := range clusters {
		sort.Slice(c.characters, func(i, j int) bool {
			return c.characters[i].strokes < c.characters[j].strokes
		})

		// sort vocabulary alphabetically
		for _, char := range c.characters {
			sort.Slice(char.vocabulary, func(i, j int) bool {
				return char.vocabulary[i].term < char.vocabulary[j].term
			})
		}
	}
}

// returns unique list of single characters within source term
func splitCharacters(term string) []string {
	charMap := make(map[string]bool)

	allCharacters := strings.Split(term, "")
	for _, char := range allCharacters {
		charMap[char] = true
	}

	characters := make([]string, 0)
	for char := range charMap {
		characters = append(characters, char)
	}

	return characters
}
