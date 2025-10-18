"use client";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { useDarkMode } from "../context/darkMod";
import { useProfile } from "../context/profile";
import { useWS } from "../context/wsContext.js";
import { useState, useEffect } from "react";
import { useChat } from "../context/chatContext.js";
import Notification from "./notofication.js";
import NOtBar from "./notfcationBar.js"
import "../styles/navbar.css";

export default function Navbar({ onCreatePost }) {
  const router = useRouter();
  const { darkMode, toggle } = useDarkMode();
  const { Profile } = useProfile();

  const [searchTerm, setSearchTerm] = useState("");
  const [searchResults, setSearchResults] = useState([]);
  const [showResults, setShowResults] = useState(false);

  const [showNotbar, chengBool] = useState(false)
  const [notData, setnot] = useState({})

  const [cont, addnotf] = useState(0);
  const [data, notif] = useState({});
  const [show, cheng] = useState(false);
  const { addListener, removeListener, connected } = useWS();
  const { activeChatID } = useChat();
  const id = useParams().id
  useEffect(() => {
    if (!connected) return; // wait for connection

    const handleNotification = (data) => {
      console.log("Notification received:", data);
      
      if (data.from !== id) return; // don't notify if chatting with the sender
      addnotf((prev) => prev + 1);
      notif(data.data || data);
      cheng(true);
      setTimeout(() => cheng(false), 4000);
    };

    addListener("notification", handleNotification);
    return () => removeListener("notification", handleNotification);
  }, [connected, addListener, removeListener, activeChatID]);


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


  useEffect(() => {
    if (!searchTerm.trim()) {
      setSearchResults([]);
      return;
    }

    const delay = setTimeout(async () => {
      try {
        const res = await fetch(`http://localhost:8080/api/searchUser?query=${encodeURIComponent(searchTerm)}`, {
          credentials: "include",
        });
        const data = await res.json();
        setSearchResults(data);
        setShowResults(true);
      } catch (err) {
        console.error(err);
      }
    }, 300); 

    return () => clearTimeout(delay);
  }, [searchTerm]);

  return (
    <div className="navbar">
      <div className="left">
        <Link href="/home">
          <span>Social-Network</span>
        </Link>

        <div className="search" style={{ position: "relative" }}>
          <i className="fa-solid fa-magnifying-glass"></i>
          <input
            type="text"
            placeholder="Search..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            onFocus={() => setShowResults(true)}
            onBlur={() => setTimeout(() => setShowResults(false), 150)}
          />

{showResults && Array.isArray(searchResults) && searchResults.length > 0 && (
            <div className="search-results" >
              {searchResults.map((user) => (
                <div
                  key={user.id}
                  onClick={() => router.push(`/profile/${user.id}`)}
                >
                  <img
                    src={`/uploads/${user.image || "default.png"}`}
                    alt=""
                    width="35"
                    height="35"
                  />
                  <span>{user.nickname || `${user.first_name} ${user.last_name}`}</span>
                </div>
              ))}
            </div>
          )}
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
