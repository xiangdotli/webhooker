webhooker from alertmanager to telegram

make build
webhooker -file=./config.json

or

docker build -t webhooker .
docker run -p 9091:9091 webhooker
