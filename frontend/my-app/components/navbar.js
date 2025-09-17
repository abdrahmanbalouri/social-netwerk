"use client";
import { useRouter } from "next/navigation";

export default function Navbar({ onLogout, onCreatePost, showSidebar, onToggleSidebar }) {
  const router = useRouter();
  
  return (
    <nav className="navbar">
      <div className="Container">
        <button 
          className="sidebar-toggle"
          onClick={onToggleSidebar}
        >
          <i className={`fa-solid ${showSidebar ? 'fa-times' : 'fa-bars'}`}></i>
        </button>
        <h2 className="log">social</h2>
      </div>
      <div className="search-bar">
        <i className="fa-solid fa-magnifying-glass"></i>
        <input type="search" placeholder="Search for someone" />
      </div>
      <div className="create">
        <button className='btn' onClick={onLogout}>Logout</button>
        <button className="btn btn-primary" onClick={onCreatePost}>
          <i className="fa-solid fa-plus"></i>
        </button>
        <div className="profile-picture">
          <a href="/profile">
            <img src="/avatar.png" alt="Profile" />
          </a>
        </div>
      </div>
    </nav>
  );
}