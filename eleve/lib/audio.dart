import 'package:audioplayers/audioplayers.dart';

class Audio {
  AudioPlayer _player = AudioPlayer();
  final AudioCache _cache = AudioCache(prefix: "lib/music/");
  List<String> _songs = [];

  Audio();

  int _currentSong = -1;

  /// set the songs to play in loop
  void setSongs(List<String> songs) {
    _currentSong = -1;
    _songs = songs;
  }

  /// launch the next song in the playlist
  Future<void> run() async {
    return _startNextSong();
  }

  /// stop the player and skip to the next song
  void pause() async {
    await _player.stop();
  }

  void _onSongDone() async {
    await _player.stop();
    await _player.dispose();
    await _startNextSong();
  }

  Future<void> _startNextSong() async {
    if (_songs.isEmpty) {
      return;
    }

    _currentSong++;

    _player = await _cache.play(_songs[_currentSong % _songs.length]);
    await _player.setReleaseMode(ReleaseMode.STOP);
    _player.onPlayerStateChanged.listen((event) {
      if (event == PlayerState.COMPLETED) {
        _onSongDone();
      }
    });
  }
}
