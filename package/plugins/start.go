package plugins

import (
	"context"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/second-state/WasmEdge-go/wasmedge"
	bindgen "github.com/second-state/wasmedge-bindgen/host/go"
	"os"
	"path/filepath"
)

type pluginInstance struct {
	path string
	vm   *wasmedge.VM
	bg   *bindgen.Bindgen
	conf *wasmedge.Configure
}

func (i *pluginInstance) release() {
	i.bg.Release()
	i.vm.Release()
	i.conf.Release()
}

//TODO: 考虑并行执行
func (i *pluginInstance) patch(gql []byte) (jsonpatch.Patch, error) {
	execute, _, err := i.bg.Execute("test", string(gql), "test")
	if err != nil {
		return nil, err
	}
	s := execute[0].(string)
	return jsonpatch.DecodePatch([]byte(s))
}

func startInstance(ctx context.Context, path string) (*pluginInstance, error) {
	conf := wasmedge.NewConfigure(wasmedge.WASI)
	vm := wasmedge.NewVMWithConfig(conf)
	wasi := vm.GetImportModule(wasmedge.WASI)
	wasi.InitWasi(
		[]string{},      // The args
		os.Environ(),    // The envs
		[]string{".:."}, // The mapping preopens
	)
	if err := vm.LoadWasmFile(path); err != nil {
		return nil, err
	}
	if err := vm.Validate(); err != nil {
		return nil, err
	}
	bg := bindgen.New(vm)
	bg.Instantiate()

	return &pluginInstance{path: path, vm: vm, bg: bg, conf: conf}, nil
}

var pluginInstances []*pluginInstance

func StartPlugins(ctx context.Context) error {
	wasmedge.SetLogDebugLevel()

	glob, err := filepath.Glob(filepath.Join(pluginDir, "*.wasm"))
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

	go func() {
		for {
			select {
			case <-ctx.Done():
				for _, instance := range pluginInstances {
					instance.release()
				}
				return
			}
		}
	}()

	return nil
}

func Patchs(gql []byte) ([]jsonpatch.Patch, error) {
	patches := make([]jsonpatch.Patch, len(pluginInstances))
	for i, instance := range pluginInstances {
		patch, err := instance.patch(gql)
		if err != nil {
			return patches, err
		}
		patches[i] = patch
	}
	return patches, nil
}
