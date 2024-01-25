module github.com/hybridgroup/tinygo-tello

go 1.18

replace tinygo.org/x/drivers => ../drivers

require (
	tinygo.org/x/drivers v0.26.1-0.20240117074700-3c5e17423a16
	tinygo.org/x/tinyfont v0.4.0
	tinygo.org/x/tinyterm v0.3.1-0.20231207163921-6842651de7e1
)

require github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
