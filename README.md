# WaterVapor - A gooder game service

A very original full stack game distribution platform built with Go, PostgreSQL, and JavaScript for CS3380.

## Project Overview

WaterVapor is a very new, never before seen concept: A digital game purchasing platform featuring user accounts, game purchases, family sharing, publisher management, wishlists, shopping cart, and refund systems.

## Services

- **PostgreSQL**: Persistent database storage with normalized schemas
- **Go Webapp**: RESTful API backend with JWT authentication
- **Nginx**: Reverse proxy serving frontend and routing API requests
- **pgAdmin**: Web-based database management interface

## Database Schema

### Third Normal Form (3NF)

Our relation schemas satisfy Third Normal Form (3NF) requirements:

**1NF Compliance:**
- All tables have atomic values (no multi-valued attributes)
- Each table has a primary key (ID fields)
- No repeating groups

**2NF Compliance:**
- All non-key attributes are fully functionally dependent on the primary key
- No partial dependencies exist

**3NF Compliance:**
- No transitive dependencies
- All non-key attributes depend only on the primary key

**Example Tables:**
- `users`: email, username, password_hash depend only on user_id
- `games`: name, price, currency, genre depend only on game_id and publisher_id (foreign key)
- `purchases`: street, city, state, zip_code, country depend on purchase_id (not on user_id or game_id)
- `families`: Separate entity eliminates transitive dependency (users don't depend on other users' data)

## Getting Started

### Database Setup and Population

1. **Clone the repository:**
   ```sh
   git clone https://github.com/bradleymt/CS3380-Project.git
   cd CS3380-Project
   ```

2. **Configure environment:**
   ```sh
   # .env file is pre-configured with:
   POSTGRES_USER=user
   POSTGRES_PASSWORD=password
   POSTGRES_DB=mydb
   ```

3. **Build and start services:**
   ```sh
   docker compose up --build
   ```

4. **Database tables are automatically created** on first run via GORM AutoMigrate in `webapp/database/database.go`

5. **Populate sample data using provided scripts:**
   ```powershell
   # Register users
   .\webapp\scripts\registerUser.sh alice alice@example.com password123
   .\webapp\scripts\registerUser.sh bob bob@example.com password123
   
   # Create publisher
   .\webapp\scripts\createPublisher.sh <auth-token> "WaterVapor Studios" "USA"
   
   # Create games
   .\webapp\scripts\createGame.sh <auth-token> "Epic RPG" 59.99 USD RPG
   
   # Add to cart and purchase
   .\webapp\scripts\addToCart.sh <auth-token> 1
   .\webapp\scripts\purchase.sh <auth-token>
   ```

## Application Features

### User Authentication & Game Library

**Feature:** Secure user registration, login, and personal game library

**Demonstration:**
1. Navigate to `http://localhost:8081/login.html`
2. Register new user with email, username, and password
3. Backend hashes password using bcrypt before storage
4. Login returns JWT token stored in browser
5. View purchased games in Library page
6. Family shared games appear separately (if in a family)

**Database Operations:**
- INSERT into `users` table with hashed password
- SELECT user credentials for authentication
- JOIN `purchases` and `games` tables for library display
- JOIN `users` and `families` for family game sharing

### Publisher System & Game Creation

**Feature:** Publishers can create and manage games

**Demonstration:**
1. Navigate to Publisher Dashboard
2. Create new publisher with studio name and country
3. Create games with name, price, currency, and genre
4. View all games published by your studio
5. Remove games from the platform

**Database Operations:**
- INSERT into `publishers` table
- UPDATE `users` SET publisher_id
- INSERT into `games` with publisher_id foreign key
- SELECT games WHERE publisher_id = current_publisher
- DELETE from `games` (cascade handled by foreign keys)

### Family Sharing & Shopping Cart

**Feature:** Family game sharing and shopping cart with discounts

**Demonstration:**
1. Create family from Profile page
2. Add other users to family by username
3. Family members can access each other's purchased games
4. Add games to cart from Store
5. Apply discounts at checkout
6. Purchase games with address and payment method

**Database Operations:**
- INSERT into `families` table
- UPDATE `users` SET family_id
- SELECT users WHERE family_id = current_family
- INSERT/DELETE from `cart_items` table
- JOIN `purchases`, `games`, `users` for family library
- INSERT into `purchases` with address details
- SELECT discounts WHERE game_id for dynamic pricing

### Refund System & Admin Management

**Feature:** Users can request refunds and admins can approve them

**Demonstration:**
1. Navigate to Library and request refund on any purchased game
2. Refunds page displays all pending and approved refund requests
3. Admin users see additional section with all pending refunds from all users
4. Admin can approve refunds which removes the game from user's library
5. Purchase and library records are permanently deleted on approval

**Database Operations:**
- INSERT into `refunds` table with purchase_id and reason
- SELECT refunds with JOIN to games for game names
- SELECT refunds with JOIN to users and games for admin view
- UPDATE refunds SET approved = true
- DELETE from `library_items` WHERE user_id and game_id
- DELETE from `purchases` WHERE purchase_id
- DELETE from `refunds` after processing

## Usage

- **Frontend**: `http://localhost:8081`
- **API**: `http://localhost:8081/api/`
- **pgAdmin**: `http://localhost:8080`
  - Email: admin@admin.com
  - Password: admin

## Key Features

- User authentication with bcrypt password hashing
- JWT token-based session management
- Publisher account system
- Game creation and management
- Shopping cart with real-time discount calculation
- Wishlist functionality
- Family sharing system (share purchased games)
- Purchase history and refund requests
- Admin refund approval system
- Dynamic game discounts
- Responsive UI design

## Project Structure

```
CS3380-Project/
├── webapp/
│   ├── database/        # GORM models and DB connection
│   ├── routes/          # API endpoints
│   ├── middleware/      # JWT authentication
│   └── main.go          # Application entry point
├── frontend/            # HTML, CSS, JavaScript
├── nginx/              # Reverse proxy configuration
├── docker-compose.yml  # Service orchestration
└── .env               # Environment variables
```

## Database Tables

- `users` - User accounts with authentication
- `publishers` - Game publisher studios
- `families` - Family sharing groups
- `games` - Game catalog with pricing
- `discounts` - Promotional discounts
- `purchases` - Transaction records
- `refunds` - Refund requests
- `library_items` - User game ownership
- `wishlist_items` - User wishlists
- `cart_items` - Shopping cart contents

## Made By
- Bradley Thornton ([@bradleymt](https://github.com/bradleymt))
- Samuel Wiseman ([@notSam25](https://github.com/notSam25))
