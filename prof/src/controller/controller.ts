import type {
  CheckMissingQuestionsOut,
  CheckParametersOut,
  LaunchSessionOut,
  ListQuestionsOut,
  LogginOut,
  Question,
  StartSessionOut,
  TrivialConfigExt,
} from "./api_gen";
import { AbstractAPI, GroupStrategyKind } from "./api_gen";

function arrayBufferToString(buffer: ArrayBuffer) {
  const uintArray = new Uint8Array(buffer);
  const encodedString = String.fromCharCode.apply(null, Array.from(uintArray));
  return decodeURIComponent(escape(encodedString));
}

class Controller extends AbstractAPI {
  protected onSuccessAskInscription(data: any): void {
    this.inRequest = false;
  }
  protected onSuccessValidateInscription(data: any): void {
    this.inRequest = false;
  }
  protected onSuccessLoggin(data: LogginOut): void {
    this.inRequest = false;
  }

  protected onSuccessDuplicateTrivialPoursuit(data: TrivialConfigExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Session dupliquée.");
    }
  }
  protected onSuccessCheckMissingQuestions(
    data: CheckMissingQuestionsOut
  ): void {
    this.inRequest = false;
  }

  protected onSuccessStopSessionTrivialPoursuit(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Session terminée.");
    }
  }
  protected onSuccessEditorDuplicateQuestion(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question dupliquée.");
    }
  }

  public onError?: (kind: string, htmlError: string) => void;
  public showMessage?: (message: string) => void;

  setToken(token: string) {
    this.authToken = token;
  }

  protected onSuccessEditorDuplicateQuestionWithDifficulty(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Groupe de question créé.");
    }
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

  protected onSuccessLaunchSessionTrivialPoursuit(
    options: LaunchSessionOut
  ): void {
    this.inRequest = false;
    if (this.showMessage) {
      if (options.GroupStrategyKind == GroupStrategyKind.RandomGroupStrategy) {
        this.showMessage(`Session de Trivial lancée : ${options.SessionID}`);
      } else {
        this.showMessage(`Session de Trivial lancée.`);
      }
    }
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

  protected onSuccessEditorSearchQuestions(data: ListQuestionsOut): void {
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

/** `IsDev` is true when the client app is served in dev mode */
export const IsDev = import.meta.env.DEV;

// when building for production, the mode for the preview client
// may actually still be "dev", for instance when served in test
export const PreviewMode = IsDev
  ? "dev"
  : window.location.origin == localhost
  ? "dev"
  : "prod";

const devLogMeta = {
  IdTeacher: 6,
  Token:
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJUZWFjaGVyIjp7ImlkIjo2LCJtYWlsIjoiMTY1MzMwNjk3MyIsInBhc3N3b3JkX2NyeXB0ZWQiOiIiLCJpc19hZG1pbiI6ZmFsc2V9LCJleHAiOjE2NTM1NjYxNzN9.tHLmCB9N9U-546O5-j1uGG7ODNWs-7dyDD-b2X3mHTI",
};

export const controller = new Controller(
  IsDev ? localhost : window.location.origin,
  IsDev ? devLogMeta.Token : ""
);
