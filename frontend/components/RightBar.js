"use client";

import { useEffect, useState } from "react";
import Link from 'next/link';
import ChatIcon from '@mui/icons-material/Chat';
import "../styles/rightbar.css"
import { useWS } from "../context/wsContext";
import { ShieldXIcon } from "lucide-react";
export default function RightBar() {


 const [toast, setToast] = useState(null)

   
    const showToast = (message, type = "error", duration = 3000) => {
    setToast({ message, type });
    setTimeout(() => {
      setToast(null);
    }, duration);
  };

  const [friends, setFriends] = useState([])
  const [users, setusers] = useState([])
  const [onlineUsers, setonlineUsers] = useState([])
  const [activeTab, setActiveTab] = useState("friends");
  const [activeTabRequests, setActiveTabRequests] = useState("followRequests");

  const [groupeInvitation, setgroupeInvitation] = useState([])

  const [followRequest, setFollowRequest] = useState([])
  const { sendMessage, addListener, removeListener } = useWS();





  useEffect(() => {
    async function fetchGroupeInvitation() {
      try {
        const res = await fetch("http://localhost:8080/api/fetchGroupInvitation", {
          method: "GET",
          credentials: "include",
        });

        if (!res.ok) {
          throw new Error("Failed to fetch posts");
        }
        const data = await res.json();

        setgroupeInvitation(data);
      } catch (err) {
        showToast("Failed to fetch group invitations");
      }
    }

    fetchGroupeInvitation();
  }, []);


  async function fetchFollowRequest() {
    try {
      const res = await fetch("http://localhost:8080/api/followRequest", {
        method: "GET",
        credentials: "include",
      });

      if (res.ok) {
        const data = await res.json();
        setFollowRequest(data);
     
      } else {
       
        return;
      }

    } catch (err) {
     
    }
  }



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
        fetchFollowRequest()
        const errMsg = await res.text();
        throw new Error("Action failed: " + errMsg);
      }


      setFollowRequest((prev) => prev.filter((user) => user.id !== userId));
    } catch (err) {
      showToast("Failed to process follow request");
    }
  }

  async function handleGroupRequest(invitaitonId, action) {
    try {
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
      setgroupeInvitation((prev) => (prev || []).filter((req) => req.invitation_id !== invitaitonId));
    } catch (err) {
      showToast("Failed to process group invitation");
    }
  }


  useEffect(() => {
    fetchFollowRequest();
  }, []);


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
showToast("Failed to fetch users");
   }
    }

    fetchusers();
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



  return (



    <div className="rightBar">
     {toast && (
          <div className={`toast ${toast.type}`}>
            <span>{toast.message}</span>
            <button onClick={() => setToast(null)} className="toast-close">×</button>
          </div>
        )}
      <div className="item">
        <div className="sections">
          <h3
            className={activeTabRequests === "followRequests" ? "active" : ""}
            onClick={() => setActiveTabRequests("followRequests")}
          >
            Follow Requests
          </h3>
          <h3
            className={activeTabRequests === "groupeInvitation" ? "active" : ""}
            onClick={() => setActiveTabRequests("groupeInvitation")}
          >
            Groupe Invitation
          </h3>
        </div>


        {activeTabRequests === "followRequests" && (

          <div>
            {!followRequest || followRequest.length === 0 ? (<p>no follow reuqets found  </p>) : (
              followRequest.map((user) => (
                <div key={user.id} className="user">
                  <div className="userInfo">
                    <div className="userDetails">
                      <Link href={`/profile/${user.id}`} className="userLink">
                        <img
                          src={user?.image ? `/uploads/${user.image}` : "/assets/default.png"}
                          alt="user avatar"
                        />
                      </Link>
                      <Link href={`/profile/${user.id}`}>
                        <span>{user.first_name + " " + user.last_name}</span>
                      </Link>
                    </div>

                    <div className="buttons">
                      <button onClick={() => { handleFollowRequest(user.id, "accept") }} >accept</button>
                      <button onClick={() => { handleFollowRequest(user.id, "reject") }} >reject</button>
                    </div>
                  </div>
                </div>
              ))
            )}


          </div>
        )}

        {activeTabRequests === "groupeInvitation" && (
          <div>

            {!groupeInvitation ? (<p> no groupe invitation found </p>) :
              (
                groupeInvitation.map((group) => (
                  <div key={group.group_id} className="user">
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
          <h3
            className={activeTab === "all" ? "active" : ""}
            onClick={() => setActiveTab("all")}
          >
            All
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
                    <div className="">
                      <Link href={`/chat/${user.id}`}><ChatIcon className="userIcon" /></Link>
                    </div>
                  </div>
                </div>
              ))
            )}
          </div>
        )}
        {activeTab === "all" && (
          <div>
            {!users ? (
              <p>No users found</p>
            ) : (
              users.map((user) => (
                <div key={user.id} className="user">
                  <div className="userInfo">
                    <div className="userDetails">
                      <Link href={`/profile/${user.id}`} className="userLink">
                        <img
                          src={user?.image ? `/uploads/${user.image}` : "/assets/default.png"}
                          alt="user avatar"
                        />
                      </Link>
                      <Link href={`/profile/${user.id}`} className="userLink">
                        <span>{user.first_name + " " + user.last_name}</span>
                      </Link>
                    </div>
                    <Link href={`/profile/${user.id}`} className="userLink">
                      <i className="fa-solid fa-user"></i>
                    </Link>
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