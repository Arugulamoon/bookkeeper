package google

import (
	"log"
	"strings"

	"google.golang.org/api/calendar/v3"
)

type AccountingGCalendar struct {
	Service   *calendar.EventsService
	Calendars map[string]string
}

func (cal *AccountingGCalendar) FindEvent(
	calendarName string,
	id string,
) *calendar.Event {
	event, err := cal.Service.Get(cal.Calendars[calendarName], id).Do()
	if err != nil {
		log.Fatalf("Unable to get event: %v", err)
	}

	log.Printf("Found Event: [Id: %s, Date: %s, Summary: %s]\n",
		event.Id, event.Start.Date, event.Summary)

	return event
}

func (cal *AccountingGCalendar) CreateAllDayEvent(
	calendarName string,
	date, summary string,
) string {
	event, err := cal.Service.Insert(cal.Calendars[calendarName],
		&calendar.Event{
			Summary:      summary,
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

func (cal *AccountingGCalendar) MarkEventPaid(
	calendarName string,
	id string,
) {
	found := cal.FindEvent(calendarName, id)
	found.Summary = strings.Replace(found.Summary, "DUE", "PAID", 1)

	event, err := cal.Service.Update(cal.Calendars[calendarName], id, found).Do()
	if err != nil {
		log.Fatalf("Unable to update event: %v", err)
	}

	log.Printf("Marked Event Paid: [Id: %s, Date: %s, Summary: %s]\n",
		event.Id, event.Start.Date, event.Summary)
}

func (cal *AccountingGCalendar) DeleteEvent(
	calendarName string,
	id string,
) {
	err := cal.Service.Delete(cal.Calendars[calendarName], id).Do()
	if err != nil {
		log.Fatalf("Unable to delete event: %v", err)
	}
}
