kakaologin:
		go fmt ./...
		go build -o bin/kakaologin lambda/kakaologin/main.go
		zip zip/kakaologin.zip bin/kakaologin