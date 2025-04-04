# BeeTimeClock


## Development

You need a postgresql database

```
docker run -d --name btc -p 5432:5432 -e POSTGRES_PASSWORD=verysecretpassword postgres:16
```

After that you can start the backend with

```
make develop-backend
```

And frontend with
```
make develop-frontend
```


Happy Coding!
