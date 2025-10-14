"use client";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useDarkMode } from "../context/darkMod";
import { useProfile } from "../context/profile";
import { useWS } from "../context/wsContext.js";
import { useState, useEffect } from "react";
import Notification from "./notofication.js";
import NOtBar from "./notfcationBar.js";

export default function Navbar({ onCreatePost }) {
  const router = useRouter();
  const { darkMode, toggle } = useDarkMode();
  const { Profile } = useProfile();
  const { ws, connected } = useWS();
  const [cont, addnotf] = useState(0);
  const [data, notif] = useState({})
  const [show, cheng] = useState(false)
  const [showNotbar, chengBool] = useState(false)
  const [notData, setnot] = useState({})

  useEffect(() => {

    if (!ws) return;
    ws.onmessage = (event) => {
      console.log(event);

      if (event.data) {
        try {
          const data = JSON.parse(event.data);
          if (data.type === "follow") {
            addnotf((prev) => prev + 1);
            notif(data)
            cheng(true)
            setTimeout(() => {
              cheng(false)

            }, 4000)

          }
        } catch (err) {
        }
      }
    };

    return () => {
      ws.onmessage = null;
    };
  }, [ws, connected]);


  const notifications = async () => {
    try {
      const res = await fetch("http://localhost:8080/notifcation", {
        method: "GET",
        credentials: "include",
      });

      if (!res.ok) {
        throw new Error(`Failed to fetch notifications: ${res.status}`);
      }

      const data = await res.json();
      setnot(data);
      addnotf(0);
      chengBool(!showNotbar);
    } catch (err) {
      console.error("Error fetching notifications:", err);
    }
  };


  return (
    <div className="navbar">
      <div className="left">
        <Link href="/home" style={{ textDecoration: "none" }}>
          <span>Social-Network</span>
        </Link>

        <div className="search">
          <i className="fa-solid fa-magnifying-glass"></i>
          <input type="text" placeholder="Search..." />
        </div>
      </div>
      {showNotbar  && <NOtBar notData={notData} />}
      {show && <Notification data={data} />}
      <div className="right">
        <i
          className={`fa-solid ${darkMode ? "fa-sun" : "fa-moon"}`}
          onClick={toggle}
        ></i>

        <div className="notification2" onClick={notifications}>
          <i className="fa-solid fa-bell"></i>
          {cont > 0 && <span className="notif-count">{cont}</span>}
        </div>

        <i className="fa-solid fa-plus" onClick={onCreatePost}></i>

        <div className="user" onClick={() => router.push("/profile/0")}>
          <img
            src={Profile?.image ? `/uploads/${Profile.image}` : "/uploads/default.png"}
            alt="user avatar"
          />
        </div>
      </div>
    </div>
  );
}
