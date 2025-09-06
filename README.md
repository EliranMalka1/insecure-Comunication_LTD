# Insecure Communication LTD 

**🚨 WARNING: This is a deliberately insecure version of the project, created for educational purposes only. Do NOT use this code in a production environment or for any real-world applications. 🚨**

## Overview

This project is an intentionally vulnerable version of the Secure Communication LTD application. It is designed to serve as a practical learning tool for identifying and understanding common web security flaws. While it mirrors the architecture of the [secure version](https://github.com/EliranMalka1/secure-Comunication_LTD), it is riddled with vulnerabilities for direct comparison and hands-on learning.

The system architecture consists of:

*   **Frontend:** A React application served by an Nginx web server (running in a Docker container).
*   **Backend:** A Go API built with the Echo framework (running in a Docker container).
*   **Database:** A MySQL 8 instance for data storage.
*   **Mail Service:** MailHog for capturing and viewing outgoing emails (e.g., for password resets).

## Key Vulnerabilities

This version of the application contains several deliberately introduced security flaws.

### 1. Login & Register (SQL Injection)

The login and registration forms are vulnerable to SQL Injection. The backend queries are constructed using unsafe string concatenation, allowing attackers to manipulate the SQL and bypass authentication.

**Example Payload (Login Bypass):**
Enter the following into the username/email field, with any password:

```sql
' OR '1'='1
```

### 2. Customer Search (SQL Injection)

The customer search functionality builds SQL queries insecurely, making it susceptible to SQL Injection. An attacker can extract, modify, or delete database information by manipulating the search input.

### 3. New Customer (Stored Cross-Site Scripting - XSS)

The "New Customer" form is vulnerable to both SQL Injection and Stored XSS.

*   **SQL Injection:** The `INSERT` query is built using string concatenation, allowing malicious SQL to be executed.
*   **Stored XSS:** The "Notes" field does not sanitize user input. Any JavaScript payload saved here will be rendered directly in the browser for anyone viewing the customer's details. The React frontend uses `dangerouslySetInnerHTML` to display this content, permitting the script to execute.

**Example Payload (XSS):**
Enter the following into the "Notes" field when creating a new customer:

```html
<script>alert('XSS attack!');</script>
```

### 4. General Vulnerabilities

*   **Missing Input Validation:** Many endpoints lack proper server-side validation, trusting user input implicitly.
*   **Verbose Error Messages:** Errors often expose sensitive information, such as database schema details or internal application paths.
*   **No CSRF Protection:** The application does not implement Cross-Site Request Forgery tokens, making it vulnerable to CSRF attacks.
*   **No Rate Limiting:** Critical endpoints like login and password reset lack rate limiting, leaving them open to brute-force attacks.

## Running the Project

### Requirements

*   **Docker:** Required to run the application with Docker Compose.
*   **Node.js (Optional):** Only needed if you wish to run the frontend development server separately.

### How to Run

The entire application stack can be launched using Docker Compose.

1.  Clone the repository.
2.  Navigate to the project's root directory.
3.  Run the following command to build and start the services:

    ```bash
    docker-compose up --build
    ```

### Exposed Services

Once running, the following services will be accessible:

*   **Frontend (React App):** `http://localhost:3001`
*   **Backend (Go API):** `http://localhost:8081`
*   **Database (MySQL):** `localhost:3307`
*   **MailHog (Web UI):** `http://localhost:8026`

## Learning Goals

*   **Identify Insecure Practices:** Recognize vulnerable coding patterns and anti-patterns in a real-world context.
*   **Exploit Vulnerabilities:** Practice exploiting common web vulnerabilities like SQL Injection and Cross-Site Scripting (XSS).
*   **Compare and Contrast:** Analyze the differences in code and behavior between this insecure version and its secure counterpart.
*   **Understand Secure Design:** Gain a deeper appreciation for the principles of secure software architecture and design.

## Disclaimer

This project is **intentionally and critically insecure**. It is provided "as is" for educational use only. The maintainers are not responsible for any damage caused by the misuse of this code. Only run it in a controlled, isolated environment for training purposes. **Under no circumstances should this code be deployed to a public-facing server or used for any purpose other than security education.**
