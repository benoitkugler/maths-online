import type {
  CheckParametersOut,
  LaunchGameOut,
  StartSessionOut
} from "./api_gen";
import { AbstractAPI } from "./api_gen";

class Controller extends AbstractAPI {
  getURL(endpoint: string) {
    return this.baseUrl + endpoint;
  }

  handleError(error: any): void {
    // TODO: real message
    console.log(`ERROR: ${error}`);
  }

  startRequest(): void {
    console.log("launching request");
  }

  protected onSuccessLaunchGame(data: LaunchGameOut): void {
    console.log(`Game started at ${data.URL}`);
  }

  protected onSuccessEditStartSession(data: StartSessionOut): void {
    console.log(data);
  }

  protected onSuccessEditSaveAndPreview(data: any): void {
    console.log("OK", data);
  }

  protected onSuccessEditCheckParameters(data: CheckParametersOut): void {
    console.log("OK", data);
  }
}

export const BuildMode = import.meta.env.DEV ? "dev" : "prod";

export const controller = new Controller(
  BuildMode == "dev" ? "http://localhost:1323" : window.location.origin,
  "",
  {}
);
