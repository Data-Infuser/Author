package constant

const ServiceDev = "dev"
const ServiceProd = "prod"

const KeyApp = "App:"
const KeyToken = "Token:"
const KeyOperation = "Op:"
const KeyAuth = "Auth:"
const KeyAppTrafficPrefix = "AppTf:"
const KeyTrafficPrefix = "Tf:"
const KEY_TRAFFIC_SET = "TrafficSet:"
const KEY_TRAFFIC_QUEUE = "TrafficQueue"

func GetTrafficUnits() []string {
	return []string{
		"hour", "day", "month",
	}
}
