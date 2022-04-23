import type {
  CheckParametersOut,
  Question,
  QuestionHeader,
  StartSessionOut,
  TrivialConfig,
  TrivialConfigExt
} from "./api_gen";
import { AbstractAPI } from "./api_gen";

class Controller extends AbstractAPI {
  protected onSuccessLaunchSessionTrivialPoursuit(data: TrivialConfig): void {
    throw new Error("Method not implemented.");
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
  requestError = "";

  getURL(endpoint: string) {
    return this.baseUrl + endpoint;
  }

  handleError(error: any): void {
    this.inRequest = false;
    this.requestError = `${error}`;
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

export const BuildMode = import.meta.env.DEV ? "dev" : "prod";

export const controller = new Controller(
  BuildMode == "dev" ? "http://localhost:1323" : window.location.origin,
  "",
  {}
);
