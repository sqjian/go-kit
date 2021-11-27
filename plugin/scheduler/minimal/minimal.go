package minimal

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plugin/container"
	"github.com/sqjian/go-kit/plugin/loader"
	"github.com/sqjian/go-kit/plugin/proto"
	"github.com/sqjian/go-kit/plugin/schema"
	"github.com/sqjian/go-kit/plugin/tool"
)

type Cfg struct {
	Viper              *viper.Viper
	Logger             log.API
	GoNativeGenerators []schema.NewPlugin
}

func NewMinimal(cfg *Cfg) (*Minimal, error) {
	minimal := &Minimal{
		viper:              cfg.Viper,
		logger:             cfg.Logger,
		goNativeGenerators: cfg.GoNativeGenerators,
	}
	if err := minimal.Init(); err != nil {
		return nil, err
	}
	return minimal, nil
}

type Minimal struct {
	logger             log.API
	viper              *viper.Viper
	goNativeGenerators []schema.NewPlugin
	container          container.Container
}

func (m *Minimal) FInit() error {
	return m.container.FInit()
}

func (m *Minimal) Process(dag Dag, msg *proto.Msg, opts ...Option) ([]byte, error) {
	var pluginTools schema.PluginTools = &tool.PluginToolsImpl{}

	schedulerOpts := newDefaultOptions()
	for _, opt := range opts {
		opt.apply(schedulerOpts)
	}
	for key, val := range schedulerOpts.kvs {
		pluginTools.Set(key, val)
	}
	m.logger.Infof("dag:%v", dag)
	steps := func() []string {
		steps := append([]string{"enter"}, dag.Steps...)
		steps = append(steps, "leave")
		return steps
	}()

	plugins, pluginsErr := m.container.Get(steps...)
	if pluginsErr != nil {
		m.logger.Errorf("container.get failed,err:%v", pluginsErr)
		return nil, pluginsErr
	}
	for _, plugin := range plugins {
		m.logger.Infof("checking if plugin:%v of dag:%v is interested in", plugin.Name(), dag)
		action, actionErr := plugin.Interest(msg, pluginTools)
		if actionErr != nil {
			m.logger.Infof("checking if plugin:%v of dag:%v is interested in failed:%v", plugin.Name(), dag, actionErr)
			return nil, actionErr
		}
		if action == schema.SKIP {
			m.logger.Warnf("plugin:%v not interested in,skip", plugin.Name())
			continue
		}
		m.logger.Infof("plugin:%v,about to PreProcessing...", plugin.Name())
		_, preProcessErr := plugin.PreProcess(msg, pluginTools)
		if preProcessErr != nil {
			m.logger.Errorf("plugin:%v,preProcess failed:%v", plugin.Name(), preProcessErr)
			return nil, preProcessErr
		}
		m.logger.Infof("plugin:%v,about to Process...", plugin.Name())
		_, processErr := plugin.Process(msg, pluginTools)
		if processErr != nil {
			m.logger.Errorf("plugin:%v,Process failed:%v", plugin.Name(), processErr)
			return nil, processErr
		}
		m.logger.Infof("plugin:%v,about to PostProcess...", plugin.Name())
		_, postProcessErr := plugin.PostProcess(msg, pluginTools)
		if postProcessErr != nil {
			m.logger.Errorf("plugin:%v,PostProcess failed:%v", plugin.Name(), postProcessErr)
			return nil, postProcessErr
		}
	}

	return pluginTools.Read(), nil
}
func (m *Minimal) Init() error {
	m.container = container.NewContainer(&container.Cfg{Viper: m.viper, Logger: m.logger})
	if err := m.container.Init(); err != nil {
		m.logger.Errorf("%v", err.Error())
		return err
	}

	loader := loader.NewLoader(&loader.Cfg{Viper: m.viper, Logger: m.logger})
	if err := loader.Init(); err != nil {
		m.logger.Errorf(err.Error())
		return err
	}

	{
		// init go naive
		plugins, pluginsErr := loader.Load(m.goNativeGenerators...)
		if pluginsErr != nil {
			m.logger.Errorf(pluginsErr.Error())
			return pluginsErr
		}

		addErr := m.container.Add(plugins)
		if addErr != nil {
			m.logger.Errorf(addErr.Error())
			return addErr
		}
	}
	return nil
}
