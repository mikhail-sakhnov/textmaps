package entities

type SingleMap struct {
	FileName string      `json:"file_name"`
	Tree     *TextMapTree `json:"tree,omitempty"`
}

type MapNode struct {
	Type      string          `json:"type"`
	Level     int             `json:"level"`
	Text      string          `json:"text"`
	Childrens NodesCollection `json:"nodes,omitempty"`
}

type TextMapTree struct {
	Childrens NodesCollection `json:"root"`
}

type NodesCollection []MapNode
type MapsCollection []SingleMap
