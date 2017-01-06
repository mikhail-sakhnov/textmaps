package maps

import (
	"strings"
)

type SingleMap struct {
	err      error
	FileName string `json:"file_name"`
	Nodes    []MapNode `json:"nodes"`
}

type MapNode struct {
	Level int `json:"level"`
	Text  string    `json:"text"`
	Nodes []MapNode `json:"nodes,omitempty"`

}

type MapsCollection []SingleMap

func parseNodes(rawContent string) ([]MapNode, error) {
	currentLevel := 0

	rawContent = strings.Replace(rawContent, "\t", "    ", -1)
	lines := strings.Split(rawContent, "\n")
	tmpResult := make([]MapNode, 0, len(lines))

	for _, line := range lines {
		node := parseLine(line, currentLevel)
		tmpResult = append(tmpResult, node)
	}
	return tmpResult, nil
}

func parseLine(line string, previousLine string, level int) MapNode {
	// Если начинается с 1. , 2. , 3. и тд - верхний уровень
	// если начинается с  * вложенный уровень, уровень по пробелам
	// если начинается просто с пробелов - комментарий к последней непробельной ноде
	return MapNode{
		Level: level,
		Text: line,
	}
}