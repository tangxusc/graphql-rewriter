package plugin_manager

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/ast"
	"path/filepath"
	"plugin"
	"strings"
)

var ConvertError = errors.New("convert error")

type pluginInstance struct {
	path          string
	rewriteQl     func(*ast.QueryDocument) error
	rewriteResult func([]byte) ([]byte, error)
}

func (i *pluginInstance) reWriteQl(gql *ast.QueryDocument) error {
	if i.rewriteQl != nil {
		return i.rewriteQl(gql)
	}
	logrus.Warnf("[plugins]%s : no rewriteQl specified", i.path)
	return nil
}

func (i *pluginInstance) reWriteResult(data []byte) ([]byte, error) {
	if i.rewriteResult != nil {
		return i.rewriteResult(data)
	}
	logrus.Warnf("[plugins]%s : no rewriteResult specified", i.path)
	return data, nil
}

func startInstance(ctx context.Context, path string) (*pluginInstance, error) {
	open, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	p := &pluginInstance{}
	lookup, err := open.Lookup("ReWriteQl")
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return nil, err
		}
	} else {
		rewriter, ok := lookup.(func(*ast.QueryDocument) error)
		if !ok {
			logrus.Errorf("[plugins]%s : convert ReWriteQl function error", path)
			return nil, ConvertError
		}
		p.rewriteQl = rewriter
	}
	lookup, err = open.Lookup("ReWriteResult")
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return nil, err
		}
	} else {
		rewriter, ok := lookup.(func([]byte) ([]byte, error))
		if !ok {
			logrus.Errorf("[plugins]%s : convert ReWriteResult function error", path)
			return nil, ConvertError
		}
		p.rewriteResult = rewriter
	}
	return p, nil
}

var pluginInstances []*pluginInstance

func StartPlugins(ctx context.Context) error {
	glob, err := filepath.Glob(filepath.Join(pluginDir, "*.so"))
	if err != nil {
		return err
	}
	pluginInstances = make([]*pluginInstance, len(glob))
	for i, path := range glob {
		instance, err := startInstance(ctx, path)
		if err != nil {
			return err
		}
		pluginInstances[i] = instance
	}
	return nil
}

func ReWriteQl(gql *ast.QueryDocument) error {
	for _, instance := range pluginInstances {
		err := instance.reWriteQl(gql)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReWriteResult(data []byte) ([]byte, error) {
	var err error
	for _, instance := range pluginInstances {
		data, err = instance.reWriteResult(data)
		if err != nil {
			return data, err
		}
	}
	return data, nil
}
