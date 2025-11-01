"use client";
import { useEffect, useRef, useState } from 'react';
import Navbar from '../../../components/Navbar.js';
import { useDarkMode } from '../../../context/darkMod.js';
<<<<<<< HEAD
import '../../../styles/profile.css';
=======
import './profile.css';
>>>>>>> postgroup
import LeftBar from '../../../components/LeftBar.js';
import RightBar from '../../../components/RightBar.js';
import { useParams, useRouter } from 'next/navigation.js';
import Post from '../../../components/Post.js';
import Comment from '../../../components/coment.js';
import { useProfile } from '../../../context/profile.js';
import ProfileCardEditor from '../../../components/ProfileCardEditor.js';
import { useWS } from "../../../context/wsContext.js";
import Link from 'next/link';
import { middleware } from '../../../middleware/middelware.js';

export default function Profile() {
  const { Profile } = useProfile();
  const { darkMode } = useDarkMode();
  const params = useParams();
  const router = useRouter();
  const userId = Number(params.id);
  const { sendMessage } = useWS();
  const [showSidebar, setShowSidebar] = useState(true);
  const [showComments, setShowComments] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [posts, setPosts] = useState([]);
  const [title, setTitle] = useState("");
  const [image, setImage] = useState(null);
  const [content, setContent] = useState("");
  const [comment, setComment] = useState([]);
  const [selectedPost, setSelectedPost] = useState(null);
  const [loading, setLoading] = useState(false);
  const [visibility, setVisibility] = useState('public');
  const [selectedUsers, setSelectedUsers] = useState([]);
  const [loadingFollowers, setLoadingFollowers] = useState(false); // Loading state for fetching followers
  const [error, setError] = useState(''); // Error state for fetching
  const [followers, setFollowers] = useState([]); // Followers list
  const [loadingcomment, setLoadingcomment] = useState(false);
  const [scroollhome, setscroolHome] = useState(0)
  const offsetpsot = useRef(0)
  const offsetcomment = useRef(0)
  const modalRef = useRef(null);
  const modalRefhome = useRef(null)
  const boleanofset = useRef(false)
  const postRefs = useRef({});
  const [showPrivacy, setShowPrivacy] = useState(false);
  const [toast, setToast] = useState(null);

  const showToast = (message, type = "error", duration = 3000) => {
    console.log(8889898);

    setToast({ message, type });
    setTimeout(() => {
      setToast(null);
    }, duration);
  };
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

  const [theprofile, setProfile] = useState(null);

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
        console.log("profile daaaataa", json);


        setProfile(json);

      }
    } catch (err) {
      console.error("loadProfile", err);
    }
  }

  useEffect(() => {
    loadProfile();

  }, []);



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


        console.log("folllllow data ", followw);

        setProfile((prevProfile) => ({
          ...prevProfile,
          privacy: followw.privacy,
          followers: followw.followers,
          following: followw.following,
          isFollowing: followw.isFollowed,
          isPending: followw.isPending
        }));
      }
    } catch (error) {
      console.error("Error following user:", error);
    }
  }




  function PrFollow() {
    if (!theprofile) return "";

    if (theprofile.privacy === "private") {
      if (theprofile.isFollowing) {
        return "Unfollow";
      } else if (theprofile.isPending) {
        return "Pending";
      } else {
        return "Request";
      }
    } else {
      if (theprofile.isFollowing) {
        return "Unfollow";
      } else {
        return "Follow";
      }
    }
  }

  function handleShowPrivacy() {
    setShowPrivacy(!showPrivacy);
  }

  function scrollToPost(postId) {

    const el = postRefs.current[postId];
    if (el) {
      el.scrollIntoView({ behavior: "smooth", block: "start" });
    }
  }
  useEffect(() => {
    async function midle() {
      try {
        const response = await fetch("http://localhost:8080/api/me", {
          credentials: "include",
          method: "GET",
        }, {});
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
  function handleUserSelect(userId) {
    setSelectedUsers((prevSelected) =>
      prevSelected.includes(userId)
        ? prevSelected.filter((id) => id !== userId)
        : [...prevSelected, userId]
    );
  }
  useEffect(() => {  
    console.log(scroollhome);
    
    console.log(loading);
    
    const reachedBottom =
      window.innerHeight + window.scrollY >= document.body.scrollHeight - 20;


    async function handlescrollhome() {
      await fetchingposts();
    
    }

    if (reachedBottom && !loading && posts.length >= 10) {
      
      handlescrollhome();
    }
  }, [scroollhome]);

  const fetchFollowers = async () => {
    setLoadingFollowers(true);
    setError('');
    try {
      const response = await fetch('http://localhost:8080/api/users/followers', {
        method: 'GET',
        credentials: 'include',
      });
      if (!response.ok) {

        throw new Error('Failed to fetch followers');
      }
      let data = await response.json();

      if (!data) {
        data = []
      }
      setFollowers(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoadingFollowers(false);
    }
  };


  function handleImageChange(e) {
    setImage(e.target.files[0]);
  }
  async function Handlelik(postId) {
    try {
      const res = await fetch(`http://localhost:8080/api/like/${postId}`, {
        method: "POST",
        credentials: "include",
      });
      const response = await res.json();
      if (response.error) {
        if (response.error == "Unauthorized") {
          router.push("/login");
          //sendMessage({ type: "logout" })
          return
        } else {
          showToast(response.error)
          return
        }

      }

      const newpost = await fetchPosts(postId)
      for (let i = 0; i < posts.length; i++) {
        if (posts[i].id == newpost.id) {


          setPosts([
            ...posts.slice(0, i),
            newpost,
            ...posts.slice(i + 1)
          ]);
          break
        }
      }

    } catch (err) {
      console.error("Error liking post:", err);
    }
  }
  async function fetchingposts() {
    if (!boleanofset.current) {
      offsetpsot.current += 10
      boleanofset.current = true
    }
    try {
      setLoading(true);
      const res = await fetch(`http://localhost:8080/api/getmypost/${params.id}/${offsetpsot.current}`, {
        method: "GET",
        credentials: "include",
      });
           console.log(res);
           
      if (!res.ok) {
        throw new Error(`Failed to fetch posts: ${res.status}`); 
      }
    

      const data = await res.json();
      


      if (!data) {
        return} 
      console.log(data,"/*/**");
      
      
      offsetpsot.current += 10

      setPosts([...posts, ...data]);
    } catch (err) {
      console.error("Error fetching posts:", err);
    } finally {
      setLoading(false);
    }
  }
  async function fetchInitialPosts() {
    try {
      setLoading(true);
      const res = await fetch(`http://localhost:8080/api/getmypost/${params.id}/${offsetpsot.current}`, {
        method: "GET",
        credentials: "include",
      });

      if (!res.ok) {
        return false
      }

      const data = await res.json();
      if (!data) return

      setPosts([...data]);
      return true
    } catch (err) {
      console.error("Error fetching posts:", err);
      return false
    } finally {
      setLoading(false);
    }
  }
  useEffect(() => {

    fetchInitialPosts();
  }, []);


  // Handle post creation
  async function handleCreatePost(e) {
    e.preventDefault();
    try {
      setLoading(true);
      const formData = new FormData();
      formData.append("title", title);
      if (image) formData.append("image", image);
      formData.append("content", content);
      formData.append("visibility", visibility);


      if (visibility === 'private') {
        formData.append("allowed_users", JSON.stringify(selectedUsers.join(',')));
      }

      const response = await fetch("http://localhost:8080/api/createpost", {
        method: "POST",
        credentials: "include",
        body: formData,
      });

       const res = await response.json();
      if (res.error) {
        if (res.error == "Unauthorized"){
        router.push("/login");
        //sendMessage({ type: "logout" })
        return 
        }else{
          showToast(res.error)
          return
        }

      }


      // Fetch the newly created post
      if (res.post_id) {
        if (Profile && Profile.id === theprofile.id) {

          const newPost = await fetchPosts(res.post_id);
          if (newPost) {
            setPosts(prevPosts => [newPost, ...prevPosts]);
          }
          offsetpsot.current++
        }
      }

    } catch (err) {
      console.error("Error creating post:", err);
      alert("Failed to create post. Please try again.");
    } finally {
      setLoading(false);
      // Reset form
      setTitle("");
      setImage(null);
      setContent("");
      setShowModal(false);
    }
  }

  // Fetch single post by ID
  async function fetchPosts(postID) {
    try {
      const res = await fetch(`http://localhost:8080/api/Getpost/${postID}`, {
        method: "GET",
        credentials: "include",
      });

      if (!res.ok) {
        throw new Error("Failed to fetch post");
      }

      const data = await res.json();
      return data;
    } catch (err) {
      console.error("Error fetching post:", err);
      return null;
    }
  }

  async function GetComments(post) {
    setLoadingcomment(true)

    try {
      setSelectedPost({
        id: post.id,
        title: post.title || post.post_title || "Post",
        image_path: post.image_path,
        content: post.content,
        author: post.first_name + " " + post.last_name
      });

      // Fetch comments
      const res = await fetch(`http://localhost:8080/api/Getcomments/${post.id}/${offsetcomment.current}`, {
        method: "GET",
        credentials: "include",
      });

      if (!res.ok) {
        return false
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
        author: comment.first_name + " " + comment.last_name || "Anonymous",
        content: comment.content || comment.text || "",
        created_at: comment.created_at || comment.createdAt || new Date().toISOString()
      }));
      setShowComments(true);

      if (comments.length == 0) {
        return false
      } else {
        offsetcomment.current += 10

      }


      setComment([...comment, ...comments]);
      return comments[0].id

    } catch (err) {
      return false
    }
    finally {
      setLoadingcomment(false);
    }
  }

  useEffect(() => {


  }, [])

  // Refresh comments after posting a new comment
  async function refreshComments(commentID) {
    if (!selectedPost?.id) return;

    try {
      const res = await fetch(`http://localhost:8080/api/getlastcomment/${commentID}`, {
        method: "GET",
        credentials: "include",
      });

      if (res.ok) {
        const data = await res.json();

        let newcomment = [];

        if (Array.isArray(data)) {
          newcomment = data;
        } else if (data && data.newcomment && Array.isArray(data.newcomment)) {
          newcomment = data.newcomment;
        } else if (data) {
          newcomment = [data];
        }


        setComment([...newcomment, ...comment]);
        offsetcomment.current++


        const potsreplace = await fetchPosts(selectedPost.id)
        for (let i = 0; i < posts.length; i++) {
          if (posts[i].id == selectedPost.id) {
            setPosts([
              ...posts.slice(0, i),
              potsreplace,
              ...posts.slice(i + 1)
            ]);
            break
          }
        }
      }
    } catch (err) {
      console.error("Error refreshing comments:", err);
    }
  }

  // Close comments modal and reset state
  function closeComments() {
    offsetcomment.current = 0
    setShowComments(false);
    setSelectedPost(null);
    setComment([]);
  }
  useEffect(() => {
    const handleScroll = () => {
      setscroolHome(window.scrollY);
    };

    window.addEventListener("scroll", handleScroll);

    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  // Loading state
  if (loading && posts.length === 0) {
    return (
      <div className={`loading-container ${darkMode ? 'theme-dark' : 'theme-light'}`}>
        <div className="loading-spinner"></div>
        <p>Loading posts...</p>
      </div>
    );
  }
  const handleVisibilityChange = (e) => {
    const newVisibility = e.target.value;
    setVisibility(newVisibility);
    if (visibility === 'private' && newVisibility !== 'private') {
      setSelectedUsers([]);
      return;
    }
    if (newVisibility === 'private') {
      fetchFollowers();
    }
  };

  if (!theprofile) {
    return (
      <div className={darkMode ? "theme-dark" : "theme-light"}>
        <Navbar
          onToggleSidebar={() => setShowSidebar(!showSidebar)}
        />
        <main className="content">
          <LeftBar showSidebar={showSidebar} />
          <div className="profile">
            <div className="images">
              <div
                className="cover"
              />
              <div
                className="profilePic"
              />
            </div>
            <div className="profileContainer">
              <div className="uInfo">
                <div className="center">
                  <span>Loading...</span>
                </div>
              </div>
            </div>
            <div className="posts">
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
      {/* Navbar */}
      <Navbar
        onToggleSidebar={() => setShowSidebar(!showSidebar)}
      />
      {toast && (
        <div className={`toast ${toast.type}`}>
          <span>{toast.message}</span>
          <button onClick={() => setToast(null)} className="toast-close">×</button>
        </div>
      )}

      <main className="content">
        <LeftBar showSidebar={showSidebar} />

        <div
          className="main-section"
         
        >
          {/* ===== Profile Section ===== */}
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
                  theprofile?.image
                    ? `/uploads/${theprofile.image}`
                    : "/assets/default.png"
                }
                alt="profile picture"
                className="profilePic"
              />
            </div>

            <div className="profileContainer">
              <div className="uInfo">
                <div className="center">
                  <div className="comb"><h1 className="nickname">{theprofile.first_name + " " + theprofile.last_name}</h1><h1 className="privacy">{theprofile.privacy}</h1></div>
                  {Profile && Profile.id !== theprofile.id && !theprofile.isFollowing && theprofile.privacy === "private" ? (
                    <div className="left">
                      <p className='disabled'>
                        following
                        <strong id="following">{theprofile.following} </strong>
                      </p>
                      <p className='disabled'>
                        followers
                        <strong id="followers"> {theprofile.followers}</strong>
                      </p>
                    </div>

                  ) : (
                    <div className="left">
                      <Link href={`/follow/${theprofile.id}?tab=following`}>

                        <p>
                          following <strong id="following">{theprofile.following}</strong>
                        </p>
                      </Link>
                      <Link href={`/follow/${theprofile.id}?tab=followers`}>
                        <p>
                          followers <strong id="followers">{theprofile.followers}</strong>
                        </p>
                      </Link>
                    </div>
                  )}
                  <hr />
                  <div className="info">
                    <div className="item">
                      <h3>About :</h3>
                      <span>{theprofile.about}</span>
                    </div>
                  </div>

                  {Profile && Profile.id !== theprofile.id && (
                    <button
                      id="FollowBtn"
                      onClick={() => {
                        if (theprofile.privacy === "private" && !theprofile.isFollowed) {
                          sendMessage({
                            type: "followRequest", receiverId: params.id,
                            messageContent: "",
                          });
                        } else if (theprofile.privacy === "public" && !theprofile.isFollowed) {
                          sendMessage({
                            type: "follow", receiverId: params.id,
                            messageContent: "",
                          });
                        }
                        followUser();
                      }}
                      style={{
                        backgroundColor:
                          !theprofile.isFollowing && !theprofile.isPending ? "blue" : "white",
                        color:
                          !theprofile.isFollowing && !theprofile.isPending ? "white" : "black",
                        border: "1px solid #ccc",
                        padding: "8px 16px",
                      }}
                    >
                      {PrFollow()}
                    </button>
                  )}
                </div>

                <div className="right">
                  {Profile && Profile.id !== theprofile.id ? (
                    ""
                  ) : (
                    <i className="fa-solid fa-pen" onClick={handleShowPrivacy}></i>
                  )}
                </div>

                {showPrivacy && (
                  <div className="privacy-overlay">
                    <div onClick={handleShowPrivacy} className="privacy-backdrop"></div>
                    <ProfileCardEditor
                      handleShowPrivacy={handleShowPrivacy}
                      initialCover={theprofile.cover}
                      initialAvatar={theprofile.image}
                      initialAbout={theprofile.about}
                      initialPrivacy={theprofile.privacy}
                    />
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* ===== Posts Section ===== */}
          <section
            className="feed"

          >
            {!posts ? (
              <p>No posts available</p>
            ) : (
              posts.map((post) => (
                <Post
                  key={post.id}
                  post={post}
                  onGetComments={GetComments}
                  ondolike={Handlelik}
                  ref={(el) => (postRefs.current[post.id] = el)}
                />
              ))
            )}
          </section>
        </div>

        <RightBar />
      </main >

      {/* ===== Create Post Modal ===== */}
      {
        showModal && (
          <div
            className={`modal-overlay ${showModal ? "is-open" : ""}`}
            onMouseDown={(e) => {
              if (e.target === e.currentTarget) setShowModal(false);
              setVisibility("public");
            }}
          >
            <div
              ref={modalRef}
              role="dialog"
              aria-modal="true"
              aria-labelledby="create-post-title"
              className="modal-content"
              onMouseDown={(e) => e.stopPropagation()}
            >
              <button
                className="modal-close"
                aria-label="Close modal"
                onClick={() => {
                  setShowModal(false);
                  setVisibility("public");
                }}
              >
                ✕
              </button>
              <h3 id="create-post-title">Create a Post</h3>
              <form onSubmit={handleCreatePost}>
                <input
                  type="text"
                  placeholder="Title"
                  className="input"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  required
                />
                <input
                  type="file"
                  className="input"
                  onChange={handleImageChange}
                  accept="image/*"
                />
                <textarea
                  placeholder="What's on your mind?"
                  className="input"
                  value={content}
                  onChange={(e) => setContent(e.target.value)}
                  required
                />

                <div className="visibility-select">
                  <label htmlFor="visibility">Visibility</label>
                  <select
                    id="visibility"
                    value={visibility}
                    onChange={handleVisibilityChange}
                  >
                    <option value="public">Public (All users)</option>
                    <option value="almost_private">Almost Private (Followers only)</option>
                    <option value="private">Private (Selected followers)</option>
                  </select>
                </div>

                {visibility === "private" && (
                  <div className="user-picker">
                    {loadingFollowers ? (
                      <p>Loading followers...</p>
                    ) : error ? (
                      <p className="error">Error: {error}</p>
                    ) : followers.length > 0 ? (
                      followers.map((follower) => (
                        <label key={follower.id} className="user-picker-item">
                          <img
                            src={follower?.image ? `/uploads/${follower.image}` : "/assets/default.png"}
                            alt={follower.nickname}
                            className="image-avatar"
                          />
                          <span>{follower.first_name + " " + follower.last_name}</span>
                          <input
                            type="checkbox"
                            checked={selectedUsers.includes(follower.id)}
                            onChange={() => handleUserSelect(follower.id)}
                          />
                        </label>
                      ))
                    ) : (
                      <p>No followers found.</p>
                    )}
                  </div>
                )}

                <div className="modal-actions">
                  <button
                    type="button"
                    className="btn cancel"
                    onClick={() => {
                      setShowModal(false);
                      setVisibility("public");
                    }}
                  >
                    Cancel
                  </button>
                  <button type="submit" className="btn submit" disabled={loadingFollowers}>
                    Post
                  </button>
                </div>
              </form>
            </div>
          </div>
        )
      }

      {/* ===== Comments Modal ===== */}
      {
        showComments && (
          <Comment
            comments={comment}
            isOpen={showComments}
            onClose={closeComments}
            postId={selectedPost?.id}
            postTitle={selectedPost?.title}
            onCommentChange={refreshComments}
            lodinggg={loadingcomment}
            ongetcomment={GetComments}
            post={selectedPost}
            showToast={showToast}
          />
        )
      }
    </div >
  );

}
