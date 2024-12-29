# Security Policy

## Supported Versions

The following versions of this project are actively supported for security updates:

| Version     | Supported          |
|-------------|--------------------|
| 1.x         | ✅ Yes            |

Older versions are no longer supported. Users are encouraged to update to the latest version for the best security and functionality.

---

## Reporting a Vulnerability

If you discover a **security vulnerability** in this project, please follow these steps:

1. **Do not open a public issue.** Security vulnerabilities should not be disclosed publicly until they are reviewed and patched.
2. Contact the maintainers via [pyd773@gmail.com](mailto:pyd773@gmail.com) with a detailed description of the vulnerability, including:

   - Steps to reproduce the issue.
   - Potential impact.
   - Any suggested fixes or patches, if applicable.

We aim to respond to security issues within **48 hours** and provide a mitigation plan within **7 days**.

---

## Dependencies and Third-Party Libraries

This project depends on several third-party libraries. We monitor their security updates regularly. Here are some key dependencies and how we address their vulnerabilities:

1. **Direct Dependencies:**
   - `github.com/gin-gonic/gin`
   - `github.com/google/generative-ai-go`
   - `google.golang.org/api`
   - `github.com/gin-gonic/gin`
   - `github.com/golang-jwt/jwt/v4`

   **Mitigation:**
   - Dependencies are updated frequently to their latest stable versions.
   - Vulnerabilities in these libraries are addressed by applying patches or upgrading as soon as fixes are available.

2. **Indirect Dependencies:**
   - Libraries such as `cloud.google.com/go`, `google.golang.org/protobuf`, `github.com/bytedance/sonic`, `golang.org/x/crypto`, and `go.opentelemetry.io` are monitored using automated tooling.

   **Mitigation:**
   - Indirect dependencies are audited using tool [GoSec](https://github.com/securego/gosec).

---

## Security Best Practices

1. **Environment Variables and Secrets Management:**
   - Ensure sensitive information (e.g., API keys, tokens) is stored securely using environment variables.
   - Do not commit `.env` files or sensitive data into the repository.

2. **Audit Code for Common Vulnerabilities:**
   - Regularly use tools like `gosec` or `staticcheck` to scan for vulnerabilities in the codebase.
   - Run dependency checks using tools like `go list -m all` and validate against known vulnerability databases (e.g., CVE).

3. **Runtime Security:**
   - Run the application in a secured environment with limited access and permissions.
   - Use containerization tools like Docker with appropriate security configurations.

---

## Contact Information

If you have questions or concerns regarding security in this project, please reach out to the maintainers at [pyd773@gmail.com](mailto:pyd773@gmail.com).
