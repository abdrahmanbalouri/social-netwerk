"use client";

import { useEffect, useState } from "react";
import Link from 'next/link';
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
    const [invitations, setInvitations] = useState([]);
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

    async function handleGroupRequest(invitaitonId, action) {
        try {
            console.log("wast l handle request dyal l groups");
            const res = await fetch("http://localhost:8080/invitations/respond", {
                method: "POST",
                credentials: "include",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    invitation_id: invitaitonId,
                    response: action,
                }),
            });

            if (!res.ok) {
                const errMsg = await res.text();
                throw new Error("Action failed: " + errMsg);
            }
            const data = await res.json();
            setFollowRequest((prev) => prev.filter((req) => req.invitation_id !== invitaitonId));
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
                setInvitations(data)
            } catch (err) {
                console.error(err);
            }
        }
        fetchFollowRequest();
    }, []);
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
    console.log("groups invitations are :", followRequest);

    return (
        <div className="rightBar">
            <div className="item">
                <span>Groupe Invitation </span>
                {!followRequest || followRequest.length === 0 ? (
                    <h1>no Invitation for now</h1>
                ) : (
                    followRequest.map((group) => (
                        <div key={group.invitation_id} className="user">
                            <div className="userInfo">
                                <div className="userDetails">
                                    <i className="fa-solid fa-people-group"></i>
                                    <span>{group.title}</span>
                                </div>

                                <div className="buttons">
                                    <button onClick={() => { handleGroupRequest(group.invitation_id, "accept") }} >accept</button>
                                    <button onClick={() => { handleGroupRequest(group.invitation_id, "reject") }} >reject</button>
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
                                            sendMessage({ type: "invite_to_group", ReceiverId: user.id, groupID: grpID })
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