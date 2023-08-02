package mini

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plugin/container"
	"github.com/sqjian/go-kit/plugin/loader"
	"github.com/sqjian/go-kit/plugin/plugin/go_native/enter"
	"github.com/sqjian/go-kit/plugin/plugin/go_native/leave"
	"github.com/sqjian/go-kit/plugin/proto"
	"github.com/sqjian/go-kit/plugin/schema"
	"github.com/sqjian/go-kit/plugin/tool"
	"github.com/sqjian/go-kit/uid"
)

type Cfg struct {
	Viper  *viper.Viper
	Logger log.Log
	Plugs  []schema.NewPlug
}

func NewMinimal(fn func(cfg *Cfg)) (*Mini, error) {
	cfg := &Cfg{}
	fn(cfg)

	mini := &Mini{
		viper:  cfg.Viper,
		logger: cfg.Logger,
		Plugs:  append(cfg.Plugs, enter.NewPlug, leave.NewPlug),
	}
	if err := mini.Init(); err != nil {
		return nil, err
	}
	return mini, nil
}

type Mini struct {
	logger    log.Log
	viper     *viper.Viper
	Plugs     []schema.NewPlug
	container *container.Container

	uniqueGenerator uid.Uid
}

func (m *Mini) FInit() error {
	return m.container.FInit()
}

func (m *Mini) Process(dag Dag, msg *proto.Msg, opts ...Opt) ([]byte, error) {
	var pluginTools schema.PlugTools = &tool.PluginToolsImpl{}

	schedulerOpts := newDefOptions()
	for _, opt := range opts {
		opt.apply(schedulerOpts)
	}
	schedulerOpts.kvs.Range(func(key, value any) bool {
		pluginTools.Set(key, value)
		return true
	})

	if len(dag.Id) == 0 {
		dag.Id, _ = m.uniqueGenerator.Gen()
		m.logger.Infof("id:%v,use snowflake algorithm to generate dag.id:%v", dag.Id, dag.Id)
	}
	pluginTools.Set("dag.id", dag.Id)

	m.logger.Infof("id:%v,dag:%#v", dag.Id, dag)

	steps := func() []string {
		steps := append([]string{"enter"}, dag.Steps...)
		steps = append(steps, "leave")
		return steps
	}()

	plugins, pluginsErr := m.container.Get(steps...)
	if pluginsErr != nil {
		m.logger.Errorf("id:%v,container.get failed,err:%v", dag.Id, pluginsErr)
		return nil, pluginsErr
	}
	for _, plugin := range plugins {
		m.logger.Infof("id:%v,checking if plug:%v of dag:%v is interested in", dag.Id, plugin.Name(), dag)
		action, actionErr := plugin.Interest(msg, pluginTools)
		if actionErr != nil {
			m.logger.Infof("id:%v,checking if plug:%v of dag:%v is interested in failed:%v", dag.Id, plugin.Name(), dag, actionErr)
			return nil, actionErr
		}
		if action == schema.SKIP {
			m.logger.Warnf("id:%v,plug:%v not interested in,skip", dag.Id, plugin.Name())
			continue
		}
		m.logger.Infof("id:%v,plug:%v,about to PreProcessing...", dag.Id, plugin.Name())
		_, preProcessErr := plugin.PreProcess(msg, pluginTools)
		if preProcessErr != nil {
			m.logger.Errorf("id:%v,plug:%v,preProcess failed:%v", dag.Id, plugin.Name(), preProcessErr)
			return nil, preProcessErr
		}
		m.logger.Infof("id:%v,plug:%v,about to Process...", dag.Id, plugin.Name())
		_, processErr := plugin.Process(msg, pluginTools)
		if processErr != nil {
			m.logger.Errorf("id:%v,plug:%v,Process failed:%v", dag.Id, plugin.Name(), processErr)
			return nil, processErr
		}
		m.logger.Infof("id:%v,plug:%v,about to PostProcess...", dag.Id, plugin.Name())
		_, postProcessErr := plugin.PostProcess(msg, pluginTools)
		if postProcessErr != nil {
			m.logger.Errorf("id:%v,plug:%v,PostProcess failed:%v", dag.Id, plugin.Name(), postProcessErr)
			return nil, postProcessErr
		}
	}

	return pluginTools.Read(), nil
}
func (m *Mini) Init() error {
	_container, _containerErr := func() (*container.Container, error) {
		c := container.NewContainer(func(cfg *container.Config) {
			cfg.Viper = m.viper
			cfg.Logger = m.logger
		})
		if err := c.Init(); err != nil {
			m.logger.Errorf("%v", err.Error())
			return nil, err
		}
		return c, nil
	}()
	if _containerErr != nil {
		m.logger.Errorf(_containerErr.Error())
		return _containerErr
	}
	m.container = _container

	_loader, _loaderErr := func() (*loader.Loader, error) {
		l := loader.NewLoader(func(cfg *loader.Cfg) {
			cfg.Viper = m.viper
			cfg.Logger = m.logger
		})
		if err := l.Init(); err != nil {
			m.logger.Errorf(err.Error())
			return nil, err
		}
		return l, nil
	}()
	if _loaderErr != nil {
		m.logger.Errorf(_loaderErr.Error())
		return _loaderErr
	}

	{
		plugins, pluginsErr := _loader.Load(m.Plugs...)
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

	{
		generator, generatorErr := uid.NewGenerator(
			uid.Snowflake,
			uid.WithSnowflakeNodeId(1),
		)
		if generatorErr != nil {
			m.logger.Errorf(generatorErr.Error())
			return generatorErr
		}
		m.uniqueGenerator = generator
	}

	return nil
}
