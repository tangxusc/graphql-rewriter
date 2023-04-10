package plugin_manager

import (
	"context"
	"github.com/vektah/gqlparser/ast"
	"path/filepath"
	"plugin"
)

type pluginInstance struct {
	path     string
	rewriter func(*ast.QueryDocument) error
}

func (i *pluginInstance) reWrite(gql *ast.QueryDocument) error {
	return i.rewriter(gql)
}

func startInstance(ctx context.Context, path string) (*pluginInstance, error) {
	open, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	lookup, err := open.Lookup("ReWrite")
	if err != nil {
		return nil, err
	}
	rewriter := lookup.(func(*ast.QueryDocument) error)
	return &pluginInstance{rewriter: rewriter}, nil
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

func ReWrite(gql *ast.QueryDocument) error {
	for _, instance := range pluginInstances {
		err := instance.reWrite(gql)
		if err != nil {
			return err
		}
	}
	return nil
}
