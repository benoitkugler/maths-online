package editor

import "github.com/benoitkugler/maths-online/maths/questions/client"

//go:generate ../../../../../structgen/structgen -source=loopback_types.go -mode=dart:../../../../eleve/lib/loopback_types.gen.dart

type LoopbackState struct {
	Question client.Question `dart-extern:"exercices/types.gen.dart"`
	IsPaused bool
}

// to keep in sync with eleve project
type loopbackClientDataKind uint8

type loopbackServerDataKind uint8
