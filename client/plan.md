1. Executive Summary
This detailed architectural report provides a comprehensive blueprint for the development of a high-performance, distributed, and cryptographically secure Poker platform. Commissioned to address the technical requirements of a modern real-time gaming environment, the system is architected using the Go programming language to leverage its superior concurrency primitives, the Gin framework for robust RESTful API delivery, and the Melody library for efficient WebSocket orchestration. The platform is designed to support the complete lifecycle of Texas Hold'em gameplay, from table creation and player seating to complex hand evaluations and financial settlement, all while maintaining a verifiable chain of custody for game integrity through Provably Fair algorithms.

The online gambling and social gaming sectors demand rigor in three critical dimensions: latency, integrity, and scalability. This report addresses these dimensions by proposing a Three-Layer Clean Architecture that ensures the decoupling of transport, logic, and persistence layers, thereby facilitating testability and future maintenance. Data persistence strategies utilize a hybrid model, employing PostgreSQL for immutable financial ledgers and relational user data, alongside Redis for sub-millisecond access to ephemeral game states and pub/sub event distribution.

Furthermore, this document serves as an exhaustive guide for the engineering team, expanding the initial development roadmap into a granular technical specification. It addresses sophisticated domain challenges such as concurrency control during betting rounds, the algorithmic implementation of multi-way side pots, the mathematical verification of deck shuffles via HMAC-SHA256, and the integration of automated end-to-end testing environments using testcontainers-go. By strictly adhering to these specifications, the resulting platform will not only meet the functional requirements of the stakeholder but also stand as a paragon of reliability and fairness in the digital gaming market.

2. Architectural Philosophy and Design Patterns
2.1. The Strategic Imperative of Clean Architecture in Go
The adoption of a Three-Layer Architecture, often synonymous with Clean Architecture or Hexagonal Architecture in the Go ecosystem, is not merely a stylistic choice but a strategic necessity for long-term project health. In highly interactive applications like poker servers, where business logic is complex and state-dependent, the entanglement of HTTP transport concerns with core game mechanics can lead to a fragile codebase known as the "Big Ball of Mud."

The proposed architecture enforces a strict dependency rule: dependencies point inwards. The core business logic—the "Domain"—knows nothing of the external world. It defines interfaces that the outer layers must satisfy.

2.1.1. The Delivery Layer (Interface Adapters)
The outermost layer, the Delivery Layer, acts as the primary interface for external actors. In this specific implementation, it encompasses the HTTP handlers managed by the Gin framework and the WebSocket event processors managed by Melody. Its responsibilities are strictly limited to protocol translation: parsing incoming JSON payloads, extracting parameters from URLs, reading authentication tokens from headers or query strings, and validating these inputs against schematic constraints. Crucially, this layer does not make business decisions. It delegates validated requests to the Service Layer and marshals the resulting domain objects back into the appropriate transfer format (e.g., JSON) for the client.   

2.1.2. The Service Layer (Use Cases)
Situated at the core of the application, the Service Layer encapsulates the application specific business rules. This is where the poker engine resides. It answers questions such as: "Is it this player's turn?", "Does the player have sufficient funds to call?", and "What is the strength of this hand?". This layer is pure Go; it does not import gin, sqlx, or redis. Instead, it relies on interfaces defined within the Domain layer to interact with databases or other external services. This isolation allows for the game logic to be unit-tested in isolation, using mocks for data persistence, ensuring that the complex rules of Poker are verified without the overhead of database connectivity.   

2.1.3. The Repository Layer (Infrastructure)
The Repository Layer provides the concrete implementations of the interfaces defined by the Domain. It is responsible for the translation of domain objects into database rows (PostgreSQL) or key-value pairs (Redis). By isolating the database logic here, the architecture supports future adaptability—swapping the underlying database technology or modifying the schema requires changes only within this layer, leaving the core business logic untouched. This layer handles the low-level details of SQL queries, transaction management, and connection pooling.   

2.2. Project Directory Layout and Organization
To support this architecture, the project structure adheres to the "Standard Go Project Layout," a widely accepted convention that promotes consistency and discoverability across Go projects.   

Plaintext
/poker-backend
├── cmd
│   └── server
│       └── main.go           # Application entry point; dependency injection wiring
├── internal
│   ├── core
│   │   ├── domain            # Core entity definitions (User, Table, Hand, Wallet)
│   │   └── ports             # Interface definitions (Repository, Service)
│   ├── service               # Business logic implementation (Poker Engine, Auth)
│   │   ├── game              # Game loop, state machine, side pot logic
│   │   └── auth              # JWT generation, password hashing
│   ├── handler               # Transport logic
│   │   ├── http              # REST endpoints (Gin controllers)
│   │   └── ws                # WebSocket event handlers (Melody)
│   └── repository            # Data access implementations
│       ├── postgres          # SQL implementations using pgx or sqlx
│       └── redis             # Redis implementations for cache and pub/sub
├── pkg
│   ├── crypto                # Provably Fair RNG logic (HMAC-SHA256)
│   ├── logger                # Structured logging wrapper (Zap/Logrus)
│   └── validator             # Input validation logic
├── migrations                # SQL database migration files (up/down)
├── templates                 # Server-side HTML templates (Landing/Lobby)
├── docker-compose.yml        # Infrastructure orchestration
├── Makefile                  # Build, test, and run commands
└── go.mod                    # Module definition
The separation of internal and pkg is deliberate. internal packages are accessible only to the project itself, preventing external imports of the core business logic, while pkg contains reusable libraries (like the RNG crypto package) that could theoretically be used by other applications.   

2.3. Concurrency Model and The Actor Pattern
Go’s concurrency model, built on goroutines and channels, is uniquely suited for real-time game servers. A naive approach might lock the entire game state mutex for every operation, but this creates a bottleneck under load. Instead, this architecture proposes a variation of the Actor Model.

Each active Poker Table functions as an independent "Hub" running in its own goroutine. This Hub owns the state of that specific table absolutely. It listens on a channel for events (e.g., PlayerAction, TimerTick). When a WebSocket handler receives a message, it does not modify the table state directly. Instead, it sends a message to the Table Hub's channel. The Hub processes messages sequentially, eliminating the need for complex locking mechanisms within the game logic itself and ensuring that state transitions are atomic and race-condition free.   

3. Data Persistence Strategy
The system utilizes a polyglot persistence strategy, selecting the optimal storage engine for the distinct data lifecycles inherent in online gaming: PostgreSQL for durable, transactional data, and Redis for high-frequency, ephemeral state.   

3.1. PostgreSQL Schema Design (Durable Storage)
PostgreSQL is chosen for its strict ACID compliance, which is non-negotiable for financial transactions and user identity management. The schema design prioritizes data integrity through foreign keys and constraints.

3.1.1. User Identity and Authentication
The users table acts as the root of the identity graph. It stores immutable credentials and profile information.

SQL
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
Insight: The use of UUID prevents enumeration attacks common with sequential integer IDs.

3.1.2. Financial Ledger: Wallets and Transactions
The financial subsystem uses a double-entry accounting pattern. The wallets table holds the current balance, but the transactions table records every modification.

SQL
CREATE TABLE wallets (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    balance DECIMAL(15, 4) NOT NULL DEFAULT 0.0000 CHECK (balance >= 0),
    version INT NOT NULL DEFAULT 1, -- For optimistic locking
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID REFERENCES wallets(user_id),
    amount DECIMAL(15, 4) NOT NULL, -- Positive for credit, negative for debit
    reference_type VARCHAR(50) NOT NULL, -- e.g., 'buy_in', 'payout', 'deposit'
    reference_id UUID NOT NULL, -- ID of the game session or payment intent
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
Insight: The version column in wallets is critical for Optimistic Locking. When updating a balance, the query must ensure the version matches the read version (UPDATE wallets SET balance =?, version = version + 1 WHERE user_id =? AND version =?). This prevents race conditions where two concurrent game processes might try to deduct funds simultaneously.   

3.1.3. Game Sessions and Configuration
SQL
CREATE TABLE game_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    small_blind DECIMAL(10, 2) NOT NULL,
    big_blind DECIMAL(10, 2) NOT NULL,
    min_buy_in DECIMAL(10, 2) NOT NULL,
    max_buy_in DECIMAL(10, 2) NOT NULL,
    max_players INT NOT NULL CHECK (max_players BETWEEN 2 AND 9),
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'closed', 'paused'
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW()
);
3.1.4. Historical Integrity: Hand Histories
To support the "Provably Fair" verification, every hand must be archived with its cryptographic seeds.

SQL
CREATE TABLE hand_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES game_sessions(id),
    hand_number INT NOT NULL,
    server_seed_hash TEXT NOT NULL,
    server_seed_revealed TEXT, -- Populated after hand completion
    client_seed TEXT NOT NULL,
    nonce INT NOT NULL,
    deck_order TEXT NOT NULL, -- The shuffled deck
    result_json JSONB NOT NULL, -- Full replay data
    completed_at TIMESTAMP DEFAULT NOW()
);
Insight: Storing the server_seed_revealed allows users to audit the fairness of past hands at any time. The JSONB column enables flexible storage of the complex game events (bets, folds, chat) without rigid schema constraints.   

3.2. Redis Data Structures (Real-Time State)
Redis is leveraged not just as a cache, but as a primary store for the active game state. Accessing Postgres for every "Check" or "Fold" action would introduce unacceptable latency and database load.

Game State Hash: game:table:{table_id}

This Hash stores atomic fields: current_player_index, pot_size, board_cards, game_stage (e.g., "RIVER"), and action_timer_expiry.

Using a Hash allows HGET and HSET operations to modify specific fields (e.g., updating the pot) without reading/writing the entire object.   

Player Session Hash: game:table:{table_id}:player:{user_id}

Stores seat_index, chips_on_table, status (Active, SittingOut, AllIn), and encrypted hole_cards.

Isolating player state allows for efficient querying of "who is at seat X".

Pub/Sub Channels: channel:table:{table_id}

Used for broadcasting events. When the Game Engine updates the state (e.g., deals a card), it publishes a payload to this channel. The WebSocket layer subscribes to this and forwards the message to the appropriate connected clients.   

4. Network and Communication Layer
The communication layer is the bridge between the client (web browser/mobile app) and the server. It must handle authentication, connection upgrades, and real-time message routing.

4.1. REST API with Gin
The Gin framework is selected for its performance and minimalist API. It handles the initial setup phases of the user journey.

Authentication Middleware with JWT: Security is enforced via JSON Web Tokens (JWT). The middleware intercepts every request, validating the Authorization: Bearer <token> header. It parses the token claims to extract the UserID. If valid, the UserID is injected into the Gin Context for downstream handlers. If invalid, the request is aborted with a 401 Unauthorized status.   

Endpoint Implementation Strategy:

POST /api/v1/poker/tables: Accepts JSON payload defining table stakes. Validates that blinds are positive integers. Creates the GameSession in Postgres.

POST /api/v1/poker/tables/:id/join: Crucial transactional endpoint. It must:

Verify the user has sufficient wallet balance.

Check if the table is full (using Redis count).

Execute a transaction: Debit Wallet -> Credit Table Escrow (conceptually) -> Add Player to Redis State.

Broadcast a "Player Joined" event via Redis Pub/Sub.

4.2. WebSocket Orchestration with Melody
WebSockets provide the persistent, bidirectional channel required for gameplay. Melody wraps the lower-level gorilla/websocket library, abstracting away the boilerplate of connection maintenance, ping/pong heartbeats, and buffer management.   

4.2.1. The Authentication Challenge
Standard WebSocket APIs in browsers do not allow setting custom headers (like Authorization) during the handshake. Therefore, the architecture dictates passing the JWT via a query parameter: ws://api.poker.com/ws?token=eyJ.... The Gin handler responsible for the upgrade must:

Extract the token from the query string.

Validate the JWT signature and expiration.

Reject the connection immediately if invalid.

If valid, upgrade the connection and attach the UserID to the Melody session instance (session.Set("user_id", id)).

4.2.2. Room-Based Broadcasting
Melody does not have a built-in concept of "rooms" in the same way Socket.IO does. Instead, we implement rooms using Broadcast Filters or by maintaining explicit session maps.   

Approach: When a user joins Table A, we set a session key: session.Set("table_id", "TableA").

Broadcasting: To send a message to Table A, we invoke m.BroadcastFilter(msg, func(q *melody.Session) bool { return q.MustGet("table_id") == "TableA" }).

Optimization: For high-throughput scenarios, iterating through all sessions in a filter is inefficient. A scalable alternative (implemented in Phase 4) involves a Map<TableID, List<Session>> managed by the Game Hub, ensuring O(1) access to the list of recipients.

5. Core Domain: The Poker Engine
The Poker Engine is the most complex component, encapsulating the rules of Texas Hold'em. It functions as a deterministic Finite State Machine (FSM).

5.1. Finite State Machine (FSM) Design
The game state transitions linearly through specific stages. The FSM prevents illegal moves (e.g., betting on the River before the Turn is dealt).   

States:

WAITING_FOR_PLAYERS: Idle state. Listens for PlayerJoin events. Transitions to PREFLOP when active_players >= 2.

PREFLOP: Blinds are posted automatically. Hole cards are dealt. The first betting round begins starting left of the Big Blind.

FLOP: Three community cards dealt. Second betting round.

TURN: Fourth community card dealt. Third betting round.

RIVER: Fifth community card dealt. Final betting round.

SHOWDOWN: If >1 player remains, hands are revealed. The winner is determined. Funds are distributed.

PAYOUT/CLEANUP: Pot is cleared. Button moves. State resets to PREFLOP or WAITING.

5.2. Hand Evaluation Logic
Evaluating poker hands efficiently is a classic computer science problem. A naive approach (sorting and checking patterns) is too slow for a server handling thousands of hands. The architecture utilizes a Lookup Table or Bitwise Algorithm (e.g., Cactus Kev's algorithm). Each card is represented by a prime number or bitmask. The combination of cards produces a unique hash that maps directly to a hand strength rank (e.g., Royal Flush = 1, specific High Card = 7462). This reduces evaluation to essentially O(1) complexity. We will utilize optimized Go libraries such as github.com/chehsunliu/poker which implement these fast evaluation algorithms.   

5.3. Complex Betting and Side Pot Algorithm
The engine must handle the "All-In" scenario where a player has fewer chips than the current bet. This requires the creation of Side Pots. Algorithm:

Collection: Collect all bets into a temporary pool.

Sort: Sort the "All-In" bet amounts in ascending order.

Iterate: For each distinct All-In amount A:

Create a pot.

Take min(committed_chips, A) from every player and add to this pot.

Reduce every player's committed_chips by the amount taken.

Mark this pot as "contested" only by players who contributed and are still live.

Main Pot: Any remaining chips form the final side pot (or main pot if no all-ins occurred).

Resolution: At showdown, evaluate hands for each pot separately, starting from the last side pot (the one with the fewest players) down to the main pot.   

5.4. Turn Timers and Concurrency
To prevent the game from stalling, every player turn must have a hard timeout.

Mechanism: When the FSM enters a WaitForPlayer state, it spawns a Go time.Timer.

Select Loop: The Hub listens on a select statement:

Go
select {
case action := <-playerActionCh:
    if!timer.Stop() { <-timer.C } // Drain timer
    processAction(action)
case <-timer.C:
    autoFoldOrCheck(currentPlayer) // Timeout logic
}
This concurrency pattern ensures that the game loop is never blocked indefinitely by a disconnected client.   

6. Cryptographic Fairness System (Provably Fair)
Trust is the currency of online gambling. The "Provably Fair" system allows players to mathematically verify that the server did not manipulate the deck after bets were placed.

6.1. The HMAC-SHA256 Construction
The system implements a standard commitment scheme.   

Server Seed: Before the hand begins, the server generates a cryptographically secure random string (32 bytes).

Commitment: The server calculates Hash = SHA256(ServerSeed) and sends this Hash to the client. The client now holds a commitment—the server cannot change the seed without invalidating the hash.

Client Seed: The client provides a seed (or the browser generates one). This ensures the server cannot pre-calculate a shuffle that favors the house, as the server does not know the Client Seed when generating the Server Seed.

Nonce: A counter (0, 1, 2...) tracked for each hand played with this seed pair.

6.2. Deterministic Deck Generation
The deck order is not random in the traditional sense; it is a deterministic function of the inputs. Algorithm:

Calculate H = HMAC_SHA256(Key = ServerSeed, Message = ClientSeed + Nonce).

Use the byte output of H to seed a Pseudo-Random Number Generator (PRNG) or consume the bytes directly to drive a Fisher-Yates Shuffle.

Implementation Detail: To strictly avoid modulo bias, we convert sections of the hash into floating point numbers `

6.3. Verification
The platform exposes an endpoint GET /api/v1/fair/verify/:hash. After the hand concludes, the server reveals the plaintext ServerSeed. The client can then:

Hash the ServerSeed to verify it matches the initial Commitment.

Run the documented Shuffle Algorithm locally using ServerSeed + ClientSeed + Nonce.

Verify that the resulting card order matches exactly what was dealt in the game.

7. Infrastructure and Scalability
7.1. Horizontal Scaling with Redis Pub/Sub
A single server instance cannot hold unlimited WebSocket connections. To scale, we deploy multiple instances of the application behind a Load Balancer. Problem: User A (connected to Server 1) plays against User B (connected to Server 2). How does Server 1 know to send a message to User B? Solution: Redis Pub/Sub.

When the Game Engine (running on any node) generates an event (e.g., "River dealt"), it publishes the payload to a Redis Channel: poker:events:table:{id}.

All Application instances subscribe to this channel pattern.

When Server 2 receives the message, it checks its local Melody hub. "Do I have any connections subscribed to table:{id}?"

If yes, it forwards the frame via WebSocket. If no, it ignores the message. This "Stateless" web tier allows the cluster to scale elastically.   

7.2. Containerization and Deployment
The application is containerized using Docker to ensure consistency across development, testing, and production environments.

Multi-Stage Build: The Dockerfile uses a golang:alpine builder stage to compile the binary, copying only the executable to a minimal runtime image (distroless or scratch). This reduces the image size (~20MB) and attack surface.

Orchestration: docker-compose.yml defines the dependencies: app, postgres, and redis. It configures networks and volumes to ensure persistence across restarts during development.

8. Quality Assurance Strategy
8.1. Integration Testing with Testcontainers
Unit tests are insufficient for a system relying heavily on DB transactions and Redis timing. We employ testcontainers-go to spin up ephemeral infrastructure for tests. Test Workflow:   

Test starts. testcontainers spins up a real Redis and Postgres container.

Migrations are applied to the temporary Postgres DB.

The Test Client creates a user, deposits funds, creates a table, and joins it via a real WebSocket client (using gorilla/websocket client).

The test asserts that the database balance is updated and the correct WebSocket frames are received.

Test ends. Containers are destroyed.

8.2. Refactoring and Code Review
Phase 4 implies a rigorous review process. Critical focus areas include:

Race Detection: Running tests with go test -race to identify concurrent memory access violations in the game loops.

RNG Audit: Verifying that the implementation of the Shuffle algorithm strictly follows the Provably Fair specification without introducing bias.

9. Implementation Roadmap
The implementation is structured into four distinct phases to manage complexity and deliver incremental value.

Phase 1: Foundation (Weeks 1–2)
Objective: Establish the development environment and basic connectivity.

Tasks:

Initialize Go module and project structure.

Configure viper for environment variable management.

Set up Docker Compose for Postgres/Redis.

Implement basic HTML templates (landing.html, lobby.html) using Go's html/template inheritance ({{define}}, {{template}}).

Define User and GameSession structs.

Phase 2: Core Services (Weeks 3–5)
Objective: Enable user identity and real-time sockets.

Tasks:

Implement User Registration/Login with bcrypt and JWT generation.

Create the Melody Hub. Implement the connection upgrade handler with JWT query param validation.

Implement Redis Pub/Sub listener loop.

Develop the pkg/crypto package for HMAC-SHA256 generation.

Phase 3: Game Engines (Weeks 6–9)
Objective: A fully playable poker game.

Tasks:

Week 6: Implement Table CRUD. Code the FSM Skeleton (Preflop/Flop transitions).

Week 7: Integrate the Hand Evaluation library. Write unit tests for hand ranking.

Week 8: Implement the Betting Logic (Call/Check/Raise/Fold) and the Side Pot algorithm.

Week 9: Connect the Game Engine to the Wallet Service. Ensure buy-ins are deducted atomically.

Phase 4: Admin, Testing, and Polish (Weeks 10–11)
Objective: Hardening the system for production.

Tasks:

Write E2E Integration tests using testcontainers.

Implement the GET /verify endpoint for Provably Fair auditing.

Conduct a code review focusing on lock contention and concurrency patterns.

Refactor logs to use structured logging (Zap) for easier debugging.

10. Conclusion
This architectural specification presents a robust, scalable, and trustworthy foundation for Madiyar's Poker platform. By strictly adhering to the Clean Architecture principles, the system ensures that the complex domain logic of Poker is isolated from external frameworks, rendering it highly testable and adaptable. The integration of Provably Fair cryptography addresses the critical need for trust in online gaming, providing users with mathematical proof of the platform's integrity. Finally, the hybrid persistence strategy utilizing PostgreSQL and Redis, combined with the concurrency power of Go, guarantees that the system can scale to meet high player demand without compromising on performance or financial accuracy. The detailed roadmap provided herein offers a clear path to execution, minimizing risk and ensuring the delivery of a commercial-grade product.