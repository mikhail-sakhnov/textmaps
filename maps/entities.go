package maps

type SingleMap struct {
	err      error
	FileName string      `json:"file_name"`
	Tree     TextMapTree `json:"tree"`
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
