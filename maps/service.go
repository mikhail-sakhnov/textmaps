package maps

import (
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"strings"
)

var ErrNotExists = errors.New("Textmap doesn't exists")

type MapService struct {
	dataDir string
	fs      afero.Afero
	parser  interface {
		Get(string) parser
	}
}

func NewService(dataDir string) MapService {
	return MapService{
		dataDir: dataDir,
		fs:      afero.Afero{afero.NewOsFs()},
		parser:  parserFactory{},
	}
}

func (s MapService) GetAllMaps() (MapsCollection, error) {
	// todo: walk for subdirectories
	var result MapsCollection
	fileInfos, err := s.fs.ReadDir(s.dataDir)
	if err != nil {
		return result, err
	}
	for _, fileInfo := range fileInfos {
		fname := fileInfo.Name()
		if !strings.HasSuffix(fname, ".tm") {
			continue
		}
		result = append(result, SingleMap{
			FileName: fname,
		})
	}
	return result, nil
}

func (s MapService) GetMapTextContent(path string) (SingleMap, error) {
	var singleMap SingleMap
	mapPath := s.dataDir + "/" + path
	fmt.Println(mapPath)
	if exists, _ := s.fs.Exists(mapPath); !exists {
		return singleMap, ErrNotExists
	}
	rawContent, err := s.fs.ReadFile(mapPath)
	if err != nil {
		return singleMap, err
	}
	singleMap.FileName = path
	parser := s.parser.Get(string(rawContent))
	if err != nil {
		return singleMap, err
	}
	tree, err := parser.Parse()
	if err != nil {
		return singleMap, err
	}
	singleMap.Tree = tree
	return singleMap, nil
}
