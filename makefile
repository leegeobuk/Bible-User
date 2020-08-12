kakaologin:
		go fmt ./...
		go test ./auth/kakao
		go build -o bin/kakaologin lambda/kakaologin/main.go
		zip zip/kakaologin.zip bin/kakaologin