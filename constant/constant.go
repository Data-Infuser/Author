package constant

import "time"

const ServiceDev = "dev"
const ServiceStage = "stage"
const ServiceProd = "prod"

const KeyApp = "App:"
const KeyToken = "Token:"
const KeyOperation = "Op:"
const KeyAuth = "Auth:"               // Auth:{TokenId}:{AppId}, 키-앱 인증 정보
const KeyAppTrafficPrefix = "AppTf:"  // AppTf:{AppId}:{Unit}, 앱의 단위시간당 트래픽 허용치
const KeyTrafficPrefix = "Tf:"        // Tf:{TokenId}:{AppId}:{unit}, 키-앱-단위 호출 횟수
const KeyTrafficDetailPrefix = "TfD:" // TfD:{TokenId}:{AppId}:{OperationId}:{unit}, 키-앱-단위 호출 횟수
const KEY_TRAFFIC_SET = "TrafficSet:"
const KeyTrafficDetailSet = "TrafficDetailSet:"
const KEY_TRAFFIC_QUEUE = "TrafficQueue"

func GetTrafficUnits() []string {
	return []string{
		"hour", "day", "month",
	}
}

const JwtSecret = "infuser-auther-jwt-secret"
const JwtExpInterval = 1 * time.Hour
const RefreshTokenExpInterval = 24 * time.Hour
