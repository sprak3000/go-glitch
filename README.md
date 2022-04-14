# go-glitch

This package is designed to help with handling errors from databases and web services. The `glitch.DataError` structure
allows the logic layer of your code to use its error code to handle each machine-readable code in a user-friendly way.

## PostgreSQL errors

You can convert a `lib/pq` error into a `glitch.DataError` using `postgres.ToDataError()`.

**Example:**

```go
query := "CALL do_work($1)"
_, err := d.conn.ExecContext(ctx, query, workID)
return postgres.ToDataError(err, fmt.Sprintf("error doing work ID %s", workID))
```

## HTTP problem responses

### Handling a problem response from an API

If a web service implements the ["Problem Details for HTTP APIs"](https://datatracker.ietf.org/doc/rfc7807)
specification, you can unmarshal the API responses into the `glitch.HTTPProblem` structure and then convert that into
a `glitch.DataError` to return to the client logic using `glitch.FromHTTPProblem`.

**Example:**

```go
// Make an API call that returns the HTTP status code as an integer along with the API response as a byte slice
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

### Returning a problem response

Your own service can return an RFC7807 problem response.

**Example:**

```go
// GetUser() will return a user structure on success or a glitch.DataError on failure.
user, err := db.GetUser(id)

if err != nil {
	var status int
	var err string

	switch err.Code() {
	case userNotFound:
		status = http.StatusNotFound
		err = "User could not be found"
	case dbConnection:
		status = http.StatusServiceUnavailable
		err = "Database error. Contact customer support."
	default:
		status = http.StatusInternalServerError
		err = "Service error. Contact customer support."
    }

	httpProblem := glitch.HTTPProblem{
        Status: status,
		Detail: err,
		// ...
    }

    ret, _ := json.Marshal(httpErr)
    w.WriteHeader(status)
    w.Write(ret)
    return
}
```
