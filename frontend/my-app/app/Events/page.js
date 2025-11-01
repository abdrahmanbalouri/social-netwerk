"use client";
import Navbar from "../../components/Navbar";
import LeftBar from "../../components/LeftBar.js";
import RightBar from "../../components/RightBar.js";
import { useDarkMode } from "../../context/darkMod.js";
import "../../styles/events-page.css";
import { useEffect, useState } from "react";


export default function EventsPage() {
    const { darkMode } = useDarkMode();
    const [myEvents, setMyEvents] = useState([]);
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
                console.log('event',data);
                
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
                <LeftBar />
                <div className="events-page" >
                    <h1>Events Page</h1>
                    <p>This is the Events page content.</p>
                </div>
                <RightBar showSidebar={true} />
            </main>
        </div>
    );
}