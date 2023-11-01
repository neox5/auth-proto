import { Component } from "@angular/core";
import { OAuthService } from "angular-oauth2-oidc";
import { authConfig } from "./auth.config";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.scss"],
})
export class AppComponent {
  title = "app-frontend";

  constructor(private _oauthService: OAuthService) {
    this._oauthService.configure(authConfig);
    this._oauthService.loadDiscoveryDocumentAndTryLogin(); // This method is trigger issuer uri
  }

  login() {
    this._oauthService.initCodeFlow();
  }

  logout() {
    this._oauthService.logOut();
  }
}
