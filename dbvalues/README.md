# DB Values

So the reason for this example is that I was struggling to write a struct to the database. Specifically it was the `r2.Point` struct from `github.com/golang/geo`. It's a simple 2d coordinate type with `float64` X and Y values. That's it. And I wanted to store it in the database.

To see an example of the error:
```bash
go run broken.go
```

You should see the following result:
```
2018/01/24 20:55:10 sql: converting argument $2 type: unsupported type r2.Point, a struct
exit status 1
```

The reason for this is that implementers of DB drivers need to support the following types:
* `int64`
* `float64`
* `bool`
* `[]byte`
* `string`
* `time.Time`

Notice that `r2.Point` isn't one of these?

So what are we going to do about it? Well the people behind the `database/sql` package thought about this a little and provided the [Valuer](https://godoc.org/database/sql/driver#Valuer) interface to help with this conversion.

We need to make sure that we try writing a type that meets this interface so that the driver can do the conversion to a type that it understands.

Take a look in `write.go` to see how we approach this. If you're new to Go, notice that we've used composition to write a type to add more behaviour to `r2.Point`. Run the file with:
```bash
go run write.go
```

Now we get _this_ error:
```
2018/01/24 21:33:12 Converted r2.Point to a string: (1.000000000000, 2.000000000000)
2018/01/24 21:33:12 sql: Scan error on column index 1: unsupported Scan, storing driver.Value type []uint8 into type *main.Location
exit status 1
```

Again, package implementers to the rescue - we can use [Scanner](https://godoc.org/database/sql#Scanner). The DB driver will use this interface to get our struct to do the converstion from the type that the DB persisted the data as to the type that we want. Wierdly, as we're using `sqlite3` this turns out to be `[]uint8`. The last file, `writeread.go` implements the `Scan` function. Run it with:
```bash
go run writeread.go
```
