package data

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

const apiTemplate = "https://api.meetup.com/%s/events?status=upcoming"

var (
	meetupNames = []string{
		"Boulder-Gophers",
		"Denver-Go-Language-User-Group",
		"Denver-Go-Programming-Language-Meetup",
	}
)

// Store contains data for the site.
type Store struct {
	pollingInterval time.Duration

	mu         sync.Mutex
	eventCache map[string][]Event
}

// NewStore creates a new store initialized with a polling interval.
func NewStore(i time.Duration) *Store {
	return &Store{
		pollingInterval: i,
	}
}

// Poll runs forever, polling the meetup API for event data and updating the
// internal cache.
func (s *Store) Poll() {
	for {
		events := s.poll()
		s.updateCache(events)
		time.Sleep(s.pollingInterval)
	}
}

func (s *Store) updateCache(events map[string][]Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.eventCache = events
}

func (s *Store) poll() map[string][]Event {
	all := make(map[string][]Event)
	for _, meetup := range meetupNames {
		eds, err := events(meetup)
		if err != nil {
			log.Printf("error fetching events for %s: %s", meetup, err)
			continue
		}
		all[meetup] = eds
	}

	for _, v := range all {
		sort.Slice(v, func(i, j int) bool {
			return v[i].Time < v[j].Time
		})
	}

	return all
}

// AllEvents returns the current meetup events in CO.
func (s *Store) AllEvents() map[string][]Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.eventCache
}

// Event contains information about a meetup event.
type Event struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Time int64  `json:"time"`
}

// HumanTime returns the time formated for the UI.
func (e Event) HumanTime() string {
	return time.Unix(e.Time/1000, 0).Format(time.RFC1123)
}

func events(name string) ([]Event, error) {
	resp, err := http.Get(fmt.Sprintf(apiTemplate, name))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var data []Event
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
