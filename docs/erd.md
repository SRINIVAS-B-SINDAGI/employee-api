# Employee Management API - Entity Relationship Diagram

## Database Schema

```
┌─────────────────────────────────────┐
│              USERS                  │
├─────────────────────────────────────┤
│ id            UUID [PK]             │
│ email         VARCHAR(255) [UNIQUE] │
│ password_hash VARCHAR(255)          │
│ created_at    TIMESTAMPTZ           │
│ updated_at    TIMESTAMPTZ           │
└─────────────────────────────────────┘


┌─────────────────────────────────────┐
│            EMPLOYEES                │
├─────────────────────────────────────┤
│ id            UUID [PK]             │
│ full_name     VARCHAR(255)          │
│ job_title     VARCHAR(100) [IDX]    │
│ country       VARCHAR(100) [IDX]    │
│ gross_salary  DECIMAL(15,2)         │
│ created_at    TIMESTAMPTZ           │
│ updated_at    TIMESTAMPTZ           │
│ deleted_at    TIMESTAMPTZ [IDX]     │
└─────────────────────────────────────┘
```

## Tables Description

### Users Table
Stores user authentication information.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique identifier |
| email | VARCHAR(255) | NOT NULL, UNIQUE | User email address |
| password_hash | VARCHAR(255) | NOT NULL | Bcrypt hashed password |
| created_at | TIMESTAMPTZ | DEFAULT CURRENT_TIMESTAMP | Record creation time |
| updated_at | TIMESTAMPTZ | DEFAULT CURRENT_TIMESTAMP | Last update time |

### Employees Table
Stores employee information with soft delete support.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | Unique identifier |
| full_name | VARCHAR(255) | NOT NULL | Employee full name |
| job_title | VARCHAR(100) | NOT NULL, INDEX | Job position |
| country | VARCHAR(100) | NOT NULL, INDEX | Country of employment |
| gross_salary | DECIMAL(15,2) | NOT NULL, CHECK >= 0 | Gross annual salary |
| created_at | TIMESTAMPTZ | DEFAULT CURRENT_TIMESTAMP | Record creation time |
| updated_at | TIMESTAMPTZ | DEFAULT CURRENT_TIMESTAMP | Last update time |
| deleted_at | TIMESTAMPTZ | INDEX, NULLABLE | Soft delete timestamp |

## Indexes

| Table | Index Name | Column(s) | Purpose |
|-------|------------|-----------|---------|
| users | users_email_key | email | Unique constraint, login lookup |
| employees | idx_employees_country | country | Salary metrics by country |
| employees | idx_employees_job_title | job_title | Salary metrics by job title |
| employees | idx_employees_deleted_at | deleted_at | Soft delete filtering |

## Tax Deduction Rules

| Country | Tax Rate | Calculation |
|---------|----------|-------------|
| India | 10% | gross_salary * 0.10 |
| United States | 12% | gross_salary * 0.12 |
| Other | 0% | 0 |

## Notes

- All primary keys use UUID v4 for distributed system compatibility
- Timestamps use TIMESTAMPTZ for timezone awareness
- Soft delete pattern allows data recovery and audit trails
- Indexes optimize the most common query patterns (country/job title metrics)
- Password hashing uses bcrypt with default cost factor
