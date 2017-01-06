package maps

import (
	"github.com/spf13/afero"
	"strings"
	"errors"
	"fmt"
)

var ErrNotExists = errors.New("Textmap doesn't exists")

type MapService struct {
	dataDir string
	fs afero.Afero
}

func NewService(dataDir string) MapService {
	return MapService{
		dataDir: dataDir,
		fs: afero.Afero{afero.NewOsFs()},
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
			FileName: fname ,
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
	nodes, err := parseNodes(string(rawContent))
	if err != nil {
		return singleMap, err
	}
	singleMap.Nodes = nodes
	return singleMap, nil
}