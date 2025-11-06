"use client";

import { useEffect, useState } from "react";
import Link from 'next/link';
import PersonAddAltIcon from "@mui/icons-material/PersonAddAlt";

import "../styles/rightbar.css"
import { useWS } from "../context/wsContext";
import { useParams } from "next/navigation";
import { Toaster, toast } from "sonner"


async function handleGroupRequest(invitationId, action, joinRequest, setJoinRequest) {
    try {
        const res = await fetch("http://localhost:8080/invitations/respond", {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                invitation_type: "join",
                invitation_id: invitationId,
                response: action,
            }),
        });

        if (!res.ok) {
            const errMsg = await res.text();
            throw new Error("Action failed: " + errMsg);
        }
        const data = await res.json();

        setJoinRequest((prev) => (prev || []).filter((req) => {
            req.InvitationID !== invitationId
        }));
    } catch (err) {
    }
}
export default function RightBarGroup({ onClick }) {
    const [friends, setFriends] = useState([])
    const [grpID, setGrpID] = useState('')
    // const [users, setusers] = useState([])
    const [onlineUsers, setonlineUsers] = useState([])
    // const [followRequest, setFollowRequest] = useState([])
    const [joinRequest, setJoinRequest] = useState([])
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



    const params = useParams();
    useEffect(() => {
        setGrpID(params.id);
    }, [params.id]);

    useEffect(() => {
        const temp = params.id
        async function fetchfriends(groupID) {
            try {
                const res = await fetch(`http://localhost:8080/api/fetchFriendsForGroups/${groupID}`, {
                    method: "GET",
                    credentials: "include",
                });
                const data = await res.json();
                if (!data) {
                    setFriends([]);
                    return;
                }
                if (data.error == "Unauthorized") {
                    window.location.href = "/login";
                    return;
                }
                setFriends(data);
            } catch (err) {
                console.error(err);
            }
        }
        fetchfriends(temp);
    }, [])


    useEffect(() => {
        if (!grpID) return;
        async function fetchJoinRequest(grpID) {
            try {
                const res = await fetch(`http://localhost:8080/api/fetchJoinRequests/${grpID}`, {
                    method: "GET",
                    credentials: "include",
                });
                const data = await res.json();
                setJoinRequest(data);
            } catch (err) {
                console.error(err);
            }
        }
        fetchJoinRequest(grpID);
    }, [grpID]);


    return (
        <div className="rightBar">
            <Toaster position="bottom-right" richColors />
            <div className="item">
                <span>Join Requests </span>
                {!joinRequest || joinRequest.length === 0 ? (
                    <h1>no Invitation for now</h1>
                ) : (
                    joinRequest.map((req) => (
                        <div key={req.InvitationID} className="user">
                            <div className="userInfo">
                                <div className="userDetails">
                                    <i className="fa-solid fa-people-group"></i>
                                    <span>{req.FirstName} {req.LastName}</span>
                                </div>
                                <div className="buttons">
                                    <button onClick={() => { handleGroupRequest(req.InvitationID, "accept", joinRequest, setJoinRequest) }} >accept</button>
                                    <button onClick={() => { handleGroupRequest(req.InvitationID, "reject") }} >reject</button>
                                </div>
                            </div>
                        </div>
                    ))
                )}
            </div>

            <div className="item">
                <div className="sectionss">
                    <span className="invitaionnn">
                        Invite Your Friends
                    </span>
                </div>
                <div>
                    {friends.length == 0 ? (
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
                                        toast.success("invitaiton sent ")
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
            </div>
        </div>
    );
}