# Movie App Endpoint

#### Supported endpoints and actions
| Method | URL Pattern | Handler | Action
| --- | --- | --- | --- |
| GET | `/v1/healthcheck` | healthcheckhandler | Show application information
| GET | `/v1/movies` | showAllMovieHander | Show the details of all movies
| POST | `/v1/movies` | createMovieHandler | Create a new movie
| GET | `/v1/movies/:id` | showMovieHandler | Show the details of a specific movie
| PATCH | `/v1/movies/:id` | patchMovieHandler | Update the details of a specific movie
| DELETE | `/v1/movies/:id` | deleteMovieHandler | Delete a specific movie
| POST | `/v1/users` | createUserHandler | Register a new user
| PUT | `/v1/users/activated` | activateUserHandler | Activate a specific user
| PUT | `/v1/users/password` | updatePasswordHandler | Update the password for a specific user
| POST | `/v1/tokens/authentication` | authenticateUserHandler | Generate a new authentication token
| POST | `/v1/tokens/password-reset` | resetUserPasswordHandler | Generate a new password-reset token
| GET | `/debug/vars` | displayMetricsHandler | Display application metrics