import type {
  CheckParametersOut,
  Question,
  QuestionHeader,
  StartSessionOut,
  TrivialConfig,
  TrivialConfigExt
} from "./api_gen";
import { AbstractAPI } from "./api_gen";

function arrayBufferToString(buffer: ArrayBuffer) {
  const uintArray = new Uint8Array(buffer);
  const encodedString = String.fromCharCode.apply(null, Array.from(uintArray));
  return decodeURIComponent(escape(encodedString));
}

class Controller extends AbstractAPI {
  public onError?: (kind: string, htmlError: string) => void;
  public showMessage?: (message: string) => void;

  protected onSuccessEditorDuplicateQuestion(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question dupliquée.");
    }
  }
  protected onSuccessLaunchSessionTrivialPoursuit(data: TrivialConfig): void {
    this.inRequest = false;
  }
  protected onSuccessGetTrivialPoursuit(data: TrivialConfigExt[] | null): void {
    this.inRequest = false;
  }
  protected onSuccessCreateTrivialPoursuit(data: TrivialConfigExt): void {
    this.inRequest = false;
  }
  protected onSuccessUpdateTrivialPoursuit(data: TrivialConfigExt): void {
    this.inRequest = false;
  }
  protected onSuccessDeleteTrivialPoursuit(data: any): void {
    this.inRequest = false;
  }
  protected onSuccessLaunchSession(data: TrivialConfig): void {
    console.log(`Game started at ${data.LaunchSessionID}`);
    this.inRequest = false;
  }
  inRequest = false;

  getURL(endpoint: string) {
    return this.baseUrl + endpoint;
  }

  handleError(error: any): void {
    this.inRequest = false;

    let kind: string,
      messageHtml: string,
      code = null;
    if (error.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      kind = `Erreur côté serveur`;
      code = error.response.status;

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
    console.log("launching request");
    this.inRequest = true;
  }

  protected onSuccessEditorStartSession(data: StartSessionOut): void {
    console.log(data);
    this.inRequest = false;
  }

  protected onSuccessEditorSaveAndPreview(data: any): void {
    this.inRequest = false;
  }

  protected onSuccessEditorCheckParameters(data: CheckParametersOut): void {
    this.inRequest = false;
  }

  protected onSuccessEditorSearchQuestions(
    data: QuestionHeader[] | null
  ): void {
    this.inRequest = false;
  }
  protected onSuccessEditorCreateQuestion(data: Question): void {
    this.inRequest = false;
  }
  protected onSuccessEditorUpdateTags(data: any): void {
    this.inRequest = false;
  }
  protected onSuccessEditorGetTags(data: any): void {
    this.inRequest = false;
  }
  protected onSuccessEditorGetQuestion(data: any): void {
    this.inRequest = false;
  }
  protected onSuccessEditorDeleteQuestion(data: any): void {
    this.inRequest = false;
  }
  protected onSuccessEditorPausePreview(data: any): void {
    this.inRequest = false;
  }
}

const localhost = "http://localhost:1323";
const buildMode = import.meta.env.DEV ? "dev" : "prod";

// when building for production, the mode for the preview client
// may actually still be "dev", for instance when served in test
export const PreviewMode =
  buildMode == "dev"
    ? "dev"
    : window.location.origin == localhost
    ? "dev"
    : "prod";

export const controller = new Controller(
  buildMode == "dev" ? localhost : window.location.origin,
  "",
  {}
);
