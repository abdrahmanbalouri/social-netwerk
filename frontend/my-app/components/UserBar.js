"use client";

import { useEffect, useState } from "react";
import Link from 'next/link';
import "../styles/userbar.css"
import { useWS } from "../context/wsContext";
export default function UserBar() {
    const [users, setusers] = useState([])
    const [onlineUsers, setonlineUsers] = useState([])
    const { sendMessage, addListener, removeListener } = useWS();
    useEffect(() => {
        const handleOlineUser = (data) => {
            setonlineUsers(data.users)
        }
        sendMessage({ type: "online_list" })
        addListener("online_list", handleOlineUser)
        return () => removeListener("online_list", handleOlineUser)
    }, [addListener, removeListener])
    useEffect(() => {
        const handleLogout = (data) => {
            let useroff = data.userID
            let arr = onlineUsers.filter((id) => {
                return id !== useroff
            })
            setonlineUsers([...arr])
        }

        addListener("logout", handleLogout)
        return () => removeListener("logout", handleLogout)
    }, [addListener, removeListener])

    useEffect(() => {
        async function fetchusers() {
            try {
                const res = await fetch("http://localhost:8080/api/communfriends", {
                    method: "GET",
                    credentials: "include",
                });
                if (!res.ok) {
                    throw new Error("Failed to fetch posts");
                }
                const data = await res.json();



                setusers(data);
            } catch (err) {
                console.error(err);
            }
        }

        fetchusers();
    }, []);

    return (
        <div className="userBar">
            <div className="usrs">
                <span>Online Friends</span>
                {!users || users.length === 0 ? (
                    <h1>no users for now</h1>
                ) : (
                    users.map((user) => (
                        <div key={user.id} className="user">
                            <div className="userInfo">
                                <Link href={`/chat/${user.id}`}>
                                    <div className="userDetails">
                                        <img
                                            src={user?.image ? `/uploads/${user.image}` : "/assets/default.png"}
                                            alt="user avatar"
                                        />
                                        <div className={onlineUsers.includes(user.id) ? "online" : "offline"} />
                                        <span>{user.first_name +" " + user.last_name}</span>
                                    </div>
                                </Link>
                            </div>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
}