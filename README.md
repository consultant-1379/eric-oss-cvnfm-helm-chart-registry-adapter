# Helm Registry Adapter

Helm Adapter is used to maintain A2A compatibility with EO-CM (if EO_CM using old onboarding compatibility with new CVNFM must be kept).
It handles old Helm Registry requests and redirects them to the OCI Container Registry in an appropriate format.

## Running locally

GoLang version 1.20 or higher is required to run the service

### Installing

From the project dir command line run the go build command to compile the code into an executable 
```
go build ./cmd/helm-registry-adapter
 ```

### Running

Now run the executable from the command line
```
./helm-registry-adapter
```

### Running the tests

From the project dir command line run the go test command
```
go test ./..
```

### Environment variables
Project default properties are stored in the ``.env`` file inside the project root filder. Run arguments or system environment variables can override default properties.
Each item takes precedence over the item below it:
* run args
* system environment variables
* .env config file