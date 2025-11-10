"use client";
import Navbar from "../../components/Navbar";
import LeftBar from "../../components/LeftBar.js";
import RightBar from "../../components/RightBar.js";
import { useDarkMode } from "../../context/darkMod.js";
import "../../styles/events-page.css";
import { useEffect, useState } from "react";
import Link from "next/link.js";
import formatTime from '../../helpers/formatTime.js';


export default function EventsPage() {
    const { darkMode } = useDarkMode();
    const [myEvents, setMyEvents] = useState([]);
    const [showSidebar, setShowSidebar] = useState(true);

    useEffect(() => {
        fetchEvents();
    }, []);
    function fetchEvents() {
        fetch('http://localhost:8080/api/myevents', {
            method: 'GET',
            credentials: 'include',
        })
            .then((response) => response.json())
            .then((data) => {

                setMyEvents(data);
            })
            .catch((error) => {
                console.error('Error fetching events:', error);
            });
    }

    return (
        <div id="div" className={darkMode ? "theme-dark" : "theme-light"}>
            <Navbar showSidebar={true} />
            <main className="content">
                <LeftBar showSidebar={showSidebar} />
                <div className="spacer" >

                    {
                        myEvents?.map((event) => (
                            <Link key={event.id} href={`/groups/${event.group_id}`} >
                                <div className="events-page" >
                                    <h1>{event.title}</h1>
                                    <p>{event.description}</p>
                                    <h3>{formatTime(event.time)}</h3>
                                    <h6>Created at: {formatTime(event.created_at)}</h6>
                                </div>
                            </Link>)
                        )
                    }
                </div>
                < RightBar showSidebar={true} />
            </main>
        </div>
    );
}