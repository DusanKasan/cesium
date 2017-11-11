package cesium

type err string

func (e err) Error() string {
	return string(e)
}

const DownstreamUnableToKeepUpError = err("Downstream is unable to keep up")
const NoEmissionOnSynchronousSinkError = err("No emissions received on the sink, epected at least one")
const TimeoutError = err("Timeout")
