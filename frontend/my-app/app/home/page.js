"use client";
import './Home.css';
import { useRouter } from "next/navigation";
import { useEffect, useState, useRef } from "react";
import Navbar from '../../components/Navbar.js';
import LeftBar from '../../components/LeftBar.js';
import RightBar from '../../components/RightBar.js';
import { useDarkMode } from '../../context/darkMod';
import Stories from '../../components/stories.js';
import Link from 'next/link.js';
import Comment from '../../components/coment.js';
import { useProfile } from '../../context/profile.js';

export default function Home() {
  const router = useRouter();
  const { darkMode } = useDarkMode();
  const { profile } = useProfile();

  // State management
  const [showSidebar, setShowSidebar] = useState(true);
  const [showComments, setShowComments] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [posts, setPosts] = useState([]);
  const [users, setusers] = useState([]);
  const [title, setTitle] = useState("");
  const [image, setImage] = useState(null);
  const [content, setContent] = useState("");
  const [comment, setComment] = useState([]);
  const [selectedPost, setSelectedPost] = useState(null);
  const [dataofonepost, setdataofonepost] = useState();
  const [loading, setLoading] = useState(false);
  
  const modalRef = useRef(null);
  const commentsModalRef = useRef(null);
  useEffect(() => {
  console.log(showModal);
  
  }, [showModal])

  // Logout function
  async function logout(e) {
    console.log("Logging out...");
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
      
      router.replace("/login");
    } catch (err) {
      console.error("Logout error:", err);
    }
  }

  // Handle image change for post creation
  function handleImageChange(e) {
    setImage(e.target.files[0]);
  }

  // Fetch initial posts on component mount
  useEffect(() => {
    async function fetchInitialPosts() {
      try {
        setLoading(true);
        const res = await fetch("http://localhost:8080/api/Getallpost", {
          method: "GET",
          credentials: "include",
        });
        
        if (!res.ok) {
          throw new Error("Failed to fetch posts");
        }
        
        const data = await res.json();
        console.log("Posts fetched:", data);
        setPosts(Array.isArray(data) ? data : []);
      } catch (err) {
        console.error("Error fetching posts:", err);
      } finally {
        setLoading(false);
      }
    }

    fetchInitialPosts();
  }, []);

  // Fetch users on component mount
  useEffect(() => {
    async function fetchusers() {
      try {
        const res = await fetch("http://localhost:8080/api/GetUsersHandler", {
          method: "GET",
          credentials: "include",
        });

        if (!res.ok) {
          throw new Error("Failed to fetch users");
        }
        
        const data = await res.json();
        console.log("Users fetched:", data);
        setusers(Array.isArray(data) ? data : []);
      } catch (err) {
        console.error("Error fetching users:", err);
      }
    }

    fetchusers();
  }, []);

  // Handle post creation
  async function handleCreatePost(e) {
    e.preventDefault();
    console.log("Creating post...");

    try {
      setLoading(true);
      const formData = new FormData();
      formData.append("title", title);
      if (image) formData.append("image", image);
      formData.append("content", content);

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
      console.log("Post created:", res);
      
      // Fetch the newly created post
      if (res.post_id) {
        const newPost = await fetchPosts(res.post_id);
        if (newPost) {
          setPosts(prevPosts => [newPost, ...prevPosts]);
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

  // Fetch comments for a specific post - IMPROVED VERSION
  async function Getcommnets(post) {
    console.log("Fetching comments for post:", post.id);
    
    try {
      // Set selected post immediately using post data we already have
      setSelectedPost({ 
        id: post.id, 
        title: post.title || post.post_title || "Post"
      });

      // Fetch comments
      const res = await fetch(`http://localhost:8080/api/Getcomments/${post.id}`, {
        method: "GET",
        credentials: "include",
      });
      
      if (!res.ok) {
        console.error("Comments fetch failed:", res.status);
        throw new Error("Failed to fetch comments");
      }

      const data = await res.json();
      console.log("Comments response:", data);
      
      // Handle different response structures
      let comments = [];
      if (Array.isArray(data)) {
        comments = data;
      } else if (data && typeof data === 'object' && data.comments && Array.isArray(data.comments)) {
        comments = data.comments;
      } else if (data && typeof data === 'object') {
        // Single comment object
        comments = [data];
      }
      
      // Ensure each comment has required properties
      comments = comments.map(comment => ({
        id: comment.id || Math.random(),
        author: comment.author || comment.username || "Anonymous",
        content: comment.content || comment.text || "",
        created_at: comment.created_at || comment.createdAt || new Date().toISOString()
      }));
      
      setComment(comments);
      setShowComments(true);
      
    } catch (err) {
      console.error("Error fetching comments:", err);
      // Set empty comments array on error
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
      console.error("Error refreshing comments:", err);
    }
  }
  

  // Close comments modal and reset state
  function closeComments() {
    setShowComments(false);
    setSelectedPost(null);
    setComment([]);
  }

  // Refresh all posts (useful after liking, commenting, etc.)
  async function refreshPosts() {
    try {
      const res = await fetch("http://localhost:8080/api/Getallpost", {
        method: "GET",
        credentials: "include",
      });
      
      if (res.ok) {
        const data = await res.json();
        setPosts(Array.isArray(data) ? data : []);
      }
    } catch (err) {
      console.error("Error refreshing posts:", err);
    }
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
           <section className="feed">
          <Stories />
          {!posts ? (
            <p>No posts available</p>
          ) : (
            posts.map((post) => (

              <div key={post.id} className="post">
                <div className="container">
                  <div className="user">
                    <div className="userInfo">
                      <img src={`/uploads/${post.profile}` || '/avatar.png'} alt="user" />
                      <div className="details">
                        <Link href={`/profile/${post.user_id}`} style={{ textDecoration: "none", color: "inherit" }} >
                          <span className="name">{post.author}</span>
                        </Link>
                        <span className="date"> {new Date(post.created_at).toLocaleString()}</span>
                      </div>
                    </div>
                  </div>
                  <div className="content">
                    <h3><p style={{ color: "#5271ff" }}> {post.title}</p>{post.content}</h3>

                    {post.image_path && (
                      <img src={post.image_path} alt="Post content" />
                    )}
                  </div>
                  <div className="info">
                    <div className="item">
                      <i className="fa-regular fa-heart"></i>
                      12 Likes
                    </div>

                    <div className="item" onClick={() => Getcommnets(post)}>
                      <i className="fa-solid fa-comment"></i>
                      12 Comments
                    </div>

                  </div>
                </div>
              </div>
            ))
          )}
        </section>

        <RightBar users={users} />
      </main>

      {/* Create Post Modal */}
       {showModal && (
        <div className={`modal-overlay ${showModal ? 'is-open' : ''}`} onMouseDown={(e) => { if (e.target === e.currentTarget) setShowModal(false); }}>
          <div ref={modalRef} role="dialog" aria-modal="true" aria-labelledby="create-post-title" className="modal-content" onMouseDown={(e) => e.stopPropagation()}>
            <button className="modal-close" aria-label="Close modal" onClick={() => setShowModal(false)}>✕</button>
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
              <div className="modal-actions">
                <button type="button" className="btn cancel" onClick={() => setShowModal(false)}>Cancel</button>
                <button type="submit" className="btn submit">Post</button>
              </div>
            </form>
          </div>
        </div>
      )}

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