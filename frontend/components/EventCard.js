"use client";

export default function EventCard({ ev, goingEvent }) {

  const d = new Date(ev.time);

  //let now = new Date()
  d.setHours(d.getHours() - 1);
  const isExpired = new Date(d) - new Date() < 0;



  return (
    <div className="events-list" key={ev.id}>
      <div className="event-card">
        <div className="event-header">
          <h3 className="event-title">{ev.title}</h3>
          <span className="event-datetime">{d.toLocaleString()}</span>
        </div>

        <p className="event-description">{ev.description}</p>

        <div className="event-actions">
          <button
            className={`btn-going ${ev.userAction === "going" ? "buttonClicked" : ""}`}
            onClick={!isExpired ? () => goingEvent("going", ev.id) : undefined}
            disabled={isExpired}
          >
            Going
          </button>
          <button
            className={`btn-not-going ${ev.userAction === "notGoing" ? "buttonClicked" : ""}`}
            onClick={!isExpired ? () => goingEvent("notGoing", ev.id) : undefined}
            disabled={isExpired}
          >
            Not Going
          </button>
        </div>

        <div className="event-stats">
          <span className="stat-going">{ev.going} going</span>
          <span className="stat-not-going">{ev.notGoing} not going</span>
        </div>
      </div>
    </div>
  );
}
