import { AuthConfig } from "angular-oauth2-oidc";

export const authConfig: AuthConfig = {
  issuer: "http://localhost:8080/realms/keycloak-test",
  redirectUri: window.location.origin, // we do this to not have it hardcoded
  clientId: "angular-pkce",
  responseType: "code",
  strictDiscoveryDocumentValidation: true, // only relevant for the OAuth2 OIDC library; if True, when the Issuer URIs are all the same
  scope: "openid profile email offline_access",
  showDebugInformation: true
}