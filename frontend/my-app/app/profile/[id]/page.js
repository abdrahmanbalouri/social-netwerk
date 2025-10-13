"use client";
import { useEffect, useState, useRef } from 'react';
import Navbar from '../../../components/Navbar.js';
import { useDarkMode } from '../../../context/darkMod.js';
import './profile.css';
import LeftBar from '../../../components/LeftBar.js';
import RightBar from '../../../components/RightBar.js';
import { useParams, useRouter } from 'next/navigation.js';
import Post from '../../../components/Post.js';
import Comment from '../../../components/coment.js';
import { useProfile } from '../../../context/profile.js';
import LanguageIcon from "@mui/icons-material/Language";
import EmailOutlinedIcon from "@mui/icons-material/EmailOutlined";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import ProfileCardEditor from '../../../components/ProfileCardEditor.js';
import { useWS } from "../../../context/wsContext.js";
import Link from 'next/link';

export default function Profile() {
  const { Profile } = useProfile();

  const { darkMode } = useDarkMode();

  const [showSidebar, setShowSidebar] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const params = useParams();
  const router = useRouter();
  const userId = Number(params.id);
  const [posts, setPosts] = useState([]);
  const [showComments, setShowComments] = useState(false);
  const [selectedPost, setSelectedPost] = useState(null);
  const [comment, setComment] = useState([]);
  const commentsModalRef = useRef(null);
  const { ws, connected } = useWS();

  const sendMsg = () => {
    const payload = {
      receiverId: params.id,
      messageContent: "",
      type: "follow",
    };

    if (connected && ws) {
      ws.send(JSON.stringify(payload));
    }
  };

  const [theprofile, setProfile] = useState(null);
  useEffect(() => {

    async function midle() {
      try {
        const response = await fetch("http://localhost:8080/api/me", {
          credentials: "include",
          method: "GET",
        });

        if (!response.ok) {
          router.replace("/login");
          return null;
        }
      } catch (error) {
        router.replace("/login");
        return null;

      }
    }
    midle()



  }, [])

  async function loadProfile() {
    try {
      const res = await fetch(
        `http://localhost:8080/api/profile?userId=${params.id}`,
        {
          method: "GET",
          credentials: "include",
        }
      );
      if (res.ok) {
        const json = await res.json();
        console.log(json);

        setProfile(json);
      }
    } catch (err) {
      console.error("loadProfile", err);
    }
  }

  useEffect(() => {
    loadProfile();
  }, []);
  // Fetch posts for this profile user
  useEffect(() => {
    async function fetchUserPosts() {
      try {
        const res = await fetch(
          `http://localhost:8080/api/Getallpost?userId=${params.id}`,
          {
            method: "GET",
            credentials: "include",
          }
        );
        if (!res.ok) {
          throw new Error("Failed to fetch user posts");
        }
        const data = await res.json();

        setPosts(Array.isArray(data) ? data : []);
      } catch (err) {
        console.error("Error fetching user posts:", err);
      }
    }
    if (params.id) {
      fetchUserPosts();
    }
  }, [params.id]);

  // Fetch comments for a specific post (like home page)
  async function GetComments(post) {
    try {
      setSelectedPost({
        id: post.id,
        title: post.title || post.post_title || "Post",
      });
      const res = await fetch(
        `http://localhost:8080/api/Getcomments/${post.id}`,
        {
          method: "GET",
          credentials: "include",
        }
      );
      if (!res.ok) {
        throw new Error("Failed to fetch comments");
      }
      const data = await res.json();

      let comments = [];
      if (Array.isArray(data)) {
        comments = data;
      } else if (
        data &&
        typeof data === "object" &&
        data.comments &&
        Array.isArray(data.comments)
      ) {
        comments = data.comments;
      } else if (data && typeof data === "object") {
        comments = [data];
      }
      comments = comments.map((comment) => ({
        id: comment.id || Math.random(),
        author: comment.author || comment.username || "Anonymous",
        content: comment.content || comment.text || "",
        created_at:
          comment.created_at || comment.createdAt || new Date().toISOString(),
      }));
      setComment(comments);
      setShowComments(true);
    } catch (err) {
      setComment([]);
      setSelectedPost({ id: post.id, title: post.title || "Post" });
      setShowComments(true);
    }
  }

  // Refresh comments after posting a new comment
  async function refreshComments() {
    if (!selectedPost?.id) return;
    try {
      const res = await fetch(
        `http://localhost:8080/api/Getcomments/${selectedPost.id}`,
        {
          method: "GET",
          credentials: "include",
        }
      );
      if (res.ok) {
        const data = await res.json();
        let comments = [];
        if (Array.isArray(data)) {
          comments = data;
        } else if (data && data.comments && Array.isArray(data.comments)) {
          comments = data.comments;
        } else if (data) {
          comments = [data];
        }
        setComment(comments);
      }
    } catch (err) {
      // ignore
    }
  }

  // Close comments modal and reset state
  function closeComments() {
    setShowComments(false);
    setSelectedPost(null);
    setComment([]);
  }
  const [showPrivacy, setShowPrivacy] = useState(false);

  function handleShowPrivacy() {
    setShowPrivacy(!showPrivacy);
  }

  async function followUser() {
    try {
      let res = await fetch(
        `http://localhost:8080/api/follow?id=${params.id}`,
        {
          method: "POST",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ userId: userId }),
        }
      );
      if (res.ok) {
        let followw = await res.json();

        // Update state directly, ma t3melsh direct DOM manipulation
        setProfile((prevProfile) => ({
          ...prevProfile,
          followers: followw.followers,
          following: followw.following,
          isFollowing: followw.isFollowed,
        }));
      }
    } catch (error) {
      console.error("Error following user:", error);
    }
  }

  if (!theprofile) {
    return (
      <div className={darkMode ? "theme-dark" : "theme-light"}>
        <Navbar
          onCreatePost={() => setShowModal(true)}
          onToggleSidebar={() => setShowSidebar(!showSidebar)}
        />
        <main className="content">
          <LeftBar showSidebar={showSidebar} />
          <div className="profile">
            <div className="images">
              <div
                className="cover"
                style={{ width: "100%", height: 300, background: "#eee" }}
              />
              <div
                className="profilePic"
                style={{
                  width: 200,
                  height: 200,
                  background: "#ccc",
                  borderRadius: "50%",
                  margin: "0 auto",
                  marginTop: -100,
                }}
              />
            </div>
            <div className="profileContainer">
              <div className="uInfo">
                <div className="center">
                  <span>Loading...</span>
                </div>
              </div>
            </div>
            <div className="posts" style={{ marginTop: 20 }}>
              {posts.length === 0 ? (
                <p>No posts available</p>
              ) : (
                posts.map((post) => (
                  <Post
                    key={post.id}
                    post={post}
                    onGetComments={() => GetComments(post)}
                  />
                ))
              )}
            </div>
          </div>
          <RightBar />
        </main>
      </div>
    );
  }

  return (
    <div className={darkMode ? "theme-dark" : "theme-light"}>
      <Navbar
        onCreatePost={() => setShowModal(true)}
        onToggleSidebar={() => setShowSidebar(!showSidebar)}
      />
      <main className="content">
        <LeftBar showSidebar={showSidebar} />
        <div className="profile">
          <div className="images">
            <img
              src={
                theprofile.cover
                  ? `/uploads/${theprofile.cover}`
                  : "https://images.pexels.com/photos/13440765/pexels-photo-13440765.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2"
              }
              alt=""
              className="cover"
            />
            <img
              src={
                theprofile.image
                  ? `/uploads/${theprofile.image}`
                  : "/uploads/default.png"
              }
              alt="profile picture"
              className="profilePic"
            />
          </div>
          <div className="profileContainer">
            <div className="uInfo">
              <div className="left">
                <Link href={`/follow/${theprofile.id}?tab=followers`}>

                  <p>
                    following
                    <strong id="following">{theprofile.following} </strong>
                  </p>
                </Link>


                <Link href={`/follow/${theprofile.id}?tab=following`}>
                  <p>
                    followers
                    <strong id="followers"> {theprofile.followers}</strong>
                  </p>

                </Link>
              </div>
              <div className="center">
                <span>{theprofile.nickname}</span>
                <div className="info">
                  <div className="item">
                    <LanguageIcon />
                    <span>{theprofile.about}</span>
                  </div>
                </div>
                {Profile && Profile.id !== theprofile.id && (
                  <button
                    id="FollowBtn"
                    onClick={followUser}
                    style={{
                      backgroundColor: theprofile.isFollowing
                        ? "blue"
                        : "white",
                      color: theprofile.isFollowing ? "white" : "black",
                      border: "1px solid #ccc",
                      padding: "8px 16px",
                    }}
                  >
                    {Profile &&
                      (theprofile.isFollowing
                        ? "Unfollow"
                        : theprofile.privacy === "private" &&
                          !theprofile.isFollowing
                          ? "Pending"
                          : "Follow")}
                  </button>
                )}
              </div>
              <div className="right">
                {Profile && Profile.id !== theprofile.id ? (
                  <EmailOutlinedIcon />
                ) : (
                  <MoreVertIcon onClick={handleShowPrivacy} />
                )}
              </div>
              {showPrivacy && <ProfileCardEditor showPrivacy={showPrivacy} />}
            </div>
            <div className="posts" style={{ marginTop: 20 }}>
              {posts.length === 0 ? (
                <p>No posts available</p>
              ) : (
                posts.map((post) => (
                  <Post
                    key={post.id}
                    post={post}
                    onGetComments={() => GetComments(post)}
                  />
                ))
              )}
            </div>
          </div>
        </div>
        <RightBar />
      </main>
      {/* Comments Modal */}
      {showComments && (
        <Comment
          comments={comment}
          isOpen={showComments}
          onClose={closeComments}
          postId={selectedPost?.id}
          postTitle={selectedPost?.title}
          onCommentChange={refreshComments}
        />
      )}
    </div>
  );
}
