import { devLogMeta } from "@/env";
import type {
  AskInscriptionOut,
  CheckExerciceParametersOut,
  CheckMissingQuestionsOut,
  CheckQuestionParametersOut,
  Classroom,
  ClassroomExt,
  Exercice,
  ExerciceExt,
  ExerciceQuestionExt,
  GenerateClassroomCodeOut,
  LaunchSessionOut,
  ListQuestionsOut,
  LogginOut,
  Question,
  RunningSessionMetaOut,
  SaveExerciceAndPreviewOut,
  SaveQuestionAndPreviewOut,
  StartSessionOut,
  Student,
  TrivialConfigExt,
  UpdateGroupTagsOut,
} from "./api_gen";
import { AbstractAPI } from "./api_gen";

function arrayBufferToString(buffer: ArrayBuffer) {
  const uintArray = new Uint8Array(buffer);
  const encodedString = String.fromCharCode.apply(null, Array.from(uintArray));
  return decodeURIComponent(escape(encodedString));
}

class Controller extends AbstractAPI {
  private isLoggedIn = false;

  public onError?: (kind: string, htmlError: string) => void;
  public showMessage?: (message: string) => void;

  logout() {
    this.isLoggedIn = false;
    this.authToken = "";
  }

  getToken() {
    return this.authToken;
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

  protected onSuccessEditorCheckQuestionParameters(
    data: CheckQuestionParametersOut
  ): void {
    this.inRequest = false;
  }
  protected onSuccessEditorSaveQuestionAndPreview(
    data: SaveQuestionAndPreviewOut
  ): void {
    this.inRequest = false;
  }

  protected onSuccessExerciceUpdateVisiblity(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Visibilité modifiée avec succès.");
    }
  }
  protected onSuccessEditorCheckExerciceParameters(
    data: CheckExerciceParametersOut
  ): void {
    this.inRequest = false;
  }
  protected onSuccessEditorSaveExerciceAndPreview(
    data: SaveExerciceAndPreviewOut
  ): void {
    this.inRequest = false;
  }

  protected onSuccessExerciceCreateQuestion(
    data: ExerciceQuestionExt[] | null
  ): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question ajoutée avec succès.");
    }
  }
  protected onSuccessExerciceUpdateQuestions(
    data: ExerciceQuestionExt[] | null
  ): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Questions modifiées avec succès.");
    }
  }

  protected onSuccessExercicesGetList(data: ExerciceExt[] | null): void {
    this.inRequest = false;
  }
  protected onSuccessExerciceCreate(data: ExerciceExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Exercice créé avec succès.");
    }
  }
  protected onSuccessExerciceDelete(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Exercice supprimé avec succès.");
    }
  }
  protected onSuccessExerciceUpdate(data: Exercice): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Exercice mis à jour avec succès.");
    }
  }

  protected onSuccessGetTrivialRunningSessions(
    data: RunningSessionMetaOut
  ): void {
    this.inRequest = false;
  }

  protected onSuccessTeacherGenerateClassroomCode(
    data: GenerateClassroomCodeOut
  ): void {
    this.inRequest = false;
  }

  protected onSuccessTeacherAddStudent(data: Student): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Elève ajouté avec succès.");
    }
  }
  protected onSuccessTeacherUpdateStudent(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Profil mis à jour avec succès.");
    }
  }
  protected onSuccessTeacherDeleteStudent(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Profil supprimé avec succès.");
    }
  }

  protected onSuccessTeacherImportStudents(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Liste importée avec succès.");
    }
  }

  protected onSuccessTeacherGetClassroomStudents(data: Student[] | null): void {
    this.inRequest = false;
  }

  protected onSuccessTeacherGetClassrooms(data: ClassroomExt[] | null): void {
    this.inRequest = false;
  }
  protected onSuccessTeacherCreateClassroom(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Classe créée avec succès.");
    }
  }
  protected onSuccessTeacherUpdateClassroom(data: Classroom): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Classe mise à jour avec succès.");
    }
  }
  protected onSuccessTeacherDeleteClassroom(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Classe supprimée avec succès.");
    }
  }

  protected onSuccessStopTrivialGame(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Partie interrompue avec succès");
    }
  }

  protected onSuccessEditorUpdateGroupTags(data: UpdateGroupTagsOut): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Etiquettes modifiées avec succès.");
    }
  }

  protected onSuccessUpdateTrivialVisiblity(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Visibilité modifiée avec succès.");
    }
  }

  protected onSuccessAskInscription(data: AskInscriptionOut): void {
    this.inRequest = false;
  }

  protected onSuccessQuestionUpdateVisiblity(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Visibilité modifiée avec succès.");
    }
  }

  protected onSuccessValidateInscription(data: any): void {
    this.inRequest = false;
  }
  protected onSuccessLoggin(data: LogginOut): void {
    this.inRequest = false;
    if (data.Error == "") {
      this.authToken = data.Token;
      this.isLoggedIn = true;
      if (this.showMessage) {
        this.showMessage("Bienvenue");
      }
    }
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

  protected onSuccessEditorDuplicateQuestion(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question dupliquée.");
    }
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
    if (this.showMessage) {
      this.showMessage("Configuration mise à jour.");
    }
  }
  protected onSuccessDeleteTrivialPoursuit(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Configuration supprimée.");
    }
  }

  protected onSuccessLaunchSessionTrivialPoursuit(
    options: LaunchSessionOut
  ): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage(`Parties lancées avec succès.`);
    }
  }

  protected onSuccessEditorStartSession(data: StartSessionOut): void {
    console.log(data);
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

export function isInscriptionValidated() {
  return window.location.search.includes("show-success-inscription");
}

export const controller = new Controller(
  IsDev ? localhost : window.location.origin,
  IsDev ? devLogMeta.Token : ""
);
