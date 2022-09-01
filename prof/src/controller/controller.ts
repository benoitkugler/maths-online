import { devLogMeta } from "@/env";
import type {
  AskInscriptionOut,
  CheckExerciceParametersOut,
  CheckMissingQuestionsOut,
  CheckQuestionParametersOut,
  Classroom,
  ClassroomExt,
  ClassroomSheets,
  Exercice,
  ExerciceExt,
  ExerciceHeader,
  GenerateClassroomCodeOut,
  LaunchSessionOut,
  ListQuestionsOut,
  LogginOut,
  Question,
  QuestiongroupExt,
  RunningSessionMetaOut,
  SaveExerciceAndPreviewOut,
  SaveQuestionAndPreviewOut,
  SheetExt,
  StartSessionOut,
  Student,
  Task,
  TrivialExt,
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
  public showMessage?: (message: string, color?: string) => void;

  public editorSessionID = "";

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

  protected onSuccessEditorDuplicateQuestiongroup(data: never): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question (et variantes) dupliquée avec succès.");
    }
  }
  protected onSuccessEditorCreateQuestiongroup(data: QuestiongroupExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question créée avec succès.");
    }
  }
  protected onSuccessEditorGetQuestions(data: Question[] | null): void {
    this.inRequest = false;
  }

  protected onSuccessHomeworkRemoveTask(data: never): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Exercice retiré avec succès.");
    }
  }
  protected onSuccessHomeworkAddTask(data: Task): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Exercice ajouté avec succès.");
    }
  }
  protected onSuccessHomeworkReorderSheetTasks(data: never): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Liste modifiée avec succès.");
    }
  }

  protected onSuccessHomeworkGetSheets(data: ClassroomSheets[] | null): void {
    this.inRequest = false;
  }
  protected onSuccessHomeworkCreateSheet(data: SheetExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Fiche ajoutée avec succès.");
    }
  }
  protected onSuccessHomeworkCopySheet(data: SheetExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Fiche dupliquée avec succès.");
    }
  }
  protected onSuccessHomeworkUpdateSheet(data: never): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Fiche modifiée avec succès.");
    }
  }
  protected onSuccessHomeworkDeleteSheet(data: never): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Fiche supprimée avec succès.");
    }
  }

  protected onSuccessExercicesGetList(data: ExerciceHeader[] | null): void {
    this.inRequest = false;
  }
  protected onSuccessExerciceGetContent(data: ExerciceExt): void {
    this.inRequest = false;
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
    if (this.showMessage && data.IsValid) {
      this.showMessage(`Question générée avec succès.`);
    }
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

  protected onSuccessExerciceCreateQuestion(data: ExerciceExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question ajoutée avec succès.");
    }
  }
  protected onSuccessExerciceUpdateQuestions(data: ExerciceExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Questions modifiées avec succès.");
    }
  }

  protected onSuccessExerciceCreate(data: ExerciceHeader): void {
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

  protected onSuccessUpdateTrivialVisiblity(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Visibilité modifiée avec succès.");
    }
  }

  protected onSuccessAskInscription(data: AskInscriptionOut): void {
    this.inRequest = false;
  }

  protected onSuccessQuestionUpdateVisiblity(): void {
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

  protected onSuccessDuplicateTrivialPoursuit(data: TrivialExt): void {
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

  protected onSuccessGetTrivialPoursuit(data: TrivialExt[] | null): void {
    this.inRequest = false;
  }
  protected onSuccessCreateTrivialPoursuit(data: TrivialExt): void {
    this.inRequest = false;
  }
  protected onSuccessUpdateTrivialPoursuit(data: TrivialExt): void {
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
    this.editorSessionID = data.ID;
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
    if (this.showMessage) {
      this.showMessage(`Etiquettes modifiées avec succès.`);
    }
  }
  protected onSuccessEditorGetTags(data: any): void {
    this.inRequest = false;
  }

  protected onSuccessEditorDeleteQuestion(data: any): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question supprimée avec succès.");
    }
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
