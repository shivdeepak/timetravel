# The Rainbow Take-Home Assignment

Please create a __private__ version of this repo, complete the objectives, and once you
are finished, send a link to your repo to us.

# The Assignment

Part of what an insurance company needs to have in its backend is a 
record system. As an insurer, we need to keep an up-to-date record of each of our policy-holder's
data points that go into the calculation of their rate. When a policy-holder updates
their information, I.E. they change addresses, or add/remove new employees to their team
we will be notified and we must keep our records up to date.

The current version of the repo is an extremely simplified version of exactly that. `GET /api/v1/record/{id}`
will retrieve a record, which is just a json mapping strings to strings. and `POST /api/v1/record/{id}`
will either create a new record or modify an existing record. However, it isn't enough to
just keep a record of the current record state but we must maintain a reference to how the state
has changed to be in full compliance.

Say that the policy-holder buys their insurance on the start of the year, and then two months later
changes the address of their business but doesn't tell us about this change until 4 months after that.
Since we were technically held liable if there was a claim event, we need to charge the customer the
difference for the 4 months since they changed addresses. To do so accurately, we need to know the
version of the records that we knew about them at the two points of time: at the time when the change happened
and at the time when we were told of the change.

In this project, you'll make a simplified version of this system. We've implemented an in-memory key-value store with no history. 
At a high-level your goal is to do two things to this existing codebase:

1. Change the storage backend to sqlite, and persist the data across turning off and on the server.
2. Add the time travel component so we can easily look up the state of each records at different timesteps.

The sections below outline these two objectives in more detail. You may use whatever libraries and tools
you like to achieve this even as far as building this in an entirely different language.

## Objective: Switch To Sqlite

The current implementation does not store the data. The data is lost once the server 
process is killed. You should change the code so that all changes are persisted on 
to sqlite.

Once you're done, the data should be persistent on to a sqlite file as the server 
is running. The server should tolerate restarting the process without data loss.

## Objective: Add Time Travel
This part is far more open-ended. You might need to make major changes across nearly
all files of the codebase. You'll be adding persistentence to the records. 

You should create a set of `/api/v2` endpoints that enable you to do run gets, creates, and updates. 
Unlike in v1, records are now versioned. Full requirements: 

- You should have endpoints that allow the api client to get records at different versions. (not just 
the latest version). 
- You should be able to add modifications on top of the latest version. 
- There should be a way to get a list of the different versions too.
- `/api/v1` should still work after these changes with identical behavior as before.

# Reccommendations

We expect you to work as if this task was a normal project at work. So please write
your code in a way that fits your intuitive notion of operating within best practices.
Additionally, you should at the very least have a different commmit for each individual objective, 
ideally more as you go through process of completing the take-home. Also we like
to see your thought process and fixes as you make changes. So don't be afraid of
committing code that you later edit. No need to squash those commits.

Many parts of the assignment is intentionally ambiguious. If you have a question, definitely
reach out. But for many of these ambiguiuties, we want to see how you independently make
software design decisions.

# FAQ
_Can I Use Another Language?_
Definitely, we've had multiple people complete this assignment in Python and Java. You can pick whatever
language you'd like although you should aim to replicate the functionality in the boilerplate. 

_Did you really end up implementing something like this at Rainbow?_
Yes, but unfortunately it wasn't as simple as this in practice. For insurance a number of requirements force us 
to maintain historic records across many different object types. So in fact we implemented this across multiple different 
tables in our database. 


# Reference -- The Current API

There are only two API endpoints `GET /api/v1/records/{id}` and `POST /api/v1/records/{id}`, all ids must be positive integers.

### `GET /api/v1/records/{id}`

This endpoint will return the record if it exists.

```bash
> GET /api/v1/records/2323 HTTP/1.1

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
{"id":2323,"data":{"david":"hey","davidx":"hey"}}
```

```bash
> GET /api/v1/records/32 HTTP/1.1

< HTTP/1.1 400 Bad Request
< Content-Type: application/json; charset=utf-8
{"error":"record of id 32 does not exist"}
```

### `POST /api/v1/records/{id}`

This endpoint will create a record if a does not exists.
Otherwise it will update the record.

The payload is a json object mapping strings to strings
and nulls. Values that are null indicate that the
backend must delete that key of the record.

```bash
# Creating a record
> POST /api/v1/records/1 HTTP/1.1
{"hello":"world"}

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
{"id":1,"data":{"hello":"world"}}


# Updating that record
> POST /api/v1/records/1 HTTP/1.1
{"hello":"world 2","status":"ok"}

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
{"id":1,"data":{"hello":"world 2","status":"ok"}}


# Deleting a field
> POST /api/v1/records/1 HTTP/1.1
{"hello":null}

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
{"id":1,"data":{"status":"ok"}}
```

# API V2 - Reference

There are only three API endpoints:

1. `GET /api/v2/records/{id}`
2. `POST /api/v2/records/{id}`
2. `GET /api/v2/records/{id}/versions`

all ids must be positive integers.

### `GET /api/v2/records/{id}`

This endpoint will return the record if it exists.

```bash
> GET /api/v2/records/2323 HTTP/1.1

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
{"id":2323,"data":{"david":"hey","davidx":"hey"}}
```

```bash
> GET /api/v2/records/32 HTTP/1.1

< HTTP/1.1 400 Bad Request
< Content-Type: application/json; charset=utf-8
{"error":"record of id 32 does not exist"}
```

This endpoint also has supports time travel, meaning you can lookup
a previous version of the record at a given point in time, if it exists.

```bash
> GET /api/v2/records/30?at=2024-08-25T16:25:00-07:00 HTTP/1.1

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
{"id":30,"data":{"created_at":"2024-08-25T16:13:02-07:00","dob":"1955-02-24T00:00:00-07:00","first_name":"Steven","last_name":"Jobs","middle_name":"","updated_at":"2024-08-25T16:19:18-07:00"}}
```

### `POST /api/v2/records/{id}`

This endpoint will create a record if a does not exists.
Otherwise it will update the record.

The payload is a json object mapping strings to strings
and nulls. Values that are null indicate that the
backend must delete that key of the record.

Note that it will only update the record if there are changes.
Otherwise it will just return a status code 200 without any
changes, i.e. it will not create a new version.

```bash
# Creating a record
> POST /api/v2/records/1 HTTP/1.1
{"hello":"world"}

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
{"id":1,"data":{"hello":"world"}}


# Updating that record
> POST /api/v2/records/1 HTTP/1.1
{"hello":"world 2","status":"ok"}

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
{"id":1,"data":{"hello":"world 2","status":"ok"}}


# Deleting a field
> POST /api/v2/records/1 HTTP/1.1
{"hello":null}

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
{"id":1,"data":{"status":"ok"}}
```

### `GET /api/v2/records/{id}/versions`

This endpoint provides ability to lookup all versions of a
records if it exists in the database.


```bash
> GET /api/v2/records/30/versions HTTP/1.1

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
[{"id":30,"data":{"created_at":"2024-08-25T16:13:02-07:00","dob":"1955-02-24T00:00:00-07:00","first_name":"Steven","last_name":"Jobs","middle_name":"","updated_at":"2024-08-25T16:19:18-07:00"}},{"id":30,"data":{"created_at":"2024-08-25T16:13:02-07:00","dob":"1955-02-24T00:00:00-07:00","first_name":"Steven","last_name":"Jobs","middle_name":"Paul","updated_at":"2024-08-25T16:14:01-07:00"}},{"id":30,"data":{"created_at":"2024-08-25T16:13:02-07:00","dob":"0001-01-01T00:00:00Z","first_name":"Steven","last_name":"Jobs","middle_name":"Paul","updated_at":"2024-08-25T16:13:21-07:00"}},{"id":30,"data":{"created_at":"2024-08-25T16:13:02-07:00","dob":"0001-01-01T00:00:00Z","first_name":"Steve","last_name":"Jobs","middle_name":"","updated_at":"2024-08-25T16:13:02-07:00"}}]
```

```bash
> GET /api/v2/records/32/versions HTTP/1.1

< HTTP/1.1 400 Bad Request
< Content-Type: application/json; charset=utf-8
{"error":"record of id 32 does not exist"}
```

# Further Improvements

### Record Versions and Audit Trail

My current implementation of versioned record updates is quite rudimentary.

I can think of making following enhancements based on requirements and needs:

1. We can store the exact changes during each update to the records table. This
   will be useful if we need to build a user interface that shows an audit
   trail of all changes that were made to a given record.

2. On top of `when` and `what` values have changed, it is also important to
   capture `who` changed the record. It could be easily implemented by capturing
   the current session user_id (when an authenticated session is available /
   implemented).

4. In the current project, there is only one resource, and we call it `record`.
   In a production level application, there are several different resources,
   and it is beneficial to create a generic system that creates and stores
   the audit trail of all the changes in the system.

3. It may be desirable to store the Audit Trail in a seperate table or even
   database for compliance reasons.

Taking all the above points into consideration, a more sophisticated audit trail
record/entry might look like this:

```javascript
{
	version_id: "12e72bb0-8ce6-4017-81d1-0dbf1c0711e5",
	entity_type: "rainbow_records", // resource name, `record` in current project
	entity_id: 30,
	actor_type: "rainbow_user", // could be admin, customer, or user of the customer
	actor_id: 3092,
        changes: [
		{type: "update", field: "first_name", before: "Steve", after: "Steven"}
		{type: "delete", field: "dob", before: "1955-02-24T00:00:00-07:00", after: null},
		{type: "create", field: "middle_name", before: null, after: "Paul"},
	],
	created_at: "2024-08-25T13:44:29-07:00"
}
```

Note that the exact data format has to also consider how other systems are going
to consume the data.

### Safer Database Schema Migrations

Here I am using [GORM's](https://gorm.io/) in-built
[auto migrate](https://gorm.io/docs/migration.html) functionality. Which doesn't
necessarily do the right thing during complex migration. I would worry that it
might even lead to data loss.

I would use a better migration mechanism, for example:
[gormigrate](https://github.com/go-gormigrate/gormigrate) that give you more
fine-grained control and visibility over the schema migrations will be running 
on the production database server.

### API Spec & API Docs

If this API is used by many users (such as other internal or external engineers
or data scientists), we can consider defining a standard API with documentation
and providing SDKs in popular programming languages.

I would use [Open API](https://spec.openapis.org/oas/latest.html) spec to build
API Documentation. It can also be used to auto-generate clients (SDKs) in
various programming languages. Open API leverages API specification defined in
a standardized `.json` file that is used for documentation, and SDK generation.

I would also consider using
[json:api](https://spec.openapis.org/oas/latest.html) spec to define the
requests' payload and responses for /api/v2. JSON:API is defined to allow an API
spec to evolve with more capabilities while being backward compatible. This will
allow us to stay with a single version for a very long time, reducing
irrevesible technical debt. Also, the attribute names use a convention that
developers can remember.

We currently don't use JSON API in /api/v2, but that could be implemented if
desired before releasing the API.

### Make it a Twelve-Factor App - [https://12factor.net/](https://12factor.net/)

It is desirable to build a 12-factor app for production which has numerous
benefits from developer experience, maintainability, to application security.

What can be improved:

1. Config: We can improve how we handle config. For instance, right now I have
   hardcoded db file location. In cloud based deployments with possibly multiple
   webservers hosting the API, it is desirable to use server based relational
   database, for example: MySQL, Postgres etc. In such cases, the database
   credentials are considered confidential and there could be other secrets and
   environment-specific configurations. Adding them through environment
   variables is more desirable. We can use
   [`godotenv`](https://github.com/joho/godotenv) library to manage environment
   variables by leveraging `.env*` files in local development and
   possibly even in production, depending on the deployment setup.
2. Build, release, run: For a production-grade application with CI/CD, we
   can think about separate build, release, and run stages. However, given the
   project's limited scope, this is not applicable.
3. Processes, Concurrency, Disposability: Ideally, a web server process should
   be stateless, but in our case, it's not because we are using SQLite which is
   a file based database.

   This can be fixed by using an server based database system such as Postgres
   or MySQL. When we do that, web server processes will become stateless and
   can be run in parallel without degrading the rate of cache hits, and we can
   achieve concurrency if desired.

   A stateless server is disposable.
4. Logs, Admin processes: We don't have any of this setup due to the project's
   limited scope. If we are to run this as a production-grade service, then
   thinking about Logging and Admin Processes is ideal.

### Tests

1. Given the limited time and scope of the project, I have not written any unit
   or integration tests. If this service were to run in production, we will
   atleast require API level integration testing. This will help us avoid
   regression issues in production.
2. API level integration testing is great to get started. But for better
   developer experience and productivity, I will also write individual module
   level unit tests so that developers can make significant changes confidently,
   and release feature quickly.
