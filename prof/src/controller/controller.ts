import type { LaunchGameOut } from "./api_gen";
import { AbstractAPI } from "./api_gen";

class Controller extends AbstractAPI {
  handleError(error: any): void {
    console.log(`ERROR: ${error}`);
  }

  startRequest(): void {
    console.log("launching request");
  }

  protected onSuccessLaunchGame(data: LaunchGameOut): void {
    console.log(`Game started at ${data.URL}`);
  }
}

export const controller = new Controller(
  import.meta.env.DEV ? "http://localhost:1323" : window.location.origin,
  "",
  {}
);
