"use client";
import './Home.css';
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
//import { vdd } from '../../../../../uploads/1686064761.jpg';

export default function Home() {
  const router = useRouter();

  const [showModal, setShowModal] = useState(false);
  const [posts, setPosts] = useState([]);
  const [users, setusers] = useState([])
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



        setPosts([data, ...posts])


      }

    } catch (err) {
      //console.log(err);

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

  return (
    <div>
      <nav className="navbar">
        <div className="Container">
          <h2 className="log">social</h2>
        </div>
        <div className="search-bar">
          <i className="fa-solid fa-magnifying-glass"></i>
          <input type="search" placeholder="Search for someone" />
        </div>
        <div className="create">
          <button className='btn' onClick={logout}>Logout</button>
          <button className="btn btn-primary" onClick={() => setShowModal(true)}>
            <i className="fa-solid fa-plus"></i>
          </button>
          <div className="profile-picture">
            <img src="/avatar.png" alt="Profile" />
          </div>
        </div>
      </nav>

      <main className="content">
        <aside className="sidebar">
          <ul>
            <li>Home</li>
            <li>Friends</li>
            <li>Groups</li>
            <li>Marketplace</li>
            <li>Watch</li>
            <li>Memories</li>
          </ul>
        </aside>

        <section className="feed">
          {!posts ? (
            <p>No posts available</p>
          ) : (
            posts.map((post) => (
              <div key={post.id} className="post">
                <div className="post-header">
                  <div className="profile-picture">
                    <img src={post.profile_picture || '/avatar.png'} alt="User" />
                  </div>
                  <div>
                    <span className="text-bold">{post.author}</span>
                    <div className="text-muted" style={{ fontSize: '0.85rem' }}>
                      {new Date(post.created_at).toLocaleString()}
                    </div>
                  </div>
                </div>

                <div className="post-title">{post.title}</div>

                <div className="post-content">{post.content}</div>

                {post.image_path && (
                  <div className="postimage">
                    <img src={post.image_path} alt="Post content" />
                  </div>
                )}

                <div className="post-actions">
                  <span>Like</span>
                  <span>Comment</span>
                  <span>Share</span>
                </div>
              </div>
            ))
          )}
        </section>
        <aside className="right-panel">
          <h3>Contacts</h3>
          <ul>
            {users.map((user) => (
              <li key={user.id}>{user.nickname}</li>
            ))}
          </ul>
        </aside>

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
