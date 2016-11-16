package binding

import (
	"k8s.io/client-go/pkg/watch"
)

func newFakeWatchBindingFunc(retErr error) (WatchBindingFunc, *watch.FakeWatcher) {
	fake := watch.NewFake()
	fn := func(namespace string) (watch.Interface, error) {
		if retErr != nil {
			return nil, retErr
		}
		return fake, nil
	}
	return fn, fake
}
