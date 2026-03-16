package google

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/api/calendar/v3"
)

type EventingGCalendar struct {
	Service   *calendar.EventsService
	Calendars map[string]string
}

func (cal *EventingGCalendar) CreateEvent(
	calendarName string,
	date string, // TODO: Change to time.Time
	startTime string, // TODO: Change to time.Time
	summary string,
	location string,
) string {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}

	start := fmt.Sprintf("%s %s:00", date, startTime)
	startDateTime, err := time.ParseInLocation(
		"2006-01-02 15:04:05", start, loc)
	if err != nil {
		log.Fatal(err)
	}
	endDateTime := startDateTime.Add(3 * time.Hour)

	event, err := cal.Service.Insert(cal.Calendars[calendarName],
		&calendar.Event{
			Summary:  summary,
			Location: location,
			Start: &calendar.EventDateTime{
				DateTime: startDateTime.Format(time.RFC3339),
			},
			End: &calendar.EventDateTime{
				DateTime: endDateTime.Format(time.RFC3339),
			},
		},
	).Do()
	if err != nil {
		log.Fatalf("Unable to insert event: %v", err)
	}

	log.Printf("Inserted Event: [Id: %s, Start: %s, End: %s, Summary: %s, Location: %s]\n",
		event.Id, event.Start.DateTime, event.End.DateTime, event.Summary, event.Location)

	return event.Id
}

func (cal *EventingGCalendar) CreateAllDayEvent(
	calendarName string,
	date string,
	summary string,
	location string,
) string {
	event, err := cal.Service.Insert(cal.Calendars[calendarName],
		&calendar.Event{
			Summary:      summary,
			Location:     location,
			Start:        &calendar.EventDateTime{Date: date},
			End:          &calendar.EventDateTime{Date: date},
			Transparency: "transparent",
		},
	).Do()
	if err != nil {
		log.Fatalf("Unable to insert all day event: %v", err)
	}

	log.Printf("Inserted All Day Event: [Id: %s, Date: %s, Summary: %s]\n",
		event.Id, event.Start.Date, event.Summary)

	return event.Id
}
