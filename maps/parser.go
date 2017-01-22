package maps

import (
	"regexp"
	"strings"
	"textmap/maps/entities"
)

type nodeType int64

const (
	firstLevel = nodeType(iota)
	embededNode
	descriptionNode
)

type ParserFactory struct{}

func (pf ParserFactory) Get(rawContent string) Parser {
	rawContent = strings.Replace(rawContent, "\t", "    ", -1)
	lines := []string{}
	for _, line := range strings.Split(rawContent, "\n") {
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return &simpleParser{
		lines: lines,
		parsed: map[int]bool{},
	}
}

type Parser interface {
	Parse() (entities.TextMapTree, error)
}

type simpleParser struct {
	lines []string
	parsed map[int]bool
}

func suffixLength(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

type lineInfo struct {
	lineNum int
	node    entities.MapNode
}

func (sp *simpleParser) Parse() (entities.TextMapTree, error) {
	var result entities.TextMapTree
	var firstLines []lineInfo
	for lineNum, line := range sp.lines {
		if line == "" {
			continue
		}
		if sp.isFirstLevel(line) {
			parsedLine, err := sp.parseFirstLine(line)
			if err != nil {
				return result, err
			}
			firstLines = append(firstLines,
				lineInfo{lineNum, parsedLine},
			)
		}
	}
	for _, firstLine := range firstLines {
		childrens, err := sp.parseChildrens(firstLine.lineNum, firstLine.node.Level)
		if err != nil {
			return result, err
		}
		firstLine.node.Childrens = childrens
		result.Childrens = append(result.Childrens, firstLine.node)
	}
	return result, nil
}

func (sp simpleParser) parseFirstLine(line string) (entities.MapNode, error) {
	return entities.MapNode{sp.typeOfForHuman(line), 1, strings.TrimLeft(line, "0123456789 ."), entities.NodesCollection{}}, nil
}

// TODO: do it more intellegent, without fixed tab size, just comparing with previous line
const tabSize = 4
func (sp simpleParser) parseLine(line string) (entities.MapNode, error) {
	return entities.MapNode{sp.typeOfForHuman(line), suffixLength(line)/tabSize + 1, strings.TrimLeft(line, " *"), entities.NodesCollection{}}, nil
}

func (sp *simpleParser) parseChildrens(fromLine int, fromLevel int) (entities.NodesCollection, error) {
	shouldStop := false
	result := entities.NodesCollection{}
	for i := fromLine+1; !shouldStop && i < len(sp.lines); i++ {
		if sp.parsed[i] {
			continue
		}
		node, err := sp.parseLine(sp.lines[i])
		if err != nil {
			return result, err
		}
		if node.Level <= fromLevel {
			break
		}
		childrens, err := sp.parseChildrens(i, node.Level)
		if err != nil {
			return result, err
		}
		node.Childrens = append(node.Childrens, childrens...)
		result = append(result, node)
		sp.parsed[i] = true
	}
	return result, nil
}

func (sp simpleParser) isFirstLevel(line string) bool {
	return sp.typeOf(line) == firstLevel
}

func (sp simpleParser) typeOf(line string) nodeType {
	rootPattern := regexp.MustCompile("^[0-9]+\\.")
	embedPattern := regexp.MustCompile("^\\s*[*]+")

	if embedPattern.MatchString(line) {
		return embededNode
	}
	if rootPattern.MatchString(line) {
		return firstLevel
	}
	return descriptionNode
}

func (sp simpleParser) typeOfForHuman(line string) string {
	return map[nodeType]string{
		embededNode: "embed",
		descriptionNode: "description",
		firstLevel: "first_level",
	}[sp.typeOf(line)]
}