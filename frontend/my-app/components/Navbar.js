"use client";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useDarkMode } from '../context/darkMod';
import { useProfile } from '../context/profile';


export default function Navbar({ onCreatePost, onToggleSidebar }) {
  const router = useRouter();
  const { darkMode, toggle } = useDarkMode();
  const { profile } = useProfile();
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
        <i className="fa-solid fa-plus" onClick={onCreatePost}></i>
        <div className="user" onClick={() => router.push("/profile")}>
          <img src={profile?.image ? `/uploads/${profile.image}` : "/avatar.png"} alt="user avatar" />
        </div>
      </div>
    </div>
  );
}