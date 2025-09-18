"use client";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useDarkMode } from '../context/darkMod';


export default function Navbar({ onLogout, onCreatePost, showSidebar, onToggleSidebar }) {
  const router = useRouter();
  const { darkMode, toggle } = useDarkMode();

  return (
    <div className="navbar">
      <div className="left">
        <Link href="/home" style={{ textDecoration: "none" }}>
          <span>Socile-Network</span>
        </Link>
        <i className={`fa-solid ${darkMode ? 'fa-sun' : 'fa-moon'}`} onClick={toggle}></i>
        <i className="fa-solid fa-bars" onClick={onToggleSidebar}></i>
        <div className="search">
          <i className="fa-solid fa-magnifying-glass"></i>
          <input type="text" placeholder="Search...." />
        </div>
      </div>
      <div className="right">
        <i className="fa-solid fa-bell"></i>
        <div className="user" onClick={() => router.push("/profile")}>
          <img src="/avatar.png" alt="Profile" />
          <span>John Doe</span>
        </div>
      </div>
    </div>
  );
}