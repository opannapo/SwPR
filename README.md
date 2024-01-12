## Initiate The Project

To start working, execute

```
make init
```

## Running

To run the project, run the following command:

```
docker-compose up --build
```

You should be able to access the API at http://localhost:1111

If you change `database.sql` file, you need to reinitate the database by running:

```
docker-compose down --volumes
```

## Testing

To run test, run the following command:

```
make test
```


# Result (Unittest Coverage 71.7%):
![Screenshot from 2023-12-23 01-48-27](https://github.com/opannapo/SwPR/assets/18698574/705f2814-d875-4407-9dd1-32838dc07d52)
![Screenshot from 2023-12-23 01-51-33](https://github.com/opannapo/SwPR/assets/18698574/e7560051-c801-49fc-b82f-520c8770d465)


# Doc : Postman Collection (dir /doc)
![Screenshot from 2023-12-22 17-06-54](https://github.com/opannapo/SwPR/assets/18698574/55df6f06-a2f9-4496-bc90-d1ee54be6fd6)
