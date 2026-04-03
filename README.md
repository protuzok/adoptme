# AdoptMe 🐾

> Animal Rescue & Adoption API built in Go. Conceived at BEST Hackathon, this pet project scales a UI concept into a robust backend architecture.

[**Live Demo (Frontend)**](https://thehaipe.github.io/volunteer/index.html) | [**Hackathon Repository**](https://github.com/thehaipe/thehaipe.github.io)

## 📖 About The Project

**AdoptMe** is a backend platform designed to connect animal shelters with volunteers and adopters. While the initial idea was born during a [BEST Hackathon](https://github.com/thehaipe/thehaipe.github.io) (initially focused on the frontend), this repository serves as a dedicated **backend pet project** to explore and demonstrate solid software engineering practices in Go.

The primary focus of this project is **Architecture**. The codebase is heavily inspired by and built upon the widely respected [`evrone/go-clean-template`](https://github.com/evrone/go-clean-template). By utilizing this template, the project strictly adheres to Uncle Bob's Clean Architecture principles, ensuring that business logic remains completely independent of frameworks, databases, or UI.

## 🏗 Architecture & Design

The application is structured into clearly separated layers:
- **Entity layer (`internal/entity`)**: Core domain objects (Animal, Shelter, Volunteer, User).
- **UseCase layer (`internal/usecase`)**: Business rules and operations (e.g., the Adoption transfer flow). Highly testable using mock interfaces.
- **Repository layer (`internal/repo`)**: Data persistence abstractions. Currently implemented for PostgreSQL.

## 🚀 Features & Roadmap

Since this is an actively evolving pet project, features are implemented incrementally.

### ✅ Currently Implemented (Core Domain)
- **Domain Entities**: Defined structures for Shelters, Volunteers, and Animals.
- **Core Business Logic (UseCases)**: Implemented the `Adoption` use case, handling the complex logic of transferring animal ownership between different actors.
- **Data Layer**: PostgreSQL repository adapters utilizing robust SQL generation.
- **Unit Testing**: Strong coverage of business logic using `gomock` and `testify`.

### 🚧 Planned / Upcoming
- [ ] **HTTP Delivery**: REST API implementation using the `Echo` framework.
- [ ] **Containerization**: `docker-compose` setup for painless local development (App + DB).
- [ ] **Database Setup**: Schema migrations using `golang-migrate/migrate`.
- [ ] **Extended Functionality**: 
  - Articles/Blog publishing for shelters and volunteers.
  - Financial donations for specific animals/shelters.
- [ ] **Advanced Tech Exploration**: Integration tests, potential gRPC endpoints, and advanced logging.

## 🛠 Tech Stack

**Currently in use:**
* [Go](https://golang.org/) (v1.26)
* [PostgreSQL](https://www.postgresql.org/) (via `jackc/pgx`)
* [Squirrel](https://github.com/Masterminds/squirrel) (SQL query builder)
* [Testify](https://github.com/stretchr/testify) & [gomock](https://github.com/uber/mock) (Testing)

**Targeting for future implementation:**
* *Echo (Router), golang-migrate, Docker.*

---
*This is a portfolio project aimed at demonstrating complex system design, clean architecture, and modern Go development practices.*
