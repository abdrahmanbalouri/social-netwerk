"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Navbar from "../../../components/Navbar";
import LeftBar from "../../../components/LeftBar";
import UserBar from "../../../components/UserBar.js";
import { useDarkMode } from "../../../context/darkMod.js";
import ChatBox from "../../../components/chatBox.js";

export default function ChatPage() {
    const { darkMode } = useDarkMode();
    const { id } = useParams();
    const [user, setUser] = useState(null);
    const router = useRouter();
    useEffect(() => {

        async function midle() {
            try {
                const response = await fetch("http://localhost:8080/api/me", {
                    credentials: "include",
                    method: "GET",
                });
                console.log(response);

                if (!response.ok) {

                    router.replace("/login");
                    return null;
                }
            } catch (error) {
                router.replace("/login");
                return null;

            }
        }
        midle()
    }, [])
    useEffect(() => {
        async function fetchUser() {
            try {
                const res = await fetch(`http://localhost:8080/api/profile?userId=${id}`, {
                    method: "GET",
                    credentials: "include",
                })
                const data = await res.json();

                setUser(data);
            } catch (err) {
                console.error("Error fetching user:", err);
            }
        }

        fetchUser();
    }, [id]);

    return (
        <div className={darkMode ? 'theme-dark' : 'theme-light'}>
            <Navbar />
            <main className="content">
                <LeftBar />
                {user ? <ChatBox user={user} /> : <p className="loading">Loading user...</p>}
                <UserBar />
            </main>
        </div>
    );
}
