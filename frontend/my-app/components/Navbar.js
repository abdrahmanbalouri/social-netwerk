"use client";
import { useRouter } from "next/navigation";
import Link from "next/link";


export default function Navbar({ onLogout, onCreatePost, showSidebar, onToggleSidebar }) {
  const router = useRouter();

  return (
    <div className="navbar">
      <div className="left">
        <Link href="/home" style={{ textDecoration: "none" }}>
          <span>Socile-Network</span>
        </Link>
        <i className="fa-solid fa-house"></i>
        <i className="fa-solid fa-moon"></i>
         <i className="fa-solid fa-bars" onClick={onToggleSidebar}></i>
        <div className="search">
          <i className="fa-solid fa-magnifying-glass"></i>
          <input type="text" placeholder="Search...." />
        </div>
      </div>
      <div className="right">
        <i className="fa-solid fa-user" onClick={onLogout}></i>
        <i className="fa-solid fa-comment" onClick={() => router.push("/chat")}></i>
        <i className="fa-solid fa-bell"></i>
        <div className="user" onClick={onCreatePost}>
          <img src="/avatar.png" alt="Profile" />
          <span>John Doe</span>
        </div>
      </div>
    </div>
  );
}