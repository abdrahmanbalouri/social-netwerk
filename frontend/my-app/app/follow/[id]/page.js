"use client";
import { useSearchParams, useRouter, useParams } from "next/navigation";
import Navbar from "../../../components/Navbar.js";
import LeftBar from "../../../components/LeftBar.js";
import RightBar from "../../../components/RightBar.js";
import { useDarkMode } from "../../../context/darkMod.js";
import { useEffect, useState } from "react";
import "./follow.css";

export default function FollowPage() {
  const [followers, setFollowers] = useState([]);
  const [following, setFollowing] = useState([]);
  const { darkMode } = useDarkMode();
  const router = useRouter();
  const searchParams = useSearchParams();
  const tab = searchParams.get("tab") || "followers";
  const params = useParams();
  const userId = params.id;

  useEffect(() => {
    async function fetchData() {
      try {
        const url = tab === "followers"
          ? `http://localhost:8080/api/followers?id=${userId}`
          : `http://localhost:8080/api/following?id=${userId}`;

        const res = await fetch(url);
        const data = await res.json();

        tab === "followers" ? setFollowers(data) : setFollowing(data);
      } catch (error) {
        console.log(error);
      }
    }
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
              <ul>
                {followers?.map((f, i) => (
                  <li key={f.nickname}>{f.image} --- {f.nickname}  </li>

                ))}
              </ul>
            ) : (
              <ul>
                {following?.map((f, i) => (
                  <li key={f.nickname}>{f.image} --- {f.nickname}  </li>

                ))}
              </ul>
            )}
          </div>
        </div>
        <RightBar />
      </main>
    </div>
  );
}
