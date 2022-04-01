import 'package:audioplayers/audioplayers.dart';
import 'package:flutter/material.dart';

enum SongCategorie { guitare, disco }

extension SongLabel on SongCategorie {
  String label() {
    switch (this) {
      case SongCategorie.disco:
        return "Disco";
      case SongCategorie.guitare:
        return "Guitare";
    }
  }
}

class Song {
  final String path;
  final SongCategorie categorie;

  String get title => path
      .splitMapJoin(RegExp('[A-Z]'), onMatch: (m) => " " + m.group(0)!)
      .trimLeft();

  const Song(this.path, this.categorie);
}

class Audio {
  /// [availableSongs] is the list of soundtracks available
  /// in the app.
  static const availableSongs = <Song>[
    Song("GrooveBow.mp3", SongCategorie.disco),
    Song("NouvelleTrajectoire.mp3", SongCategorie.disco),
    Song("AlternativeConnect.mp3", SongCategorie.disco),
    Song("AuroreBoreale.mp3", SongCategorie.disco),
    Song("DropFlow.mp3", SongCategorie.disco),
    Song("Envolées.mp3", SongCategorie.guitare),
    Song("FarUp.mp3", SongCategorie.guitare),
    Song("Forgive.mp3", SongCategorie.guitare),
    Song("PremiersPas.mp3", SongCategorie.guitare),
    Song("SetOnFire.mp3", SongCategorie.guitare),
    Song("Solitude.mp3", SongCategorie.guitare),
    Song("Suspens.mp3", SongCategorie.guitare),
    Song("Tempos.mp3", SongCategorie.guitare),
  ];

  AudioPlayer _player = AudioPlayer();
  final AudioCache _cache = AudioCache(prefix: "lib/music/");
  PlaylistController playlist = [];

  Audio();

  int _currentSong = -1;

  /// set the songs to play in loop, as indexes in [availableSongs]
  /// Note that is does not start playing music
  void setSongs(PlaylistController playlist) {
    _currentSong = -1;
    this.playlist = playlist;
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
    if (playlist.isEmpty) {
      return;
    }

    _currentSong++;

    _player = await _cache
        .play(availableSongs[playlist[_currentSong % playlist.length]].path);
    await _player.setReleaseMode(ReleaseMode.STOP);
    _player.onPlayerStateChanged.listen((event) {
      if (event == PlayerState.COMPLETED) {
        _onSongDone();
      }
    });
  }
}

typedef PlaylistController = List<int>;

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
    return ListView(
      children: List<CheckboxListTile>.generate(
          l.length,
          (index) => CheckboxListTile(
                title: Text(l[index].title),
                subtitle: Text(l[index].categorie.label()),
                selected: widget.controller.contains(index),
                value: widget.controller.contains(index),
                onChanged: (_) => setState(() {
                  widget.controller.contains(index)
                      ? widget.controller.remove(index)
                      : widget.controller.add(index);
                }),
              )),
    );
  }
}
