# OAuth 2.0 / OIDC to Descope Feature Mapping

This document maps standard OAuth 2.0 and OpenID Connect (OIDC) specification concepts to their Descope implementations.

## Grant Types

| OAuth 2.0 Grant Type | RFC | Descope Implementation | Terraform Resource |
|---|---|---|---|
| **Authorization Code** | [RFC 6749 S4.1](https://datatracker.ietf.org/doc/html/rfc6749#section-4.1) | Inbound/Third-Party Apps redirect users to Descope's `/authorize` endpoint. Descope authenticates the user via Flows and returns an authorization code. | `descope_inbound_application`, `descope_third_party_application` |
| **Authorization Code + PKCE** | [RFC 7636](https://datatracker.ietf.org/doc/html/rfc7636) | Supported for public clients (non-confidential inbound apps). PKCE parameters passed in the authorization request. | `descope_inbound_application` (non_confidential_client) |
| **Client Credentials** | [RFC 6749 S4.4](https://datatracker.ietf.org/doc/html/rfc6749#section-4.4) | Access Keys provide service-to-service authentication. The access key ID and secret are used as client credentials. | `descope_access_key` |
| **Device Authorization** | [RFC 8628](https://datatracker.ietf.org/doc/html/rfc8628) | Not directly supported as a standard grant. Can be approximated with Descope Flows. | N/A |
| **Token Exchange** | [RFC 8693](https://datatracker.ietf.org/doc/html/rfc8693) | Not directly supported. Descope uses its own session management model. | N/A |
| **JWT Bearer** | [RFC 7523](https://datatracker.ietf.org/doc/html/rfc7523) | Third-Party Apps support JWT Bearer settings for validating external tokens. | `descope_third_party_application` |

## OIDC Core Concepts

| OIDC Concept | Spec Reference | Descope Implementation |
|---|---|---|
| **ID Token** | [OIDC Core S2](https://openid.net/specs/openid-connect-core-1_0.html#IDToken) | Descope issues JWTs as session tokens. Claims are customizable via JWT Templates (configured in the project resource). |
| **UserInfo Endpoint** | [OIDC Core S5.3](https://openid.net/specs/openid-connect-core-1_0.html#UserInfo) | Available at `https://api.descope.com/oauth2/v1/userinfo`. Returns user profile data based on granted scopes. |
| **Discovery** | [OIDC Discovery](https://openid.net/specs/openid-connect-discovery-1_0.html) | Each project exposes `/.well-known/openid-configuration` with endpoints, supported scopes, and signing keys. |
| **Standard Claims** | [OIDC Core S5.1](https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims) | Descope maps user attributes (name, email, phone) to standard OIDC claims. Custom claims added via JWT Templates. |
| **Dynamic Registration** | [OIDC Dynamic Registration](https://openid.net/specs/openid-connect-registration-1_0.html) | Not supported. Applications are registered via the Management API or Terraform. |

## Token Management

| Concept | Descope Implementation | Notes |
|---|---|---|
| **Access Token** | Descope session JWT | Short-lived token containing user claims, roles, permissions |
| **Refresh Token** | Descope refresh JWT | Used to obtain new session tokens without re-authentication |
| **Token Introspection** ([RFC 7662](https://datatracker.ietf.org/doc/html/rfc7662)) | Validate tokens via Descope SDKs or JWKS endpoint | Standard JWKS-based validation |
| **Token Revocation** ([RFC 7009](https://datatracker.ietf.org/doc/html/rfc7009)) | Logout API / session management | Per-user or per-session revocation via Management API |
| **Token Lifetime** | Configurable per inbound app | Session settings include token expiration and refresh token lifetime |

## Client Types

| OAuth 2.0 Client Type | Descope Equivalent | Terraform Resource |
|---|---|---|
| **Confidential Client** | Inbound App with `non_confidential_client = false` | `descope_inbound_application` |
| **Public Client** | Inbound App with `non_confidential_client = true` | `descope_inbound_application` |
| **Third-Party Client** | Third-Party Application (with consent flow) | `descope_third_party_application` |
| **First-Party Client** | Direct Descope SDK integration (no OAuth needed) | N/A (SDK-based) |

## Scopes and Claims

| Concept | Descope Implementation |
|---|---|
| **`openid` scope** | Implicit in all OIDC flows. Returns ID token with `sub` claim. |
| **`profile` scope** | Maps to user display name, given name, family name, picture. |
| **`email` scope** | Maps to user email and email_verified. |
| **`phone` scope** | Maps to user phone and phone_verified. |
| **Custom scopes** | Defined as Permission Scopes on Inbound/Third-Party Apps. Mapped to Descope RBAC permissions. |
| **Custom claims** | Added via JWT Templates in the project configuration. |

## Federation and SSO

| Standard | Descope Feature | Terraform Resource |
|---|---|---|
| **SAML 2.0 SP** | Descope acts as SAML Service Provider. Tenants configure their SAML IdP (Okta, Azure AD, etc.). | `descope_sso` (SAML settings) |
| **OIDC RP** | Descope acts as OIDC Relying Party. Tenants configure their OIDC provider. | `descope_sso` (OIDC settings) |
| **SAML IdP** | Descope acts as SAML Identity Provider via SSO Applications. | `descope_sso_application` (blocked - enterprise) |
| **OIDC OP** | Descope acts as OIDC Provider via Inbound Applications. | `descope_inbound_application` |
| **Social Login** | Outbound Apps connect to social providers (Google, GitHub, etc.). | `descope_outbound_application` |

## Authorization

| Concept | Descope Feature | Terraform Resource |
|---|---|---|
| **OAuth 2.0 Scopes** | Permission Scopes on Inbound/Third-Party Apps | `descope_inbound_application`, `descope_third_party_application` |
| **RBAC** | Roles and Permissions (project-level and tenant-level) | `descope_role`, `descope_permission` |
| **ReBAC / FGA** | Fine-Grained Authorization with Zanzibar-style relation schema | `descope_fga_schema`, `descope_fga_check` (data source) |
| **ABAC** | User custom attributes + tenant attributes for attribute-based decisions | `descope_project` (user attributes config) |

## Descope OIDC Endpoints

For a project with ID `P123`, the OIDC endpoints are:

| Endpoint | URL |
|---|---|
| Discovery | `https://api.descope.com/P123/.well-known/openid-configuration` |
| Authorization | `https://api.descope.com/oauth2/v1/authorize` |
| Token | `https://api.descope.com/oauth2/v1/token` |
| UserInfo | `https://api.descope.com/oauth2/v1/userinfo` |
| JWKS | `https://api.descope.com/P123/.well-known/jwks.json` |

See [Descope OIDC Endpoints Quickstart](https://docs.descope.com/getting-started/oidc-endpoints) for details.
