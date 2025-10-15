"use client";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useDarkMode } from "../context/darkMod";
import { useProfile } from "../context/profile";
import { useWS } from "../context/wsContext.js";
import { useState, useEffect } from "react";
import { useChat } from "../context/chatContext.js";
import Notification from "./notofication.js";
import "../styles/navbar.css";

export default function Navbar({ onCreatePost }) {
  const router = useRouter();
  const { darkMode, toggle } = useDarkMode();
  const { Profile } = useProfile();
  const [cont, addnotf] = useState(0);
  const [data, notif] = useState({});
  const [show, cheng] = useState(false);
  const { addListener, removeListener, connected } = useWS();
  const { activeChatID } = useChat();
  useEffect(() => {
    if (!connected) return; // wait for connection

    const handleNotification = (data) => {
      console.log("Notification received:", data);
      if (activeChatID && data.from === activeChatID) return; // don't notify if chatting with the sender
      addnotf((prev) => prev + 1);
      notif(data.data || data);
      cheng(true);
      setTimeout(() => cheng(false), 4000);
    };

    addListener("notification", handleNotification);
    return () => removeListener("notification", handleNotification);
  }, [connected, addListener, removeListener, activeChatID]);

  return (
    <div className="navbar">
      <div className="left">
        <Link href="/home">
          <span>Social-Network</span>
        </Link>

        <div className="search">
          <i className="fa-solid fa-magnifying-glass"></i>
          <input type="text" placeholder="Search..." />
        </div>
      </div>
      {show && <Notification data={data} />}
      <div className="right">
        <i
          className={`fa-solid ${darkMode ? "fa-sun" : "fa-moon"}`}
          onClick={toggle}
        ></i>

        <div className="notification2">
          <i className="fa-solid fa-bell"></i>
          {cont > 0 && <span className="notif-count">{cont}</span>}
        </div>

        <i className="fa-solid fa-plus" onClick={onCreatePost}></i>

        <div className="user" onClick={() => router.push("/profile/0")}>
          <img
            src={
              Profile?.image
                ? `/uploads/${Profile.image}`
                : "/uploads/default.png"
            }
            alt="user avatar"
          />
        </div>
      </div>
    </div>
  );
}
