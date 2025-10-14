"use client";

import { useEffect, useState } from "react";
import Link from 'next/link';
import ChatIcon from '@mui/icons-material/Chat';
import "../styles/rightbar.css"
export default function RightBar() {
  const [users, setusers] = useState([])
  const [followRequest, setFollowRequest] = useState([])


  async function handleFollowRequest(userId, action) {


    console.log("fezfzelkfjzelfjlezjflezjflez", userId , action);
    
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
      console.log("✅ Action result:", data);

      setFollowRequest((prev) => prev.filter((user) => user.id !== userId));
    } catch (err) {
      console.error("❌ Error:", err);
    }
  }




  useEffect(() => {
    async function fetchFollowRequest() {
      try {
        const res = await fetch("http://localhost:8080/api/followRequest", {
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
    <div className="rightBar">
      <div className="item">
        <span>Suggestions For You</span>

        {!followRequest || followRequest.length === 0 ? (
          <h1>no users for now</h1>
        ) : (
          followRequest.map((user) => (
            <div key={user.id} className="user">
              <div className="userInfo">
                <div className="userDetails">
                  <img
                    src={user?.image ? `/uploads/${user.image}` : "/uploads/default.png"}
                    alt="user avatar"
                  />
                  <div className="online" />
                  <Link href={`/profile/${user.id}`} className="userLink">
                    <span>{user.nickname}</span>
                  </Link>
                </div>

                <div className="buttons">
                  <button onClick={() => { handleFollowRequest(user.id , "accept") }} >accept</button>
                  <button onClick={() => { handleFollowRequest(user.id  ,  "reject") }} >reject</button>
                </div>
              </div>
            </div>
          ))
        )}

      </div>

      <div className="item">
        <span>Friends</span>
        {!users || users.length === 0 ? (
          <h1>no users for now</h1>
        ) : (
          users.map((user) => (
            <div key={user.id} className="user">
              <div className="userInfo">
                <div className="userDetails">
                  <img
                    src={user?.image ? `/uploads/${user.image}` : "/uploads/default.png"}
                    alt="user avatar"
                  />
                  <div className="online" />
                  <Link href={`/profile/${user.id}`} className="userLink">
                    <span>{user.nickname}</span>
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
    </div>
  );
}