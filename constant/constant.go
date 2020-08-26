package constant

import "time"

const ServiceDev = "dev"
const ServiceStage = "stage"
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

const JwtSecret = "infuser-auther-jwt-secret"
const JwtExpInterval = 1 * time.Hour
const RefreshTokenExpInterval = 24 * time.Hour
