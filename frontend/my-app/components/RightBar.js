"use client";

import { useEffect, useState } from "react";
import Link from 'next/link';
import ChatIcon from '@mui/icons-material/Chat';
import "../styles/rightbar.css"
export default function RightBar() {
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
    <div className="rightBar">
      <div className="item">
        <span>Suggestions For You</span>
        <div className="user">
          <div className="userInfo">
            <img
              src="https://images.pexels.com/photos/4881619/pexels-photo-4881619.jpeg?auto=compress&cs=tinysrgb&w=1600"
              alt=""
            />
            <span>Jane Doe</span>
          </div>
          <div className="buttons">
            <button>follow</button>
            <button>dismiss</button>
          </div>
        </div>
        <div className="user">
          <div className="userInfo">
            <img
              src="https://images.pexels.com/photos/4881619/pexels-photo-4881619.jpeg?auto=compress&cs=tinysrgb&w=1600"
              alt=""
            />
            <span>Jane Doe</span>
          </div>
          <div className="buttons">
            <button>follow</button>
            <button>dismiss</button>
          </div>
        </div>
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
                  <Link href={`/profile/${user.id}`}>
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