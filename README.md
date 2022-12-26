# go-glitch

[![Maintainability](https://api.codeclimate.com/v1/badges/2e89b8a1b90cfdbf67da/maintainability)](https://codeclimate.com/github/sprak3000/go-glitch/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/2e89b8a1b90cfdbf67da/test_coverage)](https://codeclimate.com/github/sprak3000/go-glitch/test_coverage)

This package is designed to help with handling errors from databases and web services. The `glitch.DataError` structure
allows the logic layer of your code to use its error code to handle each machine-readable code in a user-friendly way.

Interested in making this library better? Read through our [development guide](docs/development.md).

## Using `DataError`

Create a new `glitch.DataError`.

```go
// Original error
innerErr := errors.New('root error cause')
// The error code and user-friendly error message you wish to expose
errCode := 'ErrorAPICrash'
msg := 'API experienced an error'

err := glitch.NewDataError(innerErr, errCode, msg)
```

`DataError` satisfies the `error` interface.

```go
fmt.Println(err.Error())

// Prints:
//  Code: [ErrorAPICrash] Message: [API experienced an error] Inner error: [root error cause]"
```

Access the individual error details through `err.Code()`, `err.Inner()`.

You can also wrap a `DataError` within a `DataError`.

```go
err := glitch.NewDataError(nil, "ErrorDatabaseUnavailable", "no database connection available")
newErr := glitch.NewDataError(nil, "ErrorServiceUnavailable", "service down")
newErr.Wrap(err)

// Get the wrapped error
origErr := newErr.GetCause()
```

## Handling database errors

## PostgreSQL (`lib/pq`) errors

You can convert a [`lib/pq` error](https://pkg.go.dev/github.com/lib/pq#Error) into a `glitch.DataError` using
`postgres.ToDataError()`.

```go
query := "CALL do_work($1)"
_, err := d.conn.ExecContext(ctx, query, workID)
return postgres.ToDataError(err, fmt.Sprintf("error doing work ID %s", workID))
```

## API errors

### Handling a "Problem Details for HTTP APIs" (RFC 7807) response from an API

If a web service implements the ["Problem Details for HTTP APIs"](https://datatracker.ietf.org/doc/rfc7807)
specification, you can unmarshal the API responses into the `glitch.HTTPProblem` structure and then convert that into
a `glitch.DataError` to return to the client logic using `glitch.FromHTTPProblem`.

```go
var status int
var ret []btye
status, ret := callAPI()

if status >= 400 || status < 200 {
	prob := glitch.HTTPProblem{}
	err := json.Unmarshl(ret, &prob)
	if err != nil {
        return glitch.NewDataError(err, "ErrorJSONUnmarshal", "Could not decode error response")
    }

	return glitch.FromHTTPProblem(prob, "Error calling the API")
}
```

### Creating a "Problem Details for HTTP APIs" (RFC 7807) response

Your own service can return an [RFC 7807](https://datatracker.ietf.org/doc/rfc7807) problem response. This package
defines an additional and optional `Code` field to return an API specific error code.

```go
// GetUser() will return a user structure on success or a glitch.DataError on failure.
user, err := db.GetUser(id)

if err != nil {
	var status int
	var err string
	var title string
	var code string

	switch err.Code() {
	case userNotFound:
		status = http.StatusNotFound
		err = "User could not be found"
		title = "Not Found"
		code = "ErrorNotFound"
	case dbConnection:
		status = http.StatusServiceUnavailable
		err = "Database error. Contact customer support."
		title = "Database Error"
		code = "ErrorDatabase"
	default:
		status = http.StatusInternalServerError
		err = "Service error. Contact customer support."
		title = "Internal Error"
		code = "ErrorInternal"
    }

	httpProblem := glitch.HTTPProblem{
        Status:   status,
		Detail:   err,
		Type:     "https://example.net/validation-error",
        Title:    title,
        Instance: "/foo/bar",
        Code:     code,
    }

    ret, _ := json.Marshal(httpErr)
    w.WriteHeader(status)
    w.Write(ret)
    return
}
```
