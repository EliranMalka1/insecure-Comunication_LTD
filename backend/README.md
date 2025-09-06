# Insecure Communication LTD - Intentionally Insecure Backend

**&#9888;&#65039; Warning: This is an intentionally insecure backend application. It is designed for educational purposes only and should never be used in a production environment.**

## Overview

This project is the insecure backend for the Secure Communication LTD application. It is built with the Go (Echo framework) and is designed to mirror the functionality of the secure backend, but with deliberately introduced vulnerabilities. This allows developers to learn about common security flaws in a hands-on manner.

The backend connects to a MySQL 8 database and uses MailHog for email testing. Both services are managed via Docker.

## Key Vulnerabilities in the Backend

This backend is riddled with security vulnerabilities. Here are some of the key issues you can find and exploit:

### 1. Login - SQL Injection

The login endpoint is vulnerable to SQL injection. By using a classic `' OR '1'='1'` payload, you can bypass the authentication check.

**Example Payload:**

```json
{
  "username": "' OR '1'='1",
  "password": "password"
}
```

### 2. Register - SQL Injection

The user registration endpoint is also vulnerable to SQL injection. The `INSERT` query is built using unsafe string concatenation, allowing an attacker to manipulate the query.

### 3. Customer Search - SQL Injection

The customer search functionality is vulnerable to SQL injection due to unsafe query building.

### 4. New Customer - SQL Injection and Stored XSS

The endpoint for creating a new customer has two major vulnerabilities:

*   **SQL Injection:** Similar to the registration endpoint, the `INSERT` query is vulnerable to SQL injection.
*   **Stored XSS:** The `notes` field does not sanitize user input. This allows an attacker to inject malicious scripts that will be executed when the customer's details are viewed.

**Example XSS Payload:**

```json
{
  "name": "Malicious Customer",
  "notes": "<script>alert('XSS');</script>"
}
```

### General Issues

*   **No Input Validation:** The application does not properly validate user input, which is a root cause of many of the vulnerabilities.
*   **Detailed Error Messages:** The backend returns detailed database errors in JSON responses, which can leak information about the database schema.
*   **No Rate Limiting:** There is no rate limiting on any of the endpoints, making the application vulnerable to brute-force attacks.
*   **Weak/Missing Access Controls:** The application has weak or missing access controls, allowing unauthorized users to access sensitive data and functionality.

## Running the Backend

To run the backend, you will need to have Docker and Docker Compose installed.

1.  **Start the application:**

    ```bash
    docker-compose up --build
    ```

2.  **Access the services:**

    *   **Backend:** `http://localhost:8081`
    *   **MySQL:** `localhost:3307`
    *   **MailHog:** `http://localhost:8026`

## Learning Goals

*   Observe insecure coding patterns in a real-world Go application.
*   Practice exploiting common backend vulnerabilities like SQL Injection and XSS.
*   Compare the insecure code in this repository with the secure version to understand how to prevent these vulnerabilities.

## Disclaimer

**This project is for educational purposes only. It is not intended for use in a production environment. The vulnerabilities are intentional and are meant to be used for learning and training purposes in a controlled environment.**
