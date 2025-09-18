"use client";
import './Home.css';
import { useRouter } from "next/navigation";
import { useContext, useEffect, useState } from "react";
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
  const [showModal, setShowModal] = useState(false);
  const [posts, setPosts] = useState([]);
  const [title, setTitle] = useState("");
  const [image, setImage] = useState(null);
  const [content, setContent] = useState("");

  async function logout(e) {
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
        console.log(data, '-------------------++++++++++++++++++++++++++++++++++');

        setPosts(data);
      } catch (err) {
        console.error(err);
      }
    }

    fetchInitialPosts();
  }, []);

  async function handleCreatePost(e) {
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
        throw new Error("Failed to create post");
      } else {
        let res = await response.json()

        let data = await fetchPosts(res.post_id)
        console.log(data, '-+565554+6');


        setPosts([data, ...posts])


      }

    } catch (err) {
    }
    setTitle("");
    setImage(null);
    setContent("");
    setShowModal(false);
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
  console.log(darkMode, '-----------------');
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
                    <h3>{post.title}</h3>
                    <p>{post.content}</p>
                    {post.image_path && (
                      <img src={post.image_path} alt="Post content" />
                    )}
                  </div>
                  <div className="info">
                    <div className="item">
                    <i className="fa-regular fa-heart"></i>
                      12 Likes
                    </div>
                    <div className="item">
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

      {/* Modal */}
      {showModal && (
        <div className="modal-overlay">
          <div className="modal-content">
            <h3>Create a Post</h3>
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
    </div>
  );
}
