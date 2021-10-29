package gui
//
// import (
// 	"github.com/gotk3/gotk3/glib"
// 	"github.com/gotk3/gotk3/gtk"
// 	"github.com/sirupsen/logrus"
// )
//
// type GtkBuilder struct {
// 	builder *gtk.Builder
// 	logger  *logrus.Logger
// }
//
// func NewGtkBuilder(fileName string, logger *logrus.Logger) *GtkBuilder {
// 	builder := new(GtkBuilder)
// 	builder.createBuilder(fileName)
// 	builder.logger = logger
// 	return builder
// }
//
// func (s *GtkBuilder) destroy() {
// 	s.builder = nil
// 	s.logger = nil
// }
//
// func (s *GtkBuilder) createBuilder(gladeFileName string) {
// 	gladePath, err := getResourcePath(gladeFileName)
// 	if err != nil {
// 		s.logger.Error(err)
// 		panic(err)
// 	}
//
// 	builder, err := gtk.BuilderNewFromFile(gladePath)
// 	if err != nil {
// 		s.logger.Error(err)
// 		panic(err)
// 	}
//
// 	s.builder = builder
// }
//
// func (s *GtkBuilder) getObject(name string) glib.IObject {
// 	obj, err := s.builder.GetObject(name)
// 	if err != nil {
// 		s.logger.Error(err)
// 		panic(err)
// 	}
//
// 	return obj
// }
