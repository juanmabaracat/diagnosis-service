# diagnoses-api Challenge
Diagnoses-API handles the administration of patient diagnoses.

## Goal
A solution to collect diagnostic data from patients in our systems.
This solution should be offered as a service that can be connected to any type of application for integration.
Additionally, we need to expose two endpoints:
- A diagnostic query that allows filtering by patient name and/or date.
- A writing endpoint where diagnostic data for each patient will be created.
Regarding our data structure, we need to have patients, each with their names, identification documents, and contact information. These patients will have a relationship with their diagnostic data, which will include time and date, patient information, diagnosis, and an associated prescription that may or may not exist.

Diagnoses can occur multiple times a day, as sometimes a patient may be hospitalized and monitored several times a day.

OpenAPI standard.
At least some unit tests.
Dockerize the solution.
A Postman file containing the prepared endpoints, Swagger integrated into the project for API querying from an endpoint, or a README.md with instructions on how to set up or configure your project.

# Solution
I based my solution on the concepts of clean architecture and Domain-Driven Design; that is why you can see three main
layers: domain, application, and infrastructure.

Some points that I left unattended or could be improved include:
- Input Sanitization: This could be easily added in the handler by checking formats, empty values, etc.
- Error handling: I implemented basic error handling, but there is room for improvement.
- Security: I used a default server configuration and added some middleware, but this should be revised.

Considering the context of this challenge and the time constraints, I decided to keep the domain simple, 
and I also left out some requirements such as the filter by date. Additionally, I opted to implement a quick memory repository. 
To run some examples, you can find a Postman collection and a Swagger UI. Further details can be found in the documentation section.


### Technical details
- [Go 1.22.0](https://go.dev/): Go version 1.22.0
- [go-chi](https://github.com/go-chi/chi): lightweight router
- [swaggo](https://github.com/swaggo/http-swagger): wrapper to generate RESTful API documentation with Swagger 2.0
- [Testify](https://github.com/stretchr/testify): testing tool

### How to run the application:
#### Running locally with Go installed
Clone the repository and move to the project folder:
```
git clone https://github.com/juanmabaracat/diagnosis-service.git
cd diagnosis-service
```
Run the application:
```
go run cmd/main.go
```
Run all tests (root folder):
```
go test ./...
```

#### Using docker to build and run the application:
```
$docker build -t diagnoses-api .
$docker run -p 8080:8080 diagnoses-api:latest
```

### Documentation
#### Postman Collection with examples
In docs folder there is a postman collection with two examples, one to get the diagnoses for an already created example user
and the other to create diagnoses for that user.

#### Swagger UI
I've used a tool to generate a swagger UI documentation.
After you start the server, open this url:
[open swagger](http://localhost:8080/swagger/index.html)
