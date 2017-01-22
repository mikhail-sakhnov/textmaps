package services

import (
	"errors"
	"github.com/spf13/afero"
	"strings"
	"textmap/maps"
	"textmap/maps/entities"
	"context"
)

var ErrNotExists = errors.New("Textmap doesn't exists")
var ErrEmptyPath = errors.New("Path can't be empty")

type MapService struct {
	dataDir string
	fs      afero.Afero
	parser  interface {
		Get(string) maps.Parser
	}
}

func NewService(dataDir string) MapService {
	return MapService{
		dataDir: dataDir,
		fs:      afero.Afero{afero.NewOsFs()},
		parser:  maps.ParserFactory{},
	}
}

func (s MapService) GetAllMaps(ctx context.Context) (entities.MapsCollection, error) {
	// todo: walk for subdirectories
	var result entities.MapsCollection
	fileInfos, err := s.fs.ReadDir(s.dataDir)
	if err != nil {
		return result, err
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		fname := fileInfo.Name()
		if !strings.HasSuffix(fname, ".tm") {
			continue
		}
		result = append(result, entities.SingleMap{
			FileName: fname,
		})
	}
	return result, nil
}

func (s MapService) GetMapTextContent(ctx context.Context, path string) (entities.SingleMap, error) {
	var singleMap entities.SingleMap

	if path == "" {
		return singleMap, ErrEmptyPath
	}
	mapPath := s.dataDir + "/" + path
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
	singleMap.Tree = &tree
	return singleMap, nil
}
