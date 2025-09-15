"use client";
import './Home.css';
import { useRouter } from "next/navigation";



export default function home() {
  const router = useRouter();


  async function logout(e) {
    e.preventDefault();
    const res = await fetch("http://localhost:8080/api/logout", {
      method: "POST",
      credentials: "include",
    });
    console.log(res);
    
    if (!res.ok) {
      return;
    }
    router.replace("/login"); 
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
        <label className="btn btn-primary" htmlFor="create-post">
          <i className="fa-solid fa-plus"></i>
        </label>
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
        <div className="post">
          <div className="post-header">
            <div className="profile-picture">
              <img src="/avatar.png" alt="User" />
            </div>
            <div>
              <span className="text-bold">azraji</span>
              <div className="text-muted" style={{ fontSize: '0.85rem' }}>2 hrs ago</div>
            </div>
          </div>
          <div className="post-content">
            Hello, this is my first post! <span role="img" aria-label="wave">👋</span>
          </div>
          <div className="post-actions">
            <span>Like</span>
            <span>Comment</span>
            <span>Share</span>
          </div>
        </div>
        <div className="post">
          <div className="post-header">
            <div className="profile-picture">
              <img src="/avatar.png" alt="User" />
            </div>
            <div>
              <span className="text-bold">azraji</span>
              <div className="text-muted" style={{ fontSize: '0.85rem' }}>5 hrs ago</div>
            </div>
          </div>
          <div className="post-content">
            Loving the new social network design! <span role="img" aria-label="heart">❤️</span>
          </div>
          <div className="post-actions">
            <span>Like</span>
            <span>Comment</span>
            <span>Share</span>
          </div>
        </div>
      </section>
      <aside className="right-panel">
        <h3>Contacts</h3>
        <ul>
          <li>balouri</li>
          <li>usra</li>
          <li>ahmad</li>
          <li>merwane</li>
          <li>reda</li>
        </ul>
      </aside>
    </main>
  </div>
);

}
