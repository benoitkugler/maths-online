import { devLogMeta } from "@/env";
import type {
  Classroom,
  CopyTravailToOut,
  CreateTravailOut,
  DeleteExerciceOut,
  DeleteQuestionOut,
  Exercice,
  ExercicegroupExt,
  ExerciceHeader,
  ExerciceWithPreview,
  HomeworkMarksOut,
  LaunchSessionOut,
  LogginOut,
  QuestiongroupExt,
  Review,
  SaveQuestionAndPreviewOut,
  SheetExt,
  StudentExt,
  TaskExt,
  TeacherSettings,
  TextBlock,
  Travail,
  TrivialExt,
} from "./api_gen";
import { AbstractAPI, MatiereTag } from "./api_gen";

function arrayBufferToString(buffer: ArrayBuffer) {
  const uintArray = new Uint8Array(buffer);
  const encodedString = String.fromCharCode.apply(null, Array.from(uintArray));
  return decodeURIComponent(escape(encodedString));
}

class Controller extends AbstractAPI {
  private isLoggedIn = false;

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
    return this.baseUrl + endpoint;
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
    console.log("launching request");
    this.inRequest = true;
  }

  async ensureSettings() {
    const res = await this.TeacherGetSettings();
    if (res == undefined) return;
    this.settings = res;
  }

  protected onSuccessTeacherUpdateSettings(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Paramètres modifiés avec succès.");
    }
  }

  protected onSuccessEditorDuplicateQuestiongroup(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question (et variantes) dupliquée avec succès.");
    }
  }

  protected onSuccessEditorSaveQuestionMeta(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question enregistrée avec succès.");
    }
  }

  protected onSuccessEditorUpdateQuestiongroup(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question modifiée avec succès.");
    }
  }

  protected onSuccessEditorUpdateQuestiongroupVis(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Visibilité modifiée avec succès.");
    }
  }

  protected onSuccessEditorGenerateSyntaxHint(data: TextBlock): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Conseils ajoutés avec succès.");
    }
  }

  protected onSuccessHomeworkRemoveTask(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Exercice retiré avec succès.");
    }
  }

  protected onSuccessHomeworkAddExercice(data: TaskExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Exercice ajouté avec succès.");
    }
  }
  protected onSuccessHomeworkAddMonoquestion(data: TaskExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question ajoutée avec succès.");
    }
  }

  protected onSuccessHomeworkAddRandomMonoquestion(data: TaskExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Groupe de questions ajouté avec succès.");
    }
  }

  protected onSuccessHomeworkUpdateMonoquestion(data: TaskExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Paramètres de la question modifiés avec succès.");
    }
  }

  protected onSuccessHomeworkUpdateRandomMonoquestion(data: TaskExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage(
        "Paramètres du groupe de questions modifiés avec succès."
      );
    }
  }

  protected onSuccessHomeworkReorderSheetTasks(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Liste modifiée avec succès.");
    }
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
  protected onSuccessHomeworkUpdateSheet(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Fiche modifiée avec succès.");
    }
  }
  protected onSuccessHomeworkDeleteSheet(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Fiche supprimée avec succès.");
    }
  }

  protected onSuccessHomeworkGetMarks(data: HomeworkMarksOut): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Notes chargées avec succès.");
    }
  }

  protected onSuccessHomeworkCreateTravailWith(data: Travail): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Feuille de travail ajoutée avec succès.");
    }
  }

  protected onSuccessHomeworkCreateTravail(data: CreateTravailOut): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Feuille de travail ajoutée avec succès.");
    }
  }

  protected onSuccessHomeworkUpdateTravail(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Feuille de travail modifiée avec succès.");
    }
  }
  protected onSuccessHomeworkDeleteTravail(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Feuille de travail supprimée avec succès.");
    }
  }
  protected onSuccessHomeworkCopyTravail(data: CopyTravailToOut): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Feuille de travail copiée avec succès.");
    }
  }
  protected onSuccessHomeworkSetDispense(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Permissions mises à jour avec succès.");
    }
  }

  protected onSuccessEditorSaveQuestionAndPreview(
    data: SaveQuestionAndPreviewOut
  ): void {
    this.inRequest = false;
    if (this.showMessage && data.IsValid) {
      this.showMessage(`Question générée avec succès.`);
    }
  }

  protected onSuccessEditorDuplicateExercicegroup(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Exercice (et variantes) dupliqué avec succès.");
    }
  }

  protected onSuccessStartTrivialGame(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Partie lancée avec succés.");
    }
  }

  protected onSuccessTeacherAddStudent(data: StudentExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Elève ajouté avec succès.");
    }
  }
  protected onSuccessTeacherUpdateStudent(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Profil mis à jour avec succès.");
    }
  }
  protected onSuccessTeacherDeleteStudent(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Profil supprimé avec succès.");
    }
  }

  protected onSuccessTeacherImportStudents(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Liste importée avec succès.");
    }
  }

  protected onSuccessTeacherCreateClassroom(): void {
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
  protected onSuccessTeacherDeleteClassroom(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Classe supprimée avec succès.");
    }
  }

  protected onSuccessStopTrivialGame(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Partie interrompue avec succès");
    }
  }

  protected onSuccessUpdateTrivialVisiblity(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Visibilité modifiée avec succès.");
    }
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

  protected onSuccessEditorDuplicateQuestion(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question dupliquée.");
    }
  }

  protected onSuccessCreateTrivialPoursuit(data: TrivialExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Partie d'Isy'Triv créée avec succés..");
    }
  }

  protected onSuccessDuplicateTrivialPoursuit(data: TrivialExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Session dupliquée.");
    }
  }

  protected onSuccessUpdateTrivialPoursuit(data: TrivialExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Partie mise à jour.");
    }
  }
  protected onSuccessDeleteTrivialPoursuit(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Partie supprimée.");
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

  protected onSuccessTrivialUpdateSelfaccess(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage(`Accès aux parties modifié avec succès.`);
    }
  }

  protected onSuccessEditorUpdateTags(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage(`Etiquettes modifiées avec succès.`);
    }
  }

  protected onSuccessEditorDeleteQuestion(data: DeleteQuestionOut): void {
    this.inRequest = false;
    if (this.showMessage && data.Deleted) {
      this.showMessage("Question supprimée avec succès.");
    }
  }

  protected onSuccessEditorUpdateQuestionTags(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Etiquettes modifiées avec succès.");
    }
  }

  protected onSuccessEditorUpdateExercicegroup(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Exercice modifié avec succès.");
    }
  }
  protected onSuccessEditorUpdateExerciceTags(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Etiquettes modifiées avec succès.");
    }
  }

  protected onSuccessEditorCreateExercice(data: ExercicegroupExt): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Exercice crée avec succès.");
    }
  }
  protected onSuccessEditorDeleteExercice(data: DeleteExerciceOut): void {
    this.inRequest = false;
    if (this.showMessage && data.Deleted) {
      this.showMessage("Exercice supprimé avec succès.");
    }
  }
  protected onSuccessEditorUpdateExercice(data: Exercice): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Exercice modifié avec succès.");
    }
  }
  protected onSuccessEditorDuplicateExercice(data: ExerciceHeader): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Variante dupliquée avec succès.");
    }
  }
  protected onSuccessEditorExerciceCreateQuestion(
    data: ExerciceWithPreview
  ): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question ajoutée avec succès.");
    }
  }
  protected onSuccessEditorExerciceUpdateQuestions(
    data: ExerciceWithPreview
  ): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Questions modifiées avec succès.");
    }
  }
  protected onSuccessEditorExerciceDuplicateQuestion(
    data: ExerciceWithPreview
  ): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question dupliquée avec succès.");
    }
  }
  protected onSuccessEditorUpdateExercicegroupVis(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Visibilité modifiée avec succès.");
    }
  }

  protected onSuccessEditorSaveExerciceMeta(data: Exercice): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Exercice modifié avec succès.");
    }
  }
  protected onSuccessEditorExerciceImportQuestion(
    data: ExerciceWithPreview
  ): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Question importée avec succès.");
    }
  }

  // reviews
  protected onSuccessReviewCreate(data: Review): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Demande de publication créée avec succès.");
    }
  }

  protected onSuccessReviewDelete(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Demande de publication supprimée avec succès.");
    }
  }
  protected onSuccessReviewUpdateApproval(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Vote modifié avec succès.");
    }
  }

  protected onSuccessReviewUpdateCommnents(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Commentaires modifiés avec succès.");
    }
  }
  protected onSuccessReviewAccept(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Demande de publication acceptée avec succès.");
    }
  }

  protected onSuccessEditorDeleteQuestiongroup(data: DeleteQuestionOut): void {
    if (this.showMessage && data.Deleted)
      this.showMessage("Question (et ses variantes) supprimée avec succès.");
  }
  protected onSuccessEditorDeleteExercicegroup(data: DeleteExerciceOut): void {
    if (this.showMessage && data.Deleted)
      this.showMessage("Exercice (et ses variantes) supprimé avec succès.");
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
