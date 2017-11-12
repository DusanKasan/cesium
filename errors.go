package cesium

type err string

func (e err) Error() string {
	return string(e)
}

// DownstreamUnableToKeepUpError is emitted from a Flux when using Error
// backpressure strategy and the downstream can not process items as fast as
// they are emitted.
const DownstreamUnableToKeepUpError = err("Downstream is unable to keep up")

// NoEmissionOnSynchronousSinkError is emitted from a Flux/Mono when using the
// flux.Generate/mono.Generate factory methodwhen no methods are called on the
// passed SynchronousSink.
const NoEmissionOnSynchronousSinkError = err("No emissions received on the sink, epected at least one")

// TimeoutError is returned by blocking methods (Mono.BlockTimeout,
// Flux.BlockFirstTimeout and Flux.BlockLastTimeout) when no matching items
// would be emitted in the specified timeout duration.
const TimeoutError = err("Timeout")
