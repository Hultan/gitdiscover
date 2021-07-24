package gui

import (
"github.com/gotk3/gotk3/glib"
"github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"
)

type SoftBuilder struct {
	builder *gtk.Builder
	logger *logrus.Logger
}

func SoftBuilderNew(fileName string, logger *logrus.Logger) *SoftBuilder {
	builder := new(SoftBuilder)
	builder.createBuilder(fileName)
	builder.logger = logger
	return builder
}
func (s *SoftBuilder) destroy() {
	s.builder = nil
	s.logger = nil
}

func (s *SoftBuilder) createBuilder(gladeFileName string) {
	gladePath, err := getResourcePath(gladeFileName)
	if err != nil {
		s.logger.Error(err)
		panic(err)
	}

	builder, err := gtk.BuilderNewFromFile(gladePath)
	if err != nil {
		s.logger.Error(err)
		panic(err)
	}

	s.builder = builder
}

func (s *SoftBuilder) getObject(name string) glib.IObject {
	obj, err := s.builder.GetObject(name)
	if err != nil {
		s.logger.Error(err)
		panic(err)
	}

	return obj
}


