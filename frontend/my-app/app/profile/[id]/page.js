"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Navbar from "../../../components/Navbar.js";
import LeftBar from "../../../components/LeftBar.js";
import RightBar from "../../../components/RightBar.js";
import Post from "../../../components/Post.js";
import Comment from "../../../components/coment.js";
import { useDarkMode } from "../../../context/darkMod.js";
import { useProfile } from "../../../context/profile.js";
import "./profile.css";

export default function ProfilePage() {
  const { darkMode } = useDarkMode();
  const { profile } = useProfile();
  const params = useParams();
  const router = useRouter();

  const userId = params.id;

  const [userProfile, setUserProfile] = useState(null);
  const [posts, setPosts] = useState([]);
  const [comments, setComments] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function init() {
      try {
        const me = await fetch("http://localhost:8080/api/me", {
          credentials: "include",
        });
        if (!me.ok) return router.replace("/login");

        const resProfile = await fetch(
          `http://localhost:8080/api/profile?userId=${userId}`,
          { credentials: "include" }
        );
        if (resProfile.ok) {
          setUserProfile(await resProfile.json());
        }

        const resPosts = await fetch(
          `http://localhost:8080/api/Getpostbyuser/${userId}`,
          { credentials: "include" }
        );
        if (resPosts.ok) {
          const data = await resPosts.json();
          setPosts(data || []);
        }
      } catch (err) {
        console.error("Error loading profile:", err);
        router.replace("/login");
      } finally {
        setLoading(false);
      }
    }

    init();
  }, [userId, router]);

  async function followUser() {
    try {
      await fetch(`http://localhost:8080/api/follow?id=${userId}`, {
        method: "POST",
        credentials: "include",
      });
      setUserProfile((prev) => ({
        ...prev,
        isFollowed: !prev.isFollowed,
      }));
    } catch (err) {
      console.error("Error following user:", err);
    }
  }

  async function GetComments(post) {
    try {
      const res = await fetch(`http://localhost:8080/api/comments/${post.id}`, {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        setComments(data);
      }
    } catch (err) {
      console.error("Error fetching comments:", err);
    }
  }

  async function refreshComments(post) {
    await GetComments(post);
  }

  return (
    <div className={`app ${darkMode ? "dark" : ""}`}>
      <Navbar />
      <div className="container">
        <LeftBar />

        <div className="main">
          {loading ? (
            <p>Loading profile...</p>
          ) : !userProfile ? (
            <p>User not found</p>
          ) : (
            <div className="profile-container">
              <div className="profile-header">
                <img
                  src={userProfile.image || "/default-avatar.png"}
                  alt="profile"
                  className="profile-avatar"
                />
                <h2>{userProfile.nickname}</h2>
                {profile && profile.id !== userProfile.id && (
                  <button className="follow-btn" onClick={followUser}>
                    {userProfile.isFollowed ? "Unfollow" : "Follow"}
                  </button>
                )}
              </div>

              <hr />

              <div className="posts-section">
                {posts.length === 0 ? (
                  <p>No posts yet.</p>
                ) : (
                  posts.map((post) => (
                    <div key={post.id} className="post-block">
                      <Post post={post} onGetComments={() => GetComments(post)} />
                      <Comment
                        comments={comments}
                        postId={post.id}
                        onRefresh={() => refreshComments(post)}
                      />
                    </div>
                  ))
                )}
              </div>
            </div>
          )}
        </div>

        <RightBar />
      </div>
    </div>
  );
}
