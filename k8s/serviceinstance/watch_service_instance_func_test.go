package serviceinstance

import (
	"k8s.io/client-go/pkg/watch"
)

func newFakeWatchServiceInstanceFunc(retErr error) (WatchServiceInstanceFunc, *watch.FakeWatcher) {
	fake := watch.NewFake()
	fn := func(string) (watch.Interface, error) {
		if retErr != nil {
			return nil, retErr
		}
		return fake, nil
	}
	return fn, fake
}
