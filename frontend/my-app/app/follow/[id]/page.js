"use client";
import Link from "next/link.js";
import { useSearchParams, useRouter, useParams } from "next/navigation";
import { useEffect, useState } from "react";
import Navbar from "../../../components/Navbar.js";
import LeftBar from "../../../components/LeftBar.js";
import RightBar from "../../../components/RightBar.js";
import { useDarkMode } from "../../../context/darkMod.js";
import "./follow.css";
import { middleware } from "../../../middleware/middelware.js";
import { useWS } from "../../../context/wsContext.js";


export default function FollowPage() {
  const [followers, setFollowers] = useState([]);
  const [following, setFollowing] = useState([]);
  const { darkMode } = useDarkMode();
  const params = useParams();
  const searchParams = useSearchParams();
  const router = useRouter();
  const userId = params.id;
  const tab = searchParams.get("tab") || "followers";
  const sendMessage = useWS()
  // Authentication check
  useEffect(() => {
    const checkAuth = async () => {
      const auth = await middleware();
      if (!auth) {
        router.push("/login");
        sendMessage({ type: "logout" })
      }
    }
    checkAuth();
  }, [])

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
        }
      } catch (error) {
        console.error(error);
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
              style={{
                backgroundColor: tab === "followers" ? "blue" : "white",
                color: tab === "followers" ? "white" : "black",
              }}
            >
              Followers
            </button>
            <button
              onClick={() => handleTabChange("following")}
              style={{
                backgroundColor: tab === "following" ? "blue" : "white",
                color: tab === "following" ? "white" : "black",
              }}
            >
              Following
            </button>
          </div>

          <div className="tabContent" style={{ marginTop: 20 }}>
            {tab === "followers" ? (
              <div className="itemUsers"> {followers?.map((user) => (
                <div key={user.id} className="userDiv">
                  <div className="userInfos">
                    <img src={user?.image ? `/uploads/${user.image}` : "/uploads/default.png"} alt="user avatar" />

                    <div className="onlin" />
                    <Link href={`/profile/${user.id}`} >
                      <span>{user.nickname}</span>
                    </Link>
                  </div>
                </div>
              )) || <div className="itemUsersR"> no user found </div>}
              </div>
            ) : (
              <div className="itemUsers">
                {following?.map((user) => (
                  <div key={user.id} className="userDiv">
                    <div className="userInfos">
                      <img src={user?.image ? `/uploads/${user.image}` : "/uploads/default.png"} alt="user avatar" />

                      <div className="onlin" />
                      <Link href={`/profile/${user.id}`} >
                        <span>{user.nickname}</span>
                      </Link>
                    </div>
                  </div>
                )) || <div className="itemUsersR"> no user found </div>}
              </div>
            )}
          </div>
        </div>
        <RightBar />
      </main>
    </div>
  );
}
