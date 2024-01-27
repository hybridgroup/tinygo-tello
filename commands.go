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

// FlipType is used for the various flips supported by the Tello.
type FlipType int

const (
	// FlipFront flips forward.
	FlipFront FlipType = 0

	// FlipLeft flips left.
	FlipLeft FlipType = 1

	// FlipBack flips backwards.
	FlipBack FlipType = 2

	// FlipRight flips to the right.
	FlipRight FlipType = 3

	// FlipForwardLeft flips forwards and to the left.
	FlipForwardLeft FlipType = 4

	// FlipBackLeft flips backwards and to the left.
	FlipBackLeft FlipType = 5

	// FlipBackRight flips backwards and to the right.
	FlipBackRight FlipType = 6

	// FlipForwardRight flips forwards and to the right.
	FlipForwardRight FlipType = 7
)
