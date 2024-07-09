package sys

import (
	"reflect"
	"testing"

	"github.com/fsnotify/fsnotify"
)

func TestFSNotifyHandlerMiddlewareAnyOf(t *testing.T) {
	tests := []struct {
		name      string
		event     fsnotify.Event
		allow     fsnotify.Op
		allowRest []fsnotify.Op
		want      map[fsnotify.Op]struct{}
	}{
		{
			name:  "only writes",
			event: fsnotify.Event{Op: fsnotify.Write},
			allow: fsnotify.Write,
			want:  map[fsnotify.Op]struct{}{fsnotify.Write: {}},
		},
		{
			name:  "discard chmod because it is not write",
			event: fsnotify.Event{Op: fsnotify.Chmod},
			allow: fsnotify.Write,
			want:  map[fsnotify.Op]struct{}{},
		},
		{
			name:      "allow chmod since it is on the list",
			event:     fsnotify.Event{Op: fsnotify.Chmod},
			allow:     fsnotify.Write,
			allowRest: []fsnotify.Op{fsnotify.Chmod},
			want:      map[fsnotify.Op]struct{}{fsnotify.Chmod: {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make(map[fsnotify.Op]struct{})
			handler := FSNotifyHandlerFunc(func(event fsnotify.Event) {
				got[event.Op] = struct{}{}
			})

			AnyOf(handler, tt.allow, tt.allowRest...).Notify(tt.event)

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want accept=%+v, got accept=%+v", tt.want, got)
			}
		})
	}
}
