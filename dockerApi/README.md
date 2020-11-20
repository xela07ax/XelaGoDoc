## Docker управленец

Запустим контейнер
```shell script
curl -X POST http://localhost:8180/kod/runContainer
```
Контейнер запускаетс, надо реализовать остановку контейнера по ресту, потом разобраться с подключением файлов.  

```shell script
curl -X POST http://localhost:8180/kod/stopContainer -d "8dc5e1d574b42ad8fa466ad3dc8af5aa8ff1870d4b0ef6672814a873d47b302d"
```

## О xaljer/kodexplorer
Запуск контейнера с командной строки
