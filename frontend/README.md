# Insecure Communication LTD - Intentionally Insecure Frontend

** Warning: This application is intentionally insecure and for educational purposes only. Do not use in a production environment.**

## Overview

This project is the insecure frontend for the Secure Communication LTD application. It is built with React + Vite and served using Nginx inside a Docker container.

This frontend communicates with a vulnerable backend service running on port `8081`. It is designed to mirror the functionality of the secure frontend but with intentionally weakened or disabled client-side validations to demonstrate common frontend vulnerabilities.

## Key Vulnerabilities

### 1. Login Form

*   **No Input Sanitization:** User inputs are not sanitized or trimmed. Raw payloads are sent directly to the backend.
*   **No Client-Side Validation:** The username/email field lacks validation, allowing malicious strings to be sent to the backend.
    *   **Example SQLi Payload:** `''' OR ''1''=''1'`

### 2. Register Form

*   **Weak/No Client-Side Validation:** Similar to the login form, the registration form has minimal to no validation on its input fields.
*   **SQL Injection:** A user can inject SQL payloads into the username and email fields.
    *   **Example SQLi Payload:** `test@test.com'''; DROP TABLE users; --`

### 3. Customer Search

*   **Direct Input Pass-through:** The search input is passed directly to the backend without any escaping or sanitization.
*   **Reflected SQLi:** When combined with the vulnerable backend, this allows for reflected SQL injection attacks.

### 4. New Customer Form

*   **No Validation or Sanitization:** The form for creating a new customer has no input validation or sanitization.
*   **Stored XSS:** The `notes` field is rendered using `dangerouslySetInnerHTML`, which allows for the execution of stored Cross-Site Scripting (XSS) attacks.
    *   **Example XSS Payload:** `<script>alert('''XSS''')</script>`

## Running the Frontend

### Requirements

*   Docker
*   Docker Compose

### Instructions

1.  Start the application using Docker Compose:
    ```bash
    docker-compose up --build
    ```
2.  The frontend will be served on **http://localhost:3001**.

## Learning Goals

*   **Observe Amplification:** Understand how poor frontend practices can amplify backend vulnerabilities.
*   **Frontend vs. Backend Security:** Demonstrate why client-side validation is not a substitute for backend security.
*   **Comparative Analysis:** Compare the insecure frontend implementation with its secure counterpart to understand best practices.

## Disclaimer

**This project is strictly for educational purposes. It is designed to be run in a controlled environment to learn about and demonstrate web security vulnerabilities. Under no circumstances should this code be used in a production environment.**
