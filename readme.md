# 1️⃣ What JWT Authentication *is* (conceptually)

JWT (JSON Web Token) authentication is **stateless authentication**.

Instead of:

* server storing sessions in memory or DB ❌

We do:

* server **signs a token**
* client stores it
* client sends it back on every request

The server:

* **verifies** the token
* **trusts the claims inside**

---

## JWT structure (important)

A JWT has **3 parts**:

```
HEADER.PAYLOAD.SIGNATURE
```

Example (conceptually):

```json
HEADER:
{
  "alg": "HS256",
  "typ": "JWT"
}

PAYLOAD:
{
  "sub": "user_id_123",
  "role": "admin",
  "exp": 1735689600
}

SIGNATURE:
HMACSHA256(base64(header + payload), JWT_SECRET)
```

👉 If **payload is changed**, signature breaks
👉 Only server with `JWT_SECRET` can sign valid tokens

---

# 2️⃣ JWT in *our* Go project (from `go.mod`)

### This line is the core:

```go
github.com/golang-jwt/jwt/v5
```

This library:

* creates JWTs
* signs them
* validates them
* extracts claims safely

---

## Typical JWT flow in Gin

### 🔐 Login

1. User sends email + password
2. We verify credentials (MongoDB)
3. We generate JWT
4. We send token to client

### 🔁 Protected request

1. Client sends:

   ```
   Authorization: Bearer <token>
   ```
2. Middleware validates token
3. User info is attached to context
4. Handler runs

---

# 3️⃣ Why JWT is perfect for Gin

Gin is:

* fast
* stateless
* middleware-based

JWT fits **exactly** that model.

Gin doesn’t care *who* the user is —
middleware figures it out **before** handlers run.

---

# 4️⃣ JWT Middleware (how it works internally)

In Gin, JWT auth is always **middleware**.

Conceptual flow:

```text
Request
  ↓
JWT Middleware
  ↓
Token valid? ❌ → 401
  ↓
Extract claims
  ↓
Attach to context
  ↓
Next handler
```

### Typical middleware logic (simplified)

```go
tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

token, err := jwt.Parse(tokenStr, keyFunc)

claims := token.Claims.(jwt.MapClaims)
ctx.Set("userID", claims["sub"])
ctx.Set("role", claims["role"])
```

---

# 5️⃣ RBAC (Role-Based Access Control)

We said **RBA**, but what we’re actually using is **RBAC**.

RBAC = **authorization**, not authentication.

> Auth = *Who are we?*
> RBAC = *What are we allowed to do?*

---

## RBAC in JWT (the smart way)

We **embed roles in JWT claims**.

Example payload:

```json
{
  "sub": "user_id_123",
  "role": "admin"
}
```

Now:

* no DB lookup needed
* role is cryptographically trusted
* fast authorization

---

# 6️⃣ RBAC in Gin (real-world pattern)

### Step 1: JWT middleware (auth)

Sets user data:

```go
ctx.Set("role", "admin")
```

---

### Step 2: Role middleware (authorization)

Example:

```go
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")

		if userRole != role {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
```

---

### Step 3: Apply to routes

```go
admin := r.Group("/admin")
admin.Use(JWTAuth(), RequireRole("admin"))
{
	admin.POST("/create", createUser)
}
```

✔️ Auth first
✔️ Role check second
✔️ Handler last

---

# 7️⃣ MongoDB’s role in JWT auth (from our deps)

```go
go.mongodb.org/mongo-driver
```

MongoDB is used for:

* storing users
* storing hashed passwords
* storing roles (initially)

But **NOT** for session storage.

Once JWT is issued:

* MongoDB is no longer involved per request
* JWT replaces session DB lookups

---

# 8️⃣ Why validation libs are in our `go.mod`

```go
github.com/go-playground/validator/v10
```

Used by Gin for request validation:

```go
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
```

This prevents:

* malformed requests
* empty fields
* invalid email formats

Before JWT logic even runs.

---

# 9️⃣ Why JWT_SECRET exists (critical)

From `.env`:

```env
JWT_SECRET=change_this_to_a_long_random_secret
```

Used to:

* sign tokens
* verify tokens

Rules:

* must be long
* must be random
* must never change (or all tokens break)

---

# 🔐 Security guarantees we get

With JWT + RBAC done properly:

✅ No session storage
✅ No server-side auth state
✅ Horizontal scaling is trivial
✅ Each request is self-contained
✅ Roles cannot be forged
✅ MongoDB is not hit on every request

---

# 10️⃣ Why this architecture is CLEAN

We already followed best practices:

* **Config layer** → loads secrets
* **DB layer** → data access only
* **Middleware** → auth & RBAC
* **Handlers** → business logic
* **JWT** → stateless identity

This is **production-grade backend design**.

---

JWT + RBAC in **Gin + MongoDB**

---

# 🧠 PART 0 — The Problem We’re Solving

We’re building an API.

We need to answer **two different questions** for every request:

1. **Authentication**
   👉 *Who is making this request?*

2. **Authorization**
   👉 *Are they allowed to do this?*

JWT solves **authentication**
RBAC solves **authorization**

They work together, but they are **not the same thing**.

---

# 🔐 PART 1 — Authentication (JWT)

## What authentication really means

Authentication is **identity verification**.

> “Prove we are who we say we are.”

In our app, that happens **once**:

* during login

After that, we don’t want to keep asking:

> “What’s our password?”

---

## Old way (sessions) — why it sucks

Traditional approach:

* User logs in
* Server creates a session
* Session stored in memory / Redis / DB
* Browser stores session ID in cookie
* Server looks up session on every request

Problems:

* Stateful
* Hard to scale
* Extra DB lookup every request
* Breaks easily with multiple servers

---

## Modern way (JWT)

JWT flips the model:

* Server **creates a signed token**
* Client stores it
* Client sends it with every request
* Server **verifies signature**, not session

No storage. No lookup. No state.

---

# 🧱 PART 2 — What a JWT actually is (not magic)

A JWT is **just a string**.

Example:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
.
eyJzdWIiOiI2NWFiYzEyMyIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTczNTY4OTYwMH0
.
M8n9sM5lX2vKpZJZ6jC3...
```

Three parts:

```
HEADER.PAYLOAD.SIGNATURE
```

---

## 1️⃣ Header

Describes **how** token is signed.

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

---

## 2️⃣ Payload (claims)

This is the **identity data**.

```json
{
  "sub": "user_id_123",
  "role": "admin",
  "exp": 1735689600
}
```

Important:

* `sub` = subject (user ID)
* `role` = authorization data
* `exp` = expiration time

⚠️ Payload is **base64 encoded**, NOT encrypted
Anyone can read it, but no one can **modify it safely**

---

## 3️⃣ Signature (the security)

Signature =

```
HMAC(
  base64(header) + "." + base64(payload),
  JWT_SECRET
)
```

Only our server has `JWT_SECRET`.

So:

* If payload changes → signature breaks
* If secret is wrong → verification fails

That’s the entire security model.

---

# 🔑 PART 3 — Login Flow (step by step)

### Step 1: User logs in

Client sends:

```json
POST /login
{
  "email": "user@test.com",
  "password": "123456"
}
```

---

### Step 2: Server verifies credentials (MongoDB)

* Find user by email
* Compare password hash
* Get:

  * user ID
  * role

MongoDB is **only needed here**

---

### Step 3: Server creates JWT

Server creates payload:

```go
claims := jwt.MapClaims{
	"sub": user.ID,
	"role": user.Role,
	"exp": time.Now().Add(24 * time.Hour).Unix(),
}
```

Signs it using:

```go
JWT_SECRET
```

---

### Step 4: Token is returned

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

---

# 📦 PART 4 — Client Responsibility

Client must:

* store JWT (memory / localStorage)
* attach it to every request

Example:

```
Authorization: Bearer <token>
```

That’s it.

No cookies required.

---

# 🚪 PART 5 — Protected Routes (Gin Middleware)

This is where Gin shines.

## Why middleware?

Because:

* every protected route needs auth
* duplication is bad
* middleware runs **before handlers**

---

## JWT Middleware Logic (mental model)

1. Read `Authorization` header
2. Extract token
3. Verify signature
4. Check expiration
5. Extract claims
6. Attach claims to request context

---

### Why context?

Because:

* handlers need user info
* context is request-scoped
* safe and idiomatic

Example:

```go
c.Set("userID", claims["sub"])
c.Set("role", claims["role"])
```

---

# 🧾 PART 6 — Authorization (RBAC)

Authentication answers:

> “Who are we?”

Authorization answers:

> “What can we do?”

---

## RBAC = Role-Based Access Control

Instead of:

* checking permissions everywhere

We assign:

* roles

Example roles:

* `user`
* `admin`
* `moderator`

---

## Why roles go into JWT

Because:

* roles rarely change
* roles are security-sensitive
* JWT is signed → cannot be forged

No DB lookup needed per request.

---

# 🧠 PART 7 — RBAC Middleware (the clean way)

After JWT middleware runs, role is in context.

RBAC middleware checks it.

Example logic:

```go
if role != "admin" {
	403 Forbidden
}
```

This keeps:

* handlers clean
* security centralized
* logic reusable

---

## Route protection pattern

```go
r.GET("/profile", JWTAuth(), profileHandler)

admin := r.Group("/admin")
admin.Use(JWTAuth(), RequireRole("admin"))
```

Order matters:

1. Authenticate
2. Authorize
3. Execute handler

---

# 🗃️ PART 8 — Where MongoDB fits now

MongoDB is used for:

* users
* roles
* passwords
* refresh tokens (optional)

MongoDB is **NOT** used:

* on every request
* for auth state
* for sessions

JWT removed that need.

---

# 🔄 PART 9 — Token Expiry (important)

JWT must expire.

Why?

* stolen tokens
* user logout
* role changes

Common pattern:

* Access token: short (15m – 1h)
* Refresh token: long (7–30 days)

Refresh tokens **are stored in DB**.

---

# 🧯 PART 10 — Security Guarantees

With JWT + RBAC:

✅ Identity is verifiable
✅ Roles cannot be forged
✅ No server-side sessions
✅ Horizontal scaling is trivial
✅ DB load is reduced
✅ Clean separation of concerns

---

# 🧩 FINAL MENTAL MODEL (remember this)

```
Login → Verify → Sign Token → Return Token

Request → Verify Token → Extract Claims → Check Role → Handle Request
```

JWT = **identity container**
RBAC = **gatekeeper**

---

