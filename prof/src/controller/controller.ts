import { devLogMeta } from "@/env";
import type {
  AskInscriptionOut,
  CheckExerciceParametersOut,
  CheckMissingQuestionsOut,
  CheckQuestionParametersOut,
  Classroom,
  ClassroomExt,
  CopyTravailToOut,
  CreateTravailOut,
  DeleteExerciceOut,
  DeleteQuestionOut,
  Exercice,
  ExerciceExt,
  ExercicegroupExt,
  ExerciceHeader,
  ExerciceWithPreview,
  ExportExerciceLatexOut,
  ExportQuestionLatexOut,
  GenerateClassroomCodeOut,
  HomeworkMarksOut,
  Homeworks,
  Index,
  LaunchSessionOut,
  ListExercicesOut,
  ListQuestionsOut,
  LoadTargetOut,
  LogginOut,
  MonitorOut,
  Monoquestion,
  Question,
  QuestiongroupExt,
  RandomMonoquestion,
  Review,
  ReviewExt,
  ReviewHeader,
  RunningSessionMetaOut,
  SaveExerciceAndPreviewOut,
  SaveQuestionAndPreviewOut,
  SheetExt,
  Student,
  TaskExt,
  TeacherSettings,
  TextBlock,
  Travail,
  TrivialExt,
  TrivialSelfaccess
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

  public settings: TeacherSettings = {
    Mail: "",
    HasEditorSimplified: false,
    Password: ""
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

  protected onSuccessTeacherGetSettings(data: TeacherSettings): void {
    this.inRequest = false;
  }

  protected onSuccessTeacherUpdateSettings(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Paramètres modifiés avec succès.");
    }
  }

  protected onSuccessTeacherResetPassword(): void {
    this.inRequest = false;
  }

  protected onSuccessEditorGetQuestionsIndex(data: Index): void {
    this.inRequest = false;
  }
  protected onSuccessEditorGetExercicesIndex(data: Index): void {
    this.inRequest = false;
  }

  protected onSuccessEditorDuplicateQuestiongroup(): void {
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

  protected onSuccessEditorGetQuestions(data: Question[] | null): void {
    this.inRequest = false;
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

  protected onSuccessHomeworkGetMonoquestion(data: Monoquestion): void {
    this.inRequest = false;
  }
  protected onSuccessHomeworkGetRandomMonoquestion(
    data: RandomMonoquestion
  ): void {
    this.inRequest = false;
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

  protected onSuccessHomeworkGetSheets(data: Homeworks): void {
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

  protected onSuccessEditorDuplicateExercicegroup(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Exercice (et variantes) dupliqué avec succès.");
    }
  }

  protected onSuccessEditorQuestionExportLateX(
    data: ExportQuestionLatexOut
  ): void {
    this.inRequest = false;
  }

  protected onSuccessEditorExerciceExportLateX(
    data: ExportExerciceLatexOut
  ): void {
    this.inRequest = false;
  }

  protected onSuccessGetTrivialRunningSessions(
    data: RunningSessionMetaOut
  ): void {
    this.inRequest = false;
  }

  protected onSuccessStartTrivialGame(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Partie lancée avec succés.");
    }
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

  protected onSuccessTeacherGetClassroomStudents(data: Student[] | null): void {
    this.inRequest = false;
  }

  protected onSuccessTeacherGetClassrooms(data: ClassroomExt[] | null): void {
    this.inRequest = false;
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

  protected onSuccessAskInscription(data: AskInscriptionOut): void {
    this.inRequest = false;
  }

  protected onSuccessValidateInscription(): void {
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

  protected onSuccessEditorDuplicateQuestion(): void {
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
  protected onSuccessDeleteTrivialPoursuit(): void {
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

  protected onSuccessTrivialTeacherMonitor(data: MonitorOut): void {
    this.inRequest = false;
  }

  protected onSuccessTrivialGetSelfaccess(data: TrivialSelfaccess): void {
    this.inRequest = false;
  }
  protected onSuccessTrivialUpdateSelfaccess(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage(`Accès aux parties modifié avec succès.`);
    }
  }

  protected onSuccessEditorSearchQuestions(data: ListQuestionsOut): void {
    this.inRequest = false;
  }
  protected onSuccessEditorCreateQuestion(data: Question): void {
    this.inRequest = false;
  }
  protected onSuccessEditorUpdateTags(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage(`Etiquettes modifiées avec succès.`);
    }
  }
  protected onSuccessEditorGetTags(): void {
    this.inRequest = false;
  }

  protected onSuccessEditorDeleteQuestion(data: DeleteQuestionOut): void {
    this.inRequest = false;
    console.log(data);

    if (this.showMessage && data.Deleted) {
      this.showMessage("Question supprimée avec succès.");
    }
  }
  protected onSuccessEditorPausePreview(): void {
    this.inRequest = false;
  }

  protected onSuccessEditorUpdateQuestionTags(): void {
    this.inRequest = false;
    if (this.showMessage) {
      this.showMessage("Etiquettes modifiées avec succès.");
    }
  }
  protected onSuccessEditorSearchExercices(data: ListExercicesOut): void {
    this.inRequest = false;
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
  protected onSuccessEditorGetExerciceContent(data: ExerciceExt): void {
    this.inRequest = false;
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
  protected onSuccessReviewsList(data: ReviewHeader[] | null): void {
    this.inRequest = false;
  }
  protected onSuccessReviewLoad(data: ReviewExt): void {
    this.inRequest = false;
  }

  protected onSuccessReviewLoadTarget(data: LoadTargetOut): void {
    this.inRequest = false;
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
