/// BuildMode switches between 3 build modes
enum BuildMode {
  /// connect to a remote server
  production,

  /// connect to a localhost server
  dev,

  /// no API connection, use embedded events
  debug,
}

extension Api on BuildMode {
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
}

/// buildMode returns the build mode
BuildMode buildMode() {
  const buildMode = String.fromEnvironment("mode");
  switch (buildMode) {
    case "debug":
      return BuildMode.debug;
    case "dev":
      return BuildMode.dev;
    default:
      return BuildMode.production;
  }
}
