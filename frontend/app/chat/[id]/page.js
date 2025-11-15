"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Navbar from "../../../components/Navbar";
import LeftBar from "../../../components/LeftBar";
import UserBar from "../../../components/UserBar.js";
import { useDarkMode } from "../../../context/darkMod.js";
import ChatBox from "../../../components/chatBox.js";

export default function ChatPage() {
    const router = useRouter();
    const { darkMode } = useDarkMode();
    let { id } = useParams();
    const [user, setUser] = useState(null);
    useEffect(() => {
        async function fetchusers() {
            try {
                const res = await fetch("http://localhost:8080/api/communfriends", {
                    method: "GET",
                    credentials: "include",
                });
                if (!res.ok) {
                    throw new Error("Failed to fetch users");
                }
                const data = (await res.json()) || [];

                return data.map((user) => user.id);
            } catch (err) {
                console.error(err);
                return [];
            }
        }
        fetchusers().then((data) => {
            if (!data.includes(id) && id !== "0") {
                router.push("/chat/0");
            }
        });
    }, []);


    useEffect(() => {
        async function fetchUser() {
            try {
                const res = await fetch(
                    `http://localhost:8080/api/profile?userId=${id}`,
                    {
                        method: "GET",
                        credentials: "include",
                    }
                );
                const data = await res.json();

                setUser(data);
            } catch (err) {
                console.error("Error fetching user:", err);
            }
        }

        fetchUser();
    }, [id]);
    return (
        <div id="div" className={darkMode ? "theme-dark" : "theme-light"}>
            <Navbar />
            <main className="content">
                <LeftBar showSidebar={true} />

                {user ? (
                    <ChatBox user={user} />
                ) : (
                    <p className="loading">Loading user...</p>
                )}
                <UserBar />
            </main>
        </div>
    );
}
