"use client";

import { useEffect, useState } from "react";
import Link from 'next/link';
import "../styles/userbar.css"
export default function UserBar() {
    const [users, setusers] = useState([])
    useEffect(() => {
        async function fetchusers() {
            try {
                const res = await fetch("http://localhost:8080/api/GetUsersHandler", {
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
                                <Link href={`/chat/${user.id}`} className="userLink">
                                    <div className="userDetails">
                                        <img
                                            src={user?.image ? `/uploads/${user.image}` : "/uploads/default.png"}
                                            alt="user avatar"
                                        />
                                        <div className="online" />
                                        <span>{user.nickname}</span>
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