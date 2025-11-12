"use client";
import { usePathname, useRouter } from "next/navigation";
import Link from "next/link";
import { useDarkMode } from "../context/darkMod";
import { useProfile } from "../context/profile";
import { useWS } from "../context/wsContext.js";
import { useState, useEffect } from "react";
// transient toast notification moved to GlobalNotification
import NotBar from "./notfcationBar.js"
import "../styles/navbar.css";

export default function Navbar() {
  const router = useRouter();
  const { darkMode, toggle } = useDarkMode();
  const { Profile } = useProfile();

  const [searchTerm, setSearchTerm] = useState("");
  const [searchResults, setSearchResults] = useState([]);
  const [showResults, setShowResults] = useState(false);

  const [showNotbar, chengBool] = useState(false)
  const [notData, setnot] = useState({})

  const [cont, addnotf] = useState(0);
  // notification transient data is handled globally by GlobalNotification
  const { addListener, removeListener, connected } = useWS();
  const pathname = usePathname();

  useEffect(() => {
    if (!connected) return; // wait for connection
    const handleNotification = (data) => {
      // update unread count and stored notification list for drop-down
      if (pathname !== `/chat/${data.from}` && pathname !== `/groups/${data.groupID}`) {

        addnotf((prev) => prev + 1);
        setnot(data);
      }
    };

    addListener("notification", handleNotification);
    return () => removeListener("notification", handleNotification);
  }, [connected, addListener, removeListener]);

  useEffect(() => {
    async function fetchNotifications() {
      try {
        const res = await fetch("http://localhost:8080/notifcation?bool=false", {
          method: "GET",
          credentials: "include",
        });
        if (!res.ok) {
          throw new Error(`Failed to fetch notifications: ${res.status}`);
        }

        const data = await res.json() || [];
        let t = 0
        data.map((not) => {
          if (!not.seen) {
            t++
          }
        })
        addnotf(t);
      } catch (error) {
        console.error("Error fetching notifications:", error);
      }
    }

    fetchNotifications();
  }, []);


  const notifications = async () => {
    try {
      const res = await fetch("http://localhost:8080/notifcation?bool=true", {
        method: "GET",
        credentials: "include",
      });

      if (!res.ok) {
        throw new Error(`Failed to fetch notifications: ${res.status}`);
      }

      const data = await res.json() || [];
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
    // clearTimeout(delay)    
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
  const disply = (e) => {
    const sideBar = document.getElementById("leftBar");
    const rightBar = document.getElementById("rightBar");
    console.log(e.target.id);

    if (e.target.id === "leftBtn") {
      sideBar.style.display = sideBar.style.display === "block" ? "none" : "block";
    } else if (e.target.id === "rightBtn") {
      rightBar.style.display = rightBar.style.display === "block" ? "none" : "block";
    }
    window.addEventListener('resize', () => {
      if (window.innerWidth > 768) {
        sideBar.style.display = 'block';
        rightBar.style.display = 'block';

      } else {
        sideBar.style.display = 'none';
        rightBar.style.display = 'none';
      }
    })

  };

  return (
    <div className="navbar">
      <div className="left">
        <i className="fa-solid fa-bars" id="leftBtn" onClick={(e) => disply(e)}></i>
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
          />

          {showResults && Array.isArray(searchResults) && searchResults.length > 0 && (
            <div className="search-results" >
              {searchResults.map((user) => (
                <div
                  key={user.id}
                  onClick={() => router.push(`/profile/${user.id}`)}
                >
                  <img
                    src={`${user?.image ? `/uploads/${user.image}` : "/assets/default.png"}`}
                    alt=""
                    width="35"
                    height="35"
                  />
                  <span>{`${user.first_name} ${user.last_name}`}</span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
      {showNotbar && <NotBar notData={notData} />}

      <div className="right">
        <i
          className={`fa-solid ${darkMode ? "fa-sun" : "fa-moon"}`}
          onClick={toggle}
        ></i>
        <div className="friend-requests" id="rightBtn" onClick={(e) => disply(e)}>
          <i className="fa-solid fa-user-friends" id="rightBtn"></i>
        </div>

        <div className="notification2" onClick={notifications}>
          <i className="fa-solid fa-bell"></i>
          {cont > 0 && <span className="notif-count">{cont}</span>}
        </div>


        <div className="user" onClick={() => router.push("/profile/0")}>
          <img
            src={
              Profile?.image
                ? `/uploads/${Profile.image}`
                : "/assets/default.png"
            }
            alt="user avatar"
          />
        </div>
      </div>
    </div>
  );
}
