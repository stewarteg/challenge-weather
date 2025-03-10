# challengeWeather

Para rodar local com docker:

docker build --no-cache -t challenge-weather .
docker-compose up --force-recreate

link para realizar chamada após rodar app:
http://localhost:8080/cep?cep=57062090

Apos rodar, Verificar spans em: 
http://localhost:9411/zipkin/?lookback=15m&endTs=1741643307072&limit=10

