CREATE TABLE IF NOT EXISTS SchemaMigration (
    Version INTEGER PRIMARY KEY,
    AppliedAt TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Users (
    Id TEXT PRIMARY KEY,
    Email TEXT NOT NULL UNIQUE,
    PasswordHash TEXT NOT NULL,
    DisplayName TEXT NOT NULL,
    IsAdmin INTEGER NOT NULL DEFAULT 0 CHECK (IsAdmin IN (0, 1)),
    CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Preferences (
    UserId TEXT PRIMARY KEY REFERENCES Users(Id) ON DELETE CASCADE,
    Servings INTEGER NOT NULL DEFAULT 2 CHECK (Servings > 0),
    CuisinesJson TEXT NOT NULL DEFAULT '[]',
    AllergiesJson TEXT NOT NULL DEFAULT '[]',
    AvailableMinutes INTEGER NOT NULL DEFAULT 45 CHECK (AvailableMinutes >= 0),
    CookwareJson TEXT NOT NULL DEFAULT '[]',
    PantryIngredientIdsJson TEXT NOT NULL DEFAULT '[]'
);

CREATE TABLE IF NOT EXISTS Ingredients (
    Id TEXT PRIMARY KEY,
    Name TEXT NOT NULL UNIQUE,
    Category TEXT NOT NULL,
    DefaultUnit TEXT NOT NULL CHECK (DefaultUnit IN ('g', 'kg', 'ml', 'l', '个', '份', '勺', '茶匙')),
    IsPantryStaple INTEGER NOT NULL DEFAULT 0 CHECK (IsPantryStaple IN (0, 1)),
    SubstituteGroup TEXT NOT NULL DEFAULT '',
    StoreSkuMapping TEXT
);

CREATE TABLE IF NOT EXISTS IngredientAliases (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    IngredientId TEXT NOT NULL REFERENCES Ingredients(Id) ON DELETE CASCADE,
    Alias TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Recipes (
    Id TEXT PRIMARY KEY,
    Slug TEXT NOT NULL UNIQUE,
    Name TEXT NOT NULL,
    Cuisine TEXT NOT NULL,
    Description TEXT NOT NULL,
    PrepMinutes INTEGER NOT NULL CHECK (PrepMinutes >= 0),
    CookMinutes INTEGER NOT NULL CHECK (CookMinutes >= 0),
    Servings INTEGER NOT NULL CHECK (Servings > 0),
    Difficulty INTEGER NOT NULL CHECK (Difficulty BETWEEN 1 AND 5),
    ImagePath TEXT NOT NULL DEFAULT '',
    Status TEXT NOT NULL CHECK (Status IN ('draft', 'published', 'archived')),
    AllergensJson TEXT NOT NULL DEFAULT '[]',
    CookwareJson TEXT NOT NULL DEFAULT '[]',
    CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS RecipeIngredients (
    RecipeId TEXT NOT NULL REFERENCES Recipes(Id) ON DELETE CASCADE,
    IngredientId TEXT NOT NULL REFERENCES Ingredients(Id) ON DELETE RESTRICT,
    Quantity REAL NOT NULL CHECK (Quantity > 0),
    Unit TEXT NOT NULL CHECK (Unit IN ('g', 'kg', 'ml', 'l', '个', '份', '勺', '茶匙')),
    ServingFactor REAL NOT NULL DEFAULT 1.0 CHECK (ServingFactor > 0),
    Preparation TEXT NOT NULL DEFAULT '',
    Required INTEGER NOT NULL DEFAULT 1 CHECK (Required IN (0, 1)),
    PRIMARY KEY (RecipeId, IngredientId)
);

CREATE TABLE IF NOT EXISTS RecipeSteps (
    RecipeId TEXT NOT NULL REFERENCES Recipes(Id) ON DELETE CASCADE,
    StepOrder INTEGER NOT NULL CHECK (StepOrder > 0),
    Title TEXT NOT NULL,
    Instruction TEXT NOT NULL,
    DurationSeconds INTEGER NOT NULL DEFAULT 0 CHECK (DurationSeconds >= 0),
    HasTimer INTEGER NOT NULL DEFAULT 0 CHECK (HasTimer IN (0, 1)),
    PRIMARY KEY (RecipeId, StepOrder)
);

CREATE TABLE IF NOT EXISTS MealPlans (
    Id TEXT PRIMARY KEY,
    UserId TEXT NOT NULL REFERENCES Users(Id) ON DELETE CASCADE,
    PlanDate TEXT NOT NULL,
    CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (UserId, PlanDate)
);

CREATE TABLE IF NOT EXISTS MealPlanItems (
    Id TEXT PRIMARY KEY,
    MealPlanId TEXT NOT NULL REFERENCES MealPlans(Id) ON DELETE CASCADE,
    RecipeId TEXT NOT NULL REFERENCES Recipes(Id) ON DELETE RESTRICT,
    Servings INTEGER NOT NULL CHECK (Servings > 0),
    SortOrder INTEGER NOT NULL DEFAULT 0 CHECK (SortOrder >= 0)
);

CREATE TABLE IF NOT EXISTS Reminders (
    Id TEXT PRIMARY KEY,
    UserId TEXT NOT NULL REFERENCES Users(Id) ON DELETE CASCADE,
    MealPlanItemId TEXT REFERENCES MealPlanItems(Id) ON DELETE CASCADE,
    ReminderType TEXT NOT NULL CHECK (ReminderType IN ('defrost', 'marinate', 'prep', 'cook')),
    ScheduledAt TEXT NOT NULL,
    IsEnabled INTEGER NOT NULL DEFAULT 1 CHECK (IsEnabled IN (0, 1)),
    CompletedAt TEXT
);

CREATE TABLE IF NOT EXISTS CookingSessions (
    Id TEXT PRIMARY KEY,
    UserId TEXT NOT NULL REFERENCES Users(Id) ON DELETE CASCADE,
    RecipeId TEXT NOT NULL REFERENCES Recipes(Id) ON DELETE RESTRICT,
    CurrentStepOrder INTEGER NOT NULL DEFAULT 1 CHECK (CurrentStepOrder > 0),
    State TEXT NOT NULL CHECK (State IN ('active', 'paused', 'completed', 'abandoned')),
    StartedAt TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Feedback (
    Id TEXT PRIMARY KEY,
    UserId TEXT NOT NULL REFERENCES Users(Id) ON DELETE CASCADE,
    RecipeId TEXT NOT NULL REFERENCES Recipes(Id) ON DELETE RESTRICT,
    Outcome TEXT NOT NULL CHECK (Outcome IN ('made', 'skipped')),
    TagsJson TEXT NOT NULL DEFAULT '[]',
    Comment TEXT NOT NULL DEFAULT '',
    CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS MediaAssets (
    Id TEXT PRIMARY KEY,
    RecipeId TEXT REFERENCES Recipes(Id) ON DELETE CASCADE,
    Path TEXT NOT NULL,
    Source TEXT NOT NULL,
    License TEXT NOT NULL,
    Sha256 TEXT NOT NULL,
    Width INTEGER NOT NULL CHECK (Width > 0),
    Height INTEGER NOT NULL CHECK (Height > 0)
);

CREATE TABLE IF NOT EXISTS RefreshTokens (
    Id TEXT PRIMARY KEY,
    UserId TEXT NOT NULL REFERENCES Users(Id) ON DELETE CASCADE,
    TokenHash TEXT NOT NULL UNIQUE,
    ExpiresAt TEXT NOT NULL,
    RevokedAt TEXT
);
