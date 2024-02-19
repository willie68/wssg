package interfaces

import "github.com/willie68/wssg/internal/model"

type Generator interface {
	RenderHTML(layoutName string, pg model.Page) error
}
