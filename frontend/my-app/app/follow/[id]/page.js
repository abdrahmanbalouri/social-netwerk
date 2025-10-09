"use client";
import { useSearchParams, useRouter, useParams } from "next/navigation";

import Navbar from "../../../components/Navbar.js";
import LeftBar from "../../../components/LeftBar.js";
import RightBar from "../../../components/RightBar.js";
import { useDarkMode } from "../../../context/darkMod.js";

export default function FollowPage() {
    
    const { darkMode } = useDarkMode();
  const router = useRouter();
  const searchParams = useSearchParams();
  const tab = searchParams.get("tab") || "followers";
    const params = useParams();

  const userId = params.id;

  const followers = ["Reda", "Amine", "Sara"];
  const following = ["Ali", "Youssef"];

  const handleTabChange = (newTab) => {

    router.push(`/follow/${userId}?tab=${newTab}`);
  };

  return (
    <div className={darkMode ? "theme-dark" : "theme-light"}>
      <Navbar />
      <main className="content">
        <LeftBar showSidebar={true} />

        <div className="pageContainer" style={{ padding: 20, flex: 6 }}>

          <div  className="tabButtons" >
            <button className="tabButton"
              onClick={() => handleTabChange("followers")}
              style={{
                backgroundColor: tab === "followers" ? "blue" : "white",
                color: tab === "followers" ? "white" : "black",
              }}
            >
              Followers
            </button>

            <button className="tabButton"
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
              <ul className="followersList">
                {followers.map((f, i) => (
                  <li key={i}>{f}</li>
                ))}
              </ul>
            ) : (
              <ul className="followingList">
                {following.map((f, i) => (
                  <li key={i}>{f}</li>
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
