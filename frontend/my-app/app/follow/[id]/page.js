"use client";
import Link from "next/link.js";
import { useSearchParams, useRouter, useParams } from "next/navigation";
import { useEffect, useState } from "react";
import Navbar from "../../../components/Navbar.js";
import LeftBar from "../../../components/LeftBar.js";
import RightBar from "../../../components/RightBar.js";
import { useDarkMode } from "../../../context/darkMod.js";
import "../../../styles/follow.css";
import { useWS } from "../../../context/wsContext.js";

export default function FollowPage() {
  const [toast, setToast] = useState(null);
   const showToast = (message, type = "error", duration = 3000) => {
    setToast({ message, type });
    setTimeout(() => {
      setToast(null);
    }, duration);
  };



  const [followers, setFollowers] = useState([]);
  const [following, setFollowing] = useState([]);
  const { darkMode } = useDarkMode();
  const params = useParams();
  const searchParams = useSearchParams();
  const router = useRouter();
  const userId = params.id;
  const tab = searchParams.get("tab") || "followers";

  useEffect(() => {
    const fetchData = async () => {
      try {
        const url =
          tab === "followers"
            ? `http://localhost:8080/api/followers?id=${userId}`
            : `http://localhost:8080/api/following?id=${userId}`;
        const res = await fetch(url, {
          method: "GET",
          credentials: "include",
        });


      
        if (res.ok) {
          const data = await res.json();
          tab === "followers" ? setFollowers(data) : setFollowing(data);
      
        }else {
          throw new Error("Failed to fetch data");
        }
      } catch (error) {
       showToast("Failed to fetch data. Please try again later.");
      }
    };
    fetchData();
  }, [tab, userId]);

  const handleTabChange = (newTab) => {
    router.replace(`/follow/${userId}?tab=${newTab}`);
  };

  return (
    <div className={darkMode ? "theme-dark" : "theme-light"}>
      <Navbar />
      <main className="content">
        <LeftBar showSidebar={true} />
        <div className="pageContainer" style={{ padding: 20, flex: 6 }}>
          <div className="tabButtons">
            <button
              onClick={() => handleTabChange("followers")}
              className={
                tab === "followers" ? "activeButt" : ""
              }
            >
              Followers
            </button>
            <button
              onClick={() => handleTabChange("following")}
              className={
                tab === "following" ? "activeButt" : ""
              }
            >
              Following
            </button>
          </div>

          <div className="tabContent" style={{ marginTop: 20 }}>
            
            {tab === "followers" ? (
              <div className="itemUsers"> {followers?.map((user) => (
                <div key={user.id} className="userDiv">
                  <div className="userInfos">
                    <img src={user?.image ? `/uploads/${user.image}` : "/assets/default.png"} alt="user avatar" />

                    <div className="onlin" />
                    <Link href={`/profile/${user.id}`} >
                      <span>{user.first_name + " " + user.last_name}</span>
                    </Link>
                  </div>
                </div>
              )) ||   <div className="itemUsersR"> no user found </div> }
              </div>
            ) : (
              <div className="itemUsers">
                {following?.map((user) => (
                  <div key={user.id} className="userDiv">
                    <div className="userInfos">
                      <img src={user?.image ? `/uploads/${user.image}` : "/assets/default.png"} alt="user avatar" />

                      <div className="onlin" />
                      <Link href={`/profile/${user.id}`} >
                        <span>{user.first_name + " " + user.last_name}</span>
                      </Link>
                    </div>
                  </div>
                )) || <div className="itemUsersR">  no user found </div>}
              </div>
            )}
          </div>
        </div>
        <RightBar />
      </main>

   {toast && (
        <div className={`toast ${toast.type}`}>
          <span>{toast.message}</span>
          <button onClick={() => setToast(null)} className="toast-close">×</button>
        </div>
      )}
    </div>
  );
}
