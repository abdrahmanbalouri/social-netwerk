"use client";

import { useEffect, useState } from "react";
import Link from 'next/link';
import ChatIcon from '@mui/icons-material/Chat';
import PersonAddAltIcon from "@mui/icons-material/PersonAddAlt";

import "../styles/rightbar.css"
import { useWS } from "../context/wsContext";
import { useParams } from "next/navigation";

export default function RightBarGroup({ onClick }) {
    const [friends, setFriends] = useState([])
    const [grpID, setGrpID] = useState('')
    // const [users, setusers] = useState([])
    const [onlineUsers, setonlineUsers] = useState([])
    const [activeTab, setActiveTab] = useState("friends");
    const [followRequest, setFollowRequest] = useState([])
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

    async function handleFollowRequest(userId, action) {
        try {
            const res = await fetch("http://localhost:8080/api/followRequest/action", {
                method: "POST",
                credentials: "include",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    id: userId,
                    action: action,
                }),
            });

            if (!res.ok) {
                const errMsg = await res.text();
                throw new Error("Action failed: " + errMsg);
            }

            const data = await res.json();

            setFollowRequest((prev) => prev.filter((user) => user.id !== userId));
        } catch (err) {
        }
    }

    useEffect(() => {
        async function fetchFollowRequest() {
            try {
                const res = await fetch("http://localhost:8080/api/groupeInvitation", {
                    method: "GET",
                    credentials: "include",
                });

                if (!res.ok) {
                    throw new Error("Failed to fetch posts");
                }
                const data = await res.json();

                console.log(data);
                setFollowRequest(data);
            } catch (err) {
                console.error(err);
            }
        }

        fetchFollowRequest();
    }, []);


    // useEffect(() => {
    //     async function fetchusers() {
    //         try {
    //             const res = await fetch("http://localhost:8080/api/GetUsersHandler", {
    //                 method: "GET",
    //                 credentials: "include",
    //             });

    //             if (!res.ok) {
    //                 throw new Error("Failed to fetch posts");
    //             }
    //             const data = await res.json();

    //             setusers(data);
    //         } catch (err) {
    //             console.error(err);
    //         }
    //     }

    //     fetchusers();
    // }, []);
    useEffect(() => {
        async function fetchfriends() {
            try {
                const res = await fetch("http://localhost:8080/api/communfriends", {
                    method: "GET",
                    credentials: "include",
                });
                if (!res.ok) {
                    throw new Error("Failed to fetch posts");
                }
                const data = await res.json();
                setFriends(data);
            } catch (err) {
                console.error(err);
            }
        }
        fetchfriends();
    }, [])

    const params = useParams();
    useEffect(() => {
        setGrpID(params.id);
    }, [params.id]);



    return (
        <div className="rightBar">
            <div className="item">
                <span>Groupe Invitation </span>
                {!followRequest || followRequest.length === 0 ? (
                    <h1>no Invitation for now</h1>
                ) : (
                    followRequest.map((group) => (
                        <div key={group.id} className="user">
                            <div className="userInfo">
                                <div className="userDetails">
                                    {/*  <Link href={`/profile/${group.id}`} className="userLink">  */}
                                  {/*   <img
                                        src={group?.image ? `/uploads/${group.image}` : "/uploads/default.png"}
                                        alt="user avatar"
                                    /> */}
                                                              <i className="fa-solid fa-people-group"></i>

                                    {/*</Link>*/}
                                    {/* <Link href={`/profile/${group.id}`}> */}
                                    <span>{group.title}</span>
                                    {/* </Link> */}
                                </div>

                                <div className="buttons">
                                    <button onClick={() => { handleFollowRequest(group.id, "accept") }} >accept</button>
                                    <button onClick={() => { handleFollowRequest(group.id, "reject") }} >reject</button>
                                </div>
                            </div>
                        </div>
                    ))
                )}
            </div>

            <div className="item">
                <div className="sections">
                    <h3
                        className={activeTab === "friends" ? "active" : ""}
                        onClick={() => setActiveTab("friends")}
                    >
                        Friends
                    </h3>
                </div>
                {activeTab === "friends" && (
                    <div>
                        {!friends ? (
                            <p>No friends yet</p>
                        ) : (
                            friends.map((user) => (
                                <div key={user.id} className="user">
                                    <div className="userInfo">
                                        <div className="userDetails">
                                            <Link href={`/profile/${user.id}`}>
                                                <img
                                                    src={user?.image ? `/uploads/${user.image}` : "/assets/default.png"}
                                                    alt="user avatar"
                                                />
                                            </Link>
                                            <div className={onlineUsers.includes(user.id) ? "online" : "offline"} />
                                            <Link href={`/profile/${user.id}`}>
                                                <span>{user.first_name + " " + user.last_name}</span>
                                            </Link>
                                        </div>
                                        <div onClick={() => {
                                            onClick(user.id, grpID)
                                        }}>
                                            <PersonAddAltIcon className="userIcon" />
                                        </div>
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
                )}
            </div>
        </div>
    );
}