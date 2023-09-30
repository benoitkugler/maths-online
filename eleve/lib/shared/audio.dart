import 'package:audioplayers/audioplayers.dart';
import 'package:flutter/material.dart';

enum SongCategorie { guitare, disco, alternative }

extension SongLabel on SongCategorie {
  String label() {
    switch (this) {
      case SongCategorie.disco:
        return "Disco";
      case SongCategorie.guitare:
        return "Guitare";
      case SongCategorie.alternative:
        return "Alternative";
    }
  }
}

class Song {
  final String path;
  final SongCategorie categorie;

  String get title => path
      .splitMapJoin(RegExp('[A-Z]'), onMatch: (m) => " ${m.group(0)!}")
      .trimLeft();

  const Song(this.path, this.categorie);
}

class Audio {
  /// all songs
  static final defaultPlaylist =
      PlaylistController([0, 1, 2, 3, 4, 5, 6, 7, 8, 9], true);

  /// [availableSongs] is the list of soundtracks available
  /// in the app.
  static const availableSongs = <Song>[
    Song("GrooveBow.mp3", SongCategorie.disco),
    Song("Forgive.mp3", SongCategorie.guitare),
    Song("EntreCielEtTerre.mp3", SongCategorie.alternative),
    Song("NouvelleTrajectoire.mp3", SongCategorie.disco),
    Song("AlternativeConnect.mp3", SongCategorie.disco),
    Song("AuroreBoreale.mp3", SongCategorie.disco),
    Song("DropFlow.mp3", SongCategorie.disco),
    // Song("Envolées.mp3", SongCategorie.guitare),
    Song("FarUp.mp3", SongCategorie.guitare),
    Song("PremiersPas.mp3", SongCategorie.guitare),
    // Song("SetOnFire.mp3", SongCategorie.guitare),
    // Song("Solitude.mp3", SongCategorie.guitare),
    Song("Suspens.mp3", SongCategorie.guitare),
    Song("Tempos.mp3", SongCategorie.guitare),
  ];

  PlaylistController playlist = PlaylistController([], false);

  final AudioPlayer _player = AudioPlayer();
  List<int> _songs = [];
  int _currentSong = -1;

  Audio();

  /// Set the songs to play in loop, as indexes in [availableSongs].
  /// Note that is does not start playing music.
  void setSongs(PlaylistController playlist) {
    this.playlist = playlist;

    _currentSong = -1;
    _songs = playlist.songs.toList();
    if (playlist.random) {
      _songs.shuffle();
    }
  }

  /// launch the next song in the playlist
  Future<void> run() async {
    return _startNextSong();
  }

  /// stop the player
  void pause() async {
    await _player.stop();
  }

  void _onSongDone() async {
    await _player.stop();
    // await _player.dispose();
    await _startNextSong();
  }

  Future<void> _startNextSong() async {
    if (_songs.isEmpty) {
      return;
    }

    _currentSong++; // read the next

    final songName = availableSongs[_songs[_currentSong % _songs.length]].path;
    await _player
        .setSourceAsset("music/$songName"); // the audio cache prefix is assets/
    await _player.setReleaseMode(ReleaseMode.stop);
    await _player.resume();
    _player.onPlayerStateChanged.listen((event) {
      if (event == PlayerState.completed) {
        _onSongDone();
      }
    });
  }
}

class PlaylistController {
  List<int> songs;
  bool random;

  PlaylistController(this.songs, this.random);

  Map<String, dynamic> toJson() {
    return {
      "songs": songs,
      "random": random,
    };
  }

  static PlaylistController? fromJson(dynamic json) {
    final List<int> songs;
    final bool random;
    // handle deprecated files
    if (json is List) {
      songs = json.map((e) => e as int).toList();
      random = false;
      return PlaylistController(songs, random);
    } else if (json is Map) {
      songs = (json["songs"] as List<dynamic>).map((e) => e as int).toList();
      random = json["random"] as bool;
      return PlaylistController(songs, random);
    } else {
      return null;
    }
  }
}

/// [Playlist] let the user choose the music
/// he wants to play (in loop)
/// TODO: hear the music when choosing
class Playlist extends StatefulWidget {
  final PlaylistController controller;

  const Playlist(this.controller, {Key? key}) : super(key: key);

  @override
  _PlaylistState createState() => _PlaylistState();
}

class _PlaylistState extends State<Playlist> {
  @override
  Widget build(BuildContext context) {
    const l = Audio.availableSongs;
    return Scaffold(
        appBar: AppBar(
          title: const Text("Playlist"),
        ),
        body: Column(
          children: [
            SwitchListTile(
              title: const Text("Lecture aléatoire"),
              value: widget.controller.random,
              onChanged: (b) => setState(() {
                widget.controller.random = b;
              }),
            ),
            Expanded(
              child: ListView(
                children: List<CheckboxListTile>.generate(
                    l.length,
                    (index) => CheckboxListTile(
                          title: Text(l[index].title),
                          subtitle: Text(l[index].categorie.label()),
                          selected: widget.controller.songs.contains(index),
                          value: widget.controller.songs.contains(index),
                          onChanged: (_) => setState(() {
                            widget.controller.songs.contains(index)
                                ? widget.controller.songs.remove(index)
                                : widget.controller.songs.add(index);
                          }),
                        )),
              ),
            ),
          ],
        ));
  }
}
