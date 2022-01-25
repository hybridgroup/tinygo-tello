package tello

const (
	messageStart   = 0x00cc // 204
	wifiMessage    = 0x001a // 26
	videoRateQuery = 0x0028 // 40
	lightMessage   = 0x0035 // 53
	flightMessage  = 0x0056 // 86
	logMessage     = 0x1050 // 4176

	videoEncoderRateCommand = 0x0020 // 32
	videoStartCommand       = 0x0025 // 37
	exposureCommand         = 0x0034 // 52
	timeCommand             = 0x0046 // 70
	stickCommand            = 0x0050 // 80
	takeoffCommand          = 0x0054 // 84
	landCommand             = 0x0055 // 85
	flipCommand             = 0x005c // 92
	throwtakeoffCommand     = 0x005d // 93
	palmLandCommand         = 0x005e // 94
	bounceCommand           = 0x1053 // 4179
)
