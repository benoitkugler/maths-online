/// BuildMode switches between 3 build modes
enum BuildMode {
  /// connect to a remote server
  production,

  /// connect to a localhost server
  dev,

  /// no API connection, use embedded events
  debug,
}

extension APISetting on BuildMode {
  static BuildMode fromString(String bm) {
    switch (bm) {
      case "debug":
        return BuildMode.debug;
      case "dev":
        return BuildMode.dev;
      default:
        return BuildMode.production;
    }
  }

  /// websocketURL returns url ending by the [endpoint],
  /// or an empty string
  /// [endpoint] is expected to start with a slash
  String websocketURL(String endpoint) {
    switch (this) {
      case BuildMode.production:
        return "wss://education.alwaysdata.net" + endpoint;
      case BuildMode.dev:
        return "ws://localhost:1323" + endpoint;
      case BuildMode.debug:
        return "";
    }
  }

  /// serverURL returns url ending by the [endpoint],
  /// or an empty string
  /// [endpoint] is expected to start with a slash
  String serverURL(String endpoint) {
    switch (this) {
      case BuildMode.production:
        return "https://education.alwaysdata.net" + endpoint;
      case BuildMode.dev:
        return "http://localhost:1323" + endpoint;
      case BuildMode.debug:
        return "";
    }
  }
}

/// buildMode returns the build mode
BuildMode buildMode() {
  const buildMode = String.fromEnvironment("mode");
  return APISetting.fromString(buildMode);
}
