"use client";
import FacebookTwoToneIcon from '@mui/icons-material/FacebookTwoTone';
import InstagramIcon from '@mui/icons-material/Instagram';
import TwitterIcon from '@mui/icons-material/Twitter';
import LinkedInIcon from '@mui/icons-material/LinkedIn';
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
import PlaceIcon from "@mui/icons-material/Place";
import LanguageIcon from "@mui/icons-material/Language";
import EmailOutlinedIcon from "@mui/icons-material/EmailOutlined";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import ProfileCardEditor from '../../../components/ProfileCardEditor.js';

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


  const [profile, setProfile] = useState(null);

  async function loadProfile() {
    try {
      const res = await fetch(`http://localhost:8080/api/profile?userId=${params.id}`, {
        method: "GET",
        credentials: "include",
      });
      if (res.ok) {
        const json = await res.json();
        console.log('11', json);

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
        const res = await fetch(`http://localhost:8080/api/Getallpost?userId=${params.id}`, {
          method: "GET",
          credentials: "include",
        });
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
        title: post.title || post.post_title || "Post"
      });
      const res = await fetch(`http://localhost:8080/api/Getcomments/${post.id}`, {
        method: "GET",
        credentials: "include",
      });
      if (!res.ok) {
        throw new Error("Failed to fetch comments");
      }
      const data = await res.json();
      let comments = [];
      if (Array.isArray(data)) {
        comments = data;
      } else if (data && typeof data === 'object' && data.comments && Array.isArray(data.comments)) {
        comments = data.comments;
      } else if (data && typeof data === 'object') {
        comments = [data];
      }
      comments = comments.map(comment => ({
        id: comment.id || Math.random(),
        author: comment.author || comment.username || "Anonymous",
        content: comment.content || comment.text || "",
        created_at: comment.created_at || comment.createdAt || new Date().toISOString()
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
      const res = await fetch(`http://localhost:8080/api/Getcomments/${selectedPost.id}`, {
        method: "GET",
        credentials: "include",
      });
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

  function followUser() {
    fetch(`http://localhost:8080/api/follow?id=${params.id}`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ userId: userId }),
    })
      .then((res) => {
        
      })
      .catch((err) => {
        console.error('Error following user:', err);
      });
  }

  if (!profile) {
    return (
      <div className={darkMode ? 'theme-dark' : 'theme-light'}>
        <Navbar
          onCreatePost={() => setShowModal(true)}
          onToggleSidebar={() => setShowSidebar(!showSidebar)}
        />
        <main className="content">
          <LeftBar showSidebar={showSidebar} />
          <div className="profile">
            <div className="images">
              <div className="cover" style={{ width: '100%', height: 300, background: '#eee' }} />
              <div className="profilePic" style={{ width: 200, height: 200, background: '#ccc', borderRadius: '50%', margin: '0 auto', marginTop: -100 }} />
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
  const data = profile;
  console.log(Profile);

  return (
    <div className={darkMode ? 'theme-dark' : 'theme-light'}>
      <Navbar
        onCreatePost={() => setShowModal(true)}
        onToggleSidebar={() => setShowSidebar(!showSidebar)}
      />
      <main className="content">
        <LeftBar showSidebar={showSidebar} />
        <div className="profile">
          <div className="images">
            <img

              src={data.cover ? `/uploads/${data.cover}` : "https://images.pexels.com/photos/13440765/pexels-photo-13440765.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2"}
              alt=""
              className="cover"
            />
            <img
              src={data.image ? `/uploads/${data.image}` : "/uploads/default.png"}
              alt="profile picture"
              className="profilePic"
            />
          </div>
          <div className="profileContainer">
            <div className="uInfo">
              <div className="left">
                <a href="http://facebook.com">
                  <FacebookTwoToneIcon fontSize="large" />
                </a>
                <a href="http://instagram.com">
                  <InstagramIcon fontSize="large" />
                </a>
                <a href="http://x.com">
                  <TwitterIcon fontSize="large" />
                </a>
                <a href="http://linkedin.com">
                  <LinkedInIcon fontSize="large" />
                </a>
              </div>
              <div className="center">
                <span>{data.nickname}</span>
                <div className="info">

                  <div className="item">
                    <LanguageIcon />
                    <span>{data.about}</span>
                  </div>
                </div>
                {Profile && Profile.id !== data.id && (
                  <button onClick={followUser}>follow</button>
                )}

              </div>
              <div className="right" >
                {Profile && Profile.id !== data.id ? (
                  <EmailOutlinedIcon />
                ) : (<MoreVertIcon onClick={handleShowPrivacy} />)}


              </div>
              {showPrivacy && (<ProfileCardEditor showPrivacy={showPrivacy} />)}
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