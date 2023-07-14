package uid

type OptionFunc func(*generator)

func WithSnowflakeNodeId(NodeId int64) OptionFunc {
	return func(generator *generator) {
		generator.snowflake.id = NodeId
	}
}
