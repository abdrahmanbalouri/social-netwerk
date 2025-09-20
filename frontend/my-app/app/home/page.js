"use client";
import './Home.css';
import { useRouter } from "next/navigation";
import { useContext, useEffect, useState, useRef } from "react";
import Navbar from '../../components/Navbar.js';
import LeftBar from '../../components/LeftBar.js';
import RightBar from '../../components/RightBar.js';
import { useDarkMode } from '../../context/darkMod';
import Stories from '../../components/stories.js';
import Link from 'next/link.js';

export default function Home() {
  const router = useRouter();
  const { darkMode } = useDarkMode();

  const [showSidebar, setShowSidebar] = useState(true);
  const [showComments, setShowComments] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [posts, setPosts] = useState([]);
  const [users, setusers] = useState([])
  const [title, setTitle] = useState("");
  const [image, setImage] = useState(null);
  const [content, setContent] = useState("");
  const [commentContent, setCommentContent] = useState("");
  const [comment ,  setComment] = useState([])
  const [dataofonepost , setdataofonepost] = useState()
  const modalRef = useRef(null);
  const commentsModalRef = useRef(null);
  // const previousActiveElementRef = useRef(null);


  async function logout(e) {
    console.log(222);
    
    e.preventDefault();
    const res = await fetch("http://localhost:8080/api/logout", {
      method: "POST",
      credentials: "include",
    });

    if (!res.ok) {
      return;
    }
    router.replace("/login");
  }

  function handleImageChange(e) {
    setImage(e.target.files[0]);
  }
  useEffect(() => {
    async function fetchInitialPosts() {
      try {
        const res = await fetch("http://localhost:8080/api/Getallpost", {
          method: "GET",
          credentials: "include",
        });
        if (!res.ok) {
          throw new Error("Failed to fetch posts");
        }
        const data = await res.json();


        setPosts(data);
      } catch (err) {
        console.error(err);
      }
    }

    fetchInitialPosts();
  }, []);
  useEffect(() => {
    async function fetchusers() {
      try {
        const res = await fetch("http://localhost:8080/api/GetUsersHandler", {
          method: "GET",
          credentials: "include",
        });
        console.log(res);

        if (!res.ok) {
          throw new Error("Failed to fetch posts");
        }
        const data = await res.json();
        console.log(data);


        setusers(data);
      } catch (err) {
        console.error(err);
      }
    }

    fetchusers();
  }, []);
  async function handleCreatePost(e) {
    console.log(222222);

    e.preventDefault();

    try {
      const response = await fetch("http://localhost:8080/api/createpost", {
        method: "POST",
        credentials: "include",
        body: (() => {
          const formData = new FormData();
          formData.append("title", title);
          formData.append("image", image);
          formData.append("content", content);
          return formData;
        })(),
      });
      if (!response.ok) {
        console.log(response);

        throw new Error('failed create post ');
      } else {
        let res = await response.json()
        let data = await fetchPosts(res.post_id)
        
      
        if(!posts){          
          setPosts([data])
        }else{
          setPosts([data, ...posts])
         }
      }

    } catch (err) {
      //console.log(err);

    }
    setTitle("");
    setImage(null);
    setContent("");
    setShowModal(false);
  }
  async function Getcommnets(postid) {

    try {
      const res = await fetch(`http://localhost:8080/api/Getcomments/${postid}`, {
        method: "GET",
        credentials: "include",
      });
      if (!res.ok) {
        throw new Error("Failed to fetch posts");
      }

      const data = await res.json();
      console.log(data);
      setShowComments(true)

      if(!comment){
        setComment(data)

      }else{
        setComment([data,...comment])
      }
      return data;

    } catch (err) {



    }

  }
  async function fetchPosts(postID) {
    try {
      const res = await fetch(`http://localhost:8080/api/Getpost/${postID}`, {
        method: "GET",
        credentials: "include",
      });
      if (!res.ok) {
        throw new Error("Failed to fetch posts");
      }
      const data = await res.json();
      return data;

    } catch (err) {
      console.error(err);
    }
  }
  return (
    <div className={darkMode ? 'theme-dark' : 'theme-light'}>
      <Navbar
        onLogout={logout}
        onCreatePost={() => setShowModal(true)}
        showSidebar={showSidebar}
        onToggleSidebar={() => setShowSidebar(!showSidebar)}
      />
      <main className="content">
        <LeftBar showSidebar={showSidebar} />

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
                      <img src={post.profile_picture || '/avatar.png'} alt="user" />
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

                    <div className="item" onClick={() => Getcommnets(post.id)}>
                      <i className="fa-solid fa-comment"></i>
                      12 Comments
                    </div>

                  </div>
                </div>
              </div>
            ))
          )}
        </section>

        <RightBar />
      </main>
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
 {showComments && (
        <div
          className={`modal-overlay ${showComments ? "is-open" : ""}`}
          onMouseDown={(e) => {
            if (e.target === e.currentTarget) setShowComments(false);
          }}
        >
          <div
            ref={commentsModalRef}
            role="dialog"
            aria-modal="true"
            aria-labelledby="create-post-title"
            className="modal-content"
          >
            <button
              className="modal-close"
              aria-label="Close modal"
              onClick={() => setShowComments(false)}
            >
              ✕
            </button>
            <h3 id="create-post-title">Comments</h3>

            <div className="comments-section">
             {/* // {error && <div className="error2">{error}</div>} */}
              <h2 id="popup-post-title" className="text-xl font-bold mb-4">
                {/* {postTitle} */}
              </h2>

              <div id="popup-comments-container" className="comments-container mb-4">
                {/* {comments.length === 0 ? (
                  <p className="text-gray-500">No comments yet.</p>
                ) : (
                  // comments.map((comment) => (
                  //   <div key={comment.id} className="comment p-2 mb-2 border rounded">
                  //     <p className="font-semibold">{comment.author}</p>
                  //     <p>{comment.content}</p>
                  //     <span className="text-sm text-gray-400">{comment.createdAt}</span>
                  //   </div>
                  // ))
                )} */}
              </div>

              {/* <form id="popup-comment-form" onSubmit={CreateComment}> */}
                <div className="form-group">
                  <textarea
                    id="popup-comment-content"
                    className="w-full p-2 border rounded mb-2"
                    placeholder="Write a comment..."
                    // value={commentContent}
                    onChange={(e) => setCommentContent(e.target.value)}
                    required
                  ></textarea>
                </div>
                <button
                  type="submit"
                  className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
                >
                  Post Comment
                </button>
              {/* </form> */}
            </div>
          </div>
        </div>
      )}
    
  

    </div >
  );
}
