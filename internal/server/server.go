package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/fsnotify/fsnotify"
	"github.com/zhhc99/bgen/internal/build"
)

const defaultAddr = ":8080"

type broker struct {
	mu      sync.Mutex
	clients map[chan struct{}]struct{}
}

func (b *broker) subscribe() chan struct{} {
	ch := make(chan struct{}, 1)
	b.mu.Lock()
	b.clients[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *broker) unsubscribe(ch chan struct{}) {
	b.mu.Lock()
	delete(b.clients, ch)
	b.mu.Unlock()
}

func (b *broker) broadcast() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.clients {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func Run(projectRoot string) error {
	if err := build.RunDev(projectRoot); err != nil {
		return err
	}

	b := &broker{clients: make(map[chan struct{}]struct{})}

	if err := startWatcher(projectRoot, b); err != nil {
		return err
	}

	outPath := filepath.Join(projectRoot, "output")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(outPath)))
	mux.HandleFunc("/__reload", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			return
		}
		defer conn.Close(websocket.StatusNormalClosure, "")

		ch := b.subscribe()
		defer b.unsubscribe(ch)

		ctx := r.Context()
		for {
			select {
			case <-ch:
				conn.Write(ctx, websocket.MessageText, []byte("reload"))
			case <-ctx.Done():
				return
			}
		}
	})

	fmt.Printf("bgen: serving at http://localhost%s\n", defaultAddr)
	return http.ListenAndServe(defaultAddr, mux)
}

// addDirRecursive 将 dir 下所有子目录（包括自身）加入 watcher.
// 目录不存在时静默跳过.
func addDirRecursive(w *fsnotify.Watcher, dir string) {
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}
		_ = w.Add(path)
		return nil
	})
}

func startWatcher(projectRoot string, b *broker) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating watcher: %w", err)
	}

	_ = w.Add(filepath.Join(projectRoot, "blog.yaml"))

	for _, dir := range []string{"content", "layouts", "static"} {
		addDirRecursive(w, filepath.Join(projectRoot, dir))
	}

	go func() {
		defer w.Close()
		var timer *time.Timer
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Chmod) {
					continue
				}
				// 新建目录时顺手加入监听
				if event.Has(fsnotify.Create) {
					if info, err := filepath.EvalSymlinks(event.Name); err == nil {
						_ = w.Add(info)
					}
				}
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(500*time.Millisecond, func() {
					log.Println("bgen: change detected, rebuilding...")
					if err := build.RunDev(projectRoot); err != nil {
						log.Printf("bgen: build error: %v\n", err)
						return
					}
					log.Println("bgen: rebuild complete")
					b.broadcast()
				})
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				log.Printf("bgen: watcher error: %v\n", err)
			}
		}
	}()

	return nil
}
