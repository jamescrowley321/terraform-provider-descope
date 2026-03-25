# Authentication and Authorization Flows

This document provides sequence diagrams for Descope's authentication and authorization flows.

## OAuth 2.0 Authorization Code Flow (Descope as IdP)

When an external application (Inbound App) uses Descope as its identity provider:

```mermaid
sequenceDiagram
    participant User
    participant App as Client App<br/>(Inbound App)
    participant Descope as Descope<br/>(/authorize)
    participant Flow as Descope Flow<br/>(Login UI)
    participant Token as Descope<br/>(/token)

    User->>App: Access protected resource
    App->>Descope: Redirect to /authorize<br/>(client_id, redirect_uri, scope, state)
    Descope->>Flow: Present login flow
    Flow->>User: Show login screen
    User->>Flow: Authenticate (password/OTP/magic link/SSO)
    Flow->>Descope: Authentication successful
    Note over Descope: If scopes require consent
    Descope->>User: Show consent screen
    User->>Descope: Approve scopes
    Descope->>App: Redirect to redirect_uri<br/>(code, state)
    App->>Token: POST /token<br/>(code, client_id, client_secret)
    Token->>App: Return tokens<br/>(access_token, id_token, refresh_token)
    App->>User: Grant access
```

## Third-Party App Consent Flow

When a third-party application requests access to user data:

```mermaid
sequenceDiagram
    participant User
    participant ThirdParty as Third-Party App
    participant Descope as Descope
    participant Consent as Consent Screen

    ThirdParty->>Descope: Redirect to /authorize<br/>(client_id, scopes)
    Descope->>User: Authenticate user
    User->>Descope: Credentials
    Descope->>Consent: Show permission request
    Note over Consent: "App X wants to:<br/>- Read your profile<br/>- Access your email"
    Consent->>User: Display requested scopes
    User->>Consent: Approve / Deny
    alt Approved
        Consent->>Descope: Record consent
        Descope->>ThirdParty: Authorization code
        ThirdParty->>Descope: Exchange for tokens
    else Denied
        Descope->>ThirdParty: Error: access_denied
    end
```

## SSO / SAML Federation Flow

When a tenant has SAML SSO configured:

```mermaid
sequenceDiagram
    participant User
    participant App as Application
    participant Descope as Descope (SP)
    participant IdP as Enterprise IdP<br/>(Okta/Azure AD)

    User->>App: Click "Sign in with SSO"
    App->>Descope: Initiate SSO<br/>(tenant_id)
    Descope->>IdP: SAML AuthnRequest
    IdP->>User: Show enterprise login
    User->>IdP: Enter corporate credentials
    IdP->>User: MFA challenge (if configured)
    User->>IdP: Complete MFA
    IdP->>Descope: SAML Response<br/>(assertions, attributes)
    Descope->>Descope: Validate signature,<br/>map attributes to user,<br/>apply role mappings
    Descope->>App: Session tokens<br/>(JWT with roles/permissions)
    App->>User: Authenticated session
```

## Standard Authentication Flow

Direct authentication without OAuth (SDK-based):

```mermaid
sequenceDiagram
    participant User
    participant App as Application<br/>(Descope SDK)
    participant Descope as Descope API

    alt Password Authentication
        User->>App: Enter email + password
        App->>Descope: POST /auth/password/signin
        Descope->>Descope: Validate credentials,<br/>check password policy
        Descope->>App: Session JWT + Refresh JWT
    else OTP (Email/SMS)
        User->>App: Enter email/phone
        App->>Descope: POST /auth/otp/signin
        Descope->>User: Send OTP code
        User->>App: Enter OTP code
        App->>Descope: POST /auth/otp/verify
        Descope->>App: Session JWT + Refresh JWT
    else Magic Link
        User->>App: Enter email
        App->>Descope: POST /auth/magiclink/signin
        Descope->>User: Send magic link email
        User->>App: Click magic link
        App->>Descope: Verify magic link token
        Descope->>App: Session JWT + Refresh JWT
    end

    App->>User: Authenticated session
```

## Authorization Model

How the three authorization layers compose:

```mermaid
graph TB
    subgraph "Authorization Decision"
        Decision[Access Granted?]
    end

    subgraph "Layer 1: RBAC"
        User[User] --> ProjectRoles[Project Roles]
        User --> TenantRoles[Tenant Roles]
        ProjectRoles --> Permissions[Permissions]
        TenantRoles --> Permissions
    end

    subgraph "Layer 2: FGA / ReBAC"
        FGASchema[FGA Schema<br/>types + relations] --> Relations[Relations<br/>user:alice#owner@doc:1]
        Relations --> FGACheck{FGA Check<br/>Is alice viewer of doc:1?}
    end

    subgraph "Layer 3: Lists"
        IPList[IP Allowlist] --> IPCheck{IP in list?}
        TextList[Text Denylist] --> TextCheck{Domain blocked?}
    end

    Permissions --> Decision
    FGACheck --> Decision
    IPCheck --> Decision
    TextCheck --> Decision
```

## Terraform Resource Flow

Which Terraform resources configure each part of the auth infrastructure:

```mermaid
graph LR
    subgraph "Terraform-Managed Infrastructure"
        Project[descope_project]
        Tenant[descope_tenant]
        Role[descope_role]
        Perm[descope_permission]
        SSO[descope_sso]
        InApp[descope_inbound_application]
        TPA[descope_third_party_application]
        OutApp[descope_outbound_application]
        PwdSettings[descope_password_settings]
        FGA[descope_fga_schema]
        List[descope_list]
        AK[descope_access_key]
        MK[descope_management_key]
    end

    subgraph "Runtime / Console-Managed"
        Users[Users]
        Flows[Auth Flows]
        Sessions[Sessions]
    end

    Project --> Tenant
    Project --> Role
    Project --> Perm
    Tenant --> SSO
    Role --> Perm
    Project --> InApp
    Project --> TPA
    Project --> OutApp
    Project --> PwdSettings
    Project --> FGA
    Project --> List
    Project --> AK
    Project --> MK

    InApp -.-> Users
    Flows -.-> Users
    Users -.-> Sessions
```

## Further Reading

- [Descope Authentication Methods](https://docs.descope.com/auth-methods)
- [Descope OIDC Endpoints](https://docs.descope.com/getting-started/oidc-endpoints)
- [Descope Inbound Apps](https://docs.descope.com/identity-federation/inbound-apps/using-inbound-apps)
- [Descope SSO Configuration](https://docs.descope.com/tenant-management/sso/how-authorization-works-with-sso-providers)
- [Descope RBAC](https://docs.descope.com/authorization/role-based-access-control)
- [Descope FGA/ReBAC](https://docs.descope.com/authorization/rebac)
