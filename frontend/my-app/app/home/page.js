// Home.js
"use client";
import { useRouter } from "next/navigation";
import { useEffect, useState, useRef } from "react";
import Navbar from '../../components/Navbar.js';
import LeftBar from '../../components/LeftBar.js';
import RightBar from '../../components/RightBar.js';
import { useDarkMode } from '../../context/darkMod';
import Stories from '../../components/stories.js';
import Comment from '../../components/coment.js';
import Post from '../../components/Post.js';
import { middleware } from "../../middleware/middelware.js";
import { useWS } from "../../context/wsContext.js";

export default function Home() {
  // State management
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
  const router = useRouter();
  const { darkMode } = useDarkMode();
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
  function scrollToPost(postId) {

    const el = postRefs.current[postId];
    if (el) {
      el.scrollIntoView({ behavior: "smooth", block: "start" });
    }
  }
  function handleUserSelect(userId) {
    setSelectedUsers((prevSelected) =>
      prevSelected.includes(userId)
        ? prevSelected.filter((id) => id !== userId)
        : [...prevSelected, userId]
    );
  }
  useEffect(() => {
    if (!modalRefhome.current) return;

    const modal = modalRefhome.current;

    const reachedBottom =
      modal.scrollHeight > modal.clientHeight + 10 &&
      modal.scrollTop + modal.clientHeight >= modal.scrollHeight - 50;

    async function handlescrollhome() {
      let b = await fetchingposts();
      if (b) {
        scrollToPost(b)
      }

    }
    if (reachedBottom && !loading) {
      handlescrollhome()
    }
  }, [scroollhome]);
  async function logout(e) {
    e.preventDefault();

    try {
      const res = await fetch("http://localhost:8080/api/logout", {
        method: "POST",
        credentials: "include",
      });

      if (!res.ok) {
        console.error("Logout failed");
        return;
      }

      // Instead of router.push("/home")
      window.location.href = "/login";
    } catch (err) {
      console.error("Logout error:", err);
    }
  }
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
      // const response = await res.json();
      if (res.ok) {

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
      const res = await fetch(`http://localhost:8080/api/Getallpost/${offsetpsot.current}`, {
        method: "GET",
        credentials: "include",
      });

      if (!res.ok) {
        return false
      }

      const data = await res.json();

      if (data.length !== 0) {
        offsetpsot.current += 10
      } else {
        return false
      }

      setPosts([...posts, ...data]);
      return data[0].id
    } catch (err) {
      console.error("Error fetching posts:", err);
      return false
    } finally {
      setLoading(false);
    }
  }
  async function fetchInitialPosts() {
    try {
      setLoading(true);
      const res = await fetch(`http://localhost:8080/api/Getallpost/${0}`, {
        method: "GET",
        credentials: "include",
      });

      if (!res.ok) {
        return false
      }

      const data = await res.json();



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

      if (!response.ok) {
        const errorText = await response.text();
        console.error("Create post error:", errorText);
        throw new Error('Failed to create post');
      }

      const res = await response.json();


      // Fetch the newly created post
      if (res.post_id) {
        const newPost = await fetchPosts(res.post_id);
        if (newPost) {
          setPosts(prevPosts => [newPost, ...prevPosts]);
        }
      }
      offsetpsot.current++

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
        author: post.author
      });
      setShowComments(true);
      // Fetch comments
      const res = await fetch(`http://localhost:8080/api/Getcomments/${post.id}/${offsetcomment.current}`, {
        method: "GET",
        credentials: "include",
      });

      if (!res.ok) {
        return false
      }
      const data = await res.json();



      if (data.length == 0) {
        return false
      } else {
        offsetcomment.current += 10

      }


      setComment([...comment, ...data]);
      return data[0].id

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


  return (
    <div className={darkMode ? 'theme-dark' : 'theme-light'}>
      {/* Navbar */}
      <Navbar
        onLogout={logout}
        onCreatePost={() => setShowModal(true)}
        showSidebar={showSidebar}
        onToggleSidebar={() => setShowSidebar(!showSidebar)}
      />

      {/* Main Content */}
      <main className="content">


        <LeftBar showSidebar={showSidebar} />

        {/* Feed Section */}
        <section className="feed"
          onScroll={(e) => setscroolHome(e.target.scrollTop)}
          style={{ height: "100vh", overflowY: "auto" }}

          ref={modalRefhome}
        >
          <Stories />
          {!posts ? (
            <p>No posts available</p>
          ) : (
            posts.map((post) => (
              <Post
                key={post.id}
                post={post}
                onGetComments={GetComments}
                ondolike={Handlelik}
                ref={el => commentRefs.current[post.id] = el}

              />
            ))
          )}
        </section>
        <RightBar />

      </main>

      {/* Create Post Modal */}
      {showModal && (
        <div
          className={`modal-overlay ${showModal ? 'is-open' : ''}`}
          onMouseDown={(e) => {
            if (e.target === e.currentTarget) setShowModal(false);
            setVisibility("public")
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
                setShowModal(false)
                setVisibility('public')
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
                accept="image/*,video/*"

              />
              <textarea
                placeholder="What's on your mind?"
                className="input"
                value={content}
                onChange={(e) => setContent(e.target.value)}
                required
              />
              {/* Visibility Selection */}
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

              {visibility === 'private' && (

                <div className="user-picker">
                  {loadingFollowers ? (
                    <p>Loading followers...</p>
                  ) : error ? (
                    <p className="error">Error: {error}</p>
                  ) : followers.length > 0 ? (
                    followers.map((follower) => (
                      <label key={follower.id} className="user-picker-item">
                        <img
                          src={`/uploads/${follower.image}` || "/default-avatar.png"}
                          alt={follower.nickname}
                          className="image-avatar"
                        />
                        <input
                          type="checkbox"
                          checked={selectedUsers.includes(follower.id)}
                          onChange={() => handleUserSelect(follower.id)}
                        />
                        <span>{follower.nickname}</span>
                      </label>
                    ))
                  ) : (
                    <p>No followers found.</p>
                  )}
                </div>
              )}
              <div className="modal-actions">
                <button type="button" className="btn cancel" onClick={() => {
                  setShowModal(false)
                  setVisibility('public')
                }}>
                  Cancel
                </button>
                <button type="submit" className="btn submit" disabled={loadingFollowers}>
                  Post
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Comments Modal */}
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
          />
        )
      }
    </div >

  );
}