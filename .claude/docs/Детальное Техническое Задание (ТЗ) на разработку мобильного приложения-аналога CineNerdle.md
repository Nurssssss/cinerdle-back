# Детальное Техническое Задание (ТЗ) на разработку мобильного приложения-аналога CineNerdle

## 1. Введение

Данное техническое задание (ТЗ) описывает функциональные и нефункциональные требования к мобильному приложению, являющемуся аналогом игры CineNerdle, с акцентом на соревновательный режим "Classic Battle" (1v1). Проект разработан с целью глубокого изучения современных подходов к кроссплатформенной мобильной разработке и высокопроизводительному бэкенду. Приложение будет реализовано с использованием стека Flutter для мобильного клиента и Go (чистый `net/http`) для бэкенда, с базой данных PostgreSQL.

**Цель проекта:** Создание интерактивной игры, где пользователи соревнуются в знании кино, связывая фильмы через общих актеров, режиссеров или композиторов. Проект также служит платформой для обучения и глубокого понимания принципов разработки, включая Clean Architecture, WebSockets, JWT-аутентификацию и взаимодействие с внешними API.

## 2. Стек Технологий

В проекте будут использованы следующие ключевые технологии:

*   **Мобильный клиент:** Flutter (Dart)
    *   Управление состоянием: BLoC/Cubit
    *   Сетевые запросы: `dio`
*   **Бэкенд:** Go (чистый `net/http`)
    *   Архитектура: Clean Architecture
    *   База данных: PostgreSQL
    *   Реальное время: WebSockets (`gorilla/websocket`)
*   **Внешние API:** The Movie Database (TMDb) API

## 3. Функциональные Требования

### 3.1. Пользовательская Аутентификация и Авторизация

*   **Регистрация:** Пользователь может зарегистрироваться, используя уникальный email и пароль, а также выбрать никнейм.
*   **Подтверждение Email:** После успешной регистрации на указанный email отправляется письмо с уникальной ссылкой/кодом для подтверждения аккаунта. До момента подтверждения аккаунт считается неактивным, и пользователь не может получить доступ к защищенным функциям приложения.
*   **Вход (Login):** Пользователь может войти в систему, используя зарегистрированные email и пароль. При успешном входе выдается JWT (JSON Web Token) для последующей авторизации.
*   **Выход (Logout):** Пользователь может выйти из системы, что приводит к аннулированию текущего JWT на клиенте.
*   **Управление сессиями:** Авторизация осуществляется посредством JWT. Токен, полученный при входе, должен прикрепляться к заголовку `Authorization: Bearer <JWT_TOKEN>` каждого защищенного запроса к бэкенду.

### 3.2. Игровой Режим "Classic Battle" (1v1)

*   **Поиск оппонента:** Пользователь может инициировать поиск случайного оппонента для игры 1v1. Система должна эффективно подбирать доступных игроков.
*   **Начало игры:** После успешного нахождения оппонента инициируется новая игровая сессия, и оба игрока уведомляются о начале игры через WebSocket соединение.
*   **Игровой процесс:**
    *   Игра состоит из последовательных раундов, в каждом из которых игроки по очереди называют фильмы.
    *   Первый игрок начинает, называя любой существующий фильм.
    *   Последующие игроки должны предложить фильм, который имеет общую персону (актера, режиссера или композитора) с фильмом, названным предыдущим игроком.
    *   **Проверка связи:** Бэкенд выполняет проверку корректности предложенной связи, используя кэшированные данные из TMDb. Проверка должна быть быстрой и точной.
    *   **Таймер хода:** На каждый ход игроку отводится ограниченное время (например, 30-60 секунд). Истечение времени приводит к проигрышу хода.
    *   **Пропуск хода (Skip):** Каждый игрок имеет ограниченное количество пропусков за игру (например, 1-2). Использование пропуска передает ход оппоненту без штрафа, но уменьшает доступное количество пропусков.
    *   **Подсказки (Hint):** Каждый игрок имеет ограниченное количество подсказок за игру (например, 1-2). Использование подсказки раскрывает список актеров, режиссеров или композиторов последнего названного фильма. Это помогает игроку найти связанный фильм.
    *   **Окончание игры:** Игра завершается, если один из игроков:
        *   Не смог назвать связанный фильм в отведенное время.
        *   Исчерпал все доступные пропуски и не смог сделать ход.
        *   Предложил фильм, связь которого не подтверждена бэкендом.
*   **Результаты игры:** По окончании игры обоим игрокам отображается победитель, проигравший и детальная статистика игровой сессии.
*   **Реальное время:** Все игровые события (начало игры, смена хода, действия игроков, таймер, результаты проверки, окончание игры) должны мгновенно синхронизироваться между клиентами через WebSocket соединения.

### 3.3. Управление Профилем

*   **Просмотр профиля:** Пользователь может просматривать свой никнейм, email, статус верификации, а также базовую игровую статистику (количество побед и поражений).
*   **Редактирование никнейма:** Пользователь может изменить свой никнейм, при условии, что новый никнейм уникален.

## 4. Нефункциональные Требования

*   **Производительность:**
    *   Время ответа REST API запросов (кроме запросов к TMDb) не должно превышать 200 мс.
    *   Проверка связей между фильмами на бэкенде должна занимать не более 1-2 секунд.
    *   Обновления через WebSocket должны доставляться клиентам в течение 100 мс.
*   **Масштабируемость:** Архитектура бэкенда должна быть спроектирована с учетом возможности горизонтального масштабирования для поддержки растущего числа одновременных пользователей и игровых сессий.
*   **Безопасность:**
    *   Пароли пользователей должны храниться в базе данных в виде криптографических хешей (например, bcrypt).
    *   Все коммуникации между клиентом и сервером должны осуществляться по HTTPS/WSS.
    *   Приложение должно быть защищено от распространенных уязвимостей, таких как SQL-инъекции, XSS, CSRF.
    *   JWT токены должны иметь ограниченное время жизни и быть защищены от перехвата.
*   **Надежность:**
    *   Бэкенд должен корректно обрабатывать ошибки внешних API (TMDb), сетевые сбои и ошибки базы данных.
    *   WebSocket соединения должны быть устойчивыми к кратковременным разрывам и иметь механизм переподключения.
*   **Удобство использования (UX):**
    *   Интуитивно понятный и отзывчивый пользовательский интерфейс.
    *   Четкая визуальная обратная связь о состоянии игры, таймере, доступных подсказках и пропусках.
    *   Автодополнение при вводе названий фильмов для улучшения игрового опыта.

## 5. Архитектура Системы

### 5.1. База Данных (PostgreSQL)

PostgreSQL будет использоваться как основное хранилище данных. Схема базы данных включает таблицы для пользователей, верификации email, кэшированных фильмов и персон из TMDb, а также для игровых сессий и ходов.

**SQL Schema для PostgreSQL:**

```sql
-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(50) UNIQUE NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_nickname ON users (nickname);

-- Таблица для подтверждения email
CREATE TABLE IF NOT EXISTS email_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    verification_code VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_email_verifications_user_id ON email_verifications (user_id);
CREATE INDEX IF NOT EXISTS idx_email_verifications_code ON email_verifications (verification_code);

-- Таблица фильмов (кэш из TMDb)
CREATE TABLE IF NOT EXISTS movies (
    id INTEGER PRIMARY KEY, -- TMDb ID
    title VARCHAR(255) NOT NULL,
    original_title VARCHAR(255),
    poster_path VARCHAR(255),
    release_date DATE,
    overview TEXT,
    popularity REAL,
    vote_average REAL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_movies_title ON movies (title);
CREATE INDEX IF NOT EXISTS idx_movies_release_date ON movies (release_date DESC);

-- Таблица персон (актеры, режиссеры, композиторы - кэш из TMDb)
CREATE TABLE IF NOT EXISTS persons (
    id INTEGER PRIMARY KEY, -- TMDb ID
    name VARCHAR(255) NOT NULL,
    profile_path VARCHAR(255),
    birthday DATE,
    deathday DATE,
    biography TEXT,
    popularity REAL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_persons_name ON persons (name);

-- Таблица связей между фильмами и персонами (роли)
CREATE TYPE person_role AS ENUM (
    'actor',
    'director',
    'composer'
);

CREATE TABLE IF NOT EXISTS movie_person_roles (
    movie_id INTEGER NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    person_id INTEGER NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    role person_role NOT NULL,
    PRIMARY KEY (movie_id, person_id, role),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_movie_person_roles_movie_id ON movie_person_roles (movie_id);
CREATE INDEX IF NOT EXISTS idx_movie_person_roles_person_id ON movie_person_roles (person_id);
CREATE INDEX IF NOT EXISTS idx_movie_person_roles_role ON movie_person_roles (role);

-- Таблица игровых сессий
CREATE TYPE game_status AS ENUM (
    'waiting_for_player',
    'in_progress',
    'finished',
    'cancelled'
);

CREATE TABLE IF NOT EXISTS game_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player1_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    player2_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Может быть NULL, если игрок отключился
    current_turn_user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Чей сейчас ход
    status game_status NOT NULL DEFAULT 'waiting_for_player',
    last_movie_id INTEGER REFERENCES movies(id) ON DELETE SET NULL, -- Последний названный фильм
    player1_skips_left INTEGER DEFAULT 1,
    player2_skips_left INTEGER DEFAULT 1,
    player1_hints_left INTEGER DEFAULT 1,
    player2_hints_left INTEGER DEFAULT 1,
    turn_ends_at TIMESTAMP WITH TIME ZONE, -- Время окончания текущего хода
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    winner_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_game_sessions_player1_id ON game_sessions (player1_id);
CREATE INDEX IF NOT EXISTS idx_game_sessions_player2_id ON game_sessions (player2_id);
CREATE INDEX IF NOT EXISTS idx_game_sessions_status ON game_sessions (status);

-- Таблица ходов в игре
CREATE TYPE move_type AS ENUM (
    'movie_guess',
    'skip',
    'hint'
);

CREATE TABLE IF NOT EXISTS game_moves (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES movies(id) ON DELETE SET NULL, -- NULL, если ход был пропуском/подсказкой
    move_type move_type NOT NULL,
    is_valid BOOLEAN, -- NULL для skip/hint, TRUE/FALSE для movie_guess
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_game_moves_session_id ON game_moves (session_id);
CREATE INDEX IF NOT EXISTS idx_game_moves_user_id ON game_moves (user_id);
```

### 5.2. Бэкенд (Go) - Clean Architecture

Выбор чистого `net/http` в сочетании с **Clean Architecture** [1] позволит создать гибкий, тестируемый и поддерживаемый бэкенд. Основная идея Clean Architecture — разделение на слои, где каждый слой имеет свою ответственность и зависимости направлены внутрь.

#### 5.2.1. Структура Папок

```
/cmd
  /api
    main.go                 # Точка входа для HTTP/WebSocket сервера
/internal
  /adapters
    /repository             # Реализации интерфейсов репозиториев (PostgreSQL)
      user_repository.go
      movie_repository.go
      game_session_repository.go
    /tmdb                   # Клиент для взаимодействия с TMDb API
      tmdb_client.go
    /email                  # Сервис для отправки email (например, Mailtrap)
      email_sender.go
  /domain                   # Core бизнес-сущности и правила (не зависят от других слоев)
    user.go
    movie.go
    person.go
    game_session.go
    game_move.go
    # ... другие доменные сущности
  /usecases                 # Бизнес-логика, оркестрация доменных сущностей и репозиториев
    auth_usecase.go
    game_usecase.go
    profile_usecase.go
    # ... другие use cases
  /delivery                 # Внешний интерфейс (HTTP, WebSocket)
    /http
      auth_handler.go
      game_handler.go
      profile_handler.go
      router.go             # Настройка маршрутов и middleware
    /websocket
      game_websocket_handler.go # Обработка WebSocket соединений для игр
/pkg
  /jwt                      # Утилиты для работы с JWT токенами
    jwt_manager.go
  /password                 # Утилиты для хеширования паролей
    password_hasher.go
  /errors                   # Кастомные ошибки приложения
    app_errors.go
  /config                   # Загрузка и управление конфигурацией приложения
    config.go
  /logger                   # Настройка логирования
    logger.go
```

#### 5.2.2. Слои Clean Architecture

*   **Domain Layer (Внутренний слой):** Содержит основные бизнес-сущности и правила. Этот слой не должен зависеть ни от каких внешних фреймворков, баз данных или UI. Он определяет *что* делает бизнес.
    *   **Пример структуры `User`:**
        ```go
        package domain

        import (
        	"time"
        	"github.com/google/uuid"
        )

        type User struct {
        	ID           uuid.UUID
        	Email        string
        	PasswordHash string
        	Nickname     string
        	IsVerified   bool
        	Wins         int
        	Losses       int
        	CreatedAt    time.Time
        	UpdatedAt    time.Time
        }

        // UserRepository определяет интерфейс для работы с хранилищем пользователей.
        // Реализация этого интерфейса будет находиться в слое Adapters.
        type UserRepository interface {
        	Create(user *User) error
        	FindByEmail(email string) (*User, error)
        	FindByID(id uuid.UUID) (*User, error)
        	Update(user *User) error
        	// ... другие методы, например, для обновления статистики
        }

        // EmailVerificationRepository определяет интерфейс для работы с верификациями email.
        type EmailVerificationRepository interface {
        	Create(verification *EmailVerification) error
        	FindByCode(code string) (*EmailVerification, error)
        	Delete(id uuid.UUID) error
        }

        // ... аналогичные интерфейсы для MovieRepository, PersonRepository, GameSessionRepository, GameMoveRepository
        ```

*   **Usecases Layer (Слой Применения):** Содержит бизнес-логику приложения, специфичную для конкретных вариантов использования (use cases). Он оркестрирует доменные сущности и взаимодействует с интерфейсами репозиториев, определенных в доменном слое. Этот слой определяет *как* бизнес выполняет свои операции.
    *   **Пример интерфейса `AuthUsecase`:**
        ```go
        package usecases

        import (
        	"context"
        	"time"
        	"your_project/internal/domain"
        	"github.com/google/uuid"
        )

        type AuthUsecase interface {
        	Register(ctx context.Context, email, password, nickname string) (*domain.User, error)
        	VerifyEmail(ctx context.Context, code string) error
        	Login(ctx context.Context, email, password string) (string, error) // Возвращает JWT токен
        	// ... другие методы, например, для восстановления пароля
        }

        // AuthUsecaseImpl - реализация use case
        type AuthUsecaseImpl struct {
        	userRepo              domain.UserRepository
        	emailVerificationRepo domain.EmailVerificationRepository
        	passwordHasher        domain.PasswordHasher // Интерфейс для хеширования паролей
        	jwtManager            domain.JWTManager     // Интерфейс для работы с JWT
        	emailSender           domain.EmailSender    // Интерфейс для отправки email
        	// ... другие зависимости
        }

        func NewAuthUsecase(userRepo domain.UserRepository, ... ) *AuthUsecaseImpl {
        	return &AuthUsecaseImpl{userRepo: userRepo, ...}
        }

        func (uc *AuthUsecaseImpl) Register(ctx context.Context, email, password, nickname string) (*domain.User, error) {
        	// 1. Проверить, что пользователь с таким email или никнеймом не существует
        	// 2. Хешировать пароль
        	// 3. Создать нового пользователя
        	// 4. Сохранить пользователя через userRepo.Create()
        	// 5. Создать код верификации email и сохранить через emailVerificationRepo.Create()
        	// 6. Отправить письмо с верификацией через emailSender.SendVerificationEmail()
        	// 7. Вернуть пользователя
        	return nil, nil // Заглушка
        }

        // ... другие методы AuthUsecase
        ```

*   **Adapters Layer (Слой Адаптеров):** Этот слой содержит реализации интерфейсов, определенных в доменном слое. Он адаптирует внешние детали (базы данных, сторонние API) к внутренним требованиям приложения. Здесь находятся конкретные реализации репозиториев, клиентов для внешних сервисов.
    *   **Пример реализации `UserRepository` для PostgreSQL:**
        ```go
        package repository

        import (
        	"context"
        	"database/sql"
        	"your_project/internal/domain"
        	"github.com/google/uuid"
        )

        type PostgresUserRepository struct {
        	db *sql.DB
        }

        func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
        	return &PostgresUserRepository{db: db}
        }

        func (r *PostgresUserRepository) Create(user *domain.User) error {
        	query := `INSERT INTO users (id, email, password_hash, nickname, is_verified, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
        	stmt, err := r.db.Prepare(query)
        	if err != nil {
        		return err
        	}
        	defer stmt.Close()

        	err = stmt.QueryRow(
        		user.ID,
        		user.Email,
        		user.PasswordHash,
        		user.Nickname,
        		user.IsVerified,
        		user.CreatedAt,
        		user.UpdatedAt,
        	).Scan(&user.ID)

        	return err
        }

        func (r *PostgresUserRepository) FindByEmail(email string) (*domain.User, error) {
        	user := &domain.User{}
        	query := `SELECT id, email, password_hash, nickname, is_verified, wins, losses, created_at, updated_at FROM users WHERE email = $1`
        	err := r.db.QueryRow(query, email).Scan(
        		&user.ID,
        		&user.Email,
        		&user.PasswordHash,
        		&user.Nickname,
        		&user.IsVerified,
        		&user.Wins,
        		&user.Losses,
        		&user.CreatedAt,
        		&user.UpdatedAt,
        	)
        	if err == sql.ErrNoRows {
        		return nil, domain.ErrUserNotFound
        	}
        	return user, err
        }

        // ... другие методы UserRepository
        ```

    *   **Пример клиента TMDb:**
        ```go
        package tmdb

        import (
        	"encoding/json"
        	"fmt"
        	"net/http"
        	"time"
        	"your_project/internal/domain"
        )

        const tmdbAPIURL = "https://api.themoviedb.org/3"

        type TMDbClient struct {
        	apiKey     string
        	httpClient *http.Client
        }

        func NewTMDbClient(apiKey string) *TMDbClient {
        	return &TMDbClient{
        		apiKey: apiKey,
        		httpClient: &http.Client{Timeout: 10 * time.Second},
        	}
        }

        type tmdbMovieResponse struct {
        	ID          int    `json:"id"`
        	Title       string `json:"title"`
        	PosterPath  string `json:"poster_path"`
        	ReleaseDate string `json:"release_date"`
        	Overview    string `json:"overview"`
        }

        // GetMovieByID fetches movie details from TMDb and converts to domain.Movie
        func (c *TMDbClient) GetMovieByID(ctx context.Context, movieID int) (*domain.Movie, error) {
        	url := fmt.Sprintf("%s/movie/%d?api_key=%s", tmdbAPIURL, movieID, c.apiKey)
        	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
        	if err != nil {
        		return nil, err
        	}

        	res, err := c.httpClient.Do(req)
        	if err != nil {
        		return nil, err
        	}
        	defer res.Body.Close()

        	if res.StatusCode != http.StatusOK {
        		return nil, fmt.Errorf("TMDb API returned status %d", res.StatusCode)
        	}

        	var tmdbRes tmdbMovieResponse
        	if err := json.NewDecoder(res.Body).Decode(&tmdbRes); err != nil {
        		return nil, err
        	}

        	// Преобразование в доменную сущность
        	movie := &domain.Movie{
        		ID:         tmdbRes.ID,
        		Title:      tmdbRes.Title,
        		PosterPath: tmdbRes.PosterPath,
        		Overview:   tmdbRes.Overview,
        	}
        	// Парсинг даты и т.д.

        	return movie, nil
        }

        // ... методы для получения актеров, режиссеров, композиторов
        ```

*   **Delivery Layer (Внешний слой):** Этот слой отвечает за представление данных пользователю и обработку внешних запросов. Здесь находятся HTTP-хендлеры, WebSocket-обработчики. Они преобразуют запросы из внешнего формата во внутренний (вызовы use cases) и обратно.
    *   **Пример HTTP-хендлера для регистрации:**
        ```go
        package http

        import (
        	"encoding/json"
        	"net/http"
        	"your_project/internal/usecases"
        	"your_project/pkg/errors"
        )

        type AuthHandler struct {
        	authUsecase usecases.AuthUsecase
        	// ... другие зависимости, например, логгер
        }

        func NewAuthHandler(authUsecase usecases.AuthUsecase) *AuthHandler {
        	return &AuthHandler{authUsecase: authUsecase}
        }

        type registerRequest struct {
        	Email    string `json:"email"`
        	Password string `json:"password"`
        	Nickname string `json:"nickname"`
        }

        func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
        	var req registerRequest
        	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        		http.Error(w, errors.ErrInvalidInput.Error(), http.StatusBadRequest)
        		return
        	}

        	user, err := h.authUsecase.Register(r.Context(), req.Email, req.Password, req.Nickname)
        	if err != nil {
        		// Обработка различных ошибок use case
        		http.Error(w, err.Error(), http.StatusInternalServerError)
        		return
        	}

        	w.WriteHeader(http.StatusCreated)
        	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully. Please check your email for verification.", "user_id": user.ID.String()})
        }

        // ... другие хендлеры для логина, верификации email
        ```

#### 5.2.3. Dependency Injection (Внедрение Зависимостей)

В Go внедрение зависимостей обычно реализуется вручную через конструкторы. Это означает, что при создании объекта (например, `AuthUsecaseImpl`), ему передаются все необходимые зависимости (например, `UserRepository`, `PasswordHasher`).

**Пример в `main.go`:**

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver

	"your_project/internal/adapters/email"
	"your_project/internal/adapters/repository"
	"your_project/internal/adapters/tmdb"
	"your_project/internal/delivery/http"
	"your_project/internal/usecases"
	"your_project/pkg/config"
	"your_project/pkg/jwt"
	"your_project/pkg/password"
)

func main() {
	cfg := config.LoadConfig() // Загрузка конфигурации из .env или переменных окружения

	// Инициализация базы данных
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация адаптеров (репозиториев, клиентов)
	userRepo := repository.NewPostgresUserRepository(db)
	emailVerificationRepo := repository.NewPostgresEmailVerificationRepository(db)
	movieRepo := repository.NewPostgresMovieRepository(db)
	personRepo := repository.NewPostgresPersonRepository(db)
	moviePersonRoleRepo := repository.NewPostgresMoviePersonRoleRepository(db)
	gameSessionRepo := repository.NewPostgresGameSessionRepository(db)
	gameMoveRepo := repository.NewPostgresGameMoveRepository(db)

	tmdbClient := tmdb.NewTMDbClient(cfg.TMDbAPIKey)
	emailSender := email.NewMailtrapEmailSender(cfg.MailtrapHost, cfg.MailtrapPort, cfg.MailtrapUser, cfg.MailtrapPass)

	// Инициализация pkg утилит
	passwordHasher := password.NewBcryptHasher()
	jwtManager := jwt.NewJWTManager(cfg.JWTSecret, time.Hour*24) // Токен действителен 24 часа

	// Инициализация use cases
	authUsecase := usecases.NewAuthUsecase(
		userRepo,
		emailVerificationRepo,
		passwordHasher,
		jwtManager,
		emailSender,
		cfg.AppBaseURL, // Для формирования ссылки верификации
	)
	gameUsecase := usecases.NewGameUsecase(
		gameSessionRepo,
		gameMoveRepo,
		userRepo,
		movieRepo,
		personRepo,
		moviePersonRoleRepo,
		tmdbClient,
		cfg.TurnTimeoutSeconds,
	)
	profileUsecase := usecases.NewProfileUsecase(userRepo)

	// Инициализация хендлеров Delivery слоя
	authHandler := http.NewAuthHandler(authUsecase)
	gameHandler := http.NewGameHandler(gameUsecase, jwtManager) // JWTManager для извлечения user ID из токена
	profileHandler := http.NewProfileHandler(profileUsecase, jwtManager)

	// Настройка маршрутизатора
	router := http.NewRouter(authHandler, gameHandler, profileHandler, jwtManager) // Передача хендлеров и JWTManager для middleware

	// Запуск HTTP сервера
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router))
}
```

#### 5.2.4. Обработка Ошибок

Будет использоваться централизованный подход к ошибкам с кастомными типами ошибок в `pkg/errors`, что позволит четко различать ошибки на разных слоях и корректно возвращать их клиенту.

```go
package errors

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrNicknameAlreadyExists = errors.New("nickname already exists")
	ErrVerificationCodeExpired = errors.New("verification code expired or invalid")
	ErrGameSessionNotFound  = errors.New("game session not found")
	ErrInvalidGameMove      = errors.New("invalid game move")
	ErrNotYourTurn          = errors.New("not your turn")
	ErrNoSkipsLeft          = errors.New("no skips left")
	ErrNoHintsLeft          = errors.New("no hints left")
	ErrMovieNotFound        = errors.New("movie not found")
	ErrTMDbAPIError         = errors.New("TMDb API error")
	ErrInternalServer       = errors.New("internal server error")
	ErrInvalidInput         = errors.New("invalid input")
)

// CustomError - интерфейс для кастомных ошибок с кодом статуса HTTP
type CustomError interface {
	Error() string
	StatusCode() int
}

// Implementations for CustomError for specific HTTP status codes
// ...
```

#### 5.2.5. Middleware

Для обработки JWT токенов и других общих задач (логирование, CORS) будут использоваться HTTP middleware. Пример middleware для аутентификации:

```go
package http

import (
	"context"
	"net/http"
	"strings"
	"your_project/pkg/jwt"
	"your_project/pkg/errors"
)

func AuthMiddleware(jwtManager *jwt.JWTManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, errors.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, errors.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]
		claims, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, errors.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		// Добавляем user ID в контекст запроса
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
```

### 5.3. API и Протокол WebSockets

#### 5.3.1. REST API

Все REST API запросы будут использовать JSON в качестве формата данных и будут защищены JWT токенами, за исключением эндпоинтов регистрации и логина.

**1. Аутентификация**

*   **`POST /api/v1/auth/register`**
    *   **Описание:** Регистрация нового пользователя.
    *   **Запрос (Request Body):**
        ```json
        {
          "email": "user@example.com",
          "password": "StrongPassword123!",
          "nickname": "PlayerOne"
        }
        ```
    *   **Ответ (Response Body - 201 Created):**
        ```json
        {
          "message": "User registered successfully. Please check your email for verification.",
          "user_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
        }
        ```
    *   **Ошибки (Error Response - 400 Bad Request / 409 Conflict):**
        ```json
        {
          "error": "email already exists"
        }
        ```

*   **`GET /api/v1/auth/verify-email?code={code}`**
    *   **Описание:** Подтверждение email пользователя.
    *   **Параметры запроса (Query Parameters):**
        *   `code`: Уникальный код верификации, отправленный на email.
    *   **Ответ (Response Body - 200 OK):**
        ```json
        {
          "message": "Email verified successfully."
        }
        ```
    *   **Ошибки (Error Response - 400 Bad Request / 404 Not Found):**
        ```json
        {
          "error": "verification code expired or invalid"
        }
        ```

*   **`POST /api/v1/auth/login`**
    *   **Описание:** Вход пользователя в систему.
    *   **Запрос (Request Body):**
        ```json
        {
          "email": "user@example.com",
          "password": "StrongPassword123!"
        }
        ```
    *   **Ответ (Response Body - 200 OK):**
        ```json
        {
          "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
          "user_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
        }
        ```
    *   **Ошибки (Error Response - 401 Unauthorized):**
        ```json
        {
          "error": "invalid credentials"
        }
        ```

**2. Профиль Пользователя**

*   **`GET /api/v1/profile`**
    *   **Описание:** Получение данных профиля текущего пользователя.
    *   **Заголовок (Headers):** `Authorization: Bearer <JWT_TOKEN>`
    *   **Ответ (Response Body - 200 OK):**
        ```json
        {
          "id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
          "email": "user@example.com",
          "nickname": "PlayerOne",
          "is_verified": true,
          "wins": 10,
          "losses": 5,
          "created_at": "2023-10-27T10:00:00Z"
        }
        ```

*   **`PUT /api/v1/profile`**
    *   **Описание:** Обновление никнейма пользователя.
    *   **Заголовок (Headers):** `Authorization: Bearer <JWT_TOKEN>`
    *   **Запрос (Request Body):**
        ```json
        {
          "nickname": "NewPlayerName"
        }
        ```
    *   **Ответ (Response Body - 200 OK):**
        ```json
        {
          "message": "Profile updated successfully.",
          "nickname": "NewPlayerName"
        }
        ```
    *   **Ошибки (Error Response - 409 Conflict):**
        ```json
        {
          "error": "nickname already exists"
        }
        ```

**3. Поиск Игры**

*   **`POST /api/v1/game/find-match`**
    *   **Описание:** Инициирует поиск оппонента для игры 1v1. Пользователь помещается в очередь.
    *   **Заголовок (Headers):** `Authorization: Bearer <JWT_TOKEN>`
    *   **Запрос (Request Body):** (Пустой или может содержать параметры для фильтрации, если будут)
        ```json
        {}
        ```
    *   **Ответ (Response Body - 202 Accepted):**
        ```json
        {
          "message": "Searching for opponent...",
          "user_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef"
        }
        ```

#### 5.3.2. Протокол WebSockets для Игровых Сессий

WebSockets будут использоваться для обмена сообщениями в реальном времени между клиентами и сервером во время игровой сессии. Каждая игровая сессия будет иметь свой уникальный WebSocket URL.

**`GET /ws/game/{session_id}`**
*   **Описание:** Установка WebSocket соединения для конкретной игровой сессии.
*   **Параметры пути (Path Parameters):**
    *   `session_id`: UUID игровой сессии.
*   **Заголовок (Headers):** `Authorization: Bearer <JWT_TOKEN>` (JWT токен передается как часть заголовка при установке соединения).

**1. Типы Сообщений (JSON)**

Все сообщения будут иметь поле `type` для определения типа события и `payload` для данных.

*   **Сервер -> Клиент (Server to Client):**
    *   **`game_started`**
        *   **Описание:** Уведомление о начале игры, когда найден второй игрок.
        *   **Payload:**
            ```json
            {
              "type": "game_started",
              "payload": {
                "session_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
                "player1": {"id": "uuid1", "nickname": "PlayerOne"},
                "player2": {"id": "uuid2", "nickname": "PlayerTwo"},
                "current_turn_user_id": "uuid1",
                "turn_ends_at": "2023-10-27T10:01:00Z",
                "last_movie": null,
                "player1_skips_left": 1,
                "player2_skips_left": 1,
                "player1_hints_left": 1,
                "player2_hints_left": 1
              }
            }
            ```

    *   **`turn_update`**
        *   **Описание:** Обновление информации о текущем ходе.
        *   **Payload:**
            ```json
            {
              "type": "turn_update",
              "payload": {
                "current_turn_user_id": "uuid2",
                "turn_ends_at": "2023-10-27T10:02:00Z",
                "last_movie": {
                  "id": 123,
                  "title": "Inception",
                  "poster_path": "/poster.jpg"
                },
                "player1_skips_left": 1,
                "player2_skips_left": 1,
                "player1_hints_left": 1,
                "player2_hints_left": 1
              }
            }
            ```

    *   **`game_ended`**
        *   **Описание:** Уведомление об окончании игры.
        *   **Payload:**
            ```json
            {
              "type": "game_ended",
              "payload": {
                "session_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
                "winner_id": "uuid1",
                "reason": "PlayerTwo ran out of time",
                "final_score": {"player1": 5, "player2": 4}
              }
            }
            ```

    *   **`error`**
        *   **Описание:** Сообщение об ошибке в ходе игры.
        *   **Payload:**
            ```json
            {
              "type": "error",
              "payload": {
                "code": "invalid_move",
                "message": "Movie is not connected to the previous one."
              }
            }
            ```

    *   **`hint_revealed`**
        *   **Описание:** Раскрытие подсказки (актеры/режиссеры/композиторы).
        *   **Payload:**
            ```json
            {
              "type": "hint_revealed",
              "payload": {
                "movie_id": 123,
                "persons": [
                  {"id": 1, "name": "Leonardo DiCaprio", "role": "actor"},
                  {"id": 2, "name": "Christopher Nolan", "role": "director"}
                ]
              }
            }
            ```

*   **Клиент -> Сервер (Client to Server):**
    *   **`make_move`**
        *   **Описание:** Отправка хода игрока (предложение фильма).
        *   **Payload:**
            ```json
            {
              "type": "make_move",
              "payload": {
                "movie_id": 456, // TMDb ID предложенного фильма
                "movie_title": "The Dark Knight" // Для удобства, но бэкенд будет проверять по ID
              }
            }
            ```

    *   **`skip_turn`**
        *   **Описание:** Запрос на пропуск хода.
        *   **Payload:**
            ```json
            {
              "type": "skip_turn"
            }
            ```

    *   **`request_hint`**
        *   **Описание:** Запрос на использование подсказки.
        *   **Payload:**
            ```json
            {
              "type": "request_hint"
            }
            ```

**2. Обработка Соединений**

*   Сервер должен отслеживать активные WebSocket соединения для каждой игровой сессии, используя `session_id`.
*   При отключении игрока (например, закрытие приложения или потеря соединения), сервер должен корректно завершить игру или пометить игрока как отключенного, а затем завершить игру по таймауту, если игрок не переподключится.
*   Все входящие и исходящие сообщения WebSocket должны быть валидированы на сервере для предотвращения некорректных данных и уязвимостей.

#### 5.3.3. Взаимодействие с TMDb API

*   **Получение API-ключа:** Для работы с TMDb API необходимо зарегистрироваться на [The Movie Database (TMDb)](https://www.themoviedb.org/documentation/api) и получить персональный API-ключ.
*   **Основные эндпоинты TMDb, которые будут использоваться:**
    *   `GET /movie/{movie_id}`: Получение детальной информации о фильме (название, постер, дата выхода, описание).
    *   `GET /movie/{movie_id}/credits`: Получение списка актеров, режиссеров, композиторов для конкретного фильма.
    *   `GET /search/movie`: Поиск фильмов по названию, который будет использоваться для автодополнения в клиентском приложении.
    *   `GET /person/{person_id}`: Получение детальной информации о персоне (имя, фото, биография).
*   **Кэширование:** Все данные, полученные от TMDb, должны быть кэшированы в локальной базе данных (таблицы `movies`, `persons`, `movie_person_roles`). Это критически важно для уменьшения нагрузки на внешний API, соблюдения лимитов запросов и ускорения работы приложения, особенно при проверке связей.
*   **Обработка ошибок:** Клиент TMDb в Go должен корректно обрабатывать различные ошибки API (например, лимиты запросов, неверные ID, недоступность сервиса) и возвращать соответствующие ошибки на уровень `usecases`.

#### 5.3.4. Отправка Email

*   **SMTP-сервер:** Для отправки писем подтверждения email потребуется настроить SMTP-клиент в Go. Рекомендуется использовать библиотеку `net/smtp` или сторонние решения.
*   **Сервисы:** Для разработки и тестирования рекомендуется использовать сервисы типа [Mailtrap](https://mailtrap.io/) (для тестовых писем, которые не доставляются реальным пользователям) или [Mailgun](https://www.mailgun.com/) (для реальной отправки в продакшене).
*   **Содержание письма:** Письмо должно содержать четкую инструкцию и ссылку для верификации, например: `https://your-app.com/api/v1/auth/verify-email?code={verification_code}`. Ссылка должна быть сформирована с использованием `AppBaseURL` из конфигурации бэкенда.

### 5.4. Мобильный Клиент (Flutter)

Flutter клиент будет построен с использованием архитектуры BLoC/Cubit для управления состоянием, что обеспечит чистоту кода, тестируемость и масштабируемость. Для сетевых запросов будет использоваться библиотека `dio`.

#### 5.4.1. Структура Папок (Пример)

```
/lib
  /api
    api_client.dart           # Dio клиент для REST API
    websocket_client.dart     # Клиент для WebSocket соединений
  /auth
    /bloc
      auth_bloc.dart
      auth_event.dart
      auth_state.dart
    /repository
      auth_repository.dart    # Взаимодействие с Auth API
    auth_screen.dart
  /game
    /bloc
      game_search_bloc.dart
      game_search_event.dart
      game_search_state.dart
      game_bloc.dart
      game_event.dart
      game_state.dart
    /repository
      game_repository.dart    # Взаимодействие с Game API и WebSocket
    game_search_screen.dart
    game_play_screen.dart
  /profile
    /bloc
      profile_bloc.dart
      profile_event.dart
      profile_state.dart
    /repository
      profile_repository.dart # Взаимодействие с Profile API
    profile_screen.dart
  /models
    user.dart
    movie.dart
    person.dart
    game_session.dart
    # ... другие модели данных
  /utils
    app_router.dart           # Навигация (например, GoRouter)
    constants.dart
    shared_preferences.dart   # Для хранения JWT
  main.dart
```

#### 5.4.2. Управление Состоянием с BLoC/Cubit

Каждый функциональный модуль (Auth, Game, Profile) будет иметь свой BLoC/Cubit для управления состоянием. Это позволяет изолировать логику, упростить тестирование и повысить предсказуемость поведения приложения.

**1. `AuthBloc`**

*   **События (Events):**
    *   `AuthRegisterRequested(email, password, nickname)`: Запрос на регистрацию.
    *   `AuthLoginRequested(email, password)`: Запрос на вход.
    *   `AuthLogoutRequested()`: Запрос на выход из системы.
    *   `AuthCheckStatus()`: Проверка текущего статуса аутентификации пользователя (например, при запуске приложения или возобновлении сессии).
    *   `AuthEmailVerified()`: Событие, сигнализирующее об успешной верификации email после перехода по ссылке.
*   **Состояния (States):**
    *   `AuthInitial`: Начальное состояние BLoC.
    *   `AuthLoading`: Состояние загрузки, указывающее на выполнение асинхронной операции (например, сетевого запроса).
    *   `AuthAuthenticated(user)`: Пользователь успешно аутентифицирован. Содержит объект `User` с данными профиля.
    *   `AuthUnauthenticated`: Пользователь не аутентифицирован или вышел из системы.
    *   `AuthRegistrationSuccess`: Успешная регистрация, но требуется верификация email.
    *   `AuthError(message)`: Состояние ошибки, содержащее сообщение для пользователя.

**Пример `auth_bloc.dart` (упрощенно):**

```dart
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:your_app/auth/repository/auth_repository.dart';
import 'package:your_app/models/user.dart';

// Events
abstract class AuthEvent extends Equatable { const AuthEvent(); @override List<Object> get props => []; }
class AuthRegisterRequested extends AuthEvent { final String email, password, nickname; const AuthRegisterRequested(this.email, this.password, this.nickname); @override List<Object> get props => [email, password, nickname]; }
class AuthLoginRequested extends AuthEvent { final String email, password; const AuthLoginRequested(this.email, this.password); @override List<Object> get props => [email, password]; }
class AuthLogoutRequested extends AuthEvent {} 
class AuthCheckStatus extends AuthEvent {} 
class AuthEmailVerified extends AuthEvent {} 

// States
abstract class AuthState extends Equatable { const AuthState(); @override List<Object> get props => []; }
class AuthInitial extends AuthState {} 
class AuthLoading extends AuthState {} 
class AuthAuthenticated extends AuthState { final User user; const AuthAuthenticated(this.user); @override List<Object> get props => [user]; }
class AuthUnauthenticated extends AuthState {} 
class AuthRegistrationSuccess extends AuthState {} 
class AuthError extends AuthState { final String message; const AuthError(this.message); @override List<Object> get props => [message]; }

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final AuthRepository authRepository;

  AuthBloc({required this.authRepository}) : super(AuthInitial()) {
    on<AuthRegisterRequested>(_onRegisterRequested);
    on<AuthLoginRequested>(_onLoginRequested);
    on<AuthLogoutRequested>(_onLogoutRequested);
    on<AuthCheckStatus>(_onCheckStatus);
    on<AuthEmailVerified>(_onEmailVerified);
  }

  Future<void> _onRegisterRequested(AuthRegisterRequested event, Emitter<AuthState> emit) async {
    emit(AuthLoading());
    try {
      await authRepository.register(event.email, event.password, event.nickname);
      emit(AuthRegistrationSuccess());
    } catch (e) {
      emit(AuthError(e.toString()));
    }
  }

  Future<void> _onLoginRequested(AuthLoginRequested event, Emitter<AuthState> emit) async {
    emit(AuthLoading());
    try {
      final user = await authRepository.login(event.email, event.password);
      emit(AuthAuthenticated(user));
    } catch (e) {
      emit(AuthError(e.toString()));
    }
  }

  // ... другие обработчики событий
}
```

**2. `ProfileBloc`**

*   **События (Events):**
    *   `ProfileLoadRequested()`: Запрос на загрузку данных профиля текущего пользователя.
    *   `ProfileUpdateNicknameRequested(newNickname)`: Запрос на обновление никнейма пользователя.
*   **Состояния (States):**
    *   `ProfileInitial`: Начальное состояние.
    *   `ProfileLoading`: Состояние загрузки данных профиля.
    *   `ProfileLoaded(user)`: Данные профиля успешно загружены. Содержит объект `User`.
    *   `ProfileUpdated(user)`: Профиль успешно обновлен. Содержит обновленный объект `User`.
    *   `ProfileError(message)`: Состояние ошибки при работе с профилем.

**3. `GameSearchBloc`**

*   **События (Events):**
    *   `GameSearchStartRequested()`: Инициировать поиск оппонента для игры.
    *   `GameSearchCancelRequested()`: Отменить текущий поиск оппонента.
    *   `GameMatchFound(sessionId)`: Событие, получаемое, когда бэкенд успешно нашел матч и предоставил `sessionId`.
*   **Состояния (States):**
    *   `GameSearchInitial`: Начальное состояние.
    *   `GameSearchLoading`: Активный поиск оппонента.
    *   `GameSearchMatchFound(sessionId)`: Матч найден, готов к переходу на экран игры.
    *   `GameSearchError(message)`: Ошибка при поиске игры.

**4. `GameBloc`**

Это наиболее комплексный BLoC, отвечающий за всю игровую логику и взаимодействие с WebSocket.

*   **События (Events):**
    *   `GameStarted(session)`: Игра началась. Получено от WebSocket.
    *   `GameMoveMade(movieId, movieTitle)`: Игрок сделал ход, предложив фильм.
    *   `GameSkipTurn()`: Игрок запросил пропуск хода.
    *   `GameRequestHint()`: Игрок запросил подсказку.
    *   `GameTurnUpdated(turnData)`: Обновление состояния хода (чей ход, таймер, последний фильм). Получено от WebSocket.
    *   `GameEnded(resultData)`: Игра закончилась. Получено от WebSocket.
    *   `GameError(message)`: Ошибка в игре (от WebSocket или внутренняя).
*   **Состояния (States):**
    *   `GameInitial`: Начальное состояние.
    *   `GameLoading`: Загрузка игровой сессии.
    *   `GameInProgress(session, currentMovie, player1Skips, player2Skips, ...)`: Игра активна. Содержит всю актуальную информацию об игре.
    *   `GameMoveProcessing`: Состояние, когда ход игрока обрабатывается бэкендом.
    *   `GameHintRevealed(movie, persons)`: Подсказка показана. Содержит данные о фильме и связанных персонах.
    *   `GameEnded(result)`: Игра завершена. Содержит результаты игры.
    *   `GameError(message)`: Состояние ошибки в игре.

#### 5.4.3. Сетевые Запросы (`dio`)

Библиотека `dio` будет использоваться для всех HTTP-запросов к REST API. Она будет настроена с `Interceptor` для автоматического добавления JWT токена в заголовки запросов после успешного логина, а также для централизованной обработки ошибок.

**Пример настройки `ApiClient` с `dio`:**

```dart
import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';

class ApiClient {
  final Dio _dio = Dio();

  ApiClient() {
    _dio.options.baseUrl = 'http://localhost:8080/api/v1'; // Заменить на реальный URL бэкенда
    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        final prefs = await SharedPreferences.getInstance();
        final token = prefs.getString('jwt_token');
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        return handler.next(options);
      },
      onError: (DioException e, handler) {
        // Централизованная обработка ошибок, например, 401 Unauthorized для обновления токена или выхода
        // Также можно обрабатывать специфические ошибки бэкенда и преобразовывать их в кастомные исключения Flutter
        return handler.next(e);
      },
    ));
  }

  Dio get dio => _dio;
}
```

#### 5.4.4. WebSocket Клиент

Для WebSocket соединений будет использоваться библиотека `web_socket_channel` или аналогичная. Клиент будет отвечать за установление соединения с сервером, отправку игровых действий и получение обновлений в реальном времени.

**Пример `websocket_client.dart`:**

```dart
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:shared_preferences/shared_preferences.dart';

class WebSocketGameClient {
  WebSocketChannel? _channel;
  final String sessionId;
  final String baseUrl;

  WebSocketGameClient({required this.sessionId, required this.baseUrl});

  Stream<Map<String, dynamic>> connect() async* {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('jwt_token');
    if (token == null) {
      throw Exception('JWT token not found');
    }

    final uri = Uri.parse('$baseUrl/ws/game/$sessionId');
    _channel = WebSocketChannel.connect(
      uri,
      headers: {'Authorization': 'Bearer $token'},
    );

    await _channel!.ready;
    print('WebSocket connected to $sessionId');

    await for (final message in _channel!.stream) {
      yield jsonDecode(message) as Map<String, dynamic>;
    }
  }

  void sendMessage(Map<String, dynamic> message) {
    if (_channel != null && _channel!.sink != null) {
      _channel!.sink.add(jsonEncode(message));
    }
  }

  void disconnect() {
    _channel?.sink.close();
    print('WebSocket disconnected from $sessionId');
  }
}
```

#### 5.4.5. Навигация

Для управления навигацией в приложении рекомендуется использовать пакет `go_router` или стандартный `Navigator` с именованными маршрутами. `AuthBloc` будет играть ключевую роль в защите маршрутов, автоматически перенаправляя неаутентифицированных пользователей на экран логина.

**Пример маршрутов:**

*   `/` -> `SplashScreen` (проверяет статус аутентификации и перенаправляет)
*   `/login` -> `AuthScreen` (для входа/регистрации)
*   `/home` -> `HomeScreen` (главный экран с поиском игры)
*   `/profile` -> `ProfileScreen` (экран профиля пользователя)
*   `/game/:sessionId` -> `GamePlayScreen` (экран активной игры, требует `sessionId`)

#### 5.4.6. UI/UX Замечания

*   **Индикаторы загрузки:** Использование `CircularProgressIndicator` или скелетных экранов для улучшения восприятия во время загрузки данных или ожидания ответа от сервера.
*   **Обратная связь:** Четкие и своевременные сообщения об ошибках, успехах и текущем состоянии игры (например, через `SnackBar` или диалоговые окна).
*   **Таймер:** Визуальное отображение оставшегося времени на ход в игровом процессе, с анимацией или прогресс-баром.
*   **Адаптивный дизайн:** Приложение должно корректно отображаться и быть удобным для использования на различных размерах экранов мобильных устройств (телефоны, планшеты).
*   **Поиск фильмов:** Реализовать функцию автодополнения при вводе названия фильма для хода, используя `search/movie` эндпоинт TMDb API через бэкенд. Это значительно улучшит пользовательский опыт и снизит вероятность ошибок ввода.

## 6. Игровая Логика и Взаимодействие с TMDb

1.  **Получение данных о фильме:** При первом запросе фильма (например, при поиске игроком или при начале игры) бэкенд обращается к TMDb API для получения основной информации (название, постер, дата выхода, описание) и списка всех связанных персон (актеров, режиссеров, композиторов).
2.  **Кэширование в PostgreSQL:** Полученные данные о фильме и связанных с ним персонах сохраняются в локальной базе данных (`movies`, `persons`, `movie_person_roles`). Это позволяет избежать повторных запросов к TMDb, снизить задержки и обеспечить стабильность работы при превышении лимитов внешнего API.
3.  **Проверка связи фильма:** Когда игрок предлагает новый фильм, бэкенд выполняет следующие шаги:
    *   Проверяет, существует ли предложенный фильм в локальной БД. Если нет, запрашивает его у TMDb и кэширует.
    *   Получает список всех актеров, режиссеров и композиторов для *предыдущего* фильма из таблицы `movie_person_roles`.
    *   Получает список всех актеров, режиссеров и композиторов для *нового* предложенного фильма из таблицы `movie_person_roles`.
    *   Сравнивает эти два списка на наличие общих персон. Если найдена хотя бы одна общая персона (с любой ролью), связь считается установленной, и ход признается корректным.
    *   **Оптимизация SQL-запросов:** Для эффективной проверки связей будут использоваться сложные SQL-запросы с `JOIN` и `INTERSECT` или `EXISTS` для поиска общих персон между двумя фильмами.
4.  **Обработка подсказок:** При использовании подсказки, бэкенд извлекает из кэша (или TMDb, если данных нет) список актеров, режиссеров и композиторов текущего фильма и отправляет его клиенту через WebSocket.

## 7. План Разработки (Дорожная Карта)

Рекомендуется итеративный подход, позволяющий постепенно наращивать функциональность и получать обратную связь.

### Фаза 1: Базовый Бэкенд и Аутентификация

1.  Настройка проекта Go, инициализация PostgreSQL.
2.  Реализация базовых слоев Clean Architecture (domain, usecases, adapters, delivery).
3.  Разработка модуля аутентификации: регистрация, подтверждение email, вход, JWT-авторизация.
4.  Настройка SMTP-клиента для отправки писем верификации (использование Mailtrap для тестирования).
5.  Тестирование REST API с помощью Postman/Insomnia.

### Фаза 2: Базовый Фронтенд и Аутентификация

1.  Настройка проекта Flutter, интеграция `dio`.
2.  Реализация `AuthBloc` и UI для регистрации/входа/подтверждения email.
3.  Настройка навигации (`go_router`) с защитой маршрутов.
4.  Интеграция с бэкендом для аутентификации.

### Фаза 3: Интеграция с TMDb и Кэширование

1.  Получение API-ключа TMDb.
2.  Разработка `tmdb_client` в Go для взаимодействия с TMDb API.
3.  Реализация логики кэширования фильмов и персон в PostgreSQL.
4.  Создание вспомогательных эндпоинтов для поиска фильмов (для отладки и автодополнения).

### Фаза 4: Игровая Логика (Core Battle) и WebSockets

1.  Разработка `GameUsecase` в Go для управления всей игровой логикой.
2.  Реализация WebSocket сервера (`gorilla/websocket`) для обмена сообщениями в реальном времени.
3.  Разработка логики проверки связей между фильмами.
4.  Реализация `GameSearchBloc` и `GameBloc` во Flutter.
5.  Создание UI для поиска оппонента и основного игрового поля.
6.  Интеграция WebSocket клиента во Flutter.

### Фаза 5: Дополнительный Функционал Игры

1.  Реализация механик пропусков хода и подсказок на бэкенде и фронтенде.
2.  Добавление таймера хода с синхронизацией между клиентами.
3.  Отображение детальных результатов игры и обновление статистики профиля.

### Фаза 6: Улучшения и Развертывание

1.  Разработка UI для профиля пользователя и его редактирования.
2.  Добавление базовой статистики побед/поражений.
3.  Подготовка к развертыванию: создание Docker-контейнеров для Go-сервера и PostgreSQL.
4.  Настройка CI/CD (необязательно, для продвинутых).
5.  Оптимизация производительности и безопасности.

## 8. Важные Замечания

*   **API-ключ TMDb:** Для работы с TMDb API необходимо получить бесплатный ключ на [The Movie Database (TMDb)](https://www.themoviedb.org/documentation/api).
*   **SMTP-сервер:** Для отправки писем подтверждения email потребуется доступ к SMTP-серверу. Для разработки можно использовать сервисы типа [Mailtrap](https://mailtrap.io/) или [Mailgun](https://www.mailgun.com/).
*   **Локальная разработка:** Начинать следует с локальной разработки, используя Docker Compose для удобного запуска PostgreSQL и Go-сервера.
*   **Тестирование:** Важно писать юнит- и интеграционные тесты для бэкенда и виджет-тесты для Flutter. Это обеспечит стабильность и надежность приложения.
*   **Документация:** Поддерживать актуальную документацию по API и архитектуре в процессе разработки.

## 9. Ссылки

[1]: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html "The Clean Architecture by Uncle Bob"
