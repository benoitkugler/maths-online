/// BuildMode switches between 3 build modes
enum BuildMode {
  /// connect to a remote server
  production,

  /// connect to a localhost server
  dev,

  /// no API connection, use embedded events
  debug,
}

String _withQuery(String baseUrl, Map<String, dynamic> query) {
  return Uri.parse(baseUrl).replace(queryParameters: query).toString();
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
  String websocketURL(String endpoint,
      {Map<String, dynamic> query = const {}}) {
    switch (this) {
      case BuildMode.production:
        return _withQuery("wss://isyro.fr$endpoint", query);
      case BuildMode.dev:
        return _withQuery("ws://localhost:1323$endpoint", query);
      case BuildMode.debug:
        return "";
    }
  }

  /// serverURL returns url ending by the [endpoint],
  /// or an empty string
  /// [endpoint] is expected to start with a slash
  String serverURL(String endpoint, {Map<String, dynamic> query = const {}}) {
    switch (this) {
      case BuildMode.production:
        return _withQuery("https://isyro.fr$endpoint", query);
      case BuildMode.dev:
        return _withQuery("http://localhost:1323$endpoint", query);
      case BuildMode.debug:
        return "";
    }
  }
}

/// buildMode returns the build mode
BuildMode buildModeFromEnv() {
  const buildMode = String.fromEnvironment("mode");
  return APISetting.fromString(buildMode);
}
