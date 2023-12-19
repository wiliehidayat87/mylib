module github.com/wiliehidayat87/mylib/v2

go 1.20

retract (
	v2.0.7 // Contains retractions only - failed code.
	v2.0.5 // Contains retractions only - failed code.
	v2.0.4 // Contains retractions only - failed code.
	v2.0.1 // Contains retractions only.
	v2.0.0 // Contains retractions only.
)

require github.com/mergermarket/go-pkcs7 v0.0.0-20170926155232-153b18ea13c9 // indirect
