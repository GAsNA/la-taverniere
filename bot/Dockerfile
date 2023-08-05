FROM		golang:1.20

WORKDIR		/app

COPY		go.mod go.sum ./
RUN			go mod download

COPY		src/*.go ./
COPY		src/commands/*.go ./
COPY		src/db/*.go ./
COPY		src/utils/*.go ./
COPY		src/youtube/*.go ./
COPY		src/other/*.go ./

RUN			go build -o /la-taverniere

CMD			["/la-taverniere"]
