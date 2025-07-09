import { devLogMeta } from "@/env";
import type { TeacherSettings } from "./api_gen";
import { AbstractAPI, MatiereTag } from "./api_gen";

function arrayBufferToString(buffer: ArrayBuffer) {
  const uintArray = new Uint8Array(buffer);
  const encodedString = String.fromCharCode.apply(null, Array.from(uintArray));
  return decodeURIComponent(escape(encodedString));
}

class Controller extends AbstractAPI {
  public isLoggedIn = false;

  public onError: (kind: string, htmlError: string) => void = (_, __) => {};
  public showMessage: (message: string, color?: string) => void = (_, __) => {};

  public settings: TeacherSettings = {
    Mail: "",
    HasEditorSimplified: false,
    Password: "",
    Contact: { Name: "", URL: "" },
    FavoriteMatiere: MatiereTag.Autre,
  };

  logout() {
    this.isLoggedIn = false;
    this.authToken = "";
  }

  getToken() {
    return this.authToken;
  }

  inRequest = false;

  getURL(endpoint: string) {
    return this.baseURL + endpoint;
  }

  handleError(error: any): void {
    this.inRequest = false;

    let kind: string, messageHtml: string;
    //   code = null;
    if (error.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      kind = `Erreur côté serveur`;
      //   code = error.response.status;

      messageHtml = error.response.data.message;
      if (messageHtml) {
        messageHtml = "<i>" + messageHtml + "</i>";
      } else {
        try {
          const json = arrayBufferToString(error.response.data);
          messageHtml = JSON.parse(json).message;
        } catch (error) {
          messageHtml = `Le format d'erreur du serveur n'a pu être décodé.<br/>
        Détails : <i>${error}</i>`;
        }
      }
    } else if (error.request) {
      // The request was made but no response was received
      // `error.request` is an instance of XMLHttpRequest in the browser and an instance of
      // http.ClientRequest in node.js
      kind = "Aucune réponse du serveur";
      messageHtml =
        "La requête a bien été envoyée, mais le serveur n'a donné aucune réponse...";
    } else {
      // Something happened in setting up the request that triggered an Error
      kind = "Erreur du client";
      messageHtml = `La requête n'a pu être mise en place. <br/>
                  Détails :  ${error.message} `;
    }

    if (this.onError) {
      this.onError(kind, messageHtml);
    }
  }

  startRequest(): void {
    this.inRequest = true;
  }

  async ensureSettings() {
    const res = await this.TeacherGetSettings();
    if (res == undefined) return;
    this.settings = res;
  }
}

const localhost = "http://localhost:1323";

/** `IsDev` is true when the client app is served in dev mode */
export const IsDev = import.meta.env.DEV;

// when building for production, the mode for the preview client
// may actually still be "dev", for instance when served in test
export const PreviewMode = IsDev
  ? "dev"
  : window.location.origin == localhost
  ? "dev"
  : "prod";

export function isInscriptionValidated() {
  return window.location.search.includes("show-success-inscription");
}

export const controller = new Controller(
  IsDev ? localhost : window.location.origin,
  IsDev ? devLogMeta.Token : ""
);
