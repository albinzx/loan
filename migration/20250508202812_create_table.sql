-- +goose Up
CREATE TABLE IF NOT EXISTS loan (
    id BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'loan ID',
    amount BIGINT NOT NULL DEFAULT 0 COMMENT 'principal amount',
    rate DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT 'interest rate in percentage',
    roi DECIMAL(30,10) NOT NULL DEFAULT 0 COMMENT 'return of investment',
    borrower_id BIGINT UNSIGNED COMMENT 'borrower ID',
    agreement_letter_url VARCHAR(255) COMMENT 'aggrement letter link',
    state VARCHAR(32) COMMENT 'loan state, PROPOSED/APPROVED/INVESTED/DISBURSED',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'time when record is created/updated',
    INDEX idx_borrower (borrower_id) COMMENT 'index on borrower id'
) COMMENT = 'store all loans', ENGINE InnoDB;

CREATE TABLE IF NOT EXISTS loan_investment (
    loan_id BIGINT UNSIGNED NOT NULL COMMENT 'loan ID',
    investor_id BIGINT UNSIGNED NOT NULL COMMENT 'investor ID',
    amount BIGINT NOT NULL DEFAULT 0 COMMENT 'investment amount',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'time when record is created/updated',
    PRIMARY KEY (loan_id, investor_id) COMMENT 'primary key on loan id and investor id',
    INDEX idx_investor (investor_id) COMMENT 'index on investor id'
) COMMENT = 'store all loan investments', ENGINE InnoDB;

CREATE TABLE IF NOT EXISTS loan_approval (
    loan_id BIGINT UNSIGNED NOT NULL COMMENT 'loan ID',
    employee_id BIGINT UNSIGNED NOT NULL COMMENT 'employee ID',
    approval_date DATE NOT NULL COMMENT 'date of approval/disbursment',
    action VARCHAR(32) COMMENT 'action type, APPROVE/DISBURSE',
    document_url VARCHAR(255) COMMENT 'document link',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'time when record is created/updated',
    PRIMARY KEY (loan_id, employee_id) COMMENT 'primary key on loan id and employee id'
) COMMENT = 'store all loan approvals', ENGINE InnoDB;

CREATE TABLE IF NOT EXISTS employee (
    id BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'employee ID',
    name VARCHAR(64) NOT NULL COMMENT 'employee name',
    role VARCHAR(64) NOT NULL COMMENT 'employee role',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'time when record is created/updated'
) COMMENT = 'store all employees', ENGINE InnoDB;

CREATE TABLE IF NOT EXISTS borrower (
    id BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'borrower ID',
    name VARCHAR(64) NOT NULL COMMENT 'borrower name',
    credit_limit BIGINT NOT NULL DEFAULT 0 COMMENT 'borrower credit limit',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'time when record is created/updated'
) COMMENT = 'store all borrowers', ENGINE InnoDB;

CREATE TABLE IF NOT EXISTS investor (
    id BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'investor ID',
    name VARCHAR(64) NOT NULL COMMENT 'investor name',
    email VARCHAR(64) NOT NULL COMMENT 'investor email',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'time when record is created/updated'
) COMMENT = 'store all employees', ENGINE InnoDB;

-- +goose Down
DROP TABLE IF EXISTS loan;
DROP TABLE IF EXISTS loan_investment;
DROP TABLE IF EXISTS loan_approval;
DROP TABLE IF EXISTS employee;
DROP TABLE IF EXISTS borrower;
DROP TABLE IF EXISTS investor;
