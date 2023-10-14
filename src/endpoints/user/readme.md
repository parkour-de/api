# User Authentication API Guide
Welcome to the User Authentication API! This API provides endpoints for user registration, login, and authentication management. Here's how you can use it:

## User Registration
- `GET /api/users/{key}/exist`: Check if a username is already taken (boolean response).
- `POST /api/users/{key}/claim`: Claim a username that is taken but has never been used.
- `POST /api/users/{key}/create`: Create a new user and receive a 30-minute session token.
- Valid usernames must match this regex: `^[a-zA-Z0-9_-][a-zA-Z0-9_\-.]{2,29}$`
  - The first character must be alphanumeric (a-z, A-Z, 0-9), an underscore, or a hyphen.
  - The remaining characters must be alphanumeric, an underscore, a hyphen, or a period.
  - The username must be between 3 and 30 characters long.
  - The username must not only contain digits.
  - If an empty username is provided, a random sequence of digits will be used.
    And you'll have to memorize it.
## Authentication Methods
**Choose from Various Methods:** You can set up various authentication methods:
- **Facebook:** `GET /api/users/{key}/facebook?auth={token}`
- **Google:** `GET /api/users/{key}/google?auth={token}`
- **Password:** `GET /api/users/{key}/password?password={password}`
- **TOTP (Time-Based One-Time Password):**
  - Get TOTP Setup Info: `GET /api/users/{key}/totp`
  - Enable TOTP: `POST /api/users/{key}/totp` (with JSON containing the first code)
- **Email:**
  - Request Confirmation Email: `GET /api/users/{key}/email?email={email}`
  - Confirm Email: `GET /api/users/{key}/email/{id}?code={code}`
## Token Renewal
**Renew Token:** When your session token is about to expire, use these endpoints to get a new token:
- Facebook: `GET /api/facebook`
- Google: `GET /api/google`
- Email Login: `GET /api/users/{key}/login?code={code_from_email}`
- Password Login: `POST /api/users/{key}/login` (with password in the body)
## Additional Actions (Coming Soon)
In the future, you'll be able to list or delete login methods:
- List Login Methods: `GET /api/users/{key}/logins`
- Delete Login Method: `DELETE /api/users/{key}/logins/{id}`
Please note that for security reasons, you may not be able to delete the most recently used login method.

Enjoy using our User Authentication API for your application!
